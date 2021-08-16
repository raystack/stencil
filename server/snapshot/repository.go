package snapshot

import (
	"context"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/odpf/stencil/server/store"
)

const snapshotInsertQuery = `
WITH ss(id) as (
	INSERT INTO snapshots (namespace, name, version)
	VALUES ($1, $2, $3) ON CONFLICT DO NOTHING
	RETURNING snapshots.id
)
SELECT COALESCE(
		(
			select ss.id
			from ss
		),
		(
			select id
			from snapshots
			where namespace = $1
				and name = $2
				and version = $3
		)
	)`

// Repository DB access layer
type Repository struct {
	db *store.DB
}

// ListNames returns names
func (r *Repository) ListNames(ctx context.Context, namespace string) ([]string, error) {
	var names []string
	err := pgxscan.Select(ctx, r.db, &names, `SELECT DISTINCT(name) from snapshots where namespace=$1`, namespace)
	return names, err
}

// ListVersions returns versions
func (r *Repository) ListVersions(ctx context.Context, namespace string, name string) ([]string, error) {
	var names []string
	err := pgxscan.Select(ctx, r.db, &names, `SELECT version from snapshots where namespace=$1 and name=$2`, namespace, name)
	return names, err
}

// Exists checks if mentioned version is present or not
func (r *Repository) Exists(ctx context.Context, snapshot *Snapshot) bool {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT count(id) from snapshots where namespace=$1 and name=$2 and version=$3`,
		snapshot.Namespace, snapshot.Name, snapshot.Version).Scan(&count)
	if err != nil {
		return false
	}
	return count != 0
}

// UpdateLatestVersion returns latest version number
func (r *Repository) UpdateLatestVersion(ctx context.Context, snapshot *Snapshot) error {
	return r.db.BeginFunc(ctx, func(t pgx.Tx) error {
		var previousLatestSnapshotID int64
		err := t.QueryRow(ctx, `SELECT id from snapshots where namespace=$1 and name=$2 and latest=true`, snapshot.Namespace, snapshot.Name).Scan(&previousLatestSnapshotID)
		if err != nil && err != pgx.ErrNoRows {
			return err
		}
		_, err = t.Exec(ctx, `UPDATE snapshots set latest=false where id=$1`, previousLatestSnapshotID)
		if err != nil {
			return err
		}
		_, err = t.Exec(ctx, `UPDATE snapshots set latest=true where id=$1`, snapshot.ID)
		return err
	})
}

// GetSnapshot returns full snapshot data
func (r *Repository) GetSnapshot(ctx context.Context, namespace, name, version string, latest bool) (*Snapshot, error) {
	snapshot := &Snapshot{
		Namespace: namespace,
		Name:      name,
	}
	var query strings.Builder
	var args []interface{}
	query.WriteString(`SELECT id, version, latest from snapshots where namespace=$1 and name=$2`)
	args = append(args, namespace, name)
	if latest {
		query.WriteString(` and latest=true`)
	}
	if version != "" {
		query.WriteString(` and version=$3`)
		args = append(args, version)
	}
	err := r.db.QueryRow(ctx, query.String(), args...).Scan(&snapshot.ID, &snapshot.Version, &snapshot.Latest)
	if err == pgx.ErrNoRows {
		return snapshot, ErrNotFound
	}
	return snapshot, err
}

// Create inserts snapshot data
func (r *Repository) Create(ctx context.Context, snapshot *Snapshot) error {
	return r.db.QueryRow(ctx, snapshotInsertQuery, snapshot.Namespace, snapshot.Name, snapshot.Version).Scan(&snapshot.ID)
}

// Delete deletes snapshot data
func (r *Repository) Delete(ctx context.Context, snapshot *Snapshot) error {
	_, err := r.db.Exec(ctx, `DELETE from snapshots where namespace=$1 and name=$2 and version=$3`, snapshot.Namespace, snapshot.Name, snapshot.Version)
	return err
}

// NewSnapshotRepository create instance repo
func NewSnapshotRepository(db *store.DB) *Repository {
	return &Repository{
		db,
	}
}

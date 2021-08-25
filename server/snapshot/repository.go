package snapshot

import (
	"context"
	"fmt"
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

// List returns list of snapshots
func (r *Repository) List(ctx context.Context, queryFields *Snapshot) ([]*Snapshot, error) {
	var snapshots []*Snapshot
	var query strings.Builder
	var args []interface{}
	var conditions []string
	query.WriteString(`SELECT * from snapshots`)
	if queryFields.Latest {
		conditions = append(conditions, "latest=true")
	}
	if queryFields.Namespace != "" {
		conditions = append(conditions, fmt.Sprintf("namespace=$%d", len(args)+1))
		args = append(args, queryFields.Namespace)
	}
	if queryFields.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name=$%d", len(args)+1))
		args = append(args, queryFields.Name)
	}
	if queryFields.Version != "" {
		conditions = append(conditions, fmt.Sprintf("version=$%d", len(args)+1))
		args = append(args, queryFields.Version)
	}
	if queryFields.ID != 0 {
		conditions = append(conditions, fmt.Sprintf("id=$%d", len(args)+1))
		args = append(args, queryFields.ID)
	}
	if len(conditions) > 0 {
		condition := strings.Join(conditions, " AND ")
		query.WriteString(fmt.Sprintf(` WHERE %s`, condition))
	}

	err := pgxscan.Select(ctx, r.db, &snapshots, query.String(), args...)
	return snapshots, err
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

// GetSnapshotByFields returns full snapshot data
func (r *Repository) GetSnapshotByFields(ctx context.Context, namespace, name, version string, latest bool) (*Snapshot, error) {
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

// GetSnapshotByID get snapshot by ID
func (r *Repository) GetSnapshotByID(ctx context.Context, id int64) (*Snapshot, error) {
	var s Snapshot
	err := r.db.QueryRow(ctx, `SELECT * FROM snapshots where id=$1`, id).Scan(&s.ID, &s.Namespace, &s.Name, &s.Version, &s.Latest)
	if err == pgx.ErrNoRows {
		return &s, ErrNotFound
	}
	return &s, err
}

// Exists checks if snapshot exits in DB or not
func (r *Repository) Exists(ctx context.Context, st *Snapshot) bool {
	l, err := r.List(ctx, st)
	return err == nil && len(l) > 0
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

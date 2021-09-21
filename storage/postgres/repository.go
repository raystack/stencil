package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/odpf/stencil/models"
)

// Repository DB access layer
type Store struct {
	db *DB
}

func (r *Store) Close() {
	r.db.Close()
}

// ListSnapshots returns list of snapshots.
func (r *Store) ListSnapshots(ctx context.Context, queryFields *models.Snapshot) ([]*models.Snapshot, error) {
	var snapshots []*models.Snapshot
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

// UpdateSnapshotLatestVersion returns latest version number
func (r *Store) UpdateSnapshotLatestVersion(ctx context.Context, snapshot *models.Snapshot) error {
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
		if err != nil {
			return err
		}
		snapshot.Latest = true
		return nil
	})
}

// GetSnapshotByFields returns full snapshot data
func (r *Store) GetSnapshotByFields(ctx context.Context, namespace, name, version string, latest *bool) (*models.Snapshot, error) {
	sh := &models.Snapshot{
		Namespace: namespace,
		Name:      name,
	}
	var query strings.Builder
	var args []interface{}
	query.WriteString(`SELECT id, version, latest from snapshots where namespace=$1 and name=$2`)
	args = append(args, namespace, name)
	if latest != nil {
		args = append(args, *latest)
		query.WriteString(fmt.Sprintf(` and latest=$%d`, len(args)))
	}
	if version != "" {
		args = append(args, version)
		query.WriteString(fmt.Sprintf(` and version=$%d`, len(args)))
	}
	err := r.db.QueryRow(ctx, query.String(), args...).Scan(&sh.ID, &sh.Version, &sh.Latest)
	if err == pgx.ErrNoRows {
		return sh, models.ErrNotFound
	}
	return sh, err
}

// GetSnapshotByID get snapshot by ID
func (r *Store) GetSnapshotByID(ctx context.Context, id int64) (*models.Snapshot, error) {
	var s models.Snapshot
	err := r.db.QueryRow(ctx, `SELECT * FROM snapshots where id=$1`, id).Scan(&s.ID, &s.Namespace, &s.Name, &s.Version, &s.Latest)
	if err == pgx.ErrNoRows {
		return &s, models.ErrNotFound
	}
	return &s, err
}

// ExistsSnapshot checks if snapshot exits in DB or not
func (r *Store) ExistsSnapshot(ctx context.Context, st *models.Snapshot) bool {
	l, err := r.ListSnapshots(ctx, st)
	return err == nil && len(l) > 0
}

// CreateSnapshot inserts snapshot data
func (r *Store) CreateSnapshot(ctx context.Context, snapshot *models.Snapshot) error {
	return r.db.QueryRow(ctx, snapshotInsertQuery, snapshot.Namespace, snapshot.Name, snapshot.Version).Scan(&snapshot.ID)
}

// DeleteSnapshot deletes snapshot data
func (r *Store) DeleteSnapshot(ctx context.Context, snapshot *models.Snapshot) error {
	_, err := r.db.Exec(ctx, `DELETE from snapshots where namespace=$1 and name=$2 and version=$3`, snapshot.Namespace, snapshot.Name, snapshot.Version)
	return err
}

// PutSchema inserts Schema information in DB
func (r *Store) PutSchema(ctx context.Context, snapshot *models.Snapshot, dbFiles []*models.ProtobufDBFile) error {
	return r.db.Pool.BeginFunc(ctx, func(t pgx.Tx) error {
		err := t.QueryRow(ctx, snapshotInsertQuery, snapshot.Namespace, snapshot.Name, snapshot.Version).Scan(&snapshot.ID)
		if err != nil {
			return err
		}

		batch := &pgx.Batch{}
		for _, file := range dbFiles {
			batch.Queue(fileInsertQuery, snapshot.ID, file.SearchData, file.Data)
		}
		res := t.SendBatch(ctx, batch)
		for i := 0; i < len(dbFiles); i++ {
			_, err = res.Exec()
			if err != nil {
				return err
			}
		}
		err = res.Close()
		return err
	})
}

// GetSchema Fullycontained descriptorset file given list of fully qualified message names.
// If message names are empty then whole fileDescriptorSet data returned
func (r *Store) GetSchema(ctx context.Context, snapshot *models.Snapshot, names []string) ([][]byte, error) {
	var totalData [][]byte
	var err error
	if len(names) > 0 {
		err = pgxscan.Select(ctx, r.db, &totalData, getDataForSpecificMessages, snapshot.ID, names)
	} else {
		err = pgxscan.Select(ctx, r.db, &totalData, getWholeFDS, snapshot.ID)
	}
	return totalData, err
}

const fileInsertQuery = `
WITH file_insert(id) as (
	INSERT INTO protobuf_files (search_data, data)
	VALUES ($2, $3) ON CONFLICT DO NOTHING
	RETURNING id
),
file(id) as (
	SELECT COALESCE(
			(
				SELECT id
				FROM file_insert
			),
			(
				select id
				from protobuf_files
				where search_data = $2
					and data = $3
			)
		)
)
INSERT INTO snapshots_protobuf_files(snapshot_id, file_id)
SELECT $1,file.id from file`

const getDataForSpecificMessages = `
WITH files as (
	SELECT distinct(jsonb_array_elements_text(pf.search_data->'dependencies'))
	from protobuf_files as pf
		join snapshots_protobuf_files as spf on pf.id = spf.file_id
		join snapshots s on s.id = spf.snapshot_id
	WHERE spf.snapshot_id = $1 AND pf.search_data->'messages' ?| $2
)
	SELECT pf.data
	from protobuf_files as pf
		join snapshots_protobuf_files as spf on pf.id = spf.file_id
		join snapshots s on s.id = spf.snapshot_id
	WHERE spf.snapshot_id = $1 and pf.search_data->>'path' in (select * from files);
`

const getWholeFDS = `
SELECT pf.data
	from protobuf_files as pf
		join snapshots_protobuf_files as spf on pf.id = spf.file_id
		join snapshots s on s.id = spf.snapshot_id
	WHERE spf.snapshot_id = $1
`

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

package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/odpf/stencil/models"
)

// Repository DB access layer
type Store struct {
	db *sql.DB
}

func (r *Store) Close() {
	r.db.Close()
}

// ListSnapshots returns list of snapshots.
func (r *Store) ListSnapshots(ctx context.Context, queryFields *models.Snapshot) ([]*models.Snapshot, error) {
	var query strings.Builder
	var args []interface{}
	var conditions []string
	query.WriteString(`SELECT id, namespace, name, version, latest from snapshots`)
	if queryFields.Latest {
		conditions = append(conditions, "latest=true")
	}
	if queryFields.Version != "" {
		conditions = append(conditions, fmt.Sprintf("version=$%d", len(args)+1))
		args = append(args, queryFields.Version)
	}
	if queryFields.Namespace != "" {
		conditions = append(conditions, fmt.Sprintf("namespace=$%d", len(args)+1))
		args = append(args, queryFields.Namespace)
	}
	if queryFields.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name=$%d", len(args)+1))
		args = append(args, queryFields.Name)
	}
	if queryFields.ID != 0 {
		conditions = append(conditions, fmt.Sprintf("id=$%d", len(args)+1))
		args = append(args, queryFields.ID)
	}
	if len(conditions) > 0 {
		condition := strings.Join(conditions, " AND ")
		query.WriteString(fmt.Sprintf(` WHERE %s`, condition))
	}
	rst, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	var snapshots []*models.Snapshot
	for rst.Next() {
		var s models.Snapshot
		err = rst.Scan(&s.ID, &s.Namespace, &s.Name, &s.Version, &s.Latest)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, &s)
	}
	return snapshots, nil
}

// UpdateSnapshotLatestVersion returns latest version number
func (r *Store) UpdateSnapshotLatestVersion(ctx context.Context, snapshot *models.Snapshot) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	var previousLatestSnapshotID int64
	err = tx.QueryRowContext(ctx, `SELECT id from snapshots where namespace=$1 and name=$2 and latest=true`, snapshot.Namespace, snapshot.Name).Scan(&previousLatestSnapshotID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE snapshots set latest=false where id=$1`, previousLatestSnapshotID)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE snapshots set latest=true where id=$1`, snapshot.ID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err == nil {
		snapshot.Latest = true
	}
	return err
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
	err := r.db.QueryRowContext(ctx, query.String(), args...).Scan(&sh.ID, &sh.Version, &sh.Latest)
	if err == sql.ErrNoRows {
		return sh, models.ErrNotFound
	}
	return sh, err
}

// GetSnapshotByID get snapshot by ID
func (r *Store) GetSnapshotByID(ctx context.Context, id int64) (*models.Snapshot, error) {
	var s models.Snapshot
	err := r.db.QueryRowContext(ctx, `SELECT * FROM snapshots where id=$1`, id).Scan(&s.ID, &s.Namespace, &s.Name, &s.Version, &s.Latest)
	if err == sql.ErrNoRows {
		return &s, models.ErrNotFound
	}
	return &s, err
}

// ExistsSnapshot checks if snapshot exits in DB or not
func (r *Store) ExistsSnapshot(ctx context.Context, st *models.Snapshot) bool {
	l, err := r.ListSnapshots(ctx, st)
	if err != nil {
	}
	return err == nil && len(l) > 0
}

// CreateSnapshot inserts snapshot data
func (r *Store) CreateSnapshot(ctx context.Context, snapshot *models.Snapshot) error {
	return r.db.QueryRowContext(ctx, snapshotInsertQuery, snapshot.Namespace, snapshot.Name, snapshot.Version).Scan(&snapshot.ID)
}

// DeleteSnapshot deletes snapshot data
func (r *Store) DeleteSnapshot(ctx context.Context, snapshot *models.Snapshot) error {
	_, err := r.db.ExecContext(ctx, `DELETE from snapshots where namespace=$1 and name=$2 and version=$3`, snapshot.Namespace, snapshot.Name, snapshot.Version)
	return err
}

// PutSchema inserts schema information in DB
func (r *Store) PutSchema(ctx context.Context, snapshot *models.Snapshot, dbFiles []*models.ProtobufDBFile) error {
	var result sql.Result
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	result, err = tx.ExecContext(ctx, snapshotInsertQuery, snapshot.Namespace, snapshot.Name, snapshot.Version)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		err = tx.QueryRowContext(ctx, `SELECT id FROM snapshots WHERE namespace=$1 AND name=$2 AND version=$3`, snapshot.Namespace, snapshot.Name, snapshot.Version).Scan(&snapshot.ID)
	} else {
		snapshot.ID, err = result.LastInsertId()
	}
	if err != nil {
		return err
	}
	for _, file := range dbFiles {
		var fileID int64
		searchDataJSON, err := json.Marshal(file.SearchData)
		if err != nil {
			return err
		}
		result, err = tx.ExecContext(ctx, protobufFileInsertQuery, searchDataJSON, file.Data)
		if err != nil {
			return err
		}
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			err = tx.QueryRowContext(ctx, `SELECT id FROM protobuf_files WHERE search_data=$1 AND data=$2`, searchDataJSON, file.Data).Scan(&fileID)
		} else {
			fileID, err = result.LastInsertId()
		}
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, protobufFileSnapshotMappingInsertQuery, snapshot.ID, fileID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// GetSchema Fullycontained descriptorset file given list of fully qualified message names.
// If message names are empty then whole fileDescriptorSet data returned
func (r *Store) GetSchema(ctx context.Context, snapshot *models.Snapshot, names []string) ([][]byte, error) {
	var totalData [][]byte
	var rst *sql.Rows
	var err error
	if len(names) > 0 {
		rst, err = r.db.QueryContext(ctx, getDataForSpecificMessages, snapshot.ID, names)
		if err != nil {
			return nil, err
		}
	} else {
		rst, err = r.db.QueryContext(ctx, getWholeFDS, snapshot.ID)
		if err != nil {
			return nil, err
		}
		for rst.Next() {
			var data []byte
			err = rst.Scan(&data)
			if err != nil {
				return nil, err
			}
			totalData = append(totalData, data)
		}
	}
	return totalData, nil
}

const protobufFileInsertQuery = `
INSERT INTO protobuf_files(search_data, data)
	VALUES ($1, $2) ON CONFLICT DO NOTHING
`
const protobufFileSnapshotMappingInsertQuery = `
	INSERT INTO snapshots_protobuf_files(snapshot_id, file_id)
	VALUES($1, $2) ON CONFLICT DO NOTHING
`

const snapshotInsertQuery = `
	INSERT INTO snapshots(namespace, name, version)
    VALUES($1, $2, $3) ON CONFLICT DO NOTHING
`

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

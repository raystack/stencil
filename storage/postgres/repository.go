package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/odpf/stencil/models"
	"github.com/odpf/stencil/search"
	"github.com/odpf/stencil/server/namespace"
	"github.com/odpf/stencil/server/schema"
	"github.com/odpf/stencil/storage"
)

func wrapError(err error, name string) error {
	if err == nil {
		return err
	}
	var pgErr *pgconn.PgError
	if errors.Is(err, pgx.ErrNoRows) {
		return storage.NoRowsErr.WithErr(err, name)
	}
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return storage.ConflictErr.WithErr(err, name)
		}
	}
	return storage.UnknownErr.WithErr(err, name)
}

// Store DB access layer
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

// Search returns matching message and field names from query
func (r *Store) Search(ctx context.Context, req *search.SearchRequest) ([]*search.Result, error) {
	var searchResults []*search.Result
	var err error
	err = pgxscan.Select(ctx, r.db, &searchResults, searchQuery, req.Namespace, req.Name, req.Version, req.Latest, req.Query)
	return searchResults, err
}

func (r *Store) CreateNamespace(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error) {
	newNamespace := namespace.Namespace{}
	err := pgxscan.Get(ctx, r.db, &newNamespace, namespaceInsertQuery, ns.ID, ns.Format, ns.Compatibility, ns.Description)
	return newNamespace, wrapError(err, ns.ID)
}

func (r *Store) UpdateNamespace(ctx context.Context, ns namespace.Namespace) (namespace.Namespace, error) {
	newNamespace := namespace.Namespace{}
	err := pgxscan.Get(ctx, r.db, &newNamespace, namespaceUpdateQuery, ns.ID, ns.Format, ns.Compatibility, ns.Description)
	return newNamespace, wrapError(err, ns.ID)
}

func (r *Store) GetNamespace(ctx context.Context, id string) (namespace.Namespace, error) {
	newNamespace := namespace.Namespace{}
	err := pgxscan.Get(ctx, r.db, &newNamespace, namespaceGetQuery, id)
	return newNamespace, wrapError(err, id)
}

func (r *Store) DeleteNamespace(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, namespaceDeleteQuery, id)
	return wrapError(err, id)
}

func (r *Store) ListNamespaces(ctx context.Context) ([]string, error) {
	var namespaces []string
	err := pgxscan.Select(ctx, r.db, &namespaces, namespaceListQuery)
	return namespaces, wrapError(err, "")
}

func (r *Store) CreateSchema(ctx context.Context, sc *schema.Schema) (*schema.Schema, error) {
	newSchema := &schema.Schema{}
	err := r.db.BeginFunc(ctx, func(t pgx.Tx) error {
		if err := pgxscan.Get(ctx, t, newSchema, schemaInsertQuery, sc.ID, sc.Authority, sc.Format, sc.Description, sc.NamespaceID, sc.Compatibility); err != nil {
			return err
		}
		_, err := t.Exec(ctx, versionInsertQuery, sc.ID, uuid.NewString(), sc.Data)
		return err
	})
	return newSchema, wrapError(err, sc.ID)
}

func (r *Store) ListSchemas(ctx context.Context, namespaceID string) ([]string, error) {
	var schemas []string
	err := pgxscan.Select(ctx, r.db, &schemas, schemaListQuery, namespaceID)
	return schemas, wrapError(err, "")
}

const namespaceListQuery = `
SELECT id from namespaces
`

const namespaceGetQuery = `
SELECT * from namespaces where id=$1
`

const namespaceDeleteQuery = `
DELETE from namespaces where id=$1
`

const namespaceUpdateQuery = `
UPDATE namespaces SET format=$2,compatibility=$3,description=$4,updated_at=now()
WHERE id = $1
RETURNING *
`

const namespaceInsertQuery = `
INSERT INTO namespaces (id, format, compatibility, description, created_at, updated_at)
    VALUES ($1, $2, $3, $4, now(), now())
RETURNING *
`

const schemaInsertQuery = `
INSERT INTO schemas (id, authority, format, description, namespace_id, compatibility, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, now(), now())
ON CONFLICT ON CONSTRAINT schemas_pkey DO UPDATE SET updated_at=now() RETURNING *
`

const versionInsertQuery = `
WITH max_version(value) as (
	SELECT COALESCE((SELECT MAX(vs.version) from versions as vs WHERE vs.schema_id=$1), 0)
)
INSERT INTO versions (version, schema_id, id, data, created_at)
		VALUES ((select max_version.value + 1 from max_version), $1, $2, $3, now())
`

const schemaListQuery = `
SELECT id from schemas where namespace_id=$1
`

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
				where data = $3
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

const searchQuery = `
SELECT
  pf.search_data ->> 'path' AS path,
  pf.search_data ->> 'package' AS package,
  jsonb_path_query_array(pf.search_data -> 'messages', ('$[*] ? (@ like_regex "' || $5 || '" flag "i")')::jsonpath) AS "messages",
  jsonb_path_query_array(pf.search_data -> 'fields', ('$[*] ? (@ like_regex "' || $5 || '" flag "i")')::jsonpath) AS "fields",
  jsonb_agg(jsonb_build_object('id', s.id, 'namespace', s.namespace, 'name', s.name, 'version', s.version, 'latest', s.latest)) AS "snapshots"
FROM
  protobuf_files AS pf
  JOIN snapshots_protobuf_files AS spf ON pf.id = spf.file_id
  JOIN snapshots s ON s.id = spf.snapshot_id
WHERE
  s.namespace = COALESCE(NULLIF ($1, ''), s.namespace)
  AND s.name = COALESCE(NULLIF ($2, ''), s.name)
  AND s.version = COALESCE(NULLIF ($3, ''), s.version)
  AND s.latest = COALESCE(NULLIF ($4, FALSE), s.latest)
  AND (pf.search_data -> 'messages' @? ('$[*] ? (@ like_regex "' || $5 || '" flag "i")')::jsonpath
    OR pf.search_data -> 'fields' @? ('$[*] ? (@ like_regex "' || $5 || '" flag "i")')::jsonpath)
GROUP BY
  pf.id
	`

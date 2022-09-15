package postgres

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/odpf/stencil/core/schema"
	"github.com/pkg/errors"
)

type SchemaRepository struct {
	db *DB
}

func NewSchemaRepository(dbc *DB) *SchemaRepository {
	return &SchemaRepository{
		db: dbc,
	}
}

type searchData struct {
	Types  []string
	Fields []string
}

func (r *SchemaRepository) Create(ctx context.Context, namespace string, schemaName string, metadata *schema.Metadata, versionID string, file *schema.SchemaFile) (int32, error) {
	var version int32
	err := r.db.BeginFunc(ctx, func(t pgx.Tx) error {
		vErr := t.QueryRow(ctx, getSchemaVersionByID, versionID).Scan(&version)
		if vErr == nil {
			return nil
		}
		if !errors.Is(vErr, pgx.ErrNoRows) {
			return vErr
		}
		var schemaID int32
		if err := t.QueryRow(ctx, schemaInsertQuery, schemaName, namespace, metadata.Format, metadata.Compatibility).Scan(&schemaID); err != nil {
			return err
		}
		if err := t.QueryRow(ctx, versionInsertQuery, schemaID, versionID, file.ID,
			&searchData{Types: file.Types, Fields: file.Fields}, file.Data).Scan(&version); err != nil {
			return err
		}
		return nil
	})
	return version, wrapError(err, "create schema failed for %s under%s", schemaName, namespace)
}

func (r *SchemaRepository) Get(ctx context.Context, namespaceId, schemaName string, versionNumber int32) ([]byte, error) {
	var versionID string
	var data []byte
	if err := r.db.QueryRow(ctx, getVersionIDFromSchemaNameQuery, namespaceId, schemaName, versionNumber).Scan(&versionID); err != nil {
		return nil, wrapError(err, "Get schema for %s - %s", namespaceId, schemaName)
	}
	err := r.db.QueryRow(ctx, getSchemaDataByVersionID, versionID).Scan(&data)
	return data, wrapError(err, "Get schema for %s - %s", namespaceId, schemaName)
}

func (r *SchemaRepository) GetLatestVersion(ctx context.Context, namespaceId, schemaName string) (int32, error) {
	var version int32
	if err := r.db.QueryRow(ctx, getLatestVersionIDFromSchemaNameQuery, namespaceId, schemaName).Scan(&version); err != nil {
		return version, wrapError(err, "Latest version for %s - %s", namespaceId, schemaName)
	}
	return version, nil
}

func (r *SchemaRepository) GetMetadata(ctx context.Context, namespace, sc string) (*schema.Metadata, error) {
	var meta schema.Metadata
	err := pgxscan.Get(ctx, r.db, &meta, getSchemaMetaQuery, namespace, sc)
	return &meta, wrapError(err, "meta")
}

func (r *SchemaRepository) UpdateMetadata(ctx context.Context, namespace, sc string, in *schema.Metadata) (*schema.Metadata, error) {
	var meta schema.Metadata
	err := pgxscan.Get(ctx, r.db, &meta, updateSchemaMetaQuery, namespace, sc, in.Compatibility)
	return &meta, wrapError(err, "meta")
}

func (r *SchemaRepository) List(ctx context.Context, namespaceID string) ([]schema.Schema, error) {
	var schemas []schema.Schema
	err := pgxscan.Select(ctx, r.db, &schemas, schemaListQuery, namespaceID)
	return schemas, wrapError(err, "List schemas")
}

func (r *SchemaRepository) Delete(ctx context.Context, ns string, sc string) error {
	_, err := r.db.Exec(ctx, deleteSchemaQuery, ns, sc)
	// Idempotent operation to clean orphaned data.
	r.db.Exec(ctx, deleteOrphanedData)
	return wrapError(err, "delete schema")
}

func (r *SchemaRepository) ListVersions(ctx context.Context, ns string, sc string) ([]int32, error) {
	var versions []int32
	err := pgxscan.Select(ctx, r.db, &versions, listVersionsQuery, ns, sc)
	return versions, wrapError(err, "versions")
}

func (r *SchemaRepository) DeleteVersion(ctx context.Context, ns string, sc string, version int32) error {
	_, err := r.db.Exec(ctx, deleteVersionQuery, ns, sc, version)
	// Idempotent operation to clean orphaned data.
	r.db.Exec(ctx, deleteOrphanedData)
	return wrapError(err, "delete version")
}

const schemaInsertQuery = `
INSERT INTO schemas (name, namespace_id, format, compatibility, created_at, updated_at)
    VALUES ($1, $2, $3, $4, now(), now())
ON CONFLICT ON CONSTRAINT schema_name_namespace_unique_idx DO UPDATE SET updated_at=now() RETURNING id
`

const getSchemaVersionByID = `
SELECT vs.version from versions as vs WHERE vs.id=$1
`

const versionInsertQuery = `
WITH max_version(value) as (
	SELECT COALESCE((SELECT MAX(vs.version) from versions as vs WHERE vs.schema_id=$1), 0)
),
insert_version(value) as (
	INSERT INTO versions (version, schema_id, id, created_at)
	VALUES ((select max_version.value + 1 from max_version), $1, $2, now())
	RETURNING version
),
file_insert as (
	INSERT INTO schema_files (id, search_data, data, created_at, updated_at)
	VALUES ($3, $4, $5, now(), now()) ON CONFLICT DO NOTHING
),
map_insert as (
	INSERT INTO versions_schema_files (version_id, schema_file_id) VALUES ($2, $3)
)
SELECT value from insert_version
`

const getLatestVersionIDFromSchemaNameQuery = `
WITH schema_id(id) as (
	SELECT sc.id as id from schemas as sc
	JOIN
	namespaces as ns on ns.id=sc.namespace_id
	WHERE
	ns.id=$1 AND sc.name=$2
),
max_version(value) as (
	SELECT COALESCE((SELECT MAX(vs.version) from versions as vs WHERE vs.schema_id=(select id from schema_id)), 0)
)
select value from max_version
`

const getVersionIDFromSchemaNameQuery = `
WITH schema_id(id) as (
	SELECT sc.id as id from schemas as sc
	JOIN
	namespaces as ns on ns.id=sc.namespace_id
	WHERE
	ns.id=$1 AND sc.name=$2
)
SELECT id from versions WHERE version=$3 AND schema_id=(select id from schema_id)
`

const getSchemaDataByVersionID = `
SELECT sf.data as data from schema_files as sf
JOIN
versions_schema_files as vsf ON vsf.schema_file_id=sf.id
JOIN
versions as v ON v.id=vsf.version_id
WHERE
v.id=$1
`

const getSchemaMetaQuery = `
SELECT COALESCE(sc.authority, '') as authority,  COALESCE(sc.format, '') as format, COALESCE(sc.compatibility, '') as compatibility from schemas as sc WHERE sc.namespace_id=$1 AND sc.name=$2
`
const updateSchemaMetaQuery = `
UPDATE schemas SET compatibility=$3, updated_at=now() WHERE namespace_id=$1 AND name=$2 RETURNING COALESCE(authority, '') as authority,  COALESCE(format, '') as format, COALESCE(compatibility, '') as compatibility
`

const schemaListQuery = `
SELECT name, format, compatibility, COALESCE(authority, '') as authority from schemas where namespace_id=$1
`
const listVersionsQuery = `
SELECT vs.version from versions as vs
JOIN
schemas as sc ON sc.id=vs.schema_id
WHERE sc.namespace_id=$1 AND sc.name=$2
`

const deleteSchemaQuery = `
DELETE from schemas where namespace_id=$1 AND name=$2
`

const deleteVersionQuery = `
WITH version(id) as (
	SELECT vs.id as id from versions as vs
	JOIN
	schemas as sc ON sc.id=vs.schema_id
	WHERE sc.namespace_id=$1 AND sc.name=$2 AND vs.version=$3
)
DELETE from versions where id=(select id from version)
`

const deleteOrphanedData = `
DELETE from schema_files WHERE id NOT IN (SELECT DISTINCT vsf.schema_file_id from versions_schema_files as vsf)
`

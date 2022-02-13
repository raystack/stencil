package postgres

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
INSERT INTO schemas (name, namespace_id, format, compatibility, created_at, updated_at)
    VALUES ($1, $2, $3, $4, now(), now())
ON CONFLICT ON CONSTRAINT schema_name_namespace_unique_idx DO UPDATE SET updated_at=now() RETURNING id
`

const fileInsertQuery = `
INSERT INTO schema_files (id, search_data, data, created_at, updated_at)
	VALUES ($1, $2, $3, now(), now()) ON CONFLICT DO NOTHING
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
SELECT name from schemas where namespace_id=$1
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

const searchAllQuery = `
SELECT jsonb_path_query_array(sf.search_data -> 'Fields', ('$[*] ? (@ like_regex "[A-Z0-9_\\-\\.]+\\.\\w*' || $4 || '\\w*$" flag "i")')::jsonpath) AS "fields",
       jsonb_path_query_array(sf.search_data -> 'Types', ('$[*] ? (@ like_regex "[A-Z0-9_\\-\\.]+\\.\\w*' || $4 || '\\w*$" flag "i")')::jsonpath)  AS "types",
       ns.id                                    																								   AS "namespace_id",
       s.name                                   																								   AS "schema_id",
       v.version                                																								   AS "version_id"
FROM   schema_files                             																								   AS sf
JOIN   versions_schema_files                    																								   AS vsf
ON     sf.id = vsf.schema_file_id
JOIN   versions AS v
ON     vsf.version_id = v.id
JOIN   schemas AS s
ON     s.id = v.schema_id
JOIN   namespaces AS ns
ON     s.namespace_id = ns.id
WHERE  ns.id = COALESCE(NULLIF ($1, ''), ns.id)
AND    s.name=COALESCE(NULLIF ($2, ''), s.name)
AND    v.version=COALESCE(NULLIF ($3, 0), v.version)
AND    (
              sf.search_data -> 'Fields' @? ('$[*] ? (@ like_regex "\\w+\\.\\w+\\.\\w*' || $4 || '\\w*" flag "i")')::jsonpath
       OR     sf.search_data -> 'Types' @? ('$[*] ? (@ like_regex "\\w+\\.\\w*' || $4 || '\\w*" flag "i")')::jsonpath);
`

const searchLatestQuery = `
WITH latest_version AS(
	SELECT  ns.id          AS "namespace_id",
			s.id           AS "schema_id",
			Max(v.version) AS "version_id"
	FROM     versions      AS v
	JOIN     schemas       AS s
	ON       s.id = v.schema_id
	JOIN     namespaces AS ns
	ON       s.namespace_id = ns.id
	WHERE    ns.id = COALESCE(NULLIF ($1, ''), ns.id)
	AND      s.name = COALESCE(NULLIF ($2, ''), s.name)
	GROUP BY (ns.id, s.id))
SELECT jsonb_path_query_array(sf.search_data -> 'Fields', ('$[*] ? (@ like_regex "[A-Z0-9_\\-\\.]+\\.\\w*' || $3 || '\\w*$" flag "i")')::jsonpath) AS "fields",
       jsonb_path_query_array(sf.search_data -> 'Types', ('$[*] ? (@ like_regex "[A-Z0-9_\\-\\.]+\\.\\w*' || $3 || '\\w*$" flag "i")')::jsonpath)  AS "types",
       lv.namespace_id                                                                                                                             AS "namespace_id",
       s.name                                                                                                                       		       AS "schema_id",
       lv.version_id                                                                                                                      		   AS "version_id"
FROM   schema_files                                                                                                                       		   AS sf
JOIN   versions_schema_files                                                                                                              		   AS vsf
ON     sf.id = vsf.schema_file_id
JOIN   versions AS v
ON     vsf.version_id = v.id
JOIN   latest_version AS lv 
ON     v.schema_id = lv.schema_id
AND    v.version = lv.version_id
JOIN   schemas AS s
ON     s.id = lv.schema_id
WHERE  (
              sf.search_data -> 'Fields' @? ('$[*] ? (@ like_regex "\\w+\\.\\w+\\.\\w*' || $3 || '\\w*" flag "i")')::jsonpath
       OR     sf.search_data -> 'Types' @? ('$[*] ? (@ like_regex "\\w+\\.\\w*' || $3 || '\\w*" flag "i")')::jsonpath);
`

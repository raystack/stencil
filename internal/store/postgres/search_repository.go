package postgres

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/raystack/stencil/core/search"
)

const searchAllQuery = `
SELECT jsonb_path_query_array(sf.search_data -> 'Fields', ('$[*] ? (@ like_regex "' || $4 || '" flag "i")')::jsonpath) AS "fields",
       jsonb_path_query_array(sf.search_data -> 'Types', ('$[*] ? (@ like_regex "' || $4 || '" flag "i")')::jsonpath)  AS "types",
       ns.id                                    																	   AS "namespace_id",
       s.name                                   																	   AS "schema_id",
       v.version                                																	   AS "version_id"
FROM   schema_files                             																	   AS sf
JOIN   versions_schema_files                    																	   AS vsf
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
              sf.search_data -> 'Fields' @? ('$[*] ? (@ like_regex "' || $4 || '" flag "i")')::jsonpath
       OR     sf.search_data -> 'Types' @? ('$[*] ? (@ like_regex "' || $4 || '" flag "i")')::jsonpath);
`

const searchLatestQuery = `
WITH latest_version AS(
	SELECT  ns.id          																														  AS "namespace_id",
			s.id           																														  AS "schema_id",
			Max(v.version) 																														  AS "version_id"
	FROM     versions      																														  AS v
	JOIN     schemas       																														  AS s
	ON       s.id = v.schema_id
	JOIN     namespaces AS ns
	ON       s.namespace_id = ns.id
	WHERE    ns.id = COALESCE(NULLIF ($1, ''), ns.id)
	AND      s.name = COALESCE(NULLIF ($2, ''), s.name)
	GROUP BY (ns.id, s.id))
SELECT jsonb_path_query_array(sf.search_data -> 'Fields', ('$[*] ? (@ like_regex "' || $3 || '" flag "i")')::jsonpath) 						       AS "fields",
       jsonb_path_query_array(sf.search_data -> 'Types', ('$[*] ? (@ like_regex "' || $3 || '" flag "i")')::jsonpath)  						       AS "types",
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
              sf.search_data -> 'Fields' @? ('$[*] ? (@ like_regex "' || $3 || '" flag "i")')::jsonpath
       OR     sf.search_data -> 'Types' @? ('$[*] ? (@ like_regex "' || $3 || '" flag "i")')::jsonpath);
`

type SearchRepository struct {
	db *DB
}

func NewSearchRepository(dbc *DB) *SearchRepository {
	return &SearchRepository{
		db: dbc,
	}
}

func (r *SearchRepository) Search(ctx context.Context, req *search.SearchRequest) ([]*search.SearchHits, error) {
	var searchHits []*search.SearchHits
	err := pgxscan.Select(ctx, r.db, &searchHits, searchAllQuery, req.NamespaceID, req.SchemaID, req.VersionID, req.Query)
	return searchHits, err
}

func (r *SearchRepository) SearchLatest(ctx context.Context, req *search.SearchRequest) ([]*search.SearchHits, error) {
	var searchHits []*search.SearchHits
	err := pgxscan.Select(ctx, r.db, &searchHits, searchLatestQuery, req.NamespaceID, req.SchemaID, req.Query)
	return searchHits, err
}

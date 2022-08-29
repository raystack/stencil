package search

import "context"

type Repository interface {
	Search(context.Context, *SearchRequest) ([]*SearchHits, error)
	SearchLatest(context.Context, *SearchRequest) ([]*SearchHits, error)
}

type SearchRequest struct {
	NamespaceID string
	SchemaID    string
	Query       string
	History     bool
	VersionID   int32
}

type SearchResponse struct {
	Hits []*SearchHits
}

type SearchHits struct {
	Fields      []string
	Types       []string
	Path        string
	NamespaceID string
	SchemaID    string
	VersionID   int32
}

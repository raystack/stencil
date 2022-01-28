package domain

import "context"

type SearchRepository interface {
	Search(context.Context, *SearchRequest) ([]*SearchHits, error)
	SearchLatest(context.Context, *SearchRequest) ([]*SearchHits, error)
}

type SearchService interface {
	Search(context.Context, *SearchRequest) (*SearchResponse, error)
}

type SearchRequest struct {
	NamespaceID string
	SchemaID    string
	VersionID   int32
	Query       string
	All         bool
}

type SearchResponse struct {
	Hits []*SearchHits
}

type SearchHits struct {
	Fields      []string
	Types       []string
	Keys        []string
	NamespaceID string
	SchemaID    string
	VersionID   int32
}

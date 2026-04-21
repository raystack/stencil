// Package search provides types and operations for searching across schemas.
package search

import "context"

// Repository defines the persistence interface for search operations.
type Repository interface {
	Search(context.Context, *SearchRequest) ([]*SearchHits, error)
	SearchLatest(context.Context, *SearchRequest) ([]*SearchHits, error)
}

// SearchRequest represents the parameters for a search query.
type SearchRequest struct {
	NamespaceID string
	SchemaID    string
	Query       string
	History     bool
	VersionID   int32
}

// SearchResponse represents the result of a search query.
type SearchResponse struct {
	Hits []*SearchHits
}

// SearchHits represents a single search result with matched fields and types.
type SearchHits struct {
	Fields      []string
	Types       []string
	Path        string
	NamespaceID string
	SchemaID    string
	VersionID   int32
}

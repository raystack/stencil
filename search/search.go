package search

import "context"

type SearchStore interface {
	Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)
	Index(ctx context.Context, req *IndexRequest) error
}

type SearchRequest struct {
	Namespace string
	Field     string
}

type SearchResponse struct {
	Schemas []*Schema
}

type Schema struct {
	Namespace string
	Version   string
	Message   string
	Name      string
	Latest    bool
}

type IndexRequest struct {
	Namespace string
	Version   string
	Fields    []string
	Message   string
	Name      string
	Latest    bool
}

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
	Schemas []*Schema `json:"schema"`
}

type Schema struct {
	Namespace string `json:"namespace"`
	Version   string `json:"version"`
	Message   string `json:"message"`
	Name      string `json:"name"`
	Latest    bool	 `json:"latest"`
	Package   string `json:"package"`
}

type IndexRequest struct {
	Namespace string
	Version   string
	Fields    []string
	Message   string
	Name      string
	Latest    bool
	Package   string
}

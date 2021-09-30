package search

import (
	"context"

	"github.com/odpf/stencil/models"
)

type SearchStore interface {
	SearchShema(ctx context.Context, req *SearchRequest) (*SearchResponse, error)
}

type SearchRequest struct {
	Query     string `binding:"required"`
	Namespace string
	Version   string
	Latest    bool
	Name      string
}

type SearchResponse struct {
	Results []*Result
}

type File struct {
	Path     string
	Package  string
	Messages []string
	Fields   []string
}

type Result struct {
	models.Snapshot `db:"snapshot"`
	Files           []File `db:"files"`
}

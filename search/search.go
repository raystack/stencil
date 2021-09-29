package search

import (
	"context"

	"github.com/odpf/stencil/models"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
)

type SearchStore interface {
	SearchShema(ctx context.Context, req *SearchRequest) (*SearchResponse, error)
}

type SearchRequest struct {
	Query     string
	Type      stencilv1.Search_Type
	Namespace string
	Version   string
	Latest    bool
	Name      string
}

type SearchResponse struct {
	Results []*Result
}

type Result struct {
	models.Snapshot `db:"snapshot"`
	Filepath        string
	Package         string
	MessageName     string
}

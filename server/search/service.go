package search

import (
	"context"
	"errors"

	"github.com/odpf/stencil/server/domain"
)

var (
	ErrEmptyQueryString = errors.New("query string cannot be empty")
	ErrEmptySchemaID    = errors.New("schema_id cannot be empty")
	ErrEmptyNamespaceID = errors.New("namespace_id cannot be empty")
)

type Service struct {
	Repo domain.SearchRepository
}

func (s *Service) Search(ctx context.Context, req *domain.SearchRequest) (*domain.SearchResponse, error) {
	if req.Query == "" {
		return nil, ErrEmptyQueryString
	}

	if req.SchemaID != "" && req.NamespaceID == "" {
		return nil, ErrEmptyNamespaceID
	}

	var res []*domain.SearchHits
	var err error
	if req.VersionID == 0 && !req.History {
		res, err = s.Repo.SearchLatest(ctx, req)
	} else {
		if req.VersionID > 0 && req.SchemaID == "" {
			return nil, ErrEmptySchemaID
		}
		res, err = s.Repo.Search(ctx, req)
	}

	if err != nil {
		return nil, err
	}
	return &domain.SearchResponse{
		Hits: res,
	}, nil
}

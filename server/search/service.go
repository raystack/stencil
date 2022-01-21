package search

import (
	"context"
	"errors"

	"github.com/odpf/stencil/server/domain"
)

type Service struct {
	Repo domain.SearchRepository
}

func (s *Service) SearchSchemas(ctx context.Context, req *domain.SearchSchemasRequest) (*domain.SearchSchemasResponse, error) {
	if req.Query == "" {
		return nil, errors.New("query string cannot be empty")
	}
	res, err := s.Repo.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	return &domain.SearchSchemasResponse{
		NamespaceID: req.NamespaceID,
		Hits:        res,
	}, nil
}

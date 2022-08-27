package search

import (
	"context"
	"errors"
)

var (
	ErrEmptyQueryString = errors.New("query string cannot be empty")
	ErrEmptySchemaID    = errors.New("schema_id cannot be empty")
	ErrEmptyNamespaceID = errors.New("namespace_id cannot be empty")
)

type Service struct {
	Repo SearchRepository
}

func (s *Service) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	if req.Query == "" {
		return nil, ErrEmptyQueryString
	}

	if req.SchemaID != "" && req.NamespaceID == "" {
		return nil, ErrEmptyNamespaceID
	}

	var res []*SearchHits
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
	return &SearchResponse{
		Hits: res,
	}, nil
}

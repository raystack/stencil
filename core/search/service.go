package search

import (
	"context"
	"errors"
)

// ErrEmptyQueryString indicates a missing query string in the search request.
var (
	ErrEmptyQueryString = errors.New("query string cannot be empty")
	// ErrEmptySchemaID indicates a missing schema ID when one is required.
	ErrEmptySchemaID = errors.New("schema_id cannot be empty")
	// ErrEmptyNamespaceID indicates a missing namespace ID when one is required.
	ErrEmptyNamespaceID = errors.New("namespace_id cannot be empty")
)

// Service provides search operations over schemas.
type Service struct {
	repo Repository
}

// NewService creates a new search service with the given repository.
func NewService(repository Repository) *Service {
	return &Service{
		repo: repository,
	}
}

// Search executes a search query and returns matching results.
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
		res, err = s.repo.SearchLatest(ctx, req)
	} else {
		if req.VersionID > 0 && req.SchemaID == "" {
			return nil, ErrEmptySchemaID
		}
		res, err = s.repo.Search(ctx, req)
	}

	if err != nil {
		return nil, err
	}
	return &SearchResponse{
		Hits: res,
	}, nil
}

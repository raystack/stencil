package search

import (
	"context"
)

type StoreSearch struct {
	Store
}

func (s *StoreSearch) SearchShema(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	if req.Type > 2 {
		return nil, ErrInvalidSearchType
	}

	res, err := s.Store.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	return &SearchResponse{
		Results: res,
	}, nil
}

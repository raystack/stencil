package search

import "context"

type Store interface {
	Search(ctx context.Context, req *SearchRequest) ([]*Result, error)
}

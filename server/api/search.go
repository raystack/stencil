package api

import (
	"context"

	"github.com/odpf/stencil/search"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Search returns matched message/field names given a query. Filters can be applied on snapshot
func (a *API) Search(ctx context.Context, in *stencilv1.SearchRequest) (*stencilv1.SearchResponse, error) {
	searchReq := search.SearchRequest{
		Namespace: in.GetNamespace(),
		Name:      in.GetName(),
		Version:   in.GetVersion(),
		Latest:    in.GetLatest(),
		Query:     in.GetQuery(),
	}

	if err := validate.Struct(searchReq); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := a.SearchService.SearchShema(ctx,
		&searchReq,
	)

	if err != nil {
		return nil, err
	}

	results := make([]*stencilv1.SearchResult, len(res.Results))
	for i, result := range res.Results {
		var snapshots []*stencilv1.Snapshot
		for _, sp := range result.Snapshots {
			snapshots = append(snapshots, fromSnapshotToProto(&sp))
		}
		results[i] = &stencilv1.SearchResult{
			Path:      result.Path,
			Package:   result.Package,
			Messages:  result.Messages,
			Fields:    result.Fields,
			Snapshots: snapshots,
		}
	}

	return &stencilv1.SearchResponse{
		Results: results,
	}, nil
}

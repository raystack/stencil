package api

import (
	"context"
	"log"

	"github.com/odpf/stencil/search"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
)

func (a *API) Search(ctx context.Context, in *stencilv1.SearchRequest) (*stencilv1.SearchResponse, error) {

	res, err := a.SearchService.SearchShema(ctx,
		&search.SearchRequest{
			Namespace: in.GetNamespace(),
			Name:      in.GetName(),
			Version:   in.GetVersion(),
			Latest:    in.GetLatest(),
			Query:     in.GetQuery(),
			Type:      in.GetType(),
		},
	)

	if err != nil {
		return nil, err
	}

	results := make([]*stencilv1.SearchResult, len(res.Results))
	log.Println(res.Results)
	for i, result := range res.Results {
		results[i] = &stencilv1.SearchResult{
			Snapshot: &stencilv1.Snapshot{
				Id:        result.ID,
				Namespace: result.Namespace,
				Name:      result.Name,
				Version:   result.Version,
				Latest:    result.Latest,
			},
			Filepath:    result.Filepath,
			Package:     result.Package,
			MessageName: result.MessageName,
		}
	}

	return &stencilv1.SearchResponse{
		Results: results,
	}, nil
}

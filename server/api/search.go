package api

import (
	"context"

	"github.com/odpf/stencil/search"
	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) Search(ctx context.Context, in *stencilv1.SearchRequest) (*stencilv1.SearchResponse, error) {
	searchReq := search.SearchRequest{
		Namespace: in.GetNamespace(),
		Name:      in.GetName(),
		Version:   in.GetVersion(),
		Latest:    in.GetLatest(),
		Query:     in.GetQuery(),
	}

	err := validate.Struct(searchReq)
	if err != nil {
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
		results[i] = &stencilv1.SearchResult{
			Snapshot: fromSnapshotToProto(&result.Snapshot),
			Files:    fromSearchFilesToProtoFiles(result.Files),
		}
	}

	return &stencilv1.SearchResponse{
		Results: results,
	}, nil
}

func fromSearchFilesToProtoFiles(files []search.File) []*stencilv1.SearchResult_File {
	var searchFiles []*stencilv1.SearchResult_File
	for _, f := range files {
		searchFiles = append(searchFiles, &stencilv1.SearchResult_File{
			Path:     f.Path,
			Package:  f.Package,
			Messages: f.Messages,
			Fields:   f.Fields,
		})
	}
	return searchFiles
}

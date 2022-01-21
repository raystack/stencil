package api

import (
	"context"

	"github.com/odpf/stencil/server/domain"
	stencilv1beta1 "github.com/odpf/stencil/server/odpf/stencil/v1beta1"
)

func (a *API) SearchSchemas(ctx context.Context, in *stencilv1beta1.SearchSchemasRequest) (*stencilv1beta1.SearchSchemasResponse, error) {
	searchReq := &domain.SearchSchemasRequest{
		NamespaceID: in.GetNamespaceId(),
		Query:       in.GetQuery(),
		VersionID:   in.GetVersionId(),
	}
	res, err := a.Search.SearchSchemas(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	hits := make([]*stencilv1beta1.SearchHits, 0)
	for _, hit := range res.Hits {
		hits = append(hits, &stencilv1beta1.SearchHits{
			Schema:    hit.Schema,
			VersionId: hit.VersionID,
			Fields:    hit.Fields,
			Types:     hit.Types,
		})
	}
	return &stencilv1beta1.SearchSchemasResponse{
		NamepsaceId: res.NamespaceID,
		Hits:        hits,
	}, nil
}

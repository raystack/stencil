package api

import (
	"context"
	"fmt"

	"github.com/goto/stencil/core/search"
	stencilv1beta1 "github.com/goto/stencil/proto/v1beta1"
)

func (a *API) Search(ctx context.Context, in *stencilv1beta1.SearchRequest) (*stencilv1beta1.SearchResponse, error) {
	searchReq := &search.SearchRequest{
		NamespaceID: in.GetNamespaceId(),
		Query:       in.GetQuery(),
		SchemaID:    in.GetSchemaId(),
	}

	switch v := in.GetVersion().(type) {
	case *stencilv1beta1.SearchRequest_VersionId:
		searchReq.VersionID = v.VersionId
	case *stencilv1beta1.SearchRequest_History:
		searchReq.History = v.History
	}

	res, err := a.search.Search(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	hits := make([]*stencilv1beta1.SearchHits, 0)
	for _, hit := range res.Hits {
		hits = append(hits, &stencilv1beta1.SearchHits{
			SchemaId:    hit.SchemaID,
			VersionId:   hit.VersionID,
			Fields:      hit.Fields,
			Types:       hit.Types,
			NamespaceId: hit.NamespaceID,
			Path:        fmt.Sprintf("/v1beta1/namespaces/%s/schemas/%s/versions/%d", hit.NamespaceID, hit.SchemaID, hit.VersionID),
		})
	}
	return &stencilv1beta1.SearchResponse{
		Hits: hits,
		Meta: &stencilv1beta1.SearchMeta{
			Total: uint32(len(hits)),
		},
	}, nil
}

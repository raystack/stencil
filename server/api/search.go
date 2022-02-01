package api

import (
	"context"
	"strings"

	"github.com/odpf/stencil/server/domain"
	stencilv1beta1 "github.com/odpf/stencil/server/odpf/stencil/v1beta1"
)

func (a *API) Search(ctx context.Context, in *stencilv1beta1.SearchRequest) (*stencilv1beta1.SearchResponse, error) {
	searchReq := &domain.SearchRequest{
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

	res, err := a.SearchService.Search(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	hits := make([]*stencilv1beta1.SearchHits, 0)
	for _, hit := range res.Hits {
		for _, k := range hit.Keys {
			li := strings.LastIndex(k, ".")
			if strings.HasPrefix(k[:li], "m_") {
				hit.Types = append(hit.Types, k[2:li])
			} else {
				hit.Fields = append(hit.Fields, k[2:li])
			}

		}
		hits = append(hits, &stencilv1beta1.SearchHits{
			SchemaId:    hit.SchemaID,
			VersionId:   hit.VersionID,
			Fields:      hit.Fields,
			Types:       hit.Types,
			NamespaceId: hit.NamespaceID,
		})
	}
	return &stencilv1beta1.SearchResponse{
		Hits: hits,
		Meta: &stencilv1beta1.SearchMeta{
			Total: uint32(len(hits)),
		},
	}, nil
}

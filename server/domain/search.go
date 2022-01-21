package domain

import "context"

type SearchRepository interface {
	Search(context.Context, *SearchSchemasRequest) ([]*SearchHits, error)
	GetLatestVersion(ctx context.Context, namespaceId, schemaName string) (int32, error)
}

type SearchService interface {
	SearchSchemas(context.Context, *SearchSchemasRequest) (*SearchSchemasResponse, error)
}

type SearchSchemasRequest struct {
	NamespaceID string
	VersionID   int32
	Query       string
	Latest      bool
}

type SearchSchemasResponse struct {
	NamespaceID string
	Hits        []*SearchHits
}

type SearchHits struct {
	Fields      []string
	Types       []string
	NamespaceID string
	Schema      string
	VersionID   int32
}

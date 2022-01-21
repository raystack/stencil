package api_test

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/stencil/mocks"
	"github.com/odpf/stencil/server/api"
)

func setup() (*mocks.NamespaceService, *mocks.SchemaService, *mocks.SearchService, *runtime.ServeMux, *api.API) {
	nsService := &mocks.NamespaceService{}
	schemaService := &mocks.SchemaService{}
	searchService := &mocks.SearchService{}
	mux := runtime.NewServeMux()
	v1beta1 := &api.API{
		Namespace: nsService,
		Schema:    schemaService,
		Search:    searchService,
	}
	v1beta1.RegisterSchemaHandlers(mux)
	return nsService, schemaService, searchService, mux, v1beta1
}

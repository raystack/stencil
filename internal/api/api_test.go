package api_test

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/stencil/internal/api"
	"github.com/odpf/stencil/internal/api/mocks"
)

func setup() (*mocks.NamespaceService, *mocks.SchemaService, *mocks.SearchService, *runtime.ServeMux, *api.API) {
	nsService := &mocks.NamespaceService{}
	schemaService := &mocks.SchemaService{}
	searchService := &mocks.SearchService{}
	mux := runtime.NewServeMux()
	v1beta1 := api.NewAPI(nsService, schemaService, searchService)
	v1beta1.RegisterSchemaHandlers(mux, nil)
	return nsService, schemaService, searchService, mux, v1beta1
}

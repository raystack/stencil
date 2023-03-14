package api_test

import (
	"github.com/goto/stencil/internal/api"
	"github.com/goto/stencil/internal/api/mocks"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

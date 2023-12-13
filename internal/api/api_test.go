package api_test

import (
	"github.com/goto/stencil/internal/api"
	"github.com/goto/stencil/internal/api/mocks"
	mocks2 "github.com/goto/stencil/pkg/newrelic/mocks"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func setup() (*mocks.NamespaceService, *mocks.SchemaService, *mocks.SearchService, *runtime.ServeMux, *api.API, *mocks2.NewRelic) {
	nsService := &mocks.NamespaceService{}
	schemaService := &mocks.SchemaService{}
	searchService := &mocks.SearchService{}
	newRelic := &mocks2.NewRelic{}
	mux := runtime.NewServeMux()
	v1beta1 := api.NewAPI(nsService, schemaService, searchService, newRelic)
	v1beta1.RegisterSchemaHandlers(mux, nil)
	return nsService, schemaService, searchService, mux, v1beta1, newRelic
}

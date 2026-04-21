package api_test

import (
	"net/http"

	"github.com/raystack/stencil/internal/api"
	"github.com/raystack/stencil/internal/api/mocks"
)

func setup() (*mocks.NamespaceService, *mocks.SchemaService, *mocks.SearchService, *http.ServeMux, *api.API) {
	nsService := &mocks.NamespaceService{}
	schemaService := &mocks.SchemaService{}
	searchService := &mocks.SearchService{}
	mux := http.NewServeMux()
	v1beta1 := api.NewAPI(nsService, schemaService, searchService)
	v1beta1.RegisterSchemaHandlers(mux)
	return nsService, schemaService, searchService, mux, v1beta1
}

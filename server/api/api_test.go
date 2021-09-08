package api_test

import (
	"net/http"

	"github.com/odpf/stencil/server"
	"github.com/odpf/stencil/server/config"

	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/api/mocks"
)

func setup() (http.Handler, *mocks.StoreService, *mocks.MetadataService, *api.API) {
	mockService := &mocks.StoreService{}
	mockMetadataService := &mocks.MetadataService{}
	v1 := &api.API{
		Store:    mockService,
		Metadata: mockMetadataService,
	}
	router := server.Router(v1, &config.Config{})
	return router, mockService, mockMetadataService, v1
}

package api_test

import (
	"github.com/odpf/stencil/server/config"
	server2 "github.com/odpf/stencil/server/server"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/api/mocks"
)

func setup() (*gin.Engine, *mocks.StoreService, *mocks.MetadataService, *api.API) {
	mockService := &mocks.StoreService{}
	mockMetadataService := &mocks.MetadataService{}
	v1 := &api.API{
		Store:    mockService,
		Metadata: mockMetadataService,
	}
	router := server2.Router(v1, &config.Config{})
	return router, mockService, mockMetadataService, v1
}

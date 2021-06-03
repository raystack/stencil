package server

import (
	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/config"
	"github.com/odpf/stencil/server/service"
	"github.com/odpf/stencil/server/store"
)

// Router returns server router
func Router(api *api.API, config *config.Config) *gin.Engine {
	router := gin.New()
	addMiddleware(router, config)
	registerCustomValidations(router)
	registerRoutes(router, api)
	return router
}

// Start Entry point to start the server
func Start() {
	config := config.LoadConfig()
	store := store.New(config)
	dService := &service.DescriptorService{Store: store}
	api := &api.API{
		Store: dService,
	}
	router := Router(api, config)

	runWithGracefulShutdown(config.Port, router, store.Close)
}

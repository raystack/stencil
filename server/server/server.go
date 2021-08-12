package server

import (
	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/config"
	"github.com/odpf/stencil/server/proto"
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
	gcsStore := store.New(config)
	db := store.NewDBStore(config)
	repo := proto.NewProtoRepository(db)
	ser := proto.NewService(repo)
	dService := &service.DescriptorService{Store: gcsStore, ProtoService: ser}
	api := &api.API{
		Store: dService,
	}
	router := Router(api, config)

	runWithGracefulShutdown(config, router, func() {
		gcsStore.Close()
		db.Pool.Close()
	})
}

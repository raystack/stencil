package server

import (
	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/api"
	"github.com/odpf/stencil/server/config"
	"github.com/odpf/stencil/server/proto"
	"github.com/odpf/stencil/server/snapshot"
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
	db := store.NewDBStore(config)
	stRepo := snapshot.NewSnapshotRepository(db)
	stSvc := snapshot.NewSnapshotService(stRepo)
	protoRepo := proto.NewProtoRepository(db)
	protoService := proto.NewService(protoRepo, stRepo)
	api := &api.API{
		Store:    protoService,
		Metadata: stSvc,
	}
	router := Router(api, config)

	runWithGracefulShutdown(config, router, func() {
		db.Pool.Close()
	})
}

package server

import (
	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/api"
)

func registerRoutes(router *gin.Engine, handlers *api.API) {
	apiV1 := router.Group("/v1/namespaces/:namespace")
	router.NoRoute(api.NoRoute)
	router.GET("/ping", api.Ping)
	apiV1.POST("/descriptors", handlers.Upload)
	apiV1.GET("/descriptors", handlers.ListNames)
	apiV1.GET("/descriptors/:name/versions", handlers.ListVersions)
	apiV1.GET("/descriptors/:name/versions/:version", handlers.Download)
	apiV1.GET("/metadata/:name", handlers.GetVersion)
	apiV1.POST("/metadata", handlers.UpdateLatestVersion)
}

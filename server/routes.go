package server

import (
	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/api"
)

func registerRoutes(router *gin.Engine, handlers *api.API) {
	apiV1 := router.Group("/v1")
	apiV1.Use(orgHeaderCheck())
	router.GET("/ping", api.Ping)
	apiV1.POST("/descriptors", handlers.Upload)
	apiV1.GET("/descriptors", handlers.ListNames)
	apiV1.GET("/descriptors/:name", handlers.ListVersions)
	apiV1.GET("/descriptors/:name/:version", handlers.Download)
}

package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/stencil/server/api"
)

func proxyToGin(e *gin.Engine) func(http.ResponseWriter, *http.Request, map[string]string) {
	return func(rw http.ResponseWriter, r *http.Request, m map[string]string) {
		e.ServeHTTP(rw, r)
	}
}

func registerRoutes(router *gin.Engine, mux *runtime.ServeMux, handlers *api.API) {
	apiV1 := router.Group("/v1/namespaces/:namespace")
	router.GET("/ping", api.Ping)
	apiV1.POST("/descriptors", handlers.HTTPUpload)
	apiV1.GET("/descriptors/:name/versions/:version", handlers.HTTPDownload)
	apiV1.PATCH("/descriptors/:name/versions/:version", handlers.HTTPMerge)
	mux.HandlePath("GET", "/ping", proxyToGin(router))
	mux.HandlePath("GET", "/v1/namespaces/{namespace}/descriptors/{name}/versions/{version}", proxyToGin(router))
	mux.HandlePath("POST", "/v1/namespaces/{namespace}/descriptors", proxyToGin(router))
	mux.HandlePath("PATCH", "/v1/namespaces/{namespace}/descriptors/{name}/versions/{version}", proxyToGin(router))
}

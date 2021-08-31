package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/odpf/stencil/server/config"
	"github.com/odpf/stencil/server/models"
	"google.golang.org/grpc/status"
)

func errorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		ginErr := c.Errors.Last()
		if ginErr == nil {
			return
		}
		if err, ok := status.FromError(ginErr.Err); ok {
			code := runtime.HTTPStatusFromCode(err.Code())
			msg := err.Message()
			if code >= 500 {
				msg = "Internal error"
			}
			c.AbortWithStatusJSON(code, gin.H{"message": msg})
			return
		}
		if err, ok := ginErr.Err.(models.APIError); ok {
			c.AbortWithStatusJSON(err.Code(), gin.H{"message": err.Message()})
			return
		}

		if err, ok := ginErr.Meta.(models.APIError); ok {
			c.AbortWithStatusJSON(err.Code(), gin.H{"message": err.Message()})
			return
		}
	}
}

func getLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] \"%s %s\"\t%d %s %s\n",
			params.TimeStamp.UTC().Format(time.RFC3339), params.Method, params.Path,
			params.StatusCode, params.Latency, params.ErrorMessage)
	})
}

func addMiddleware(router *gin.Engine, config *config.Config) {
	router.Use(getNewRelicMiddleware(config))
	router.Use(gin.Recovery())
	router.Use(getLogger())
	router.Use(errorHandle())
}

package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/models"
)

func errorHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		ginErr := c.Errors.Last()
		if ginErr == nil {
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

func addMiddleware(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(getLogger())
	router.Use(errorHandle())
}

package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/odpf/stencil/server/models"
)

func orgHeaderCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := c.GetHeader("x-scope-orgid")
		if orgID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "x-scope-orgid header should be present"})
			return
		}
		c.Next()
	}
}

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

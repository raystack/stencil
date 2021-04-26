package api

import (
	"github.com/gin-gonic/gin"
)

//NoRoute default response for no route
func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"message": "page not found"})
}

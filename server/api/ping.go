package api

import (
	"github.com/gin-gonic/gin"
)

//Ping handler
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

//NoRoute default response for no route
func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"message": "page not found"})
}

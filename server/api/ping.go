package api

import (
	"github.com/gin-gonic/gin"
)

//Ping handler
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

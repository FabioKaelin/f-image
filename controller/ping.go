package controller

import (
	"github.com/gin-gonic/gin"
)

// Ping is a simple ping/pong endpoint
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

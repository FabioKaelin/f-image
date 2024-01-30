package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fabiokaelin/f-image/config"
	"github.com/fabiokaelin/f-image/controller"
	"github.com/gin-gonic/gin"
)

var server *gin.Engine

func init() {
	startconfig, err := config.LoadConfig(".")
	config.StartConfig = startconfig
	if err != nil {
		fmt.Println("Error", err)
	}

	server = gin.Default()

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("CORSMiddleware")
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin, Cookie, caches, Pragma, Expires")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		// Vary: Origin
		c.Writer.Header().Set("Vary", "Origin")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func init() {
	if _, err := os.Stat("public/images"); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll("public/images", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {

	// server.Use(cors.New(corsConfig))
	server.Use(CORSMiddleware())

	router := server.Group("/api")
	router.POST("/users/:userid", controller.PostProfileImage)
	router.GET("/users/:userid", controller.GetProfileImage)

	router.StaticFS("/images", http.Dir("public/images"))
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Route Not Found"})
	})

	log.Fatal(server.Run(":" + "8002"))
}

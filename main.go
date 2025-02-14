package main

import (
	"github.com/fabiokaelin/f-image/config"
	"github.com/fabiokaelin/f-image/controllers"
	"github.com/fabiokaelin/fcommon/pkg/baseHandler"
	"github.com/fabiokaelin/fcommon/pkg/cors"
	"github.com/fabiokaelin/fcommon/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	ferr := config.Load("")
	if ferr != nil {
		panic(ferr.Error())
	}

	logger.InitLogger()

	gin.SetMode(config.GinMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.GetGinLogger(false))

	router.ForwardedByClientIP = true
	router.HandleMethodNotAllowed = true
	router.Use(cors.CORSMiddleware())

	baseHandler.InitBaseHandler(router, baseHandler.IgnoreChecks{}, "Welcome to the API for images")

	apiGroup := router.Group("/api")

	apiGroup.POST("/users/:userid", controllers.PostProfileImage)
	apiGroup.GET("/users/:userid", controllers.GetProfileImage)

	logger.Log.Info("Server Version " + config.FVersion + " is running on port 8002")
	router.Run("0.0.0.0:8002")
}

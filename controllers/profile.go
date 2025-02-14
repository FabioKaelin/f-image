package controllers

import (
	"errors"
	"net/http"
	"os"

	"github.com/fabiokaelin/fcommon/pkg/logger"
	"github.com/gin-gonic/gin"
)

// ProfileRouter defines the routes for the profile
func ProfileRouter(apiGroup *gin.RouterGroup) {
	profileGroup := apiGroup.Group("/profile")
	{
		profileGroup.GET("/:userid", profileGet)
		profileGroup.POST("/:userid", profilePost)
	}
}

func profileGet(c *gin.Context) {
	userID := c.Param("userid")
	logger.Log.Debug("userID: " + userID)

	if _, err := os.Stat("public/images/profiles/" + userID + ".png"); err == nil {
		c.Writer.Header().Set("Content-Type", "image/png")
		c.Status(http.StatusOK)
		c.File("public/images/profiles/" + userID + ".png")
	} else if errors.Is(err, os.ErrNotExist) {
		logger.Log.Debug("file does not exist, returning default image")

		c.Writer.Header().Set("Content-Type", "image/png")
		c.Status(http.StatusOK)
		c.File("public/default.png")
	} else {
		logger.Log.Error(err.Error())
		logger.Log.Error("error durring checking if file exists, returning default image")

		c.Writer.Header().Set("Content-Type", "image/png")
		c.Status(http.StatusOK)
		c.File("public/default.png")
	}
}

func profilePost(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "success", "message": "profile post"})
}

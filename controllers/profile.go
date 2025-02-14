package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/fabiokaelin/f-image/pkg/coder"
	"github.com/fabiokaelin/f-image/pkg/save"
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

	if _, err := os.Stat("public/dynamic/profiles/" + userID + ".png"); err == nil {
		c.Writer.Header().Set("Content-Type", "image/png")
		c.Status(http.StatusOK)
		c.File("public/dynamic/profiles/" + userID + ".png")
	} else if errors.Is(err, os.ErrNotExist) {
		logger.Log.Debug("file does not exist, returning default image")

		c.Writer.Header().Set("Content-Type", "image/png")
		c.Status(http.StatusOK)
		c.File("public/static/profiles/default.png")
	} else {
		logger.Log.Error(err.Error())
		logger.Log.Error("error durring checking if file exists, returning default image")

		c.Writer.Header().Set("Content-Type", "image/png")
		c.Status(http.StatusOK)
		c.File("public/static/profiles/default.png")
	}
}

func profilePost(c *gin.Context) {
	userID := c.Param("userid")
	logger.Log.Debug("userID: " + userID)

	file, _, err := c.Request.FormFile("image")
	if err != nil {
		logger.Log.Warn("file err: " + err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "file err: " + err.Error()})
		return
	}

	img, ferr := coder.Decode(file)
	if ferr != nil {
		logger.Log.Error("decode err: " + ferr.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "decode err: " + ferr.UserMsg()})
		return
	}

	logger.Log.Debug("X-Size: " + fmt.Sprint(img.Bounds().Max.X))
	logger.Log.Debug("Y-Size: " + fmt.Sprint(img.Bounds().Max.Y))

	pngImg, ferr := coder.ConvertToPng(img)
	if ferr != nil {
		logger.Log.Error("convert to png err: " + ferr.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "convert to png err: " + ferr.UserMsg()})
		return
	}

	ferr = save.ResizeSave(pngImg, "public/dynamic/profiles/"+userID+".png", 200, 200)
	if ferr != nil {
		logger.Log.Error("resize save err: " + ferr.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "resize save err: " + ferr.UserMsg()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "success", "worked": true})
}

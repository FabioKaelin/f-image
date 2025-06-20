package controllers

import (
	"net/http"
	"os"
	"strings"

	"github.com/fabiokaelin/fcommon/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	dynamicGalleryImagePath = "public/dynamic/gallery/"
)

func GalleryRouter(apiGroup *gin.RouterGroup) {
	galleryGroup := apiGroup.Group("/gallery")
	{
		galleryGroup.GET("/:imageid", galleryGetByID)
		galleryGroup.GET("", galleryGetAll)
	}
}

func galleryGetByID(c *gin.Context) {
	imageID := c.Param("imageid")
	// c.Writer.Header().Set("Content-Type", "image/png")
	if strings.HasSuffix(imageID, ".jpg") || strings.HasSuffix(imageID, ".jpeg") {
		c.Header("Content-Type", "image/jpeg")
		c.Writer.Header().Set("Content-Type", "image/jpeg")
		logger.Log.Debug("Image is JPEG")
	} else if strings.HasSuffix(imageID, ".gif") {
		c.Header("Content-Type", "image/gif")
		c.Writer.Header().Set("Content-Type", "image/gif")
		logger.Log.Debug("Image is GIF")
	} else if strings.HasSuffix(imageID, ".webp") {
		c.Header("Content-Type", "image/webp")
		c.Writer.Header().Set("Content-Type", "image/webp")
		logger.Log.Debug("Image is WebP")
	} else if strings.HasSuffix(imageID, ".svg") {
		c.Header("Content-Type", "image/svg+xml")
		c.Writer.Header().Set("Content-Type", "image/svg+xml")
		logger.Log.Debug("Image is SVG")
	} else if strings.HasSuffix(imageID, ".png") {
		c.Header("Content-Type", "image/png")
		c.Writer.Header().Set("Content-Type", "image/png")
		logger.Log.Debug("Image is PNG")
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "unsupported image format"})
		return
	}
	c.Status(200)

	logger.Log.Debug("Content-Type: " + c.GetHeader("Content-Type"))
	c.File(dynamicGalleryImagePath + imageID)

}

func galleryGetAll(c *gin.Context) {
	// make a simple ls in the dynamicGalleryImagePath directory
	files, err := os.ReadDir(dynamicGalleryImagePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to read gallery images"})
		return
	}
	var images []string
	for _, file := range files {
		if !file.IsDir() {
			images = append(images, file.Name())
		}
	}
	c.JSON(200, images)
}

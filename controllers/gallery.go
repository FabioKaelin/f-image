package controllers

import (
	"os"

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
	c.Writer.Header().Set("Content-Type", "image/png")
	c.Status(200)

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

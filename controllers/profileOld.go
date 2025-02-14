package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

func PostProfileImage(ctx *gin.Context) {
	fmt.Println("----------------------------- PostProfileImage ----------------------------- Start")
	userID := ctx.Param("userid")
	fmt.Println("userID", userID)

	file, _, err := ctx.Request.FormFile("image")
	if err != nil {
		fmt.Println("file err", err)
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}

	newFileName := "profileimage-" + userID + ".png"
	fmt.Println("newFileName", newFileName)

	imageFile, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error durring decoding image"})
		return
	}

	fmt.Println("X-Size:", imageFile.Bounds().Max.X)
	fmt.Println("Y-Size:", imageFile.Bounds().Max.Y)

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, imageFile); err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error durring encoding image"})
		return
	}

	pngFile, err := png.Decode(buf)
	if err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error durring decoding image"})
		return
	}

	src := imaging.Fill(pngFile, 200, 200, imaging.Center, imaging.Lanczos)
	// src := imaging.Fill(pngFile, 400, 400, imaging.Center, imaging.Lanczos)
	// src := imaging.Fill(pngFile, 100, 100, imaging.Center, imaging.Lanczos)
	// src := imaging.Resize(pngFile, 1000, 0, imaging.Lanczos)
	err = imaging.Save(src, fmt.Sprintf("public/images/%v", newFileName))
	if err != nil {
		fmt.Println("error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error durring saving image"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"worked": true})
	fmt.Println("----------------------------- PostProfileImage ----------------------------- End")
}

func GetProfileImage(ctx *gin.Context) {
	userID := ctx.Param("userid")
	fmt.Println("userID", userID)
	// check if public/images/profileimage-userID.png exists
	// if not, return default image (public/default.png)
	// if yes, return public/images/profileimage-userID.png

	if _, err := os.Stat("public/images/profileimage-" + userID + ".png"); err == nil {
		fmt.Println("file exists")
		ctx.Status(http.StatusOK)
		ctx.Writer.Header().Set("Content-Type", "image/png")
		ctx.File("public/images/profileimage-" + userID + ".png")

	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Println("file does not exist")
		ctx.Writer.Header().Set("Content-Type", "image/png")
		ctx.Status(http.StatusOK)
		ctx.File("public/default.png")
	} else {
		fmt.Println(err)
		fmt.Println("error durring checking if file exists")
		// ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "error durring checking if file exists"})
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		ctx.Writer.Header().Set("Content-Type", "image/png")
		ctx.Status(http.StatusOK)
		ctx.File("public/default.png")
	}

}

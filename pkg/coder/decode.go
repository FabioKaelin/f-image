package coder

import (
	"bytes"
	"image"
	"image/png"
	"mime/multipart"

	"github.com/fabiokaelin/ferror"
)

func Decode(file multipart.File) (image.Image, ferror.FError) {
	imageFile, _, err := image.Decode(file)
	if err != nil {
		ferr := ferror.FromError(err)
		ferr.SetUserMsg("error durring decoding image")
		return nil, ferr
	}
	return imageFile, nil
}

func ConvertToPng(img image.Image) (image.Image, ferror.FError) {
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		ferr := ferror.FromError(err)
		ferr.SetUserMsg("error durring encoding image")
		return nil, ferr
	}

	pngFile, err := png.Decode(buf)
	if err != nil {
		ferr := ferror.FromError(err)
		ferr.SetUserMsg("error durring decoding image")
		return nil, ferr
	}

	return pngFile, nil
}

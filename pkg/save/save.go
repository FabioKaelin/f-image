package save

import (
	"image"

	"github.com/disintegration/imaging"
	"github.com/fabiokaelin/ferror"
)

func ResizeSave(img image.Image, filename string, width, height int) ferror.FError {
	src := imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
	err := imaging.Save(src, filename)
	if err != nil {
		ferr := ferror.FromError(err)
		ferr.SetUserMsg("Error saving image")
		return ferr
	}
	return nil
}

func Save(img image.Image, filename string) ferror.FError {
	err := imaging.Save(img, filename)
	if err != nil {
		ferr := ferror.FromError(err)
		ferr.SetUserMsg("Error saving image")
		return ferr
	}
	return nil
}

package controllers

import (
	"errors"
	"fmt"
	"hash/fnv"
	"image"
	"image/color"
	"image/png"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fabiokaelin/f-image/pkg/coder"
	"github.com/fabiokaelin/f-image/pkg/save"
	"github.com/fabiokaelin/fcommon/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	// defaultProfileImagePath is the path to the default profile image
	defaultProfileImagePath = "public/static/profiles/default.png"
	// dynamicProfileImagePath is the path to the dynamic profile image
	dynamicProfileImagePath = "public/dynamic/profiles/"
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
	c.Writer.Header().Set("Content-Type", "image/png")
	c.Status(http.StatusOK)
	if userID == "default" {
		c.File(defaultProfileImagePath)
		return
	}

	if _, err := os.Stat(dynamicProfileImagePath + userID + ".png"); err == nil {
		c.File(dynamicProfileImagePath + userID + ".png")
	} else if errors.Is(err, os.ErrNotExist) {
		if genErr := createUserDefaultProfileImage(userID); genErr != nil {
			logger.Log.Error("failed to generate user default image: " + genErr.Error())
			c.File(defaultProfileImagePath)
			return
		}
		c.File(dynamicProfileImagePath + userID + ".png")
	} else {
		logger.Log.Error(err.Error())
		logger.Log.Error("error durring checking if file exists, returning default image")

		c.File(defaultProfileImagePath)
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

	ferr = save.ResizeSave(pngImg, dynamicProfileImagePath+userID+".png", 200, 200)
	if ferr != nil {
		logger.Log.Error("resize save err: " + ferr.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "resize save err: " + ferr.UserMsg()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "success", "worked": true})
}

func createUserDefaultProfileImage(userID string) error {
	if err := os.MkdirAll(dynamicProfileImagePath, 0o755); err != nil {
		return err
	}

	defaultImageFile, err := os.Open(defaultProfileImagePath)
	if err != nil {
		return err
	}
	defer defaultImageFile.Close()

	defaultImage, _, err := image.Decode(defaultImageFile)
	if err != nil {
		return err
	}

	targetHue := deterministicHue(userID)
	tinted := tintAvatarImage(defaultImage, targetHue)

	targetPath := filepath.Join(dynamicProfileImagePath, userID+".png")
	outFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return png.Encode(outFile, tinted)
}

func deterministicHue(seed string) float64 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(seed))
	return float64(h.Sum32()%360) / 360.0
}

func tintAvatarImage(src image.Image, targetHue float64) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := color.RGBAModel.Convert(src.At(x, y)).(color.RGBA)
			if rgba.A == 0 {
				dst.SetRGBA(x, y, rgba)
				continue
			}

			r := float64(rgba.R) / 255.0
			g := float64(rgba.G) / 255.0
			b := float64(rgba.B) / 255.0
			_, s, v := rgbToHSV(r, g, b)

			// Keep white and near-white areas untouched.
			if s < 0.08 && v > 0.9 {
				dst.SetRGBA(x, y, rgba)
				continue
			}

			nr, ng, nb := hsvToRGB(targetHue, s, v)
			dst.SetRGBA(x, y, color.RGBA{
				R: uint8(clamp01(nr) * 255),
				G: uint8(clamp01(ng) * 255),
				B: uint8(clamp01(nb) * 255),
				A: rgba.A,
			})
		}
	}

	return dst
}

func rgbToHSV(r, g, b float64) (float64, float64, float64) {
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min

	h := 0.0
	if delta != 0 {
		switch max {
		case r:
			h = math.Mod((g-b)/delta, 6)
		case g:
			h = ((b-r)/delta + 2)
		case b:
			h = ((r-g)/delta + 4)
		}
		h /= 6
		if h < 0 {
			h += 1
		}
	}

	s := 0.0
	if max != 0 {
		s = delta / max
	}

	v := max
	return h, s, v
}

func hsvToRGB(h, s, v float64) (float64, float64, float64) {
	if s == 0 {
		return v, v, v
	}

	h = math.Mod(h, 1.0) * 6
	i := math.Floor(h)
	f := h - i
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	switch int(i) {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	default:
		return v, p, q
	}
}

func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
)

const quality = 100

func streamImage(c *gin.Context, img image.Image, contentType string) {
	var err error
	c.Header("Content-Type", contentType)

	switch contentType {
	case "image/jpeg":
		err = jpeg.Encode(c.Writer, img, &jpeg.Options{Quality: quality})
	case "image/png":
		err = png.Encode(c.Writer, img)
	case "image/gif":
		err = png.Encode(c.Writer, img)
	default:
		validateError(c, fmt.Errorf("content type: '%s' for stream image not show", contentType))
	}

	if err != nil {
		log.Error(err)
	}
}

func validateError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
}

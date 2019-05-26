package v1

import (
	"github.com/gin-gonic/gin"
	"image/internal/forms"
)

// Разрешенный тип контента
const contentType = "image/jpeg"

func (h *Handler) Resize(c *gin.Context) {
	form := forms.Resize{ContentType: contentType}

	if err := c.BindQuery(&form); err != nil {
		validateError(c, err)
		return
	}

	img, err := form.Image()
	if err != nil {
		validateError(c, err)
		return
	}

	c.Header("Cache-Control", "public, max-age=3600, must-revalidate")
	streamImage(c, img, contentType)
}

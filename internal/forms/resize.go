package forms

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

type Resize struct {
	Url         string `form:"url" validate:"required,url"`
	Width       uint   `form:"width" validate:"required,gt=0,lt=10000"`
	Height      uint   `form:"height" validate:"required,gt=0,lt=10000"`
	ContentType string `form:"-"`
}

func (r *Resize) Image() (img image.Image, err error) {
	rd, contentType, err := streamByUrl(r.Url)
	if err != nil {
		return
	}

	defer func() {
		err := rd.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	if r.ContentType != contentType {
		err = fmt.Errorf("content type: `%s` not allowed for url %s", contentType, r.Url)
		return
	}

	img, _, err = image.Decode(rd)
	img = resize.Resize(r.Width, r.Height, img, resize.Lanczos3)

	return
}

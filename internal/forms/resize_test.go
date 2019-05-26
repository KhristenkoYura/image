package forms

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestResizeHeightWidthValidate тестирования пограничных состояний у размера картинки
func TestResize_HeightWidthValidate(t *testing.T) {
	var err error

	sizeError := []int{0, -1, 10000}
	sizeFields := []string{"Width", "Height"}

	f := Resize{}

	for _, size := range sizeError {
		for _, field := range sizeFields {
			switch field {
			case "Width":
				f.Width = uint(size)
			case "Height":
				f.Height = uint(size)
			}
			err = validator.New().StructPartial(f, field)
			assert.Error(t, err)
		}
	}
}

// TestResize_RequiredValidate проверка  на заполненость всех поле
func TestResize_RequiredValidate(t *testing.T) {
	var err error
	fields := []string{"Width", "Height", "Url"}

	// Пустые значения
	fe := Resize{}
	//Заполненные
	fn := Resize{ //
		Width:  1,
		Height: 1,
		Url:    "http://localhost/",
	}

	for _, field := range fields {
		err = validator.New().StructPartial(fe, field)
		assert.Error(t, err)

		err = validator.New().StructPartial(fn, field)
		assert.NoError(t, err)
	}
}

// TestResize_Image стнадартная проверка на ресайз картинки
func TestResize_Image(t *testing.T) {
	url, cancel := fileServer("1.jpg")
	defer cancel()

	f := Resize{
		Width:       200,
		Height:      100,
		Url:         url,
		ContentType: "image/jpeg",
	}

	img, err := f.Image()
	assert.NoError(t, err)

	assert.Equal(t, f.Width, uint(img.Bounds().Max.X))
	assert.Equal(t, f.Height, uint(img.Bounds().Max.Y))
}

// TestResize_ImageError ошибки с другим типом картинки
func TestResize_ImageError(t *testing.T) {
	url, cancel := fileServer("1.gif")
	defer cancel()

	f := Resize{
		Width:       200,
		Height:      100,
		Url:         url,
		ContentType: "image/jpeg",
	}

	_, err := f.Image()
	assert.Error(t, err)
}

// TestStreamByUrlNotFound - проверка в случаи если нет файла
func TestStreamByUrlNotFound(t *testing.T) {
	url, cancel := fileServer("not_found.gif")
	defer cancel()

	_, _, err := streamByUrl(url)

	assert.Error(t, err)
}

// TestStreamByUrlContentType - проверка на соотвествие типов
func TestStreamByUrlContentType(t *testing.T) {

	files := map[string]string{
		"image/jpeg": "1.jpg",
		"image/gif":  "1.gif",
		"image/png":  "1.png",
	}

	for contentType, file := range files {
		url, cancel := fileServer(file)
		defer cancel()

		rc, ct, err := streamByUrl(url)
		defer rc.Close()

		assert.Equal(t, contentType, ct)
		assert.NoError(t, err)
	}
}

// fileServer тестовый файловый сервер: папка test/
func fileServer(file string) (url string, cancel func()) {
	wd, _ := os.Getwd()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(wd + "/../../test/" + file)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = io.Copy(w, bytes.NewReader(data))
	}))
	url = ts.URL
	cancel = func() {
		ts.Close()
	}
	return
}

package forms

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// количество байт для определения типа контента
const sniffLen = 512

//Проверяем на соотвествие интерфейса
var _ io.ReadCloser = readerClose{}

// streamByUrl загружаем файл с автоматическим опеделением его content-type
func streamByUrl(url string) (r io.ReadCloser, contentType string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad http status: %d for url %s", resp.StatusCode, url)
		return
	}

	stiffType := make([]byte, sniffLen)
	_, err = resp.Body.Read(stiffType)

	if err != nil {
		return
	}

	contentType = http.DetectContentType(stiffType)
	r = readerClose{
		r:     io.MultiReader(bytes.NewReader(stiffType), resp.Body),
		close: resp.Body.Close,
	}
	return
}

// Структура для преобразования (проксирование) io.Reader в io.ReadCloser
type readerClose struct {
	r     io.Reader
	close func() error
}

func (rc readerClose) Read(p []byte) (n int, err error) {
	n, err = rc.r.Read(p)
	return
}

func (rc readerClose) Close() error {
	return rc.close()
}

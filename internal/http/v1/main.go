package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/olebedev/gin-cache"
	"time"
)

type Handler struct {
}

func (h *Handler) AddHandlers(r *gin.RouterGroup) {
	r.GET("resize", h.Resize)
}

func (h *Handler) AddMiddleware(r *gin.RouterGroup) {

	//Кеш для статики
	r.Use(cache.New(cache.Options{
		// set expire duration
		// by default zero, it means that cached content won't drop
		Expire: 1 * time.Hour,
		// it uses slice listed below as default to calculate
		// key, if `Header` slice is not specified
		Headers: []string{
			"User-Agent",
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Cookie",
			"User-Agent",
		},
		Store:         cache.NewInMemory(),
		DoNotUseAbort: false,
	}))

}

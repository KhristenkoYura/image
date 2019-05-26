package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"image/internal/app"
	"image/internal/http/v1"
)

type Web struct {
	r      *gin.Engine
	app    *app.App
	ctx    context.Context
	cancel context.CancelFunc
}

// Http сервер
func Server() *Web {
	ctx, cancel := context.WithCancel(context.Background())

	return &Web{
		r:      gin.New(),
		app:    app.New(),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (w *Web) Run() {
	defer w.cancel()
	w.init()
	w.initRouter()
	w.startServer()
}

func (w *Web) Close() {
	w.cancel()
}

func (w *Web) init() {
	//Update validator
	binding.Validator = &defaultValidator{}
}

func (w *Web) initRouter() {
	apiv1 := w.r.Group("/api/v1")

	v1h := v1.Handler{}
	v1h.AddHandlers(apiv1)
	v1h.AddMiddleware(apiv1)
}

func (w *Web) startServer() {
	err := w.r.Run(w.app.Config.App.AddrBind)
	if err != nil {
		panic(err)
	}
}

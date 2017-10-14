package main

//go:generate go-bindata-assetfs frontend

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(h *Handlers, assets string) *httprouter.Router {
	router := httprouter.New()

	// API
	router.PUT("/api/ai/:uuid", h.New)
	router.POST("/api/ai/:uuid/update", h.Update)
	router.POST("/api/ai/:uuid/update/image", h.UpdateImage)
	//router.POST("/api/ai/:uuid/error", h.Error)
	router.DELETE("/api/ai/:uuid", h.Delete)

	// Client
	router.GET("/stream", h.Stream)
	if len(assets) == 0 {
		router.ServeFiles("/*filepath", assetFS())
	} else {
		router.ServeFiles("/*filepath", http.Dir(assets))
	}

	return router
}

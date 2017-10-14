package main

import (
	"github.com/julienschmidt/httprouter"
)

func NewRouter(h *Handlers) *httprouter.Router {
	router := httprouter.New()

	// API
	router.PUT("/api/ai/:uuid", h.New)
	router.POST("/api/ai/:uuid/update", h.Update)
	router.POST("/api/ai/:uuid/image", h.Image)
	router.DELETE("/api/ai/:uuid", h.Delete)

	// Client
	router.GET("/stream", h.Stream)

	return router
}

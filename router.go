package main

//go:generate go-bindata-assetfs frontend/...

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(h *Handlers, assets string) *httprouter.Router {
	router := httprouter.New()

	// API
	router.GET("/api/ai", h.Index)
	router.PUT("/api/ai/:uuid", h.New)
	router.GET("/api/ai/:uuid", h.Get)
	router.POST("/api/ai/:uuid/update", h.Update)
	router.POST("/api/ai/:uuid/update/image", h.UpdateImage)
	//router.POST("/api/ai/:uuid/error", h.Error)
	router.DELETE("/api/ai/:uuid", h.Delete)

	// Dashboard Stream
	router.GET("/api/stream", h.Stream)

	// Dashboard
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.Redirect(w, r, "/dashboard/", http.StatusFound)
	})
	if assets == "builtin" {
		router.ServeFiles("/dashboard/*filepath", assetFS())
	} else {
		router.ServeFiles("/dashboard/*filepath", http.Dir(assets))
	}

	return router
}

package main

//go:generate go-bindata-assetfs frontend/...

import (
	"log"
	"net/http"

	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(h *Handlers, assets string, debug bool) *httprouter.Router {
	router := httprouter.New()

	// Wapper
	router.GET("/api/ai", h.Index)
	router.PUT("/api/ai/:uuid", h.New)
	router.GET("/api/ai/:uuid", h.Get)
	router.POST("/api/ai/:uuid/update", h.Update)
	router.POST("/api/ai/:uuid/update/image", h.UpdateImage)
	//router.POST("/api/ai/:uuid/error", h.Error)
	router.DELETE("/api/ai/:uuid", h.Delete)

	// Docker
	router.PUT("/api/docker", h.DockerNew)

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

	// Debug
	if debug {
		log.Println("router:", "Debug routes enabled!")
		router.HandlerFunc("GET", "/debug/pprof/", pprof.Index)
		router.GET("/debug/pprof/:profile", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			switch p.ByName("profile") {
			case "cmdline":
				pprof.Cmdline(w, r)
			case "profile":
				pprof.Profile(w, r)
			case "symbol":
				pprof.Symbol(w, r)
			case "trace":
				pprof.Trace(w, r)
			default:
				pprof.Handler(p.ByName("profile")).ServeHTTP(w, r)
			}
		})
	}

	return router
}

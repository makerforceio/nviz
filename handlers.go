package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handlers struct {
	manager *Manager
}

// API

func (h *Handlers) New(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse JSON
	log.Println(ps)

	var d NewData
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(d)

	// Print
	fmt.Fprintln(w, "OK")
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Don't need to parse JSON, just forward directly to clients
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	// Print
	fmt.Fprintln(w, "OK")
}

func (h *Handlers) Image(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Push image

	// Print
	fmt.Fprintln(w, "OK")
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Push image

	// Print
	fmt.Fprintln(w, "OK")
}

// Client

func (h *Handlers) Stream(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	m.streamer.ServeHTTP(w http.ResponseWriter, _ *http.Request)
}

func NewHandlers(manager *Manager) *Handlers {
	return &Handlers{
		manager,
	}
}

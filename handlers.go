package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handler

type Handlers struct {
	manager *Manager
}

// API

func (h *Handlers) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	instances, err := h.manager.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(instances)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("handlers:", "Index")
	// Return
	w.Write(data)
}

func (h *Handlers) New(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	// Create
	var d Instance
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.manager.New(ps.ByName("uuid"), d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("handlers:", "New:", ps.ByName("uuid"), d)
	// Print
	fmt.Fprintln(w, "OK")
}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	instances, err := h.manager.Get(ps.ByName("uuid"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(instances)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("handlers:", "Get:", ps.ByName("uuid"))
	// Return
	w.Write(data)
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	// Update
	var d Update
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.manager.Update(ps.ByName("uuid"), d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("handlers:", "Update:", ps.ByName("uuid"), d)
	// Print
	fmt.Fprintln(w, "OK")
}

func (h *Handlers) UpdateImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	// Push image
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	d := buf.Bytes()

	err := h.manager.UpdateImage(ps.ByName("uuid"), UpdateImage{d})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("handlers:", "Update:", ps.ByName("uuid"), len(d))
	// Print
	fmt.Fprintln(w, "OK")
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := h.manager.Delete(ps.ByName("uuid"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("handlers:", "Delete:", ps.ByName("uuid"))
	// Print
	fmt.Fprintln(w, "OK")
}

// Client

func (h *Handlers) Stream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.manager.Streamer.ServeHTTP(w, r)
}

func NewHandlers(manager *Manager) *Handlers {
	return &Handlers{
		manager,
	}
}

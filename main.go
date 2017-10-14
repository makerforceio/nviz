package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port int
var assets string

func init() {
	flag.IntVar(&port, "port", 8080, "listen on port")
	flag.StringVar(&assets, "assets", "builtin", "serve files from assets instead of built-in assets")
}

func main() {
	// Parse commandline flags
	flag.Parse()

	manager := NewManager()
	// Setup handlers
	handlers := NewHandlers(manager)
	// Create router
	router := NewRouter(handlers, assets)

	// Listen
	log.Println("main: Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

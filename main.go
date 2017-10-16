package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port int
var url string
var assets string

func init() {
	flag.IntVar(&port, "port", 8080, "listen on port")
	flag.StringVar(&url, "url", "http://172.17.0.1:8080", "url that is acessible by containers on $DOCKER_HOST")
	flag.StringVar(&assets, "assets", "builtin", "serve files from assets instead of built-in assets")
}

func main() {
	// Parse commandline flags
	flag.Parse()

	manager := NewManager(url)
	// Setup handlers
	handlers := NewHandlers(manager)
	// Create router
	router := NewRouter(handlers, assets)

	// Listen
	log.Println("main: Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

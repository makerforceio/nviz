package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8080, "listen on port")
}

func main() {
	// Parse commandline flags
	flag.Parse()

	manager := NewManager()
	// Setup handlers
	handlers := NewHandlers(manager)
	// Create router
	router := NewRouter(handlers)

	// Listen
	log.Println("main: Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

package main

import (
	"log"

	"github.com/julienschmidt/sse"
)

type Manager struct {
	Streamer []sse.Streamer
}

func NewManager() *Manager {
	return &Manager{
		Streamer: sse.New()
	}
}

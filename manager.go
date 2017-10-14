package main

import (
	"errors"
	"encoding/base64"

	"github.com/julienschmidt/sse"
)

// Structs

type Instance struct {
	Name            string      `json:"name"`
	Args            interface{} `json:"args"`
	LastUpdate      Update      `json:"lastupdate"`
	LastUpdateImage []byte `json:"-"`
}

type Update struct {
	Epoch uint64      `json:"epoch"`
	Loss  float64     `json:"loss"`
	Args  interface{} `json:"args"`
}

// Manager

type Manager struct {
	Streamer  *sse.Streamer
	instances map[string]*Instance
}

func (m *Manager) New(uuid string, instance Instance) error {
	_, ok := m.instances[uuid]
	if ok {
		return errors.New("stuff exists")
	}
	m.instances[uuid] = &instance
	m.Streamer.SendJSON(uuid, "New", instance)
	return nil
}

func (m *Manager) Update(uuid string, update Update) error {
	_, ok := m.instances[uuid]
	if !ok {
		return errors.New("stuff doesn't exist")
	}
	m.instances[uuid].LastUpdate = update
	m.Streamer.SendJSON(uuid, "Update", update)
	return nil
}

func (m *Manager) UpdateImage(uuid string, image []byte) error {
	_, ok := m.instances[uuid]
	if !ok {
		return errors.New("stuff doesn't exist")
	}
	m.instances[uuid].LastUpdateImage = image
	image64 := make([]byte, base64.StdEncoding.EncodedLen(len(image)))
	base64.StdEncoding.Encode(image64, image)
	m.Streamer.SendBytes(uuid, "UpdateImage", image64)
	return nil
}

func (m *Manager) Delete(uuid string) error {
	_, ok := m.instances[uuid]
	if !ok {
		return errors.New("stuff doesn't exist")
	}
	delete(m.instances, uuid)
	m.Streamer.SendJSON(uuid, "Delete", nil)
	return nil
}

func NewManager() *Manager {
	return &Manager{
		Streamer:  sse.New(),
		instances: make(map[string]*Instance),
	}
}

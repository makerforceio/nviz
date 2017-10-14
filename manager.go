package main

import (
	"errors"

	"github.com/julienschmidt/sse"
)

// Structs

type Instance struct {
	Name            string      `json:"name"`
	Args            interface{} `json:"args"`
	LastUpdate      Update      `json:"lastupdate"`
	LastUpdateImage UpdateImage `json:"lastupdateimage"`
}

type Update struct {
	Epoch uint64      `json:"epoch"`
	Loss  float64     `json:"loss"`
	Stats interface{} `json:"stats"`
}

type UpdateImage struct {
	Image []byte `json:"image"`
}

type Event struct {
	UUID string      `json:"uuid"`
	Data interface{} `json:"data"`
}

// Manager

type Manager struct {
	Streamer  *sse.Streamer
	instances map[string]*Instance
}

func (m *Manager) Index() (map[string]*Instance, error) {
	return m.instances, nil
}

func (m *Manager) New(uuid string, instance Instance) error {
	_, ok := m.instances[uuid]
	if ok {
		return errors.New("stuff exists")
	}
	m.instances[uuid] = &instance
	m.Streamer.SendJSON(uuid, "New", Event{uuid, instance})
	return nil
}

func (m *Manager) Get(uuid string) (*Instance, error) {
	instance, ok := m.instances[uuid]
	if !ok {
		return nil, errors.New("stuff doesn't exist")
	}
	return instance, nil
}

func (m *Manager) Update(uuid string, update Update) error {
	_, ok := m.instances[uuid]
	if !ok {
		return errors.New("stuff doesn't exist")
	}
	m.instances[uuid].LastUpdate = update
	m.Streamer.SendJSON(uuid, "Update", Event{uuid, update})
	return nil
}

func (m *Manager) UpdateImage(uuid string, updateImage UpdateImage) error {
	_, ok := m.instances[uuid]
	if !ok {
		return errors.New("stuff doesn't exist")
	}
	m.instances[uuid].LastUpdateImage = updateImage
	m.Streamer.SendJSON(uuid, "UpdateImage", Event{uuid, updateImage})
	return nil
}

func (m *Manager) Delete(uuid string) error {
	_, ok := m.instances[uuid]
	if !ok {
		return errors.New("stuff doesn't exist")
	}
	delete(m.instances, uuid)
	m.Streamer.SendJSON(uuid, "Delete", Event{uuid, nil})
	return nil
}

func NewManager() *Manager {
	return &Manager{
		Streamer:  sse.New(),
		instances: make(map[string]*Instance),
	}
}

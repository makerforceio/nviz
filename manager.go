package main

import (
	"errors"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/julienschmidt/sse"
	"golang.org/x/net/context"
)

// Structs

type Instance struct {
	Name            string      `json:"name"`
	Args            interface{} `json:"args"`
	LastUpdate      Update      `json:"lastupdate"`
	LastUpdateImage UpdateImage `json:"lastupdateimage"`
}

type Update struct {
	Epoch        uint64      `json:"epoch"`
	TrainingLoss float64     `json:"training_loss"`
	Stats        interface{} `json:"stats"`
}

type UpdateImage struct {
	Image []byte `json:"image"`
}

type Event struct {
	UUID string      `json:"uuid"`
	Data interface{} `json:"data"`
}

type DockerNewParams struct {
	Image string `json:"docker-image"`
	Args  string `json:"docker-args"`
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

// Docker
func (m *Manager) DockerNew(params DockerNewParams) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return errors.New("something broke")
	}

	_, err = cli.ImagePull(ctx, params.Image, types.ImagePullOptions{})
	if err != nil {
		return errors.New("container image likely not found")
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: params.Image,
		Cmd:   strings.Split(params.Args, " "),
		Env:   []string{"URL=http://localhost:8080"},
	}, nil, nil, "")
	if err != nil {
		return errors.New("cannot create container image")
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return errors.New("cannot start container image")
	}

	return nil
}

func NewManager() *Manager {
	return &Manager{
		Streamer:  sse.New(),
		instances: make(map[string]*Instance),
	}
}

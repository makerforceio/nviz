package main

import (
	"errors"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/julienschmidt/sse"
	"golang.org/x/net/context"
)

// Structs

type Instance struct {
	sync.Mutex
	Name             string                 `json:"name"`
	Args             interface{}            `json:"args"`
	LastUpdate       Update                 `json:"lastupdate"`
	LastUpdateImages map[uint64]UpdateImage `json:"lastupdateimages"`
}

type Update struct {
	Epoch        uint64      `json:"epoch"`
	TrainingLoss float64     `json:"training_loss"`
	Stats        interface{} `json:"stats"`
}

type UpdateImage struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"`
	Url  []byte `json:"url"`
}

type Event struct {
	UUID string      `json:"uuid"`
	Data interface{} `json:"data"`
}

type DockerNewParams struct {
	ShouldPull bool   `json:"docker-shouldpull"`
	Image      string `json:"docker-image"`
	Args       string `json:"docker-args"`
}

// Manager

type Manager struct {
	sync.Mutex
	Streamer  *sse.Streamer
	instances map[string]*Instance
	url       string
}

func (m *Manager) Index() (map[string]*Instance, error) {
	return m.instances, nil
}

func (m *Manager) New(uuid string, instance Instance) error {
	_, ok := m.instances[uuid]
	if ok {
		return errors.New("stuff exists")
	}
	m.Lock()
	m.instances[uuid] = &instance
	m.instances[uuid].LastUpdateImages = make(map[uint64]UpdateImage)
	m.Unlock()
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
	instance, ok := m.instances[uuid]
	if !ok {
		return errors.New("stuff doesn't exist")
	}
	instance.Lock()
	instance.LastUpdate = update
	instance.Unlock()
	m.Streamer.SendJSON(uuid, "Update", Event{uuid, update})
	return nil
}

func (m *Manager) UpdateImage(uuid string, updateImage UpdateImage) error {
	instance, ok := m.instances[uuid]
	if !ok {
		return errors.New("stuff doesn't exist")
	}
	instance.Lock()
	instance.LastUpdateImages[updateImage.ID] = updateImage
	instance.Unlock()
	m.Streamer.SendJSON(uuid, "UpdateImage", Event{uuid, updateImage})
	return nil
}

func (m *Manager) Delete(uuid string) error {
	_, ok := m.instances[uuid]
	if !ok {
		return errors.New("stuff doesn't exist")
	}
	m.Lock()
	delete(m.instances, uuid)
	m.Unlock()
	m.Streamer.SendJSON(uuid, "Delete", Event{uuid, nil})
	return nil
}

// Docker
func (m *Manager) DockerNew(params DockerNewParams) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	if params.ShouldPull {
		_, err := cli.ImagePull(ctx, params.Image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: params.Image,
		Cmd:   strings.Split(params.Args, " "),
		Env:   []string{"URL=" + m.url},
	}, &container.HostConfig{
		NetworkMode: "host",
	}, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func (m *Manager) DockerDelete(uuid string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, c := range containers {
		if c.Labels["uuid"] == uuid {
			if err := cli.ContainerStop(ctx, c.ID, nil); err != nil {
				return err
			}
		}
	}

	return nil
}

func NewManager(url string) *Manager {
	return &Manager{
		Streamer:  sse.New(),
		instances: make(map[string]*Instance),
		url:       url,
	}
}

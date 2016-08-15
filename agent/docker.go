package main

import (
	"github.com/fsouza/go-dockerclient"
	"log"
)

type DockerManager struct {
	client   *docker.Client
	logger   *log.Logger
	listener chan *docker.APIEvents
}

func NewDockerManager(host string, logger *log.Logger) (*DockerManager, error) {
	client, err := docker.NewClient(host)
	if err != nil {
		return nil, err
	}
	return &DockerManager{
		client:   client,
		logger:   logger,
		listener: make(chan *docker.APIEvents, 10),
	}, nil
}

func (d *DockerManager) RunningContainers() []string {
	ids := []string{}
	containers, err := d.client.ListContainers(docker.ListContainersOptions{All: false})
	if err != nil {
		d.logger.Printf("Docker error: %s", err)
		return ids
	}

	for _, container := range containers {
		ids = append(ids, container.ID)
	}

	return ids

}

func (d *DockerManager) start(rmanager *RabbitMQManager) error {
	go d.manageEvents(rmanager)
	err := d.client.AddEventListener(d.listener)
	return err
}

func (d *DockerManager) stop() error {
	if d.listener == nil {
		return nil
	}
	err := d.client.RemoveEventListener(d.listener)
	return err
}

func (d *DockerManager) manageEvents(rmanager *RabbitMQManager) {
	for {
		event := <-d.listener
		if event.Type == "container" && (event.Status == "start" || event.Status == "destroy") {
			d.logger.Printf("Event: %v", event)
			rmanager.Register()
		}
	}
}

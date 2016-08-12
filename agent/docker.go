package main

import (
	"github.com/fsouza/go-dockerclient"
	"log"
)

type DockerManager struct {
	client *docker.Client
	logger *log.Logger
}

func NewDockerManager(host string, logger *log.Logger) (*DockerManager, error) {
	client, err := docker.NewClient(host)
	if err != nil {
		return nil, err
	}
	return &DockerManager{
		client: client,
		logger: logger,
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

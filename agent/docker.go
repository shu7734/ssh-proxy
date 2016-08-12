package main

import (
	"github.com/fsouza/go-dockerclient"
)

type DockerManager struct {
	client *docker.Client
}

func NewDockerManager(host string) *DockerManager {
	return &DockerManager{
		client: docker.NewClient(host),
	}
}

func (d *DockerManager) RunningContainers(prefix string) []string {
	return []string{"TestMe"}
}

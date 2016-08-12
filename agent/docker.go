package main

type DockerManager struct {
	Host string
}

func NewDockerManager(host string) *DockerManager {
	return &DockerManager{
		Host: host,
	}
}

func (d *DockerManager) RunningContainers(prefix string) []string {
	return []string{}
}

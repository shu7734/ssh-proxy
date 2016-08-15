package main

import (
	"github.com/docker/libchan"
	"github.com/vexor/ssh-proxy/commands"
)

type DockerReceiver struct {
	transport libchan.Transport
}

func NewDockerReceiver(transport *libchan.Transport) *DockerReceiver {
	return &DockerReceiver{
		transport: *transport,
	}
}

func (r *DockerReceiver) Process() error {
	return nil
}

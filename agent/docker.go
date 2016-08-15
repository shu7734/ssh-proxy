package main

import (
	"fmt"
	"github.com/docker/libchan/spdy"
	"github.com/fsouza/go-dockerclient"
	"log"
	"net"
)

type DockerManager struct {
	client   *docker.Client
	logger   *log.Logger
	listener chan *docker.APIEvents
	port     int
}

func NewDockerManager(host string, port int, logger *log.Logger) (*DockerManager, error) {
	client, err := docker.NewClient(host)
	if err != nil {
		return nil, err
	}
	return &DockerManager{
		client:   client,
		logger:   logger,
		port:     port,
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
	go d.listenCommands()
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

func (d *DockerManager) listenCommands() {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", d.port))
	if err != nil {
		d.logger.Fatal(err)
	}

	for {
		c, err := listener.Accept()
		if err != nil {
			d.logger.Print(err)
			break
		}
		p, err := spdy.NewSpdyStreamProvider(c, true)
		if err != nil {
			d.logger.Print(err)
			break
		}
		t := spdy.NewTransport(p)
		go func() {
			for {
				err := NewDockerReceiver(&t).Process()
				if err != nil {
					d.logger.Print(err)
					break
				}

			}
		}()
	}
}

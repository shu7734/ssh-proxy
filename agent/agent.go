package main

import (
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

type Agent struct {
	config   *Config
	app      *cli.App
	logger   *log.Logger
	hostname string
	RManager *RabbitMQManager
	DManager *DockerManager
}

func (a Agent) Logf(format string, v ...interface{}) {
	a.logger.Printf(format, v)
}

func (a Agent) Fatalf(format string, v ...interface{}) {
	a.Logf(format, v)
	os.Exit(1)
}

func (a Agent) failOnError(err error, message string) {
	if err != nil {
		a.Fatalf("%s: %s", message, err)
	}
}

func startAgent(c *cli.Context, config *Config) {
	logger := log.New(c.App.Writer, "AGENT: ", log.Ldate|log.Ltime|log.Lshortfile)
	agent := Agent{
		config: config,
		app:    c.App,
		logger: logger,
	}
	agent.RManager = &RabbitMQManager{
		Url:    config.RabbitmqUrl,
		Agent:  &agent,
		logger: logger,
	}
	docker_manager, err := NewDockerManager(config.DockerHost, logger)
	if err != nil {
		agent.Fatalf("Cannot get docker info: %s", err)
	}
	agent.DManager = docker_manager

	agent.start()
}

func (a *Agent) start() {
	err := a.RManager.Register()
	if err != nil {
		a.Fatalf("RabbitMQ registration error: %v", err)

	}
	a.DManager.start(a.RManager)
	// TODO: rewrite infinity loop
	for {

	}
}

func (a *Agent) ContainerIds() []string {
	return a.DManager.RunningContainers()
}

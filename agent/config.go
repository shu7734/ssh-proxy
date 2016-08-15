package main

import (
	"gopkg.in/urfave/cli.v1"
)

type Config struct {
	BindAddress     string
	BindCommandPort int
	RabbitmqUrl     string
	DockerHost      string
}

func newConfig(c *cli.Context) *Config {
	return &Config{
		BindAddress:     c.String("bind-address"),
		BindCommandPort: c.Int("bind-command-port"),
		RabbitmqUrl:     c.String("rabbitmq-url"),
		DockerHost:      c.String("docker-host"),
	}
}

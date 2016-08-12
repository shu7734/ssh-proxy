package main

import (
	"gopkg.in/urfave/cli.v1"
)

type Config struct {
	BindAddress string
	RabbitmqUrl string
	DockerHost  string
}

func newConfig(c *cli.Context) *Config {
	return &Config{
		BindAddress: c.String("bind-address"),
		RabbitmqUrl: c.String("rabbitmq-url"),
		DockerHost:  c.String("docker-host"),
	}
}

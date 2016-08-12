package main

import (
	"gopkg.in/urfave/cli.v1"
)

func initFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "rabbitmq-url",
			Value:  "amqp://guest:guest@localhost:5672/",
			Usage:  "RabbitMQ url",
			EnvVar: "RABBITMQ_URL",
		},
		cli.StringFlag{
			Name:   "bind-address",
			Value:  "0.0.0.0:2222",
			Usage:  "Bind address",
			EnvVar: "BIND_ADDRESS",
		},
		cli.StringFlag{
			Name:   "docker-host",
			Value:  "unix:///var/run/docker.sock",
			Usage:  "Docker host path or url",
			EnvVar: "DOCKER_HOST",
		},
	}
}

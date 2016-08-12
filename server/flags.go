package main

import (
	"gopkg.in/urfave/cli.v1"
)

func initFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "api-host",
			Value:  "http://localhost:3004",
			Usage:  "Api host",
			EnvVar: "API_HOST",
		},
		cli.StringFlag{
			Name:   "bind-address",
			Value:  "0.0.0.0:2222",
			Usage:  "Bind address",
			EnvVar: "BIND_ADDRESS",
		},
		cli.StringFlag{
			Name:   "shell",
			Value:  "/bin/bash",
			Usage:  "Login command",
			EnvVar: "SHELL",
		},
		cli.StringFlag{
			Name:   "host-key",
			Value:  "built-in",
			Usage:  "Host Key",
			EnvVar: "HOST_SSH_KEY",
		},
		cli.StringFlag{
			Name:   "rabbitmq-url",
			Value:  "amqp://guest:guest@localhost:5672/",
			Usage:  "RabbitMQ url",
			EnvVar: "RABBITMQ_URL",
		},
	}
}

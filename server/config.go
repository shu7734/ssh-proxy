package main

import (
	"gopkg.in/urfave/cli.v1"
)

type Config struct {
	ApiHost     string
	BindAddress string
	Shell       string
}

func newConfig(c *cli.Context) *Config {
	return &Config{
		ApiHost:     c.String("api-host"),
		BindAddress: c.String("bind-address"),
		Shell:       c.String("shell"),
	}
}

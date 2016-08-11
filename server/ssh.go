package main

import (
	"golang.org/x/crypto/ssh"
)

func initSSHConfig(config *Config) *ssh.ServerConfig {
	return &ssh.ServerConfig{}
}

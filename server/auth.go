package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

func (s Server) KeyboardInteractiveCallback(conn ssh.ConnMetadata, challenge ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
	username := conn.User()
	s.Logf("KeyboardInteractiveCallback: %q", username)
	return nil, nil
}

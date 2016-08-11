package main

import (
	"golang.org/x/crypto/ssh"
)

func (s Server) KeyboardInteractiveCallback(conn ssh.ConnMetadata, challenge ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
	username := conn.User()
	s.Logf("KeyboardInteractiveCallback: %q", username)
	// Check agents count
	// Parse username as jobID and token
	// Check if it running and SSH enabled
	return nil, nil
}

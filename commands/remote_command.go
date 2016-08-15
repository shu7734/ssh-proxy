package commands

import (
	"github.com/docker/libchan"
	"io"
)

// RemoteCommand is the received command parameters to execute and return
type RemoteCommand struct {
	Cmd        string
	Args       []string
	Stdin      io.Reader
	Stdout     io.WriteCloser
	Stderr     io.WriteCloser
	StatusChan libchan.Sender
}

// CommandResponse is the response struct to return to the client
type CommandResponse struct {
	Status int
}

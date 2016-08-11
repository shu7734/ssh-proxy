package main

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/urfave/cli.v1"
	"log"
	"net"
)

type Server struct {
	config    *Config
	app       *cli.App
	sshConfig *ssh.ServerConfig
	logger    *log.Logger
}

func startServer(c *cli.Context, config *Config) {
	server := Server{
		app:       c.App,
		config:    config,
		sshConfig: initSSHConfig(config),
		logger:    log.New(c.App.Writer, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	server.start()
}

func (s Server) start() {
	bindAddress := s.config.BindAddress

	listener, err := net.Listen("tcp", bindAddress)
	if err != nil {
		s.logger.Fatalf("Failed to start listener on %q: %v", bindAddress, err)
	}
	s.logger.Printf("listener started on %q", bindAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Printf("Accept failed: %v", err)
			continue

		}
		go s.Handle(conn)
	}
}

func (s Server) Handle(netConn net.Conn) error {
	s.logger.Printf("Server.Handle netConn=%v", netConn)
	return nil
}

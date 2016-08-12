package main

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

type Server struct {
	Config    *Config
	app       *cli.App
	SshConfig *ssh.ServerConfig
	logger    *log.Logger
	Banner    string
}

func startServer(c *cli.Context, config *Config) {
	server := Server{
		app:    c.App,
		Config: config,
		logger: log.New(c.App.Writer, "SERVER: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	server.SshConfig = &ssh.ServerConfig{
		KeyboardInteractiveCallback: server.KeyboardInteractiveCallback,
		PasswordCallback:            nil,
	}

	// Register the SSH host key
	hostKey := c.String("host-key")
	switch hostKey {
	case "built-in":
		hostKey = DefaultHostKey
	}

	err := server.AddHostKey(hostKey)
	if err != nil {
		server.Logf("Cannot add host key: %v", err)
	}

	server.start()
}

func (s Server) Logf(format string, v ...interface{}) {
	s.logger.Printf(format, v)
}

func (s Server) Fatalf(format string, v ...interface{}) {
	s.Logf(format, v)
	os.Exit(1)
}

func (s Server) failOnError(err error, message string) {
	if err != nil {
		s.Fatalf("%s: %s", message, err)
	}
}

func (s Server) receiveAgents() {
	agents_receiver := NewAgentsReceiver(&s)
	go agents_receiver.receive()
}

func (s Server) initListener() {
	bindAddress := s.Config.BindAddress
	listener, err := net.Listen("tcp", bindAddress)
	if err != nil {
		s.Logf("Failed to start listener on %q: %v", bindAddress, err)
	}
	s.Logf("listener started on %q", bindAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.Logf("Accept failed: %v", err)
			continue

		}
		go s.Handle(conn)
	}
}

func (s Server) start() {
	s.receiveAgents()
	s.initListener()
}

func (s Server) Handle(netConn net.Conn) error {
	s.Logf("Server.Handle netConn=%v", netConn)

	conn, chans, reqs, err := ssh.NewServerConn(netConn, s.SshConfig)
	if err != nil {
		s.Logf("Received disconnect from %s: 11: Bye Bye [preauth]", netConn.RemoteAddr().String())
		s.Logf("Error: %v", err)
		return err
	}
	client := NewClient(conn, chans, reqs, &s)

	// Handle requests
	if err = client.HandleRequests(); err != nil {
		return err
	}

	// Handle channels
	if err = client.HandleChannels(); err != nil {
		return err
	}

	return nil
}

func (s *Server) AddHostKey(keystring string) error {
	// Check if keystring is a key path or a key string
	keypath := os.ExpandEnv(strings.Replace(keystring, "~", "$HOME", 2))
	_, err := os.Stat(keypath)
	var keybytes []byte
	if err == nil {
		keybytes, err = ioutil.ReadFile(keypath)
		if err != nil {
			return err
		}
	} else {
		keybytes = []byte(keystring)
	}

	// Parse SSH priate key
	hostkey, err := ssh.ParsePrivateKey(keybytes)
	if err != nil {
		return err
	}

	// Register key to the server
	s.SshConfig.AddHostKey(hostkey)
	return nil
}

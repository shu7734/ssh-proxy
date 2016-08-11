package main

import (
	"fmt"
	"github.com/flynn/go-shlex"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var clientCounter = 0

type Client struct {
	Idx        int
	ChannelIdx int
	Conn       *ssh.ServerConn
	Chans      <-chan ssh.NewChannel
	Reqs       <-chan *ssh.Request
	Server     *Server
	Pty, Tty   *os.File
	ClientID   string
	Config     *ClientConfig
}

type ClientConfig struct {
	RemoteUser string            `json:"remote-user,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	Command    []string          `json:"command,omitempty"`
	User       string            `json:"user,omitempty"`
	Keys       []string          `json:"keys,omitempty"`
	EntryPoint string            `json:"entrypoint,omitempty"`
	Allowed    bool              `json:"allowed,omitempty"`
	IsLocal    bool              `json:"is-local,omitempty"`
	UseTTY     bool              `json:"use-tty,omitempty"`
}

func NewClient(conn *ssh.ServerConn, chans <-chan ssh.NewChannel, reqs <-chan *ssh.Request, server *Server) *Client {
	client := Client{
		Idx:        clientCounter,
		ClientID:   conn.RemoteAddr().String(),
		ChannelIdx: 0,
		Conn:       conn,
		Chans:      chans,
		Reqs:       reqs,
		Server:     server,
	}
	// Default ClientConfig, will be overwritten if a hook is used
	client.Config = &ClientConfig{
		Env:     make(map[string]string),
		Command: make([]string, 0),
	}

	clientCounter++
	remoteAddr := strings.Split(client.ClientID, ":")
	server.Logf("Accepted for %s from %s port %s", conn.User(), remoteAddr[0], remoteAddr[1])
	return &client
}

func (c *Client) Logf(format string, v ...interface{}) {
	c.Server.Logf(format, v)
}

// HandleRequests handles SSH requests
func (c *Client) HandleRequests() error {
	go func(in <-chan *ssh.Request) {
		for req := range in {
			c.Logf("HandleRequest: %v", req)
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}(c.Reqs)
	return nil
}

// HandleChannels handles SSH channels
func (c *Client) HandleChannels() error {
	for newChannel := range c.Chans {
		if err := c.HandleChannel(newChannel); err != nil {
			return err
		}
	}
	return nil
}

// HandleChannel handles one SSH channel
func (c *Client) HandleChannel(newChannel ssh.NewChannel) error {
	if newChannel.ChannelType() != "session" {
		c.Logf("Unknown channel type: %s", newChannel.ChannelType())
		newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
		return nil
	}

	channel, requests, err := newChannel.Accept()
	if err != nil {
		c.Logf("newChannel.Accept failed: %v", err)
		return err
	}
	c.ChannelIdx++
	c.Logf("HandleChannel.channel (client=%d channel=%d)", c.Idx, c.ChannelIdx)

	c.Logf("Creating pty...")
	c.Pty, c.Tty, err = pty.Open()
	if err != nil {
		c.Logf("pty.Open failed: %v", err)
		return nil
	}

	c.HandleChannelRequests(channel, requests)

	return nil
}

func (c *Client) HandleChannelRequests(channel ssh.Channel, requests <-chan *ssh.Request) {
	go func(in <-chan *ssh.Request) {
		defer c.Tty.Close()
		for req := range in {
			ok := false
			switch req.Type {
			case "shell":
				c.Logf("HandleChannelRequests.req shell")
				if len(req.Payload) != 0 {
					break
				}
				ok = true

				entrypoint := ""
				if c.Config.EntryPoint != "" {
					entrypoint = c.Config.EntryPoint
				}

				var args []string
				if c.Config.Command != nil {
					args = c.Config.Command
				}

				if entrypoint == "" && len(args) == 0 {
					//args = []string{c.Server.DefaultShell}
					args = []string{"/bin/bash"}
				}

				c.runCommand(channel, entrypoint, args)

			case "exec":
				command := string(req.Payload[4:])
				c.Logf("HandleChannelRequests.req exec: %q", command)
				ok = true

				args, err := shlex.Split(command)
				if err != nil {
					c.Logf("Failed to parse command %q: %v", command, args)
				}
				c.runCommand(channel, c.Config.EntryPoint, args)

			case "pty-req":
				ok = true
				c.Config.UseTTY = true
				termLen := req.Payload[3]
				c.Config.Env["TERM"] = string(req.Payload[4 : termLen+4])
				c.Config.Env["USE_TTY"] = "1"
				// w, h := ttyhelper.ParseDims(req.Payload[termLen+4:])
				// ttyhelper.SetWinsize(c.Pty.Fd(), w, h)
				// log.Debugf("HandleChannelRequests.req pty-req: TERM=%q w=%q h=%q", c.Config.Env["TERM"], int(w), int(h))
				c.Logf("HandleChannelRequests.req pty-req: TERM=%q", c.Config.Env["TERM"])

			// case "window-change":
			// 	w, h := ttyhelper.ParseDims(req.Payload)
			// 	ttyhelper.SetWinsize(c.Pty.Fd(), w, h)
			// 	continue

			case "env":
				keyLen := req.Payload[3]
				key := string(req.Payload[4 : keyLen+4])
				valueLen := req.Payload[keyLen+7]
				value := string(req.Payload[keyLen+8 : keyLen+8+valueLen])
				c.Logf("HandleChannelRequets.req 'env': %s=%q", key, value)
				c.Config.Env[key] = value

			default:
				c.Logf("Unhandled request type: %q: %v", req.Type, req)
			}

			if req.WantReply {
				if !ok {
					c.Logf("Declining %s request...", req.Type)
				}
				req.Reply(ok, nil)
			}
		}

	}(requests)
}

func (c *Client) runCommand(channel ssh.Channel, entrypoint string, command []string) {
	c.Logf("Command: %s", strings.Join(command, " "))
	var cmd *exec.Cmd
	cmd = exec.Command(entrypoint, command...)

	fmt.Fprintf(channel, "%s\n\r", "Welcome to Vexor.io!")

	cmd.Stdout = channel
	cmd.Stdin = channel
	cmd.Stderr = channel
	var wg sync.WaitGroup
	if c.Config.UseTTY {
		cmd.Stdout = c.Tty
		cmd.Stdin = c.Tty
		cmd.Stderr = c.Tty

		wg.Add(1)
		go func() {
			io.Copy(channel, c.Pty)
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			io.Copy(c.Pty, channel)
			wg.Done()
		}()
		defer wg.Wait()
	}

	channel.Close()
	c.Logf("cmd.Wait done")
}

package main

import (
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
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
	c.Logf("HandleChannelRequests: %v", requests)
}

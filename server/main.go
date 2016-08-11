package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	server := cli.NewApp()
	server.Usage = "Run ssh proxy server"
	server.Flags = initFlags()
	server.Action = func(c *cli.Context) error {
		config := newConfig(c)
		startServer(c, config)
		return nil
	}
	server.Run(os.Args)
}

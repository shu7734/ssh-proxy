package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	server := cli.NewApp()
	server.Usage = "Run ssh proxy server"
	server.Run(os.Args)
}

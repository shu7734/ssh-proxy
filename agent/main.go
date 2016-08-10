package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	agent := cli.NewApp()
	agent.Usage = "Ssh agent"
	agent.Run(os.Args)
}

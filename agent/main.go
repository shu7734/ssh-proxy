package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
)

func main() {
	agent := cli.NewApp()
	agent.Usage = "Ssh agent"
	agent.Flags = initFlags()
	agent.Action = func(c *cli.Context) error {
		config := newConfig(c)
		startAgent(c, config)
		return nil
	}
	agent.Run(os.Args)
}

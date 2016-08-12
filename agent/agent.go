package main

import (
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

type Agent struct {
	config   *Config
	app      *cli.App
	logger   *log.Logger
	hostname string
	Manager  *RabbitMQManager
}

func (a Agent) Logf(format string, v ...interface{}) {
	a.logger.Printf(format, v)
}

func (a Agent) Fatalf(format string, v ...interface{}) {
	a.Logf(format, v)
	os.Exit(1)
}

func startAgent(c *cli.Context, config *Config) {
	agent := Agent{
		config: config,
		app:    c.App,
		logger: log.New(c.App.Writer, "AGENT: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	agent.Manager = &RabbitMQManager{
		Url:   config.RabbitmqUrl,
		Agent: &agent,
	}
	agent.start()
}

func (a *Agent) start() {
	err := a.Manager.Register()
	if err != nil {
		a.Fatalf("RabbitMQ registration error: %v", err)

	}
}

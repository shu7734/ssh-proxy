package main

import (
	"github.com/streadway/amqp"
)

type RabbitMQManager struct {
	Url        string
	connection *amqp.Connection
	Agent      *Agent
}

func (m *RabbitMQManager) Connect() error {
	connection, err := amqp.Dial(m.Url)
	if err != nil {
		return err
	}
	m.connection = connection
	return nil
}

func (m *RabbitMQManager) Disconnect() {
	if m.connection != nil {
		m.connection.Close()
	}
}

func (m *RabbitMQManager) PushNodeInformation() error {
	return nil
}

func (m *RabbitMQManager) Register() error {
	if err := m.Connect(); err != nil {
		return err
	}
	defer m.Disconnect()
	m.Agent.Logf("RammitMQ connection SUCCESSED")
	m.PushNodeInformation()

	return nil
}

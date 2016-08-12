package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"github.com/vexor/ssh-proxy/messages"
	"net"
	"os"
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
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	hosts, err := net.LookupHost(hostname)
	if err != nil {
		m.Agent.Logf("Error: %v", err)
		return err
	}
	if len(hosts) < 1 {
		m.Agent.Fatalf("Can't get ip address by hostname (%s)", hostname)
	}

	node_info := messages.NodeInfo{
		Host: proto.String(hosts[0]),
	}

	m.Agent.Logf("Hostname: %v", node_info.GetHost())

	m.publish(&node_info)
	return nil
}

func (m *RabbitMQManager) publish(node_info *messages.NodeInfo) {
	ch, err := m.connection.Channel()
	m.Agent.failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"ssh-proxy-agents", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	m.Agent.failOnError(err, "Failed to declare a queue")

	out, err := proto.Marshal(node_info)
	m.Agent.failOnError(err, "Failed to encode address book:")

	body := out
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	m.Agent.failOnError(err, "Failed to publish a message")

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

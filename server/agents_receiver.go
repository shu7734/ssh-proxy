package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	"github.com/vexor/ssh-proxy/messages"
)

type AgentsReceiver struct {
	connection *amqp.Connection
	server     *Server
}

func NewAgentsReceiver(server *Server) *AgentsReceiver {
	conn, err := amqp.Dial(server.Config.RabbitmqUrl)
	server.failOnError(err, "Cannot connect to RabbitMQ")

	agents_receiver := AgentsReceiver{
		connection: conn,
		server:     server,
	}
	return &agents_receiver
}

func (r AgentsReceiver) receive() {
	r.server.Logf("Try to crate channel...")
	ch, err := r.connection.Channel()
	r.server.failOnError(err, "Failed to open a channel")

	defer r.connection.Close()
	defer ch.Close()
	r.server.Logf("Connected...")

	q, err := ch.QueueDeclare(
		"ssh-proxy-agents", // name
		false,              // durable
		false,              // delete when usused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	r.server.failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	r.server.failOnError(err, "Failed to register a consumer")

	agents_consumer := &AgentsConsumer{
		server:       r.server,
		node_storage: NewNodeStorage(),
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			node_info := &messages.NodeInfo{}
			if err := proto.Unmarshal(d.Body, node_info); err != nil {
				r.server.Fatalf("Failed to parse NodeInfo:", err)
			}
			agents_consumer.Apply(node_info)
		}
	}()
	r.server.Logf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

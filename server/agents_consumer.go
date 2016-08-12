package main

import (
	"github.com/vexor/ssh-proxy/messages"
)

type AgentsConsumer struct {
	server       *Server
	node_storage *NodeStorage
}

func (c AgentsConsumer) Apply(node_info *messages.NodeInfo) {
	c.server.Logf("Received message: %s", node_info.GetHost())
	c.node_storage.Put(node_info)
	c.server.Logf("Nodes: %v", c.node_storage.store)
}

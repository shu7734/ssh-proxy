package main

import (
	"github.com/vexor/ssh-proxy/messages"
)

type NodeStorage struct {
	store map[string][]string
}

func NewNodeStorage() *NodeStorage {
	return &NodeStorage{
		store: make(map[string][]string),
	}
}

func (ns NodeStorage) Put(node_info *messages.NodeInfo) {
	ns.store[node_info.GetHost()] = node_info.GetContainers()
}

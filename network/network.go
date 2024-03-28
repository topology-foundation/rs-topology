package network

import (
	"context"
	"fmt"

	"github.com/topology-gg/gram/config"
	"github.com/topology-gg/gram/execution"
	"github.com/topology-gg/gram/storage"
)

type NetworkModule struct {
	ctx       context.Context
	execution execution.Execution
	storage   storage.Storage
	config    *config.NetworkConfig
	p2p       *P2P
	rpc       *RPC
	grpc      *GRPC
}

type NetworkMediator interface {
	MessageHandler(message string, source Source)
}

type Source int

const (
	SourceP2P Source = iota
	SourceRPC
)

func NewNetwork(ctx context.Context, execution execution.Execution, storage storage.Storage, config *config.NetworkConfig) *NetworkModule {
	return &NetworkModule{
		ctx:       ctx,
		execution: execution,
		storage:   storage,
		config:    config,
		p2p:       nil,
		rpc:       nil,
		grpc:      nil,
	}
}

func (network *NetworkModule) Start() {
	p2p := NewP2P(network.ctx, network, network.config.Namespace, network.config.MaxPeers)
	network.p2p = p2p

	fmt.Println("(Network) Host ID:", p2p.host.ID())
	fmt.Println("(Network) Host addresses:", p2p.host.Addrs())

	go p2p.JoinNetwork()
	go p2p.SubscribeTopics(network.config.Topics)

	grpc := NewGRPC(network.ctx)
	network.grpc = grpc

	go grpc.Start()

	rpc := NewRPC(network.ctx, network)
	network.rpc = rpc

	rpc.Start()
}

func (network *NetworkModule) MessageHandler(message string, source Source) {
	if source == SourceRPC {
		message = fmt.Sprintf("%s: %s", network.p2p.host.ID().String(), message)
	}

	fmt.Printf("(Network) %s", message)

	network.execution.Execute(message)

	if source == SourceRPC {
		network.p2p.Publish(message)
	}
}

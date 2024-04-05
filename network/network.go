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
	p2p       *P2P
	grpc      *GRPC
	rpc       *RPC
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
	network := &NetworkModule{
		ctx:       ctx,
		execution: execution,
		storage:   storage,
	}

	network.p2p = NewP2P(ctx, network, &config.P2p)
	network.grpc = NewGRPC(ctx, &config.Grpc)
	network.rpc = NewRPC(ctx, network, &config.Rpc)

	return network
}

func (network *NetworkModule) Start() {
	p2p := network.p2p
	grpc := network.grpc
	rpc := network.rpc

	fmt.Println("(Network) Host ID:", p2p.host.ID())
	fmt.Println("(Network) Host addresses:", p2p.host.Addrs())

	go p2p.Start()
	go grpc.Start()
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

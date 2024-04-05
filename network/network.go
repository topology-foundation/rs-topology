package network

import (
	"context"
	"fmt"
	"os"

	"github.com/topology-gg/gram/config"
	"github.com/topology-gg/gram/execution"
	"github.com/topology-gg/gram/storage"
)

type NetworkModule struct {
	ctx       context.Context
	errCh     chan error
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

func NewNetwork(ctx context.Context, errCh chan error, execution execution.Execution, storage storage.Storage, config *config.NetworkConfig) (*NetworkModule, error) {
	network := &NetworkModule{
		ctx:       ctx,
		errCh:     errCh,
		execution: execution,
		storage:   storage,
	}

	p2p, err := NewP2P(ctx, errCh, network, &config.P2p)
	if err != nil {
		return nil, err
	}

	grpc, err := NewGRPC(ctx, errCh, &config.Grpc)
	if err != nil {
		return nil, err
	}

	rpc, err := NewRPC(ctx, errCh, network, &config.Rpc)
	if err != nil {
		return nil, err
	}

	network.p2p = p2p
	network.grpc = grpc
	network.rpc = rpc

	return network, nil
}

func (network *NetworkModule) Start() {
	go network.p2p.Start()
	go network.grpc.Start()
	go network.rpc.Start()
}

// Shutdown gracefuly shutdowns network modules
func (network *NetworkModule) Shutdown() error {
	if err := network.rpc.Shutdown(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if err := network.grpc.Shutdown(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return nil
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

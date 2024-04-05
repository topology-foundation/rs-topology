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
	ctx        context.Context
	errCh      chan error
	execution  execution.Execution
	storage    storage.Storage
	networkCfg *config.NetworkConfig
	grpcCfg    *config.GrpcConfig
	p2p        *P2P
	rpc        *RPC
	grpc       *GRPC
}

type NetworkMediator interface {
	MessageHandler(message string, source Source)
}

type Source int

const (
	SourceP2P Source = iota
	SourceRPC
)

func NewNetwork(ctx context.Context, errCh chan error, execution execution.Execution, storage storage.Storage, config *config.NetworkConfig, grpcCfg *config.GrpcConfig) (*NetworkModule, error) {
	return &NetworkModule{
		ctx:        ctx,
		errCh:      errCh,
		execution:  execution,
		storage:    storage,
		networkCfg: config,
		grpcCfg:    grpcCfg,
		p2p:        nil,
		rpc:        nil,
		grpc:       nil,
	}, nil
}

func (network *NetworkModule) Start() {
	p2p, err := NewP2P(network.ctx, network.errCh, network, network.networkCfg.Namespace, network.networkCfg.MaxPeers, network.networkCfg.Port)
	if err != nil {
		network.errCh <- err
		return
	}

	network.p2p = p2p
	go network.p2p.JoinNetwork()
	go network.p2p.SubscribeTopics(network.networkCfg.Topics)

	grpc, err := NewGRPC(network.ctx, network.errCh, network.grpcCfg)
	if err != nil {
		network.errCh <- err
		return
	}

	network.grpc = grpc
	go network.grpc.Start()

	rpc, err := NewRPC(network.ctx, network)
	if err != nil {
		network.errCh <- err
		return
	}

	network.rpc = rpc
	go network.rpc.Start()
}

// Shutdown gracefuly shutdowns network modules
func (network *NetworkModule) Shutdown() error {
	// TODO: add other modules shutdown
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

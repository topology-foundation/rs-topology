package network

import (
	"context"
	"fmt"
	"os"

	"github.com/topology-gg/gram/config"
	"github.com/topology-gg/gram/execution"
	"github.com/topology-gg/gram/log"
	grpc "github.com/topology-gg/gram/network/grpc"
	p2p "github.com/topology-gg/gram/network/p2p"
	rpc "github.com/topology-gg/gram/network/rpc"
	"github.com/topology-gg/gram/storage"
)

type NetworkModule struct {
	ctx       context.Context
	errCh     chan error
	execution execution.Execution
	storage   storage.Storage
	p2p       *p2p.P2P
	grpc      *grpc.GRPC
	rpc       *rpc.RPC
}

func NewNetwork(ctx context.Context, errCh chan error, execution execution.Execution, storage storage.Storage, config *config.NetworkConfig) (*NetworkModule, error) {
	network := &NetworkModule{
		ctx:       ctx,
		errCh:     errCh,
		execution: execution,
		storage:   storage,
	}

	p2p, err := p2p.NewP2P(ctx, errCh, execution, &config.P2p)
	if err != nil {
		return nil, err
	}

	grpc, err := grpc.NewGRPC(ctx, errCh, &config.Grpc)
	if err != nil {
		return nil, err
	}

	rpc, err := rpc.NewRPC(ctx, errCh, execution, &config.Rpc, p2p)
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
func (network *NetworkModule) Shutdown() {
	if err := network.rpc.Shutdown(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if err := network.grpc.Shutdown(); err != nil {
		log.Error("(Network) GRPC shutdown", "error", err)
	}

	if err := network.p2p.Shutdown(); err != nil {
		log.Error("(Network) P2P shutdown", "error", err)
		fmt.Fprintln(os.Stderr, err)
	}
}

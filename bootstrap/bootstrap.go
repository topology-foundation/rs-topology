package bootstrap

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
)

type Bootstrap struct {
	ctx  context.Context
	host host.Host
}

func New(ctx context.Context, listenAddr string) (*Bootstrap, error) {
	multiaddr, err := multiaddr.NewMultiaddr(listenAddr)
	if err != nil {
		return nil, fmt.Errorf("creating multiaddr failed: %w", err)
	}

	nodeHost, err := libp2p.New(libp2p.ListenAddrs(multiaddr))
	if err != nil {
		panic(err)
	}

	return &Bootstrap{
		ctx:  ctx,
		host: nodeHost,
	}, nil
}

func (b *Bootstrap) Start() error {
	kademliaDHT, err := dht.New(b.ctx, b.host, dht.Mode(dht.ModeServer))
	if err != nil {
		return fmt.Errorf("creating DHT failed: %w", err)
	}

	if err := kademliaDHT.Bootstrap(b.ctx); err != nil {
		return fmt.Errorf("bootstrapping DHT failed: %w", err)
	}

	for _, addr := range b.host.Addrs() {
		fullAddr := addr.Encapsulate(multiaddr.StringCast(fmt.Sprintf("/p2p/%s", b.host.ID())))
		fmt.Printf("I am available at: %s\n", fullAddr)
	}

	return nil
}

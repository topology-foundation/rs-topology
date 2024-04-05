package network

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/topology-gg/gram/config"
)

type P2P struct {
	ctx       context.Context
	mediator  NetworkMediator
	host      host.Host
	namespace string
	maxPeers  int
	port      int
	pubsub    *pubsub.PubSub
	streams   []Stream
}

type Stream struct {
	name         string
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
}

func NewP2P(ctx context.Context, mediator NetworkMediator, cfg *config.P2pConfig) *P2P {
	namespace, maxPeers, port := cfg.Namespace, cfg.MaxPeers, cfg.Port

	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)

	host, err := libp2p.New(libp2p.ListenAddrStrings(listenAddr))
	if err != nil {
		panic(err)
	}

	gossipsub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}

	streams := make([]Stream, len(cfg.Topics))

	for i, name := range cfg.Topics {
		streams[i] = Stream{name: name}
	}

	return &P2P{
		ctx:       ctx,
		mediator:  mediator,
		host:      host,
		namespace: namespace,
		maxPeers:  maxPeers,
		port:      port,
		pubsub:    gossipsub,
		streams:   streams,
	}
}

func (p2p *P2P) Start() {
	p2p.joinNetwork()
	p2p.subscribeTopics()
}

func (p2p *P2P) Publish(message string) {
	for i := range p2p.streams {
		if err := p2p.streams[i].topic.Publish(p2p.ctx, []byte(message)); err != nil {
			fmt.Println("(Network) Failed to publish to topic:", p2p.streams[i].name, err)
		}
	}
}

func (p2p *P2P) joinNetwork() {
	kademliaDHT := p2p.getKademliaDHT()
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(p2p.ctx, routingDiscovery, p2p.namespace)
	p2p.connectPeers(routingDiscovery)

	fmt.Println("(Network) Successfully joinned network:", p2p.namespace)
}

func (p2p *P2P) subscribeTopics() {
	for i := range p2p.streams {
		topic, err := p2p.pubsub.Join(p2p.streams[i].name)
		if err != nil {
			panic(err)
		}

		p2p.streams[i].topic = topic

		subscription, err := topic.Subscribe()
		if err != nil {
			panic(err)
		}

		p2p.streams[i].subscription = subscription

		go p2p.p2pMessageHandler(subscription)

		fmt.Println("(Network) Successfully subscribed to gossipsub topic:", p2p.streams[i].name)
	}
}

func (p2p *P2P) getKademliaDHT() *dht.IpfsDHT {
	kademliaDHT, err := dht.New(p2p.ctx, p2p.host)
	if err != nil {
		panic(err)
	}

	if err := kademliaDHT.Bootstrap(p2p.ctx); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for i := range dht.DefaultBootstrapPeers {
		peerInfo, err := peer.AddrInfoFromP2pAddr(dht.DefaultBootstrapPeers[i])
		if err != nil {
			panic(err)
		}

		go func() {
			defer wg.Done()

			if err := p2p.host.Connect(p2p.ctx, *peerInfo); err != nil {
				fmt.Println("(Network) Failed to connect to bootstrap node:", err)
			} else {
				fmt.Println("(Network) Successfully connected to bootstrap node:", peerInfo)
			}
		}()

		wg.Add(1)
	}
	wg.Wait()

	return kademliaDHT
}

func (p2p *P2P) connectPeers(routingDiscovery *drouting.RoutingDiscovery) {
	peers := 0
	isConnected := false

	for !isConnected {
		fmt.Println("(Network) Searching for peers to connect...")

		peerInfoChan, err := routingDiscovery.FindPeers(p2p.ctx, p2p.namespace)
		if err != nil {
			panic(err)
		}

		for peerInfo := range peerInfoChan {
			if peerInfo.ID == p2p.host.ID() {
				continue
			}

			if err := p2p.host.Connect(p2p.ctx, peerInfo); err != nil {
				fmt.Println("(Network) Failed to connect to peer:", err)
			} else {
				peers += 1
				isConnected = true
				fmt.Println("(Network) Successfully connected to peer:", peerInfo)
			}

			if peers >= p2p.maxPeers {
				break
			}
		}
	}

	fmt.Println("(Network) Connecting peers is completed")
}

func (p2p *P2P) p2pMessageHandler(subscription *pubsub.Subscription) {
	for {
		message, err := subscription.Next(p2p.ctx)
		if err != nil {
			panic(err)
		}

		sender := message.ReceivedFrom
		if sender == p2p.host.ID() {
			continue
		}

		p2p.mediator.MessageHandler(string(message.Message.Data), SourceP2P)
	}
}

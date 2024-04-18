package p2p

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
	ex "github.com/topology-gg/gram/execution"
	"github.com/topology-gg/gram/log"
	pbcodec "github.com/topology-gg/gram/proto"
	"github.com/topology-gg/gram/proto/gen/gram/base"
)

type P2P struct {
	ctx        context.Context
	errCh      chan error
	executor   ex.Execution
	host       host.Host
	namespace  string
	maxPeers   int
	port       int
	pubsub     *pubsub.PubSub
	streams    []Stream
	serializer pbcodec.Serializer
}

type Stream struct {
	name         string
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
}

func NewP2P(ctx context.Context, errCh chan error, executor ex.Execution, cfg *config.P2pConfig) (*P2P, error) {
	namespace, maxPeers, port := cfg.Namespace, cfg.MaxPeers, cfg.Port

	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)

	host, err := libp2p.New(libp2p.ListenAddrStrings(listenAddr))
	if err != nil {
		return nil, err
	}

	fmt.Println("(Network) Host ID:", host.ID())
	fmt.Println("(Network) Host addresses:", host.Addrs())

	gossipsub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}

	streams := make([]Stream, len(cfg.Topics))

	for i, name := range cfg.Topics {
		streams[i] = Stream{name: name}
	}

	serializer := &pbcodec.ProtoBufSerializer{}

	return &P2P{
		ctx:        ctx,
		errCh:      errCh,
		executor:   executor,
		host:       host,
		namespace:  namespace,
		maxPeers:   maxPeers,
		port:       port,
		pubsub:     gossipsub,
		streams:    streams,
		serializer: serializer,
	}, nil
}

func (p2p *P2P) Start() {
	p2p.joinNetwork()
	p2p.subscribeTopics()
}

func (p2p *P2P) Publish(message string) {
	msg, err := p2p.serializer.Marshal(&base.HelloRequest{Name: message})
	if err != nil {
		log.Error("(Network) Failed to serialize message", "message", err)
		return
	}

	for i := range p2p.streams {
		if err := p2p.streams[i].topic.Publish(p2p.ctx, msg); err != nil {
			log.Error("(Network) Failed to publish to topic", "topic", p2p.streams[i].name, "error", err)

		}
	}
}

func (p2p *P2P) joinNetwork() {
	kademliaDHT, err := p2p.getKademliaDHT()
	if err != nil {
		p2p.errCh <- err
		return
	}

	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(p2p.ctx, routingDiscovery, p2p.namespace)

	err = p2p.connectPeers(routingDiscovery)
	if err != nil {
		p2p.errCh <- err
		return
	}

	log.Info("(Network) Successfully joined network", "namespace", p2p.namespace)
}

func (p2p *P2P) subscribeTopics() {
	for i := range p2p.streams {
		topic, err := p2p.pubsub.Join(p2p.streams[i].name)
		if err != nil {
			p2p.errCh <- err
			return
		}

		p2p.streams[i].topic = topic

		subscription, err := topic.Subscribe()
		if err != nil {
			p2p.errCh <- err
			return
		}

		p2p.streams[i].subscription = subscription

		go p2p.p2pMessageHandler(subscription)

		log.Info("(Network) Successfully subscribed to gossipsub topic", "topic", p2p.streams[i].name)
	}
}

func (p2p *P2P) getKademliaDHT() (*dht.IpfsDHT, error) {
	kademliaDHT, err := dht.New(p2p.ctx, p2p.host)
	if err != nil {
		return nil, err
	}

	if err := kademliaDHT.Bootstrap(p2p.ctx); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for i := range dht.DefaultBootstrapPeers {
		peerInfo, err := peer.AddrInfoFromP2pAddr(dht.DefaultBootstrapPeers[i])
		if err != nil {
			return nil, err
		}

		go func() {
			defer wg.Done()

			if err := p2p.host.Connect(p2p.ctx, *peerInfo); err != nil {
				log.Error("(Network) Failed to connect to bootstrap node", "error", err)
			} else {
				log.Info("(Network) Successfully connected to bootstrap node", "peerInfo", peerInfo)
			}
		}()

		wg.Add(1)
	}
	wg.Wait()

	return kademliaDHT, nil
}

func (p2p *P2P) connectPeers(routingDiscovery *drouting.RoutingDiscovery) error {
	peers := 0
	isConnected := false

	for !isConnected {
		log.Info("(Network) Searching for peers to connect...")

		peerInfoChan, err := routingDiscovery.FindPeers(p2p.ctx, p2p.namespace)
		if err != nil {
			return err
		}

		for peerInfo := range peerInfoChan {
			if peerInfo.ID == p2p.host.ID() {
				continue
			}

			if err := p2p.host.Connect(p2p.ctx, peerInfo); err != nil {
				log.Error("(Network) Failed to connect to peer", "error", err)
			} else {
				peers++
				isConnected = true
				log.Info("(Network) Successfully connected to peer", "peerInfo", peerInfo)
			}

			if peers >= p2p.maxPeers {
				break
			}
		}
	}

	log.Info("(Network) Connecting peers is completed")
	return nil
}

func (p2p *P2P) p2pMessageHandler(subscription *pubsub.Subscription) {
	for {
		message, err := subscription.Next(p2p.ctx)
		if err != nil {
			log.Error("(Network) Error handling P2P message", "error", err)
			continue
		}

		sender := message.ReceivedFrom
		if sender == p2p.host.ID() {
			continue
		}

		var msg base.HelloRequest
		if err := p2p.serializer.Unmarshal(message.Data, &msg); err != nil {
			log.Error("(Network) Failed to deserialize message", "message", err)
			continue
		}

		p2p.executor.Execute(string(message.Message.Data))
	}
}

// Shutdown gracefuly shutdowns p2p communication
func (p2p *P2P) Shutdown() error {
	for i := range p2p.streams {
		if p2p.streams[i].subscription != nil {
			p2p.streams[i].subscription.Cancel()
		}

		if p2p.streams[i].topic != nil {
			if err := p2p.streams[i].topic.Close(); err != nil {
				// just log the error here, since we need to try to close other topics
				log.Error("(Network) Error closing topic", "error", err)
			}
		}
	}

	if err := p2p.host.Close(); err != nil {
		return err
	}

	log.Info("(Network) P2P host successfully shutted down")
	return nil
}

func (p2p *P2P) HostId() string {
	return p2p.host.ID().String()
}

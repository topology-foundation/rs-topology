package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/topology-gg/gram/config"
	"github.com/topology-gg/gram/log"
	helloPb "github.com/topology-gg/gram/proto/gen/gram/base"
	"google.golang.org/grpc"
)

// GRPC represents a struct for GRPC server
type GRPC struct {
	ctx    context.Context
	errCh  chan error
	config *config.GrpcConfig
	server *grpc.Server
}

// NewGRPC creates a new grpc server struct
func NewGRPC(ctx context.Context, errCh chan error, config *config.GrpcConfig) (*GRPC, error) {
	return &GRPC{
		ctx:    ctx,
		errCh:  errCh,
		config: config,
	}, nil
}

// Start starts a GRPC server
func (g *GRPC) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g.config.Port))
	if err != nil {
		g.errCh <- err
		return
	}

	g.server = grpc.NewServer()

	helloPb.RegisterServiceServer(g.server, &helloServer{})
	log.Info("(GRPC Server)", "address", listener.Addr())

	if err := g.server.Serve(listener); err != nil {
		g.errCh <- err
	}
}

// Shutdown gracefuly shutdowns grpc server
func (g *GRPC) Shutdown() error {
	g.server.GracefulStop()

	log.Info("(GRPC Server) successfully shutted down")
	return nil
}

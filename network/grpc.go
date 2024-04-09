package network

import (
	"context"
	"fmt"
	"net"

	"github.com/topology-gg/gram/config"
	helloPb "github.com/topology-gg/gram/proto/gen/gram/base"
	"google.golang.org/grpc"
)

// TODO: move server implementation to module specific folders
type helloServer struct {
	helloPb.UnimplementedServiceServer
}

func (s *helloServer) SayHello(ctx context.Context, in *helloPb.HelloRequest) (*helloPb.HelloResponse, error) {
	fmt.Printf("Received greet from %s", in.GetName())
	return &helloPb.HelloResponse{Name: "Hello " + in.GetName()}, nil
}

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
	fmt.Printf("GRPC Server is listening on port: %v\n", listener.Addr())

	if err := g.server.Serve(listener); err != nil {
		g.errCh <- err
	}
}

// Shutdown gracefuly shutdowns grpc server
func (g *GRPC) Shutdown() error {
	g.server.GracefulStop()

	fmt.Println("GRPC server successfully shutted down")
	return nil
}

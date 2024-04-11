package grpc

import (
	"context"
	"fmt"

	helloPb "github.com/topology-gg/gram/proto/gen/gram/base"
)

type helloServer struct {
	helloPb.UnimplementedServiceServer
}

func (s *helloServer) SayHello(ctx context.Context, in *helloPb.HelloRequest) (*helloPb.HelloResponse, error) {
	fmt.Printf("Received greet from %s", in.GetName())
	return &helloPb.HelloResponse{Name: "Hello " + in.GetName()}, nil
}

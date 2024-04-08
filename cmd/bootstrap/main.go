package main

import (
	"context"
	"fmt"
	"os"

	"github.com/topology-gg/gram/bootstrap"
)

func main() {
	ctx := context.Background()

	defaultPort := "4001"
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", port)

	bootstrapNode, err := bootstrap.New(ctx, listenAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create bootstrap node: %v\n", err)
		os.Exit(1)
	}

	if err := bootstrapNode.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start bootstrap node: %v\n", err)
		os.Exit(1)
	}

	select {}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/topology-gg/gram/bootstrap"
)

func main() {
	ctx := context.Background()

	port := flag.String("p", "4001", "Port number")

	flag.Parse()

	fmt.Println(*port)

	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", *port)
	fmt.Printf("Configured to listen on %s\n", listenAddr)

	bootstrapNode, err := bootstrap.New(ctx, listenAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create bootstrap node: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Bootstrap node created successfully.")
	}

	if err := bootstrapNode.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start bootstrap node: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Bootstrap node started successfully.")
	}

	select {}
}

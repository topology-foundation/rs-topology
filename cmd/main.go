package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/topology-gg/gram/config"
	"github.com/topology-gg/gram/execution"
	"github.com/topology-gg/gram/internal/app"
	"github.com/topology-gg/gram/network"
	"github.com/topology-gg/gram/storage"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.json", "Path to configuration file")

	flag.Parse()

	app := app.NewApp()

	app.Name = "gram"
	app.Description = "The official Go implementation of the RAM network"
	app.Action = func() error {
		return gram(configPath)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func gram(configPath string) error {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	storage := storage.NewStorage(ctx, &cfg.Storage)
	execution := execution.NewExecution(ctx, storage, &cfg.Execution)
	network := network.NewNetwork(ctx, execution, storage, &cfg.Network)

	network.Start()

	return nil
}

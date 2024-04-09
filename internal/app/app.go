package app

import (
	"context"
	"fmt"
	"os"

	"github.com/topology-gg/gram/config"
	"github.com/topology-gg/gram/execution"
	"github.com/topology-gg/gram/network"
	"github.com/topology-gg/gram/storage"
	"github.com/topology-gg/gram/util"
)

// Gram configures initial values and bootstraps the project
func Gram() {
	fmt.Println("Starting gram node")

	ctx := context.Background()

	// load configuration from file
	cfg, err := config.LoadConfig()
	logErrorAndPanic(err)

	ch := make(chan error)

	// instantiate modules
	storage, err := storage.NewStorage(ctx, &cfg.Storage)
	logErrorAndPanic(err)

	execution, err := execution.NewExecution(ctx, storage, &cfg.Execution)
	logErrorAndPanic(err)

	network, err := network.NewNetwork(ctx, ch, execution, storage, &cfg.Network)
	logErrorAndPanic(err)

	// run modules
	go network.Start()
	go util.ListenSigint(ch)

	// if we are getting an error from one of the modules
	// we should shutdown everything properly, log received error and terminate the app
	moduleErr := <-ch
	fmt.Fprintln(os.Stderr, moduleErr)

	network.Shutdown()

	if err := storage.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println("Shutting down gram")
}

// This is a private fucntion that is used only during app setup
// For any other case the app should never panic but handle errors
func logErrorAndPanic(err error) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err)
	panic(err)
}

package execution

import (
	"context"
	"fmt"
	"strings"

	"github.com/topology-gg/gram/config"
	"github.com/topology-gg/gram/storage"
)

type ExecutionModule struct {
	ctx     context.Context
	storage storage.Storage
	config  *config.ExecutionConfig
}

func NewExecution(ctx context.Context, storage storage.Storage, config *config.ExecutionConfig) (*ExecutionModule, error) {
	return &ExecutionModule{
		ctx:     ctx,
		storage: storage,
		config:  config,
	}, nil
}

func (execution *ExecutionModule) Execute(message string) {
	// TODO: Proper message handling comes here.

	fmt.Printf("(Execution) %s", message)
	kv := strings.Split(message, ": ")
	_ = execution.storage.Set([]byte(kv[0]), []byte(kv[1]))
}

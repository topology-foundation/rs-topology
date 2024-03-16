package storage

import (
	"context"
	"fmt"

	"github.com/topology-gg/gram/config"
)

type StorageModule struct {
	ctx    context.Context
	config config.StorageConfig
}

func NewStorage(ctx context.Context, config *config.StorageConfig) *StorageModule {
	return &StorageModule{
		ctx:    ctx,
		config: *config,
	}
}

func (storage *StorageModule) Set(key, value []byte) error {
	// TODO: Store the key-value pair.

	fmt.Printf("(Storage) %s: %s", key, value)
	return nil
}

package storage

import (
	"context"
	"fmt"
	"slices"

	"github.com/topology-gg/gram/config"

	"github.com/cockroachdb/pebble"
)

type StorageModule struct {
	ctx    context.Context
	config config.StorageConfig
	db     *pebble.DB
}

func NewStorage(ctx context.Context, config *config.StorageConfig) (*StorageModule, error) {
	db, err := pebble.Open(config.DatabasePath, &pebble.Options{})
	if err != nil {
		return nil, err
	}

	return &StorageModule{
		ctx:    ctx,
		config: *config,
		db:     db,
	}, nil
}

func (storage *StorageModule) Has(key []byte) (bool, error) {
	_, closer, err := storage.db.Get(key)
	if err != nil {
		return false, err
	}
	defer closer.Close()

	return true, nil
}

func (storage *StorageModule) Get(key []byte) ([]byte, error) {
	value, closer, err := storage.db.Get(key)
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	// Return a copy of `value` to guarantee its validity after closing the Closer.
	return slices.Clone(value), nil
}

func (storage *StorageModule) Set(key, value []byte) error {
	fmt.Printf("(Storage) %s: %s", key, value)

	return storage.db.Set(key, value, pebble.Sync)
}

func (storage *StorageModule) Delete(key []byte) error {
	return storage.db.Delete(key, pebble.Sync)
}

func (storage *StorageModule) Close() error {
	if err := storage.db.Close(); err != nil {
		return err
	}

	fmt.Println("DB connection successfully closed")
	return nil
}

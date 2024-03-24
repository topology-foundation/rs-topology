package config

import (
	"encoding/json"
	"os"
)

type ExecutionConfig struct {
}

type NetworkConfig struct {
	Namespace string   `json:"namespace"`
	Topics    []string `json:"topics"`
	MaxPeers  int      `json:"maxPeers"`
	Port      int      `json:"port"`
}

type StorageConfig struct {
}

type AppConfig struct {
	Execution ExecutionConfig `json:"executionConfig"`
	Network   NetworkConfig   `json:"networkConfig"`
	Storage   StorageConfig   `json:"storageConfig"`
}

func DefaultExecutionConfig() *ExecutionConfig {
	return &ExecutionConfig{}
}

func DefaultNetworkConfig() *NetworkConfig {
	return &NetworkConfig{
		Namespace: "gram-namespace",
		Topics:    []string{"gram-topic"},
		MaxPeers:  1,
		Port:      1319,
	}
}

func DefaultStorageConfig() *StorageConfig {
	return &StorageConfig{}
}

func LoadConfig(path string) (*AppConfig, error) {
	cfg := &AppConfig{
		Execution: *DefaultExecutionConfig(),
		Network:   *DefaultNetworkConfig(),
		Storage:   *DefaultStorageConfig(),
	}

	file, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

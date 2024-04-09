package config

import (
	"encoding/json"
	"flag"
	"os"
)

var (
	configPath = flag.String("config", "config.json", "The path to the configuration file")
)

type ExecutionConfig struct {
}

type NetworkConfig struct {
	P2p  P2pConfig  `json:"p2pConfig"`
	Grpc GrpcConfig `json:"grpcConfig"`
	Rpc  RpcConfig  `json:"rpcConfig"`
}

type StorageConfig struct {
	DatabasePath string `json:"databasePath"`
}

type P2pConfig struct {
	Namespace string   `json:"namespace"`
	Topics    []string `json:"topics"`
	MaxPeers  int      `json:"maxPeers"`
	Port      int      `json:"port"`
}

type GrpcConfig struct {
	Port int `json:"port"`
}

type RpcConfig struct {
	Port int `json:"port"`
}

type AppConfig struct {
	Execution ExecutionConfig `json:"executionConfig"`
	Network   NetworkConfig   `json:"networkConfig"`
	Storage   StorageConfig   `json:"storageConfig"`
}

func LoadConfig[T any]() (*T, error) {
	flag.Parse()

	file, err := os.Open(*configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg T
	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

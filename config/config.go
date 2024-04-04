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

type P2pConfig struct {
	Namespace string   `json:"namespace"`
	Topics    []string `json:"topics"`
	MaxPeers  int      `json:"maxPeers"`
	Port      int      `json:"port"`
}

type StorageConfig struct {
	DatabasePath string `json:"databasePath"`
}

type GrpcConfig struct {
	Port int `json:"port"`
}

type RpcConfig struct {
	Port int `json:"port"`
}

type AppConfig struct {
	Execution ExecutionConfig `json:"executionConfig"`
	Storage   StorageConfig   `json:"storageConfig"`
	Network   NetworkConfig   `json:"networkConfig"`
}

func LoadConfig() (*AppConfig, error) {
	flag.Parse()

	cfg := &AppConfig{}

	file, err := os.Open(*configPath)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

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

type LogConfig struct {
	LogLevel string `json:"logLevel"`
}

type AppConfig struct {
	Execution ExecutionConfig `json:"executionConfig"`
	Network   NetworkConfig   `json:"networkConfig"`
	Storage   StorageConfig   `json:"storageConfig"`
	Log       LogConfig       `json:"logConfig"`
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

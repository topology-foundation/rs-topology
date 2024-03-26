package config

import (
	"encoding/json"
	"flag"
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

type ConfigOverrides struct {
	Port     int
	MaxPeers int
}

func LoadConfig() (*AppConfig, error) {
	configPath, overrides := parseFlags()

	cfg := &AppConfig{}

	file, err := os.Open(configPath)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(cfg); err != nil {
		return cfg, err
	}

	applyOverrides(cfg, overrides)

	return cfg, nil
}

func parseFlags() (string, ConfigOverrides) {
	var configPath string
	var overrides ConfigOverrides

	flag.StringVar(&configPath, "config", "./config.json", "Path to configuration file")
	flag.StringVar(&configPath, "c", "./config.json", "Path to configuration file (shorthand)")

	flag.IntVar(&overrides.Port, "port", 0, "Port to override the one in the configuration file")
	flag.IntVar(&overrides.Port, "p", 0, "Port to override the one in the configuration file")

	flag.IntVar(&overrides.MaxPeers, "peers", 0, "MaxPeers to override the one in the configuration file")

	flag.Parse()

	return configPath, overrides
}

func applyOverrides(cfg *AppConfig, overrides ConfigOverrides) {
	if overrides.Port > 0 {
		cfg.Network.Port = overrides.Port
	}
	if overrides.MaxPeers > 0 {
		cfg.Network.MaxPeers = overrides.MaxPeers
	}
}

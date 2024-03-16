package config

type ExecutionConfig struct {
}

type NetworkConfig struct {
	Namespace string
	Topics    []string
	MaxPeers  int
}

type StorageConfig struct {
}

func DefaultExecutionConfig() *ExecutionConfig {
	return &ExecutionConfig{}
}

func DefaultNetworkConfig() *NetworkConfig {
	return &NetworkConfig{
		Namespace: "gram-namespace",
		Topics:    []string{"gram-topic"},
		MaxPeers:  1,
	}
}

func DefaultStorageConfig() *StorageConfig {
	return &StorageConfig{}
}

package config

import (
	"sync"
)

// AppConfig keeps configurations of application.
type AppConfig struct {
	RunAddr         string `json:"address"`
	GRPCRunAddr     string `json:"grpc_address"`
	Key             string
	Session         string
	LogLevel        string
	MigrationsPath  string
	TrustedSubnet   string
	PrivateKeyPath  string `json:"crypto_key"`
	FileStoragePath string `json:"store_file"`
	DatabaseDsn     string `json:"database_dsn"`
	ConfigPath      string
	StoreInterval   int  `json:"store_interval"`
	Restore         bool `json:"restore"`
	EnableGRPC      bool `json:"enable_grpc"`

	mux sync.Mutex
}

// NewConfig creates new configuration.
func NewConfig() *AppConfig {
	c := AppConfig{
		mux: sync.Mutex{},
	}
	c.initFlags()

	return &c
}

// ParseConfig parses existing configuration.
func (c *AppConfig) ParseConfig() {
	c.parseJSON()
	c.parseEnv()
	c.parseFlags()
}

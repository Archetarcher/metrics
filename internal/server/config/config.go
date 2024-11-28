package config

import (
	"sync"
)

// AppConfig keeps configurations of application.
type AppConfig struct {
	RunAddr         string `json:"address"`
	Key             string
	Session         string
	LogLevel        string
	MigrationsPath  string
	PrivateKeyPath  string `json:"crypto_key"`
	FileStoragePath string `json:"store_file"`
	DatabaseDsn     string `json:"database_dsn"`
	StoreInterval   int    `json:"store_interval"`
	ConfigPath      string
	Restore         bool `json:"restore"`

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
	c.parseFlags()
	c.parseEnv()
	c.parseJson()
}

package config

import (
	"sync"

	"github.com/Archetarcher/metrics.git/internal/server/store"
)

// AppConfig keeps configurations of application.
type AppConfig struct {
	Store    store.Config
	RunAddr  string
	Key      string
	LogLevel string
	mux      sync.Mutex
}

// NewConfig creates new configuration.
func NewConfig(store store.Config) *AppConfig {
	c := AppConfig{
		mux:   sync.Mutex{},
		Store: store,
	}
	c.initFlags()

	return &c
}

// ParseConfig parses existing configuration.
func (c *AppConfig) ParseConfig() {
	c.parseFlags()
	c.parseEnv()
}

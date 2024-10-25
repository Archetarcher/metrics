package config

import (
	"sync"

	"github.com/Archetarcher/metrics.git/internal/server/store"
)

// AppConfig keeps configurations of application.
type AppConfig struct {
	mux      sync.Mutex
	RunAddr  string
	LogLevel string
	Key      string
	Store    store.Config
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

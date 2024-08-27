package config

import (
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"sync"
)

func (c *AppConfig) ParseConfig() {
	c.parseFlags()
	c.parseEnv()
}

type AppConfig struct {
	mux      sync.Mutex
	RunAddr  string
	LogLevel string
	Key      string
	Store    store.Config
}

func NewConfig(store store.Config) *AppConfig {
	var c = AppConfig{
		mux:   sync.Mutex{},
		Store: store,
	}
	c.initFlags()

	return &c
}

package config

import "github.com/Archetarcher/metrics.git/internal/server/store"

func (c *AppConfig) ParseConfig() {
	c.parseFlags()
	c.parseEnv()
}

type AppConfig struct {
	RunAddr  string
	LogLevel string
	Store    store.Config
}

func NewConfig(store store.Config) *AppConfig {
	c := AppConfig{
		Store: store,
	}
	c.initFlags()

	return &c
}

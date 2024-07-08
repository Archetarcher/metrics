package config

func (c *AppConfig) ParseConfig() {
	c.parseFlags()
	c.parseEnv()
}

type AppConfig struct {
	RunAddr         string
	LogLevel        string
	FileStoragePath string
	StoreInterval   int
	Restore         bool
}

func NewConfig() *AppConfig {
	var c AppConfig
	c.initFlags()

	return &c
}

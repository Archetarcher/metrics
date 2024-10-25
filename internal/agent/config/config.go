package config

// AppConfig keeps configurations of application.
type AppConfig struct {
	ServerRunAddr  string
	ReportInterval int
	PollInterval   int
	LogLevel       string
	RateLimit      int
	Key            string
}

// NewConfig creates new configuration.
func NewConfig() *AppConfig {
	var c AppConfig
	c.initFlags()

	return &c
}

// ParseConfig parses existing configuration.
func (c *AppConfig) ParseConfig() {
	c.parseFlags()
	c.parseEnv()
}

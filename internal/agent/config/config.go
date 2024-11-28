package config

// AppConfig keeps configurations of application.
type AppConfig struct {
	ServerRunAddr  string
	LogLevel       string
	Key            string
	PublicKeyPath  string
	Session        Session
	ReportInterval int
	PollInterval   int
	RateLimit      int
}

type Session struct {
	Key       string
	RetryConn int
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

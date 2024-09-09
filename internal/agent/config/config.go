package config

func (c *AppConfig) ParseConfig() {
	c.parseFlags()
	c.parseEnv()
}

type AppConfig struct {
	ServerRunAddr  string
	ReportInterval int
	PollInterval   int
	LogLevel       string
	RateLimit      int
	Key            string
}

func NewConfig() *AppConfig {
	var c AppConfig
	c.initFlags()

	return &c
}

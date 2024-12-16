package config

// AppConfig keeps configurations of application.
type AppConfig struct {
	ServerRunAddr  string `json:"address"`
	GRPCRunAddr    string `json:"grpc_address"`
	LogLevel       string
	Key            string
	PublicKeyPath  string `json:"crypto_key"`
	ConfigPath     string
	Session        Session
	ReportInterval int `json:"report_interval"`
	PollInterval   int `json:"poll_interval"`
	RateLimit      int
	EnableGRPC     bool `json:"enable_grpc"`
}

// Session keeps session data
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
	c.parseJSON()
	c.parseEnv()
	c.parseFlags()

}

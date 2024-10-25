package config

import (
	"flag"
)

const (
	flagServerRunAddrName  = "a"
	flagReportIntervalName = "r"
	flagPollIntervalName   = "p"
	flagLogLevelName       = "l"
	flagKeyName            = "k"
	flagRateLimitName      = "rl"
)

func (c *AppConfig) initFlags() {
	flag.StringVar(&c.ServerRunAddr, flagServerRunAddrName, "localhost:8080", "address and port where server is running")
	flag.IntVar(&c.ReportInterval, flagReportIntervalName, 10, "interval in seconds for report to server")
	flag.IntVar(&c.PollInterval, flagPollIntervalName, 2, "interval in seconds for poll ")
	flag.StringVar(&c.LogLevel, flagLogLevelName, "info", "log level")
	flag.StringVar(&c.Key, flagKeyName, "", "key")
	flag.IntVar(&c.RateLimit, flagRateLimitName, 3, "rate limit")
}

func (c *AppConfig) parseFlags() {
	flag.Parse()
}

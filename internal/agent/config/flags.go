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
	flagPublicKeyPathName  = "crypto-key"
	flagSessionRetryConn   = "rc"
	flagConfigPathName     = "c config"
	flagGRPCRunAddrName    = "ga"
	flagEnableGRPCName     = "eg"
)

func (c *AppConfig) initFlags() {
	flag.StringVar(&c.ServerRunAddr, flagServerRunAddrName, "localhost:8080", "address and port where server is running")
	flag.IntVar(&c.ReportInterval, flagReportIntervalName, 3, "interval in seconds for report to server")
	flag.IntVar(&c.PollInterval, flagPollIntervalName, 1, "interval in seconds for poll ")
	flag.StringVar(&c.LogLevel, flagLogLevelName, "info", "log level")
	flag.StringVar(&c.Key, flagKeyName, "", "key")
	flag.StringVar(&c.PublicKeyPath, flagPublicKeyPathName, "public.pem", "crypto-key")
	flag.StringVar(&c.ConfigPath, flagConfigPathName, "agent-config.json", "config file")

	flag.IntVar(&c.Session.RetryConn, flagSessionRetryConn, 5, "connection retry count")
	flag.IntVar(&c.RateLimit, flagRateLimitName, 3, "rate limit")

	flag.BoolVar(&c.EnableGRPC, flagEnableGRPCName, false, "run grpc server or not")
	flag.StringVar(&c.GRPCRunAddr, flagGRPCRunAddrName, ":3200", "address and port to run grpc server")

}

func (c *AppConfig) parseFlags() {
	flag.Parse()
}

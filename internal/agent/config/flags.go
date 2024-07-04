package config

import (
	"flag"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
)

const (
	flagServerRunAddrName  = "a"
	flagReportIntervalName = "r"
	flagPollIntervalName   = "p"
	flagLogLevelName       = "l"
)

func parseFlags() {
	flag.StringVar(&domain.ServerRunAddr, flagServerRunAddrName, "localhost:8080", "address and port where server is running")
	flag.IntVar(&domain.ReportInterval, flagReportIntervalName, 10, "interval in seconds for report to server")
	flag.IntVar(&domain.PollInterval, flagPollIntervalName, 2, "interval in seconds for poll ")
	flag.StringVar(&domain.LogLevel, flagLogLevelName, "info", "log level")
	flag.Parse()
}

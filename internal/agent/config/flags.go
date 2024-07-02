package config

import (
	"flag"
	"github.com/Archetarcher/metrics.git/internal/agent/models"
)

const (
	flagServerRunAddrName  = "a"
	flagReportIntervalName = "r"
	flagPollIntervalName   = "p"
	flagLogLevelName       = "l"
)

func parseFlags() {
	flag.StringVar(&models.ServerRunAddr, flagServerRunAddrName, "localhost:8080", "address and port where server is running")
	flag.IntVar(&models.ReportInterval, flagReportIntervalName, 10, "interval in seconds for report to server")
	flag.IntVar(&models.PollInterval, flagPollIntervalName, 2, "interval in seconds for poll ")
	flag.StringVar(&models.LogLevel, flagLogLevelName, "info", "log level")

	flag.Parse()
}

package handlers

import (
	"flag"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
)

const (
	flagServerRunAddrName  = "a"
	flagReportIntervalName = "r"
	flagPollIntervalName   = "p"
)

func parseFlags() {
	flag.StringVar(&domain.FlagServerRunAddr, flagServerRunAddrName, "localhost:8080", "address and port where server is running")
	flag.IntVar(&domain.ReportInterval, flagReportIntervalName, 10, "interval in seconds for report to server")
	flag.IntVar(&domain.PollInterval, flagPollIntervalName, 2, "interval in seconds for poll ")
	flag.Parse()

}

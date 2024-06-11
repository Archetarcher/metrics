package handlers

import (
	"flag"
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"os"
	"strconv"
)

const (
	flagServerRunAddrName  = "a"
	flagReportIntervalName = "r"
	flagPollIntervalName   = "p"

	envServerRunAddrName  = "ADDRESS"
	envReportIntervalName = "REPORT_INTERVAL"
	envPollIntervalName   = "p"
)

func parseFlags() {
	flag.StringVar(&domain.ServerRunAddr, flagServerRunAddrName, "localhost:8080", "address and port where server is running")
	flag.IntVar(&domain.ReportInterval, flagReportIntervalName, 10, "interval in seconds for report to server")
	flag.IntVar(&domain.PollInterval, flagPollIntervalName, 2, "interval in seconds for poll ")
	flag.Parse()

	if envRunAddr := os.Getenv(envServerRunAddrName); envRunAddr != "" {
		domain.ServerRunAddr = envRunAddr
	}

	if envReportInt := os.Getenv(envReportIntervalName); envReportInt != "" {

		if i, err := strconv.Atoi(envReportInt); err == nil {
			domain.ReportInterval = i

		}
	}
	if envPollInt := os.Getenv(envPollIntervalName); envPollInt != "" {

		if i, err := strconv.Atoi(envPollInt); err == nil {
			domain.PollInterval = i

		}
	}

}

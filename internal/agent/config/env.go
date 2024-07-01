package config

import (
	"github.com/Archetarcher/metrics.git/internal/agent/domain"
	"os"
	"strconv"
)

const (
	envServerRunAddrName  = "ADDRESS"
	envReportIntervalName = "REPORT_INTERVAL"
	envPollIntervalName   = "p"
)

func parseEnv() {
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

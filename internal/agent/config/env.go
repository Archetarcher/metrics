package config

import (
	"github.com/Archetarcher/metrics.git/internal/agent/models"
	"os"
	"strconv"
)

const (
	envServerRunAddrName  = "ADDRESS"
	envReportIntervalName = "REPORT_INTERVAL"
	envPollIntervalName   = "p"
	envLogLevelName       = "LOG_LEVEL"
)

func parseEnv() {
	if envRunAddr := os.Getenv(envServerRunAddrName); envRunAddr != "" {
		models.ServerRunAddr = envRunAddr
	}

	if envReportInt := os.Getenv(envReportIntervalName); envReportInt != "" {

		if i, err := strconv.Atoi(envReportInt); err == nil {
			models.ReportInterval = i

		}
	}
	if envPollInt := os.Getenv(envPollIntervalName); envPollInt != "" {

		if i, err := strconv.Atoi(envPollInt); err == nil {
			models.PollInterval = i

		}
	}
	if envLogLevel := os.Getenv(envLogLevelName); envLogLevel != "" {
		models.LogLevel = envLogLevel
	}
}

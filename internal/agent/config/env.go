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
	envLogLevelName       = "LOG_LEVEL"
)

func getEnvOrDefault(env string, def string) string {
	val := os.Getenv(env)
	if val == "" {
		return def
	}
	return val
}

func parseEnv() {
	domain.ServerRunAddr = getEnvOrDefault(envServerRunAddrName, domain.ServerRunAddr)
	domain.LogLevel = getEnvOrDefault(envLogLevelName, domain.LogLevel)

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

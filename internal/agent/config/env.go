package config

import (
	"os"
	"strconv"
)

const (
	envServerRunAddrName  = "ADDRESS"
	envReportIntervalName = "REPORT_INTERVAL"
	envPollIntervalName   = "p"
	envLogLevelName       = "LOG_LEVEL"
	envKeyName            = "KEY"
	envRateLimitName      = "RATE_LIMIT"
)

func getEnvOrDefault(env string, def any, t int) any {
	val := os.Getenv(env)
	if val == "" {
		return def
	}

	switch t {
	case 1:
		return val
	case 2:
		if i, err := strconv.Atoi(val); err == nil {
			return i

		}
		return def
	case 3:
		if i, err := strconv.ParseBool(val); err == nil {
			return i
		}
		return def
	default:
		return def
	}
}

func (c *AppConfig) parseEnv() {
	c.ServerRunAddr = getEnvOrDefault(envServerRunAddrName, c.ServerRunAddr, 1).(string)
	c.LogLevel = getEnvOrDefault(envLogLevelName, c.LogLevel, 1).(string)
	c.ReportInterval = getEnvOrDefault(envReportIntervalName, c.ReportInterval, 2).(int)
	c.PollInterval = getEnvOrDefault(envPollIntervalName, c.PollInterval, 2).(int)
	c.Key = getEnvOrDefault(envKeyName, c.Key, 1).(string)
	c.RateLimit = getEnvOrDefault(envRateLimitName, c.RateLimit, 2).(int)
}

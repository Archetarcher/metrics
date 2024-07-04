package config

import (
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"os"
	"strconv"
)

const (
	envRunAddrName         = "ADDRESS"
	envLogLevelName        = "LOG_LEVEL"
	envFileStoragePathName = "FILE_STORAGE_PATH"
	envStoreIntervalName   = "STORE_INTERVAL"
	envRestoreName         = "RESTORE"
)

func getEnvOrDefault(env string, def string) string {
	val := os.Getenv(env)
	if val == "" {
		return def
	}
	return val
}

func parseEnv() {
	domain.RunAddr = getEnvOrDefault(envRunAddrName, domain.RunAddr)
	domain.LogLevel = getEnvOrDefault(envLogLevelName, domain.LogLevel)
	domain.FileStoragePath = getEnvOrDefault(envFileStoragePathName, domain.FileStoragePath)

	if envStoreInterval := os.Getenv(envStoreIntervalName); envStoreInterval != "" {
		if i, err := strconv.Atoi(envStoreInterval); err == nil {
			domain.StoreInterval = i

		}
	}
	if envRestore := os.Getenv(envRestoreName); envRestore != "" {
		if i, err := strconv.ParseBool(envRestore); err == nil {
			domain.Restore = i
		}
	}
}

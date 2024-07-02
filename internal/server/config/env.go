package config

import (
	"github.com/Archetarcher/metrics.git/internal/server/models"
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

func parseEnv() {
	if envRunAddr := os.Getenv(envRunAddrName); envRunAddr != "" {
		models.RunAddr = envRunAddr
	}
	if envLogLevel := os.Getenv(envLogLevelName); envLogLevel != "" {
		models.LogLevel = envLogLevel
	}
	if envFileStoragePath := os.Getenv(envFileStoragePathName); envFileStoragePath != "" {
		models.FileStoragePath = envFileStoragePath
	}
	if envStoreInterval := os.Getenv(envStoreIntervalName); envStoreInterval != "" {
		if i, err := strconv.Atoi(envStoreInterval); err == nil {
			models.StoreInterval = i

		}
	}
	if envRestore := os.Getenv(envRestoreName); envRestore != "" {
		if i, err := strconv.ParseBool(envRestore); err == nil {
			models.Restore = i
		}

	}
}

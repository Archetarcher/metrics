package config

import (
	"os"
	"strconv"
)

const (
	envRunAddrName                = "ADDRESS"
	envLogLevelName               = "LOG_LEVEL"
	envKeyName                    = "KEY"
	envFileStoragePathName        = "FILE_STORAGE_PATH"
	envStoreIntervalName          = "STORE_INTERVAL"
	envRestoreName                = "RESTORE"
	envDatabaseDsnName            = "DATABASE_DSN"
	envDatabaseMigrationsPathName = "DATABASE_MIGRATIONS_PATH"
	envPrivateKeyPathName         = "CRYPTO_KEY"
	envConfigPathName             = "CONFIG"
)

func (c *AppConfig) parseEnv() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.RunAddr = getEnvOrDefault(envRunAddrName, c.RunAddr, 1).(string)
	c.LogLevel = getEnvOrDefault(envLogLevelName, c.LogLevel, 1).(string)
	c.Key = getEnvOrDefault(envKeyName, c.Key, 1).(string)

	c.FileStoragePath = getEnvOrDefault(envFileStoragePathName, c.FileStoragePath, 1).(string)
	c.StoreInterval = getEnvOrDefault(envStoreIntervalName, c.StoreInterval, 2).(int)
	c.Restore = getEnvOrDefault(envRestoreName, c.Restore, 3).(bool)

	c.DatabaseDsn = getEnvOrDefault(envDatabaseDsnName, c.DatabaseDsn, 1).(string)
	c.MigrationsPath = getEnvOrDefault(envDatabaseMigrationsPathName, c.MigrationsPath, 1).(string)

	c.PrivateKeyPath = getEnvOrDefault(envPrivateKeyPathName, c.PrivateKeyPath, 1).(string)
	c.ConfigPath = getEnvOrDefault(envConfigPathName, c.ConfigPath, 1).(string)

}

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

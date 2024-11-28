package config

import (
	"flag"
)

const (
	flagRunAddrName                = "a"
	flagLogLevelName               = "l"
	flagKeyName                    = "k"
	flagFileStoragePathName        = "f"
	flagStoreIntervalName          = "i"
	flagRestoreName                = "r"
	flagDatabaseDsnName            = "d"
	flagDatabaseMigrationsPathName = "m"
	flagPrivateKeyPathName         = "crypto-key"
	flagConfigPathName             = "c config"
)

func (c *AppConfig) initFlags() {
	c.mux.Lock()
	defer c.mux.Unlock()

	flag.StringVar(&c.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.StringVar(&c.LogLevel, flagLogLevelName, "info", "log level")
	flag.StringVar(&c.Key, flagKeyName, "", "key")

	flag.StringVar(&c.FileStoragePath, flagFileStoragePathName, "/tmp/metrics-pgx.json", "file storage path")
	flag.IntVar(&c.StoreInterval, flagStoreIntervalName, 300, "seconds to save data to file")
	flag.BoolVar(&c.Restore, flagRestoreName, true, "load data from file")

	flag.StringVar(&c.DatabaseDsn, flagDatabaseDsnName, "", "dsn")
	flag.StringVar(&c.MigrationsPath, flagDatabaseMigrationsPathName, "internal/server/migrations", "migrations")

	flag.StringVar(&c.PrivateKeyPath, flagPrivateKeyPathName, "private.pem", "crypto-key")
	flag.StringVar(&c.ConfigPath, flagConfigPathName, "server-config.json", "config file")

}

func (c *AppConfig) parseFlags() {
	flag.Parse()
}

package config

import (
	"flag"
)

const (
	flagRunAddrName                = "a"
	flagLogLevelName               = "l"
	flagFileStoragePathName        = "f"
	flagStoreIntervalName          = "i"
	flagRestoreName                = "r"
	flagDatabaseDsnName            = "d"
	flagDatabaseMigrationsPathName = "m"
)

func (c *AppConfig) parseFlags() {
	flag.Parse()
}
func (c *AppConfig) initFlags() {

	flag.StringVar(&c.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.StringVar(&c.LogLevel, flagLogLevelName, "info", "log level")

	flag.StringVar(&c.Store.Memory.FileStoragePath, flagFileStoragePathName, "/tmp/metrics-pgx.json", "file storage path")
	flag.IntVar(&c.Store.Memory.StoreInterval, flagStoreIntervalName, 300, "seconds to save data to file")
	flag.BoolVar(&c.Store.Memory.Restore, flagRestoreName, true, "load data from file")

	flag.StringVar(&c.Store.Pgx.DatabaseDsn, flagDatabaseDsnName, "", "dsn")
	flag.StringVar(&c.Store.Pgx.MigrationsPath, flagDatabaseMigrationsPathName, "internal/server/migrations", "migrations")

}

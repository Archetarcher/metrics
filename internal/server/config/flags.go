package config

import (
	"flag"
)

const (
	flagRunAddrName         = "a"
	flagLogLevelName        = "l"
	flagFileStoragePathName = "f"
	flagStoreIntervalName   = "i"
	flagRestoreName         = "r"
)

func (c *AppConfig) parseFlags() {
	flag.Parse()
}
func (c *AppConfig) initFlags() {

	flag.StringVar(&c.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.StringVar(&c.LogLevel, flagLogLevelName, "info", "log level")
	flag.StringVar(&c.FileStoragePath, flagFileStoragePathName, "/tmp/metrics-db.json", "file storage path")
	flag.IntVar(&c.StoreInterval, flagStoreIntervalName, 300, "seconds to save data to file")
	flag.BoolVar(&c.Restore, flagRestoreName, true, "load data from file")

}

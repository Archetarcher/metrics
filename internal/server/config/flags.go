package config

import (
	"flag"
	"github.com/Archetarcher/metrics.git/internal/server/models"
)

const (
	flagRunAddrName         = "a"
	flagLogLevelName        = "l"
	flagFileStoragePathName = "f"
	flagStoreIntervalName   = "i"
	flagRestoreName         = "r"
)

func parseFlags() {
	flag.StringVar(&models.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.StringVar(&models.LogLevel, flagLogLevelName, "info", "log level")
	flag.StringVar(&models.FileStoragePath, flagFileStoragePathName, "/tmp/metrics-db.json", "file storage path")
	flag.IntVar(&models.StoreInterval, flagStoreIntervalName, 300, "seconds to save data to file")
	flag.BoolVar(&models.Restore, flagRestoreName, true, "load data from file")
	flag.Parse()
}

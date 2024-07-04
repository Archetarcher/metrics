package config

import (
	"flag"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
)

const (
	flagRunAddrName         = "a"
	flagLogLevelName        = "l"
	flagFileStoragePathName = "f"
	flagStoreIntervalName   = "i"
	flagRestoreName         = "r"
)

func parseFlags() {
	flag.StringVar(&domain.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.StringVar(&domain.LogLevel, flagLogLevelName, "info", "log level")
	flag.StringVar(&domain.FileStoragePath, flagFileStoragePathName, "/tmp/metrics-db.json", "file storage path")
	flag.IntVar(&domain.StoreInterval, flagStoreIntervalName, 300, "seconds to save data to file")
	flag.BoolVar(&domain.Restore, flagRestoreName, true, "load data from file")
	flag.Parse()
}

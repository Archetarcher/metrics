package config

import (
	"flag"
	"github.com/Archetarcher/metrics.git/internal/server/models"
)

const (
	flagRunAddrName  = "a"
	flagLogLevelName = "l"
)

func parseFlags() {
	flag.StringVar(&models.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.StringVar(&models.LogLevel, flagLogLevelName, "info", "log level")
	flag.Parse()
}

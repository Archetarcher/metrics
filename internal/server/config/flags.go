package config

import (
	"flag"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
)

const (
	flagRunAddrName  = "a"
	flagLogLevelName = "l"
)

func parseFlags() {
	flag.StringVar(&domain.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.StringVar(&domain.LogLevel, flagLogLevelName, "info", "log level")
	flag.Parse()
}

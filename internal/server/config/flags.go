package config

import (
	"flag"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
)

const (
	flagRunAddrName = "a"
)

func parseFlags() {
	flag.StringVar(&domain.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.Parse()
}

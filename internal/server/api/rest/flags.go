package rest

import (
	"flag"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"os"
)

const (
	flagRunAddrName = "a"
	envRunAddrName  = "ADDRESS"
)

func parseFlags() {
	flag.StringVar(&domain.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.Parse()

	if envRunAddr := os.Getenv(envRunAddrName); envRunAddr != "" {
		domain.RunAddr = envRunAddr
	}

}

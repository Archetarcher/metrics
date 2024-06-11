package rest

import (
	"flag"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
)

func parseFlags() {
	flag.StringVar(&domain.FlagRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
}

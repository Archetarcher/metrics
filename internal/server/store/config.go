package store

import (
	"github.com/Archetarcher/metrics.git/internal/server/store/memory"
	"github.com/Archetarcher/metrics.git/internal/server/store/pgx"
)

type Config struct {
	Memory *memory.Config
	Pgx    *pgx.Config
}

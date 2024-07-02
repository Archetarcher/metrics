package main

import (
	"github.com/Archetarcher/metrics.git/internal/server/api/rest"
	"github.com/Archetarcher/metrics.git/internal/server/config"
	"github.com/Archetarcher/metrics.git/internal/server/store"
)

func main() {
	config.ParseConfig()

	storage := store.NewStorage()

	api, err := rest.NewMetricAPI(storage)

	if err != nil {
		panic(err)
	}

	if err := api.Run(); err != nil {
		panic(err)
	}
}

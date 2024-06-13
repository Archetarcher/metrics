package main

import (
	"github.com/Archetarcher/metrics.git/internal/server/api/rest"
	"github.com/Archetarcher/metrics.git/internal/server/store"
)

func main() {
	storage := store.NewStorage()
	api := rest.NewAPI(storage)
	if err := api.Run(); err != nil {
		panic(err)
	}
}

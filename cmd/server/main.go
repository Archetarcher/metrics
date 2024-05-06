package main

import (
	"github.com/Archetarcher/metrics.git/internal/server/api/rest"
)

func main() {

	api := rest.NewApi()
	if err := api.Run(); err != nil {
		panic(err)
	}
}

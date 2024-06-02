package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"github.com/Archetarcher/metrics.git/internal/server/store"
	"net/http"
)

type API struct {
	server *http.ServeMux
}

func NewAPI(storage *store.MemStorage) API {

	mux := http.NewServeMux()
	handler := handlers.MetricsHandler{Storage: storage}

	mux.HandleFunc(`/update/{type}/{name}/{value}`, handler.Update)

	return API{
		server: mux,
	}
}

func (a API) Run() error {
	return http.ListenAndServe(`:8080`, a.server)
}

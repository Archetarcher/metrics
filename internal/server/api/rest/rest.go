package rest

import (
	"github.com/Archetarcher/metrics.git/internal/server/handlers"
	"net/http"
)

type Api struct {
	server *http.ServeMux
}

func NewApi() Api {

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/{type}/{name}/{value}`, handlers.UpdateMetrics)

	return Api{
		server: mux,
	}
}

func (a Api) Run() error {
	return http.ListenAndServe(`:8080`, a.server)
}

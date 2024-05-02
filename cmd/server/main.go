package main

import (
	"net/http"
)

const (
	gauge   = "gauge"
	counter = "counter"
)

func main() {

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/{type}/{name}/{value}`, handler)

	return http.ListenAndServe(`:8080`, mux)
}

func handler(w http.ResponseWriter, r *http.Request) {

}

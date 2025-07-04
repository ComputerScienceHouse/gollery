package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"CSH/gollery/internal/endpoints/v1"
)

func main() {
	println("Starting server on :8000")

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	v1.RegisterAPIRoutes(api)

	// Optionally, a root handler
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the API"))
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}

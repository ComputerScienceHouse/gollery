package v1

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterFileRoutes(router *mux.Router) {
	router.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("File endpoint"))
	}).Methods("GET")
}

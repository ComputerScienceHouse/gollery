package v1

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterDirectoryRoutes(router *mux.Router) {
	router.HandleFunc("/directory", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Directory endpoint"))
	}).Methods("GET")
}

package v1

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterDownloadRoutes(router *mux.Router) {
	router.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Download endpoint"))
	}).Methods("GET")
}

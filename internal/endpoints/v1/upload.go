package v1

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterUploadRoutes(router *mux.Router) {
	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Upload endpoint"))
	}).Methods("POST")
}

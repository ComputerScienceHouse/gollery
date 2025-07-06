package v1

import (
	"encoding/json"
	"log"
	"net/http"

	srv "github.com/ComputerScienceHouse/gollery/internal/services"
	"github.com/gorilla/mux"
)

type directory struct {
	Did        int    `json:"id"`
	Name       string `json:"name"`
	Creator    string `json:"creator"`
	CreateDate string `json:"string"`
}

func RegisterDirectoryRoutes(router *mux.Router) {
	router.HandleFunc("/directory", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Directory endpoint"))
	}).Methods("GET")

	router.HandleFunc("/directory/{did}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		did := vars["did"]
		if err := srv.DB.QueryRow("SELECT name, creator, creat_date FROM directory WHERE did = $1", did); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error inserting book into books table %v", err)
			return
		}
	}).Methods("GET")

	router.HandleFunc("/directory/create", func(w http.ResponseWriter, r *http.Request) {
		var body directory
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error decoding request body into CreateBookBody struct %v", err)
			return
		}

		if err := srv.DB.QueryRow("INSERT INTO directory (name, creator, create_date) VALUES ($1, $2, $3)", body.Did, body.Creator, body.CreateDate); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error inserting book into books table %v", err)
			return
		}
	}).Methods("POST")
}

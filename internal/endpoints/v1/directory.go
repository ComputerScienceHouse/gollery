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

type maskedDirectory struct {
	Name       string `json:"name"`
	Creator    string `json:"creator"`
	CreateDate string `json:"string"`
}

func RegisterDirectoryRoutes(router *mux.Router) {
	router.HandleFunc("/directory", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Directory endpoint"))
	}).Methods("GET")

	router.HandleFunc("/directory/{did}", func(w http.ResponseWriter, r *http.Request) {
		var dir maskedDirectory
		vars := mux.Vars(r)
		did := vars["did"]
		log.Printf("vars: %v\n", vars)
		if err := srv.DB.QueryRow("SELECT name, creator, create_date FROM directory WHERE did = $1;", did).Scan(&dir.Name, &dir.Creator, &dir.CreateDate); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error selecting from the directory table %v", err)
			return
		}

		j, err := json.Marshal(dir)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error marshalling books into json %v", err)
			return
		}
		w.Write(j)

	}).Methods("GET")

	router.HandleFunc("/directory/create", func(w http.ResponseWriter, r *http.Request) {
		var body directory
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error decoding request body into CreateBookBody struct %v", err)
			return
		}

		sqlRes, err := srv.DB.Exec("INSERT INTO directory (name, creator) VALUES ($1, $2);", body.Did, body.Creator)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error inserting directory into directory table %v", err)
			return
		}
		log.Print(sqlRes)
		// var dir maskedDirectory
		// did, err := sqlRes.LastInsertId()
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	log.Printf("error inserting directory into directory table %v", err)
		// 	return
		// }
		// if err := srv.DB.QueryRow("SELECT name, creator, create_date FROM directory WHERE did = $1;", did).Scan(&dir.Name, &dir.Creator, &dir.CreateDate); err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	log.Printf("error selecting from the directory table %v", err)
		// 	return
		// }

		// j, err := json.Marshal(dir)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	log.Printf("error marshalling books into json %v", err)
		// 	return
		// }
		// w.Write(j)

	}).Methods("POST")
}

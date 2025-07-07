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

	//grabs all the directories
	router.HandleFunc("/directory/all", func(w http.ResponseWriter, r *http.Request) {
		rows, err := srv.DB.Query("SELECT name, creator, create_date FROM directory;")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error selecting from the directory table %v", err)
			return
		}
		defer rows.Close()

		//iterates through all of the rows that returned to put them in an array of maskedDirectory structs
		var dirs []maskedDirectory
		for rows.Next() {
			var dir maskedDirectory
			if err := rows.Scan(&dir.Name, &dir.Creator, &dir.CreateDate); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("error scanning directory row %v", err)
				return
			}
			dirs = append(dirs, dir)
		}
		if err := rows.Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error iterating directory rows %v", err)
			return
		}

		j, err := json.Marshal(dirs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error marshalling directories into json %v", err)
			return
		}
		log.Printf("j = %v\n", j)
		w.Write(j)
		w.Header().Set("Content-Type", "application/json")
	}).Methods("GET")

	//Returns the directory with the directory id did from within the http
	router.HandleFunc("/directory/{did:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		var dir maskedDirectory
		vars := mux.Vars(r)
		did := vars["did"]
		log.Printf("vars: %v\n", vars)
		dir = selectFromById(did, w)
		marshalIntoMaskedDirectory(dir, w)

	}).Methods("GET")

	//Creates directory with data from the body of the HTTP request
	router.HandleFunc("/directory/create", func(w http.ResponseWriter, r *http.Request) {
		var body directory
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error decoding request body into Directory struct %v", err)
			return
		}
		log.Printf("body: %v\n", body)

		if err := srv.DB.QueryRow("INSERT INTO directory (name, creator, create_date) VALUES ($1, $2, CURRENT_DATE);", body.Name, body.Creator).Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error inserting directory into directory table %v", err)
			return
		}
		//Can grab the highest value from the table which should be the most recently created thing
		//need to test where things are running at the same time to see if each query through go is atomic
		// 		 and so it could in theory return the wrong value if the requests are sent at the same time
		//TODO: Setup this query using things like BeginTx which allows you to prepare statements and execute them
		//  	 so you know that that block of the db is locked but idk if that actually fixes things
		var dir maskedDirectory
		if err := srv.DB.QueryRow(`
SELECT a.name, a.creator, a.create_date
FROM directory a
LEFT OUTER JOIN directory b
    ON a.did < b.did
WHERE b.did IS NULL;`).Scan(&dir.Name, &dir.Creator, &dir.CreateDate); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error inserting directory into directory table %v", err)
			return
		}

		marshalIntoMaskedDirectory(dir, w)
	}).Methods("POST")

	//Updates the directory with Directory id (did) did from within the HTTP request, updated values is within the body of the request
	router.HandleFunc("/directory/update/{did:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		var dir maskedDirectory
		vars := mux.Vars(r)
		did := vars["did"]
		log.Printf("vars: %v\n", vars)
		var body directory
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error decoding request body into Directory struct %v", err)
			return
		}

		//updating the directory based on the body of the request with the did as param
		if err := srv.DB.QueryRow("UPDATE directory SET name = $1, creator = $2, create_date = CURRENT_DATE WHERE did = $3", body.Name, body.Creator, did); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error updating the directory %v, %v, %v, %v, with error %v", body.Name, body.Creator, body.CreateDate, did, err)
			return
		}

		//returning the updated value below
		dir = selectFromById(did, w)
		marshalIntoMaskedDirectory(dir, w)

	}).Methods("PUT")

	//Deletes the directory with Directory id (did) did from within the HTTP request
	router.HandleFunc("/directory/{did:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		did := vars["did"]
		log.Print("Delete called")
		if err := srv.DB.QueryRow("DELETE from directory WHERE did = $1;", did); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error deleting from the directory table %v", err)
			return
		}
		// temp := 1
		// if err := srv.DB.QueryRow("SELECT COALESCE((SELECT COUNT(*) from directory WHERE did = $1), 0)", did).Scan(&temp); err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	log.Printf("error selecting from the directory table %v", err)
		// 	return
		// }
		// log.Printf("temp = %d", temp)
		// if temp != 0 {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	log.Printf("Error deleting from the directory table, record still exists, %v", temp)
		// 	return
		// }
		w.WriteHeader(http.StatusOK)
		log.Printf("Deleted row with ID %v", did)
	}).Methods("DELETE")

}

// selectFromById selects a maskedDirectory from a did, will likely be modularized later on once I understand go interfaces
func selectFromById(id string, w http.ResponseWriter) maskedDirectory {
	var dir maskedDirectory
	if err := srv.DB.QueryRow("SELECT name, creator, create_date FROM directory WHERE did = $1;", id).Scan(&dir.Name, &dir.Creator, &dir.CreateDate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error selecting from the directory table %v", err)
		return dir
	}
	return dir
}

// marchalIntoMaskedDirectory takes a maskedDirectory and puts it into a JSON object j and prints it to the response and the log
func marshalIntoMaskedDirectory(dir maskedDirectory, w http.ResponseWriter) {
	j, err := json.Marshal(dir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling directory into json %v", err)
		return
	}
	log.Printf("j = %v\n", j)
	w.Write(j)
}

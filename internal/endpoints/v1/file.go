package v1

import (
	"encoding/json"
	"log"
	"net/http"

	srv "github.com/ComputerScienceHouse/gollery/internal/services"
	"github.com/gorilla/mux"
)

type file struct {
	Fid         int    `json:"id"`
	Name        string `json:"name"`
	Creator     string `json:"creator"`
	Upload_time string `json:"upload_time"`
	S3_key      string `json:"s3_key"`
	Did         int    `json:"did"`
}
type maskedFile struct {
	Name        string `json:"name"`
	Creator     string `json:"creator"`
	Upload_time string `json:"upload_time"`
	S3_key      string `json:"s3_key"`
	Did         int    `json:"did"`
}

func RegisterFileRoutes(router *mux.Router) {
	router.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("File endpoint"))
	}).Methods("GET")

	//Returns the File with the directory id did from within the http
	router.HandleFunc("/file/{fid:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		var file maskedFile
		vars := mux.Vars(r)
		fid := vars["fid"]
		log.Printf("vars: %v\n", vars)
		file = selectFileById(fid, w)
		marshalIntoMaskedFile(file, w)
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	router.HandleFunc("/file/all", func(w http.ResponseWriter, r *http.Request) {
		rows, err := srv.DB.Query("SELECT name, creator, upload_time, s3_key, did FROM file;")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error selecting from the file table %v", err)
			return
		}
		defer rows.Close()

		//iterates through all of the rows that returned to put them in an array of maskedDirectory structs
		var files []maskedFile
		for rows.Next() {
			var file maskedFile
			if err := rows.Scan(&file.Creator, &file.Upload_time, &file.S3_key, &file.Did); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("error scanning file row %v", err)
				return
			}
			files = append(files, file)
		}
		if err := rows.Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error iterating file rows %v", err)
			return
		}

		j, err := json.Marshal(files)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error marshalling files into json %v", err)
			return
		}
		log.Printf("j = %v\n", j)
		w.Write(j)
		w.Header().Set("Content-Type", "application/json")
	}).Methods("GET")

	//Creates directory with data from the body of the HTTP request
	router.HandleFunc("/directory/create", func(w http.ResponseWriter, r *http.Request) {
		var body file
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error decoding request body into File struct %v", err)
			return
		}
		log.Printf("body: %v\n", body)
		// TODO: MAKE S3 LINK HERE AND INSERT
		if err := srv.DB.QueryRow("INSERT INTO directory (name, creator, upload_time, s3_key) VALUES ($1, $2, CURRENT_DATE, $3);", body.Name, body.Creator).Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error inserting directory into directory table %v", err)
			return
		}

		//Can grab the highest value from the table which should be the most recently created thing
		//need to test where things are running at the same time to see if each query through go is atomic
		// 		 and so it could in theory return the wrong value if the requests are sent at the same time
		//TODO: Setup this query using things like BeginTx which allows you to prepare statements and execute them
		//  	 so you know that that block of the db is locked but idk if that actually fixes things
		var file maskedFile
		if err := srv.DB.QueryRow(`
SELECT a.creator, a.upload_time, a.s3_key, a.did
FROM file a
LEFT OUTER JOIN file b
    ON a.did < b.did
WHERE b.did IS NULL;`).Scan(&file.Name, &file.Creator, &file.Upload_time, &file.S3_key, &file.Did); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error inserting directory into directory table %v", err)
			return
		}

		marshalIntoMaskedFile(file, w)
	}).Methods("POST")

} // end RegisterFileRoutes

// selectFromById selects a maskedDirectory from a did, will likely be modularized later on once I understand go interfaces
func selectFileById(id string, w http.ResponseWriter) maskedFile {
	var file maskedFile
	if err := srv.DB.QueryRow("SELECT name, creator, upload_time, s3_key, did FROM file WHERE fid = $1;", id).Scan(&file.Name, &file.Creator, &file.Upload_time, &file.S3_key, &file.Did); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error selecting from the directory table %v", err)
		return file
	}
	return file
}

// marchalIntoMaskedDirectory takes a maskedDirectory and puts it into a JSON object j and prints it to the response and the log
func marshalIntoMaskedFile(file maskedFile, w http.ResponseWriter) {
	j, err := json.Marshal(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error marshalling directory into json %v", err)
		return
	}
	log.Printf("j = %v\n", j)
	w.Write(j)
}

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
	Creator     string `json:"creator"`
	Upload_time string `json:"upload_time"`
	S3_key      string `json:"s3_key"`
	Did         int    `json:"did"`
}
type maskedFile struct {
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
} // end RegisterFileRoutes

// selectFromById selects a maskedDirectory from a did, will likely be modularized later on once I understand go interfaces
func selectFileById(id string, w http.ResponseWriter) maskedFile {
	var file maskedFile
	if err := srv.DB.QueryRow("SELECT creator, upload_time, s3_key, did FROM file WHERE fid = $1;", id).Scan(&file.Creator, &file.Upload_time, &file.S3_key, file.Did); err != nil {
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

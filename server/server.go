package server

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/ComputerScienceHouse/gollery/internal/endpoints/v1"
)

const (
	listAddr = ":8000"
	webDir   = "web"
)

var (
	//go:embed web
	webFS embed.FS
)

// handleStatic serves te embeded web directory
func handleStatic(w http.ResponseWriter, r *http.Request) {
	fsys, err := fs.Sub(webFS, webDir)
	if err != nil {
		panic(err)
	}
	if r.URL.Path != "/" {
		file, err := fsys.Open(r.URL.Path)
		if errors.Is(err, fs.ErrNotExist) {
			r.URL.Path = "/"
		}
		if file != nil {
			if err = file.Close(); err != nil {
				log.Println("failed to close file")
			}
		}
	}
	handlers.CompressHandler(http.FileServerFS(fsys)).ServeHTTP(w, r)
}

// Serve starts the web server
func Serve() error {
	fmt.Printf("Starting server on %s\n", listAddr)

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	v1.RegisterAPIRoutes(api)

	// If we don't have a route for it, probably frontend
	r.NotFoundHandler = http.HandlerFunc(handleStatic)

	return http.ListenAndServe(listAddr, r)
}

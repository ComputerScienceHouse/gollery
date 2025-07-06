package main

import (
	"log"

	db "github.com/ComputerScienceHouse/gollery/internal/services"
	"github.com/ComputerScienceHouse/gollery/server"
)

func main() {
	db_err := db.ConnectToDB()
	if db_err != nil {
		log.Printf("error connecting to postgresql db %v", db_err)
	}
	defer db.DisconnectFromDB()

	err := server.Serve()
	if err != nil {
		log.Fatalln(err)
	}
}

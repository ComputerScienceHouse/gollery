package main

import (
	"log"

	"github.com/ComputerScienceHouse/gollery/server"
)

func main() {
	err := server.Serve()
	if err != nil {
		log.Fatalln(err)
	}
}

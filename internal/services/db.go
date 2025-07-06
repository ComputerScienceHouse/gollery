package services

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectToDB() error {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", "postgresql", 5432, "postgres", "postgres", "postgres")
	DB, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Print("Cannot connect to postgres dumbass")
		return err
	}
	return nil
}

func DisconnectFromDB() error {
	return DB.Close()
}

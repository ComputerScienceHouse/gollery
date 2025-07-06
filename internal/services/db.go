package services

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectToDB() error {
	var err error
	DB, err = sqlx.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return err
	}
	return nil
}

func DisconnectFromDB() error {
	return DB.Close()
}

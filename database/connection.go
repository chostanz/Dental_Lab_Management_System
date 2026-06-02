package database

import (
	"log"

	sqlx "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB = Connection()

func Connection() *sqlx.DB {
	conn, err := sqlx.Connect("postgres", "user=postgres password=dewisuperkeren dbname=dental_lab sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

package database

import (
	"log"
	"os" // Wajib tambah ini untuk membaca variabel environment

	sqlx "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB = Connection()

func Connection() *sqlx.DB {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		dsn = "host=localhost user=postgres password=dewisuperkeren dbname=dental_lab sslmode=disable"
	}

	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

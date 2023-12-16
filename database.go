package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	Db *sql.DB
)

func ConnectDB() {
	hostname := "localhost"
	username := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	database := os.Getenv("DBNAME")
	sslMode := "disable"

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", username, password, hostname, database, sslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	Db = db

	/*
		err = Db.Ping()
		if err != nil {
			log.Fatal(err)
		}*/
}

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	Db *sql.DB
)

func ConnectDB() {
	hostname := "localhost"
	username := "dvolkov"
	password := "unity020kek"
	database := "course-tracker"
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

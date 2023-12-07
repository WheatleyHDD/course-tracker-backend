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

func connectDB() {
	hostname := "localhost"
	username := "dvolkov"
	password := ""
	database := "course-tracker"
	sslMode := "disable"

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", username, password, hostname, database, sslMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

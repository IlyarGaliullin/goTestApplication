package databaseDriver

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "admin"
	database = "test_db"
)

func NewConnection() *sql.DB {

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)

	var err error
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

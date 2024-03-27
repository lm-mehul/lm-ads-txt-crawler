package models

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

func SetupSQLConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "lemma:admin@tcp(127.0.0.1:3306)/lemma_crawler")
	if err != nil {
		log.Printf("Could not connect to database: %v", err)
	}
	return db, err
}

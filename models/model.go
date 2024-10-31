package models

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

func SetupSQLConn() (*sql.DB, error) {

	db, err := sql.Open("mysql", "lemma_rw:Lemm@r0cks!@tcp(23.108.100.104)/lm_teda_crawler_v2")
	if err != nil {
		log.Printf("Could not connect to database: %v", err)
	}

	// db, err := sql.Open("mysql", "lemma:admin@tcp(127.0.0.1:3306)/lemma_crawler")
	// if err != nil {
	// 	log.Printf("Could not connect to database: %v", err)
	// }
	return db, err
}

func SetupLMTedaSQLConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "lemma_rw:Lemm@r0cks!@tcp(23.108.100.104)/lm_teda")
	if err != nil {
		log.Printf("Could not connect to database: %v", err)
	}
	return db, err
}

func SetupLmAdsTxtSQLConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "lemma_rw:Lemm@r0cks!@tcp(23.108.100.104)/lm_ads_txt")
	if err != nil {
		log.Printf("Could not connect to database: %v", err)
	}
	return db, err
}

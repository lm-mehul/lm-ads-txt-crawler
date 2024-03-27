package main

import (
	"log"
	"os"

	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/service"
)

func main() {
	db, err := models.SetupSQLConn()
	if err != nil {
		log.Printf("Could not connect to database: %v", err)
		db.Close() // Example of closing a database connection
		os.Exit(1) // Exit after cleanup with a non-zero status code
	}
	defer db.Close()

	s := service.NewService(db)
	s.Start()
}

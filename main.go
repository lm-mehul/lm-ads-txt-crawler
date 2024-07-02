package main

import (
	"log"
	"os"

	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/server"
)

// Run script in the format: go run <filename>.go --script_type <1/2/3/4/5/6/7>
func main() {
	db, err := models.SetupSQLConn()
	if err != nil {
		log.Printf("Could not connect to database: %v", err)
		db.Close()
		os.Exit(1) // Exit after cleanup with a non-zero status code
	}
	defer db.Close()

	constant.InitConstants()

	s := server.NewService(db)
	s.Start()
}

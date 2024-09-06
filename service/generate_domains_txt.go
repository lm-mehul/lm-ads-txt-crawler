package service

import (
	"database/sql"
	"log"

	"github.com/lemmamedia/ads-txt-crawler/utils"
)

func FetchDomains(db *sql.DB, parserType string) ([]string, string, error) {
	var domainFileName, pageType string
	var domainsList []string
	if parserType == "app-ads" {
		domainFileName = "app-ads_txt_domains.txt"
		pageType = "app-ads.txt"
	} else {
		domainFileName = "ads_txt_domains.txt"
		pageType = "ads.txt"
	}

	// Run WriteStringArrayToFile in background as a goroutine

	err := utils.WriteStringArrayToFile(domainFileName, domainsList)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
	log.Printf("Data written to %s successfully.", domainFileName)

	return domainsList, pageType, nil
}

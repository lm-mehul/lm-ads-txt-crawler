package service

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lemmamedia/ads-txt-crawler/utils"
)

func GenerateDomainTxt(db *sql.DB) {
	fmt.Println("Fetching domains from Database...")
	_, _, err := FetchDomains(db, "ads")
	if nil != err {
		log.Printf("Error fetching domains from database with error : %v", err)
		return
	}
	_, _, err = FetchDomains(db, "app-ads")
	if nil != err {
		log.Printf("Error fetching domains from database with error : %v", err)
		return
	}
}

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

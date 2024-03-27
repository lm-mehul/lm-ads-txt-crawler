package service

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/models"
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

		androidBundles, err := models.GetDomainsFromDB(db, constant.BUNDLE_MOBILE_ANDROID)
		if err != nil {
			log.Printf("Error fetching : %v bundles from database with error : %v", constant.BUNDLE_MOBILE_ANDROID, err)
			return domainsList, pageType, err
		}
		domainsList = append(domainsList, androidBundles...)
		ctvBundles, err := models.GetDomainsFromDB(db, constant.BUNDLE_CTV)
		if err != nil {
			log.Printf("Error fetching : %v bundles from database with error : %v", constant.BUNDLE_CTV, err)
			return domainsList, pageType, err
		}
		domainsList = append(domainsList, ctvBundles...)
		iOSBundles, err := models.GetDomainsFromDB(db, constant.BUNDLE_MOBILE_IOS)
		if err != nil {
			log.Printf("Error fetching : %v bundles from database with error : %v", constant.BUNDLE_MOBILE_IOS, err)
			return domainsList, pageType, err
		}
		domainsList = append(domainsList, iOSBundles...)

	} else {
		domainFileName = "ads_txt_domains.txt"
		pageType = "ads.txt"

		webBundles, err := models.GetDomainsFromDB(db, constant.BUNDLE_WEB)
		if err != nil {
			log.Printf("Error fetching : %v bundles from database with error : %v", constant.BUNDLE_WEB, err)
			return domainsList, pageType, err
		}
		domainsList = append(domainsList, webBundles...)
	}

	// Run WriteStringArrayToFile in background as a goroutine
	go func() {
		err := utils.WriteStringArrayToFile(domainFileName, domainsList)
		if err != nil {
			log.Fatalf("Error writing to file: %v", err)
		}
		log.Printf("Data written to %s successfully.", domainFileName)
	}()

	return domainsList, pageType, nil
}

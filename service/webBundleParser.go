package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/models"
)

func webParser(db *sql.DB) {
	webBundles, err := models.GetBundlesFromDB(db, constant.BUNDLE_WEB)
	if err != nil {
		log.Printf("Error fetching : %v bundles from database with error : %v", constant.BUNDLE_WEB, err)
		return
	}

	fmt.Println("Executing Web bundle parser...")
	var bundles []models.BundleInfo
	var bundle models.BundleInfo
	batchCount := 0

	for _, webBundle := range webBundles {
		bundle.Bundle = webBundle
		bundle.Category = constant.BUNDLE_WEB
		bundle.Domain = extractDomainForWebParser(webBundle)

		bundles = append(bundles, bundle)
		batchCount++

		// If batch size is reached, insert the batch into the database
		if batchCount == constant.BATCH_SIZE {
			err := models.SaveCrawledBundlesInDB(db, bundles)
			if nil != err {
				log.Fatal("Failed to save bundles in DB")
				continue
			}

			// Reset batch count and values
			batchCount = 0
			bundles = []models.BundleInfo{}
		}
	}
	// Insert the remaining batch
	if batchCount > 0 {
		err = models.SaveCrawledBundlesInDB(db, bundles)
		if err != nil {
			log.Printf("Error inserting %v bundles into database with error : %v", constant.BUNDLE_MOBILE_ANDROID, err)
		}
	}
}

func extractDomainForWebParser(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("Error processing URL '%s': %s\n", rawURL, err)
		return ""
	}

	if strings.Contains(rawURL, "/") {
		fmt.Printf("Parsed URL: %+v\n", parsedURL)
		fmt.Printf("parsedURL.Host: %s\n", parsedURL.Host)
		return parsedURL.Host
	}

	return rawURL
}

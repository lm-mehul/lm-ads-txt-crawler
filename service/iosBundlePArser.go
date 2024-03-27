package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/utils"
)

func iosBundleParser(db *sql.DB) {
	iOSBundles, err := models.GetBundlesFromDB(db, constant.BUNDLE_MOBILE_IOS)
	if err != nil {
		log.Printf("Error fetching : %v bundles from database with error : %v", constant.BUNDLE_MOBILE_IOS, err)
		return
	}
	fmt.Println("Executing iOS bundle parser...")
	var bundles []models.BundleInfo
	var bundle models.BundleInfo
	batchCount := 0

	for _, iOSBundle := range iOSBundles {
		var appleStoreURL string
		url := fmt.Sprintf("https://apps.apple.com/us/app/%s/id%s", iOSBundle, iOSBundle)

		response, err := http.Head(url)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			appleStoreURL = response.Request.URL.String()
		} else {
			utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, fmt.Sprintf(" Error : %d\n", response.StatusCode))
			continue
		}

		if appleStoreURL != "" {
			response, err := http.Get(appleStoreURL)
			if err != nil {
				// fmt.Printf("Error: %s\n", err)
				utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
				continue
			}
			defer response.Body.Close()

			if response.StatusCode == 200 {
				doc, err := goquery.NewDocumentFromReader(response.Body)
				if err != nil {
					// fmt.Printf("Error: %s\n", err)
					utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
					continue
				}

				websiteElement := doc.Find("a.link.icon.icon-after.icon-external")

				if websiteElement.Length() > 0 {
					associatedWebsiteURL, _ := websiteElement.Attr("href")
					bundle.Website = associatedWebsiteURL
					bundle.Bundle = iOSBundle
					bundle.Category = constant.BUNDLE_MOBILE_IOS
					bundle.Domain = extractDomainFromBundleURL(bundle.Website)

					bundles = append(bundles, bundle)

				} else {
					utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "No associated website")
					continue
				}
			} else {
				// fmt.Printf("Error: %d\n", response.StatusCode)
				utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
				continue
			}
		}
		batchCount++

		// If batch size is reached, insert the batch into the database
		if batchCount == constant.BATCH_SIZE {
			err := models.SaveCrawledBundlesInDB(db, bundles)
			if nil != err {
				log.Fatal("Failed to save bundles in DB")
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

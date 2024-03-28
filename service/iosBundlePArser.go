package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

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

	var wg sync.WaitGroup
	batchSize := constant.BATCH_SIZE
	numBatches := (len(iOSBundles) + batchSize - 1) / batchSize // Calculate number of batches

	for i := 0; i < numBatches; i++ {
		startIndex := i * batchSize
		endIndex := (i + 1) * batchSize
		if endIndex > len(iOSBundles) {
			endIndex = len(iOSBundles)
		}
		batch := iOSBundles[startIndex:endIndex]

		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			processIOSBatch(db, batch)
		}(batch)
	}

	wg.Wait()

	// Handle remaining bundles
	if numBatches*batchSize < len(iOSBundles) {
		remaining := iOSBundles[numBatches*batchSize:]
		processIOSBatch(db, remaining)
	}
}

func processIOSBatch(db *sql.DB, batch []string) {
	var bundles []models.BundleInfo
	for _, iOSBundle := range batch {
		var bundle models.BundleInfo

		url := fmt.Sprintf("https://apps.apple.com/us/app/%s/id%s", iOSBundle, iOSBundle)
		response, err := http.Head(url)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, fmt.Sprintf("Error: %d", response.StatusCode))
			continue
		}

		appleStoreURL := response.Request.URL.String()

		response, err = http.Get(appleStoreURL)
		if err != nil {
			utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
			continue
		}

		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
			continue
		}

		websiteElement := doc.Find("a.link.icon.icon-after.icon-external")
		if websiteElement.Length() == 0 {
			utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "No associated website")
			continue
		}

		associatedWebsiteURL, _ := websiteElement.Attr("href")
		bundle.Website = strings.TrimSpace(associatedWebsiteURL)
		bundle.Bundle = iOSBundle
		bundle.Category = constant.BUNDLE_MOBILE_IOS
		bundle.Domain = extractDomainFromBundleURL(strings.TrimSpace(bundle.Website))

		bundles = append(bundles, bundle)
	}

	// Save bundles and uncrawled domains in the database
	err := models.SaveCrawledBundlesInDB(db, bundles)
	if err != nil {
		log.Printf("Error inserting bundles into database: %v", err)
	}
	err = models.SaveUnCrawledDomainsInDB(db, bundles)
	if err != nil {
		log.Printf("Error saving uncrawled domains into database: %v", err)
	}
}

// func iosBundleParser(db *sql.DB) {
// 	iOSBundles, err := models.GetBundlesFromDB(db, constant.BUNDLE_MOBILE_IOS)
// 	if err != nil {
// 		log.Printf("Error fetching : %v bundles from database with error : %v", constant.BUNDLE_MOBILE_IOS, err)
// 		return
// 	}
// 	fmt.Println("Executing iOS bundle parser...")
// 	var bundles []models.BundleInfo
// 	var bundle models.BundleInfo
// 	batchCount := 0

// 	for _, iOSBundle := range iOSBundles {
// 		var appleStoreURL string
// 		url := fmt.Sprintf("https://apps.apple.com/us/app/%s/id%s", iOSBundle, iOSBundle)

// 		response, err := http.Head(url)
// 		if err != nil {
// 			fmt.Printf("Error: %s\n", err)
// 			continue
// 		}
// 		defer response.Body.Close()

// 		if response.StatusCode == 200 {
// 			appleStoreURL = response.Request.URL.String()
// 		} else {
// 			utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, fmt.Sprintf(" Error : %d\n", response.StatusCode))
// 			continue
// 		}

// 		if appleStoreURL != "" {
// 			response, err := http.Get(appleStoreURL)
// 			if err != nil {
// 				// fmt.Printf("Error: %s\n", err)
// 				utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
// 				continue
// 			}
// 			defer response.Body.Close()

// 			if response.StatusCode == 200 {
// 				doc, err := goquery.NewDocumentFromReader(response.Body)
// 				if err != nil {
// 					// fmt.Printf("Error: %s\n", err)
// 					utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
// 					continue
// 				}

// 				websiteElement := doc.Find("a.link.icon.icon-after.icon-external")

// 				if websiteElement.Length() > 0 {
// 					associatedWebsiteURL, _ := websiteElement.Attr("href")
// 					bundle.Website = associatedWebsiteURL
// 					bundle.Bundle = iOSBundle
// 					bundle.Category = constant.BUNDLE_MOBILE_IOS
// 					bundle.Domain = extractDomainFromBundleURL(bundle.Website)

// 					bundles = append(bundles, bundle)

// 				} else {
// 					utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "No associated website")
// 					continue
// 				}
// 			} else {
// 				// fmt.Printf("Error: %d\n", response.StatusCode)
// 				utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
// 				continue
// 			}
// 		}
// 		batchCount++

// 		// If batch size is reached, insert the batch into the database
// 		if batchCount == constant.BATCH_SIZE {
// 			err := models.SaveCrawledBundlesInDB(db, bundles)
// 			if nil != err {
// 				log.Fatal("Failed to save bundles in DB")
// 				continue
// 			}
// 			err = models.SaveUnCrawledDomainsInDB(db, bundles)
// 			if nil != err {
// 				log.Fatal("Failed to save bundles in DB")
// 				continue
// 			}

// 			// Reset batch count and values
// 			batchCount = 0
// 			bundles = []models.BundleInfo{}
// 		}

// 	}
// 	// Insert the remaining batch
// 	if batchCount > 0 {
// 		err = models.SaveCrawledBundlesInDB(db, bundles)
// 		if err != nil {
// 			log.Printf("Error inserting %v bundles into database with error : %v", constant.BUNDLE_MOBILE_ANDROID, err)
// 		}
// 		err = models.SaveUnCrawledDomainsInDB(db, bundles)
// 		if nil != err {
// 			log.Fatal("Failed to save bundles in DB")
// 		}
// 	}
// }

package parsers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/utils"
)

func IosBundleParser(db *sql.DB) {
	// iOSBundles, err := repository.GetBundlesFromDB(db, constant.BUNDLE_MOBILE_IOS)
	// if err != nil {
	// 	log.Printf("Error fetching : %v bundles from database with error : %v", constant.BUNDLE_MOBILE_IOS, err)
	// 	return
	// }
	fmt.Println("Executing iOS bundle parser...")

	processIOSBatch(db, models.IOSBundles)
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

		fmt.Printf("iOS - Bundle: %s, Website: %s, Domain: %s\n", bundle.Bundle, bundle.Website, bundle.Domain)

		bundles = append(bundles, bundle)
	}

	// Save bundles and uncrawled domains in the database
	// err := repository.SaveCrawledBundlesInDB(db, bundles)
	// if err != nil {
	// 	log.Printf("Error inserting bundles into database: %v", err)
	// }
	// err = repository.SaveUnCrawledDomainsInDB(db, bundles)
	// if err != nil {
	// 	log.Printf("Error saving uncrawled domains into database: %v", err)
	// }
}

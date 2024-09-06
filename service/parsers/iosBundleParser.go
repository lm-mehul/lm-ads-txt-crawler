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
)

func IosBundleParser(db *sql.DB) {
	fmt.Println("Executing iOS bundle parser...")
	ProcessIOSBundle(db, models.IOSBundles[0])
}

func ProcessIOSBundle(db *sql.DB, iOSBundle string) models.BundleInfo {
	var bundle models.BundleInfo

	url := fmt.Sprintf("https://apps.apple.com/us/app/%s/id%s", iOSBundle, iOSBundle)
	response, err := http.Head(url)
	if err != nil {
		// fmt.Printf("Error: %s\n", err)
		return bundle
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		// utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, fmt.Sprintf("Error: %d", response.StatusCode))
		return bundle
	}

	appleStoreURL := response.Request.URL.String()

	response, err = http.Get(appleStoreURL)
	if err != nil {
		// utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
		return bundle
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		// utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
		return bundle
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		// utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "Invalid iOS Bundle")
		return bundle
	}

	websiteElement := doc.Find("a.link.icon.icon-after.icon-external")
	if websiteElement.Length() == 0 {
		// utils.LogBundleError(iOSBundle, constant.BUNDLE_MOBILE_IOS, "No associated website")
		return bundle
	}

	associatedWebsiteURL, _ := websiteElement.Attr("href")
	bundle.Website = strings.TrimSpace(associatedWebsiteURL)
	bundle.Bundle = iOSBundle
	bundle.Category = constant.BUNDLE_MOBILE_IOS
	bundle.Domain = extractDomainFromBundleURL(strings.TrimSpace(bundle.Website))

	fmt.Printf("iOS - Bundle: %s, Website: %s, Domain: %s\n", bundle.Bundle, bundle.Website, bundle.Domain)
	return bundle
}

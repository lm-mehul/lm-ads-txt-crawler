package parsers

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

func WebParser(db *sql.DB) {

	fmt.Println("Executing Web bundle parser...")
	var bundles []models.BundleInfo
	var bundle models.BundleInfo
	webBundles := []string{"com.google.android.apps.maps", "com.google.android.apps.docs", "com.google.android.apps.photos"}

	for _, webBundle := range webBundles {
		bundle.Bundle = webBundle
		bundle.Category = constant.BUNDLE_WEB
		bundle.Domain = extractDomainForWebParser(webBundle)

		bundles = append(bundles, bundle)

	}
}

func extractDomainForWebParser(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("Error processing URL '%s': %s\n", rawURL, err)
		return ""
	}

	if strings.Contains(rawURL, "/") {
		fmt.Printf("Parsed URL: %+v\n", parsedURL)
		fmt.Printf("parsedURL.Host: %s\n", parsedURL.Host)
		return strings.TrimSpace(parsedURL.Host)
	}

	return strings.TrimSpace(rawURL)
}

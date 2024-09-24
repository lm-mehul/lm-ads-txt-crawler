package parsers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/models"
)

// webBundles := []string{"com.google.android.apps.maps", "com.google.android.apps.docs", "com.google.android.apps.photos"}

func ProcessWebBundle(db *sql.DB, webBundle string) (models.BundleInfo, error) {

	var bundle models.BundleInfo
	var err error

	bundle.Bundle = webBundle
	bundle.Category = constant.BUNDLE_WEB
	bundle.Domain, err = extractDomainForWebParser(webBundle)
	if err != nil {
		return bundle, errors.New("Error extracting domain for web parser")
	}
	return bundle, nil
}

func extractDomainForWebParser(rawURL string) (string, error) {

	// Remove unwanted escape sequences like %20 (space) and %09 (tab)
	rawURL = strings.ReplaceAll(rawURL, "%20", "")
	rawURL = strings.ReplaceAll(rawURL, "%09", "")

	if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
		parsedURL, err := url.ParseRequestURI(rawURL)
		if err != nil {
			return "", errors.New("Error processing URL")
		}

		if strings.Contains(rawURL, "/") {
			fmt.Printf("Host: %v\n", parsedURL.Host)
			return strings.TrimSpace(parsedURL.Host), nil
		}
	}

	return strings.TrimSpace(rawURL), nil
}

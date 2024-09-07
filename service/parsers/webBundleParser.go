package parsers

import (
	"database/sql"
	"errors"
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
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("Error processing URL")
	}

	if strings.Contains(rawURL, "/") {
		return strings.TrimSpace(parsedURL.Host), nil
	}

	return strings.TrimSpace(rawURL), nil
}

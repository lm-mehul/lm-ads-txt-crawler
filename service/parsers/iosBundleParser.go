package parsers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/models"
)

func ProcessIOSBundle(db *sql.DB, iOSBundle string) (models.BundleInfo, error) {
	var bundle models.BundleInfo

	url := fmt.Sprintf("https://apps.apple.com/us/app/%s/id%s", iOSBundle, iOSBundle)
	response, err := http.Head(url)
	if err != nil {
		return bundle, errors.New("Invalid iOS Bundle")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return bundle, errors.New("Invalid iOS Bundle")
	}

	appleStoreURL := response.Request.URL.String()

	response, err = http.Get(appleStoreURL)
	if err != nil {
		return bundle, errors.New("Invalid iOS Bundle")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return bundle, errors.New("Invalid iOS Bundle")
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return bundle, errors.New("Invalid iOS Bundle")
	}

	websiteElement := doc.Find("a.link.icon.icon-after.icon-external")
	if websiteElement.Length() == 0 {
		return bundle, errors.New("No associated website")
	}

	associatedWebsiteURL, _ := websiteElement.Attr("href")
	bundle.Website = strings.TrimSpace(associatedWebsiteURL)
	bundle.Bundle = iOSBundle
	bundle.Category = constant.BUNDLE_MOBILE_IOS
	bundle.Domain = extractDomainFromBundleURL(strings.TrimSpace(bundle.Website))

	return bundle, nil
}

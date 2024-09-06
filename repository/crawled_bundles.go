package repository

import (
	"bytes"
	"database/sql"
	"log"
	"strings"

	"github.com/lemmamedia/ads-txt-crawler/models"
)

// SaveCrawledBundlesInDB inserts a bundle domain and category into the crawled_bundles table.
func SaveCrawledBundlesInDB(db *sql.DB, bundles []models.BundleInfo) error {

	var buff bytes.Buffer
	buff.WriteString("INSERT IGNORE INTO crawled_bundles(bundle, category, website, domain, ads_txt_URL, app_ads_txt_URL, ads_txt_Hash, app_ads_txt_Hash) VALUES ")
	values := make([]interface{}, 0)
	validBundleCount := int64(0)
	for index := range bundles {
		values = append(values,
			strings.TrimSpace(bundles[index].Bundle),
			bundles[index].Category,
			strings.TrimSpace(bundles[index].Website),
			strings.TrimSpace(bundles[index].Domain),
			strings.TrimSpace(bundles[index].AdsTxtURL),
			strings.TrimSpace(bundles[index].AppAdsTxtURL),
			strings.TrimSpace(bundles[index].AdsTxtHash),
			strings.TrimSpace(bundles[index].AppAdsTxtHash),
		)
		validBundleCount++
	}
	placeholder := strings.Repeat("(?,?,?,?,?,?,?,?), ", int(validBundleCount))
	if validBundleCount > 0 {
		placeholder = placeholder[:len(placeholder)-2] // Remove the trailing comma and space
	}
	buff.WriteString(placeholder)

	// Execute the query with the collected values
	_, err := db.Exec(buff.String(), values...)
	if err != nil {
		log.Printf("Error executing crawled_bundles data insert: %v", err)
		return err
	}
	return nil
}

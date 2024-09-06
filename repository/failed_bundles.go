package repository

import (
	"bytes"
	"database/sql"
	"log"
	"strings"

	"github.com/lemmamedia/ads-txt-crawler/models"
)

// SaveFailedBundlesInDB inserts bundles and categories into the failed_bundles table.
func SaveFailedBundlesInDB(db *sql.DB, bundles []models.BundleInfo) error {
	var buff bytes.Buffer
	buff.WriteString("INSERT IGNORE INTO failed_bundles(bundle, category) VALUES ")
	values := make([]interface{}, 0)
	validBundleCount := int64(0)
	for index := range bundles {
		values = append(values,
			strings.TrimSpace(bundles[index].Bundle),
			bundles[index].Category,
		)
		validBundleCount++
	}
	placeholder := strings.Repeat("(?,?), ", int(validBundleCount))
	if validBundleCount > 0 {
		placeholder = placeholder[:len(placeholder)-2] // Remove the trailing comma and space
	}
	buff.WriteString(placeholder)

	// Execute the query with the collected values
	_, err := db.Exec(buff.String(), values...)
	if err != nil {
		log.Printf("Error executing failed_bundles data insert: %v", err)
		return err
	}
	return nil
}

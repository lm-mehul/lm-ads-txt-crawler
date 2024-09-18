package repository

import (
	"bytes"
	"database/sql"
	"log"
	"strings"

	"github.com/lemmamedia/ads-txt-crawler/models"
)

// SaveBundlesInDB inserts a bundle domain and category
func SaveBundlesInDB(db *sql.DB, bundles []models.BundleInfo) error {
	var buff bytes.Buffer
	buff.WriteString("INSERT IGNORE INTO bundles (bundle, category) VALUES ")
	values := make([]any, 0)
	validBundleCount := int64(0)
	for index, _ := range bundles {
		values = append(values, strings.TrimSpace(bundles[index].Bundle), bundles[index].Category)
		validBundleCount++
	}
	placeHolder := strings.Repeat("(?,?), ", int(validBundleCount))
	if validBundleCount > 0 {
		placeHolder = placeHolder[:len(placeHolder)-2]
	}
	buff.WriteString(placeHolder)
	_, err := db.Exec(buff.String(), values...)
	if err != nil {
		log.Printf("Error executing bundles data insert: %v", err)
		return err
	}
	return nil
}

func GetBundlesFromDB(db *sql.DB, limit, offset int) ([]models.BundleInfo, error) {
	var bundles []models.BundleInfo
	query := "SELECT bundle,category FROM bundles"
	params := make([]any, 0)

	if limit > 0 {
		query += "LIMIT ? OFFSET ?"
		params = append(params, limit, offset)
	}

	rows, err := db.Query(query, params...)
	if err != nil {
		return bundles, err
	}
	defer rows.Close()

	for rows.Next() {
		var bundle models.BundleInfo
		if err := rows.Scan(&bundle.Bundle, &bundle.Category); err != nil {
			return bundles, err
		}
		bundles = append(bundles, bundle)
	}

	return bundles, nil
}

func GetBundlesCount(db *sql.DB) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM bundles"
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

// SaveBundlesFromMasterSheet inserts distinct bundles and categories from the bundle_mastersheet table into the bundles table.
func SaveBundlesFromMasterSheet(db *sql.DB) error {
	query := `
		INSERT INTO bundles (bundle, category)
		SELECT DISTINCT
			bundle,
			CASE
				WHEN device = 'ios' THEN 'Mobile App IOS'
				WHEN device = 'android' THEN 'Mobile App Android'
				ELSE 'CTV'
			END AS category
		FROM lm_teda.bundle_mastersheet
		WHERE bundle IS NOT NULL;
	`

	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error executing bundles data insert: %v", err)
		return err
	}
	return nil
}

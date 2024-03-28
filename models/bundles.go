package models

import (
	"bytes"
	"database/sql"
	"log"
	"strings"
	"sync"
)

type BundleInfo struct {
	Bundle   string
	Category string
	Website  string
	Domain   string
}

// SaveBundlesInDB inserts a bundle domain and category
func SaveBundlesInDB(db *sql.DB, bundles []BundleInfo) error {
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

var dbMutex sync.Mutex // Declare a mutex variable

// SaveCrawledBundlesInDB inserts a bundle domain and category
func SaveCrawledBundlesInDB(db *sql.DB, bundles []BundleInfo) error {
	dbMutex.Lock()         // Lock the mutex before accessing the shared resource
	defer dbMutex.Unlock() // Unlock the mutex when done, even if an error occurs

	var buff bytes.Buffer
	buff.WriteString("INSERT IGNORE INTO crawled_bundles(bundle, category, website, domain) VALUES ")
	values := make([]interface{}, 0)
	validBundleCount := int64(0)
	for index := range bundles {
		values = append(values, strings.TrimSpace(bundles[index].Bundle), bundles[index].Category, strings.TrimSpace(bundles[index].Website), strings.TrimSpace(bundles[index].Domain))
		validBundleCount++
	}
	placeholder := strings.Repeat("(?,?,?,?), ", int(validBundleCount))
	if validBundleCount > 0 {
		placeholder = placeholder[:len(placeholder)-2]
	}
	buff.WriteString(placeholder)
	_, err := db.Exec(buff.String(), values...)
	if err != nil {
		log.Printf("Error executing crawled_bundles data insert: %v", err)
		return err
	}
	return nil
}

func GetBundlesFromDB(db *sql.DB, category string) ([]string, error) {
	var bundles []string
	query := "SELECT bundle FROM bundles WHERE category = ? AND is_deleted = 0"
	rows, err := db.Query(query, category)
	if err != nil {
		return bundles, err
	}
	defer rows.Close()

	for rows.Next() {
		var bundle string
		if err := rows.Scan(&bundle); err != nil {
			return bundles, err
		}
		bundles = append(bundles, bundle)
	}
	if err := rows.Err(); err != nil {
		return bundles, err
	}
	return bundles, nil
}

package models

import (
	"bytes"
	"database/sql"
	"log"
	"strings"
)

func GetDomainsFromDB(db *sql.DB, category string) ([]string, error) {
	var bundles []string
	query := "SELECT domain FROM un_crawled_domains WHERE category = ? AND is_deleted = 0"
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

// SaveUnCrawledDomainsInDB inserts a bundle domain and category
func SaveUnCrawledDomainsInDB(db *sql.DB, bundles []BundleInfo) error {
	dbMutex.Lock()         // Lock the mutex before accessing the shared resource
	defer dbMutex.Unlock() // Unlock the mutex when done, even if an error occurs

	var buff bytes.Buffer
	buff.WriteString("INSERT IGNORE INTO un_crawled_domains(domain, category) VALUES ")
	values := make([]interface{}, 0)
	validBundleCount := int64(0)
	for index := range bundles {
		values = append(values, strings.TrimSpace(bundles[index].Domain), bundles[index].Category)
		validBundleCount++
	}
	placeholder := strings.Repeat("(?,?), ", int(validBundleCount))
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

func IsDomainCrawled(domainName string, hash string, db *sql.DB) (bool, error) {
	query := `
        UPDATE crawled_bundles 
        SET ads_txt_page_hash = ?
        WHERE domain = ?
    `

	result, err := db.Exec(query, hash, domainName)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	// If no rows were affected, it means no update was performed, hence return true
	if rowsAffected == 0 {
		return true, nil
	}

	// If rows were affected, it means an update was performed, hence return false
	return false, nil
}

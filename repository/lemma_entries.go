package repository

import (
	"bytes"
	"database/sql"
	"log"
	"strings"

	"github.com/lemmamedia/ads-txt-crawler/models"
)

// SaveLemmaEntriesInDB inserts bundles, categories, and lemma information into the lemma_entries table.
func SaveLemmaEntriesInDB(db *sql.DB, entries []models.LemmaEntry) error {
	var buff bytes.Buffer
	buff.WriteString("INSERT IGNORE INTO lemma_entries(bundle, category, Lemma_Direct, Lemma_Reseller, ads_page_url, page_type) VALUES ")
	values := make([]interface{}, 0)
	validEntryCount := int64(0)
	for index := range entries {
		values = append(values,
			strings.TrimSpace(entries[index].Bundle),
			entries[index].Category,
			strings.TrimSpace(entries[index].LemmaDirect),
			strings.TrimSpace(entries[index].LemmaReseller),
			entries[index].AdsPageURL,
			entries[index].PageType,
		)
		validEntryCount++
	}
	placeholder := strings.Repeat("(?,?,?,?,?,?), ", int(validEntryCount))
	if validEntryCount > 0 {
		placeholder = placeholder[:len(placeholder)-2] // Remove the trailing comma and space
	}
	buff.WriteString(placeholder)

	// Execute the query with the collected values
	_, err := db.Exec(buff.String(), values...)
	if err != nil {
		log.Printf("Error executing lemma_entries data insert: %v", err)
		return err
	}
	return nil
}

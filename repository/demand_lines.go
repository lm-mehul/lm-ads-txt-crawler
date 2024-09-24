package repository

import (
	"bytes"
	"database/sql"
	"log"
	"strings"

	"github.com/lemmamedia/ads-txt-crawler/models"
)

func SaveDemandLinesResultInDB(db *sql.DB, bundles []models.DemandLinesEntry) error {
	var buff bytes.Buffer
	buff.WriteString("INSERT INTO bundle_demand_lines(bundle_id, category, domain, demand_line, ads_page_url, page_type) VALUES ")
	values := make([]interface{}, 0)
	validBundleCount := int64(0)

	for index := range bundles {
		values = append(values,
			bundles[index].Bundle, // Assuming bundle_id comes from another function
			bundles[index].Category,
			bundles[index].Domain,
			bundles[index].DemandLine,
			bundles[index].AdsPageURL,
			bundles[index].PageType,
		)
		validBundleCount++
	}

	placeholder := strings.Repeat("(?,?,?,?,?,?), ", int(validBundleCount))
	if validBundleCount > 0 {
		placeholder = placeholder[:len(placeholder)-2] // Remove the trailing comma and space
	}
	buff.WriteString(placeholder)

	// Execute the query with the collected values
	_, err := db.Exec(buff.String(), values...)
	if err != nil {
		log.Printf("Error executing bundle_demand_lines data insert: %v", err)
		return err
	}
	return nil
}

func GetDemandLinesFromDB(db *sql.DB) ([]string, error) {
	var demandLines []string
	query := "SELECT demand_line FROM ads_txt_demand_lines"

	rows, err := db.Query(query)
	if err != nil {
		return demandLines, err
	}
	defer rows.Close()

	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return demandLines, err
		}
		demandLines = append(demandLines, line)
	}

	return demandLines, nil
}

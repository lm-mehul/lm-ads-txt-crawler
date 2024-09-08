package handler

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/repository"
	"github.com/xuri/excelize/v2"
)

func PopulateBundles(db *sql.DB) {

	const batchSize = 1000 // Adjust this batch size as needed

	// Path to the Excel file
	filePath := "/home/lemma/Desktop/CMS-2/crawler/lm-ads-txt-crawler/resources/domains/bundles.xlsx"

	// Read data from the Excel file
	bundles, err := readBundlesExcel(filePath)
	if err != nil {
		log.Fatalf("Error reading bundles from Excel: %v", err)
	}

	// Process bundles in batches
	totalBundles := len(bundles)
	for i := 0; i < totalBundles; i += batchSize {
		end := i + batchSize
		if end > totalBundles {
			end = totalBundles
		}

		batch := bundles[i:end]

		// Save the current batch to the database
		if err := repository.SaveBundlesInDB(db, batch); err != nil {
			log.Fatalf("Error saving bundles to database: %v", err)
		}
	}

	fmt.Printf("Total %v bundles successfully saved to the database.\n", totalBundles)
}

// readBundlesExcel reads bundles from an Excel file and returns a slice of Bundle structs
func readBundlesExcel(filePath string) ([]models.BundleInfo, error) {
	// Open the Excel file
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer file.Close()

	// Get the first sheet name
	sheetName := file.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	// Read all rows from the first sheet
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows from sheet: %w", err)
	}

	// Slice to hold Bundle structs
	var bundles []models.BundleInfo

	// Read data rows, skipping the header
	for i, row := range rows {
		if i == 0 {
			continue // Skip the header row
		}
		if len(row) < 2 {
			continue // Skip rows with insufficient columns
		}
		bundle := models.BundleInfo{
			Bundle:   strings.TrimSpace(row[0]),
			Category: strings.TrimSpace(row[1]),
		}
		bundles = append(bundles, bundle)
	}

	return bundles, nil
}

package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/tealeg/xlsx"
)

func BundleParser(db *sql.DB) {
	fmt.Println("Executing Bundle parser...")

	readCSVFile(db)
	fmt.Print("Successfully conpleted parsing bundles CSV file...\n")

	// WaitGroup to wait for all parsers to finish
	var wg sync.WaitGroup

	// Add the number of parsers to WaitGroup
	wg.Add(4)

	// Run each parser in its own goroutine
	go func() {
		defer wg.Done()
		androidBundleParser(db)
	}()

	go func() {
		defer wg.Done()
		iosBundleParser(db)
	}()

	go func() {
		defer wg.Done()
		ctvBundleParser(db)
	}()

	go func() {
		defer wg.Done()
		webParser(db)
	}()

	// Wait for all parsers to finish
	wg.Wait()

	fmt.Println("All parsers have finished.")
	fmt.Println("Bundle Adstxt Parser completed.")

}

func readCSVFile(db *sql.DB) error {

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return err
	}
	// Open the XLSX file
	xlFile, err := xlsx.OpenFile(dir + "/resources/domains/bundles.xlsx")
	if err != nil {
		log.Fatalf("Error opening XLSX file: %v", err)
	}

	// Assume there's only one sheet
	sheet := xlFile.Sheets[0]

	batchCount := 0
	batchSize := 5000
	var bundles []models.BundleInfo

	// Iterate through each row in the sheet
	for _, row := range sheet.Rows {
		// Read data from each column in the row
		var bundleInfo models.BundleInfo
		for idx, cell := range row.Cells {
			if idx == 0 {
				bundleInfo.Bundle = cell.String()
			} else if idx == 1 {
				bundleInfo.Category = cell.String()
			}
		}
		bundles = append(bundles, bundleInfo)

		batchCount++

		// If batch size is reached, insert the batch into the database
		if batchCount == batchSize {
			err := models.SaveBundlesInDB(db, bundles)
			if nil != err {
				log.Fatal("Failed to save bundles in DB")
				return err
			}

			// Reset batch count and values
			batchCount = 0
			bundles = []models.BundleInfo{}
		}
	}

	// Insert the remaining batch
	if batchCount > 0 {
		err := models.SaveBundlesInDB(db, bundles)
		if err != nil {
			log.Printf("Error inserting remaining batch into database: %v", err)
		}
	}
	fmt.Println("Data insertion complete.")
	return nil
}

func extractDomainFromBundleURL(urlStr string) string {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error processing URL '%s': %v\n", urlStr, r)
		}
	}()

	if strings.Contains(urlStr, "/") {
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			panic(err)
		}
		return parsedURL.Hostname()
	} else {
		return strings.TrimSpace(urlStr)
	}
}

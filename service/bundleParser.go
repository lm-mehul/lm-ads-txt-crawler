package service

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func BundleParser(db *sql.DB) {
	fmt.Println("Executing Bundle parser...")

	// Run each parser in its own goroutine
	parser.AndroidBundleParser(db)
	parser.IosBundleParser(db)
	parser.CTVBundleParser(db)
	parser.WebParser(db)

	fmt.Println("All parsers have finished.")

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

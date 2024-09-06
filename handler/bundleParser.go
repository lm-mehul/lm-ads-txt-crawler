package handler

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func BundleParser(db *sql.DB) {
	fmt.Println("Executing Bundle parser...")

	// Run each parser in its own goroutine
	// parsers.AndroidBundleParser(db)
	// parsers.IosBundleParser(db)
	// parsers.CTVBundleParser(db)
	// parsers.WebParser(db)

	// domains := []string{"americasvoicenews.com", "americasvoice.news", "www.paltalk.com"}

	// for i, domain := range domains {
	// 	adsTxtPage, err := CrawlDomain(domain, "ads.txt")
	// 	if err != nil {
	// 		log.Printf("Error crawling domain %s: %v", domain, err)
	// 	}
	// 	fmt.Printf("hello %v\n%v\n\n", i, string(adsTxtPage))
	// }

	fmt.Println("All parsers have finished.")

}

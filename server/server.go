package server

import (
	"database/sql"
	"flag"
	"fmt"
	"time"

	"github.com/lemmamedia/ads-txt-crawler/parsers"
)

type Service struct {
	TotalErrors int // A global counter for errors
	db          *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Start() {
	start := time.Now()
	fmt.Println("------------Parser Started------------\n")
	scriptType := flag.Int("script_type", 0, "Specify script type (1 to 7)")

	flag.Parse()

	segregateBundle()

	switch *scriptType {
	case 1:
		// Run all Bundle+adtxt+category
		bundleParser(s.db)
		categoryParser(s.db)
		masterSheetCreation(s.db, "A")
	case 2:
		// Run Bundle+adtxt
		parsers.BundleParser(s.db)
		bundleParser(s.db)
	case 3:
		// Run only adtxt parser(ads)
		adsTxtParser(s.db, "ads")
	case 4:
		// Run only adtxt parser(app-ads)
		adsTxtParser(s.db, "app-ads")
	case 5:
		// Run adtxt parser(app-ads) + category
		adsTxtParser(s.db, "app-ads")
		categoryParser(s.db)
		masterSheetCreation(s.db, "B")
	case 6:
		// Run adtxt parser(ads) + category
		adsTxtParser(s.db, "ads")
		categoryParser(s.db)
		masterSheetCreation(s.db, "B")
	case 7:
		// Run only category
		fmt.Println("script 7 is running category parser started.........................")
		categoryParser(s.db)
		masterSheetCreation(s.db, "B")
	default:
		fmt.Println("Invalid script_type. Please provide 1 to 7.")
	}

	elapsed := time.Since(start)
	fmt.Printf("Total execution time: %s\n", elapsed)
}

// func (s *Service) Start() {

// 	switch *scriptType {
// 	case 1:
// 		BundleParser(s.db)
// 	case 2:
// 		AdsTxtLineCheck(s.db, "ads")
// 	case 3:
// 		AdsTxtLineCheck(s.db, "app-ads")
// 	case 4:
// 		FetchLemmaDirectsAndResellerInventory(s.db, "ads")
// 	case 5:
// 		FetchLemmaDirectsAndResellerInventory(s.db, "app-ads")
// 	default:
// 		fmt.Println("Invalid script_type. Please provide 1 or 2.")

// 	}

// 	elapsed := time.Since(start)
// 	fmt.Printf("Total execution time: %s\n", elapsed)
// }

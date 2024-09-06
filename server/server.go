package server

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lemmamedia/ads-txt-crawler/handler"
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

	// switch *scriptType {
	// case 1:
	// 	AdsTxtLineCheck(s.db, "ads")
	// case 2:
	// 	AdsTxtLineCheck(s.db, "app-ads")
	// case 3:
	handler.FetchLemmaDirectsAndResellerInventory(s.db)
	// case 4:
	//    BundlePArser(s.db)
	// case 5:
	//  generate Report
	// case 6:
	// populate bundle data
	// handler.PopulateBundles(s.db)
	// default:
	// 	fmt.Println("Invalid script_type. Please provide 1 or 2.")
	// 	fmt.Println("Run script in the format: go run <filename>.go --script_type <1/2>")
	// }

	elapsed := time.Since(start)
	fmt.Printf("Total execution time: %s\n", elapsed)
}

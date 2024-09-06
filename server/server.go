package server

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lemmamedia/ads-txt-crawler/service"
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

	service.BundleParser(s.db)
	// switch *scriptType {
	// case 1:
	// 	BundleParser(s.db)
	// case 2:
	// 	AdsTxtLineCheck(s.db, "ads")
	// case 3:
	// 	AdsTxtLineCheck(s.db, "app-ads")
	// case 4:
	// 	FetchLemmaDirectsAndResellerInventory(s.db, "ads")
	// case 5:
	// 	FetchLemmaDirectsAndResellerInventory(s.db, "app-ads")
	// default:
	// 	fmt.Println("Invalid script_type. Please provide 1 or 2.")
	// 	fmt.Println("Run script in the format: go run <filename>.go --script_type <1/2>")
	// }

	elapsed := time.Since(start)
	fmt.Printf("Total execution time: %s\n", elapsed)
}

package server

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	reader := bufio.NewReader(os.Stdin)

	for i := 1; ; i++ {
		// Display the menu
		fmt.Println("---------------------------------------------------------------------------------")
		fmt.Println("Select an option to run the corresponding script:")
		fmt.Println("1: AdsTxtLineCheck with 'ads'")
		fmt.Println("2: AdsTxtLineCheck with 'app-ads'")
		fmt.Println("3: Fetch Lemma Directs and Reseller Inventory")
		fmt.Println("4: Bundle Parser")
		fmt.Println("5: Generate Report")
		fmt.Println("6: Populate Bundle Data from CSV file")
		fmt.Println("0: Exit")
		fmt.Println("---------------------------------------------------------------------------------")

		// Read user input
		fmt.Print("Enter the script number to run: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Trim and parse input
		input = strings.TrimSpace(input)
		scriptType, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number between 0 and 6.")
			continue
		}

		// Exit the loop if user selects 0
		if scriptType == 0 {
			fmt.Println("Exiting the script.")
			break
		}

		start := time.Now()

		// Execute the selected operation based on the script type
		switch scriptType {
		case 1:
			// s.AdsTxtLineCheck(s.db, "ads")
		case 2:
			// s.AdsTxtLineCheck(s.db, "app-ads")
		case 3:
			handler.FetchLemmaDirectsAndResellerInventory(s.db)
		case 4:
			handler.BundleParser(s.db)
		case 5:
			// s.GenerateReport()
		case 6:
			handler.PopulateBundlesFromExcel(s.db)
		default:
			fmt.Println("Invalid option. Please select a number between 0 and 6.")
		}

		// Print the total execution time
		elapsed := time.Since(start)
		fmt.Printf("Total execution time: %s\n", elapsed)
		fmt.Printf("---------------------------------------------------------------------------------\n")
	}

}

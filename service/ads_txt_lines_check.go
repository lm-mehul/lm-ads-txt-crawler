package service

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/lemmamedia/ads-txt-crawler/utils"
)

var presentDomains []string

func AdsTxtLineCheck(db *sql.DB, parserType string) {
	fmt.Println("Executing adstxt parser...")
	fmt.Println("Fetching domains from Database...")

	// domainsList, pageType, err := FetchDomains(db, parserType)
	// if err != nil {
	// 	log.Printf("Error fetching domains from database: %v", err)
	// 	return
	// }

	// presentDomains = make([]string, 0)
	// fmt.Println("Page type:", pageType)

	// // Define batch size and concurrency limit
	// batchSize := 1000      // Adjust as needed
	// concurrencyLimit := 10 // Adjust as needed

	// var wg sync.WaitGroup
	// domainCh := make(chan string, batchSize)
	// resultCh := make(chan [][]string, 10) // Buffer result channel

	// // Start workers
	// for i := 0; i < concurrencyLimit; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		var adsTxtParserList [][]string

	// 		for domain := range domainCh {
	// 			adstxtSingleList := []string{domain}
	// 			adsTxtPage, _, err := CrawlDomain(domain, parserType)
	// 			if err != nil {
	// 				log.Printf("Error crawling domain %s: %v", domain, err)
	// 				continue
	// 			}
	// 			presenceList := IsAdsTxtLinePresent(domain, string(adsTxtPage))
	// 			adstxtSingleList = append(adstxtSingleList, presenceList...)
	// 			adsTxtParserList = append(adsTxtParserList, adstxtSingleList)
	// 			presentDomains = append(presentDomains, domain)
	// 		}

	// 		resultCh <- adsTxtParserList
	// 	}()
	// }

	// // Send domains to workers in batches
	// batch := make([]string, 0, batchSize)
	// for _, domain := range domainsList {
	// 	batch = append(batch, domain)
	// 	if len(batch) == batchSize {
	// 		for _, domain := range batch {
	// 			domainCh <- domain
	// 		}
	// 		batch = batch[:0]
	// 	}
	// }

	// // Send the remaining domains if any
	// for _, domain := range batch {
	// 	domainCh <- domain
	// }

	// close(domainCh)

	// // Wait for workers to finish and collect results
	// go func() {
	// 	wg.Wait()
	// 	close(resultCh)
	// }()

	// // Collect results from workers
	// var adsTxtParserList [][]string
	// for result := range resultCh {
	// 	adsTxtParserList = append(adsTxtParserList, result...)
	// }

	// // Display the data in adsTxtParserList
	// fmt.Println("AdsTxtParserList:")
	// for _, record := range adsTxtParserList {
	// 	fmt.Println(record)
	// }
	// // Display the data in adsTxtParserList
	// fmt.Println("Domains :")
	// for _, record := range presentDomains {
	// 	fmt.Println(record)
	// }
}

func ReadAdsTxtDemandLines() []string {
	var adsTxtDemandLines []string
	fmt.Println("Fetching ads.txt demand lines from adstxt_lines.txt file...")

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return []string{}
	}

	adsTxtLinesList := utils.ReadLinesFromFile(dir + "/resources/domains/adstxt_lines.txt")

	for _, adsTxtLine := range adsTxtLinesList {
		line := strings.TrimSpace(strings.ReplaceAll(strings.ToLower(adsTxtLine), " ", ""))
		adsTxtDemandLines = append(adsTxtDemandLines, line)
	}

	return adsTxtDemandLines
}

func IsAdsTxtLinePresent(domain, adsTxtPage string, adsTxtLinesList []string) (bool, []string) {

	pageResponse := strings.ToLower(strings.ReplaceAll(adsTxtPage, " ", ""))
	pageResponse = strings.ReplaceAll(pageResponse, "\xa0", "")

	present := make([]string, 0)

	for _, adsTxtLine := range adsTxtLinesList {
		searchLine := strings.ReplaceAll(strings.ToLower(adsTxtLine), " ", "")
		if strings.Contains(pageResponse, searchLine) {
			present = append(present, adsTxtLine)
			presentDomains = append(presentDomains, domain)
		}
	}
	return true, present
}

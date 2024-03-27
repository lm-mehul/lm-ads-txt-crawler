package service

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/utils"
)

func FetchLemmaDirectsAndResellerInventory(db *sql.DB, parserType string) {

	fmt.Println("Executing adstxt parser...")
	fmt.Println("Fetching domains from Database...")

	domainsList, pageType, err := FetchDomains(db, parserType)
	if err != nil {
		log.Printf("Error fetching domains from database: %v", err)
		return
	}

	// Define batch size and concurrency limit
	batchSize := 1000      // Adjust as needed
	concurrencyLimit := 10 // Adjust as needed

	var wg sync.WaitGroup
	domainCh := make(chan string, batchSize)
	resultCh := make(chan [][]string, 10) // Buffer result channel

	// Start workers
	for i := 0; i < concurrencyLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var adsTxtParserList [][]string

			for domain := range domainCh {
				adstxtSingleList := []string{domain}
				adsTxtPage, err := crawlDomain(domain, pageType)
				if err != nil {
					log.Printf("Error crawling domain %s: %v", domain, err)
					continue
				}
				hash := utils.GenerateHash(adsTxtPage)
				status, err := models.IsDomainCrawled(domain, string(hash), db)
				if nil != err {
					log.Printf("Error checking domain hash from DB for domain : %s: %v", domain, err)
					continue
				}
				if status {
					log.Printf("Domain already crawled and no changes were there in ads txt page : %v", domain)
					continue
				}
				presenceList := lemmaDirectsAndResellerInventory(string(adsTxtPage))
				adstxtSingleList = append(adstxtSingleList, presenceList...)
				adsTxtParserList = append(adsTxtParserList, adstxtSingleList)
			}

			resultCh <- adsTxtParserList
		}()
	}

	// Send domains to workers in batches
	batch := make([]string, 0, batchSize)
	for _, domain := range domainsList {
		batch = append(batch, domain)
		if len(batch) == batchSize {
			for _, domain := range batch {
				domainCh <- domain
			}
			batch = batch[:0]
		}
	}
	close(domainCh)

	// Wait for workers to finish and collect results
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results from workers
	var adsTxtParserList [][]string
	for result := range resultCh {
		adsTxtParserList = append(adsTxtParserList, result...)
	}

	// Display the data in adsTxtParserList
	fmt.Println("Domains :")
	for _, record := range adsTxtParserList {
		fmt.Println(record)
	}

}

func lemmaDirectsAndResellerInventory(adsTxtPage string) []string {

	present := make([]string, 0)

	lemmaDirectPubs := make(map[string]struct{})
	lemmaResellerPubs := make(map[string]struct{})

	pageResponse := strings.ToLower(strings.ReplaceAll(adsTxtPage, " ", ""))
	pageResponse = strings.ReplaceAll(pageResponse, "\xa0", "")

	adstxtPageLineSplits := strings.Split(pageResponse, "\n")

	for _, line := range adstxtPageLineSplits {
		if strings.Contains(line, "lemmatechnologies.com") {
			parts := strings.Split(line, ",")
			if len(parts) > 1 {
				publisherID := strings.TrimSpace(parts[1])
				if strings.Contains(line, "direct") {
					lemmaDirectPubs[publisherID] = struct{}{}
				} else if strings.Contains(line, "reseller") {
					lemmaResellerPubs[publisherID] = struct{}{}
				}
			}
		}
	}

	lemmaDirectPubsStr, lemmaResellerPubsStr := mapToStrings(lemmaDirectPubs), mapToStrings(lemmaResellerPubs)
	present = append(present, lemmaDirectPubsStr, lemmaResellerPubsStr)
	return present
}

func mapToStrings(m map[string]struct{}) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

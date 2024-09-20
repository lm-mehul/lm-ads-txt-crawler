package handler

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/logger"
	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/repository"
	"github.com/lemmamedia/ads-txt-crawler/service"
	"github.com/lemmamedia/ads-txt-crawler/service/parsers"
)

var AdsTxtDemandLines []string

type DemandLinesResult struct {
	DemandLines []models.DemandLinesEntry
}

func workerDemandLines(db *sql.DB, jobs <-chan models.BundleInfo, results chan<- DemandLinesResult) {
	for bundle := range jobs {
		var demandLines []models.DemandLinesEntry

		// Process the bundle based on its category
		switch bundle.Category {
		case constant.BUNDLE_MOBILE_ANDROID:
			fetchedBundle, err := parsers.ProcessAndroidBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
			} else {
				// Process crawled bundle
				if fetchedBundle.Domain != "" {
					demand := processFetchedBundle(db, fetchedBundle, bundle)
					demandLines = append(demandLines, demand...)

				}
			}
		case constant.BUNDLE_MOBILE_IOS:
			fetchedBundle, err := parsers.ProcessIOSBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_IOS, err.Error())
			} else {
				if fetchedBundle.Domain != "" {
					demand := processFetchedBundle(db, fetchedBundle, bundle)
					demandLines = append(demandLines, demand...)
				}
			}
		case constant.BUNDLE_CTV:
			fetchedBundle, err := parsers.ProcessCTVBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_CTV, err.Error())
			} else {
				if fetchedBundle.Domain != "" {
					demand := processFetchedBundle(db, fetchedBundle, bundle)
					demandLines = append(demandLines, demand...)
				}
			}
		case constant.BUNDLE_WEB:
			fetchedBundle, err := parsers.ProcessWebBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_WEB, err.Error())
			} else {
				if fetchedBundle.Domain != "" {
					demand := processFetchedBundle(db, fetchedBundle, bundle)
					demandLines = append(demandLines, demand...)
				}
			}
		default:
			logger.Info(bundle.Bundle, constant.BUNDLE_CTV, "Invalid bundle category")
		}

		results <- DemandLinesResult{
			DemandLines: demandLines,
		}
	}
}

func processFetchedBundle(db *sql.DB, fetchedBundle models.BundleInfo, bundle models.BundleInfo) []models.DemandLinesEntry {

	var demandLines []models.DemandLinesEntry

	// Crawling and processing logic
	adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
	if err == nil {
		presenceList := service.IsAdsTxtLinePresent(string(adsTxtPage), AdsTxtDemandLines)

		if len(presenceList) > 0 {
			for _, demandLine := range presenceList {
				adsTxtDemand := models.DemandLinesEntry{
					Bundle:     fetchedBundle.Bundle,
					Domain:     fetchedBundle.Domain,
					Category:   fetchedBundle.Category,
					AdsPageURL: url,
					PageType:   constant.ADS_TXT_pageType,
					DemandLine: demandLine,
				}
				demandLines = append(demandLines, adsTxtDemand)
			}
		}
	} else {
		logger.Info(bundle.Bundle, constant.FAILED_DOMAIN_CRAWLING, err.Error())
	}

	appAdsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.APP_ADS_TXT_pageType)
	if err == nil {
		presenceList := service.IsAdsTxtLinePresent(string(appAdsTxtPage), AdsTxtDemandLines)

		if len(presenceList) > 0 {
			for _, demandLine := range presenceList {
				adsTxtDemand := models.DemandLinesEntry{
					Bundle:     fetchedBundle.Bundle,
					Domain:     fetchedBundle.Domain,
					Category:   fetchedBundle.Category,
					AdsPageURL: url,
					PageType:   constant.APP_ADS_TXT_pageType,
					DemandLine: demandLine,
				}
				demandLines = append(demandLines, adsTxtDemand)
			}
		}
	} else {
		logger.Info(bundle.Bundle, constant.FAILED_DOMAIN_CRAWLING, err.Error())
	}

	return demandLines
}

func FetchDemandLinesInventory(db *sql.DB) {

	fmt.Printf("---------------------------------------------------------------------------------\n")
	fmt.Printf("Demand Lines fetching started...\n")
	fmt.Printf("---------------------------------------------------------------------------------\n")

	AdsTxtDemandLines = service.ReadAdsTxtDemandLines()

	fmt.Printf("\nAds.txt Demand Lines loaded successfully in memory...\n%v\n", AdsTxtDemandLines)

	fmt.Printf("---------------------------------------------------------------------------------\n")
	fmt.Printf("Demand Lines loaded successfully in memory...\n")
	fmt.Printf("---------------------------------------------------------------------------------\n")

	const batchSize = 1000
	tempBundles := models.PopulateSampleBundles()

	// tempBundles, err := repository.GetBundlesFromDB(db, 0, 0)
	// if err != nil {
	// 	logger.Error("Error fetching bundles from DB: %v", err)
	// 	return
	// }

	totalBundles := len(tempBundles)

	jobs := make(chan models.BundleInfo, batchSize)
	results := make(chan DemandLinesResult, batchSize)
	var wg sync.WaitGroup

	// Start worker pool
	numWorkers := 50 // Number of workers can be adjusted based on system capability
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerDemandLines(db, jobs, results)
		}()
	}

	// Distribute jobs to workers
	go func() {
		for i := 0; i < totalBundles; i++ {
			jobs <- tempBundles[i]
		}
		close(jobs)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregating results
	var demandLinesCount int
	var allDemandLines []models.DemandLinesEntry

	for result := range results {
		demandLinesCount += len(result.DemandLines)
		allDemandLines = append(allDemandLines, result.DemandLines...)
	}

	if len(allDemandLines) > 0 {
		err := repository.SaveDemandLinesResultInDB(db, allDemandLines)
		if err != nil {
			logger.Error("Error saving demand lines in DB: %v", err)
		}
	}

	// Print summary

	fmt.Printf("---------------------------------------------------------------------------------\n")
	fmt.Printf("Total bundles: %d\n", totalBundles)
	fmt.Printf("Total demand lines: %d\n", demandLinesCount)
	fmt.Printf("\n\n---------------------------------------------------------------------------------\n")
	fmt.Printf("Total Request Timeout Count: %d\n", constant.RequestTimeoutCount)
	fmt.Printf("---------------------------------------------------------------------------------\n")
}

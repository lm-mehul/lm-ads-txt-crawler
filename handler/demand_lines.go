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
	"github.com/lemmamedia/ads-txt-crawler/utils"
)

type DemandLinesResult struct {
	CrawledBundles []models.BundleInfo
	FailedBundles  []models.BundleInfo
	DemandLines    []models.DemandLinesEntry
}

func workerDemandLines(db *sql.DB, jobs <-chan models.BundleInfo, results chan<- DemandLinesResult) {
	for bundle := range jobs {
		var crawledBundles []models.BundleInfo
		var failedBundles []models.BundleInfo
		var demandLines []models.DemandLinesEntry

		// Process the bundle based on its category
		switch bundle.Category {
		case constant.BUNDLE_MOBILE_ANDROID:
			fetchedBundle, err := parsers.ProcessAndroidBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
				failedBundles = append(failedBundles, bundle)
			} else {
				// Process crawled bundle
				if fetchedBundle.Domain != "" {
					newBundle, demand, failed := processFetchedBundle(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					demandLines = append(demandLines, demand...)
					crawledBundles = append(crawledBundles, newBundle)
					failedBundles = append(failedBundles, failed...)
				}
			}
		case constant.BUNDLE_MOBILE_IOS:
			fetchedBundle, err := parsers.ProcessIOSBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_IOS, err.Error())
				failedBundles = append(failedBundles, bundle)
			} else {
				if fetchedBundle.Domain != "" {
					newBundle, demand, failed := processFetchedBundle(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					demandLines = append(demandLines, demand...)
					crawledBundles = append(crawledBundles, newBundle)
					failedBundles = append(failedBundles, failed...)
				}
			}
		case constant.BUNDLE_CTV:
			fetchedBundle, err := parsers.ProcessCTVBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_CTV, err.Error())
				failedBundles = append(failedBundles, bundle)
			} else {
				if fetchedBundle.Domain != "" {
					newBundle, demand, failed := processFetchedBundle(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					demandLines = append(demandLines, demand...)
					crawledBundles = append(crawledBundles, newBundle)
					failedBundles = append(failedBundles, failed...)
				}
			}
		case constant.BUNDLE_WEB:
			fetchedBundle, err := parsers.ProcessWebBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_WEB, err.Error())
				failedBundles = append(failedBundles, bundle)
			} else {
				if fetchedBundle.Domain != "" {
					newBundle, demand, failed := processFetchedBundle(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					demandLines = append(demandLines, demand...)
					crawledBundles = append(crawledBundles, newBundle)
					failedBundles = append(failedBundles, failed...)
				}
			}
		default:
			logger.Info(bundle.Bundle, constant.BUNDLE_CTV, "Invalid bundle category")
			failedBundles = append(failedBundles, bundle)
		}

		results <- DemandLinesResult{
			CrawledBundles: crawledBundles,
			FailedBundles:  failedBundles,
			DemandLines:    demandLines,
		}
	}
}

func processFetchedBundle(db *sql.DB, fetchedBundle models.BundleInfo, bundle models.BundleInfo) (models.BundleInfo, []models.DemandLinesEntry, []models.BundleInfo) {

	var demandLines []models.DemandLinesEntry
	var failedBundles []models.BundleInfo
	isCrawled := 0

	// Crawling and processing logic
	adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
	if err == nil {
		isCrawled++
		// presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))
		// if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
		// 	presenceList.Bundle = fetchedBundle.Bundle
		// 	presenceList.Category = fetchedBundle.Category
		// 	presenceList.AdsPageURL = url
		// 	presenceList.PageType = constant.ADS_TXT_pageType
		// 	demandLines = append(demandLines, presenceList)
		// }
		bundle.AdsTxtURL = url
		bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)
	} else {
		logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
	}

	appAdsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.APP_ADS_TXT_pageType)
	if err == nil {
		isCrawled++
		// presenceList := service.LemmaDirectsAndResellerInventory(string(appAdsTxtPage))
		// if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
		// 	presenceList.Bundle = fetchedBundle.Bundle
		// 	presenceList.Category = fetchedBundle.Category
		// 	presenceList.AdsPageURL = url
		// 	presenceList.PageType = constant.APP_ADS_TXT_pageType
		// 	demandLines = append(demandLines, presenceList)
		// }
		bundle.AppAdsTxtURL = url
		bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
	} else {
		logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
	}

	if isCrawled == 0 {
		failedBundles = append(failedBundles, bundle)
	}

	return bundle, demandLines, failedBundles
}

func FetchDemandLinesInventory(db *sql.DB) {
	const batchSize = 1000
	tempBundles := models.PopulateSampleBundles()
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
	var crawledBundlesCount, failedBundlesCount, demandLinesCount int
	var allCrawledBundles, allFailedBundles []models.BundleInfo
	var allDemandLines []models.DemandLinesEntry

	for result := range results {
		crawledBundlesCount += len(result.CrawledBundles)
		failedBundlesCount += len(result.FailedBundles)
		demandLinesCount += len(result.DemandLines)

		allCrawledBundles = append(allCrawledBundles, result.CrawledBundles...)
		allFailedBundles = append(allFailedBundles, result.FailedBundles...)
		allDemandLines = append(allDemandLines, result.DemandLines...)
	}

	// Save results to DB
	if len(allCrawledBundles) > 0 {
		err := repository.SaveCrawledBundlesInDB(db, allCrawledBundles)
		if err != nil {
			logger.Error("Error saving crawled bundles in DB: %v", err)
		}
	}

	if len(allFailedBundles) > 0 {
		err := repository.SaveFailedBundlesInDB(db, allFailedBundles)
		if err != nil {
			logger.Error("Error saving failed bundles in DB: %v", err)
		}
	}

	// if len(allDemandLines) > 0 {
	// 	err := repository.SaveLemmaEntriesInDB(db, allDemandLines)
	// 	if err != nil {
	// 		logger.Error("Error saving demand lines in DB: %v", err)
	// 	}
	// }

	// Print summary

	fmt.Printf("---------------------------------------------------------------------------------\n")
	fmt.Printf("Total bundles: %d\n", totalBundles)
	fmt.Printf("Total crawled bundles: %d\n", crawledBundlesCount)
	fmt.Printf("Total failed bundles: %d\n", failedBundlesCount)
	fmt.Printf("Total demand lines: %d\n", demandLinesCount)
	fmt.Printf("\n\n---------------------------------------------------------------------------------\n")
	fmt.Printf("Total Request Timeout Count: %d\n", constant.RequestTimeoutCount)
	fmt.Printf("---------------------------------------------------------------------------------\n")
}

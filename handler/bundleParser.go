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

type BundleParserResult struct {
	CrawledBundles []models.BundleInfo
	FailedBundles  []models.BundleInfo
}

func workerbundleParser(db *sql.DB, jobs <-chan models.BundleInfo, results chan<- BundleParserResult) {
	for bundle := range jobs {
		var crawledBundles []models.BundleInfo
		var failedBundles []models.BundleInfo

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
					newBundle, failed := processBundleParserRequests(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website

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
					newBundle, failed := processBundleParserRequests(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website

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
					newBundle, failed := processBundleParserRequests(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website

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
					newBundle, failed := processBundleParserRequests(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website

					crawledBundles = append(crawledBundles, newBundle)
					failedBundles = append(failedBundles, failed...)
				}
			}
		default:
			logger.Info(bundle.Bundle, constant.BUNDLE_CTV, "Invalid bundle category")
			failedBundles = append(failedBundles, bundle)
		}

		results <- BundleParserResult{
			CrawledBundles: crawledBundles,
			FailedBundles:  failedBundles,
		}
	}
}

func processBundleParserRequests(db *sql.DB, fetchedBundle models.BundleInfo, bundle models.BundleInfo) (models.BundleInfo, []models.BundleInfo) {

	var failedBundles []models.BundleInfo
	isCrawled := 0

	// Crawling and processing logic
	adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
	if err == nil {
		isCrawled++
		bundle.AdsTxtURL = url
		bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)
	} else {
		logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
	}

	appAdsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.APP_ADS_TXT_pageType)
	if err == nil {
		isCrawled++
		bundle.AppAdsTxtURL = url
		bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
	} else {
		logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
	}

	if isCrawled == 0 {
		failedBundles = append(failedBundles, bundle)
	}

	return bundle, failedBundles
}

func BundleParser(db *sql.DB) {

	fmt.Println("Executing Bundle parser...")

	const batchSize = 1000
	tempBundles := models.PopulateSampleBundles()
	totalBundles := len(tempBundles)

	jobs := make(chan models.BundleInfo, batchSize)
	results := make(chan BundleParserResult, batchSize)
	var wg sync.WaitGroup

	// Start worker pool
	numWorkers := 50 // Number of workers can be adjusted based on system capability
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerbundleParser(db, jobs, results)
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

	fmt.Printf("Bundles are cralwed. Please wait for the results to be stored in database...\n")

	// Aggregating results
	var crawledBundlesCount, failedBundlesCount int
	var allCrawledBundles, allFailedBundles []models.BundleInfo

	for result := range results {
		crawledBundlesCount += len(result.CrawledBundles)
		failedBundlesCount += len(result.FailedBundles)

		allCrawledBundles = append(allCrawledBundles, result.CrawledBundles...)
		allFailedBundles = append(allFailedBundles, result.FailedBundles...)
	}

	// // Save results to DB
	// if len(allCrawledBundles) > 0 {
	// 	err := repository.SaveCrawledBundlesInDB(db, allCrawledBundles)
	// 	if err != nil {
	// 		logger.Error("Error saving crawled bundles in DB: %v", err)
	// 	}
	// }

	// if len(allFailedBundles) > 0 {
	// 	err := repository.SaveFailedBundlesInDB(db, allFailedBundles)
	// 	if err != nil {
	// 		logger.Error("Error saving failed bundles in DB: %v", err)
	// 	}
	// }

	// Define your batch size
	const dbBatchSize = 1000

	// Save crawled bundles in batches
	if len(allCrawledBundles) > 0 {
		models.BatchSave(db, allCrawledBundles, batchSize, repository.SaveCrawledBundlesInDB, "crawled bundles")
	}

	// Save failed bundles in batches
	if len(allFailedBundles) > 0 {
		models.BatchSave(db, allFailedBundles, batchSize, repository.SaveFailedBundlesInDB, "failed bundles")
	}

	fmt.Println("All parsers have finished.")

	// Print summary

	fmt.Printf("---------------------------------------------------------------------------------\n")
	fmt.Printf("Total bundles: %d\n", totalBundles)
	fmt.Printf("Total crawled bundles: %d\n", crawledBundlesCount)
	fmt.Printf("Total failed bundles: %d\n", failedBundlesCount)
	fmt.Printf("\n\n---------------------------------------------------------------------------------\n")
	fmt.Printf("Total Request Timeout Count: %d\n", constant.RequestTimeoutCount)
	fmt.Printf("---------------------------------------------------------------------------------\n")
}

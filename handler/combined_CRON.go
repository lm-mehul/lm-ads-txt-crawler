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

type CombinedCRONResult struct {
	CrawledBundles []models.BundleInfo
	FailedBundles  []models.BundleInfo
	LemmaLines     []models.LemmaEntry
	DemandLines    []models.DemandLinesEntry
}

func workercombinedLines(db *sql.DB, jobs <-chan models.BundleInfo, results chan<- CombinedCRONResult) {
	for bundle := range jobs {
		var crawledBundles []models.BundleInfo
		var failedBundles []models.BundleInfo
		var lemmaLines []models.LemmaEntry
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
					newBundle, lemma, failed, demand := processFetchedBundlesForCombinedCRON(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					lemmaLines = append(lemmaLines, lemma...)
					crawledBundles = append(crawledBundles, newBundle)
					failedBundles = append(failedBundles, failed...)
					demandLines = append(demandLines, demand...)
				}
			}
		case constant.BUNDLE_MOBILE_IOS:
			fetchedBundle, err := parsers.ProcessIOSBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_IOS, err.Error())
				failedBundles = append(failedBundles, bundle)
			} else {
				if fetchedBundle.Domain != "" {
					newBundle, lemma, failed, demand := processFetchedBundlesForCombinedCRON(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					lemmaLines = append(lemmaLines, lemma...)
					crawledBundles = append(crawledBundles, newBundle)
					failedBundles = append(failedBundles, failed...)
					demandLines = append(demandLines, demand...)
				}
			}
		case constant.BUNDLE_CTV:
			fetchedBundle, err := parsers.ProcessCTVBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_CTV, err.Error())
				failedBundles = append(failedBundles, bundle)
			} else {
				if fetchedBundle.Domain != "" {
					newBundle, lemma, failed, demand := processFetchedBundlesForCombinedCRON(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					lemmaLines = append(lemmaLines, lemma...)
					crawledBundles = append(crawledBundles, newBundle)
					failedBundles = append(failedBundles, failed...)
					demandLines = append(demandLines, demand...)
				}
			}
		case constant.BUNDLE_WEB:
			fetchedBundle, err := parsers.ProcessWebBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_WEB, err.Error())
				failedBundles = append(failedBundles, bundle)
			} else {
				if fetchedBundle.Domain != "" {
					newBundle, lemma, failed, demand := processFetchedBundlesForCombinedCRON(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					lemmaLines = append(lemmaLines, lemma...)
					crawledBundles = append(crawledBundles, newBundle)
					failedBundles = append(failedBundles, failed...)
					demandLines = append(demandLines, demand...)
				}
			}
		default:
			logger.Info(bundle.Bundle, constant.BUNDLE_CTV, "Invalid bundle category")
			failedBundles = append(failedBundles, bundle)
		}

		results <- CombinedCRONResult{
			CrawledBundles: crawledBundles,
			FailedBundles:  failedBundles,
			LemmaLines:     lemmaLines,
			DemandLines:    demandLines,
		}
	}
}

func processFetchedBundlesForCombinedCRON(db *sql.DB, fetchedBundle models.BundleInfo, bundle models.BundleInfo) (models.BundleInfo, []models.LemmaEntry, []models.BundleInfo, []models.DemandLinesEntry) {

	var lemmaLines []models.LemmaEntry
	var failedBundles []models.BundleInfo
	var demandLines []models.DemandLinesEntry

	isCrawled := 0

	// Crawling and processing logic
	adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
	if err == nil {
		isCrawled++
		presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

		if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
			presenceList.Bundle = fetchedBundle.Bundle
			presenceList.Category = fetchedBundle.Category
			presenceList.AdsPageURL = url
			presenceList.PageType = constant.ADS_TXT_pageType

			lemmaLines = append(lemmaLines, presenceList)
		}
		bundle.AdsTxtURL = url
		bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)

		adsTxtPresenceList := service.IsAdsTxtLinePresent(string(adsTxtPage), AdsTxtDemandLines)

		if len(adsTxtPresenceList) > 0 {
			for _, demandLine := range adsTxtPresenceList {
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
		isCrawled++

		presenceList := service.LemmaDirectsAndResellerInventory(string(appAdsTxtPage))
		if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
			presenceList.Bundle = fetchedBundle.Bundle
			presenceList.Category = fetchedBundle.Category
			presenceList.AdsPageURL = url
			presenceList.PageType = constant.APP_ADS_TXT_pageType
			lemmaLines = append(lemmaLines, presenceList)
		}
		bundle.AppAdsTxtURL = url
		bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)

		appAdsTxtPresenceList := service.IsAdsTxtLinePresent(string(appAdsTxtPage), AdsTxtDemandLines)

		if len(appAdsTxtPresenceList) > 0 {
			for _, demandLine := range appAdsTxtPresenceList {
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

	if isCrawled == 0 {
		failedBundles = append(failedBundles, bundle)
	}

	return bundle, lemmaLines, failedBundles, demandLines
}

func ScheduleCombinedCRON(db *sql.DB) {

	repository.ClearTableData(db, "crawled_bundles")
	repository.ClearTableData(db, "failed_bundles")
	repository.ClearTableData(db, "lemma_entries")
	repository.ClearTableData(db, "bundle_demand_lines")

	AdsTxtDemandLines = service.ReadAdsTxtDemandLines()

	fmt.Printf("---------------------------------------------------------------------------------\n")
	fmt.Printf("Fetching ScheduleCombinedCRON...\n")
	fmt.Printf("---------------------------------------------------------------------------------\n")

	const batchSize = 1000

	// Fetch bundles from DB
	// tempBundles := models.PopulateSampleBundles()

	tempBundles, err := repository.GetBundlesFromDB(db, 0, 0)
	if err != nil {
		logger.Error("Error fetching bundles from DB: %v", err)
		return
	}

	totalBundles := len(tempBundles)

	jobs := make(chan models.BundleInfo, batchSize)
	results := make(chan CombinedCRONResult, batchSize)
	var wg sync.WaitGroup

	// Start worker pool
	numWorkers := 10 // Number of workers can be adjusted based on system capability
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workercombinedLines(db, jobs, results)
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
	var crawledBundlesCount, failedBundlesCount, lemmaLinesCount int
	var allCrawledBundles, allFailedBundles []models.BundleInfo
	var allLemmaLines []models.LemmaEntry
	var allDemandLines []models.DemandLinesEntry

	for result := range results {
		crawledBundlesCount += len(result.CrawledBundles)
		failedBundlesCount += len(result.FailedBundles)
		lemmaLinesCount += len(result.LemmaLines)
		allCrawledBundles = append(allCrawledBundles, result.CrawledBundles...)
		allFailedBundles = append(allFailedBundles, result.FailedBundles...)
		allLemmaLines = append(allLemmaLines, result.LemmaLines...)
		allDemandLines = append(allDemandLines, result.DemandLines...)
	}

	fmt.Printf("Bundle crawling completed. Please wait for the results to be stored in the database...\n")

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

	// Save lemma entries in batches
	if len(allLemmaLines) > 0 {
		models.BatchSave(db, allLemmaLines, batchSize, repository.SaveLemmaEntriesInDB, "lemma lines inventory")
	}

	// Save demand line entries in batches
	if len(allDemandLines) > 0 {
		models.BatchSave(db, allDemandLines, batchSize, repository.SaveDemandLinesResultInDB, "Demand Lines inventory")
	}

	// Print summary

	fmt.Printf("---------------------------------------------------------------------------------\n")
	fmt.Printf("Total bundles: %d\n", totalBundles)
	fmt.Printf("Total crawled bundles: %d\n", crawledBundlesCount)
	fmt.Printf("Total failed bundles: %d\n", failedBundlesCount)
	fmt.Printf("Total lemma lines: %d\n", lemmaLinesCount)
	fmt.Printf("\n\n---------------------------------------------------------------------------------\n")
	fmt.Printf("Total Request Timeout Count: %d\n", constant.RequestTimeoutCount)
	fmt.Printf("---------------------------------------------------------------------------------\n")
}

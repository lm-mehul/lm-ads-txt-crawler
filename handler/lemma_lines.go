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

type LemmaLinesResult struct {
	CrawledBundles []models.BundleInfo
	FailedBundles  []models.BundleInfo
	LemmaLines     []models.LemmaEntry
}

func workerLemmaLines(db *sql.DB, jobs <-chan models.BundleInfo, results chan<- LemmaLinesResult) {
	for bundle := range jobs {
		var crawledBundles []models.BundleInfo
		var failedBundles []models.BundleInfo
		var lemmaLines []models.LemmaEntry

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
					newBundle, lemma, failed := processFetchedBundlesForLemmaLines(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					lemmaLines = append(lemmaLines, lemma...)
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
					newBundle, lemma, failed := processFetchedBundlesForLemmaLines(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					lemmaLines = append(lemmaLines, lemma...)
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
					newBundle, lemma, failed := processFetchedBundlesForLemmaLines(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					lemmaLines = append(lemmaLines, lemma...)
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
					newBundle, lemma, failed := processFetchedBundlesForLemmaLines(db, fetchedBundle, bundle)

					newBundle.Domain = fetchedBundle.Domain
					newBundle.Website = fetchedBundle.Website
					lemmaLines = append(lemmaLines, lemma...)
					crawledBundles = append(crawledBundles, newBundle)
					failedBundles = append(failedBundles, failed...)
				}
			}
		default:
			logger.Info(bundle.Bundle, constant.BUNDLE_CTV, "Invalid bundle category")
			failedBundles = append(failedBundles, bundle)
		}

		results <- LemmaLinesResult{
			CrawledBundles: crawledBundles,
			FailedBundles:  failedBundles,
			LemmaLines:     lemmaLines,
		}
	}
}

func processFetchedBundlesForLemmaLines(db *sql.DB, fetchedBundle models.BundleInfo, bundle models.BundleInfo) (models.BundleInfo, []models.LemmaEntry, []models.BundleInfo) {
	var lemmaLines []models.LemmaEntry
	var failedBundles []models.BundleInfo
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
	} else {
		logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
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
	} else {
		logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
	}

	if isCrawled == 0 {
		failedBundles = append(failedBundles, bundle)
	}

	return bundle, lemmaLines, failedBundles
}

func FetchLemmaDirectsAndResellerInventory(db *sql.DB) {

	fmt.Printf("---------------------------------------------------------------------------------\n")
	fmt.Printf("Fetching lemma directs and reseller inventory...\n")
	fmt.Printf("---------------------------------------------------------------------------------\n")

	const batchSize = 1000

	// Fetch bundles from DB
	tempBundles := models.PopulateSampleBundles()

	// tempBundles, err := repository.GetBundlesFromDB(db, 0, 0)
	// if err != nil {
	// 	logger.Error("Error fetching bundles from DB: %v", err)
	// 	return
	// }

	totalBundles := len(tempBundles)

	jobs := make(chan models.BundleInfo, batchSize)
	results := make(chan LemmaLinesResult, batchSize)
	var wg sync.WaitGroup

	// Start worker pool
	numWorkers := 3 // Number of workers can be adjusted based on system capability
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workerLemmaLines(db, jobs, results)
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

	for result := range results {
		crawledBundlesCount += len(result.CrawledBundles)
		failedBundlesCount += len(result.FailedBundles)
		lemmaLinesCount += len(result.LemmaLines)
		allCrawledBundles = append(allCrawledBundles, result.CrawledBundles...)
		allFailedBundles = append(allFailedBundles, result.FailedBundles...)
		allLemmaLines = append(allLemmaLines, result.LemmaLines...)
	}

	fmt.Printf("Bundle crawling completed. Please wait for the results to be stored in the database...\n")

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

	// if len(allLemmaLines) > 0 {
	// 	err := repository.SaveLemmaEntriesInDB(db, allLemmaLines)
	// 	if err != nil {
	// 		logger.Error("Error saving lemma lines in DB: %v", err)
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

	// Save lemma entries in batches
	if len(allLemmaLines) > 0 {
		models.BatchSave(db, allLemmaLines, batchSize, repository.SaveLemmaEntriesInDB, "lemma lines inventory")
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

// -----------------------------------------------------------------------------

// package handler

// import (
// 	"database/sql"
// 	"fmt"

// 	"github.com/lemmamedia/ads-txt-crawler/constant"
// 	"github.com/lemmamedia/ads-txt-crawler/logger"
// 	"github.com/lemmamedia/ads-txt-crawler/models"
// 	"github.com/lemmamedia/ads-txt-crawler/repository"
// 	"github.com/lemmamedia/ads-txt-crawler/service"
// 	"github.com/lemmamedia/ads-txt-crawler/service/parsers"
// 	"github.com/lemmamedia/ads-txt-crawler/utils"
// )

// func FetchLemmaDirectsAndResellerInventory(db *sql.DB) {

// 	// var tempBundles []models.BundleInfo
// 	const batchSize = 1000
// 	crawledBundlesCount := 0
// 	failedBundlesCount := 0
// 	lemmaLinesCount := 0

// 	var err error

// 	tempBundles := models.PopulateSampleBundles()

// 	// tempBundles, err := repository.GetBundlesFromDB(db, 0, 0)
// 	// if err != nil {
// 	// 	log.Printf("Error fetching bundles from DB: %v", err)
// 	// 	return
// 	// }

// 	// Process bundles in batches
// 	totalBundles := len(tempBundles)
// 	for i := 0; i < totalBundles; i += batchSize {
// 		end := i + batchSize
// 		if end > totalBundles {
// 			end = totalBundles
// 		}

// 		batch := tempBundles[i:end]

// 		var crawledBundles []models.BundleInfo
// 		var failedBundles []models.BundleInfo
// 		var lemmaLines []models.LemmaEntry

// 		for _, bundle := range batch {

// 			// status, err := models.IsDomainCrawled(domain, string(hash), db)
// 			// if nil != err {
// 			// 	log.Printf("Error checking domain hash from DB for domain : %s: %v", domain, err)
// 			// 	continue
// 			// }
// 			// if status {
// 			// 	log.Printf("Domain already crawled and no changes were there in ads txt page : %v", domain)
// 			// 	continue
// 			// }

// 			isCrawled := 0

// 			switch bundle.Category {
// 			case constant.BUNDLE_MOBILE_ANDROID:
// 				fetchedBundle, err := parsers.ProcessAndroidBundle(db, bundle.Bundle)
// 				if err != nil {
// 					logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
// 					failedBundles = append(failedBundles, bundle)
// 				} else {
// 					if fetchedBundle.Domain != "" {
// 						adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
// 						if err != nil {
// 							logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
// 						} else {
// 							isCrawled++

// 							presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

// 							if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
// 								presenceList.Bundle = fetchedBundle.Bundle
// 								presenceList.Category = fetchedBundle.Category
// 								presenceList.AdsPageURL = url
// 								presenceList.PageType = constant.ADS_TXT_pageType
// 								lemmaLines = append(lemmaLines, presenceList)
// 							}

// 							bundle.AdsTxtURL = url
// 							bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)
// 						}

// 						appAdsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.APP_ADS_TXT_pageType)
// 						if err != nil {
// 							logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
// 						} else {
// 							isCrawled++

// 							presenceList := service.LemmaDirectsAndResellerInventory(string(appAdsTxtPage))

// 							if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
// 								presenceList.Bundle = fetchedBundle.Bundle
// 								presenceList.Category = fetchedBundle.Category
// 								presenceList.AdsPageURL = url
// 								presenceList.PageType = constant.APP_ADS_TXT_pageType
// 								lemmaLines = append(lemmaLines, presenceList)
// 							}

// 							bundle.AppAdsTxtURL = url
// 							bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
// 						}

// 						if isCrawled > 0 {
// 							bundle.Domain = fetchedBundle.Domain
// 							bundle.Website = fetchedBundle.Website
// 							crawledBundles = append(crawledBundles, bundle)
// 						} else {
// 							failedBundles = append(failedBundles, bundle)
// 						}
// 					} else {
// 						failedBundles = append(failedBundles, bundle)
// 					}
// 				}

// 			case constant.BUNDLE_MOBILE_IOS:
// 				fetchedBundle, err := parsers.ProcessIOSBundle(db, bundle.Bundle)
// 				if err != nil {
// 					logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_IOS, err.Error())
// 					failedBundles = append(failedBundles, bundle)
// 				} else {
// 					if fetchedBundle.Domain != "" {
// 						adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
// 						if err != nil {
// 							logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_IOS, err.Error())
// 						} else {
// 							isCrawled++

// 							presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

// 							if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
// 								presenceList.Bundle = fetchedBundle.Bundle
// 								presenceList.Category = fetchedBundle.Category
// 								presenceList.AdsPageURL = url
// 								presenceList.PageType = constant.ADS_TXT_pageType
// 								lemmaLines = append(lemmaLines, presenceList)
// 							}

// 							bundle.AdsTxtURL = url
// 							bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)
// 						}

// 						appAdsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.APP_ADS_TXT_pageType)
// 						if err != nil {
// 							logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_IOS, err.Error())
// 						} else {
// 							isCrawled++

// 							presenceList := service.LemmaDirectsAndResellerInventory(string(appAdsTxtPage))

// 							if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
// 								presenceList.Bundle = fetchedBundle.Bundle
// 								presenceList.Category = fetchedBundle.Category
// 								presenceList.AdsPageURL = url
// 								presenceList.PageType = constant.APP_ADS_TXT_pageType
// 								lemmaLines = append(lemmaLines, presenceList)
// 							}

// 							bundle.AppAdsTxtURL = url
// 							bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
// 						}

// 						if isCrawled > 0 {
// 							bundle.Domain = fetchedBundle.Domain
// 							bundle.Website = fetchedBundle.Website
// 							crawledBundles = append(crawledBundles, bundle)
// 						} else {
// 							failedBundles = append(failedBundles, bundle)
// 						}
// 					} else {
// 						failedBundles = append(failedBundles, bundle)
// 					}
// 				}
// 			case constant.BUNDLE_CTV:
// 				fetchedBundle, err := parsers.ProcessCTVBundle(db, bundle.Bundle)
// 				if err != nil {
// 					logger.Info(bundle.Bundle, constant.BUNDLE_CTV, err.Error())
// 					failedBundles = append(failedBundles, bundle)
// 				} else {
// 					if fetchedBundle.Domain != "" {
// 						adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
// 						if err != nil {
// 							logger.Info(bundle.Bundle, constant.BUNDLE_CTV, err.Error())
// 						} else {
// 							isCrawled++

// 							presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

// 							if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
// 								presenceList.Bundle = fetchedBundle.Bundle
// 								presenceList.Category = fetchedBundle.Category
// 								presenceList.AdsPageURL = url
// 								presenceList.PageType = constant.ADS_TXT_pageType
// 								lemmaLines = append(lemmaLines, presenceList)
// 							}

// 							bundle.AdsTxtURL = url
// 							bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)
// 						}

// 						appAdsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.APP_ADS_TXT_pageType)
// 						if err != nil {
// 							logger.Info(bundle.Bundle, constant.BUNDLE_CTV, err.Error())
// 						} else {
// 							isCrawled++

// 							presenceList := service.LemmaDirectsAndResellerInventory(string(appAdsTxtPage))

// 							if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
// 								presenceList.Bundle = fetchedBundle.Bundle
// 								presenceList.Category = fetchedBundle.Category
// 								presenceList.AdsPageURL = url
// 								presenceList.PageType = constant.APP_ADS_TXT_pageType
// 								lemmaLines = append(lemmaLines, presenceList)
// 							}

// 							bundle.AppAdsTxtURL = url
// 							bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
// 						}

// 						if isCrawled > 0 {
// 							bundle.Domain = fetchedBundle.Domain
// 							bundle.Website = fetchedBundle.Website
// 							crawledBundles = append(crawledBundles, bundle)
// 						} else {
// 							failedBundles = append(failedBundles, bundle)
// 						}

// 					} else {
// 						failedBundles = append(failedBundles, bundle)
// 					}
// 				}
// 			case constant.BUNDLE_WEB:

// 				fetchedBundle, err := parsers.ProcessWebBundle(db, bundle.Bundle)
// 				if err != nil {
// 					logger.Info(bundle.Bundle, constant.BUNDLE_CTV, err.Error())
// 					failedBundles = append(failedBundles, bundle)
// 				} else {
// 					if fetchedBundle.Domain != "" {

// 						adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
// 						if err != nil {
// 							logger.Info(bundle.Bundle, constant.BUNDLE_WEB, err.Error())
// 							failedBundles = append(failedBundles, bundle)
// 						} else {

// 							isCrawled++

// 							presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

// 							if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
// 								presenceList.Bundle = fetchedBundle.Bundle
// 								presenceList.Category = fetchedBundle.Category
// 								presenceList.AdsPageURL = url
// 								presenceList.PageType = constant.ADS_TXT_pageType
// 								lemmaLines = append(lemmaLines, presenceList)
// 							}

// 							bundle.AdsTxtURL = url
// 							bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)
// 						}

// 						appAdsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.APP_ADS_TXT_pageType)
// 						if err != nil {
// 							logger.Info(bundle.Bundle, constant.BUNDLE_CTV, err.Error())
// 						} else {
// 							isCrawled++

// 							presenceList := service.LemmaDirectsAndResellerInventory(string(appAdsTxtPage))

// 							if presenceList.LemmaDirect != "" || presenceList.LemmaReseller != "" {
// 								presenceList.Bundle = fetchedBundle.Bundle
// 								presenceList.Category = fetchedBundle.Category
// 								presenceList.AdsPageURL = url
// 								presenceList.PageType = constant.APP_ADS_TXT_pageType
// 								lemmaLines = append(lemmaLines, presenceList)
// 							}

// 							bundle.AppAdsTxtURL = url
// 							bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
// 						}

// 						if isCrawled > 0 {
// 							bundle.Domain = fetchedBundle.Domain
// 							bundle.Website = fetchedBundle.Website
// 							crawledBundles = append(crawledBundles, bundle)
// 						} else {
// 							failedBundles = append(failedBundles, bundle)
// 						}

// 					} else {
// 						failedBundles = append(failedBundles, bundle)
// 					}
// 				}
// 			default:
// 				logger.Info(bundle.Bundle, constant.BUNDLE_CTV, "Invalid bundle category")
// 				failedBundles = append(failedBundles, bundle)
// 			}
// 		}

// 		if len(crawledBundles) > 0 {
// 			err = repository.SaveCrawledBundlesInDB(db, crawledBundles)
// 			if err != nil {
// 				logger.Error("Error saving crawled bundles in DB: %v", err)
// 			}
// 		}

// 		if len(failedBundles) > 0 {
// 			err = repository.SaveFailedBundlesInDB(db, failedBundles)
// 			if err != nil {
// 				logger.Error("Error saving failed bundles in DB: %v", err)
// 			}
// 		}

// 		if len(lemmaLines) > 0 {
// 			err = repository.SaveLemmaEntriesInDB(db, lemmaLines)
// 			if err != nil {
// 				logger.Error("Error saving lemma lines in DB: %v", err)
// 			}
// 		}

// 		crawledBundlesCount = len(crawledBundles)

// 		failedBundlesCount = len(failedBundles)
// 		lemmaLinesCount = len(lemmaLines)
// 	}

// 	fmt.Printf("Total bundles: %d\n", len(tempBundles))
// 	fmt.Printf("Total crawled bundles: %d\n", crawledBundlesCount)
// 	fmt.Printf("Total failed bundles: %d\n", failedBundlesCount)
// 	fmt.Printf("Total lemma lines: %d\n", lemmaLinesCount)
// }

// -----------------------------------------------------------------------------

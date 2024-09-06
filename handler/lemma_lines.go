package handler

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/logger"
	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/repository"
	"github.com/lemmamedia/ads-txt-crawler/service"
	"github.com/lemmamedia/ads-txt-crawler/service/parsers"
	"github.com/lemmamedia/ads-txt-crawler/utils"
)

var tempBundles []models.BundleInfo

func FetchLemmaDirectsAndResellerInventory(db *sql.DB) {

	populateBundles()

	var crawledBundles []models.BundleInfo
	var failedBundles []models.BundleInfo
	var lemmaLines []models.LemmaEntry

	for _, bundle := range tempBundles {

		// status, err := models.IsDomainCrawled(domain, string(hash), db)
		// if nil != err {
		// 	log.Printf("Error checking domain hash from DB for domain : %s: %v", domain, err)
		// 	continue
		// }
		// if status {
		// 	log.Printf("Domain already crawled and no changes were there in ads txt page : %v", domain)
		// 	continue
		// }

		isCrawled := 0

		switch bundle.Category {
		case constant.BUNDLE_MOBILE_ANDROID:
			fetchedBundle, err := parsers.ProcessAndroidBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_ANDROID, err.Error())
				failedBundles = append(failedBundles, bundle)
			} else {
				if fetchedBundle.Domain == "" {
					adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
					if err != nil {
						log.Printf("Error crawling domain %s: %v", fetchedBundle.Domain, err)
					} else {
						isCrawled++
						fmt.Printf("Crawled domain: %s\n", fetchedBundle.Domain)

						presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

						presenceList.Bundle = fetchedBundle.Bundle
						presenceList.Category = fetchedBundle.Category

						lemmaLines = append(lemmaLines, presenceList)

						bundle.AdsTxtURL = url
						bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)
					}

					appAdsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.APP_ADS_TXT_pageType)
					if err != nil {
						log.Printf("Error crawling domain %s: %v", fetchedBundle.Domain, err)
					} else {
						isCrawled++
						fmt.Printf("Crawled domain: %s\n", fetchedBundle.Domain)

						presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

						presenceList.Bundle = fetchedBundle.Bundle
						presenceList.Category = fetchedBundle.Category

						lemmaLines = append(lemmaLines, presenceList)

						bundle.AppAdsTxtURL = url
						bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
					}

					if isCrawled > 0 {
						crawledBundles = append(crawledBundles, bundle)
					} else {
						failedBundles = append(failedBundles, bundle)
					}
				} else {
					failedBundles = append(failedBundles, bundle)
				}
			}

		case constant.BUNDLE_MOBILE_IOS:
			fetchedBundle, err := parsers.ProcessIOSBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_MOBILE_IOS, err.Error())
				failedBundles = append(failedBundles, bundle)
			} else {
				if fetchedBundle.Domain == "" {
					adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
					if err != nil {
						log.Printf("Error crawling domain %s: %v", fetchedBundle.Domain, err)
					} else {
						isCrawled++
						fmt.Printf("Crawled domain: %s\n", fetchedBundle.Domain)

						presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

						presenceList.Bundle = fetchedBundle.Bundle
						presenceList.Category = fetchedBundle.Category

						lemmaLines = append(lemmaLines, presenceList)

						bundle.AdsTxtURL = url
						bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)
					}

					appAdsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.APP_ADS_TXT_pageType)
					if err != nil {
						log.Printf("Error crawling domain %s: %v", fetchedBundle.Domain, err)
					} else {
						isCrawled++
						fmt.Printf("Crawled domain: %s\n", fetchedBundle.Domain)

						presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

						presenceList.Bundle = fetchedBundle.Bundle
						presenceList.Category = fetchedBundle.Category

						lemmaLines = append(lemmaLines, presenceList)

						bundle.AppAdsTxtURL = url
						bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
					}

					if isCrawled > 0 {
						crawledBundles = append(crawledBundles, bundle)
					} else {
						failedBundles = append(failedBundles, bundle)
					}
				} else {
					failedBundles = append(failedBundles, bundle)
				}
			}
		case constant.BUNDLE_CTV:
			fetchedBundle, err := parsers.ProcessCTVBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_CTV, err.Error())
				failedBundles = append(failedBundles, bundle)
			} else {
				if fetchedBundle.Domain == "" {
					adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
					if err != nil {
						log.Printf("Error crawling domain %s: %v", fetchedBundle.Domain, err)
					} else {
						isCrawled++
						fmt.Printf("Crawled domain: %s\n", fetchedBundle.Domain)

						presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

						presenceList.Bundle = fetchedBundle.Bundle
						presenceList.Category = fetchedBundle.Category

						lemmaLines = append(lemmaLines, presenceList)

						bundle.AdsTxtURL = url
						bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)
					}

					appAdsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.APP_ADS_TXT_pageType)
					if err != nil {
						log.Printf("Error crawling domain %s: %v", fetchedBundle.Domain, err)
					} else {
						isCrawled++
						fmt.Printf("Crawled domain: %s\n", fetchedBundle.Domain)

						presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

						presenceList.Bundle = fetchedBundle.Bundle
						presenceList.Category = fetchedBundle.Category

						lemmaLines = append(lemmaLines, presenceList)

						bundle.AppAdsTxtURL = url
						bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
					}

					if isCrawled > 0 {
						crawledBundles = append(crawledBundles, bundle)
					} else {
						failedBundles = append(failedBundles, bundle)
					}

				} else {
					failedBundles = append(failedBundles, bundle)
				}
			}
		case constant.BUNDLE_WEB:
			fetchedBundle, err := parsers.ProcessWebBundle(db, bundle.Bundle)
			if err != nil {
				logger.Info(bundle.Bundle, constant.BUNDLE_CTV, err.Error())
				failedBundles = append(failedBundles, bundle)
			} else {
				if fetchedBundle.Domain == "" {
					adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
					if err != nil {
						log.Printf("Error crawling domain %s: %v", fetchedBundle.Domain, err)
						failedBundles = append(failedBundles, bundle)
					} else {

						fmt.Printf("Crawled domain: %s\n", fetchedBundle.Domain)

						presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

						presenceList.Bundle = fetchedBundle.Bundle
						presenceList.Category = fetchedBundle.Category

						lemmaLines = append(lemmaLines, presenceList)

						bundle.AdsTxtURL = url
						bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)

						crawledBundles = append(crawledBundles, bundle)
					}

				} else {
					failedBundles = append(failedBundles, bundle)
				}
			}
		default:
			logger.Info(bundle.Bundle, constant.BUNDLE_CTV, "Invalid bundle category")
			failedBundles = append(failedBundles, bundle)
		}
	}

	err := repository.SaveCrawledBundlesInDB(db, crawledBundles)
	if err != nil {
		log.Printf("Error saving crawled bundles in DB: %v", err)
	}

	err = repository.SaveFailedBundlesInDB(db, failedBundles)
	if err != nil {
		log.Printf("Error saving failed bundles in DB: %v", err)
	}

	err = repository.SaveLemmaEntriesInDB(db, lemmaLines)
	if err != nil {
		log.Printf("Error saving lemma lines in DB: %v", err)
	}

	fmt.Printf("Total bundles: %d\n", len(tempBundles))
	fmt.Printf("Total crawled bundles: %d\n", len(crawledBundles))
	fmt.Printf("Total failed bundles: %d\n", len(failedBundles))
	fmt.Printf("Total lemma lines: %d\n", len(lemmaLines))
}

func populateBundles() {
	for _, bundle := range models.AndroidBundles {
		bundleInfo := models.BundleInfo{
			Bundle:   bundle,
			Category: constant.BUNDLE_MOBILE_ANDROID,
		}
		tempBundles = append(tempBundles, bundleInfo)
	}
	for _, bundle := range models.IOSBundles {
		bundleInfo := models.BundleInfo{
			Bundle:   bundle,
			Category: constant.BUNDLE_MOBILE_IOS,
		}
		tempBundles = append(tempBundles, bundleInfo)
	}
	for _, bundle := range models.CTVBundles {
		bundleInfo := models.BundleInfo{
			Bundle:   bundle,
			Category: constant.BUNDLE_CTV,
		}
		tempBundles = append(tempBundles, bundleInfo)
	}
	for _, bundle := range models.WebBundles {
		bundleInfo := models.BundleInfo{
			Bundle:   bundle,
			Category: constant.BUNDLE_WEB,
		}
		tempBundles = append(tempBundles, bundleInfo)
	}
}

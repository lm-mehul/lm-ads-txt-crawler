package handler

import (
	"database/sql"
	"log"

	"github.com/lemmamedia/ads-txt-crawler/constant"
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

		// (ios, android,CTV) bundle

		switch bundle.Category {
		case constant.BUNDLE_MOBILE_ANDROID:
			fetchedBundle := parsers.ProcessAndroidBundle(db, bundle.Bundle)
			if fetchedBundle.Domain == "" {
				adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
				if err != nil {
					log.Printf("Error crawling domain %s: %v", fetchedBundle.Domain, err)
				} else {
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
					presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

					presenceList.Bundle = fetchedBundle.Bundle
					presenceList.Category = fetchedBundle.Category

					lemmaLines = append(lemmaLines, presenceList)

					bundle.AppAdsTxtURL = url
					bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
				}

			} else {
				failedBundles = append(failedBundles, bundle)
			}

		case constant.BUNDLE_MOBILE_IOS:
			fetchedBundle := parsers.ProcessIOSBundle(db, bundle.Bundle)
			if fetchedBundle.Domain == "" {
				adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
				if err != nil {
					log.Printf("Error crawling domain %s: %v", fetchedBundle.Domain, err)
				} else {
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
					presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

					presenceList.Bundle = fetchedBundle.Bundle
					presenceList.Category = fetchedBundle.Category

					lemmaLines = append(lemmaLines, presenceList)

					bundle.AppAdsTxtURL = url
					bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
				}

			} else {
				failedBundles = append(failedBundles, bundle)
			}
		case constant.BUNDLE_CTV:
			fetchedBundle := parsers.ProcessCTVBundle(db, bundle.Bundle)
			if fetchedBundle.Domain == "" {
				adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
				if err != nil {
					log.Printf("Error crawling domain %s: %v", fetchedBundle.Domain, err)
				} else {
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
					presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

					presenceList.Bundle = fetchedBundle.Bundle
					presenceList.Category = fetchedBundle.Category

					lemmaLines = append(lemmaLines, presenceList)

					bundle.AppAdsTxtURL = url
					bundle.AppAdsTxtHash = utils.GenerateHash(appAdsTxtPage)
				}

			} else {
				failedBundles = append(failedBundles, bundle)
			}
		case constant.BUNDLE_WEB:
			fetchedBundle := parsers.ProcessWebBundle(db, bundle.Bundle)
			if fetchedBundle.Domain == "" {
				adsTxtPage, url, err := service.CrawlDomain(fetchedBundle.Domain, constant.ADS_TXT_pageType)
				if err != nil {
					log.Printf("Error crawling domain %s: %v", fetchedBundle.Domain, err)
				} else {
					presenceList := service.LemmaDirectsAndResellerInventory(string(adsTxtPage))

					presenceList.Bundle = fetchedBundle.Bundle
					presenceList.Category = fetchedBundle.Category

					lemmaLines = append(lemmaLines, presenceList)

					bundle.AdsTxtURL = url
					bundle.AdsTxtHash = utils.GenerateHash(adsTxtPage)
				}

			} else {
				failedBundles = append(failedBundles, bundle)
			}
		default:
			// log.Printf("Invalid bundle category: %v", bundle.Category)
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

package parsers

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/utils"
	"golang.org/x/net/html"
)

func AndroidBundleParser(db *sql.DB) {
	// androidBundles, err := repository.GetBundlesFromDB(db, constant.BUNDLE_MOBILE_ANDROID)
	// if err != nil {
	// 	log.Printf("Error fetching : %v bundles from database with error : %v", constant.BUNDLE_MOBILE_ANDROID, err)
	// 	return
	// }

	fmt.Println("Executing android bundle parser...")
	processBatch(db, models.AndroidBundles)
}

func processBatch(db *sql.DB, batch []string) {
	var bundles []models.BundleInfo
	var bundle models.BundleInfo
	for _, androidBundle := range batch {
		playStoreURL := fmt.Sprintf("https://play.google.com/store/apps/details?id=%s&hl=en", androidBundle)
		response, err := http.Get(playStoreURL)
		if err != nil {
			utils.LogBundleError(androidBundle, constant.BUNDLE_MOBILE_ANDROID, "Invalid Google Bundle")
			continue
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				log.Printf("Error reading response body: %v", err)
				continue
			}

			bundle.Website = strings.TrimSpace(findWebsiteInHTML(body, "a", "Si6A0c RrSxVb"))
			if bundle.Website == "" {
				utils.LogBundleError(androidBundle, constant.BUNDLE_MOBILE_ANDROID, "Website not found in parser html response.")
				continue
			}
			bundle.Bundle = androidBundle
			bundle.Category = constant.BUNDLE_MOBILE_ANDROID
			bundle.Domain = extractDomainFromBundleURL(bundle.Website)
			fmt.Printf("Android - Bundle: %s, Website: %s, Domain: %s\n", bundle.Bundle, bundle.Website, bundle.Domain)

			bundles = append(bundles, bundle)
		} else {
			utils.LogBundleError(androidBundle, constant.BUNDLE_MOBILE_ANDROID, "Bundle not in Google Playstore")
			continue
		}
	}

	// Save bundles in the database
	// err := repository.SaveCrawledBundlesInDB(db, bundles)
	// if err != nil {
	// 	log.Printf("Error inserting bundles into database: %v", err)
	// }

	// Save uncrawled domains in the database
	// err = repository.SaveUnCrawledDomainsInDB(db, bundles)
	// if err != nil {
	// 	log.Printf("Error saving uncrawled domains into database: %v", err)
	// }
}

func findWebsiteInHTML(body []byte, tagName, classVal string) string {
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Println("Error parsing HTML")
		return ""
	}
	var b func(*html.Node) string
	b = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == tagName {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == classVal {
					// Assuming the href attribute is in the same <a> tag
					for _, attr := range n.Attr {
						if attr.Key == "href" {
							return attr.Val
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if href := b(c); href != "" {
				return href
			}
		}
		return ""
	}
	website := b(doc)
	if website == "" {
		return ""
	}
	return website
}

func extractDomainFromBundleURL(urlStr string) string {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error processing URL '%s': %v\n", urlStr, r)
		}
	}()

	if strings.Contains(urlStr, "/") {
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			panic(err)
		}
		return parsedURL.Hostname()
	} else {
		return strings.TrimSpace(urlStr)
	}
}

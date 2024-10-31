package parsers

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/models"
	"golang.org/x/net/html"
)

// Custom client with redirect count check
var client = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 3 { // Limit redirects to 3
			return fmt.Errorf("stopped after 3 redirects")
		}
		return nil // Allow redirect
	},
}

func ProcessAndroidBundle(db *sql.DB, androidBundle string) (models.BundleInfo, error) {
	var bundle models.BundleInfo

	// Sanitize the androidBundle using url.QueryEscape to handle spaces and special characters
	androidBundle = url.QueryEscape(androidBundle)
	playStoreURL := fmt.Sprintf("https://play.google.com/store/apps/details?id=%s&hl=en", androidBundle)

	response, err := client.Get(playStoreURL)
	if err != nil {
		return bundle, errors.New("Invalid Google Bundle")
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return bundle, errors.New("Error reading response body")
		}

		bundle.Website = strings.TrimSpace(findWebsiteInHTML(body, "a", "Si6A0c RrSxVb"))
		if bundle.Website == "" {
			return bundle, errors.New("Website not found in parser html response")
		}

		bundle.Bundle = androidBundle
		bundle.Category = constant.BUNDLE_MOBILE_ANDROID
		bundle.Domain = extractDomainFromBundleURL(bundle.Website)

		return bundle, nil
	} else {
		return bundle, errors.New("Invalid Google Bundle")
	}

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

	// Remove unwanted escape sequences like %20 (space) and %09 (tab)
	urlStr = strings.ReplaceAll(urlStr, "%20", "")
	urlStr = strings.ReplaceAll(urlStr, "%09", "")

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error processing URL '%s': %v\n", urlStr, r)
		}
	}()

	// if strings.Contains(urlStr, "/") {
	// 	parsedURL, err := url.Parse(urlStr)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	return parsedURL.Hostname()
	// } else {
	// 	return strings.TrimSpace(urlStr)
	// }

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("Invalid URL format after cleanup: %s\n", urlStr)
		return ""
	}

	return parsedURL.Hostname()

}

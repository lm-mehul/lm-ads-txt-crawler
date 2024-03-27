package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/utils"
)

type AlgoliaRequest struct {
	Requests []RequestItem `json:"requests"`
}

type RequestItem struct {
	IndexName string `json:"indexName"`
	Params    string `json:"params"`
}

func ctvBundleParser(db *sql.DB) {
	ctvBundles, err := models.GetBundlesFromDB(db, constant.BUNDLE_CTV)
	if err != nil {
		log.Printf("Error fetching : %v bundles from database with error : %v", constant.BUNDLE_CTV, err)
		return
	}
	fmt.Println("Executing CTV bundle parser...")
	var bundles []models.BundleInfo
	var bundle models.BundleInfo
	batchCount := 0

	algoliaURL := `https://awy63wpylf-1.algolianet.com/1/indexes/*/queries?x-algolia-agent=Algolia%20for%20JavaScript%20(4.14.2)%3B%20Browser%20(lite)%3B%20angular%20(12.0.5)%3B%20angular-instantsearch%20(4.3.0)%3B%20instantsearch.js%20(4.44.0)%3B%20JS%20Helper%20(3.11.0)&x-algolia-api-key=471f4e22aa833a11ef213cd30c540344&x-algolia-application-id=AWY63WPYLF`
	commonParams := `attributesToSnippet=%5B%22description%3A10%22%5D&facets=%5B%22developerCountryName%22%2C%22appStorePrimaryCategories%22%2C%22appStoreSecondaryCategories%22%2C%22iabPrimaryCategory%22%2C%22iabSubCategory%22%2C%22releaseDate%22%2C%22lastUpdatedDate%22%2C%22delistedDate%22%2C%22coppaTargetChildren%22%2C%22hasPrivacyLink%22%2C%22hasAppAdsTxt%22%2C%22hasAds%22%5D&filters=&highlightPostTag=__%2Fais-highlight__&highlightPreTag=__ais-highlight__&hitsPerPage=9&maxValuesPerFacet=10&page=0`

	for _, ctvBundle := range ctvBundles {

		payload, err := constructPayload(commonParams, ctvBundle)
		if err != nil {
			log.Printf("Error creating payload request with error : %v\n", err)
			continue
		}
		// Headers for the request
		headers := map[string]string{"Content-Type": "text/plain"}

		response, err := postData(algoliaURL, payload, headers)
		if err != nil {
			log.Printf("Error making request: %v\n", err)
			continue
		}

		if response.StatusCode == http.StatusOK {
			var jsonResponse map[string]interface{}
			if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
				log.Printf("Error decoding JSON response: %v\n", err)
				continue
			}

			hits := jsonResponse["results"].([]interface{})[0].(map[string]interface{})["hits"].([]interface{})
			if len(hits) > 0 {
				hit := hits[0].(map[string]interface{})
				// appID := hit["appId"].(string)

				if nil == hit["publisherWebsite"] {
					utils.LogBundleError(ctvBundle, constant.BUNDLE_MOBILE_ANDROID, "Website Not Found")
					continue
				}
				bundle.Website = hit["publisherWebsite"].(string)

				bundle.Bundle = ctvBundle
				bundle.Category = constant.BUNDLE_CTV
				bundle.Domain = extractDomainFromBundleURL(bundle.Website)

				bundles = append(bundles, bundle)
			} else {
				fmt.Println("Website Not Found")
				utils.LogBundleError(ctvBundle, constant.BUNDLE_MOBILE_ANDROID, "Website Not Found")
				continue
			}
		} else {
			fmt.Printf("Error: %d, %s\n", response.StatusCode, response.Status)
			fmt.Println("Website Not Found")
			utils.LogBundleError(ctvBundle, constant.BUNDLE_MOBILE_ANDROID, "Website Not Found")
			continue
		}
		batchCount++

		// If batch size is reached, insert the batch into the database
		if batchCount == constant.BATCH_SIZE {
			err := models.SaveCrawledBundlesInDB(db, bundles)
			if nil != err {
				log.Fatal("Failed to save bundles in DB")
			}

			// Reset batch count and values
			batchCount = 0
			bundles = []models.BundleInfo{}
		}
		batchCount++

		// If batch size is reached, insert the batch into the database
		if batchCount == constant.BATCH_SIZE {
			err := models.SaveCrawledBundlesInDB(db, bundles)
			if nil != err {
				log.Fatal("Failed to save bundles in DB")
			}

			// Reset batch count and values
			batchCount = 0
			bundles = []models.BundleInfo{}
		}
	}
	// Insert the remaining batch
	if batchCount > 0 {
		err = models.SaveCrawledBundlesInDB(db, bundles)
		if err != nil {
			log.Printf("Error inserting %v bundles into database with error : %v", constant.BUNDLE_MOBILE_ANDROID, err)
		}
	}
}

func constructPayload(commonParams, ctvBundle string) (string, error) {
	// Construct payload
	payload := AlgoliaRequest{
		Requests: []RequestItem{
			{
				IndexName: "prod_v2_apps",
				Params:    fmt.Sprintf("%s&query=%s", commonParams, ctvBundle),
			},
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error encoding payload: %v\n", err)
		return "", err
	}

	// Convert payload to string
	return string(payloadJSON), nil
}

func postData(url, payload string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}

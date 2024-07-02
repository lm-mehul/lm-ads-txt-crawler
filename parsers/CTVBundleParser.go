package parsers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

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
		log.Printf("Error fetching %v bundles from the database: %v", constant.BUNDLE_CTV, err)
		return
	}
	fmt.Println("Executing CTV bundle parser...")

	var wg sync.WaitGroup
	batchSize := constant.BATCH_SIZE
	numBatches := (len(ctvBundles) + batchSize - 1) / batchSize // Calculate number of batches

	for i := 0; i < numBatches; i++ {
		startIndex := i * batchSize
		endIndex := (i + 1) * batchSize
		if endIndex > len(ctvBundles) {
			endIndex = len(ctvBundles)
		}
		batch := ctvBundles[startIndex:endIndex]

		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			processCTVBatch(db, batch)
		}(batch)
	}

	wg.Wait()

	// Handle remaining bundles
	if numBatches*batchSize < len(ctvBundles) {
		remaining := ctvBundles[numBatches*batchSize:]
		processCTVBatch(db, remaining)
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

func processCTVBatch(db *sql.DB, batch []string) {
	var bundles []models.BundleInfo
	algoliaURL := `https://awy63wpylf-1.algolianet.com/1/indexes/*/queries?x-algolia-agent=Algolia%20for%20JavaScript%20(4.14.2)%3B%20Browser%20(lite)%3B%20angular%20(12.0.5)%3B%20angular-instantsearch%20(4.3.0)%3B%20instantsearch.js%20(4.44.0)%3B%20JS%20Helper%20(3.11.0)&x-algolia-api-key=471f4e22aa833a11ef213cd30c540344&x-algolia-application-id=AWY63WPYLF`
	commonParams := `attributesToSnippet=%5B%22description%3A10%22%5D&facets=%5B%22developerCountryName%22%2C%22appStorePrimaryCategories%22%2C%22appStoreSecondaryCategories%22%2C%22iabPrimaryCategory%22%2C%22iabSubCategory%22%2C%22releaseDate%22%2C%22lastUpdatedDate%22%2C%22delistedDate%22%2C%22coppaTargetChildren%22%2C%22hasPrivacyLink%22%2C%22hasAppAdsTxt%22%2C%22hasAds%22%5D&filters=&highlightPostTag=__%2Fais-highlight__&highlightPreTag=__ais-highlight__&hitsPerPage=9&maxValuesPerFacet=10&page=0`

	for _, ctvBundle := range batch {
		payload, err := constructPayload(commonParams, ctvBundle)
		if err != nil {
			log.Printf("Error creating payload request for bundle %s: %v\n", ctvBundle, err)
			continue
		}

		headers := map[string]string{"Content-Type": "text/plain"}
		response, err := postData(algoliaURL, payload, headers)
		if err != nil {
			log.Printf("Error making request for bundle %s: %v\n", ctvBundle, err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Printf("Error: %d, %s\n", response.StatusCode, response.Status)
			utils.LogBundleError(ctvBundle, constant.BUNDLE_CTV, "Website Not Found")
			continue
		}

		var jsonResponse map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
			log.Printf("Error decoding JSON response for bundle %s: %v\n", ctvBundle, err)
			continue
		}

		hits := jsonResponse["results"].([]interface{})[0].(map[string]interface{})["hits"].([]interface{})
		if len(hits) > 0 {
			hit := hits[0].(map[string]interface{})
			if publisherWebsite, ok := hit["publisherWebsite"].(string); ok {
				bundle := models.BundleInfo{
					Website:  strings.TrimSpace(publisherWebsite),
					Bundle:   ctvBundle,
					Category: constant.BUNDLE_CTV,
					Domain:   extractDomainFromBundleURL(strings.TrimSpace(publisherWebsite)),
				}
				bundles = append(bundles, bundle)
			} else {
				utils.LogBundleError(ctvBundle, constant.BUNDLE_CTV, "Website Not Found")
				continue
			}
		} else {
			utils.LogBundleError(ctvBundle, constant.BUNDLE_CTV, "Website Not Found")
			continue
		}
	}

	// Save bundles and uncrawled domains in the database
	err := models.SaveCrawledBundlesInDB(db, bundles)
	if err != nil {
		log.Printf("Error saving bundles in database: %v\n", err)
	}
	err = models.SaveUnCrawledDomainsInDB(db, bundles)
	if err != nil {
		log.Fatal("Failed to save bundles in DB")
	}
}

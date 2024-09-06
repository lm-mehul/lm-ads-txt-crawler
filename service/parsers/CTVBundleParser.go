package parsers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lemmamedia/ads-txt-crawler/constant"
	"github.com/lemmamedia/ads-txt-crawler/models"
)

func ProcessCTVBundle(db *sql.DB, ctvBundle string) (models.BundleInfo, error) {

	var bundle models.BundleInfo
	algoliaURL := `https://awy63wpylf-1.algolianet.com/1/indexes/*/queries?x-algolia-agent=Algolia%20for%20JavaScript%20(4.14.2)%3B%20Browser%20(lite)%3B%20angular%20(12.0.5)%3B%20angular-instantsearch%20(4.3.0)%3B%20instantsearch.js%20(4.44.0)%3B%20JS%20Helper%20(3.11.0)&x-algolia-api-key=471f4e22aa833a11ef213cd30c540344&x-algolia-application-id=AWY63WPYLF`
	commonParams := `attributesToSnippet=%5B%22description%3A10%22%5D&facets=%5B%22developerCountryName%22%2C%22appStorePrimaryCategories%22%2C%22appStoreSecondaryCategories%22%2C%22iabPrimaryCategory%22%2C%22iabSubCategory%22%2C%22releaseDate%22%2C%22lastUpdatedDate%22%2C%22delistedDate%22%2C%22coppaTargetChildren%22%2C%22hasPrivacyLink%22%2C%22hasAppAdsTxt%22%2C%22hasAds%22%5D&filters=&highlightPostTag=__%2Fais-highlight__&highlightPreTag=__ais-highlight__&hitsPerPage=9&maxValuesPerFacet=10&page=0`

	payload, err := constructPayload(commonParams, ctvBundle)
	if err != nil {
		return bundle, errors.New("Error creating payload request")
	}

	headers := map[string]string{"Content-Type": "text/plain"}
	response, err := postData(algoliaURL, payload, headers)
	if err != nil {
		return bundle, errors.New("Error making request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return bundle, errors.New("HTTP request failed")
	}

	var jsonResponse map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		return bundle, errors.New("Error decoding JSON response")
	}

	hits := jsonResponse["results"].([]interface{})[0].(map[string]interface{})["hits"].([]interface{})
	if len(hits) > 0 {
		hit := hits[0].(map[string]interface{})
		if publisherWebsite, ok := hit["publisherWebsite"].(string); ok {

			bundle.Website = strings.TrimSpace(publisherWebsite)
			bundle.Bundle = ctvBundle
			bundle.Category = constant.BUNDLE_CTV
			bundle.Domain = extractDomainFromBundleURL(strings.TrimSpace(publisherWebsite))

			return bundle, nil
		} else {
			return bundle, errors.New("Website not found in response")
		}
	} else {
		return bundle, errors.New("Website not found in response")
	}
}

type AlgoliaRequest struct {
	Requests []RequestItem `json:"requests"`
}

type RequestItem struct {
	IndexName string `json:"indexName"`
	Params    string `json:"params"`
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

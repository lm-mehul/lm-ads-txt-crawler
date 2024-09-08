package service

import (
	"strings"

	"github.com/lemmamedia/ads-txt-crawler/models"
	"github.com/lemmamedia/ads-txt-crawler/utils"
)

func LemmaDirectsAndResellerInventory(adsTxtPage string) models.LemmaEntry {

	present := models.LemmaEntry{}

	lemmaDirectPubs := make(map[string]struct{})
	lemmaResellerPubs := make(map[string]struct{})

	pageResponse := strings.ToLower(strings.ReplaceAll(adsTxtPage, " ", ""))
	pageResponse = strings.ReplaceAll(pageResponse, "\xa0", "")

	adstxtPageLineSplits := strings.Split(pageResponse, "\n")

	for _, line := range adstxtPageLineSplits {
		if strings.Contains(line, "lemmatechnologies.com") {

			parts := strings.Split(line, ",")

			if len(parts) > 1 {

				publisherID := strings.TrimSpace(parts[1])

				if strings.Contains(line, "direct") {
					lemmaDirectPubs[publisherID] = struct{}{}

				} else if strings.Contains(line, "reseller") {
					lemmaResellerPubs[publisherID] = struct{}{}
				}
			}
		}
	}

	lemmaDirectPubsStr, lemmaResellerPubsStr := utils.MapToStrings(lemmaDirectPubs), utils.MapToStrings(lemmaResellerPubs)

	present.LemmaDirect = lemmaDirectPubsStr
	present.LemmaReseller = lemmaResellerPubsStr
	return present
}

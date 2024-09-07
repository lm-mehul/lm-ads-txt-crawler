package service

import (
	"fmt"
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
			fmt.Printf("Lemma Line: %s\n", line)
			parts := strings.Split(line, ",")
			fmt.Printf("Parts: %v\n", parts)
			if len(parts) > 1 {
				publisherID := strings.TrimSpace(parts[1])
				if strings.Contains(line, "direct") {
					lemmaDirectPubs[publisherID] = struct{}{}
					fmt.Printf("Direct: %s\n", publisherID)
				} else if strings.Contains(line, "reseller") {
					fmt.Printf("Reseller: %s\n", publisherID)
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

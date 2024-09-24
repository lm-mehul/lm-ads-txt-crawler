package service

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lemmamedia/ads-txt-crawler/logger"
	"github.com/lemmamedia/ads-txt-crawler/repository"
)

func ReadAdsTxtDemandLines(db *sql.DB) []string {
	var adsTxtDemandLines []string
	fmt.Println("Fetching ads.txt demand lines from adstxt_lines.txt file...")

	// dir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println("Error getting current directory:", err)
	// 	return []string{}
	// }

	// adsTxtLinesList := utils.ReadLinesFromFile(dir + "/resources/domains/adstxt_lines.txt")

	adsTxtLinesList, err := repository.GetDemandLinesFromDB(db)
	if err != nil {
		logger.Error("Error fetching demand lines from DB: %v", err)
		return adsTxtDemandLines
	}

	for _, adsTxtLine := range adsTxtLinesList {
		line := strings.TrimSpace(strings.ReplaceAll(strings.ToLower(adsTxtLine), " ", ""))
		adsTxtDemandLines = append(adsTxtDemandLines, line)
	}

	return adsTxtDemandLines
}

func IsAdsTxtLinePresent(adsTxtPage string, adsTxtLinesList []string) []string {

	pageResponse := strings.ToLower(strings.ReplaceAll(adsTxtPage, " ", ""))
	pageResponse = strings.ReplaceAll(pageResponse, "\xa0", "")

	present := make([]string, 0)

	for _, adsTxtLine := range adsTxtLinesList {
		searchLine := strings.ReplaceAll(strings.ToLower(adsTxtLine), " ", "")
		if strings.Contains(pageResponse, searchLine) {
			present = append(present, adsTxtLine)
		}
	}
	return present
}

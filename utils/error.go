package utils

import (
	"github.com/lemmamedia/ads-txt-crawler/logger"
)

// logError logs crawling and parsing errors to the database.
func LogBundleError(bundle, bundleType, errMsg string) {
	logger.InfoLogger.Printf("[ %v ] Bundle : %v \t Error : %v \t", bundleType, bundle, errMsg)
}

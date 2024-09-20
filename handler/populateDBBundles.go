package handler

import (
	"database/sql"

	"github.com/lemmamedia/ads-txt-crawler/repository"
)

func MigrateBundlesFromMasterSheet(db *sql.DB) {
	// Populate bundles from the master sheet

	repository.ClearTableData(db, "bundles")

	repository.SaveWebBundlesFromMasterSheet(db)
	repository.SaveBundlesFromMasterSheet(db)
}

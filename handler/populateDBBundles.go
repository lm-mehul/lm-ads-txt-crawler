package handler

import (
	"database/sql"
	"log"

	"github.com/lemmamedia/ads-txt-crawler/repository"
)

func MigrateBundlesFromMasterSheet(db *sql.DB) {
	// Populate bundles from the master sheet
	repository.SaveBundlesFromMasterSheet(db)
}

func CleanOldBackupTables(db *sql.DB) {
	tables := []string{"bundles", "domains", "publishers"}

	for _, table := range tables {
		if err := repository.DropTableIfExists(db, table); err != nil {
			log.Fatalf("Error dropping table %s: %v", table, err)
		}
	}
}

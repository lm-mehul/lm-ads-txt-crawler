package models

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/lemmamedia/ads-txt-crawler/logger"
)

// Generic BatchSave function to save data in batches
func BatchSave[T any](db *sql.DB, data []T, batchSize int, saveFunc func(*sql.DB, []T) error, dataType string) {
	totalCount := len(data)
	batchCount := int(math.Ceil(float64(totalCount) / float64(batchSize)))

	for i := 0; i < batchCount; i++ {
		start := i * batchSize
		end := start + batchSize

		if end > totalCount {
			end = totalCount
		}

		batch := data[start:end]
		err := saveFunc(db, batch)
		if err != nil {
			logger.Error(fmt.Sprintf("Error saving %s batch %d: %v\n", dataType, i+1, err))
		}
	}
}

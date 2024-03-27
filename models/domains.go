package models

import "database/sql"

func GetDomainsFromDB(db *sql.DB, category string) ([]string, error) {
	var bundles []string
	query := "SELECT domain FROM crawled_bundles WHERE category = ? AND is_deleted = 0"
	rows, err := db.Query(query, category)
	if err != nil {
		return bundles, err
	}
	defer rows.Close()

	for rows.Next() {
		var bundle string
		if err := rows.Scan(&bundle); err != nil {
			return bundles, err
		}
		bundles = append(bundles, bundle)
	}
	if err := rows.Err(); err != nil {
		return bundles, err
	}
	return bundles, nil
}

func IsDomainCrawled(domainName string, hash string, db *sql.DB) (bool, error) {
	query := `
        UPDATE crawled_bundles 
        SET ads_txt_page_hash = ?
        WHERE domain = ?
    `

	result, err := db.Exec(query, hash, domainName)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	// If no rows were affected, it means no update was performed, hence return true
	if rowsAffected == 0 {
		return true, nil
	}

	// If rows were affected, it means an update was performed, hence return false
	return false, nil
}

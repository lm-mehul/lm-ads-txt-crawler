package constant

const (
	TABLE_BUNDLES             = "bundles"
	TABLE_FAILED_BUNDLES      = "failed_bundles"
	TABLE_CRAWLED_BUNDLES     = "crawled_bundles"
	TABLE_LEMMA_ENTRIES       = "lemma_entries"
	TABLE_BUNDLE_DEMAND_LINES = "bundle_demand_lines"

	// Schema for the bundles table
	SCHEMA_BUNDLES = `
		bundle VARCHAR(255) PRIMARY KEY,
		category VARCHAR(255)
	`

	// Schema for the failed_bundles table
	SCHEMA_FAILED_BUNDLES = `
		bundle VARCHAR(255) PRIMARY KEY,
		category VARCHAR(255),
		error_message TEXT
	`

	// Schema for the crawled_bundles table
	SCHEMA_CRAWLED_BUNDLES = `
		bundle VARCHAR(255) PRIMARY KEY,
		category VARCHAR(255),
		created_at TIMESTAMP
	`

	// Schema for the lemma_entries table
	SCHEMA_LEMMA_ENTRIES = `
		id INT AUTO_INCREMENT PRIMARY KEY,
		lemma VARCHAR(255),
		created_at TIMESTAMP
	`

	// Schema for the bundle_demand_lines table
	SCHEMA_BUNDLE_DEMAND_LINES = `
		id INT AUTO_INCREMENT PRIMARY KEY,
		bundle VARCHAR(255),
		domain VARCHAR(255),
		seller_account_id VARCHAR(255),
		account_type VARCHAR(255),
		cert_auth_id VARCHAR(255),
		created_at TIMESTAMP
	`
)

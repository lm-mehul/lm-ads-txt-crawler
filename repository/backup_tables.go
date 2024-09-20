package repository

import (
	"bytes"
	"database/sql"
	"log"
)

// DropTableIfExists drops a table if it exists in the database.
func ClearTableData(db *sql.DB, tableName string) error {
	// Initialize a buffer to build the SQL query
	var buff bytes.Buffer
	buff.WriteString("DELETE FROM ")
	buff.WriteString(tableName)

	// Prepare the query for execution
	query := buff.String()

	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error deleting data from table %s: %v", tableName, err)
		return err
	}

	log.Printf("Table %s data cleared successfully if it existed.", tableName)
	return nil
}

// CreateTable creates a table in the database with the specified schema.
func CreateTable(db *sql.DB, tableName string, schema string) error {
	// Initialize a buffer to build the SQL query
	var buff bytes.Buffer
	buff.WriteString("CREATE TABLE IF NOT EXISTS ")
	buff.WriteString(tableName)
	buff.WriteString(" (")
	buff.WriteString(schema)
	buff.WriteString(")")

	// Prepare the query for execution
	query := buff.String()

	// Execute the query
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error creating table %s: %v", tableName, err)
		return err
	}

	log.Printf("Table %s created successfully if it did not exist.", tableName)
	return nil
}

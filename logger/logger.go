package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	// Create logs directory if it doesn't exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Initialize InfoLogger for bundle logs
	infoFile, err := os.OpenFile("logs/bundle_logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening info log file: %v", err)
	}
	InfoLogger = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize ErrorLogger for general errors
	errorFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening error log file: %v", err)
	}
	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info logs informational messages to bundle_logs.log
func Info(v ...interface{}) {
	InfoLogger.Println(v...)
}

// Error logs error messages to app.log
func Error(v ...interface{}) {
	ErrorLogger.Println(v...)
}

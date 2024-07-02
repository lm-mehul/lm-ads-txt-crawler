package logger

import (
	"log"
	"os"
)

// Global loggers
var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	// Initialize InfoLogger
	infoFile, err := os.OpenFile("logs/bundle_logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening info log file: %v", err)
	}

	// Initialize ErrorLogger
	errorFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening error log file: %v", err)
	}

	InfoLogger = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info logs informational messages
func Info(v ...interface{}) {
	InfoLogger.Println(v...)
}

// Info logs informational messages
func Error(v ...interface{}) {
	ErrorLogger.Println(v...)
}

package logger

import (
	"log"
	"os"
)

var (
	InfoLogger *log.Logger
)

func init() {

	// Initialize InfoLogger
	infoFile, err := os.OpenFile("logs/bundle_logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error opening info log file: %v", err)
	}

	// Initialize ErrorLogger
	errorFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("error opening error log file: %v", err)
	}

	InfoLogger = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	// Set the output of logs to the file
	log.SetOutput(errorFile)
}

// Info logs informational messages
func Info(v ...interface{}) {
	InfoLogger.Println(v...)
}

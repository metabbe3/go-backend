package utils

import (
	"log"
	"os"
	"time"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

// InitLogger initializes loggers
func InitLogger() {
	logFile := "app.log"

	// Open or create log file
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set up loggers
	infoLogger = log.New(file, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(file, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(file, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info logs informational messages
func Info(message string) {
	infoLogger.Println(formatLogMessage("INFO", message))
}

// Warning logs warning messages
func Warning(message string) {
	warningLogger.Println(formatLogMessage("WARNING", message))
}

// Error logs error messages
func Error(message string) {
	errorLogger.Println(formatLogMessage("ERROR", message))
}

// formatLogMessage adds a timestamp to the log message
func formatLogMessage(level, message string) string {
	return time.Now().Format("2006-01-02 15:04:05") + " [" + level + "] " + message
}

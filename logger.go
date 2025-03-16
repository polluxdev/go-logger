package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const IMPLEMENT_ME = "implement me!"

// Interface -.
type Interface interface {
	Debug(message interface{}, fields ...interface{})
	Info(message string, fields ...interface{})
	Warn(message string, fields ...interface{})
	Error(message interface{}, fields ...interface{})
	Fatal(message interface{}, fields ...interface{})
}

// Logger -.
type Logger struct {
	mu       sync.Mutex
	logFile  *lumberjack.Logger
	filePath string
}

// initLogger
func (l *Logger) initLogger(level string) {
	// Set log level
	l.setLogLevel(level)

	panic(IMPLEMENT_ME)
}

// setLogLevel -.
func (l *Logger) setLogLevel(level string) {
	panic(IMPLEMENT_ME)
}

// initRotate -.
func (l *Logger) initRotate() {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Ensure the logs directory exists
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		return
	}

	// Define the log file with the current date
	currentDate := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s/%s.log", logDir, currentDate)

	// Define log rotation settings
	l.logFile = &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10, // MB
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   false,
	}

	// Save the current log file path for compression later
	l.filePath = filename
}

// compressLogFile -.
func (l *Logger) compressLogFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	oldLogFile := l.filePath
	compressedFile := fmt.Sprintf("%s.gz", oldLogFile)

	inputFile, err := os.Open(oldLogFile)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(compressedFile)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := gzip.NewWriter(outputFile)
	defer writer.Close()

	_, err = io.Copy(writer, inputFile)
	if err != nil {
		return err
	}

	// Remove the original uncompressed log file after compression
	return os.Remove(oldLogFile)
}

// autoRotateLogs -.
func (l *Logger) autoRotateLogs(level string) {
	for {
		// Get the current time
		now := time.Now()

		// Calculate time until midnight
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		duration := time.Until(nextMidnight)

		fmt.Printf("Next log rotation at: %s (in %s)\n", nextMidnight.Format(time.DateTime), duration)

		// Wait until midnight
		time.Sleep(duration)

		// Compress the old log file before rotating
		if err := l.compressLogFile(); err != nil {
			fmt.Printf("Failed to compress log file: %v\n", err)
		}

		// Reinitialize the logger with the new day's log file
		l.initRotate()
		l.initLogger(level)
	}
}

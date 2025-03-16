package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// Logrus -.
type Logrus struct {
	goLogger
	logger *logrus.Logger
	level  logrus.Level
}

var _ Interface = (*Logrus)(nil)

// NewLogrus -.
func NewLogrus(level string) *Logrus {
	logger := &Logrus{}
	logger.initRotate()
	logger.initLogger(level)

	// Start a goroutine to monitor date changes and rotate logs at midnight
	go logger.autoRotateLogs(level)

	return logger
}

// initLogger
func (l *Logrus) initLogger(level string) {
	// Set log level
	l.setLogLevel(level)

	// Initialize the logger
	l.logger = logrus.New()

	// Set output to multi-writer: log to file and console
	l.logger.SetOutput(io.MultiWriter(l.logFile, os.Stdout))

	// Set the formatter to JSON
	l.logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339, // Custom timestamp format
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// Customize the caller information
			return "", fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
		},
	})

	// Add caller information
	l.logger.SetReportCaller(true)
}

// setLogLevel -.
func (l *Logrus) setLogLevel(level string) {
	switch level {
	case "debug":
		l.level = logrus.DebugLevel
	case "info":
		l.level = logrus.InfoLevel
	case "warn":
		l.level = logrus.WarnLevel
	case "error":
		l.level = logrus.ErrorLevel
	case "fatal":
		l.level = logrus.FatalLevel
	case "panic":
		l.level = logrus.PanicLevel
	default:
		l.level = logrus.InfoLevel // Default to info level
	}

	logrus.SetLevel(l.level)
}

// Debug -.
func (l *Logrus) Debug(message interface{}, fields ...interface{}) {
	l.msg("debug", message, fields...)
}

// Info -.
func (l *Logrus) Info(message string, fields ...interface{}) {
	l.msg("info", message, fields...)
}

// Warn -.
func (l *Logrus) Warn(message string, fields ...interface{}) {
	l.msg("warn", message, fields...)
}

// Error -.
func (l *Logrus) Error(message interface{}, fields ...interface{}) {
	l.msg("error", message, fields...)
}

// Fatal -.
func (l *Logrus) Fatal(message interface{}, fields ...interface{}) {
	l.msg("fatal", message, fields...)

	os.Exit(1)
}

// fields -.
func (l *Logrus) fields(fields ...interface{}) logrus.Fields {
	// Convert variadic fields into logrus.Fields
	result := make(logrus.Fields)
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				result[key] = fields[i+1]
			}
		}
	}

	return result
}

// log -.
func (l *Logrus) log(message string, fields ...interface{}) {
	l.logger.WithFields(l.fields(fields...)).Info(message)
}

// msg -.
func (l *Logrus) msg(level string, message interface{}, fields ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), fields...)
	case string:
		l.log(msg, fields...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), fields...)
	}
}

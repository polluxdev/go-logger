package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Zerolog -.
type Zerolog struct {
	goLogger
	logger zerolog.Logger
	level  zerolog.Level
}

var _ Interface = (*Zerolog)(nil)

// NewZerolog -.
func NewZerolog(level string) *Zerolog {
	logger := &Zerolog{}
	logger.initRotate()
	logger.initLogger(level)

	// Start a goroutine to monitor date changes and rotate logs at midnight
	go logger.autoRotateLogs(level)

	return logger
}

// initLogger
func (l *Zerolog) initLogger(level string) {
	// Set log level
	l.setLogLevel(level)

	// Multi-writer: log to file and console
	multiWriter := zerolog.MultiLevelWriter(l.logFile, os.Stdout)

	// Initialize the logger
	skipFrameCount := 3
	l.logger = zerolog.New(multiWriter).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()
}

// setLogLevel -.
func (l *Zerolog) setLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		l.level = zerolog.DebugLevel
	case "info":
		l.level = zerolog.InfoLevel
	case "warn":
		l.level = zerolog.WarnLevel
	case "error":
		l.level = zerolog.ErrorLevel
	case "fatal":
		l.level = zerolog.FatalLevel
	default:
		l.level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l.level)
}

// Debug -.
func (l *Zerolog) Debug(message interface{}, fields ...interface{}) {
	l.msg("debug", message, fields...)
}

// Info -.
func (l *Zerolog) Info(message string, fields ...interface{}) {
	l.msg("info", message, fields...)
}

// Warn -.
func (l *Zerolog) Warn(message string, fields ...interface{}) {
	l.msg("warn", message, fields...)
}

// Error -.
func (l *Zerolog) Error(message interface{}, fields ...interface{}) {
	l.msg("error", message, fields...)
}

// Fatal -.
func (l *Zerolog) Fatal(message interface{}, fields ...interface{}) {
	l.msg("fatal", message, fields...)

	os.Exit(1)
}

// log -.
func (l *Zerolog) log(message string, fields ...interface{}) {
	l.logger.Info().Fields(fields).Msg(message)
}

// msg -.
func (l *Zerolog) msg(level string, message interface{}, fields ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), fields...)
	case string:
		l.log(msg, fields...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), fields...)
	}
}

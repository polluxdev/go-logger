package gologger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zap -.
type Zap struct {
	goLogger
	logger *zap.Logger
	level  zapcore.Level
}

var _ Interface = (*Zap)(nil)

// NewZap -.
func NewZap(level string) *Zap {
	logger := &Zap{}
	logger.initRotate()
	logger.initLogger(level)

	// Start a goroutine to monitor date changes and rotate logs at midnight
	go logger.autoRotateLogs(level)

	return logger
}

// initLogger
func (l *Zap) initLogger(level string) {
	// Set log level
	l.setLogLevel(level)

	// Configure the file writer
	fileWriter := zapcore.AddSync(l.logFile)

	// Configure the console writer
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Create a JSON encoder configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	// Create a multi-writer (file and console)
	multiWriter := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			fileWriter,
			l.level,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			consoleWriter,
			l.level,
		),
	)

	// Add caller information and stacktrace
	skipFrameCount := 3
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(skipFrameCount),
		zap.AddStacktrace(zapcore.ErrorLevel), // Add stacktrace for errors and above
	}

	// Initialize the logger
	l.logger = zap.New(multiWriter, options...)
}

// setLogLevel -.
func (l *Zap) setLogLevel(level string) {
	switch level {
	case "debug":
		l.level = zapcore.DebugLevel
	case "info":
		l.level = zapcore.InfoLevel
	case "warn":
		l.level = zapcore.WarnLevel
	case "error":
		l.level = zapcore.ErrorLevel
	case "fatal":
		l.level = zapcore.FatalLevel
	case "panic":
		l.level = zapcore.PanicLevel
	default:
		l.level = zapcore.InfoLevel // Default to info level
	}
}

// Debug -.
func (l *Zap) Debug(message interface{}, fields ...interface{}) {
	l.msg("debug", message, fields...)
}

// Info -.
func (l *Zap) Info(message string, fields ...interface{}) {
	l.msg("info", message, fields...)
}

// Warn -.
func (l *Zap) Warn(message string, fields ...interface{}) {
	l.msg("warn", message, fields...)
}

// Error -.
func (l *Zap) Error(message interface{}, fields ...interface{}) {
	l.msg("error", message, fields...)
}

// Fatal -.
func (l *Zap) Fatal(message interface{}, fields ...interface{}) {
	l.msg("fatal", message, fields...)

	os.Exit(1)
}

// fields -.
func (l *Zap) fields(fields ...interface{}) []zap.Field {
	// Convert variadic fields into zap.Field slice
	var result []zap.Field
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				result = append(result, zap.Any(key, fields[i+1]))
			}
		}
	}

	return result
}

// log -.
func (l *Zap) log(message string, fields ...interface{}) {
	l.logger.Info(message, l.fields(fields...)...)
}

// msg -.
func (l *Zap) msg(level string, message interface{}, fields ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), fields...)
	case string:
		l.log(msg, fields...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), fields...)
	}
}

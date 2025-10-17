package togomq

import (
	"fmt"
	"log"
	"strings"
)

// LogLevel represents the logging level
type LogLevel int

const (
	// LogLevelDebug enables all logs including debug messages
	LogLevelDebug LogLevel = iota
	// LogLevelInfo enables info, warning, and error messages
	LogLevelInfo
	// LogLevelWarn enables warning and error messages
	LogLevelWarn
	// LogLevelError enables only error messages
	LogLevelError
	// LogLevelNone disables all logging
	LogLevelNone
)

// ParseLogLevel converts a string to LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "none":
		return LogLevelNone
	default:
		return LogLevelInfo
	}
}

// Logger provides logging functionality for the SDK
type Logger struct {
	level LogLevel
}

// NewLogger creates a new logger with the specified level
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level: level,
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= LogLevelInfo {
		log.Printf("[INFO] "+format, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= LogLevelWarn {
		log.Printf("[WARN] "+format, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= LogLevelError {
		log.Printf("[ERROR] "+format, args...)
	}
}

// Errorf logs an error message and returns a formatted error
func (l *Logger) Errorf(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	l.Error(err.Error())
	return err
}

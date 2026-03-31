package goop

import (
	"log"
	"time"
)

// Debug levels for different verbosity levels
type DebugLevel int

const (
	DebugOff     DebugLevel = 0
	DebugBasic   DebugLevel = 1
	DebugVerbose DebugLevel = 2
	DebugTrace   DebugLevel = 3
)

// Logger interface for custom debug logging
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Tracef(format string, args ...interface{})
}

// Global debug configuration
var (
	debugLevel  DebugLevel = DebugOff
	debugLogger Logger     = &defaultLogger{}
)

// Initialize cache with default config
func init() {
	SetCacheConfig(DefaultCacheConfig)
}

// defaultLogger implements Logger using standard log package
type defaultLogger struct{}

func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *defaultLogger) Tracef(format string, args ...interface{}) {
	log.Printf("[TRACE] "+format, args...)
}

// SetDebugLevel sets debug verbosity level
func SetDebugLevel(level DebugLevel) {
	debugLevel = level
}

// SetDebugLogger sets a custom logger for debug output
func SetDebugLogger(logger Logger) {
	debugLogger = logger
}

// GetDebugLevel returns current debug level
func GetDebugLevel() DebugLevel {
	return debugLevel
}

// debugLog logs messages based on current debug level
func debugLog(level DebugLevel, format string, args ...interface{}) {
	if debugLevel >= level {
		switch level {
		case DebugBasic:
			debugLogger.Debugf(format, args...)
		case DebugVerbose:
			debugLogger.Infof(format, args...)
		case DebugTrace:
			debugLogger.Tracef(format, args...)
		}
	}
}

// timing helper for performance measurement
type operationTimer struct {
	operation string
	start     time.Time
	level     DebugLevel
}

// startTimer begins timing an operation
func startTimer(operation string, level DebugLevel) *operationTimer {
	debugLog(level, "Starting operation: %s", operation)
	return &operationTimer{
		operation: operation,
		start:     time.Now(),
		level:     level,
	}
}

// finish ends timing and logs duration
func (t *operationTimer) finish() {
	duration := time.Since(t.start)
	debugLog(t.level, "Completed operation: %s (took %v)", t.operation, duration)
}

// logHTTPRequest logs HTTP request details
func logHTTPRequest(method, url string, headers map[string]string) {
	if debugLevel >= DebugTrace {
		debugLog(DebugTrace, "HTTP Request: %s %s", method, url)
		for k, v := range headers {
			debugLog(DebugTrace, "  Header: %s: %s", k, v)
		}
	}
}

// logHTTPResponse logs HTTP response details
func logHTTPResponse(statusCode int, contentLength int) {
	if debugLevel >= DebugTrace {
		debugLog(DebugTrace, "HTTP Response: %d (%d bytes)", statusCode, contentLength)
	}
}

// logDOMOperation logs DOM traversal operations
func logDOMOperation(operation, selector string, foundCount int) {
	if debugLevel >= DebugVerbose {
		debugLog(DebugVerbose, "DOM Operation: %s '%s' -> %d elements found", operation, selector, foundCount)
	}
}

package goop

import (
	"time"
)

// Default timeout for HTTP requests
var DefaultTimeout = 30 * time.Second

// Timeout configuration for HTTP requests
type TimeoutConfig struct {
	Timeout time.Duration
}

// SetTimeout sets the global default timeout for all HTTP requests
func SetTimeout(timeout time.Duration) {
	DefaultTimeout = timeout
}

// GetTimeout returns the current global default timeout
func GetTimeout() time.Duration {
	return DefaultTimeout
}

// NewTimeoutConfig creates a new timeout configuration
func NewTimeoutConfig(timeout time.Duration) TimeoutConfig {
	return TimeoutConfig{Timeout: timeout}
}

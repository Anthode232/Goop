package goop

import (
	"fmt"
	"math"
	"time"
)

// RetryConfig defines retry behavior
type RetryConfig struct {
	MaxAttempts int              // Maximum number of retry attempts
	BaseDelay   time.Duration    // Base delay between retries
	MaxDelay    time.Duration    // Maximum delay between retries
	Multiplier  float64          // Backoff multiplier
	Retryable   func(error) bool // Function to determine if error is retryable
}

// DefaultRetryConfig provides sensible defaults
var DefaultRetryConfig = RetryConfig{
	MaxAttempts: 3,
	BaseDelay:   1 * time.Second,
	MaxDelay:    30 * time.Second,
	Multiplier:  2.0,
	Retryable:   isRetryableError,
}

// Global retry configuration
var globalRetryConfig = DefaultRetryConfig

// SetRetryConfig sets the global retry configuration
func SetRetryConfig(config RetryConfig) {
	globalRetryConfig = config
}

// GetRetryConfig returns the current global retry configuration
func GetRetryConfig() RetryConfig {
	return globalRetryConfig
}

// NewRetryConfig creates a new retry configuration
func NewRetryConfig(maxAttempts int, baseDelay, maxDelay time.Duration) RetryConfig {
	return RetryConfig{
		MaxAttempts: maxAttempts,
		BaseDelay:   baseDelay,
		MaxDelay:    maxDelay,
		Multiplier:  2.0,
		Retryable:   isRetryableError,
	}
}

// isRetryableError determines if an error should be retried
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's a Goop error
	if goopErr, ok := err.(Error); ok {
		switch goopErr.Type {
		case ErrInGetRequest, ErrTimeout:
			return true
		case ErrCreatingGetRequest, ErrUnableToParse, ErrElementNotFound:
			return false
		default:
			return true // Retry other errors by default
		}
	}

	// For non-Goop errors, assume they're retryable (network errors, etc.)
	return true
}

// calculateDelay calculates the delay for a given attempt using exponential backoff
func calculateDelay(attempt int, config RetryConfig) time.Duration {
	if attempt <= 0 {
		return 0
	}

	// Exponential backoff: delay = baseDelay * (multiplier ^ (attempt - 1))
	delay := float64(config.BaseDelay) * math.Pow(config.Multiplier, float64(attempt-1))

	// Apply jitter (±25% random variation)
	jitter := delay * 0.25 * (2*float64(time.Now().UnixNano()%1000)/1000 - 1)
	delay += jitter

	// Ensure delay doesn't exceed max delay
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	return time.Duration(delay)
}

// RetryGet performs a GET request with retry logic
func RetryGet(url string, timeout time.Duration) (string, error) {
	return RetryGetWithConfig(url, timeout, globalRetryConfig)
}

// RetryGetWithConfig performs a GET request with custom retry configuration
func RetryGetWithConfig(url string, timeout time.Duration, config RetryConfig) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		debugLog(DebugVerbose, "HTTP GET attempt %d/%d for %s", attempt, config.MaxAttempts, url)

		// Apply rate limiting
		waitForRateLimit()

		resp, err := GetWithTimeout(url, timeout)
		if err == nil {
			debugLog(DebugVerbose, "HTTP GET succeeded on attempt %d for %s", attempt, url)
			return resp, nil
		}

		lastErr = err

		// Check if error is retryable
		if !config.Retryable(err) {
			debugLog(DebugBasic, "HTTP GET failed with non-retryable error: %v", err)
			return "", err
		}

		// If this is the last attempt, don't wait
		if attempt == config.MaxAttempts {
			break
		}

		// Calculate delay and wait
		delay := calculateDelay(attempt, config)
		debugLog(DebugVerbose, "HTTP GET failed (attempt %d), retrying in %v: %v", attempt, delay, err)
		time.Sleep(delay)
	}

	debugLog(DebugBasic, "HTTP GET failed after %d attempts: %v", config.MaxAttempts, lastErr)
	return "", fmt.Errorf("failed after %d attempts: %w", config.MaxAttempts, lastErr)
}

// RetryPost performs a POST request with retry logic
func RetryPost(url string, bodyType string, body interface{}, timeout time.Duration) (string, error) {
	return RetryPostWithConfig(url, bodyType, body, timeout, globalRetryConfig)
}

// RetryPostWithConfig performs a POST request with custom retry configuration
func RetryPostWithConfig(url string, bodyType string, body interface{}, timeout time.Duration, config RetryConfig) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		debugLog(DebugVerbose, "HTTP POST attempt %d/%d for %s", attempt, config.MaxAttempts, url)

		// Apply rate limiting
		waitForRateLimit()

		resp, err := PostWithTimeout(url, bodyType, body, timeout)
		if err == nil {
			debugLog(DebugVerbose, "HTTP POST succeeded on attempt %d for %s", attempt, url)
			return resp, nil
		}

		lastErr = err

		// Check if error is retryable
		if !config.Retryable(err) {
			debugLog(DebugBasic, "HTTP POST failed with non-retryable error: %v", err)
			return "", err
		}

		// If this is the last attempt, don't wait
		if attempt == config.MaxAttempts {
			break
		}

		// Calculate delay and wait
		delay := calculateDelay(attempt, config)
		debugLog(DebugVerbose, "HTTP POST failed (attempt %d), retrying in %v: %v", attempt, delay, err)
		time.Sleep(delay)
	}

	debugLog(DebugBasic, "HTTP POST failed after %d attempts: %v", config.MaxAttempts, lastErr)
	return "", fmt.Errorf("failed after %d attempts: %w", config.MaxAttempts, lastErr)
}

// RetryOperation executes any function with retry logic
func RetryOperation(operation func() error, config RetryConfig) error {
	var lastErr error

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		debugLog(DebugVerbose, "Operation attempt %d/%d", attempt, config.MaxAttempts)

		err := operation()
		if err == nil {
			debugLog(DebugVerbose, "Operation succeeded on attempt %d", attempt)
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !config.Retryable(err) {
			debugLog(DebugBasic, "Operation failed with non-retryable error: %v", err)
			return err
		}

		// If this is the last attempt, don't wait
		if attempt == config.MaxAttempts {
			break
		}

		// Calculate delay and wait
		delay := calculateDelay(attempt, config)
		debugLog(DebugVerbose, "Operation failed (attempt %d), retrying in %v: %v", attempt, delay, err)
		time.Sleep(delay)
	}

	debugLog(DebugBasic, "Operation failed after %d attempts: %v", config.MaxAttempts, lastErr)
	return fmt.Errorf("failed after %d attempts: %w", config.MaxAttempts, lastErr)
}

// RetryWithBackoff is a convenience function for simple retry scenarios
func RetryWithBackoff(maxAttempts int, baseDelay time.Duration, operation func() error) error {
	config := RetryConfig{
		MaxAttempts: maxAttempts,
		BaseDelay:   baseDelay,
		MaxDelay:    baseDelay * time.Duration(math.Pow(2, float64(maxAttempts-1))),
		Multiplier:  2.0,
		Retryable:   isRetryableError,
	}

	return RetryOperation(operation, config)
}

package goop

import (
	"sync"
	"time"
)

// RateLimiter controls the rate of requests to avoid overwhelming servers
type RateLimiter struct {
	interval    time.Duration
	lastRequest time.Time
	mutex       sync.Mutex
}

// Global rate limiter instance
var globalRateLimiter *RateLimiter

// SetRateLimit sets a global rate limit for all HTTP requests
// interval specifies the minimum time between requests
func SetRateLimit(interval time.Duration) {
	globalRateLimiter = &RateLimiter{
		interval: interval,
	}
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{
		interval: interval,
	}
}

// Wait blocks until the rate limit allows the next request
func (rl *RateLimiter) Wait() {
	if rl == nil || rl.interval <= 0 {
		return
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	if rl.lastRequest.IsZero() {
		rl.lastRequest = now
		return
	}

	elapsed := now.Sub(rl.lastRequest)
	if elapsed < rl.interval {
		sleepTime := rl.interval - elapsed
		debugLog(DebugVerbose, "Rate limiting: waiting %v before next request", sleepTime)
		time.Sleep(sleepTime)
	}

	rl.lastRequest = time.Now()
}

// GetRateLimit returns the current global rate limit interval
func GetRateLimit() time.Duration {
	if globalRateLimiter == nil {
		return 0
	}
	return globalRateLimiter.interval
}

// ResetRateLimit clears the global rate limit
func ResetRateLimit() {
	globalRateLimiter = nil
}

// waitForRateLimit applies global rate limiting if configured
func waitForRateLimit() {
	if globalRateLimiter != nil {
		globalRateLimiter.Wait()
	}
}

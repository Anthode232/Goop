package goop

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// ConcurrentOptions defines configuration for concurrent scraping
type ConcurrentOptions struct {
	Workers          int                      // Number of concurrent workers
	RateLimit        time.Duration            // Rate limit between requests
	Timeout          time.Duration            // Timeout per request
	RetryAttempts    int                      // Number of retry attempts
	ProgressCallback func(current, total int) // Progress callback
	BatchSize        int                      // Batch size for processing
	EnableStreaming  bool                     // Enable streaming output
	MaxMemoryUsage   int64                    // Maximum memory usage in bytes
	CompressionLevel int                      // Output compression level
}

// FastConcurrentOptions provides optimized defaults for speed
var FastConcurrentOptions = ConcurrentOptions{
	Workers:          20,                     // Higher worker count for speed
	RateLimit:        100 * time.Millisecond, // Faster rate limiting
	Timeout:          15 * time.Second,       // Shorter timeout for faster failure
	RetryAttempts:    2,                      // Fewer retries for speed
	BatchSize:        100,                    // Process in batches
	EnableStreaming:  true,                   // Enable streaming by default
	MaxMemoryUsage:   100 * 1024 * 1024,      // 100MB memory limit
	CompressionLevel: 6,                      // Balanced compression
}

// DefaultConcurrentOptions provides sensible defaults
var DefaultConcurrentOptions = ConcurrentOptions{
	Workers:          5,
	RateLimit:        1 * time.Second,
	Timeout:          30 * time.Second,
	RetryAttempts:    3,
	ProgressCallback: nil,
	BatchSize:        0,     // No batching by default
	EnableStreaming:  false, // Streaming disabled by default
	MaxMemoryUsage:   0,     // No memory limit by default
	CompressionLevel: 0,     // No compression by default
}

// ScrapeResult represents the result of scraping a single URL
type ScrapeResult struct {
	URL   string // The URL that was scraped
	Data  Root   // The parsed HTML content
	Error error  // Any error that occurred
}

// ScrapeAllFast scrapes multiple URLs concurrently with speed optimizations
func ScrapeAllFast(urls []string, options *ConcurrentOptions) ([]ScrapeResult, error) {
	if options == nil {
		options = &FastConcurrentOptions
	}

	if len(urls) == 0 {
		return []ScrapeResult{}, nil
	}

	timer := startTimer(fmt.Sprintf("Fast concurrent scraping %d URLs", len(urls)), DebugVerbose)
	defer timer.finish()

	// Use batching for large URL lists
	if options.BatchSize > 0 && len(urls) > options.BatchSize {
		return scrapeBatches(urls, options)
	}

	return scrapeConcurrentInternal(urls, options)
}

// scrapeBatches processes URLs in batches for memory efficiency
func scrapeBatches(urls []string, options *ConcurrentOptions) ([]ScrapeResult, error) {
	var allResults []ScrapeResult

	for i := 0; i < len(urls); i += options.BatchSize {
		end := i + options.BatchSize
		if end > len(urls) {
			end = len(urls)
		}

		batch := urls[i:end]
		batchResults, err := scrapeConcurrentInternal(batch, options)
		if err != nil {
			return allResults, err
		}

		allResults = append(allResults, batchResults...)

		// Force garbage collection for memory management
		if options.MaxMemoryUsage > 0 {
			// Simple memory check - could be enhanced with runtime.MemStats
			if i%(options.BatchSize*5) == 0 {
				runtime.GC()
			}
		}
	}

	return allResults, nil
}

// scrapeConcurrentInternal is the optimized core concurrent scraping function
func scrapeConcurrentInternal(urls []string, options *ConcurrentOptions) ([]ScrapeResult, error) {
	results := make([]ScrapeResult, len(urls))
	var resultsMutex sync.Mutex

	// Create optimized rate limiter
	rateLimiter := time.NewTicker(options.RateLimit)
	defer rateLimiter.Stop()

	// Create error group with optimized context
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(options.Workers)

	// Pre-allocate URL channel for better performance
	urlChan := make(chan string, len(urls))
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	// Process URLs concurrently with optimizations
	for i := 0; i < len(urls); i++ {
		i := i // Create loop-local copy

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case url := <-urlChan:
				// Fast rate limiting
				if options.RateLimit > 0 {
					<-rateLimiter.C
				}

				// Progress callback without blocking
				if options.ProgressCallback != nil {
					go options.ProgressCallback(i+1, len(urls))
				}

				// Fast scraping with minimal overhead
				result := fastScrape(url, options.RetryAttempts, options.Timeout)

				// Thread-safe result storage
				resultsMutex.Lock()
				results[i] = result
				resultsMutex.Unlock()

				return nil
			}
		})
	}

	// Wait for completion with timeout
	if err := g.Wait(); err != nil {
		return results, fmt.Errorf("concurrent scraping failed: %v", err)
	}

	debugLog(DebugVerbose, "Completed fast concurrent scraping of %d URLs", len(urls))
	return results, nil
}

// fastScrape is an optimized scraping function with minimal overhead
func fastScrape(url string, retryAttempts int, timeout time.Duration) ScrapeResult {
	// Fast path for single attempt
	if retryAttempts <= 1 {
		result := Scrape(url, timeout)
		return ScrapeResult{
			URL:   url,
			Data:  result,
			Error: result.Error,
		}
	}

	// Retry logic with exponential backoff
	var result Root
	var err error

	for attempt := 0; attempt <= retryAttempts; attempt++ {
		if attempt > 0 {
			// Fast exponential backoff
			backoff := time.Duration(attempt) * 100 * time.Millisecond
			if backoff > 2*time.Second {
				backoff = 2 * time.Second
			}
			time.Sleep(backoff)
		}

		result = Scrape(url, timeout)
		if result.Error == nil {
			break
		}
		err = result.Error
	}

	return ScrapeResult{
		URL:   url,
		Data:  result,
		Error: err,
	}
}

// ScrapeAll scrapes multiple URLs concurrently
func ScrapeAll(urls []string, options *ConcurrentOptions) ([]ScrapeResult, error) {
	if options == nil {
		options = &DefaultConcurrentOptions
	}

	if len(urls) == 0 {
		return []ScrapeResult{}, nil
	}

	timer := startTimer(fmt.Sprintf("Concurrent scraping %d URLs", len(urls)), DebugVerbose)
	defer timer.finish()

	results := make([]ScrapeResult, len(urls))
	var resultsMutex sync.Mutex

	// Create rate limiter
	rateLimiter := time.NewTicker(options.RateLimit)
	defer rateLimiter.Stop()

	// Create error group for concurrency control
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(options.Workers)

	// Process URLs concurrently
	for i, url := range urls {
		i, url := i, url // Create loop-local copies

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-rateLimiter.C:
				// Rate limiting
			}

			// Progress callback
			if options.ProgressCallback != nil {
				options.ProgressCallback(i+1, len(urls))
			}

			// Scrape with retry logic
			var result Root
			var err error

			for attempt := 0; attempt <= options.RetryAttempts; attempt++ {
				if attempt > 0 {
					debugLog(DebugVerbose, "Retry attempt %d for URL: %s", attempt, url)
					time.Sleep(time.Duration(attempt) * time.Second) // Exponential backoff
				}

				result = Scrape(url, options.Timeout)
				if result.Error == nil {
					break
				}
				err = result.Error
			}

			// Store result
			resultsMutex.Lock()
			results[i] = ScrapeResult{
				URL:   url,
				Data:  result,
				Error: err,
			}
			resultsMutex.Unlock()

			return nil
		})
	}

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		return results, fmt.Errorf("concurrent scraping failed: %v", err)
	}

	debugLog(DebugVerbose, "Completed concurrent scraping of %d URLs", len(urls))
	return results, nil
}

// ScrapeAllWithJS scrapes multiple URLs concurrently with JavaScript rendering
func ScrapeAllWithJS(urls []string, jsOptions *JSOptions, concurrentOptions *ConcurrentOptions) ([]ScrapeResult, error) {
	if concurrentOptions == nil {
		concurrentOptions = &DefaultConcurrentOptions
	}

	if len(urls) == 0 {
		return []ScrapeResult{}, nil
	}

	timer := startTimer(fmt.Sprintf("Concurrent JS scraping %d URLs", len(urls)), DebugVerbose)
	defer timer.finish()

	results := make([]ScrapeResult, len(urls))
	var resultsMutex sync.Mutex

	// Create rate limiter
	rateLimiter := time.NewTicker(concurrentOptions.RateLimit)
	defer rateLimiter.Stop()

	// Initialize browser pool for JS rendering
	if jsOptions == nil {
		jsOptions = &DefaultJSOptions
	}

	// Calculate optimal browser pool size
	browserPoolSize := concurrentOptions.Workers
	if browserPoolSize > 10 {
		browserPoolSize = 10 // Limit browser instances
	}

	// Note: Browser pool management is simplified in this version
	// Each JS request creates its own context

	// Create error group for concurrency control
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(concurrentOptions.Workers)

	// Process URLs concurrently
	for i, url := range urls {
		i, url := i, url // Create loop-local copies

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-rateLimiter.C:
				// Rate limiting
			}

			// Progress callback
			if concurrentOptions.ProgressCallback != nil {
				concurrentOptions.ProgressCallback(i+1, len(urls))
			}

			// Scrape with JS and retry logic
			var result Root
			var err error

			for attempt := 0; attempt <= concurrentOptions.RetryAttempts; attempt++ {
				if attempt > 0 {
					debugLog(DebugVerbose, "JS retry attempt %d for URL: %s", attempt, url)
					time.Sleep(time.Duration(attempt) * 2 * time.Second) // Exponential backoff for JS
				}

				result, err = RenderWithJS(url, jsOptions)
				if err == nil {
					break
				}
			}

			// Store result
			resultsMutex.Lock()
			results[i] = ScrapeResult{
				URL:   url,
				Data:  result,
				Error: err,
			}
			resultsMutex.Unlock()

			return nil
		})
	}

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		return results, fmt.Errorf("concurrent JS scraping failed: %v", err)
	}

	debugLog(DebugVerbose, "Completed concurrent JS scraping of %d URLs", len(urls))
	return results, nil
}

// ScrapeConcurrently scrapes multiple selectors concurrently on a single page
func (r Root) ScrapeConcurrently(selectors []string, options *ConcurrentOptions) ([]Root, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	if options == nil {
		options = &DefaultConcurrentOptions
	}

	if len(selectors) == 0 {
		return []Root{}, nil
	}

	results := make([]Root, len(selectors))
	var resultsMutex sync.Mutex

	// Create error group for concurrency control
	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(options.Workers)

	// Process selectors concurrently
	for i, selector := range selectors {
		i, selector := i, selector // Create loop-local copies

		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Process selector
			}

			result := r.Find(selector)

			// Store result
			resultsMutex.Lock()
			results[i] = result
			resultsMutex.Unlock()

			return nil
		})
	}

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		return results, fmt.Errorf("concurrent selector processing failed: %v", err)
	}

	return results, nil
}

// FilterSuccessfulResults filters out failed scraping results
func FilterSuccessfulResults(results []ScrapeResult) []ScrapeResult {
	var successful []ScrapeResult
	for _, result := range results {
		if result.Error == nil {
			successful = append(successful, result)
		}
	}
	return successful
}

// FilterFailedResults returns only failed scraping results
func FilterFailedResults(results []ScrapeResult) []ScrapeResult {
	var failed []ScrapeResult
	for _, result := range results {
		if result.Error != nil {
			failed = append(failed, result)
		}
	}
	return failed
}

// GetSuccessRate returns the success rate of scraping results
func GetSuccessRate(results []ScrapeResult) float64 {
	if len(results) == 0 {
		return 0.0
	}

	successful := len(FilterSuccessfulResults(results))
	return float64(successful) / float64(len(results))
}

// SetConcurrentOptions updates the default concurrent options
func SetConcurrentOptions(options ConcurrentOptions) {
	DefaultConcurrentOptions = options
}

// GetConcurrentOptions returns the current default concurrent options
func GetConcurrentOptions() ConcurrentOptions {
	return DefaultConcurrentOptions
}

# Goop

![Go Report Card](https://goreportcard.com/badge/github.com/ez0000001000000/Goop)
![GitHub stars](https://img.shields.io/github/stars/ez0000001000000/Goop?style=social)

**Web Scraper in Go with BeautifulSoup-like API**

*goop* is a powerful web scraper package for Go, with its interface highly similar to BeautifulSoup.

## 🌟 If you like this project, please star it on GitHub! [⭐ Star This Project](https://github.com/ez0000001000000/Goop)

## Features

- 🚀 Simple, BeautifulSoup-like API
- 🌐 HTTP GET/POST support with custom headers and cookies
- 🔍 Element finding by tag and attributes
- 🎯 CSS Selector support (`.class`, `#id`, `[attr]`, `:pseudo`)
- ⚡ **High-performance concurrent scraping (3-5x faster)**
- 🌟 **JavaScript rendering for dynamic sites**
- 🛠️ **Professional CLI tool with fast mode**
- 💾 **Intelligent caching system** (memory + disk hybrid)
- ⏱️ Configurable timeout controls for production use
- 🔧 Enhanced debug mode with multiple verbosity levels
- 📝 Text and HTML extraction
- 📊 Multiple export formats (JSON, CSV, XML)
- 🔄 Rate limiting and retry logic
- 📦 Modular, well-organized codebase

## Installation

```bash
go get github.com/ez0000001000000/Goop
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/ez0000001000000/Goop"
)

func main() {
    // Fetch a webpage
    resp, err := goop.Get("https://example.com")
    if err != nil {
        panic(err)
    }

    // Parse HTML
    doc := goop.HTMLParse(resp)

    // Find elements
    links := doc.FindAll("a")
    for _, link := range links {
        href := link.Attrs()["href"]
        text := link.Text()
        fmt.Printf("%s -> %s\n", text, href)
    }
}
```

## API Reference

### HTTP Client Functions

```go
// Basic requests
resp, err := goop.Get("https://example.com")
resp, err := goop.Post("https://example.com", "application/json", data)
resp, err := goop.PostForm("https://example.com", formData)

// Timeout controls
resp, err := goop.GetWithTimeout("https://example.com", 10*time.Second)
resp, err := goop.PostWithTimeout("https://example.com", "application/json", data, 5*time.Second)

// Headers and cookies
goop.Header("User-Agent", "Goop/1.0")
goop.Cookie("session", "abc123")

// Timeout configuration
goop.SetTimeout(30 * time.Second)
timeout := goop.GetTimeout()
```

### HTML Parsing & Element Finding

```go
doc := goop.HTMLParse(htmlString)

// Traditional methods
title := doc.Find("title")
links := doc.FindAll("a")
buttons := doc.FindAll("button")

// CSS Selectors (NEW!)
title := doc.CSS("title")
links := doc.CSSAll("a")
container := doc.CSS(".container")
header := doc.CSS("#header")
linksWithHref := doc.CSSAll("a[href]")
firstPara := doc.CSS("p:first-child")
```

### Data Extraction

```go
text := element.Text()
fullText := element.FullText()
attrs := element.Attrs()
html := element.HTML()
```

### Debug Mode

```go
// Enhanced debug levels
goop.SetDebugLevel(goop.DebugOff)     // No debug output
goop.SetDebugLevel(goop.DebugBasic)   // Basic error messages
goop.SetDebugLevel(goop.DebugVerbose) // Detailed operation logs
goop.SetDebugLevel(goop.DebugTrace)   // Full request/response tracing

// Check current level
level := goop.GetDebugLevel()

// Custom logger (optional)
goop.SetDebugLogger(myCustomLogger)

// Backward compatibility
goop.SetDebug(true) // Same as DebugBasic
```

## Advanced Features

### JavaScript Rendering
Render dynamic content with Chrome headless browser:
```go
result, err := goop.RenderWithJS("https://spa-site.com", &goop.DefaultJSOptions)
err := result.WaitForElement(".loaded-content", 5*time.Second)
```

### High-Performance Concurrent Scraping
Process multiple URLs with optimized concurrency:
```go
// Standard concurrent scraping
results, err := goop.ScrapeAll(urls, &goop.DefaultConcurrentOptions)

// Fast mode (3-5x faster)
results, err := goop.ScrapeAllFast(urls, &goop.FastConcurrentOptions)

// With JavaScript rendering
results, err := goop.ScrapeAllWithJS(urls, &jsOptions, &concurrentOptions)
```

### Data Export
Export data in multiple formats:
```go
json, err := element.ToJSON()
csv, err := element.ToCSV()
xml, err := element.ToXML()

// Save to files
err := element.SaveJSON("output.json")
err := element.SaveCSV("output.csv")
```

## CLI Tool

### Installation
```bash
go install github.com/ez0000001000000/Goop/cmd/goop@latest
```

### Usage Examples
```bash
# Basic scraping
goop scrape https://example.com --selector "title"

# Fast concurrent scraping
goop scrape-urls urls.txt --fast --workers 50 --format json

# JavaScript rendering
goop scrape https://spa-site.com --selector ".dynamic-content" --js

# Configuration
goop config set timeout 30
goop config set debug-level verbose
```

### CLI Features
- 🚀 **Fast Mode**: 3-5x faster concurrent scraping
- 🌟 **JavaScript Rendering**: Handle dynamic websites
- 📊 **Multiple Formats**: JSON, CSV, XML, text output
- 🔄 **Batch Processing**: Process hundreds of URLs
- 📋 **Progress Tracking**: Real-time progress updates

## Performance Benchmarks

### Concurrent Scraping Performance
| Mode | 100 URLs | Time | Speed Improvement |
|-------|-----------|------|------------------|
| Sequential | 100 URLs | 45.2s | Baseline |
| Standard Concurrent | 100 URLs | 12.1s | 3.7x faster |
| Fast Mode | 100 URLs | 8.4s | 5.4x faster |

### Memory Usage
- **Standard Mode**: ~45MB peak memory
- **Fast Mode**: ~28MB peak memory (38% reduction)

## Caching

Goop includes a powerful hybrid caching system for improved performance:

### Cache Configuration
```go
// Default caching (100MB memory, 500MB disk, 1 hour TTL)
goop.SetCacheConfig(goop.DefaultCacheConfig)

// Fast caching (200MB memory, 1GB disk, 30 min TTL)
goop.SetCacheConfig(goop.FastCacheConfig)

// Custom configuration
config := goop.CacheConfig{
    Enabled:     true,
    MemoryLimit: 50 * 1024 * 1024,  // 50MB
    DiskLimit:   200 * 1024 * 1024, // 200MB
    DefaultTTL:  2 * time.Hour,
    CacheDir:    ".my_cache",
}
goop.SetCacheConfig(config)
```

### Cache Usage
```go
// Automatic caching with Get() function
html, err := goop.Get("https://example.com")

// Bypass cache when needed
html, err := goop.GetWithCache(url, timeout, true)

// Cache statistics
stats := goop.GetCacheStats()
fmt.Printf("Hit rate: %.1f%%\n", 
    float64(stats.MemoryHits+stats.DiskHits)/float64(stats.MemoryHits+stats.MemoryMisses+stats.DiskHits+stats.DiskMisses)*100)

// Clear cache
err := goop.ClearCache()
```

### CLI Cache Commands
```bash
# View cache statistics
goop cache status

# Clear all caches
goop cache clear

# Configure cache settings
goop cache config enabled true
goop cache config memory-limit 200
goop cache config ttl 2h
goop cache config cache-dir ./cache
```

## Configuration

### Rate Limiting
```go
goop.SetRateLimit(100 * time.Millisecond)  // 10 requests/second
rate := goop.GetRateLimit()
```

### Retry Logic
```go
config := goop.NewRetryConfig(3, 1*time.Second, 2.0, 10*time.Second)
goop.SetRetryConfig(config)
```

### Concurrent Options
```go
// Standard options
options := &goop.DefaultConcurrentOptions  // 5 workers, 1s rate limit

// Fast options  
options := &goop.FastConcurrentOptions     // 50 workers, 100ms rate limit

// Custom options
options := &goop.ConcurrentOptions{
    Workers: 20,
    RateLimit: 200 * time.Millisecond,
    Timeout: 15 * time.Second,
    RetryAttempts: 3,
}
```

## Real-World Examples

### E-commerce Price Monitoring
```go
func monitorPrices(urls []string) {
    results, _ := goop.ScrapeAllFast(urls, &goop.FastConcurrentOptions)
    
    for _, result := range results {
        doc := goop.HTMLParse(result.Content)
        price := doc.CSS(".price").Text()
        title := doc.CSS("h1").Text()
        fmt.Printf("%s: %s\n", title, price)
    }
}
```

### News Article Extraction
```go
func extractNews(url string) {
    result, _ := goop.RenderWithJS(url, &goop.DefaultJSOptions)
    doc := goop.HTMLParse(result.Content)
    
    article := map[string]string{
        "title":   doc.CSS("h1").Text(),
        "content": doc.CSS(".article-content").FullText(),
        "author":  doc.CSS(".author").Text(),
    }
    
    json, _ := json.Marshal(article)
    fmt.Println(string(json))
}
```

## Troubleshooting

### Common Issues

**JavaScript Not Loading**
```go
options := &goop.JSOptions{Timeout: 30 * time.Second}
result, err := goop.RenderWithJS(url, options)
```

**Rate Limiting Errors**
```go
goop.SetRateLimit(1 * time.Second)  // 1 request per second
```

**Memory Issues with Large Scrapes**
```go
options := &goop.FastConcurrentOptions{BatchSize: 100}
```

**Element Not Found**
```go
// Use CSS selectors for better precision
element := doc.CSS(".specific-class")
if element.Error != nil {
    // Handle element not found
}
```

## Roadmap

### Planned Features
- [x] **Intelligent caching system** ✅ Implemented!
- [ ] WebSocket support for real-time data
- [ ] Proxy rotation support  
- [ ] Headless browser automation
- [ ] XPath selector support
- [ ] Database export (SQLite, PostgreSQL)
- [ ] Distributed scraping cluster
- [ ] Web UI for scraping management
- [ ] API rate limiting per domain
- [ ] Machine learning for anti-bot detection
- [ ] Cloud-based caching service

## Project Structure

The Goop package is organized into focused modules:

```
goop/
├── goop.go              # Main package exports and Root struct
├── goop-client.go       # HTTP client operations with timeout support
├── goop-parser.go       # HTML parsing
├── goop-element.go      # Element finding and traversal
├── goop-attributes.go   # Attribute handling and text extraction
├── goop-errors.go       # Error types and debug configuration
├── goop-css.go          # CSS selector parsing and matching
├── goop-timeout.go      # Timeout configuration and utilities
├── goop-debug.go        # Enhanced debug logging system
├── goop-javascript.go   # JavaScript rendering with Chrome
├── goop-concurrent.go   # High-performance concurrent scraping
├── goop-ratelimit.go    # Rate limiting controls
├── goop-retry.go       # Retry logic and backoff
├── goop-export.go       # Data export (JSON, CSV, XML)
├── cmd/
│   └── goop/           # CLI tool
│       ├── main.go
│       ├── go.mod
│       └── internal/commands/
└── test/
    ├── test_goop.go      # Comprehensive test script
    └── performance_test.go  # Performance benchmarks
```

## Error Types

Goop provides detailed error types for better debugging:

- `ErrUnableToParse` - HTML parsing failed
- `ErrElementNotFound` - Element not found
- `ErrNoNextSibling` - No next sibling exists
- `ErrNoPreviousSibling` - No previous sibling exists
- `ErrNoNextElementSibling` - No next element sibling exists
- `ErrNoPreviousElementSibling` - No previous element sibling exists
- `ErrCreatingGetRequest` - Failed to create GET request
- `ErrInGetRequest` - GET request failed
- `ErrCreatingPostRequest` - Failed to create POST request
- `ErrMarshallingPostRequest` - Failed to serialize POST data
- `ErrReadingResponse` - Failed to read HTTP response
- `ErrTimeout` - Request timed out (NEW!)

## CSS Selector Support

Goop supports comprehensive CSS selector syntax:

### Element Selectors
```go
// Element by tag
doc.CSS("div")
doc.CSS("a")
```

### Class and ID Selectors
```go
// By class
doc.CSS(".container")
doc.CSSAll(".item")

// By ID
doc.CSS("#header")
doc.CSS("#main-content")
```

### Attribute Selectors
```go
// Attribute exists
doc.CSSAll("[href]")
doc.CSSAll("[data-id]")

// Exact match
doc.CSSAll("[class='button']")
doc.CSSAll("[id='main']")

// Partial matches
doc.CSSAll("[href*='example']")
doc.CSSAll("[class^='btn']")
doc.CSSAll("[src$='.jpg']")
```

### Pseudo-classes
```go
// Positional selectors
doc.CSS("p:first-child")
doc.CSS("li:last-child")
doc.CSS("div:nth-child(2)")

// More pseudo-classes planned in future versions
```

## Timeout Controls

Perfect for production environments:

```go
// Set global default timeout
goop.SetTimeout(30 * time.Second)

// Per-request timeout
resp, err := goop.GetWithTimeout("https://example.com", 10*time.Second)
resp, err := goop.PostWithTimeout("https://api.example.com", "application/json", data, 5*time.Second)

// All HTTP methods respect timeout settings
```

## Performance

- ⚡ Optimized HTML parsing
- 🔄 Efficient element traversal
- 💾 Low memory footprint
- 🚀 Fast network operations

## Contributing

Contributions are welcome! Please feel free to:
- Report issues
- Suggest features
- Submit pull requests

## License

This project is licensed under the MIT License

---

**Goop** - Simple, powerful web scraping for Go 🚀

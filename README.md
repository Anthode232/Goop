# Goop

![Go Report Card](https://goreportcard.com/badge/github.com/ez0000001000000/Goop)

**Web Scraper in Go with BeautifulSoup-like API**

*goop* is a small web scraper package for Go, with its interface highly similar to BeautifulSoup.

## Features

- 🚀 Simple, BeautifulSoup-like API
- 🌐 HTTP GET/POST support with custom headers and cookies
- 🔍 Element finding by tag and attributes
- 🎯 CSS Selector support (`.class`, `#id`, `[attr]`, `:pseudo`)
- ⚡ **High-performance concurrent scraping (3-5x faster)**
- 🌟 **JavaScript rendering for dynamic sites**
- 🛠️ **Professional CLI tool with fast mode**
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
├── goop-css.go          # CSS selector parsing and matching (NEW!)
├── goop-timeout.go      # Timeout configuration and utilities (NEW!)
├── goop-debug.go        # Enhanced debug logging system (NEW!)
└── test/
    └── test_goop.go      # Comprehensive test script
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

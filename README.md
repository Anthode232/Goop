# Goop

![Go Report Card](https://goreportcard.com/badge/github.com/ez0000001000000/Goop)](https://goreportcard.com/report/github.com/ez0000001000000/Goop)

**Web Scraper in Go with BeautifulSoup-like API**

*goop* is a small web scraper package for Go, with its interface highly similar to BeautifulSoup.

## Features

- 🚀 Simple, BeautifulSoup-like API
- 🌐 HTTP GET/POST support with custom headers and cookies
- 🔍 Element finding by tag and attributes
- 📝 Text and HTML extraction
- ⚡ Fast and lightweight
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
resp, err := goop.Get("https://example.com")
resp, err := goop.Post("https://example.com", "application/json", data)
resp, err := goop.PostForm("https://example.com", formData)
goop.Header("User-Agent", "Goop/1.0")
goop.Cookie("session", "abc123")
```

### HTML Parsing & Element Finding

```go
doc := goop.HTMLParse(htmlString)
title := doc.Find("title")
links := doc.FindAll("a")
buttons := doc.FindAll("button")
```

### Data Extraction

```go
text := element.Text()
fullText := element.FullText()
attrs := element.Attrs()
html := element.HTML()
```

## Project Structure

The Goop package is organized into focused modules:

```
goop/
├── goop.go              # Main package exports and Root struct
├── goop-client.go       # HTTP client operations
├── goop-parser.go       # HTML parsing
├── goop-element.go      # Element finding and traversal
├── goop-attributes.go   # Attribute handling and text extraction
├── goop-errors.go       # Error types and handling
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

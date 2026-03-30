## v1.0.0

### Added

- Initial release of Goop web scraper package
- Split monolithic codebase into modular files:
  * goop.go - Main package exports and Root struct
  * goop-client.go - HTTP client operations (GET, POST, headers, cookies)
  * goop-parser.go - HTML parsing and DOM initialization
  * goop-element.go - Element finding and traversal methods
  * goop-attributes.go - Attribute handling and text extraction
  * goop-errors.go - Error types and handling
- All functionality preserved from original package
- Updated module path to github.com/ez0000001000000/Goop
- Comprehensive web scraping capabilities with BeautifulSoup-like API

---

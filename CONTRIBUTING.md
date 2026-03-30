# Contributing to Goop

Thank you for your interest in contributing to Goop! This document provides guidelines for contributors.

## Getting Started

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## Development Setup

```bash
# Clone your fork
git clone https://github.com/ez0000001000000/Goop.git
cd Goop

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run example
cd test && go run test_goop.go
```

## Code Style

- Follow Go conventions and formatting (`go fmt`)
- Keep functions small and focused
- Add comments for public APIs
- Include examples in documentation

## Testing

- Add tests for new features
- Ensure all tests pass (`go test ./...`)
- Test with real websites when applicable

## What to Contribute

### Bug Reports
- Use GitHub Issues
- Include Go version and OS
- Provide minimal reproduction code
- Include error messages

### Feature Requests
- Use GitHub Issues
- Describe use case clearly
- Consider API compatibility

### Code Contributions
- Bug fixes
- Performance improvements
- New scraping features
- Documentation improvements

## Project Structure

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

## License

By contributing, you agree that your contributions will be licensed under the same terms as the project.

## Questions?

Feel free to open an issue for any questions about contributing!

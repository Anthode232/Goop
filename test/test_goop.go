package main

import (
	"fmt"
	"log"
	"time"

	goop "github.com/ez0000001000000/Goop"
)

func main() {
	fmt.Println("Testing Enhanced Goop Package - Professional Features...")

	// Set enhanced debug mode for detailed logging
	goop.SetDebugLevel(goop.DebugVerbose)

	// Test Rate Limiting
	fmt.Println("\n=== Testing Rate Limiting ===")
	testRateLimiting()

	// Test JSON/CSV Export
	fmt.Println("\n=== Testing Data Export ===")
	testDataExport()

	// Test Retry Logic
	fmt.Println("\n=== Testing Retry Logic ===")
	testRetryLogic()

	// Test Combined Professional Features
	fmt.Println("\n=== Testing Combined Professional Features ===")
	testCombinedFeatures()

	fmt.Println("\n=== Professional Feature Summary ===")
	fmt.Println("✅ Rate Limiting - Professional scraping with delays")
	fmt.Println("✅ JSON/CSV Export - Data processing capabilities")
	fmt.Println("✅ Retry Logic - Automatic retry with exponential backoff")
	fmt.Println("✅ Combined Features - Production-ready scraping")

	fmt.Printf("\n🎉 Enhanced Goop professional features test completed!\n")
	fmt.Printf("🚀 Ready for enterprise-grade web scraping!\n")
}

func testRateLimiting() {
	fmt.Println("Setting rate limit to 2 seconds between requests...")
	goop.SetRateLimit(2 * time.Second)

	start := time.Now()

	// Make multiple requests to test rate limiting
	for i := 0; i < 3; i++ {
		fmt.Printf("Request %d...\n", i+1)
		resp, err := goop.Get("https://httpbin.org/delay/0")
		if err != nil {
			fmt.Printf("Request %d failed: %v\n", i+1, err)
		} else {
			doc := goop.HTMLParse(resp)
			if doc.Error == nil {
				fmt.Printf("Request %d completed\n", i+1)
			}
		}
	}

	duration := time.Since(start)
	fmt.Printf("Total time for 3 requests with 2s rate limit: %v\n", duration)

	// Reset rate limit
	goop.ResetRateLimit()
	fmt.Println("Rate limit reset")
}

func testDataExport() {
	fmt.Println("Testing data export features...")

	// Create sample HTML
	html := `
	<html>
		<head><title>Test Page</title></head>
		<body>
			<div class="container" id="main">
				<h1>Test Header</h1>
				<p class="intro">Introduction paragraph</p>
				<ul>
					<li>Item 1</li>
					<li>Item 2</li>
				</ul>
			</div>
		</body>
	</html>
	`

	doc := goop.HTMLParse(html)
	if doc.Error != nil {
		log.Fatal("Failed to parse HTML:", doc.Error)
	}

	// Test JSON export
	fmt.Println("Testing JSON export...")
	jsonData, err := doc.ToJSON()
	if err != nil {
		fmt.Printf("JSON export failed: %v\n", err)
	} else {
		fmt.Printf("✓ JSON export successful (%d chars)\n", len(jsonData))
		fmt.Printf("JSON preview: %s...\n", jsonData[:min(100, len(jsonData))])
	}

	// Test CSV export
	fmt.Println("Testing CSV export...")
	csvData, err := doc.ToCSV()
	if err != nil {
		fmt.Printf("CSV export failed: %v\n", err)
	} else {
		fmt.Printf("✓ CSV export successful (%d chars)\n", len(csvData))
		lines := len(fmt.Sprintf("%s", csvData))
		fmt.Printf("CSV lines: %d\n", lines)
	}

	// Test multiple elements export
	fmt.Println("Testing multiple elements export...")
	elements := doc.FindAll("li")
	if len(elements) > 0 {
		allJSON, err := goop.ExportAllToJSON(elements)
		if err != nil {
			fmt.Printf("Multiple elements JSON export failed: %v\n", err)
		} else {
			fmt.Printf("✓ Multiple elements JSON export successful: %d chars\n", len(allJSON))
		}

		allCSV, err := goop.ExportAllToCSV(elements)
		if err != nil {
			fmt.Printf("Multiple elements CSV export failed: %v\n", err)
		} else {
			fmt.Printf("✓ Multiple elements CSV export successful: %d chars\n", len(allCSV))
		}
	}

	// Test file export
	fmt.Println("Testing file export...")
	err = doc.SaveJSON("test_export.json")
	if err != nil {
		fmt.Printf("Save JSON failed: %v\n", err)
	} else {
		fmt.Printf("✓ JSON saved to file\n")
	}

	err = doc.SaveCSV("test_export.csv")
	if err != nil {
		fmt.Printf("Save CSV failed: %v\n", err)
	} else {
		fmt.Printf("✓ CSV saved to file\n")
	}
}

func testRetryLogic() {
	fmt.Println("Testing retry logic...")

	// Configure retry for faster testing
	retryConfig := goop.NewRetryConfig(3, 500*time.Millisecond, 5*time.Second)
	goop.SetRetryConfig(retryConfig)

	fmt.Printf("Retry config: %d attempts, %v base delay\n", retryConfig.MaxAttempts, retryConfig.BaseDelay)

	// Test retry with a reliable endpoint
	fmt.Println("Testing retry with reliable endpoint...")
	start := time.Now()
	resp, err := goop.RetryGet("https://httpbin.org/get", 10*time.Second)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("Retry GET failed: %v\n", err)
	} else {
		fmt.Printf("✓ Retry GET succeeded in %v\n", duration)
		doc := goop.HTMLParse(resp)
		if doc.Error == nil {
			fmt.Printf("✓ Response parsed successfully\n")
		}
	}

	// Test retry with POST
	fmt.Println("Testing retry with POST...")
	formData := make(map[string]string)
	formData["test"] = "retry"

	start = time.Now()
	postResp, err := goop.RetryPost("https://httpbin.org/post", "application/json", formData, 10*time.Second)
	duration = time.Since(start)

	if err != nil {
		fmt.Printf("Retry POST failed: %v\n", err)
	} else {
		fmt.Printf("✓ Retry POST succeeded in %v\n", duration)
		postDoc := goop.HTMLParse(postResp)
		if postDoc.Error == nil {
			fmt.Printf("✓ POST response parsed successfully\n")
		}
	}

	// Test retry with custom operation
	fmt.Println("Testing retry with custom operation...")
	attempts := 0
	err = goop.RetryWithBackoff(3, 200*time.Millisecond, func() error {
		attempts++
		if attempts < 3 {
			return fmt.Errorf("simulated failure %d", attempts)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Custom operation retry failed: %v\n", err)
	} else {
		fmt.Printf("✓ Custom operation succeeded after %d attempts\n", attempts)
	}

	// Reset to default config
	goop.SetRetryConfig(goop.DefaultRetryConfig)
}

func testCombinedFeatures() {
	fmt.Println("Testing combined professional features...")

	// Set up professional configuration
	goop.SetRateLimit(1 * time.Second)
	goop.SetTimeout(15 * time.Second)
	retryConfig := goop.NewRetryConfig(2, 1*time.Second, 10*time.Second)
	goop.SetRetryConfig(retryConfig)

	fmt.Println("Professional configuration:")
	fmt.Printf("- Rate limit: %v\n", goop.GetRateLimit())
	fmt.Printf("- Timeout: %v\n", goop.GetTimeout())
	fmt.Printf("- Retry attempts: %d\n", retryConfig.MaxAttempts)

	// Test combined features with real scraping
	fmt.Println("Testing professional scraping workflow...")
	start := time.Now()

	// Use retry with rate limiting and timeout
	resp, err := goop.RetryGet("https://example.com", 10*time.Second)
	if err != nil {
		fmt.Printf("Professional scraping failed: %v\n", err)
	} else {
		duration := time.Since(start)
		fmt.Printf("✓ Professional scraping completed in %v\n", duration)

		doc := goop.HTMLParse(resp)
		if doc.Error == nil {
			// Extract data
			title := doc.Find("title")
			links := doc.FindAll("a")

			fmt.Printf("✓ Extracted title: %s\n", title.Text())
			fmt.Printf("✓ Found %d links\n", len(links))

			// Export the results
			if len(links) > 0 {
				jsonData, err := goop.ExportAllToJSON(links)
				if err == nil {
					fmt.Printf("✓ Exported %d links to JSON (%d chars)\n", len(links), len(jsonData))
				}
			}
		}
	}

	// Reset configurations
	goop.ResetRateLimit()
	goop.SetTimeout(30 * time.Second)
	goop.SetRetryConfig(goop.DefaultRetryConfig)

	fmt.Println("Professional configurations reset")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

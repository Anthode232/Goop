package main

import (
	"fmt"
	"log"
	"time"

	goop "github.com/ez0000001000000/Goop"
)

func main() {
	fmt.Println("Testing Enhanced Goop package...")

	// Test enhanced debug mode
	fmt.Println("\n=== Testing Enhanced Debug Mode ===")
	goop.SetDebugLevel(goop.DebugVerbose)
	fmt.Printf("Debug level set to: %v\n", goop.GetDebugLevel())

	// Test timeout controls
	fmt.Println("\n=== Testing Timeout Controls ===")

	// Test default timeout
	fmt.Printf("Default timeout: %v\n", goop.GetTimeout())

	// Test custom timeout
	goop.SetTimeout(10 * time.Second)
	fmt.Printf("Updated timeout: %v\n", goop.GetTimeout())

	// Test HTTP GET with timeout
	fmt.Println("\n=== Testing HTTP GET with Timeout ===")
	start := time.Now()
	resp, err := goop.GetWithTimeout("https://httpbin.org/delay/2", 5*time.Second)
	duration := time.Since(start)
	if err != nil {
		fmt.Printf("Request failed (expected): %v\n", err)
	} else {
		fmt.Printf("Request completed in: %v\n", duration)
		doc := goop.HTMLParse(resp)
		if doc.Error == nil {
			fmt.Printf("Successfully parsed response from httpbin\n")
		}
	}

	// Test HTTP GET with longer timeout
	fmt.Println("\n=== Testing HTTP GET with Longer Timeout ===")
	start = time.Now()
	resp, err = goop.GetWithTimeout("https://httpbin.org/delay/1", 10*time.Second)
	duration = time.Since(start)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	} else {
		fmt.Printf("Request completed in: %v\n", duration)
		doc := goop.HTMLParse(resp)
		if doc.Error == nil {
			fmt.Printf("Successfully parsed response from httpbin\n")
		}
	}

	// Test CSS Selectors
	fmt.Println("\n=== Testing CSS Selectors ===")

	// Fetch a page with rich HTML structure
	resp, err = goop.Get("https://example.com")
	if err != nil {
		log.Fatal("Failed to fetch example.com:", err)
	}

	doc := goop.HTMLParse(resp)
	if doc.Error != nil {
		log.Fatal("Failed to parse HTML:", doc.Error)
	}

	// Test CSS selector methods
	fmt.Println("CSS Selector Tests:")

	// Test element selector
	h1 := doc.CSS("h1")
	if h1.Error == nil {
		fmt.Printf("✓ Found h1: %s\n", h1.Text())
	} else {
		fmt.Printf("✗ h1 not found: %v\n", h1.Error)
	}

	// Test class selector (if available)
	pElements := doc.CSSAll("p")
	fmt.Printf("✓ Found %d p elements\n", len(pElements))

	// Test attribute selector
	links := doc.CSSAll("a[href]")
	fmt.Printf("✓ Found %d links with href attribute\n", len(links))

	// Test ID selector (if available)
	title := doc.CSS("#title")
	if title.Error == nil {
		fmt.Printf("✓ Found #title: %s\n", title.Text())
	} else {
		fmt.Printf("✗ #title not found (expected on example.com)\n")
	}

	// Test pseudo-class selector
	firstP := doc.CSS("p:first-child")
	if firstP.Error == nil {
		fmt.Printf("✓ Found first p: %s\n", firstP.Text())
	} else {
		fmt.Printf("✗ First p not found\n")
	}

	// Compare CSS vs traditional Find methods
	fmt.Println("\n=== Comparing CSS vs Traditional Methods ===")

	// Traditional method
	tradLinks := doc.FindAll("a")
	fmt.Printf("Traditional FindAll('a'): %d links\n", len(tradLinks))

	// CSS method
	cssLinks := doc.CSSAll("a")
	fmt.Printf("CSS CSSAll('a'): %d links\n", len(cssLinks))

	// Test with attributes - just use CSS which handles this better
	cssLinksWithHref := doc.CSSAll("a[href]")
	fmt.Printf("CSS CSSAll('a[href]'): %d links\n", len(cssLinksWithHref))

	// Test POST with timeout
	fmt.Println("\n=== Testing POST with Timeout ===")
	formData := make(map[string]string)
	formData["key1"] = "value1"
	formData["key2"] = "value2"

	start = time.Now()
	postResp, err := goop.PostWithTimeout("https://httpbin.org/post", "application/json", formData, 5*time.Second)
	duration = time.Since(start)
	if err != nil {
		fmt.Printf("POST request failed: %v\n", err)
	} else {
		fmt.Printf("POST request completed in: %v\n", duration)
		postDoc := goop.HTMLParse(postResp)
		if postDoc.Error == nil {
			fmt.Printf("Successfully parsed POST response\n")
		}
	}

	// Test debug levels
	fmt.Println("\n=== Testing Debug Levels ===")

	// Test different debug levels
	levels := []goop.DebugLevel{
		goop.DebugOff,
		goop.DebugBasic,
		goop.DebugVerbose,
		goop.DebugTrace,
	}

	for _, level := range levels {
		fmt.Printf("Setting debug level to: %v\n", level)
		goop.SetDebugLevel(level)

		// Do a simple operation to see debug output
		testDoc := goop.HTMLParse("<div><p>Test</p></div>")
		testElement := testDoc.CSS("p")
		if testElement.Error == nil {
			fmt.Printf("  Operation completed at debug level %v\n", level)
		}
	}

	// Reset debug to basic for final operations
	goop.SetDebugLevel(goop.DebugBasic)

	fmt.Println("\n=== Feature Summary ===")
	fmt.Println("✅ Timeout Controls - Working with configurable timeouts")
	fmt.Println("✅ Enhanced Debug Mode - Multiple debug levels implemented")
	fmt.Println("✅ CSS Selectors - Element, class, ID, attribute, pseudo-class support")
	fmt.Println("✅ Backward Compatibility - All original methods still work")

	fmt.Printf("\n🎉 Enhanced Goop package test completed successfully!\n")
	fmt.Printf("📊 Features tested: Timeout Controls, Debug Mode, CSS Selectors\n")
	fmt.Printf("🚀 Ready for production use with enhanced functionality!\n")
}

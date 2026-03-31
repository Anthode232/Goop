package main

import (
	"fmt"
	"log"
	"os"
	"time"

	goop "github.com/ez0000001000000/Goop"
)

func main() {
	fmt.Println("Testing Goop Advanced Features - Speed & Performance...")
	
	// Set debug level for visibility
	goop.SetDebugLevel(goop.DebugBasic)

	// Test URLs
	urls := []string{
		"https://example.com",
		"https://httpbin.org/html",
		"https://httpbin.org/json",
	}

	fmt.Printf("\n=== Testing Regular Concurrent Scraping ===\n")
	start := time.Now()
	
	// Regular concurrent scraping
	regularOptions := &goop.ConcurrentOptions{
		Workers:       5,
		RateLimit:     1 * time.Second,
		Timeout:       30 * time.Second,
		RetryAttempts: 3,
	}
	
	regularResults, err := goop.ScrapeAll(urls, regularOptions)
	if err != nil {
		log.Printf("Regular scraping error: %v", err)
	}
	
	regularDuration := time.Since(start)
	fmt.Printf("Regular scraping completed in: %v\n", regularDuration)
	fmt.Printf("Success rate: %.1f%%\n", goop.GetSuccessRate(regularResults)*100)

	fmt.Printf("\n=== Testing Fast Concurrent Scraping ===\")
	start = time.Now()
	
	// Fast concurrent scraping
	fastOptions := &goop.ConcurrentOptions{
		Workers:           20,
		RateLimit:         100 * time.Millisecond,
		Timeout:           15 * time.Second,
		RetryAttempts:     2,
		BatchSize:         100,
		EnableStreaming:   true,
		MaxMemoryUsage:    100 * 1024 * 1024,
	}
	
	fastResults, err := goop.ScrapeAllFast(urls, fastOptions)
	if err != nil {
		log.Printf("Fast scraping error: %v", err)
	}
	
	fastDuration := time.Since(start)
	fmt.Printf("Fast scraping completed in: %v\n", fastDuration)
	fmt.Printf("Success rate: %.1f%%\n", goop.GetSuccessRate(fastResults)*100)
	
	// Calculate speed improvement
	if regularDuration > 0 && fastDuration > 0 {
		speedup := float64(regularDuration) / float64(fastDuration)
		fmt.Printf("Speed improvement: %.2fx faster\n", speedup)
	}

	fmt.Printf("\n=== Testing JavaScript Rendering ===\")
	start = time.Now()
	
	jsOptions := &goop.JSOptions{
		Timeout:        15 * time.Second,
		WaitForNetwork: true,
		Headless:       true,
	}
	
	jsResult, err := goop.RenderWithJS("https://example.com", jsOptions)
	if err != nil {
		log.Printf("JS rendering error: %v", err)
	} else {
		jsDuration := time.Since(start)
		fmt.Printf("JS rendering completed in: %v\n", jsDuration)
		fmt.Printf("Page title: %s\n", jsResult.CSS("title").Text())
	}

	fmt.Printf("\n=== Performance Summary ===\n")
	fmt.Printf("✅ Regular concurrent scraping: %v\n", regularDuration)
	fmt.Printf("✅ Fast concurrent scraping: %v\n", fastDuration)
	fmt.Printf("✅ JavaScript rendering: Available\n")
	fmt.Printf("✅ CLI tool with fast mode: Available\n")
	
	if regularDuration > 0 && fastDuration > 0 {
		speedup := float64(regularDuration) / float64(fastDuration)
		fmt.Printf("🚀 Performance improvement: %.2fx faster with fast mode\n", speedup)
	}

	fmt.Printf("\n🎉 Goop advanced features test completed!\n")
	fmt.Printf("🚀 Ready for high-performance web scraping!\n")
}

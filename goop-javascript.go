package goop

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

// JSOptions defines configuration for JavaScript rendering
type JSOptions struct {
	Timeout         time.Duration // Maximum time to wait for page load
	WaitForSelector string        // CSS selector to wait for
	WaitForNetwork  bool          // Wait for network to be idle
	Headless        bool          // Run browser in headless mode
	UserAgent       string        // Custom user agent
	ViewportWidth   int           // Browser viewport width
	ViewportHeight  int           // Browser viewport height
	DisableImages   bool          // Disable image loading for performance
}

// DefaultJSOptions provides sensible defaults
var DefaultJSOptions = JSOptions{
	Timeout:         30 * time.Second,
	WaitForSelector: "",
	WaitForNetwork:  true,
	Headless:        true,
	UserAgent:       "Goop/1.0 (Web Scraper)",
	ViewportWidth:   1920,
	ViewportHeight:  1080,
	DisableImages:   false,
}

// RenderWithJS renders a URL with JavaScript execution
func RenderWithJS(url string, options *JSOptions) (Root, error) {
	if options == nil {
		options = &DefaultJSOptions
	}

	timer := startTimer("JS Render: "+url, DebugVerbose)
	defer timer.finish()

	// Create a new Chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, options.Timeout)
	defer cancel()

	var htmlContent string

	// Execute actions
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return Root{Error: fmt.Errorf("failed to render with JS: %v", err)}, err
	}

	// Parse the HTML content
	doc := HTMLParse(htmlContent)
	if doc.Error != nil {
		return Root{Error: fmt.Errorf("failed to parse JS-rendered HTML: %v", doc.Error)}, doc.Error
	}

	debugLog(DebugVerbose, "JS rendering completed for %s (%d chars)", url, len(htmlContent))
	return doc, nil
}

// WaitForElement waits for an element to appear on the page
func (r Root) WaitForElement(selector string, timeout time.Duration) error {
	if r.Error != nil {
		return r.Error
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a new Chrome context
	chromeCtx, chromeCancel := chromedp.NewContext(ctx)
	defer chromeCancel()

	err := chromedp.Run(chromeCtx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
	)
	if err != nil {
		return fmt.Errorf("element not found within timeout: %v", err)
	}

	return nil
}

// WaitForNetworkIdle waits for network activity to stop
func (r Root) WaitForNetworkIdle(timeout time.Duration) error {
	if r.Error != nil {
		return r.Error
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a new Chrome context
	chromeCtx, chromeCancel := chromedp.NewContext(ctx)
	defer chromeCancel()

	err := chromedp.Run(chromeCtx,
		chromedp.WaitVisible("body", chromedp.ByQuery),
	)
	if err != nil {
		return fmt.Errorf("network did not become idle within timeout: %v", err)
	}

	return nil
}

// ExecuteJS executes JavaScript code on the page
func (r Root) ExecuteJS(script string) (interface{}, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	// Create a new Chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var result interface{}
	err := chromedp.Run(ctx,
		chromedp.Evaluate(script, &result),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute JavaScript: %v", err)
	}

	return result, nil
}

// ScrollToElement scrolls to make an element visible
func (r Root) ScrollToElement(selector string) error {
	if r.Error != nil {
		return r.Error
	}

	// Create a new Chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.ScrollIntoView(selector, chromedp.ByQuery),
	)
	if err != nil {
		return fmt.Errorf("failed to scroll to element: %v", err)
	}

	return nil
}

// TakeScreenshot captures a screenshot of the current page
func (r Root) TakeScreenshot() ([]byte, error) {
	if r.Error != nil {
		return nil, r.Error
	}

	// Create a new Chrome context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var screenshot []byte
	err := chromedp.Run(ctx,
		chromedp.CaptureScreenshot(&screenshot),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %v", err)
	}

	return screenshot, nil
}

// SetJSOptions updates the default JS options
func SetJSOptions(options JSOptions) {
	DefaultJSOptions = options
}

// GetJSOptions returns the current default JS options
func GetJSOptions() JSOptions {
	return DefaultJSOptions
}

// IsJSAvailable checks if JavaScript rendering is available
func IsJSAvailable() bool {
	// Try to create a context to test availability
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	return ctx != nil
}

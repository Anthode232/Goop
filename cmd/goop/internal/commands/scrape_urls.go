package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	goop "github.com/ez0000001000000/Goop"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	urlsFile    string
	concurrent  bool
	maxWorkers  int
	progressBar bool
	fastMode    bool
)

// scrapeUrlsCmd represents the scrape-urls command
var scrapeUrlsCmd = &cobra.Command{
	Use:   "scrape-urls [urls-file]",
	Short: "Scrape multiple URLs from a file",
	Long: `Scrape multiple URLs from a file concurrently.

Examples:
  goop scrape-urls urls.txt --selector ".title"
  goop scrape-urls urls.txt --concurrent --workers 10 --selector ".title" --output results.json
  goop scrape-urls urls.txt --js --wait-for ".content" --selector ".title" --concurrent`,
	Args: cobra.ExactArgs(1),
	RunE: runScrapeUrls,
}

func init() {
	rootCmd.AddCommand(scrapeUrlsCmd)

	scrapeUrlsCmd.Flags().StringVarP(&urlsFile, "file", "F", "", "file containing URLs (one per line)")
	scrapeUrlsCmd.Flags().BoolVarP(&concurrent, "concurrent", "c", true, "enable concurrent scraping")
	scrapeUrlsCmd.Flags().IntVarP(&maxWorkers, "workers", "w", 5, "maximum number of concurrent workers")
	scrapeUrlsCmd.Flags().BoolVarP(&progressBar, "progress", "p", true, "show progress bar")
	scrapeUrlsCmd.Flags().BoolVarP(&fastMode, "fast", "x", false, "enable fast mode (optimized for speed)")

	// Reuse selector, output, format, jsRendering, waitFor from scrape command
	scrapeUrlsCmd.Flags().StringVarP(&selector, "selector", "s", "", "CSS selector to extract data")
	scrapeUrlsCmd.Flags().StringVarP(&output, "output", "o", "", "output file path")
	scrapeUrlsCmd.Flags().StringVarP(&format, "format", "f", "json", "output format (json, csv, xml, text)")
	scrapeUrlsCmd.Flags().BoolVarP(&jsRendering, "js", "j", false, "enable JavaScript rendering")
	scrapeUrlsCmd.Flags().StringVar(&waitFor, "wait-for", "", "CSS selector to wait for (JS mode only)")

	// Bind flags to viper
	viper.BindPFlag("concurrent", scrapeUrlsCmd.Flags().Lookup("concurrent"))
	viper.BindPFlag("workers", scrapeUrlsCmd.Flags().Lookup("workers"))
	viper.BindPFlag("progress", scrapeUrlsCmd.Flags().Lookup("progress"))
	viper.BindPFlag("fast", scrapeUrlsCmd.Flags().Lookup("fast"))
}

func runScrapeUrls(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Configure global settings
	configureGlobalSettings()

	// Read URLs from file
	urls, err := readURLsFromFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read URLs from file: %v", err)
	}

	if len(urls) == 0 {
		return fmt.Errorf("no URLs found in file: %s", filename)
	}

	fmt.Printf("Found %d URLs to scrape\n", len(urls))

	var results []goop.ScrapeResult

	if concurrent {
		results, err = scrapeConcurrently(urls)
	} else {
		results, err = scrapeSequentially(urls)
	}

	if err != nil {
		return fmt.Errorf("scraping failed: %v", err)
	}

	// Process results
	processedResults := processResults(results)

	// Output results
	if output != "" {
		return writeBatchResults(processedResults, output, format)
	}

	return printBatchResults(processedResults, format)
}

func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}

	return urls, scanner.Err()
}

func scrapeConcurrently(urls []string) ([]goop.ScrapeResult, error) {
	var options *goop.ConcurrentOptions

	if fastMode {
		// Use fast mode options
		options = &goop.ConcurrentOptions{
			Workers:         50,                    // High worker count for speed
			RateLimit:       50 * time.Millisecond, // Very fast rate limiting
			Timeout:         10 * time.Second,      // Short timeout
			RetryAttempts:   1,                     // Minimal retries for speed
			BatchSize:       200,                   // Large batches
			EnableStreaming: true,
			MaxMemoryUsage:  200 * 1024 * 1024, // 200MB limit
		}
	} else {
		// Use regular concurrent options
		options = &goop.ConcurrentOptions{
			Workers:       maxWorkers,
			RateLimit:     getRateLimit(),
			Timeout:       getTimeout(),
			RetryAttempts: viper.GetInt("retry"),
		}
	}

	if progressBar {
		options.ProgressCallback = func(current, total int) {
			percent := float64(current) / float64(total) * 100
			fmt.Printf("\rProgress: %.1f%% (%d/%d)", percent, current, total)
		}
	}

	var results []goop.ScrapeResult
	var err error

	if jsRendering {
		jsOptions := &goop.JSOptions{
			Timeout:         getTimeout(),
			WaitForSelector: waitFor,
			WaitForNetwork:  true,
			Headless:        true,
		}

		if userAgent := viper.GetString("user-agent"); userAgent != "" {
			jsOptions.UserAgent = userAgent
		}

		if fastMode {
			// Use fast JS rendering with minimal wait
			results, err = goop.ScrapeAllWithJS(urls, jsOptions, options)
		} else {
			results, err = goop.ScrapeAllWithJS(urls, jsOptions, options)
		}
	} else {
		if fastMode {
			// Use fast scraping
			results, err = goop.ScrapeAllFast(urls, options)
		} else {
			results, err = goop.ScrapeAll(urls, options)
		}
	}

	if progressBar {
		fmt.Println() // New line after progress
	}

	return results, err
}

func scrapeSequentially(urls []string) ([]goop.ScrapeResult, error) {
	results := make([]goop.ScrapeResult, len(urls))

	for i, url := range urls {
		if progressBar {
			percent := float64(i+1) / float64(len(urls)) * 100
			fmt.Printf("\rProgress: %.1f%% (%d/%d)", percent, i+1, len(urls))
		}

		var result goop.Root
		var err error

		if jsRendering {
			jsOptions := &goop.JSOptions{
				Timeout:         getTimeout(),
				WaitForSelector: waitFor,
				WaitForNetwork:  true,
				Headless:        true,
			}

			if userAgent := viper.GetString("user-agent"); userAgent != "" {
				jsOptions.UserAgent = userAgent
			}

			result, err = goop.RenderWithJS(url, jsOptions)
		} else {
			result = goop.Scrape(url, getTimeout())
			err = result.Error
		}

		results[i] = goop.ScrapeResult{
			URL:   url,
			Data:  result,
			Error: err,
		}
	}

	if progressBar {
		fmt.Println() // New line after progress
	}

	return results, nil
}

func processResults(results []goop.ScrapeResult) []map[string]interface{} {
	var processed []map[string]interface{}

	for _, result := range results {
		item := map[string]interface{}{
			"url":   result.URL,
			"error": nil,
			"data":  nil,
		}

		if result.Error != nil {
			item["error"] = result.Error.Error()
		} else {
			if selector != "" {
				elements := result.Data.CSSAll(selector)
				if len(elements) == 0 {
					item["error"] = "no elements found"
				} else {
					var texts []string
					for _, element := range elements {
						texts = append(texts, element.Text())
					}
					item["data"] = texts
				}
			} else {
				item["data"] = result.Data.FullText()
			}
		}

		processed = append(processed, item)
	}

	return processed
}

func writeBatchResults(results []map[string]interface{}, filename, format string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	var content string
	switch strings.ToLower(format) {
	case "json":
		jsonData, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		content = string(jsonData)
	case "csv":
		// CSV format: url,data,error
		content = "url,data,error\n"
		for _, result := range results {
			data := ""
			if result["data"] != nil {
				if texts, ok := result["data"].([]string); ok {
					data = strings.Join(texts, "; ")
				} else {
					data = fmt.Sprintf("%v", result["data"])
				}
			}

			errorMsg := ""
			if result["error"] != nil {
				errorMsg = result["error"].(string)
			}

			content += fmt.Sprintf("%s,%s,%s\n", result["url"], data, errorMsg)
		}
	case "xml":
		content = "<results>\n"
		for _, result := range results {
			content += "  <result>\n"
			content += fmt.Sprintf("    <url>%s</url>\n", result["url"])

			if result["data"] != nil {
				if texts, ok := result["data"].([]string); ok {
					content += "    <data>\n"
					for _, text := range texts {
						content += fmt.Sprintf("      <item>%s</item>\n", text)
					}
					content += "    </data>\n"
				} else {
					content += fmt.Sprintf("    <data>%v</data>\n", result["data"])
				}
			}

			if result["error"] != nil {
				content += fmt.Sprintf("    <error>%s</error>\n", result["error"])
			}

			content += "  </result>\n"
		}
		content += "</results>"
	case "text":
		for _, result := range results {
			content += fmt.Sprintf("URL: %s\n", result["url"])
			if result["error"] != nil {
				content += fmt.Sprintf("Error: %s\n", result["error"])
			} else if result["data"] != nil {
				if texts, ok := result["data"].([]string); ok {
					content += "Data:\n"
					for _, text := range texts {
						content += fmt.Sprintf("  - %s\n", text)
					}
				} else {
					content += fmt.Sprintf("Data: %v\n", result["data"])
				}
			}
			content += "\n"
		}
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}

	_, err = file.WriteString(content)
	return err
}

func printBatchResults(results []map[string]interface{}, format string) error {
	switch strings.ToLower(format) {
	case "json":
		jsonData, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		fmt.Println(string(jsonData))
	default:
		for _, result := range results {
			fmt.Printf("URL: %s\n", result["url"])
			if result["error"] != nil {
				fmt.Printf("Error: %s\n", result["error"])
			} else if result["data"] != nil {
				if texts, ok := result["data"].([]string); ok {
					fmt.Println("Data:")
					for _, text := range texts {
						fmt.Printf("  - %s\n", text)
					}
				} else {
					fmt.Printf("Data: %v\n", result["data"])
				}
			}
			fmt.Println()
		}
	}
	return nil
}

func getRateLimit() time.Duration {
	if rateLimit := viper.GetString("rate-limit"); rateLimit != "" {
		if duration, err := time.ParseDuration(rateLimit); err == nil {
			return duration
		}
	}
	return 1 * time.Second
}

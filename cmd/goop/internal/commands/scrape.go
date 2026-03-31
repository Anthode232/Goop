package commands

import (
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
	selector    string
	output      string
	format      string
	jsRendering bool
	waitFor     string
	workers     int
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape [url]",
	Short: "Scrape a single URL for data",
	Long: `Scrape a single URL and extract data using CSS selectors.

Examples:
  goop scrape https://example.com --selector ".title"
  goop scrape https://example.com --selector ".title" --output results.json
  goop scrape https://spa-example.com --js --wait-for ".content" --selector ".title"`,
	Args: cobra.ExactArgs(1),
	RunE: runScrape,
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.Flags().StringVarP(&selector, "selector", "s", "", "CSS selector to extract data")
	scrapeCmd.Flags().StringVarP(&output, "output", "o", "", "output file path")
	scrapeCmd.Flags().StringVarP(&format, "format", "f", "json", "output format (json, csv, xml, text)")
	scrapeCmd.Flags().BoolVarP(&jsRendering, "js", "j", false, "enable JavaScript rendering")
	scrapeCmd.Flags().StringVar(&waitFor, "wait-for", "", "CSS selector to wait for (JS mode only)")
	scrapeCmd.Flags().IntVarP(&workers, "workers", "w", 5, "number of concurrent workers (for multiple selectors)")

	// Bind flags to viper
	viper.BindPFlag("selector", scrapeCmd.Flags().Lookup("selector"))
	viper.BindPFlag("output", scrapeCmd.Flags().Lookup("output"))
	viper.BindPFlag("format", scrapeCmd.Flags().Lookup("format"))
	viper.BindPFlag("js", scrapeCmd.Flags().Lookup("js"))
	viper.BindPFlag("wait-for", scrapeCmd.Flags().Lookup("wait-for"))
	viper.BindPFlag("workers", scrapeCmd.Flags().Lookup("workers"))
}

func runScrape(cmd *cobra.Command, args []string) error {
	url := args[0]

	// Configure global settings
	configureGlobalSettings()

	var result goop.Root
	var err error

	if jsRendering {
		// JavaScript rendering mode
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
		// Static HTML mode
		result = goop.Scrape(url, getTimeout())
		err = result.Error
	}

	if err != nil {
		return fmt.Errorf("failed to scrape %s: %v", url, err)
	}

	// Extract data if selector provided
	var outputData interface{}
	if selector != "" {
		elements := result.CSS(selector)
		if elements.Error != nil {
			return fmt.Errorf("failed to find elements with selector %s: %v", selector, elements.Error)
		}

		// Extract text content from all matching elements
		text := elements.Text()
		outputData = []string{text}
	} else {
		// Return full HTML if no selector
		outputData = result.FullText()
	}

	// Output results
	if output != "" {
		return writeToFile(outputData, output, format)
	}

	// Print to stdout
	return printOutput(outputData, format)
}

func configureGlobalSettings() {
	// Set debug level
	if viper.GetBool("debug") {
		goop.SetDebugLevel(goop.DebugVerbose)
	} else if viper.GetBool("verbose") {
		goop.SetDebugLevel(goop.DebugBasic)
	}

	// Set user agent
	if userAgent := viper.GetString("user-agent"); userAgent != "" {
		goop.Header("User-Agent", userAgent)
	}

	// Set rate limit
	if rateLimit := viper.GetString("rate-limit"); rateLimit != "" {
		if duration, err := time.ParseDuration(rateLimit); err == nil {
			goop.SetRateLimit(duration)
		}
	}

	// Set timeout
	if timeout := viper.GetString("timeout"); timeout != "" {
		if duration, err := time.ParseDuration(timeout); err == nil {
			goop.SetTimeout(duration)
		}
	}
}

func getTimeout() time.Duration {
	if timeout := viper.GetString("timeout"); timeout != "" {
		if duration, err := time.ParseDuration(timeout); err == nil {
			return duration
		}
	}
	return 30 * time.Second
}

func writeToFile(data interface{}, filename, format string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	var content string
	switch strings.ToLower(format) {
	case "json":
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		content = string(jsonData)
	case "csv":
		if texts, ok := data.([]string); ok {
			content = strings.Join(texts, "\n")
		} else {
			return fmt.Errorf("CSV format only supported for text arrays")
		}
	case "xml":
		// Simple XML format for text arrays
		if texts, ok := data.([]string); ok {
			content = "<results>\n"
			for _, text := range texts {
				content += fmt.Sprintf("  <item>%s</item>\n", text)
			}
			content += "</results>"
		} else {
			content = fmt.Sprintf("<result>%s</result>", data)
		}
	case "text":
		if texts, ok := data.([]string); ok {
			content = strings.Join(texts, "\n")
		} else {
			content = fmt.Sprintf("%v", data)
		}
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}

	_, err = file.WriteString(content)
	return err
}

func printOutput(data interface{}, format string) error {
	switch strings.ToLower(format) {
	case "json":
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		fmt.Println(string(jsonData))
	case "csv", "text":
		if texts, ok := data.([]string); ok {
			for _, text := range texts {
				fmt.Println(text)
			}
		} else {
			fmt.Printf("%v\n", data)
		}
	case "xml":
		if texts, ok := data.([]string); ok {
			fmt.Println("<results>")
			for _, text := range texts {
				fmt.Printf("  <item>%s</item>\n", text)
			}
			fmt.Println("</results>")
		} else {
			fmt.Printf("<result>%v</result>\n", data)
		}
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
	return nil
}

package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	setConfig   []string
	showConfig  bool
	resetConfig bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long: `Manage Goop configuration settings.

Examples:
  goop config --show                    # Show current configuration
  goop config --set timeout=60s         # Set timeout to 60 seconds
  goop config --set user-agent="Custom Agent"  # Set custom user agent
  goop config --reset                   # Reset to defaults`,
	RunE: runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringSliceVar(&setConfig, "set", []string{}, "set configuration value (key=value)")
	configCmd.Flags().BoolVar(&showConfig, "show", false, "show current configuration")
	configCmd.Flags().BoolVar(&resetConfig, "reset", false, "reset configuration to defaults")
}

func runConfig(cmd *cobra.Command, args []string) error {
	if resetConfig {
		return resetConfiguration()
	}

	if showConfig {
		return showConfiguration()
	}

	if len(setConfig) > 0 {
		return setConfigurationValues(setConfig)
	}

	return showConfiguration()
}

func resetConfiguration() error {
	// Remove config file if it exists
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	configFile := fmt.Sprintf("%s/.goop.yaml", home)
	if _, err := os.Stat(configFile); err == nil {
		if err := os.Remove(configFile); err != nil {
			return fmt.Errorf("failed to remove config file: %v", err)
		}
		fmt.Println("Configuration reset to defaults")
	} else {
		fmt.Println("No configuration file found")
	}

	return nil
}

func showConfiguration() error {
	fmt.Println("Current Goop Configuration:")
	fmt.Println("==========================")

	// Show all configuration values
	settings := map[string]interface{}{
		"verbose":     viper.GetBool("verbose"),
		"debug":       viper.GetBool("debug"),
		"user-agent":  viper.GetString("user-agent"),
		"timeout":     viper.GetString("timeout"),
		"rate-limit":  viper.GetString("rate-limit"),
		"retry":       viper.GetInt("retry"),
		"config-file": viper.ConfigFileUsed(),
	}

	for key, value := range settings {
		fmt.Printf("%-12s: %v\n", key, value)
	}

	// Show default values for reference
	fmt.Println("\nDefault Values:")
	fmt.Println("===============")
	defaults := map[string]interface{}{
		"verbose":     false,
		"debug":       false,
		"user-agent":  "",
		"timeout":     "30s",
		"rate-limit":  "1s",
		"retry":       3,
	}

	for key, value := range defaults {
		fmt.Printf("%-12s: %v\n", key, value)
	}

	return nil
}

func setConfigurationValues(values []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	configFile := fmt.Sprintf("%s/.goop.yaml", home)

	// Set configuration values
	for _, value := range values {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid configuration format: %s (expected key=value)", value)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Validate key
		validKeys := map[string]bool{
			"verbose":    true,
			"debug":      true,
			"user-agent": true,
			"timeout":    true,
			"rate-limit": true,
			"retry":      true,
		}

		if !validKeys[key] {
			return fmt.Errorf("unknown configuration key: %s", key)
		}

		// Validate and set value
		switch key {
		case "verbose", "debug":
			if val != "true" && val != "false" {
				return fmt.Errorf("value for %s must be true or false", key)
			}
			viper.Set(key, val == "true")
		case "retry":
			var retry int
			if _, err := fmt.Sscanf(val, "%d", &retry); err != nil || retry < 0 {
				return fmt.Errorf("value for %s must be a non-negative integer", key)
			}
			viper.Set(key, retry)
		default:
			viper.Set(key, val)
		}

		fmt.Printf("Set %s = %s\n", key, val)
	}

	// Write configuration to file
	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write configuration: %v", err)
	}

	fmt.Printf("Configuration saved to: %s\n", configFile)
	return nil
}

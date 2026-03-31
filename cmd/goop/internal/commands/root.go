package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goop",
	Short: "Goop is a powerful web scraping tool",
	Long: `Goop is a professional web scraping library and CLI tool that supports:
- Static HTML parsing with CSS selectors
- JavaScript rendering for dynamic sites  
- Concurrent scraping for performance
- Multiple export formats (JSON, CSV, XML)
- Rate limiting and retry logic
- Professional debugging capabilities`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goop.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug output")
	rootCmd.PersistentFlags().StringP("user-agent", "u", "", "custom user agent")
	rootCmd.PersistentFlags().StringP("timeout", "t", "30s", "request timeout")
	rootCmd.PersistentFlags().StringP("rate-limit", "r", "1s", "rate limit between requests")
	rootCmd.PersistentFlags().IntP("retry", "R", 3, "number of retry attempts")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("user-agent", rootCmd.PersistentFlags().Lookup("user-agent"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("rate-limit", rootCmd.PersistentFlags().Lookup("rate-limit"))
	viper.BindPFlag("retry", rootCmd.PersistentFlags().Lookup("retry"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".goop" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".goop")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

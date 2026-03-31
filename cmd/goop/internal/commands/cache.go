package commands

import (
	"fmt"
	"os"
	"strconv"
	"time"

	goop "github.com/ez0000001000000/Goop"
	"github.com/spf13/cobra"
)

// CacheCmd represents the cache command
var CacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage caching configuration",
	Long:  `Configure and manage Goop's caching system.`,
}

func init() {
	// Add cache subcommands
	CacheCmd.AddCommand(cacheStatusCmd)
	CacheCmd.AddCommand(cacheClearCmd)
	CacheCmd.AddCommand(cacheConfigCmd)
}

// cacheStatusCmd shows cache statistics
var cacheStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cache statistics",
	Long:  `Display current cache performance statistics and configuration.`,
	Run:   cacheStatus,
}

// cacheClearCmd clears the cache
var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all caches",
	Long:  `Empty both memory and disk caches.`,
	Run:   cacheClear,
}

// cacheConfigCmd configures cache settings
var cacheConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure cache settings",
	Long:  `Set cache configuration options.`,
	Run:   cacheConfig,
}

func cacheStatus(cmd *cobra.Command, args []string) {
	stats := goop.GetCacheStats()
	config := goop.GetCacheConfig()

	fmt.Println("=== Cache Configuration ===")
	fmt.Printf("Enabled: %v\n", config.Enabled)
	fmt.Printf("Memory Limit: %d MB\n", config.MemoryLimit/1024/1024)
	fmt.Printf("Disk Limit: %d MB\n", config.DiskLimit/1024/1024)
	fmt.Printf("Default TTL: %v\n", config.DefaultTTL)
	fmt.Printf("Cache Directory: %s\n", config.CacheDir)

	fmt.Println("\n=== Cache Statistics ===")
	fmt.Printf("Memory Hits: %d\n", stats.MemoryHits)
	fmt.Printf("Memory Misses: %d\n", stats.MemoryMisses)
	fmt.Printf("Disk Hits: %d\n", stats.DiskHits)
	fmt.Printf("Disk Misses: %d\n", stats.DiskMisses)
	fmt.Printf("Memory Usage: %.2f MB\n", float64(stats.MemorySize)/1024/1024)
	fmt.Printf("Disk Usage: %.2f MB\n", float64(stats.DiskSize)/1024/1024)
	fmt.Printf("Total Entries: %d\n", stats.TotalEntries)

	// Calculate hit rates
	totalRequests := stats.MemoryHits + stats.MemoryMisses + stats.DiskHits + stats.DiskMisses
	if totalRequests > 0 {
		memoryHitRate := float64(stats.MemoryHits) / float64(stats.MemoryHits+stats.MemoryMisses) * 100
		diskHitRate := float64(stats.DiskHits) / float64(stats.DiskHits+stats.DiskMisses) * 100
		totalHitRate := float64(stats.MemoryHits+stats.DiskHits) / float64(totalRequests) * 100

		fmt.Printf("Memory Hit Rate: %.1f%%\n", memoryHitRate)
		fmt.Printf("Disk Hit Rate: %.1f%%\n", diskHitRate)
		fmt.Printf("Overall Hit Rate: %.1f%%\n", totalHitRate)
	}
}

func cacheClear(cmd *cobra.Command, args []string) {
	fmt.Println("Clearing cache...")
	err := goop.ClearCache()
	if err != nil {
		fmt.Printf("Error clearing cache: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Cache cleared successfully!")
}

func cacheConfig(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: goop cache config <key> <value>")
		fmt.Println("\nAvailable configuration options:")
		fmt.Println("  enabled <true|false>     Enable/disable caching")
		fmt.Println("  memory-limit <MB>        Set memory limit in MB")
		fmt.Println("  disk-limit <MB>          Set disk limit in MB")
		fmt.Println("  ttl <duration>           Set default TTL (e.g., 1h, 30m)")
		fmt.Println("  cache-dir <path>         Set cache directory")
		return
	}

	key := args[0]
	value := args[1]

	config := goop.GetCacheConfig()

	switch key {
	case "enabled":
		if value == "true" || value == "false" {
			config.Enabled = value == "true"
		} else {
			fmt.Println("Error: enabled must be 'true' or 'false'")
			os.Exit(1)
		}
	case "memory-limit":
		if mb, err := strconv.ParseInt(value, 10, 64); err == nil {
			config.MemoryLimit = mb * 1024 * 1024
		} else {
			fmt.Println("Error: memory-limit must be a number in MB")
			os.Exit(1)
		}
	case "disk-limit":
		if mb, err := strconv.ParseInt(value, 10, 64); err == nil {
			config.DiskLimit = mb * 1024 * 1024
		} else {
			fmt.Println("Error: disk-limit must be a number in MB")
			os.Exit(1)
		}
	case "ttl":
		if duration, err := time.ParseDuration(value); err == nil {
			config.DefaultTTL = duration
		} else {
			fmt.Println("Error: ttl must be a valid duration (e.g., 1h, 30m, 10s)")
			os.Exit(1)
		}
	case "cache-dir":
		config.CacheDir = value
	default:
		fmt.Printf("Error: Unknown configuration key '%s'\n", key)
		os.Exit(1)
	}

	// Apply new configuration
	goop.SetCacheConfig(config)
	fmt.Printf("Cache configuration updated: %s = %s\n", key, value)
}

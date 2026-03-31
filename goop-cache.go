package goop

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// CacheEntry represents a cached response
type CacheEntry struct {
	Key       string        `json:"key"`
	Content   string        `json:"content"`
	Headers   string        `json:"headers"`
	Timestamp time.Time     `json:"timestamp"`
	TTL       time.Duration `json:"ttl"`
	Size      int64         `json:"size"`
}

// IsExpired checks if the cache entry has expired
func (c *CacheEntry) IsExpired() bool {
	return time.Since(c.Timestamp) > c.TTL
}

// CacheConfig defines caching behavior
type CacheConfig struct {
	Enabled     bool          `json:"enabled"`
	MemoryLimit int64         `json:"memory_limit"` // bytes
	DiskLimit   int64         `json:"disk_limit"`   // bytes
	DefaultTTL  time.Duration `json:"default_ttl"`
	CacheDir    string        `json:"cache_dir"`
	Compression bool          `json:"compression"`
}

// DefaultCacheConfig returns sensible default caching settings
var DefaultCacheConfig = CacheConfig{
	Enabled:     true,
	MemoryLimit: 100 * 1024 * 1024, // 100MB
	DiskLimit:   500 * 1024 * 1024, // 500MB
	DefaultTTL:  1 * time.Hour,
	CacheDir:    ".goop_cache",
	Compression: true,
}

// FastCacheConfig optimized for speed
var FastCacheConfig = CacheConfig{
	Enabled:     true,
	MemoryLimit: 200 * 1024 * 1024,  // 200MB
	DiskLimit:   1024 * 1024 * 1024, // 1GB
	DefaultTTL:  30 * time.Minute,
	CacheDir:    ".goop_cache_fast",
	Compression: false,
}

// CacheStats provides cache performance metrics
type CacheStats struct {
	MemoryHits   int64 `json:"memory_hits"`
	MemoryMisses int64 `json:"memory_misses"`
	DiskHits     int64 `json:"disk_hits"`
	DiskMisses   int64 `json:"disk_misses"`
	MemorySize   int64 `json:"memory_size"`
	DiskSize     int64 `json:"disk_size"`
	TotalEntries int64 `json:"total_entries"`
}

// Cache interface defines cache operations
type Cache interface {
	Get(key string) (*CacheEntry, bool)
	Set(key string, entry *CacheEntry) error
	Delete(key string) error
	Clear() error
	Stats() CacheStats
}

// MemoryCache provides in-memory caching
type MemoryCache struct {
	config CacheConfig
	cache  map[string]*CacheEntry
	mutex  sync.RWMutex
	stats  CacheStats
}

// NewMemoryCache creates a new memory cache
func NewMemoryCache(config CacheConfig) *MemoryCache {
	return &MemoryCache{
		config: config,
		cache:  make(map[string]*CacheEntry),
		stats:  CacheStats{},
	}
}

// Get retrieves an entry from memory cache
func (m *MemoryCache) Get(key string) (*CacheEntry, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entry, exists := m.cache[key]
	if !exists {
		m.stats.MemoryMisses++
		return nil, false
	}

	if entry.IsExpired() {
		delete(m.cache, key)
		m.stats.MemoryMisses++
		return nil, false
	}

	m.stats.MemoryHits++
	return entry, true
}

// Set stores an entry in memory cache
func (m *MemoryCache) Set(key string, entry *CacheEntry) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check memory limit
	if m.getCurrentMemorySize()+entry.Size > m.config.MemoryLimit {
		m.evictOldest()
	}

	m.cache[key] = entry
	m.stats.MemorySize += entry.Size
	m.stats.TotalEntries = int64(len(m.cache))
	return nil
}

// Delete removes an entry from memory cache
func (m *MemoryCache) Delete(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if entry, exists := m.cache[key]; exists {
		delete(m.cache, key)
		m.stats.MemorySize -= entry.Size
		m.stats.TotalEntries = int64(len(m.cache))
	}
	return nil
}

// Clear empties the memory cache
func (m *MemoryCache) Clear() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.cache = make(map[string]*CacheEntry)
	m.stats.MemorySize = 0
	m.stats.TotalEntries = 0
	return nil
}

// Stats returns cache statistics
func (m *MemoryCache) Stats() CacheStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.stats
}

// getCurrentMemorySize calculates current memory usage
func (m *MemoryCache) getCurrentMemorySize() int64 {
	var size int64
	for _, entry := range m.cache {
		size += entry.Size
	}
	return size
}

// evictOldest removes the oldest entry
func (m *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range m.cache {
		if oldestKey == "" || entry.Timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.Timestamp
		}
	}

	if oldestKey != "" {
		if entry, exists := m.cache[oldestKey]; exists {
			delete(m.cache, oldestKey)
			m.stats.MemorySize -= entry.Size
		}
	}
}

// DiskCache provides persistent disk-based caching
type DiskCache struct {
	config CacheConfig
	mutex  sync.RWMutex
	stats  CacheStats
}

// NewDiskCache creates a new disk cache
func NewDiskCache(config CacheConfig) *DiskCache {
	// Ensure cache directory exists
	os.MkdirAll(config.CacheDir, 0755)

	return &DiskCache{
		config: config,
		stats:  CacheStats{},
	}
}

// Get retrieves an entry from disk cache
func (d *DiskCache) Get(key string) (*CacheEntry, bool) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	filename := filepath.Join(d.config.CacheDir, key+".json")

	data, err := os.ReadFile(filename)
	if err != nil {
		d.stats.DiskMisses++
		return nil, false
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		d.stats.DiskMisses++
		return nil, false
	}

	if entry.IsExpired() {
		os.Remove(filename)
		d.stats.DiskMisses++
		return nil, false
	}

	d.stats.DiskHits++
	return &entry, true
}

// Set stores an entry in disk cache
func (d *DiskCache) Set(key string, entry *CacheEntry) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check disk limit
	if d.getCurrentDiskSize()+entry.Size > d.config.DiskLimit {
		d.evictOldestDisk()
	}

	filename := filepath.Join(d.config.CacheDir, key+".json")
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	d.stats.DiskSize += entry.Size
	return nil
}

// Delete removes an entry from disk cache
func (d *DiskCache) Delete(key string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	filename := filepath.Join(d.config.CacheDir, key+".json")

	// Get file size before deletion
	if info, err := os.Stat(filename); err == nil {
		d.stats.DiskSize -= info.Size()
	}

	return os.Remove(filename)
}

// Clear empties the disk cache
func (d *DiskCache) Clear() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	err := os.RemoveAll(d.config.CacheDir)
	if err != nil {
		return err
	}

	os.MkdirAll(d.config.CacheDir, 0755)
	d.stats.DiskSize = 0
	return nil
}

// Stats returns cache statistics
func (d *DiskCache) Stats() CacheStats {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.stats
}

// getCurrentDiskSize calculates current disk usage
func (d *DiskCache) getCurrentDiskSize() int64 {
	var size int64
	filepath.Walk(d.config.CacheDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

// evictOldestDisk removes the oldest file from disk
func (d *DiskCache) evictOldestDisk() {
	var oldestFile string
	var oldestTime time.Time

	filepath.Walk(d.config.CacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		if oldestFile == "" || info.ModTime().Before(oldestTime) {
			oldestFile = path
			oldestTime = info.ModTime()
		}
		return nil
	})

	if oldestFile != "" {
		if info, err := os.Stat(oldestFile); err == nil {
			os.Remove(oldestFile)
			d.stats.DiskSize -= info.Size()
		}
	}
}

// HybridCache combines memory and disk caching
type HybridCache struct {
	memoryCache *MemoryCache
	diskCache   *DiskCache
	config      CacheConfig
	mutex       sync.RWMutex
}

// NewHybridCache creates a new hybrid cache
func NewHybridCache(config CacheConfig) *HybridCache {
	return &HybridCache{
		memoryCache: NewMemoryCache(config),
		diskCache:   NewDiskCache(config),
		config:      config,
	}
}

// Get retrieves an entry from hybrid cache (memory first, then disk)
func (h *HybridCache) Get(key string) (*CacheEntry, bool) {
	// Try memory cache first
	if entry, found := h.memoryCache.Get(key); found {
		return entry, true
	}

	// Try disk cache
	if entry, found := h.diskCache.Get(key); found {
		// Promote to memory cache
		h.memoryCache.Set(key, entry)
		return entry, true
	}

	return nil, false
}

// Set stores an entry in both memory and disk cache
func (h *HybridCache) Set(key string, entry *CacheEntry) error {
	// Store in both caches
	if err := h.memoryCache.Set(key, entry); err != nil {
		return err
	}

	return h.diskCache.Set(key, entry)
}

// Delete removes an entry from both caches
func (h *HybridCache) Delete(key string) error {
	h.memoryCache.Delete(key)
	return h.diskCache.Delete(key)
}

// Clear empties both caches
func (h *HybridCache) Clear() error {
	h.memoryCache.Clear()
	return h.diskCache.Clear()
}

// Stats returns combined cache statistics
func (h *HybridCache) Stats() CacheStats {
	memStats := h.memoryCache.Stats()
	diskStats := h.diskCache.Stats()

	return CacheStats{
		MemoryHits:   memStats.MemoryHits,
		MemoryMisses: memStats.MemoryMisses,
		DiskHits:     diskStats.DiskHits,
		DiskMisses:   diskStats.DiskMisses,
		MemorySize:   memStats.MemorySize,
		DiskSize:     diskStats.DiskSize,
		TotalEntries: memStats.TotalEntries,
	}
}

// Global cache instance
var globalCache Cache
var cacheConfig CacheConfig
var cacheMutex sync.RWMutex

// SetCacheConfig configures the global cache
func SetCacheConfig(config CacheConfig) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	cacheConfig = config

	if config.Enabled {
		globalCache = NewHybridCache(config)
	} else {
		globalCache = nil
	}
}

// GetCacheConfig returns current cache configuration
func GetCacheConfig() CacheConfig {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	return cacheConfig
}

// GetCacheStats returns cache performance statistics
func GetCacheStats() CacheStats {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	if globalCache == nil {
		return CacheStats{}
	}

	return globalCache.Stats()
}

// ClearCache empties all caches
func ClearCache() error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if globalCache == nil {
		return nil
	}

	return globalCache.Clear()
}

// generateCacheKey creates a unique cache key
func generateCacheKey(url string, method string, headers map[string]string) string {
	h := md5.New()
	h.Write([]byte(url))
	h.Write([]byte(method))

	// Add headers to key for cache differentiation
	for k, v := range headers {
		h.Write([]byte(k + ":" + v))
	}

	return hex.EncodeToString(h.Sum(nil))
}

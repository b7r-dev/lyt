package build

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cache stores file modification times for incremental builds
type Cache struct {
	mu      sync.RWMutex
	FileMod map[string]int64 `json:"file_mod"`
	Path    string           `json:"-"`
}

// NewCache creates or loads a build cache
func NewCache(cachePath string) *Cache {
	c := &Cache{
		FileMod: make(map[string]int64),
		Path:    cachePath,
	}

	// Try to load existing cache
	data, err := os.ReadFile(cachePath)
	if err == nil {
		if err := json.Unmarshal(data, c); err != nil {
			// Ignore cache load errors, start fresh
		}
	}
	return c
}

// Save writes the cache to disk
func (c *Cache) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.Path, data, 0644)
}

// GetModTime returns the cached modification time, or 0 if not found
func (c *Cache) GetModTime(path string) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.FileMod[path]
}

// SetModTime updates the modification time for a path
func (c *Cache) SetModTime(path string, modTime int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.FileMod[path] = modTime
}

// HasChanged checks if a file has changed since last build
func (c *Cache) HasChanged(path string) bool {
	c.mu.RLock()
	oldTime, exists := c.FileMod[path]
	c.mu.RUnlock()

	if !exists {
		return true // New file
	}

	// Get current modification time
	info, err := os.Stat(path)
	if err != nil {
		return true // File doesn't exist or can't be stat
	}

	return info.ModTime().Unix() > oldTime
}

// Touch marks a file as processed in the cache
func (c *Cache) Touch(path string) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	c.SetModTime(path, info.ModTime().Unix())
}

// TouchDir processes all files in a directory and updates the cache
func (c *Cache) TouchDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			c.Touch(path)
		}
		return nil
	})
}

// ShouldRebuild checks if any file in the given directories has changed
func (c *Cache) ShouldRebuild(dirs ...string) bool {
	for _, dir := range dirs {
		changed := false
		_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip on error
			}
			if !info.IsDir() && c.HasChanged(path) {
				changed = true
				return filepath.SkipDir // Exit walk
			}
			return nil
		})
		if changed {
			return true
		}
	}
	return false
}

// Invalidate removes all entries from the cache
func (c *Cache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.FileMod = make(map[string]int64)
}

// CacheDuration returns how long since the last build
func (c *Cache) CacheDuration() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.FileMod) == 0 {
		return 0
	}

	// Return a reasonable default - in production we'd track build time
	return time.Hour
}

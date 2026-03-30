package build

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCache(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "cache.json")

	c := NewCache(cachePath)

	// Test that new cache has no entries
	if len(c.FileMod) != 0 {
		t.Errorf("expected empty cache, got %d entries", len(c.FileMod))
	}

	// Test SetModTime and GetModTime
	testPath := "/test/file.txt"
	c.SetModTime(testPath, 12345)
	if got := c.GetModTime(testPath); got != 12345 {
		t.Errorf("expected 12345, got %d", got)
	}

	// Test HasChanged for new file
	if !c.HasChanged("/nonexistent") {
		t.Error("expected new file to have changed")
	}

	// Save and reload
	if err := c.Save(); err != nil {
		t.Fatal(err)
	}
	c2 := NewCache(cachePath)
	if c2.GetModTime(testPath) != 12345 {
		t.Error("cache not persisted correctly")
	}
}

func TestCacheShouldRebuild(t *testing.T) {
	tmpDir := t.TempDir()
	contentDir := filepath.Join(tmpDir, "content")
	if err := os.Mkdir(contentDir, 0755); err != nil {
		t.Fatal(err)
	}

	cachePath := filepath.Join(tmpDir, "cache.json")

	// Create a file
	testFile := filepath.Join(contentDir, "test.yaml")
	if err := os.WriteFile(testFile, []byte("test: true"), 0644); err != nil {
		t.Fatal(err)
	}

	// New cache - should rebuild (cache is empty)
	c := NewCache(cachePath)
	if !c.ShouldRebuild(contentDir) {
		t.Error("expected rebuild when cache is empty")
	}

	// Save the cache with the file
	c.Touch(testFile)
	if err := c.Save(); err != nil {
		t.Fatal(err)
	}

	// Reload - no changes, should not rebuild
	c2 := NewCache(cachePath)
	if c2.ShouldRebuild(contentDir) {
		t.Error("expected no rebuild when files haven't changed")
	}
}

func TestCacheInvalidate(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "cache.json")

	c := NewCache(cachePath)
	c.SetModTime("/test/file.txt", 12345)
	c.Invalidate()

	if len(c.FileMod) != 0 {
		t.Errorf("expected empty cache after invalidate, got %d entries", len(c.FileMod))
	}
}

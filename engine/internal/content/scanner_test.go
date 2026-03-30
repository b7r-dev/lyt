package content

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanner(t *testing.T) {
	// Create a temporary content directory structure
	tmpDir := t.TempDir()

	// Create pages directory
	pagesDir := filepath.Join(tmpDir, "pages")
	if err := os.Mkdir(pagesDir, 0755); err != nil {
		t.Fatalf("failed to create pages dir: %v", err)
	}

	// Create a test page
	pageContent := `
meta:
  title: "Test Page"
  slug: "/test"
  description: "A test page"
sections:
  - id: "hero"
    type: "hero"
    title: "Welcome"
`
	if err := os.WriteFile(filepath.Join(pagesDir, "test.yaml"), []byte(pageContent), 0644); err != nil {
		t.Fatalf("failed to write test page: %v", err)
	}

	// Create blog directory
	blogDir := filepath.Join(tmpDir, "blog")
	if err := os.Mkdir(blogDir, 0755); err != nil {
		t.Fatalf("failed to create blog dir: %v", err)
	}

	// Create a test blog post
	postContent := `
meta:
  title: "Test Post"
  slug: "test-post"
  description: "A test post"
  date: "2026-03-29"
  published: true
body: |
  ## Hello
  This is a test post.
`
	if err := os.WriteFile(filepath.Join(blogDir, "test-post.yaml"), []byte(postContent), 0644); err != nil {
		t.Fatalf("failed to write test post: %v", err)
	}

	// Scan the content
	scanner := NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Verify pages
	if len(collection.Pages) != 1 {
		t.Errorf("expected 1 page, got %d", len(collection.Pages))
	}

	// Verify blog posts
	if len(collection.Blog) != 1 {
		t.Errorf("expected 1 blog post, got %d", len(collection.Blog))
	}

	// Check page content
	if len(collection.Pages) > 0 {
		page := collection.Pages[0]
		meta, ok := page.Data["meta"].(map[string]interface{})
		if !ok {
			t.Error("expected meta to be a map")
		}
		if meta["title"] != "Test Page" {
			t.Errorf("expected title 'Test Page', got %v", meta["title"])
		}
	}
}

func TestScannerSkipsDataDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a data directory (should be skipped)
	dataDir := filepath.Join(tmpDir, "data")
	if err := os.Mkdir(dataDir, 0755); err != nil {
		t.Fatal(err)
	}

	dataContent := `
items:
  - name: "Test"
`
	if err := os.WriteFile(filepath.Join(dataDir, "items.yaml"), []byte(dataContent), 0644); err != nil {
		t.Fatal(err)
	}

	scanner := NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Data directory should be skipped - no pages
	if len(collection.Pages) != 0 {
		t.Errorf("expected 0 pages (data dir skipped), got %d", len(collection.Pages))
	}
}

func TestContentFileTypes(t *testing.T) {
	// Test that content file type detection works
	// Based on file location in the collection
	tmpDir := t.TempDir()
	pagesDir := filepath.Join(tmpDir, "pages")
	blogDir := filepath.Join(tmpDir, "blog")
	os.Mkdir(pagesDir, 0755)
	os.Mkdir(blogDir, 0755)

	// Create test files
	os.WriteFile(filepath.Join(pagesDir, "test.yaml"), []byte("meta:\n  title: test"), 0644)
	os.WriteFile(filepath.Join(blogDir, "test.yaml"), []byte("meta:\n  title: test"), 0644)

	scanner := NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(collection.Pages) != 1 {
		t.Errorf("expected 1 page, got %d", len(collection.Pages))
	}
	if len(collection.Blog) != 1 {
		t.Errorf("expected 1 blog post, got %d", len(collection.Blog))
	}
}

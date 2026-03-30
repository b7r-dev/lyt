package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/b7r-dev/lyt/engine/internal/content"
)

// TestRenderNavActiveState tests that nav links have correct aria-current
func TestRenderNavActiveState(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config with nav
	configDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "site.yaml"), []byte(`
meta:
  title: "Test Site"
nav:
  - label: "Home"
    href: "/"
  - label: "Docs"
    href: "/docs"
  - label: "About"
    href: "/about"
`), 0644)

	// Create pages
	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	os.WriteFile(filepath.Join(pagesDir, "index.yaml"), []byte(`
meta:
  title: "Home"
  slug: "/"
sections:
  - id: "hero"
    type: "hero"
    title: "Welcome"
`), 0644)

	os.WriteFile(filepath.Join(pagesDir, "docs.yaml"), []byte(`
meta:
  title: "Docs"
  slug: "/docs"
sections:
  - id: "intro"
    title: "Documentation"
`), 0644)

	scanner := content.NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	renderer := NewRenderer(collection, tmpDir, false)

	tests := []struct {
		name        string
		pageRelPath string
		currentSlug string
		wantActive  string // href that should be active
	}{
		{
			name:        "home page active on /",
			pageRelPath: "pages/index.yaml",
			currentSlug: "/",
			wantActive:  "/",
		},
		{
			name:        "docs page active on /docs",
			pageRelPath: "pages/docs.yaml",
			currentSlug: "/docs",
			wantActive:  "/docs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cf content.ContentFile
			for _, p := range collection.Pages {
				if p.RelPath == tt.pageRelPath {
					cf = p
					break
				}
			}

			html, err := renderer.RenderPage(cf)
			if err != nil {
				t.Fatalf("RenderPage failed: %v", err)
			}

			// Should have aria-current on the active link
			// Note: Both desktop and mobile nav render, so we expect 2 (one of each)
			expectedAttr := `aria-current="page"`
			if !strings.Contains(html, expectedAttr) {
				t.Errorf("Expected aria-current attribute, HTML: %s", html)
			}

			// Should have 3 aria-current (desktop nav + mobile nav + home brand link on home page)
			// Or 2 on other pages (desktop + mobile nav)
			count := strings.Count(html, `aria-current="page"`)
			if count < 2 || count > 3 {
				t.Errorf("Expected 2-3 aria-current, got %d", count)
			}
		})
	}
}

// TestRenderHeroSection tests hero section rendering
func TestRenderHeroSection(t *testing.T) {
	tmpDir := t.TempDir()

	configDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "site.yaml"), []byte(`
meta:
  title: "Test Site"
`), 0644)

	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	os.WriteFile(filepath.Join(pagesDir, "test.yaml"), []byte(`
meta:
  title: "Test"
  slug: "/test"
sections:
  - id: "hero"
    type: "hero"
    title: "Hello World"
    subtitle: "A subtitle"
    body: "Welcome to the site"
`), 0644)

	scanner := content.NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	renderer := NewRenderer(collection, tmpDir, false)

	html, err := renderer.RenderPage(collection.Pages[0])
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	// Check hero elements
	if !strings.Contains(html, `<h1 class="hero-title">Hello World</h1>`) {
		t.Error("Hero title not found")
	}
	if !strings.Contains(html, `<p class="hero-subtitle">A subtitle</p>`) {
		t.Error("Hero subtitle not found")
	}
	if !strings.Contains(html, `<div class="hero-body">`) {
		t.Error("Hero body not found")
	}
}

// TestRenderFeaturesSection tests features/cards section rendering
func TestRenderFeaturesSection(t *testing.T) {
	tmpDir := t.TempDir()

	configDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "site.yaml"), []byte(`
meta:
  title: "Test Site"
`), 0644)

	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	os.WriteFile(filepath.Join(pagesDir, "test.yaml"), []byte(`
meta:
  title: "Test"
  slug: "/test"
sections:
  - id: "features"
    type: "features"
    title: "Features"
    cards:
      - title: "Feature One"
        body: "Description one"
      - title: "Feature Two"
        body: "Description two"
`), 0644)

	scanner := content.NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	renderer := NewRenderer(collection, tmpDir, false)

	html, err := renderer.RenderPage(collection.Pages[0])
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	// Check features elements
	if !strings.Contains(html, `<h2 class="section-title">Features</h2>`) {
		t.Error("Features title not found")
	}
	if !strings.Contains(html, `<div class="features-grid">`) {
		t.Error("Features grid not found")
	}
	if !strings.Contains(html, `<h3 class="feature-title">Feature One</h3>`) {
		t.Error("First feature title not found")
	}
	if !strings.Contains(html, `<h3 class="feature-title">Feature Two</h3>`) {
		t.Error("Second feature title not found")
	}
}

// TestRenderCalloutSection tests callout section rendering
func TestRenderCalloutSection(t *testing.T) {
	tmpDir := t.TempDir()

	configDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "site.yaml"), []byte(`
meta:
  title: "Test Site"
`), 0644)

	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	os.WriteFile(filepath.Join(pagesDir, "test.yaml"), []byte(`
meta:
  title: "Test"
  slug: "/test"
sections:
  - id: "note"
    type: "callout"
    variant: "tip"
    title: "Note"
    body: "This is a tip"
`), 0644)

	scanner := content.NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	renderer := NewRenderer(collection, tmpDir, false)

	html, err := renderer.RenderPage(collection.Pages[0])
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	// Check callout
	if !strings.Contains(html, `<aside class="callout callout-tip"`) {
		t.Error("Callout with variant not found")
	}
	if !strings.Contains(html, `<strong class="callout-title">Note</strong>`) {
		t.Error("Callout title not found")
	}
}

// TestRenderCTASection tests CTA section rendering
func TestRenderCTASection(t *testing.T) {
	tmpDir := t.TempDir()

	configDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "site.yaml"), []byte(`
meta:
  title: "Test Site"
`), 0644)

	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	os.WriteFile(filepath.Join(pagesDir, "test.yaml"), []byte(`
meta:
  title: "Test"
  slug: "/test"
sections:
  - id: "cta"
    type: "cta"
    title: "Get Started"
    body: "Start building today"
    button_text: "Click Here"
    button_href: "/signup"
`), 0644)

	scanner := content.NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	renderer := NewRenderer(collection, tmpDir, false)

	html, err := renderer.RenderPage(collection.Pages[0])
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	// Check CTA
	if !strings.Contains(html, `<h2 class="cta-title">Get Started</h2>`) {
		t.Error("CTA title not found")
	}
	if !strings.Contains(html, `<a href="/signup" class="button button-primary">Click Here</a>`) {
		t.Error("CTA button not found")
	}
}

// TestRenderDivider tests divider section rendering
func TestRenderDivider(t *testing.T) {
	tmpDir := t.TempDir()

	configDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "site.yaml"), []byte(`
meta:
  title: "Test Site"
`), 0644)

	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	os.WriteFile(filepath.Join(pagesDir, "test.yaml"), []byte(`
meta:
  title: "Test"
  slug: "/test"
sections:
  - id: "break"
    type: "divider"
    icon: "***"
`), 0644)

	scanner := content.NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	renderer := NewRenderer(collection, tmpDir, false)

	html, err := renderer.RenderPage(collection.Pages[0])
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	// Check divider
	if !strings.Contains(html, `<div class="section-break" id="break">`) {
		t.Error("Divider not found")
	}
	if !strings.Contains(html, `<span class="section-break-icon">***</span>`) {
		t.Error("Divider icon not found")
	}
}

// TestRenderPullQuote tests pull-quote section rendering
func TestRenderPullQuote(t *testing.T) {
	tmpDir := t.TempDir()

	configDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "site.yaml"), []byte(`
meta:
  title: "Test Site"
`), 0644)

	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	os.WriteFile(filepath.Join(pagesDir, "test.yaml"), []byte(`
meta:
  title: "Test"
  slug: "/test"
sections:
  - type: "pull-quote"
    quote: "This is a quote"
    attribution: "John Doe"
`), 0644)

	scanner := content.NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	renderer := NewRenderer(collection, tmpDir, false)

	html, err := renderer.RenderPage(collection.Pages[0])
	if err != nil {
		t.Fatalf("RenderPage failed: %v", err)
	}

	// Check pull quote
	if !strings.Contains(html, `<blockquote class="pull-quote">`) {
		t.Error("Pull quote not found")
	}
	if !strings.Contains(html, `<span class="attribution">John Doe</span>`) {
		t.Error("Attribution not found")
	}
}

// TestSlugToClass tests URL slug to CSS class conversion
func TestSlugToClass(t *testing.T) {
	tests := []struct {
		slug string
		want string
	}{
		{"/", ""},
		{"/docs", "docs"},
		{"/docs/getting-started", "docs-getting-started"},
		{"docs/", "docs"},
		{"docs", "docs"},
	}

	for _, tt := range tests {
		t.Run(tt.slug, func(t *testing.T) {
			result := slugToClass(tt.slug)
			if result != tt.want {
				t.Errorf("slugToClass(%q) = %q, want %q", tt.slug, result, tt.want)
			}
		})
	}
}

// TestRenderBlogPost tests blog post rendering
func TestRenderBlogPost(t *testing.T) {
	tmpDir := t.TempDir()

	configDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "site.yaml"), []byte(`
meta:
  title: "Test Site"
`), 0644)

	blogDir := filepath.Join(tmpDir, "blog")
	os.MkdirAll(blogDir, 0755)
	os.WriteFile(filepath.Join(blogDir, "test-post.yaml"), []byte(`
meta:
  title: "Test Post"
  slug: "test-post"
  description: "A test post"
  date: "2026-03-29"
sections:
  - type: "default"
    body: "Hello world"
`), 0644)

	scanner := content.NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	renderer := NewRenderer(collection, tmpDir, false)

	if len(collection.Blog) != 1 {
		t.Fatalf("expected 1 blog post, got %d", len(collection.Blog))
	}

	html, err := renderer.RenderBlogPost(collection.Blog[0])
	if err != nil {
		t.Fatalf("RenderBlogPost failed: %v", err)
	}

	// Check blog post elements
	if !strings.Contains(html, `<article class="blog-post">`) {
		t.Error("Blog post article not found")
	}
	if !strings.Contains(html, `<h1 class="post-title">Test Post</h1>`) {
		t.Error("Post title not found")
	}
	if !strings.Contains(html, `<p class="post-meta">2026-03-29</p>`) {
		t.Error("Post date not found")
	}
	if !strings.Contains(html, `<p class="post-description">A test post</p>`) {
		t.Error("Post description not found")
	}
}

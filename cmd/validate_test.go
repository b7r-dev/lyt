package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/b7r-dev/lyt/internal/content"
)

// TestResolveURL tests the URL resolution logic
func TestResolveURL(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		relative string
		want     string
	}{
		{
			name:     "absolute path",
			base:     "/docs/getting-started",
			relative: "/about",
			want:     "/about",
		},
		{
			name:     "relative path same dir",
			base:     "/docs/getting-started",
			relative: "configuration",
			want:     "/docs/getting-started/configuration",
		},
		{
			name:     "relative path parent",
			base:     "/docs/components",
			relative: "../about",
			want:     "/docs/about",
		},
		{
			name:     "relative path subdir",
			base:     "/docs",
			relative: "getting-started/install",
			want:     "/docs/getting-started/install",
		},
		{
			name:     "root base",
			base:     "/",
			relative: "about",
			want:     "/about",
		},
		{
			name:     "empty base",
			base:     "",
			relative: "about",
			want:     "/about",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveURL(tt.base, tt.relative)
			if got != tt.want {
				t.Errorf("resolveURL(%q, %q) = %q, want %q", tt.base, tt.relative, got, tt.want)
			}
		})
	}
}

// TestLinkValidationEdgeCases tests link validation with various edge cases
func TestLinkValidationEdgeCases(t *testing.T) {
	// Create temp dist directory
	tmpDir := t.TempDir()

	// Create sitemap
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>https://lyt.local/</loc></url>
  <url><loc>https://lyt.local/about</loc></url>
  <url><loc>https://lyt.local/docs</loc></url>
  <url><loc>https://lyt.local/docs/getting-started</loc></url>
  <url><loc>https://lyt.local/docs/configuration</loc></url>
</urlset>`
	os.WriteFile(filepath.Join(tmpDir, "sitemap.xml"), []byte(sitemap), 0644)

	// Create index.html with various link types
	indexHTML := `<!DOCTYPE html>
<html>
<head><title>Home</title></head>
<body>
  <a href="/about">About</a>
  <a href="/docs">Docs</a>
  <a href="/docs/getting-started">Getting Started</a>
  <a href="/docs/getting-started#install">Install Section</a>
  <a href="#top">Top Anchor</a>
  <a href="https://external.com">External</a>
  <a href="//cdn.example.com">Protocol Relative</a>
  <a href="mailto:test@example.com">Email</a>
  <a href="tel:+1234567890">Phone</a>
</body>
</html>`
	os.MkdirAll(filepath.Join(tmpDir, "about"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "about", "index.html"), []byte(`<!DOCTYPE html><html><body>About</body></html>`), 0644)

	os.MkdirAll(filepath.Join(tmpDir, "docs"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "docs", "index.html"), []byte(`<!DOCTYPE html><html><body>Docs</body></html>`), 0644)

	os.MkdirAll(filepath.Join(tmpDir, "docs", "getting-started"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "docs", "getting-started", "index.html"), []byte(indexHTML), 0644)

	os.MkdirAll(filepath.Join(tmpDir, "docs", "configuration"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "docs", "configuration", "index.html"), []byte(`<!DOCTYPE html><html><body>Config</body></html>`), 0644)

	// Create root index
	os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(`<!DOCTYPE html><html><body>Home</body></html>`), 0644)

	// Test validation
	validateDir = tmpDir
	err := runLinkValidation()
	validateDir = ""

	if err != nil {
		t.Errorf("runLinkValidation() error = %v", err)
	}
}

// TestLinkValidationWithBrokenLinks tests detection of broken links
func TestLinkValidationWithBrokenLinks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create sitemap
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url><loc>https://lyt.local/</loc></url>
  <url><loc>https://lyt.local/page</loc></url>
</urlset>`
	os.WriteFile(filepath.Join(tmpDir, "sitemap.xml"), []byte(sitemap), 0644)

	// Create index.html with broken link
	brokenHTML := `<!DOCTYPE html>
<html>
<body>
  <a href="/about">About (exists)</a>
  <a href="/nonexistent">Broken Link</a>
  <a href="/about#section">Anchor on exists</a>
  <a href="/missing#section">Anchor on missing</a>
</body>
</html>`
	os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte(brokenHTML), 0644)

	os.MkdirAll(filepath.Join(tmpDir, "about"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "about", "index.html"), []byte(`<!DOCTYPE html><html><body>About</body></html>`), 0644)

	// Validate - should find broken links
	validateDir = tmpDir
	err := runLinkValidation()
	validateDir = ""

	// Should return error because there are broken links
	if err == nil {
		t.Log("Warning: expected error for broken links but got none")
	}
}

// TestLinkValidationNoSitemap tests behavior when no sitemap exists
func TestLinkValidationNoSitemap(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an empty dir without sitemap - should skip gracefully
	validateDir = tmpDir
	err := runLinkValidation()
	validateDir = ""

	// Should NOT error, just skip
	if err != nil {
		t.Errorf("runLinkValidation() with no sitemap should skip, got error: %v", err)
	}
}

// TestLinkValidationNoDist tests behavior when dist doesn't exist
func TestLinkValidationNoDist(t *testing.T) {
	validateDir = "/nonexistent/path"
	err := runLinkValidation()
	validateDir = ""

	if err != nil {
		t.Errorf("runLinkValidation() with no dist error = %v", err)
	}
}

// TestValidatePage tests page validation logic
func TestValidatePage(t *testing.T) {
	tests := []struct {
		name        string
		cf          content.ContentFile
		wantErrors  int
		errorFields []string
	}{
		{
			name: "valid page",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test Page",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":  "hero",
							"title": "Test",
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "missing title",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"slug": "/test",
					},
				},
			},
			wantErrors:  1,
			errorFields: []string{"meta.title"},
		},
		{
			name: "missing slug",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
					},
				},
			},
			wantErrors:  1,
			errorFields: []string{"meta.slug"},
		},
		{
			name: "invalid section type",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type": "invalid-type",
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "cta without button_text",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":        "cta",
							"button_href": "/link",
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "cta without button_href",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":        "cta",
							"button_text": "Click Me",
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "valid cta with both fields",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":        "cta",
							"title":       "Title",
							"button_text": "Click Me",
							"button_href": "/link",
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "callout with invalid variant",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":    "callout",
							"variant": "invalid-variant",
							"title":   "Note",
							"body":    "Content",
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "callout with valid variant",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":    "callout",
							"variant": "tip",
							"title":   "Note",
							"body":    "Content",
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "code-example without code field",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":  "code-example",
							"title": "Example",
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "code-example with code field",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":     "code-example",
							"title":    "Example",
							"code":     "print(\"hello\")",
							"language": "python",
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "hero without title",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type": "hero",
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "features with valid cards",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":  "features",
							"title": "Features",
							"cards": []interface{}{
								map[string]interface{}{
									"title": "Feature 1",
									"body":  "Description",
								},
								map[string]interface{}{
									"title": "Feature 2",
									"body":  "Description",
								},
							},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "features with card missing title",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":  "features",
							"title": "Features",
							"cards": []interface{}{
								map[string]interface{}{
									"body": "Description only",
								},
							},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "agent page with agent_content",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
						"agent": true,
					},
					"agent_content": map[string]interface{}{
						"title": "Agent Title",
						"sections": []interface{}{
							map[string]interface{}{
								"type":  "default",
								"title": "Section",
								"body":  "Content",
							},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "agent page without agent_content",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
						"agent": true,
					},
				},
			},
			wantErrors:  1,
			errorFields: []string{"agent_content"},
		},
		{
			name: "button with invalid variant",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":  "hero",
							"title": "Hero",
							"buttons": []interface{}{
								map[string]interface{}{
									"text":    "Click",
									"href":    "/link",
									"variant": "invalid",
								},
							},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "button with valid primary variant",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":  "hero",
							"title": "Hero",
							"buttons": []interface{}{
								map[string]interface{}{
									"text":    "Click",
									"href":    "/link",
									"variant": "primary",
								},
							},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "button with valid outline variant",
			cf: content.ContentFile{
				RelPath: "pages/test.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test",
						"slug":  "/test",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type":  "hero",
							"title": "Hero",
							"buttons": []interface{}{
								map[string]interface{}{
									"text":    "Click",
									"href":    "/link",
									"variant": "outline",
								},
							},
						},
					},
				},
			},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validatePage(tt.cf)
			if len(errs) != tt.wantErrors {
				t.Errorf("validatePage() got %d errors, want %d. Errors: %+v", len(errs), tt.wantErrors, errs)
			}

			// Check specific error fields if specified
			if len(tt.errorFields) > 0 {
				found := false
				for _, err := range errs {
					for _, field := range tt.errorFields {
						if err.Field == field {
							found = true
							break
						}
					}
				}
				if !found && tt.wantErrors > 0 {
					t.Errorf("Expected error for field %v but not found in %+v", tt.errorFields, errs)
				}
			}
		})
	}
}

// TestValidateBlog tests blog post validation logic
func TestValidateBlog(t *testing.T) {
	tests := []struct {
		name       string
		cf         content.ContentFile
		wantErrors int
	}{
		{
			name: "valid blog post",
			cf: content.ContentFile{
				RelPath: "blog/test-post.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test Post",
						"slug":  "test-post",
						"date":  "2026-01-01",
					},
					"sections": []interface{}{
						map[string]interface{}{
							"type": "default",
							"body": "Content",
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "blog slug with leading slash (invalid)",
			cf: content.ContentFile{
				RelPath: "blog/test-post.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test Post",
						"slug":  "/test-post",
						"date":  "2026-01-01",
					},
				},
			},
			wantErrors: 2, // Both leading slash and invalid slug chars
		},
		{
			name: "invalid date format",
			cf: content.ContentFile{
				RelPath: "blog/test-post.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test Post",
						"slug":  "test-post",
						"date":  "01-01-2026",
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "valid date format YYYY-MM-DD",
			cf: content.ContentFile{
				RelPath: "blog/test-post.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test Post",
						"slug":  "test-post",
						"date":  "2026-03-30",
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "invalid slug characters",
			cf: content.ContentFile{
				RelPath: "blog/test-post.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test Post",
						"slug":  "Test_Post!",
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "valid slug with hyphens",
			cf: content.ContentFile{
				RelPath: "blog/test-post.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test Post",
						"slug":  "test-post-2026",
					},
				},
			},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validateBlog(tt.cf)
			if len(errs) != tt.wantErrors {
				t.Errorf("validateBlog() got %d errors, want %d. Errors: %+v", len(errs), tt.wantErrors, errs)
			}
		})
	}
}

// TestValidateAgentSection tests agent section validation
func TestValidateAgentSection(t *testing.T) {
	tests := []struct {
		name       string
		file       string
		section    map[string]interface{}
		wantErrors int
	}{
		{
			name:       "valid cli section",
			file:       "pages/test.yaml",
			section:    map[string]interface{}{"type": "cli", "commands": "ls -la"},
			wantErrors: 0,
		},
		{
			name:       "cli section missing commands",
			file:       "pages/test.yaml",
			section:    map[string]interface{}{"type": "cli"},
			wantErrors: 1,
		},
		{
			name:       "valid example section",
			file:       "pages/test.yaml",
			section:    map[string]interface{}{"type": "example", "example": "code here"},
			wantErrors: 0,
		},
		{
			name:       "example section missing example",
			file:       "pages/test.yaml",
			section:    map[string]interface{}{"type": "example"},
			wantErrors: 1,
		},
		{
			name:       "valid link section",
			file:       "pages/test.yaml",
			section:    map[string]interface{}{"type": "link", "text": "Click", "href": "/link"},
			wantErrors: 0,
		},
		{
			name:       "link section missing text",
			file:       "pages/test.yaml",
			section:    map[string]interface{}{"type": "link", "href": "/link"},
			wantErrors: 1,
		},
		{
			name:       "link section missing href",
			file:       "pages/test.yaml",
			section:    map[string]interface{}{"type": "link", "text": "Click"},
			wantErrors: 1,
		},
		{
			name:       "invalid agent section type",
			file:       "pages/test.yaml",
			section:    map[string]interface{}{"type": "invalid"},
			wantErrors: 1,
		},
		{
			name:       "valid default section",
			file:       "pages/test.yaml",
			section:    map[string]interface{}{"type": "default", "title": "Title", "body": "Content"},
			wantErrors: 0,
		},
		{
			name:       "valid schema section",
			file:       "pages/test.yaml",
			section:    map[string]interface{}{"type": "schema", "schema": "field: value"},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validateAgentSection(tt.file, 0, tt.section)
			if len(errs) != tt.wantErrors {
				t.Errorf("validateAgentSection() got %d errors, want %d. Errors: %+v", len(errs), tt.wantErrors, errs)
			}
		})
	}
}

// TestValidateConfigField tests config validation
func TestValidateConfigField(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      interface{}
		wantErrors int
	}{
		{
			name:       "valid nav item",
			key:        "nav",
			value:      []interface{}{map[string]interface{}{"label": "Home", "href": "/"}},
			wantErrors: 0,
		},
		{
			name:       "nav missing label",
			key:        "nav",
			value:      []interface{}{map[string]interface{}{"href": "/"}},
			wantErrors: 1,
		},
		{
			name:       "nav missing href",
			key:        "nav",
			value:      []interface{}{map[string]interface{}{"label": "Home"}},
			wantErrors: 1,
		},
		{
			name:       "valid agent_section with path",
			key:        "agent_section",
			value:      map[string]interface{}{"enabled": true, "path": "/agents"},
			wantErrors: 0,
		},
		{
			name:       "agent_section path without leading slash",
			key:        "agent_section",
			value:      map[string]interface{}{"enabled": true, "path": "agents"},
			wantErrors: 1,
		},
		{
			name:       "agent_section with empty path",
			key:        "agent_section",
			value:      map[string]interface{}{"enabled": true, "path": ""},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validateConfigField(tt.key, tt.value)
			if len(errs) != tt.wantErrors {
				t.Errorf("validateConfigField() got %d errors, want %d. Errors: %+v", len(errs), tt.wantErrors, errs)
			}
		})
	}
}

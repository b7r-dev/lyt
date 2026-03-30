package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/b7r-dev/lyt/engine/internal/content"
)

// TestAgentLinks verifies that agent links appear on pages correctly
func TestAgentLinks(t *testing.T) {
	// Create a temp directory for test content
	tmpDir := t.TempDir()

	// Create minimal config
	configDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "site.yaml"), []byte(`
meta:
  title: "Test Site"
agent_section:
  path: "/agents"
`), 0644)

	// Create a page with agent: true and agent_content
	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	os.WriteFile(filepath.Join(pagesDir, "index.yaml"), []byte(`
meta:
  title: "Home"
  slug: "/"
  agent: true

sections:
  - id: "hero"
    type: "hero"
    title: "Welcome"
`), 0644)

	// Create a page with agent: true but NO agent_content
	os.WriteFile(filepath.Join(pagesDir, "about.yaml"), []byte(`
meta:
  title: "About"
  slug: "/about"
  agent: true

sections:
  - id: "intro"
    title: "About Us"
`), 0644)

	// Create a page without agent flag
	os.WriteFile(filepath.Join(pagesDir, "contact.yaml"), []byte(`
meta:
  title: "Contact"
  slug: "/contact"

sections:
  - id: "intro"
    title: "Get in Touch"
`), 0644)

	// Create a page with agent: true AND agent_content (should link to its own agent page)
	os.WriteFile(filepath.Join(pagesDir, "docs.yaml"), []byte(`
meta:
  title: "Docs"
  slug: "/docs"
  agent: true

agent_content:
  title: "Documentation"
  sections:
    - type: "default"
      body: "Agent docs here"

sections:
  - id: "intro"
    title: "Documentation"
`), 0644)

	// Scan content
	scanner := content.NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Test HasAgentPage
	tests := []struct {
		name        string
		relPath     string
		wantAgent   bool
		wantHubLink bool
		wantLink    string
	}{
		{
			name:        "page with agent_content",
			relPath:     "pages/docs.yaml",
			wantAgent:   true,
			wantHubLink: true,
			wantLink:    "/agents/docs",
		},
		{
			name:        "page with agent true but no agent_content",
			relPath:     "pages/about.yaml",
			wantAgent:   false,
			wantHubLink: true,
			wantLink:    "/agents",
		},
		{
			name:        "page without agent flag",
			relPath:     "pages/contact.yaml",
			wantAgent:   false,
			wantHubLink: false,
			wantLink:    "",
		},
		{
			name:        "home page with agent true",
			relPath:     "pages/index.yaml",
			wantAgent:   false,
			wantHubLink: true,
			wantLink:    "/agents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cf *content.ContentFile
			for i := range collection.Pages {
				if collection.Pages[i].RelPath == tt.relPath {
					cf = &collection.Pages[i]
					break
				}
			}
			if cf == nil {
				t.Fatalf("Page not found: %s", tt.relPath)
			}

			gotAgent := collection.HasAgentPage(*cf)
			if gotAgent != tt.wantAgent {
				t.Errorf("HasAgentPage() = %v, want %v", gotAgent, tt.wantAgent)
			}

			gotHubLink := collection.ShowAgentHubLink(*cf)
			if gotHubLink != tt.wantHubLink {
				t.Errorf("ShowAgentHubLink() = %v, want %v", gotHubLink, tt.wantHubLink)
			}
		})
	}

	// Now test actual HTML rendering
	renderer := NewRenderer(collection, tmpDir, false)

	for _, tt := range tests {
		t.Run(tt.name+"_render", func(t *testing.T) {
			var cf content.ContentFile
			for i := range collection.Pages {
				if collection.Pages[i].RelPath == tt.relPath {
					cf = collection.Pages[i]
					break
				}
			}

			html, err := renderer.RenderPage(cf)
			if err != nil {
				t.Fatalf("RenderPage failed: %v", err)
			}

			hasAgentLink := strings.Contains(html, `class="agent-link"`)
			hasAgentLinkToExpected := strings.Contains(html, `href="`+tt.wantLink+`"`)

			if tt.wantHubLink && !hasAgentLink {
				t.Errorf("Expected agent link in HTML, but not found")
			}

			if tt.wantHubLink && !hasAgentLinkToExpected {
				t.Errorf("Expected agent link to %q, but not found. HTML: %s", tt.wantLink, html)
			}

			if !tt.wantHubLink && hasAgentLink {
				t.Errorf("Did not expect agent link, but found one")
			}
		})
	}
}

// TestAgentPageGeneration verifies agent pages are generated correctly
func TestAgentPageGeneration(t *testing.T) {
	// Create a temp directory for test content
	tmpDir := t.TempDir()

	// Create minimal config
	configDir := filepath.Join(tmpDir, "config")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "site.yaml"), []byte(`
meta:
  title: "Test Site"
agent_section:
  path: "/agents"
`), 0644)

	// Create a page with agent_content
	pagesDir := filepath.Join(tmpDir, "pages")
	os.MkdirAll(pagesDir, 0755)
	os.WriteFile(filepath.Join(pagesDir, "docs.yaml"), []byte(`
meta:
  title: "Docs"
  slug: "/docs"
  agent: true

agent_content:
  title: "Documentation for Agents"
  description: "Agent-specific content"
  sections:
    - type: "cli"
      title: "Commands"
      commands: |
        lyt build

sections:
  - id: "intro"
    title: "Docs"
`), 0644)

	// Scan content
	scanner := content.NewScanner(tmpDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	renderer := NewRenderer(collection, tmpDir, false)

	// Render agent page
	var cf content.ContentFile
	for i := range collection.Pages {
		if collection.Pages[i].RelPath == "pages/docs.yaml" {
			cf = collection.Pages[i]
			break
		}
	}

	agentHTML, err := renderer.RenderAgentPage(cf)
	if err != nil {
		t.Fatalf("RenderAgentPage failed: %v", err)
	}

	// Verify agent content is rendered (not human content)
	if !strings.Contains(agentHTML, "Documentation for Agents") {
		t.Error("Agent page should contain agent_content title")
	}

	if !strings.Contains(agentHTML, "lyt build") {
		t.Error("Agent page should contain CLI commands")
	}

	// Should NOT contain human content
	if strings.Contains(agentHTML, "Docs") && strings.Contains(agentHTML, "section-hero") {
		t.Error("Agent page should NOT contain human hero section")
	}
}

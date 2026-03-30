package content

import (
	"testing"
)

// TestGetMeta tests the GetMeta helper function
func TestGetMeta(t *testing.T) {
	tests := []struct {
		name      string
		cf        ContentFile
		wantTitle string
		wantSlug  string
	}{
		{
			name: "page with meta",
			cf: ContentFile{
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Test Title",
						"slug":  "/test",
					},
				},
			},
			wantTitle: "Test Title",
			wantSlug:  "/test",
		},
		{
			name: "page without meta",
			cf: ContentFile{
				Data: map[string]interface{}{},
			},
			wantTitle: "",
			wantSlug:  "",
		},
		{
			name: "page with empty meta",
			cf: ContentFile{
				Data: map[string]interface{}{
					"meta": map[string]interface{}{},
				},
			},
			wantTitle: "",
			wantSlug:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := GetMeta(tt.cf)

			gotTitle := ""
			if v, ok := meta["title"].(string); ok {
				gotTitle = v
			}
			if gotTitle != tt.wantTitle {
				t.Errorf("title = %q, want %q", gotTitle, tt.wantTitle)
			}
			gotSlug := ""
			if v, ok := meta["slug"].(string); ok {
				gotSlug = v
			}
			if gotSlug != tt.wantSlug {
				t.Errorf("slug = %q, want %q", gotSlug, tt.wantSlug)
			}
		})
	}
}

// TestGetSection tests the GetSection helper function
func TestGetSection(t *testing.T) {
	cf := ContentFile{
		Data: map[string]interface{}{
			"sections": []interface{}{
				map[string]interface{}{
					"id":    "hero",
					"title": "Welcome",
				},
				map[string]interface{}{
					"id":    "about",
					"title": "About Us",
				},
			},
		},
	}

	// Test finding existing section
	hero := GetSection(cf, "hero")
	if hero == nil {
		t.Error("expected to find hero section")
	}
	if hero["title"] != "Welcome" {
		t.Errorf("hero title = %q, want %q", hero["title"], "Welcome")
	}

	// Test finding non-existent section
	footer := GetSection(cf, "footer")
	if footer != nil {
		t.Error("expected nil for non-existent section")
	}
}

// TestToStrings tests the ToStrings helper function
func TestToStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name:     "string slice",
			input:    []interface{}{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "mixed slice",
			input:    []interface{}{"a", 1, "b"},
			expected: []string{"a", "1", "b"},
		},
		{
			name:     "non-slice",
			input:    "single",
			expected: nil,
		},
		{
			name:     "nil",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToStrings(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("ToStrings() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCollectionMethods tests collection helper methods
func TestCollectionMethods(t *testing.T) {
	// Create a minimal collection for testing
	c := &Collection{
		Config: map[string]interface{}{
			"meta": map[string]interface{}{
				"title": "Test Site",
			},
			"nav": []interface{}{
				map[string]interface{}{
					"label": "Home",
					"href":  "/",
				},
			},
			"agent_section": map[string]interface{}{
				"enabled": true,
				"path":    "/agents",
			},
		},
		Pages: []ContentFile{
			{
				RelPath: "pages/index.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "Home",
						"slug":  "/",
						"agent": true,
					},
					"agent_content": map[string]interface{}{
						"title": "Agent Home",
					},
				},
			},
			{
				RelPath: "pages/about.yaml",
				Data: map[string]interface{}{
					"meta": map[string]interface{}{
						"title": "About",
						"slug":  "/about",
					},
				},
			},
		},
	}

	// Test GetSiteTitle
	if c.GetSiteTitle() != "Test Site" {
		t.Errorf("GetSiteTitle() = %q, want %q", c.GetSiteTitle(), "Test Site")
	}

	// Test GetNav
	nav := c.GetNav()
	if len(nav) != 1 {
		t.Errorf("GetNav() length = %d, want 1", len(nav))
	}

	// Test AgentSectionEnabled
	if !c.AgentSectionEnabled() {
		t.Error("AgentSectionEnabled() = false, want true")
	}

	// Test GetAgentPath
	if c.GetAgentPath() != "/agents" {
		t.Errorf("GetAgentPath() = %q, want %q", c.GetAgentPath(), "/agents")
	}

	// Test HasAgentPage - page with agent_content
	if !c.HasAgentPage(c.Pages[0]) {
		t.Error("HasAgentPage(index) = false, want true (has agent_content)")
	}

	// Test HasAgentPage - page without agent flag
	if c.HasAgentPage(c.Pages[1]) {
		t.Error("HasAgentPage(about) = true, want false (no agent flag)")
	}

	// Test ShowAgentHubLink - page with agent: true
	if !c.ShowAgentHubLink(c.Pages[0]) {
		t.Error("ShowAgentHubLink(index) = false, want true")
	}

	// Test ShowAgentHubLink - page without agent flag
	if c.ShowAgentHubLink(c.Pages[1]) {
		t.Error("ShowAgentHubLink(about) = true, want false")
	}

	// Test GetAgentContent
	ac := GetAgentContent(c.Pages[0])
	if ac == nil {
		t.Error("GetAgentContent(index) = nil, want agent_content")
	}
	if ac["title"] != "Agent Home" {
		t.Errorf("agent_content title = %q, want %q", ac["title"], "Agent Home")
	}

	// Test GetAgentContent - no agent_content
	ac2 := GetAgentContent(c.Pages[1])
	if ac2 != nil {
		t.Error("GetAgentContent(about) = want nil")
	}
}

// TestCollectionDefaults tests default values
func TestCollectionDefaults(t *testing.T) {
	c := &Collection{
		Config: map[string]interface{}{},
	}

	// Test defaults when no config
	if c.GetSiteTitle() != "lyt" {
		t.Errorf("default GetSiteTitle() = %q, want %q", c.GetSiteTitle(), "lyt")
	}

	if c.GetAgentPath() != "/agents" {
		t.Errorf("default GetAgentPath() = %q, want %q", c.GetAgentPath(), "/agents")
	}

	if c.AgentSectionEnabled() {
		t.Error("default AgentSectionEnabled() = true, want false")
	}

	if len(c.GetNav()) != 0 {
		t.Errorf("default GetNav() = %v, want empty", c.GetNav())
	}
}

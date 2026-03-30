package cmd

import "testing"

func TestPageSlugToPath(t *testing.T) {
	tests := []struct {
		name    string
		relPath string
		want    string
	}{
		{
			name:    "root index",
			relPath: "pages/index.yaml",
			want:    "dist/index.html",
		},
		{
			name:    "about page",
			relPath: "pages/about.yaml",
			want:    "dist/about/index.html",
		},
		{
			name:    "docs index",
			relPath: "pages/docs.yaml",
			want:    "dist/docs/index.html",
		},
		{
			name:    "docs getting-started - multi-segment",
			relPath: "pages/docs/getting-started.yaml",
			want:    "dist/docs/getting-started/index.html",
		},
		{
			name:    "docs deployment - multi-segment",
			relPath: "pages/docs/deployment.yaml",
			want:    "dist/docs/deployment/index.html",
		},
		{
			name:    "components page",
			relPath: "pages/components.yaml",
			want:    "dist/components/index.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pageSlugToPath(tt.relPath, "dist")
			if got != tt.want {
				t.Errorf("pageSlugToPath(%q) = %q, want %q", tt.relPath, got, tt.want)
			}
		})
	}
}

func TestBlogSlugToPath(t *testing.T) {
	tests := []struct {
		name    string
		relPath string
		want    string
	}{
		{
			name:    "blog post",
			relPath: "blog/community-over-capital.yaml",
			want:    "dist/blog/community-over-capital/index.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := blogSlugToPath(tt.relPath, "dist")
			if got != tt.want {
				t.Errorf("blogSlugToPath(%q) = %q, want %q", tt.relPath, got, tt.want)
			}
		})
	}
}

func TestAgentSlugToPath(t *testing.T) {
	tests := []struct {
		name            string
		relPath         string
		agentPathPrefix string
		want            string
	}{
		{
			name:            "agent hub",
			relPath:         "pages/agents.yaml",
			agentPathPrefix: "agents",
			want:            "dist/agents/index.html",
		},
		{
			name:            "docs page under agents",
			relPath:         "pages/docs/getting-started.yaml",
			agentPathPrefix: "agents",
			want:            "dist/agents/docs/getting-started/index.html",
		},
		{
			name:            "docs configuration under agents",
			relPath:         "pages/docs/configuration.yaml",
			agentPathPrefix: "agents",
			want:            "dist/agents/docs/configuration/index.html",
		},
		{
			name:            "custom agent path prefix",
			relPath:         "pages/docs/deployment.yaml",
			agentPathPrefix: "/ai",
			want:            "dist/ai/docs/deployment/index.html",
		},
		{
			name:            "blog post under agents (flattened)",
			relPath:         "blog/test-post.yaml",
			agentPathPrefix: "agents",
			want:            "dist/agents/test-post/index.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := agentSlugToPath(tt.relPath, "dist", tt.agentPathPrefix)
			if got != tt.want {
				t.Errorf("agentSlugToPath(%q, %q) = %q, want %q", tt.relPath, tt.agentPathPrefix, got, tt.want)
			}
		})
	}
}

func TestGetSlug(t *testing.T) {
	tests := []struct {
		name string
		data map[string]interface{}
		want string
	}{
		{
			name: "basic slug",
			data: map[string]interface{}{
				"meta": map[string]interface{}{
					"slug": "/about",
				},
			},
			want: "/about",
		},
		{
			name: "docs slug",
			data: map[string]interface{}{
				"meta": map[string]interface{}{
					"slug": "/docs/getting-started",
				},
			},
			want: "/docs/getting-started",
		},
		{
			name: "no meta",
			data: map[string]interface{}{},
			want: "",
		},
		{
			name: "no slug field",
			data: map[string]interface{}{
				"meta": map[string]interface{}{},
			},
			want: "",
		},
		{
			name: "slug not a string",
			data: map[string]interface{}{
				"meta": map[string]interface{}{
					"slug": 123,
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSlug(tt.data)
			if got != tt.want {
				t.Errorf("getSlug() = %q, want %q", got, tt.want)
			}
		})
	}
}

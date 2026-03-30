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

package markdown

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "heading",
			input:    "# Hello World",
			expected: "<h1",
		},
		{
			name:     "paragraph",
			input:    "This is a paragraph.",
			expected: "<p>This is a paragraph.</p>",
		},
		{
			name:     "bold",
			input:    "This is **bold** text",
			expected: "<strong>bold</strong>",
		},
		{
			name:     "italic",
			input:    "This is *italic* text",
			expected: "<em>italic</em>",
		},
		{
			name:     "code",
			input:    "`inline code`",
			expected: "<code>inline code</code>",
		},
		{
			name:     "link",
			input:    "[link text](https://example.com)",
			expected: "<a href=\"https://example.com\">link text</a>",
		},
		{
			name:     "list",
			input:    "- item 1\n- item 2",
			expected: "<ul>",
		},
		{
			name:     "gfm strikethrough",
			input:    "~~deleted~~",
			expected: "<del>deleted</del>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := Render(tt.input)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Render(%q) = %q, want to contain %q", tt.input, output, tt.expected)
			}
		})
	}
}

func TestRenderEmpty(t *testing.T) {
	output, err := Render("")
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if output != "" {
		t.Errorf("Render(\"\") = %q, want empty string", output)
	}
}

func TestRenderHeadingIDs(t *testing.T) {
	output, err := Render("## Test Heading")
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	// Goldmark adds IDs to headings
	if !strings.Contains(output, "id=") {
		t.Errorf("expected heading ID, got: %s", output)
	}
}

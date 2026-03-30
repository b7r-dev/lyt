package tokens

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessTokens(t *testing.T) {
	// Create a temporary YAML file with test tokens
	content := `
colors:
  base:
    bg: "#f5f5f0"
    text: "#3d3d3d"
typography:
  font_family:
    body: "Georgia, serif"
  font_size:
    sm: "0.875rem"
z:
  base: "0"
  top: "100"
`
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "tokens.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	output, err := ProcessTokens(tmpFile, false)
	if err != nil {
		t.Fatalf("ProcessTokens failed: %v", err)
	}

	// Verify output contains expected CSS custom properties
	tests := []string{
		"--color-base-bg",
		"--color-base-text",
		"--font-body",
		"--text-sm",
		"--z-base",
		"--z-top",
	}

	for _, expected := range tests {
		if !contains(output, expected) {
			t.Errorf("expected output to contain %s, got:\n%s", expected, output)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestProcessTokensFileNotFound(t *testing.T) {
	_, err := ProcessTokens("/nonexistent/path/tokens.yaml", false)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]interface{}{
		"z": "2",
		"a": "1",
		"m": "3",
	}

	keys := sortedKeys(m)
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}

	// Keys should be sorted alphabetically
	if keys[0] != "a" || keys[1] != "m" || keys[2] != "z" {
		t.Errorf("expected keys to be sorted, got %v", keys)
	}
}

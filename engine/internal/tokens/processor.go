package tokens

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func ProcessTokens(path string, verbose bool) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("/* lyt design tokens */\n")
	sb.WriteString(":root {\n")

	// Colors
	if colors, ok := raw["colors"].(map[string]interface{}); ok {
		sb.WriteString("  /* Colors */\n")
		writeColors(&sb, colors, "")
	}

	// Spacing
	if spacing, ok := raw["spacing"].(map[string]interface{}); ok {
		sb.WriteString("  /* Spacing */\n")
		writeNumbers(&sb, spacing, "--space")
	}

	// Typography
	if typ, ok := raw["typography"].(map[string]interface{}); ok {
		sb.WriteString("  /* Typography */\n")
		writeTypography(&sb, typ)
	}

	// Z-indices (the 3 planes)
	if z, ok := raw["z"].(map[string]interface{}); ok {
		sb.WriteString("  /* Z Planes */\n")
		writeNumbers(&sb, z, "--z")
	}

	sb.WriteString("}\n")
	return sb.String(), nil
}

func writeColors(sb *strings.Builder, m map[string]interface{}, prefix string) {
	keys := sortedKeys(m)
	for _, k := range keys {
		v := m[k]
		name := k
		if prefix != "" {
			name = prefix + "-" + k
		}
		switch x := v.(type) {
		case string:
			sb.WriteString(fmt.Sprintf("  --color-%s: %s;\n", name, x))
		case map[string]interface{}:
			writeColors(sb, x, name)
		}
	}
}

func writeNumbers(sb *strings.Builder, m map[string]interface{}, prefix string) {
	keys := sortedKeys(m)
	for _, k := range keys {
		v := m[k]
		switch x := v.(type) {
		case string:
			sb.WriteString(fmt.Sprintf("  %s-%s: %s;\n", prefix, k, x))
		case int:
			sb.WriteString(fmt.Sprintf("  %s-%s: %dpx;\n", prefix, k, x))
		case float64:
			sb.WriteString(fmt.Sprintf("  %s-%s: %.2fpx;\n", prefix, k, x))
		}
	}
}

func writeTypography(sb *strings.Builder, m map[string]interface{}) {
	if ff, ok := m["font_family"].(map[string]interface{}); ok {
		for k, v := range ff {
			sb.WriteString(fmt.Sprintf("  --font-%s: %s;\n", k, v))
		}
	}
	if fs, ok := m["font_size"].(map[string]interface{}); ok {
		writeNumbers(sb, fs, "--text")
	}
	if fw, ok := m["font_weight"].(map[string]interface{}); ok {
		writeNumbers(sb, fw, "--weight")
	}
	if lh, ok := m["line_height"].(map[string]interface{}); ok {
		writeNumbers(sb, lh, "--leading")
	}
	if ls, ok := m["letter_spacing"].(map[string]interface{}); ok {
		writeNumbers(sb, ls, "--tracking")
	}
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

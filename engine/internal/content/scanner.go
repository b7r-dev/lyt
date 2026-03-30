package content

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ContentFile represents a parsed YAML content file
type ContentFile struct {
	Path     string
	RelPath  string
	Type     string // "page", "blog", "component", "config"
	Data     map[string]interface{}
	Markdown map[string]string // key → resolved @file.md content
}

// Collection holds all content
type Collection struct {
	Pages      []ContentFile
	Blog       []ContentFile
	Components []ContentFile
	Config     map[string]interface{}
	Tokens     map[string]interface{}
}

// Scanner scans the content directory
type Scanner struct {
	dir     string
	verbose bool
}

func NewScanner(dir string, verbose bool) *Scanner {
	return &Scanner{dir: dir, verbose: verbose}
}

func (s *Scanner) Scan() (*Collection, error) {
	c := &Collection{
		Pages:      []ContentFile{},
		Blog:       []ContentFile{},
		Components: []ContentFile{},
		Config:     make(map[string]interface{}),
		Tokens:     make(map[string]interface{}),
	}

	err := filepath.Walk(s.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !s.isYAML(path) {
			return nil
		}

		rel, _ := filepath.Rel(s.dir, path)
		rel = filepath.ToSlash(rel)

		// Skip data files
		if strings.HasPrefix(rel, "data/") || strings.HasPrefix(rel, "_") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		var raw map[string]interface{}
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return nil
		}

		if raw == nil || len(raw) == 0 {
			return nil
		}

		cf := ContentFile{
			Path:     path,
			RelPath:  rel,
			Data:     raw,
			Markdown: s.resolveMarkdown(raw, filepath.Dir(path)),
		}

		switch {
		case strings.HasPrefix(rel, "blog/"):
			cf.Type = "blog"
			c.Blog = append(c.Blog, cf)
		case strings.HasPrefix(rel, "components/"):
			cf.Type = "component"
			c.Components = append(c.Components, cf)
		case strings.HasPrefix(rel, "config/"):
			cf.Type = "config"
			for k, v := range raw {
				c.Config[k] = v
			}
		case strings.HasPrefix(rel, "tokens.yaml"):
			c.Tokens = raw
		default:
			cf.Type = "page"
			c.Pages = append(c.Pages, cf)
		}

		return nil
	})

	return c, err
}

func (s *Scanner) isYAML(path string) bool {
	return strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")
}

// resolveMarkdown scans string fields for @file.md references and loads them
func (s *Scanner) resolveMarkdown(data map[string]interface{}, baseDir string) map[string]string {
	resolved := make(map[string]string)
	s.resolveMap(data, baseDir, resolved)
	return resolved
}

func (s *Scanner) resolveMap(data map[string]interface{}, baseDir string, resolved map[string]string) {
	for _, v := range data {
		switch x := v.(type) {
		case string:
			if strings.HasPrefix(x, "@") {
				path := strings.TrimPrefix(x, "@")
				if content, err := os.ReadFile(filepath.Join(baseDir, path)); err == nil {
					resolved[x] = string(content)
				}
			}
		case map[string]interface{}:
			s.resolveMap(x, baseDir, resolved)
		case []interface{}:
			s.resolveSlice(x, baseDir, resolved)
		}
	}
}

func (s *Scanner) resolveSlice(data []interface{}, baseDir string, resolved map[string]string) {
	for _, v := range data {
		switch x := v.(type) {
		case string:
			// field refs like "@filename.md" at list level
		case map[string]interface{}:
			s.resolveMap(x, baseDir, resolved)
		case []interface{}:
			s.resolveSlice(x, baseDir, resolved)
		}
	}
}

func (s *Scanner) GetField(cf ContentFile, path string) interface{} {
	parts := strings.Split(path, ".")
	var val interface{} = cf.Data
	for _, part := range parts {
		if m, ok := val.(map[string]interface{}); ok {
			val = m[part]
		} else {
			return nil
		}
	}
	return val
}

// GetMeta extracts meta from a ContentFile
func GetMeta(cf ContentFile) map[string]interface{} {
	if meta, ok := cf.Data["meta"].(map[string]interface{}); ok {
		return meta
	}
	return make(map[string]interface{})
}

// GetSection extracts a named section from page data
func GetSection(cf ContentFile, name string) map[string]interface{} {
	if sections, ok := cf.Data["sections"].([]interface{}); ok {
		for _, s := range sections {
			if m, ok := s.(map[string]interface{}); ok {
				if id, ok := m["id"].(string); ok && id == name {
					return m
				}
			}
		}
	}
	return nil
}

// ToStrings converts interface slice to string slice
func ToStrings(v interface{}) []string {
	if slice, ok := v.([]interface{}); ok {
		out := make([]string, len(slice))
		for i, x := range slice {
			if s, ok := x.(string); ok {
				out[i] = s
			}
		}
		return out
	}
	return nil
}

// ToMapStringString converts interface map to map[string]string
func ToMapStringString(v interface{}) map[string]string {
	if m, ok := v.(map[string]interface{}); ok {
		out := make(map[string]string)
		for k, x := range m {
			if s, ok := x.(string); ok {
				out[k] = s
			}
		}
		return out
	}
	return nil
}

func (s *Collection) GetComponent(name string) *ContentFile {
	for _, c := range s.Components {
		if n, ok := c.Data["name"].(string); ok && n == name {
			return &c
		}
	}
	return nil
}

func (s *Collection) GetNav() []interface{} {
	if nav, ok := s.Config["nav"].([]interface{}); ok {
		return nav
	}
	return []interface{}{}
}

func (s *Collection) GetSiteTitle() string {
	if meta, ok := s.Config["meta"].(map[string]interface{}); ok {
		if title, ok := meta["title"].(string); ok {
			return title
		}
	}
	return "lyt"
}

// AgentSectionEnabled returns true if agent-specific content generation is enabled
func (s *Collection) AgentSectionEnabled() bool {
	if agentSection, ok := s.Config["agent_section"].(map[string]interface{}); ok {
		if enabled, ok := agentSection["enabled"].(bool); ok && enabled {
			return true
		}
	}
	return false
}

// GetAgentPath returns the URL path for agent content (default: /agents)
func (s *Collection) GetAgentPath() string {
	if agentSection, ok := s.Config["agent_section"].(map[string]interface{}); ok {
		if path, ok := agentSection["path"].(string); ok && path != "" {
			return path
		}
	}
	return "/agents"
}

// GetURL returns the site URL from config
func (s *Collection) GetURL() string {
	if meta, ok := s.Config["meta"].(map[string]interface{}); ok {
		if url, ok := meta["url"].(string); ok && url != "" {
			return url
		}
	}
	return "https://lyt.b7r.dev"
}

// GetCopyright returns the copyright string from config
func (s *Collection) GetCopyright() string {
	if meta, ok := s.Config["meta"].(map[string]interface{}); ok {
		if copyright, ok := meta["copyright"].(string); ok && copyright != "" {
			return copyright
		}
	}
	return "Copyright © 2026 Aggressively Beige Holdings, LLC"
}

// GetLicense returns the license from config
func (s *Collection) GetLicense() string {
	if meta, ok := s.Config["meta"].(map[string]interface{}); ok {
		if license, ok := meta["license"].(string); ok && license != "" {
			return license
		}
	}
	return "MIT"
}

// HasAgentPage returns true if a content file should have its own agent version
// Requires both: meta.agent: true AND agent_content section exists
func (s *Collection) HasAgentPage(cf ContentFile) bool {
	// Must have meta.agent: true
	meta := GetMeta(cf)
	if agent, ok := meta["agent"].(bool); !ok || !agent {
		return false
	}
	// Must have agent_content section defined
	if GetAgentContent(cf) != nil {
		return true
	}
	return false
}

// ShowAgentHubLink returns true if a page should link to the agent hub (/agents)
// True if meta.agent: true (even without agent_content, link to hub)
func (s *Collection) ShowAgentHubLink(cf ContentFile) bool {
	meta := GetMeta(cf)
	if agent, ok := meta["agent"].(bool); ok && agent {
		return true
	}
	return false
}

// GetAgentContent extracts agent-specific content from a content file
// Returns the agent_content section if present
func GetAgentContent(cf ContentFile) map[string]interface{} {
	if agentContent, ok := cf.Data["agent_content"].(map[string]interface{}); ok {
		return agentContent
	}
	return nil
}

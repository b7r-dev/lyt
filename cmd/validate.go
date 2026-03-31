package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/b7r-dev/lyt/internal/content"
	"github.com/spf13/cobra"
)

// ValidationError holds validation failure information
type ValidationError struct {
	File    string
	Line    int
	Field   string
	Type    string
	Message string
}

var (
	validateStrict bool
	validateFix    bool
	validateSchema bool
	validateLinks  bool
	validateDir    string
)

// Valid section types
var validSectionTypes = map[string]bool{
	"hero":         true,
	"default":      true,
	"features":     true,
	"callout":      true,
	"cta":          true,
	"warning":      true,
	"pull-quote":   true,
	"citation":     true,
	"code-example": true,
	"divider":      true,
	"about":        true,
}

// Valid agent section types
var validAgentSectionTypes = map[string]bool{
	"default": true,
	"cli":     true,
	"schema":  true,
	"example": true,
	"link":    true,
}

// Valid button variants
var validButtonVariants = map[string]bool{
	"primary": true,
	"outline": true,
}

// Valid callout variants
var validCalloutVariants = map[string]bool{
	"info": true,
	"tip":  true,
	"note": true,
}

// Valid warning variants
var validWarningVariants = map[string]bool{
	"warning": true,
	"error":   true,
	"info":    true,
	"success": true,
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate content (schema and/or links)",
	Long: `Validate lyt content.
	
Without flags, validates both schema and links.
	
Examples:
  lyt validate              # Validate schema and links
  lyt validate --schema    # Just schema validation
  lyt validate --links     # Just link validation
  lyt validate --schema --links  # Explicitly both

Exit code 0 = all valid
Exit code 1 = validation errors found`,
	RunE: runValidate,
}

func init() {
	validateCmd.Flags().BoolVarP(&validateStrict, "strict", "s", false, "Treat warnings as errors")
	validateCmd.Flags().BoolVarP(&validateFix, "fix", "f", false, "Attempt to fix common issues")
	validateCmd.Flags().BoolVar(&validateSchema, "schema", false, "Validate content schema")
	validateCmd.Flags().BoolVar(&validateLinks, "links", false, "Validate internal links")
	validateCmd.Flags().StringVar(&validateDir, "dir", "", "Directory with dist files to check links (default: ./dist)")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	// Default: run both schema and links validation
	runSchema := validateSchema
	runLinks := validateLinks
	if !runSchema && !runLinks {
		runSchema = true
		runLinks = true
	}

	var hasErrors bool

	// Schema validation
	if runSchema {
		if err := runSchemaValidation(); err != nil {
			hasErrors = true
		}
	}

	// Link validation
	if runLinks {
		if err := runLinkValidation(); err != nil {
			hasErrors = true
		}
	}

	if hasErrors {
		return fmt.Errorf("validation failed")
	}

	return nil
}

// runSchemaValidation validates content against schema
func runSchemaValidation() error {
	projectRoot, err := detectProjectRoot()
	if err != nil {
		return fmt.Errorf("not a valid lyt project: %w", err)
	}

	contentDir := filepath.Join(projectRoot, "content")

	// Check if schema exists
	schemaPath := filepath.Join(contentDir, "schema.yaml")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		fmt.Printf("⚠️  No schema.yaml found in content/ - skipping schema validation\n")
		fmt.Println("   Create content/schema.yaml to enable validation")
		return nil
	}

	fmt.Println("🔍 Validating content against schema...")

	scanner := content.NewScanner(contentDir, false)
	collection, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	var errors []ValidationError

	// Validate pages
	for _, page := range collection.Pages {
		if errs := validatePage(page); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Validate blog posts
	for _, post := range collection.Blog {
		if errs := validateBlog(post); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Validate config
	for k, v := range collection.Config {
		if errs := validateConfigField(k, v); len(errs) > 0 {
			errors = append(errors, errs...)
		}
	}

	// Print results
	if len(errors) == 0 {
		fmt.Printf("✅ Schema valid (%d pages, %d blog posts)\n",
			len(collection.Pages), len(collection.Blog))
		return nil
	}

	fmt.Printf("❌ Found %d schema error(s):\n\n", len(errors))
	for _, e := range errors {
		loc := e.File
		if e.Line > 0 {
			loc = fmt.Sprintf("%s:%d", e.File, e.Line)
		}
		fmt.Printf("  ✗ %s\n", loc)
		if e.Field != "" {
			fmt.Printf("      Field: %s\n", e.Field)
		}
		fmt.Printf("      %s\n\n", e.Message)
	}

	if validateStrict {
		fmt.Println("💥 Strict mode: treating warnings as errors")
		return fmt.Errorf("schema validation failed")
	}

	return fmt.Errorf("schema validation failed with %d error(s)", len(errors))
}

func validatePage(cf content.ContentFile) []ValidationError {
	var errors []ValidationError

	meta := content.GetMeta(cf)

	// Check required meta fields
	if _, ok := meta["title"]; !ok {
		errors = append(errors, ValidationError{
			File:    cf.RelPath,
			Field:   "meta.title",
			Type:    "required",
			Message: "title is required",
		})
	}

	if _, ok := meta["slug"]; !ok {
		errors = append(errors, ValidationError{
			File:    cf.RelPath,
			Field:   "meta.slug",
			Type:    "required",
			Message: "slug is required",
		})
	} else {
		// Validate slug format
		slug, _ := meta["slug"].(string)
		if !strings.HasPrefix(slug, "/") && cf.RelPath != "pages/agents.yaml" {
			errors = append(errors, ValidationError{
				File:    cf.RelPath,
				Field:   "meta.slug",
				Type:    "format",
				Message: fmt.Sprintf("slug should start with '/' (got %q)", slug),
			})
		}
	}

	// Check agent_content if meta.agent is true
	agent, _ := meta["agent"].(bool)
	if agent {
		if cf.Data["agent_content"] == nil {
			errors = append(errors, ValidationError{
				File:    cf.RelPath,
				Field:   "agent_content",
				Type:    "required",
				Message: "agent_content required when meta.agent is true",
			})
		}
	}

	// Validate sections
	if sections, ok := cf.Data["sections"].([]interface{}); ok {
		for i, s := range sections {
			if sec, ok := s.(map[string]interface{}); ok {
				if errs := validateSection(cf.RelPath, i, sec); len(errs) > 0 {
					errors = append(errors, errs...)
				}
			}
		}
	}

	// Validate agent_content sections if present
	if ac, ok := cf.Data["agent_content"].(map[string]interface{}); ok {
		if sections, ok := ac["sections"].([]interface{}); ok {
			for i, s := range sections {
				if sec, ok := s.(map[string]interface{}); ok {
					if errs := validateAgentSection(cf.RelPath, i, sec); len(errs) > 0 {
						errors = append(errors, errs...)
					}
				}
			}
		}
	}

	return errors
}

func validateBlog(cf content.ContentFile) []ValidationError {
	var errors []ValidationError

	meta := content.GetMeta(cf)

	// Check required meta fields
	if _, ok := meta["title"]; !ok {
		errors = append(errors, ValidationError{
			File:    cf.RelPath,
			Field:   "meta.title",
			Type:    "required",
			Message: "title is required",
		})
	}

	if _, ok := meta["slug"]; !ok {
		errors = append(errors, ValidationError{
			File:    cf.RelPath,
			Field:   "meta.slug",
			Type:    "required",
			Message: "slug is required",
		})
	} else {
		// Validate slug format (no leading slash for blog)
		slug, _ := meta["slug"].(string)
		if strings.HasPrefix(slug, "/") {
			errors = append(errors, ValidationError{
				File:    cf.RelPath,
				Field:   "meta.slug",
				Type:    "format",
				Message: "blog slug should not start with '/'",
			})
		}
		// Check for valid slug characters
		validSlug := regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)
		if !validSlug.MatchString(slug) {
			errors = append(errors, ValidationError{
				File:    cf.RelPath,
				Field:   "meta.slug",
				Type:    "format",
				Message: "slug should be lowercase alphanumeric with hyphens",
			})
		}
	}

	// Validate date format if present
	if date, ok := meta["date"].(string); ok {
		validDate := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		if !validDate.MatchString(date) {
			errors = append(errors, ValidationError{
				File:    cf.RelPath,
				Field:   "meta.date",
				Type:    "format",
				Message: "date should be YYYY-MM-DD format",
			})
		}
	}

	// Validate sections
	if sections, ok := cf.Data["sections"].([]interface{}); ok {
		for i, s := range sections {
			if sec, ok := s.(map[string]interface{}); ok {
				if errs := validateSection(cf.RelPath, i, sec); len(errs) > 0 {
					errors = append(errors, errs...)
				}
			}
		}
	}

	return errors
}

func validateSection(file string, index int, sec map[string]interface{}) []ValidationError {
	var errors []ValidationError

	// Check for type field
	secType, ok := sec["type"].(string)
	if !ok {
		errors = append(errors, ValidationError{
			File:    file,
			Line:    index + 1,
			Field:   fmt.Sprintf("sections[%d].type", index),
			Type:    "required",
			Message: "section type is required",
		})
		return errors
	}

	// Validate section type
	if !validSectionTypes[secType] {
		errors = append(errors, ValidationError{
			File:    file,
			Line:    index + 1,
			Field:   fmt.Sprintf("sections[%d].type", index),
			Type:    "invalid",
			Message: fmt.Sprintf("invalid section type %q", secType),
		})
	}

	// Type-specific validation
	switch secType {
	case "cta":
		if _, ok := sec["button_text"]; !ok {
			errors = append(errors, ValidationError{
				File:    file,
				Line:    index + 1,
				Field:   fmt.Sprintf("sections[%d]", index),
				Type:    "required",
				Message: "cta section requires button_text",
			})
		}
		if _, ok := sec["button_href"]; !ok {
			errors = append(errors, ValidationError{
				File:    file,
				Line:    index + 1,
				Field:   fmt.Sprintf("sections[%d]", index),
				Type:    "required",
				Message: "cta section requires button_href",
			})
		}

	case "callout":
		if variant, ok := sec["variant"].(string); ok && !validCalloutVariants[variant] {
			errors = append(errors, ValidationError{
				File:    file,
				Line:    index + 1,
				Field:   fmt.Sprintf("sections[%d].variant", index),
				Type:    "invalid",
				Message: fmt.Sprintf("invalid callout variant %q", variant),
			})
		}

	case "warning":
		if variant, ok := sec["variant"].(string); ok && !validWarningVariants[variant] {
			errors = append(errors, ValidationError{
				File:    file,
				Line:    index + 1,
				Field:   fmt.Sprintf("sections[%d].variant", index),
				Type:    "invalid",
				Message: fmt.Sprintf("invalid warning variant %q", variant),
			})
		}

	case "code-example":
		if _, ok := sec["code"]; !ok {
			errors = append(errors, ValidationError{
				File:    file,
				Line:    index + 1,
				Field:   fmt.Sprintf("sections[%d]", index),
				Type:    "required",
				Message: "code-example section requires code field",
			})
		}

	case "features":
		if cards, ok := sec["cards"].([]interface{}); ok {
			for ci, card := range cards {
				if c, ok := card.(map[string]interface{}); ok {
					if _, ok := c["title"]; !ok {
						errors = append(errors, ValidationError{
							File:    file,
							Line:    index + 1,
							Field:   fmt.Sprintf("sections[%d].cards[%d]", index, ci),
							Type:    "required",
							Message: "feature card requires title",
						})
					}
				}
			}
		}

	case "hero":
		if _, ok := sec["title"]; !ok {
			errors = append(errors, ValidationError{
				File:    file,
				Line:    index + 1,
				Field:   fmt.Sprintf("sections[%d]", index),
				Type:    "required",
				Message: "hero section requires title",
			})
		}
	}

	// Validate buttons
	if buttons, ok := sec["buttons"].([]interface{}); ok {
		for bi, button := range buttons {
			if btn, ok := button.(map[string]interface{}); ok {
				if _, ok := btn["text"]; !ok {
					errors = append(errors, ValidationError{
						File:    file,
						Line:    index + 1,
						Field:   fmt.Sprintf("sections[%d].buttons[%d]", index, bi),
						Type:    "required",
						Message: "button requires text",
					})
				}
				if _, ok := btn["href"]; !ok {
					errors = append(errors, ValidationError{
						File:    file,
						Line:    index + 1,
						Field:   fmt.Sprintf("sections[%d].buttons[%d]", index, bi),
						Type:    "required",
						Message: "button requires href",
					})
				}
				if variant, ok := btn["variant"].(string); ok && !validButtonVariants[variant] {
					errors = append(errors, ValidationError{
						File:    file,
						Line:    index + 1,
						Field:   fmt.Sprintf("sections[%d].buttons[%d].variant", index, bi),
						Type:    "invalid",
						Message: fmt.Sprintf("invalid button variant %q", variant),
					})
				}
			}
		}
	}

	return errors
}

func validateAgentSection(file string, index int, sec map[string]interface{}) []ValidationError {
	var errors []ValidationError

	// Check for type field
	secType, ok := sec["type"].(string)
	if !ok {
		errors = append(errors, ValidationError{
			File:    file,
			Line:    index + 1,
			Field:   fmt.Sprintf("agent_content.sections[%d].type", index),
			Type:    "required",
			Message: "agent section type is required",
		})
		return errors
	}

	// Validate agent section type
	if !validAgentSectionTypes[secType] {
		errors = append(errors, ValidationError{
			File:    file,
			Line:    index + 1,
			Field:   fmt.Sprintf("agent_content.sections[%d].type", index),
			Type:    "invalid",
			Message: fmt.Sprintf("invalid agent section type %q", secType),
		})
	}

	// Type-specific validation
	switch secType {
	case "cli":
		if _, ok := sec["commands"]; !ok {
			errors = append(errors, ValidationError{
				File:    file,
				Line:    index + 1,
				Field:   fmt.Sprintf("agent_content.sections[%d]", index),
				Type:    "required",
				Message: "cli section requires commands field",
			})
		}

	case "example":
		if _, ok := sec["example"]; !ok {
			errors = append(errors, ValidationError{
				File:    file,
				Line:    index + 1,
				Field:   fmt.Sprintf("agent_content.sections[%d]", index),
				Type:    "required",
				Message: "example section requires example field",
			})
		}

	case "link":
		if _, ok := sec["text"]; !ok {
			errors = append(errors, ValidationError{
				File:    file,
				Line:    index + 1,
				Field:   fmt.Sprintf("agent_content.sections[%d]", index),
				Type:    "required",
				Message: "link section requires text field",
			})
		}
		if _, ok := sec["href"]; !ok {
			errors = append(errors, ValidationError{
				File:    file,
				Line:    index + 1,
				Field:   fmt.Sprintf("agent_content.sections[%d]", index),
				Type:    "required",
				Message: "link section requires href field",
			})
		}
	}

	return errors
}

func validateConfigField(key string, value interface{}) []ValidationError {
	var errors []ValidationError

	switch key {
	case "nav":
		if nav, ok := value.([]interface{}); ok {
			for i, item := range nav {
				if n, ok := item.(map[string]interface{}); ok {
					if _, ok := n["label"]; !ok {
						errors = append(errors, ValidationError{
							File:    "config/site.yaml",
							Field:   fmt.Sprintf("nav[%d]", i),
							Type:    "required",
							Message: "nav item requires label",
						})
					}
					if _, ok := n["href"]; !ok {
						errors = append(errors, ValidationError{
							File:    "config/site.yaml",
							Field:   fmt.Sprintf("nav[%d]", i),
							Type:    "required",
							Message: "nav item requires href",
						})
					}
				}
			}
		}

	case "agent_section":
		if as, ok := value.(map[string]interface{}); ok {
			if enabled, ok := as["enabled"].(bool); ok && enabled {
				if path, ok := as["path"].(string); ok && path != "" {
					if !strings.HasPrefix(path, "/") {
						errors = append(errors, ValidationError{
							File:    "config/site.yaml",
							Field:   "agent_section.path",
							Type:    "format",
							Message: "agent_section.path should start with '/'",
						})
					}
				}
			}
		}
	}

	return errors
}

// LinkValidationError represents a broken link
type LinkValidationError struct {
	SourceFile string
	Line       int
	Link       string
	Message    string
}

// runLinkValidation validates internal links in built HTML files
func runLinkValidation() error {
	// Determine dist directory
	distDir := validateDir
	if distDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}
		distDir = filepath.Join(cwd, "dist")
	}

	// Check if dist exists
	if _, err := os.Stat(distDir); os.IsNotExist(err) {
		fmt.Printf("⚠️  No dist directory found at %s - skipping link validation\n", distDir)
		fmt.Println("   Run 'lyt build' first, or use --dir to specify a different directory")
		return nil
	}

	// Read sitemap.xml to get URLs
	sitemapPath := filepath.Join(distDir, "sitemap.xml")
	sitemapData, err := os.ReadFile(sitemapPath)
	if err != nil {
		fmt.Printf("⚠️  No sitemap.xml found - skipping link validation\n")
		return nil
	}

	// Parse URLs from sitemap (simple parsing)
	var urls []string
	lines := strings.Split(string(sitemapData), "\n")
	for _, line := range lines {
		if strings.Contains(line, "<loc>") {
			start := strings.Index(line, "<loc>") + 5
			end := strings.Index(line, "</loc>")
			if start > 4 && end > start {
				url := line[start:end]
				// Strip base URL to get path
				url = strings.TrimPrefix(url, "https://lyt.local")
				url = strings.TrimPrefix(url, "http://lyt.local")
				urls = append(urls, url)
			}
		}
	}

	if len(urls) == 0 {
		fmt.Println("⚠️  No URLs found in sitemap - skipping link validation")
		return nil
	}

	fmt.Printf("🔗 Checking %d pages for broken links...\n", len(urls))

	var errors []LinkValidationError

	// Check each HTML file for broken links
	for _, url := range urls {
		// Convert URL to file path
		filePath := filepath.Join(distDir, url)
		if url == "/" || url == "" {
			filePath = filepath.Join(distDir, "index.html")
		} else {
			filePath = filepath.Join(filePath, "index.html")
		}

		// Read the HTML file
		htmlData, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip if can't read
		}

		html := string(htmlData)

		// Find all links in the HTML
		// Match href="..." or href='...'
		linkRegex := regexp.MustCompile(`href=["']([^"']+)["']`)
		matches := linkRegex.FindAllStringSubmatchIndex(html, -1)

		for _, match := range matches {
			if len(match) < 4 {
				continue
			}
			link := html[match[2]:match[3]]

			// Skip external links, anchors, and mailto
			if strings.HasPrefix(link, "http") ||
				strings.HasPrefix(link, "//") ||
				strings.HasPrefix(link, "#") ||
				strings.HasPrefix(link, "mailto:") ||
				strings.HasPrefix(link, "tel:") {
				continue
			}

			// Handle ./ prefix (same directory)
			link = strings.TrimPrefix(link, "./")

			// Handle fragment identifiers (e.g., /docs/content#section)
			fragment := ""
			if idx := strings.Index(link, "#"); idx != -1 {
				fragment = link[idx:]
				link = link[:idx]
			}

			// Skip empty links after removing fragment
			if link == "" {
				continue
			}

			// Resolve relative link
			resolved := resolveURL(url, link)

			// Check if the resolved path exists
			targetPath := filepath.Join(distDir, resolved)
			if !strings.HasSuffix(resolved, ".html") && !strings.HasSuffix(resolved, "/") {
				targetPath = filepath.Join(targetPath, "index.html")
			}

			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				// Try with .html
				if !strings.HasSuffix(resolved, ".html") {
					targetPath = resolved + ".html"
				}
				if _, err := os.Stat(targetPath); os.IsNotExist(err) {
					errors = append(errors, LinkValidationError{
						SourceFile: filePath,
						Link:       link + fragment,
						Message:    fmt.Sprintf("broken link: %s -> %s (resolved: %s)", url, link, resolved),
					})
				}
			}
		}
	}

	// Print results
	if len(errors) == 0 {
		fmt.Printf("✅ All links valid (%d pages checked)\n", len(urls))
		return nil
	}

	fmt.Printf("❌ Found %d broken link(s):\n\n", len(errors))
	for _, e := range errors {
		relPath, _ := filepath.Rel(distDir, e.SourceFile)
		fmt.Printf("  ✗ %s\n", relPath)
		fmt.Printf("      Broken link: %s\n\n", e.Link)
	}

	return fmt.Errorf("link validation failed with %d broken link(s)", len(errors))
}

// resolveURL resolves a relative URL against a base URL
func resolveURL(base, relative string) string {
	// Handle absolute links
	if strings.HasPrefix(relative, "/") {
		return relative
	}

	// Normalize base - remove trailing slash
	base = strings.TrimSuffix(base, "/")

	// Handle empty or root base
	if base == "" || base == "/" {
		if relative != "" {
			return "/" + relative
		}
		return "/"
	}

	// Handle parent directory references
	if strings.HasPrefix(relative, "../") {
		// Count how many levels up we need to go
		parts := strings.Split(base, "/")
		relParts := strings.Split(relative, "/")

		upCount := 0
		for _, part := range relParts {
			if part == ".." {
				upCount++
			}
		}

		// Go up the required number of levels
		newParts := parts[:len(parts)-upCount]
		if len(newParts) == 0 {
			newParts = []string{""}
		}

		// Append any remaining path parts
		for _, part := range relParts {
			if part != ".." && part != "" {
				newParts = append(newParts, part)
			}
		}

		result := strings.Join(newParts, "/")
		if !strings.HasPrefix(result, "/") {
			result = "/" + result
		}
		return result
	}

	// Simple relative path - just join with base directory
	// Get the directory part of the base path
	// For "/docs", we want to keep "/docs" as the base
	if base == "/" || base == "" {
		return "/" + relative
	}

	// Base is like "/docs" or "/docs/components"
	// Find the last slash
	idx := strings.LastIndex(base, "/")
	if idx >= 0 {
		// There may be a slash in the path
		return base + "/" + relative
	}

	return "/" + relative
}

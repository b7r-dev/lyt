package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/b7r-dev/lyt/engine/internal/build"
	"github.com/b7r-dev/lyt/engine/internal/content"
	"github.com/b7r-dev/lyt/engine/internal/render"
	tokens "github.com/b7r-dev/lyt/engine/internal/tokens"
	"github.com/spf13/cobra"
)

// Project directories that must exist for a valid lyt project
var requiredDirs = []string{"content", "templates"}

// detectProjectRoot finds the project root by looking for required directories.
// Returns the project root path and an error if not a valid project.
func detectProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}

	// Check if current directory is a valid project
	valid := true
	for _, dir := range requiredDirs {
		if _, err := os.Stat(filepath.Join(cwd, dir)); os.IsNotExist(err) {
			valid = false
			break
		}
	}

	if valid {
		return cwd, nil
	}

	// Check parent directory (in case user runs from engine/ subdir)
	parent := filepath.Dir(cwd)
	valid = true
	for _, dir := range requiredDirs {
		if _, err := os.Stat(filepath.Join(parent, dir)); os.IsNotExist(err) {
			valid = false
			break
		}
	}

	if valid {
		return parent, nil
	}

	return "", fmt.Errorf("not a valid lyt project: missing content/ or templates/ directory")
}

var buildOutput string
var forceRebuild bool

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the static site",
	RunE:  runBuild,
}

func runBuild(cmd *cobra.Command, args []string) error {
	// Detect project root
	projectRoot, err := detectProjectRoot()
	if err != nil {
		return err
	}

	start := time.Now()

	// Use current directory as project root for output
	// This allows "cd /path/to/my-blog && lyt build" to output to ./dist
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	contentDir := filepath.Join(projectRoot, "content")
	templatesDir := filepath.Join(projectRoot, "templates")
	publicDir := filepath.Join(projectRoot, "public")

	outputDir := buildOutput
	if outputDir == "" {
		// Default: output to ./dist in current working directory
		outputDir = filepath.Join(cwd, "dist")
	}

	fmt.Println("🔨 lyt build")
	if verbose && projectRoot != cwd {
		fmt.Printf("   Project: %s\n   Output: %s\n", projectRoot, outputDir)
	}

	// Initialize build cache for incremental builds
	cachePath := filepath.Join(outputDir, ".build.cache")
	cache := build.NewCache(cachePath)

	// Check if we can skip the build (incremental)
	if !forceRebuild && !cache.ShouldRebuild(contentDir, templatesDir) {
		fmt.Println("✅ No changes detected, skipping build")
		return nil
	}

	// Step 1: Scan and validate content
	fmt.Println("📝 Scanning content...")
	scanner := content.NewScanner(contentDir, verbose)
	collection, err := scanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if verbose {
		fmt.Printf("   Found %d pages, %d blog posts, %d components\n",
			len(collection.Pages), len(collection.Blog), len(collection.Components))
	}

	// Step 2: Create output dir
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("mkdir failed: %w", err)
	}

	// Step 3: Process design tokens → CSS
	fmt.Println("🎨 Processing design tokens...")
	tokensCSS, err := tokens.ProcessTokens(filepath.Join(contentDir, "tokens.yaml"), verbose)
	if err != nil {
		fmt.Printf("   ⚠️  tokens skipped: %v\n", err)
		tokensCSS = "/* no tokens */"
	}
	tokensPath := filepath.Join(outputDir, "tokens.css")
	if err := os.WriteFile(tokensPath, []byte(tokensCSS), 0644); err != nil {
		return fmt.Errorf("write tokens failed: %w", err)
	}
	cache.Touch(filepath.Join(contentDir, "tokens.yaml"))

	// Step 4: Copy base CSS
	baseCSS, err := os.ReadFile(filepath.Join(templatesDir, "base.css"))
	if err != nil {
		baseCSS = []byte("/* no base.css */")
	}
	if err := os.WriteFile(filepath.Join(outputDir, "base.css"), baseCSS, 0644); err != nil {
		return fmt.Errorf("write base.css failed: %w", err)
	}
	cache.Touch(filepath.Join(templatesDir, "base.css"))

	// Step 5: Generate pages
	fmt.Println("📄 Generating pages...")
	renderer := render.NewRenderer(collection, contentDir, verbose)

	pagesGenerated := 0
	agentPathPrefix := collection.GetAgentPath()
	for _, page := range collection.Pages {
		html, err := renderer.RenderPage(page)
		if err != nil {
			fmt.Printf("   ⚠️  %s: %v\n", page.RelPath, err)
			continue
		}
		outPath := pageSlugToPath(page.RelPath, outputDir)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("mkdir page failed: %w", err)
		}
		if err := os.WriteFile(outPath, []byte(html), 0644); err != nil {
			return fmt.Errorf("write page failed: %w", err)
		}
		pagesGenerated++
		cache.Touch(filepath.Join(contentDir, page.RelPath))

		// Generate agent page if this page has agent content
		if collection.HasAgentPage(page) {
			agentHTML, err := renderer.RenderAgentPage(page)
			if err == nil {
				agentOutPath := agentSlugToPath(page.RelPath, outputDir, agentPathPrefix)
				if err := os.MkdirAll(filepath.Dir(agentOutPath), 0755); err != nil {
					fmt.Printf("   ⚠️  agent mkdir failed: %v\n", err)
				} else if err := os.WriteFile(agentOutPath, []byte(agentHTML), 0644); err != nil {
					fmt.Printf("   ⚠️  agent write failed: %v\n", err)
				}
			}
		}
	}

	// Step 6: Generate blog posts
	for _, post := range collection.Blog {
		html, err := renderer.RenderBlogPost(post)
		if err != nil {
			fmt.Printf("   ⚠️  %s: %v\n", post.RelPath, err)
			continue
		}
		outPath := blogSlugToPath(post.RelPath, outputDir)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return fmt.Errorf("mkdir blog failed: %w", err)
		}
		if err := os.WriteFile(outPath, []byte(html), 0644); err != nil {
			return fmt.Errorf("write blog failed: %w", err)
		}
		pagesGenerated++
		cache.Touch(filepath.Join(contentDir, post.RelPath))

		// Generate agent page for blog post if it has agent content
		if collection.HasAgentPage(post) {
			agentHTML, err := renderer.RenderAgentPage(post)
			if err == nil {
				agentOutPath := agentSlugToPath(post.RelPath, outputDir, agentPathPrefix)
				if err := os.MkdirAll(filepath.Dir(agentOutPath), 0755); err != nil {
					fmt.Printf("   ⚠️  agent mkdir failed: %v\n", err)
				} else if err := os.WriteFile(agentOutPath, []byte(agentHTML), 0644); err != nil {
					fmt.Printf("   ⚠️  agent write failed: %v\n", err)
				}
			}
		}
	}

	// Step 7: Copy assets
	fmt.Println("📦 Copying assets...")
	if err := copyDir(publicDir, outputDir, verbose); err != nil {
		fmt.Printf("   ⚠️  assets skipped: %v\n", err)
	}

	// Step 8: Generate sitemap
	fmt.Println("🗺️  Generating sitemap...")
	sitemap := generateSitemap(collection, outputDir)
	if err := os.WriteFile(filepath.Join(outputDir, "sitemap.xml"), []byte(sitemap), 0644); err != nil {
		return fmt.Errorf("write sitemap failed: %w", err)
	}

	// Save the cache
	_ = cache.Save()

	elapsed := time.Since(start)
	fmt.Printf("✅ Built %d pages → %s (%.2fs)\n", pagesGenerated, outputDir, elapsed.Seconds())
	return nil
}

func pageSlugToPath(relPath, outputDir string) string {
	name := strings.TrimSuffix(relPath, ".yaml")
	name = strings.TrimPrefix(name, "pages/")

	// Special cases:
	// - "/" (root) -> index.html at root
	// - "/agents" -> agents/index.html (the agent hub)
	// - "/docs" -> docs/index.html
	// - Any other path like "/docs/getting-started" -> docs/getting-started.html

	// Handle root index
	if name == "" || name == "index" {
		return filepath.Join(outputDir, "index.html")
	}

	// For paths like "agents", "docs", etc. that don't end in a filename,
	// output as directory/index.html
	if !strings.Contains(name, "/") {
		// Single path segment - could be "agents", "docs", etc.
		// These should be directory-style URLs
		return filepath.Join(outputDir, name, "index.html")
	}

	// Multi-segment path like "docs/getting-started" -> docs/getting-started/index.html
	return filepath.Join(outputDir, name, "index.html")
}

func blogSlugToPath(relPath, outputDir string) string {
	name := strings.TrimSuffix(relPath, ".yaml")
	name = strings.TrimPrefix(name, "blog/")
	return filepath.Join(outputDir, "blog", name, "index.html")
}

func agentSlugToPath(relPath, outputDir, agentPathPrefix string) string {
	// Strip the directory prefix (pages/ or blog/)
	name := strings.TrimSuffix(relPath, ".yaml")
	name = strings.TrimPrefix(name, "pages/")
	name = strings.TrimPrefix(name, "blog/")

	// Clean the agent path prefix (remove leading/trailing slashes)
	agentPathPrefix = strings.Trim(agentPathPrefix, "/")

	// If the source page is already the agent hub itself (slug: "/agents"),
	// output to the root agent path, not /agents/agents
	if name == "agents" || name == "" || name == "index" {
		return filepath.Join(outputDir, agentPathPrefix, "index.html")
	}

	// Build the agent output path: /agents/docs/getting-started
	return filepath.Join(outputDir, agentPathPrefix, name, "index.html")
}

func generateSitemap(c *content.Collection, outputDir string) string {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	base := "https://lyt.local"
	sb.WriteString("  <url><loc>" + base + "/</loc></url>\n")

	for _, p := range c.Pages {
		if slug := getSlug(p.Data); slug != "" {
			sb.WriteString("  <url><loc>" + base + slug + "</loc></url>\n")
		}
	}
	for _, b := range c.Blog {
		if slug := getSlug(b.Data); slug != "" {
			sb.WriteString("  <url><loc>" + base + "/blog/" + slug + "</loc></url>\n")
		}
	}

	sb.WriteString("</urlset>\n")
	return sb.String()
}

func getSlug(data map[string]interface{}) string {
	if meta, ok := data["meta"].(map[string]interface{}); ok {
		if slug, ok := meta["slug"].(string); ok {
			return slug
		}
	}
	return ""
}

func copyDir(src, dst string, verbose bool) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		dest := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dest, data, info.Mode()); err != nil {
			return err
		}
		if verbose {
			fmt.Printf("   📄 %s\n", rel)
		}
		return nil
	})
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", "", "Output directory (default: ./dist)")
	buildCmd.Flags().BoolVarP(&forceRebuild, "force", "f", false, "Force rebuild even when no changes detected")
}

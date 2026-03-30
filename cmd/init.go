package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initForce bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new lyt project",
	Long: `Create a new lyt project with default content and templates.

This command creates a complete project structure with:
- content/pages/index.yaml - Home page with hero and features
- content/pages/about.yaml - About page template  
- content/blog/hello-world.yaml - Sample blog post
- content/tokens.yaml - Design tokens
- templates/base.css - Stylesheet
- public/ - Assets directory

Run from an empty directory or use --force to overwrite existing files.`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Check for existing project files
	hasExisting := checkExistingProject(cwd)
	if hasExisting && !initForce {
		fmt.Printf("⚠️  Directory already contains files. Use --force to overwrite.\n")
		return nil
	}

	fmt.Println("🎨 lyt init")

	// Create directory structure
	dirs := []string{
		"content/pages",
		"content/blog",
		"content/config",
		"templates",
		"public",
	}

	for _, dir := range dirs {
		path := filepath.Join(cwd, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
		if verbose {
			fmt.Printf("   📁 %s/\n", dir)
		}
	}

	// Create default content files
	files := map[string]string{
		"content/pages/index.yaml":      defaultIndexYAML,
		"content/pages/about.yaml":      defaultAboutYAML,
		"content/pages/components.yaml": defaultComponentsYAML,
		"content/blog/hello-world.yaml": defaultBlogYAML,
		"content/tokens.yaml":           defaultTokensYAML,
		"content/config/site.yaml":      defaultSiteYAML,
		"templates/base.css":            defaultCSS,
	}

	for path, content := range files {
		fullPath := filepath.Join(cwd, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("write file %s: %w", path, err)
		}
		if verbose {
			fmt.Printf("   📄 %s\n", path)
		}
	}

	fmt.Printf("✅ Project initialized in %s\n", cwd)
	fmt.Println("   Run 'lyt build' to create your site")
	fmt.Println("   Run 'lyt serve' to preview")

	return nil
}

func checkExistingProject(cwd string) bool {
	required := []string{"content", "templates"}
	for _, dir := range required {
		if _, err := os.Stat(filepath.Join(cwd, dir)); err == nil {
			return true
		}
	}
	return false
}

const defaultIndexYAML = `# Home page

meta:
  title: "Welcome"
  slug: "/"
  description: "A minimal static site built with lyt"

sections:
  - id: "hero"
    type: "hero"
    title: "Hello, World"
    subtitle: "Built with lyt"
    body: "A minimal static site generator. Edit YAML, run lyt build, deploy anywhere."
    buttons:
      - text: "Learn More"
        href: "/about"
        variant: "primary"
      - text: "View Source"
        href: "https://github.com/b7r-dev/lyt"
        variant: "outline"

  - id: "features"
    type: "features"
    title: "Why lyt?"
    cards:
      - title: "YAML + Markdown"
        body: "Content lives in YAML files. Markdown for prose. Clean, structured, version-controlled."
      - title: "Zero JS"
        body: "No runtime JavaScript in output. Pure HTML and CSS. Fast, accessible, no hydration."
      - title: "Go Speed"
        body: "Built in Go. Sub-second builds. Incremental compilation. Instant server restarts."
      - title: "Deploy Anywhere"
        body: "Static output. Works on Netlify, Vercel, Cloudflare Pages, or rsync to a VPS."

  - id: "cta"
    type: "cta"
    title: "Get Started"
    body: "Create your first project with lyt init. It's that simple."
    button_text: "Initialize Project"
    button_href: "/about"
`

const defaultAboutYAML = `# About page

meta:
  title: "About"
  slug: "/about"
  description: "About this site"

sections:
  - id: "hero"
    type: "hero"
    title: "About"
    subtitle: "Who we are"
    body: "This is an about page created by lyt init. Edit content/pages/about.yaml to customize it."

  - id: "mission"
    type: "default"
    title: "Our Mission"
    body: |
      We believe in:
      
      - **Simplicity** — The best tool is the one you don't notice
      - **Speed** — Waiting breaks flow state
      - **Control** — Your content, your structure, your output

  - id: "team"
    type: "features"
    title: "The Team"
    cards:
      - title: "You"
        body: "The content creator. The writer. The one with something to say."
      - title: "lyt"
        body: "The tool that gets out of your way and lets you write."
`

const defaultBlogYAML = `# Hello World blog post

meta:
  title: "Hello, World"
  slug: "hello-world"
  date: 2024-01-15
  tags:
    - intro
    - lyt
  draft: false

sections:
  - id: "intro"
    type: "default"
    title: "Hello, World"
    body: "Welcome to your new lyt blog! This post was created by lyt init.\n\nHere's what you can do:\n\n- Edit this file in content/blog/hello-world.yaml\n- Run lyt build to generate HTML\n- Run lyt serve to preview\n\n## Features\n\nlyt supports rich content components out of the box:"
      
  - id: "quote"
    type: "pull-quote"
    quote: "The best way to predict the future is to invent it."
    attribution: "Alan Kay"

  - id: "warning"
    type: "warning"
    title: "Work in Progress"
    body: "This site is under construction. More content coming soon."

  - id: "cta"
    type: "cta"
    title: "Enjoying this?"
    body: "Subscribe to get notified when new posts are published."
    button_text: "Subscribe"
    button_href: "/subscribe"

  - id: "citation"
    type: "citation"
    text: "The Pragmatic Programmer"
    url: "https://pragprog.com/titles/tpp20/the-pragmatic-programmer-20th-anniversary-edition/"
    publisher: "Addison-Wesley, 2019"
`

const defaultComponentsYAML = `# Components Gallery

meta:
  title: "Components"
  slug: "/components"
  description: "A showcase of all available lyt components"

sections:
  - id: "intro"
    type: "default"
    title: "Component Gallery"
    body: "Every component in lyt is defined in YAML and rendered to semantic HTML. No shortcodes, no templates - just declare what you need."

  - id: "hero-example"
    type: "hero"
    title: "Hero Section"
    subtitle: "Full-width opener"
    body: "The hero component creates an impactful opening section with title, subtitle, body text, and optional buttons."
    buttons:
      - text: "Primary"
        href: "#"
        variant: "primary"
      - text: "Outline"
        href: "#"
        variant: "outline"

  - id: "features-example"
    type: "features"
    title: "Features Grid"
    cards:
      - title: "Feature One"
        body: "Cards display in a responsive grid that adapts to screen size."
      - title: "Feature Two"
        body: "Cards have subtle depth with layered shadows."
      - title: "Feature Three"
        body: "Add as many cards as you need."

  - id: "callout-example"
    type: "callout"
    variant: "info"
    title: "Info Callout"
    body: "Use callouts to highlight tips, notes, or important information."

  - id: "quote-example"
    type: "pull-quote"
    quote: "Simplicity is the ultimate sophistication."
    attribution: "Leonardo da Vinci"

  - id: "cta-example"
    type: "cta"
    title: "Call to Action"
    body: "Encourage users to take action with a styled CTA section."
    button_text: "Get Started"
    button_href: "/"

  - id: "warning-example"
    type: "warning"
    variant: "warning"
    title: "Warning"
    body: "Alert users to important information with styled warning boxes."

  - id: "code-example"
    type: "code-example"
    title: "example.yaml"
    language: "yaml"
    code: |
      sections:
        - type: hero
          title: Hello
`

const defaultTokensYAML = `# Design tokens
# Customize colors, typography, and spacing

colors:
  # Primary palette - warm, desaturated beige
  background: "#f5f2eb"
  surface: "#ebe7de"
  surface_elevated: "#ffffff"
  
  # Text colors
  text_primary: "#2d2a26"
  text_secondary: "#5c5650"
  text_muted: "#8a837a"
  
  # Accent colors
  accent: "#c4a77d"
  accent_hover: "#b8956a"
  
  # Semantic colors
  link: "#8b6914"
  link_hover: "#6b5110"
  
  # Component colors
  warning_bg: "#fff3cd"
  warning_text: "#856404"
  cta_bg: "#2d2a26"
  cta_text: "#f5f2eb"

typography:
  # Font families
  font_sans: "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
  font_mono: "'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', monospace"
  
  # Font sizes (rem)
  text_xs: "0.75"
  text_sm: "0.875"
  text_base: "1"
  text_lg: "1.125"
  text_xl: "1.25"
  text_2xl: "1.5"
  text_3xl: "1.875"
  text_4xl: "2.25"
  
  # Line heights
  leading_tight: "1.25"
  leading_normal: "1.5"
  leading_relaxed: "1.625"

spacing:
  # Spacing scale (rem)
  space_0: "0"
  space_1: "0.25"
  space_2: "0.5"
  space_3: "0.75"
  space_4: "1"
  space_6: "1.5"
  space_8: "2"
  space_12: "3"
  space_16: "4"
  space_24: "6"

# Z-plane layers: 8 total (0: ground, 1-4: userland, 5-8: system)
z:
  # Userland (1-4)
  content: "1"      # Default user content
  elevated: "2"    # Elevated user content
  floating: "3"    # Floating user content (buttons, CTAs)
  top: "4"         # Highest userland
  
  # System (5-8)
  overlay: "5"      # Overlays
  nav: "6"          # Navigation (sticky header)
  tooltip: "7"      # Tooltips
  modal: "8"        # Modals, dialogs

layout:
  # Max content width
  max_width: "65rem"
  
  # Container padding
  padding_x: "1.5rem"
  
  # Border radius
  radius: "0.5rem"
  
  # Shadows (z-plane depth)
  shadow_sm: "0 1px 2px rgba(45, 42, 38, 0.05)"
  shadow_md: "0 4px 6px rgba(45, 42, 38, 0.07), 0 2px 4px rgba(45, 42, 38, 0.05)"
  shadow_lg: "0 10px 15px rgba(45, 42, 38, 0.1), 0 4px 6px rgba(45, 42, 38, 0.05)"
`

const defaultSiteYAML = `# Site configuration

meta:
  title: "My lyt Site"
  description: "Built with lyt"
  author: "Your Name"
  url: "https://example.com"

nav:
  - label: "Home"
    href: "/"
  - label: "About"
    href: "/about"
  - label: "Components"
    href: "/components"

build:
  # Output directory (relative to project root)
  output: "./dist"
  
  # Whether to minify HTML output
  minify: false
  
  # Whether to generate a sitemap
  sitemap: true

serve:
  # Development server port
  port: 8080
  
  # Whether to open browser automatically
  open: false
`

const defaultCSS = `/* lyt base stylesheet */
/* Edit this file to customize your site's appearance */

:root {
  /* Colors */
  --color-bg: #f5f2eb;
  --color-surface: #ebe7de;
  --color-surface-elevated: #ffffff;
  --color-text: #2d2a26;
  --color-text-secondary: #5c5650;
  --color-text-muted: #8a837a;
  --color-accent: #c4a77d;
  --color-accent-hover: #b8956a;
  --color-link: #8b6914;
  --color-link-hover: #6b5110;
  
  /* Typography */
  --font-sans: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  --font-mono: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', monospace;
  
  /* Spacing */
  --space-1: 0.25rem;
  --space-2: 0.5rem;
  --space-3: 0.75rem;
  --space-4: 1rem;
  --space-6: 1.5rem;
  --space-8: 2rem;
  --space-12: 3rem;
  
  /* Layout */
  --max-width: 65rem;
  --radius: 0.5rem;
}

/* Reset */
*, *::before, *::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

html {
  font-size: 16px;
  scroll-behavior: smooth;
}

body {
  font-family: var(--font-sans);
  background: var(--color-bg);
  color: var(--color-text);
  line-height: 1.5;
  min-height: 100vh;
}

/* Layout */
.container {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 0 var(--space-6);
}

main {
  padding: var(--space-8) 0;
  min-height: 80vh;
}

/* Typography */
h1, h2, h3, h4, h5, h6 {
  font-weight: 600;
  line-height: 1.25;
  margin-bottom: var(--space-4);
}

h1 { font-size: 2.25rem; }
h2 { font-size: 1.875rem; }
h3 { font-size: 1.5rem; }
h4 { font-size: 1.25rem; }

p {
  margin-bottom: var(--space-4);
}

a {
  color: var(--color-link);
  text-decoration: none;
}

a:hover {
  color: var(--color-link-hover);
  text-decoration: underline;
}

/* Navigation */
nav {
  background: var(--color-surface);
  padding: var(--space-4) 0;
  border-bottom: 1px solid rgba(0,0,0,0.05);
}

nav .container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

nav a {
  color: var(--color-text);
  margin-left: var(--space-4);
}

nav a:hover {
  color: var(--color-accent);
}

/* Buttons */
.btn {
  display: inline-block;
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius);
  font-weight: 500;
  text-decoration: none;
  cursor: pointer;
  border: none;
  transition: background 0.2s, color 0.2s;
}

.btn-primary {
  background: var(--color-accent);
  color: #fff;
}

.btn-primary:hover {
  background: var(--color-accent-hover);
  text-decoration: none;
}

.btn-outline {
  background: transparent;
  border: 2px solid var(--color-accent);
  color: var(--color-accent);
}

.btn-outline:hover {
  background: var(--color-accent);
  color: #fff;
  text-decoration: none;
}

/* Hero section */
.hero {
  text-align: center;
  padding: var(--space-12) var(--space-4);
  background: var(--color-surface);
  border-radius: var(--radius);
  margin-bottom: var(--space-8);
}

.hero h1 {
  font-size: 3rem;
  margin-bottom: var(--space-2);
}

.hero .subtitle {
  font-size: 1.25rem;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-4);
}

.hero .body {
  max-width: 40rem;
  margin: 0 auto var(--space-6);
}

.hero .buttons {
  display: flex;
  gap: var(--space-4);
  justify-content: center;
}

/* Features grid */
.features {
  margin: var(--space-8) 0;
}

.features h2 {
  text-align: center;
  margin-bottom: var(--space-8);
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: var(--space-6);
}

.feature-card {
  background: var(--color-surface-elevated);
  padding: var(--space-6);
  border-radius: var(--radius);
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.feature-card h3 {
  margin-bottom: var(--space-2);
}

/* CTA section */
.cta {
  background: var(--color-text);
  color: var(--color-bg);
  padding: var(--space-12) var(--space-6);
  border-radius: var(--radius);
  text-align: center;
  margin: var(--space-8) 0;
}

.cta h2 {
  color: var(--color-bg);
  margin-bottom: var(--space-2);
}

.cta .btn {
  margin-top: var(--space-4);
}

/* Pull quote */
.pull-quote {
  border-left: 4px solid var(--color-accent);
  padding-left: var(--space-6);
  margin: var(--space-8) 0;
  font-style: italic;
  font-size: 1.25rem;
  color: var(--color-text-secondary);
}

.pull-quote cite {
  display: block;
  margin-top: var(--space-2);
  font-size: 0.875rem;
  font-style: normal;
  color: var(--color-text-muted);
}

/* Warning */
.warning {
  background: #fff3cd;
  border: 1px solid #ffeeba;
  color: #856404;
  padding: var(--space-4);
  border-radius: var(--radius);
  margin: var(--space-4) 0;
}

/* Footer */
footer {
  background: var(--color-surface);
  padding: var(--space-8) 0;
  text-align: center;
  color: var(--color-text-secondary);
  margin-top: var(--space-12);
}

/* Utility */
.text-center { text-align: center; }
.mt-4 { margin-top: var(--space-4); }
.mb-4 { margin-bottom: var(--space-4); }

/* Responsive */
@media (max-width: 640px) {
  .hero h1 { font-size: 2rem; }
  .hero .buttons { flex-direction: column; }
  .container { padding: 0 var(--space-4); }
}

/* Dark Mode: Night sky with fireflies */
@media (prefers-color-scheme: dark) {
  :root {
    --color-bg: #0b0f14;
    --color-surface: #131920;
    --color-surface-elevated: #1a2029;
    --color-text: #e6e2db;
    --color-text-secondary: #a8a295;
    --color-muted: #6e6a63;
    --color-accent: #d4a574;
    --color-accent-hover: #f0b429;
    --color-link: #f0b429;
    --color-link-hover: #fbdb6b;
  }
  
  .card, .feature-card, .callout, .pull-quote, .cta, .warning {
    box-shadow: 2px 2px 0 0 #000, 4px 4px 8px rgba(0,0,0,0.6);
  }
  
  .warning {
    background: color-mix(in srgb, #f0b429 12%, #0b0f14);
    border-color: #f0b429;
    color: #e8e4dc;
    box-shadow: 0 0 12px color-mix(in srgb, #f0b429 30%, transparent);
  }
}
`

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Force initialization even if directory is not empty")
}

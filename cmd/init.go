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
		"content/pages/docs.yaml":       defaultDocsYAML,
		"content/pages/components.yaml": defaultComponentsYAML,
		"content/pages/blog.yaml":       defaultBlogIndexYAML,
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
  title: "Home"
  slug: "/"
  description: "A minimal static site built with lyt"
  agent: true

agent_content:
  title: "lyt — Agent Documentation"
  description: "lyt is a minimal static site generator written in Go."
  
  sections:
    - type: "default"
      title: "What is lyt?"
      body: |
        lyt transforms YAML content files into static HTML pages.
        
        Use lyt to:
        - Generate documentation sites
        - Build blogs with minimal dependencies
        - Create simple, fast static sites

    - type: "cli"
      title: "Installation"
      commands: |
        go install github.com/b7r-dev/lyt@latest
        lyt --version

    - type: "cli"
      title: "Quick Start"
      commands: |
        lyt init my-site
        cd my-site
        lyt build
        lyt serve

    - type: "schema"
      title: "Content Structure"
      schema: |
        content/pages/*.yaml:
          meta:
            title: "Page Title"
            slug: "/page-slug"
            description: "Description"
          
          sections:
            - type: "hero"
              title: "Title"

    - type: "link"
      text: "→ Full Documentation"
      href: "/agents/docs/getting-started"

    - type: "cli"
      title: "Validate Content"
      commands: |
        lyt validate          # Check content files
        lyt build --validate  # Build with validation

sections:
  - id: "hero"
    type: "hero"
    title: "lyt"
    subtitle: "yaml · markdown · templates"
    body: "A minimal static site generator that stays out of your way. Edit content in YAML and Markdown. Build with Go. Deploy anywhere."
    buttons:
      - text: "Get Started"
        href: "/docs"
        variant: "primary"
      - text: "View on GitHub"
        href: "https://github.com/b7r-dev/lyt"
        variant: "outline"

  - id: "principles"
    type: "features"
    title: "Principles"
    cards:
      - title: "Separation"
        body: "Engine, content, markup, styling — each in its place. Never the twain shall meet."
      - title: "Simplicity"
        body: "YAML for structure. Markdown for prose. Go for everything else. No runtime JS, no build step for templates."
      - title: "Depth"
        body: "3 z-planes give depth without weight. Shadows and blur, rendered with CSS alone."
      - title: "Beige"
        body: "Warm, desaturated, depth-present. A palette that recedes so your content advances."

  - id: "components"
    type: "default"
    title: "Components"
    body: |
      lyt includes composable components for rich content:

  - id: "cta"
    type: "cta"
    title: "Try lyt today"
    body: A minimal SSG that gets out of your way. Zero dependencies, pure Go.
    button_text: "Read the Docs"
    button_href: "/docs"
`

const defaultAboutYAML = `# About page

meta:
  title: "About"
  slug: "/about"
  description: "About this site"
  agent: true

agent_content:
  title: "About lyt"
  description: "Learn about lyt's philosophy and design"
  
  sections:
    - type: "default"
      title: "About This Site"
      body: |
        This page was created by lyt init. Edit content/pages/about.yaml to customize it.

    - type: "cli"
      title: "Customize Content"
      commands: |
        # Edit the about page
        vim content/pages/about.yaml
        lyt build

    - type: "link"
      text: "← Back to Home"
      href: "/"

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

const defaultBlogIndexYAML = `# Blog index

meta:
  title: "Blog"
  slug: "/blog"
  description: "Writing on tools, craft, and the web"

sections:
  - id: "intro"
    type: "hero"
    title: "Blog"
    body: |
      Writing on tools, craft, and the web.

  - id: "posts"
    type: "default"
    title: "Posts"
    cards: []
`

const defaultDocsYAML = `# Documentation

meta:
  title: "Documentation"
  slug: "/docs"
  description: "Complete lyt documentation — static sites without the ceremony"

sections:
  - id: "intro"
    type: "hero"
    title: "Documentation"
    subtitle: "lyt"
    body: "Everything you need to build a site with lyt. From installation to deployment."

  - id: "overview"
    type: "default"
    title: "Overview"
    body: |
      lyt is a minimal static site generator. Write content in YAML, render to HTML, deploy anywhere.
      
      **Why lyt?**
      
      - Zero runtime JS — Pure HTML and CSS output
      - YAML + Markdown — Structured content, easy versioning
      - Go-powered — Sub-second builds, single binary
      - Components — Rich content without shortcodes
      
      [Getting Started →](/docs/getting-started)

  - id: "content"
    type: "default"
    title: "Content"
    body: |
      Every page is a YAML file. Define structure with sections, render rich content with components.
      
      Learn about: frontmatter fields, section types, component library, markdown integration.

  - id: "configuration"
    type: "default"
    title: "Configuration"
    body: |
      Customize your site with design tokens and site configuration.
      
      Control: colors, typography, spacing, navigation, site metadata.

  - id: "deployment"
    type: "default"
    title: "Deployment"
    body: |
      lyt outputs pure static files. Deploy to any hosting provider.
      
      Supports: Netlify, Vercel, Cloudflare Pages, GitHub Pages, rsync/VPS.
`

const defaultTokensYAML = `# lyt design tokens
# beige, depth-present, un-opinionated

colors:
  # Base palette - warm, desaturated
  base:
    bg: "#f5f5f0"         # Warm off-white (primary background)
    surface: "#ecece6"     # Slightly darker beige (cards, pre blocks)
    border: "#d8d8d0"      # Muted warm gray
    muted: "#7a7a6e"      # Secondary text
    text: "#3d3d3d"       # Primary text (softened black)
    heading: "#2a2a25"    # Dark headings
    link: "#5c5c4a"       # Links
    link-hover: "#3d3d3d" # Link hover

  # Accent: Warm amber - works in both light and dark mode
  accent:
    1: "#f5e6c8"  # Subtle - backgrounds, hover states
    2: "#d4b896"  # Muted - borders, secondary elements
    3: "#c4a77d"  # Default - buttons, links, highlights
    4: "#a08050"  # Strong - active states, emphasis
    5: "#6b552f"  # Intense - focus rings, important CTAs

  # System
  system:
    error: "#8b3a3a"
    success: "#3a6b3a"
    warning: "#8b6a30"

spacing:
  0: "0"
  1: "0.25rem"
  2: "0.5rem"
  3: "0.75rem"
  4: "1rem"
  5: "1.25rem"
  6: "1.5rem"
  8: "2rem"
  10: "2.5rem"
  12: "3rem"
  16: "4rem"
  20: "5rem"
  24: "6rem"
  32: "8rem"

typography:
  font_family:
    body: "system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif"
    heading: "Georgia, 'Times New Roman', serif"
    brand: "Georgia, 'Times New Roman', serif"
    mono: "'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', monospace"

  font_size:
    xs: "0.75rem"
    sm: "0.8125rem"
    base: "1.0625rem"
    lg: "1.1875rem"
    xl: "1.375rem"
    2xl: "1.625rem"
    3xl: "2rem"
    4xl: "2.5rem"
    5xl: "3rem"
    6xl: "3.75rem"

  font_weight:
    regular: "400"
    medium: "500"
    semibold: "600"
    bold: "700"

  line_height:
    none: "1"
    tight: "1.2"
    snug: "1.3"
    normal: "1.45"
    relaxed: "1.55"
    loose: "1.7"

  letter_spacing:
    tight: "-0.015em"
    normal: "0"
    wide: "0.015em"
    wider: "0.025em"
    widest: "0.05em"

# Z-plane layers: 8 total
# 0: ground (default)
# 1-4: userland (content components)
# 5-8: system (nav, tooltips, modals)
z:
  # Userland (1-4)
  content: "1"       # Default user content
  elevated: "2"     # Elevated user content
  floating: "3"      # Floating user content (buttons, CTAs)
  top: "4"          # Highest userland

  # System (5-8)
  overlay: "5"      # Overlays
  nav: "6"           # Navigation (sticky header)
  tooltip: "7"       # Tooltips
  modal: "8"         # Modals, dialogs

# Layout
layout:
  max_width: "720px"
  content_width: "65ch"
`

const defaultSiteYAML = `# Site config

meta:
  title: "My lyt Site"
  description: "Built with lyt"
  url: "https://example.com"
  copyright: "Copyright © 2026"
  license: "MIT"

nav:
  - label: "Home"
    href: "/"
  - label: "Blog"
    href: "/blog"
  - label: "Docs"
    href: "/docs"
  - label: "Components"
    href: "/components"
  - label: "About"
    href: "/about"

# Agent path prefix (pages opt-in via meta.agent: true)
agent_section:
  path: "/agents"
`

const defaultCSS = `/* lyt base stylesheet */
/* Minimal, un-opinionated, depth-present, beige */
/* Mobile-first, 3-plane z-index metaphor */

/* ─── Reset ─────────────────────────────────────────── */
*, *::before, *::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

html {
  font-size: 16px;
  scroll-behavior: smooth;
  -webkit-text-size-adjust: 100%;
  color-scheme: light dark;
  /* Background on html prevents white flash before body loads */
  background-color: var(--color-bg, #f5f5f0);
  background-image: 
    radial-gradient(ellipse 60% 30% at 50% -10%, rgba(220, 180, 120, 0.06) 0%, transparent 60%),
    linear-gradient(180deg, var(--color-bg, #f5f5f0) 0%, var(--color-bg, #f5f5f0) 100%);
  background-attachment: fixed;
}

body {
  font-family: var(--font-body, system-ui, sans-serif);
  font-size: var(--text-base, 1rem);
  line-height: var(--leading-normal, 1.5);
  color: var(--color-text, #3d3d3d);
  background-color: var(--color-bg, #f5f5f0);
  /* Universal depth gradient - warm glow works in both light and dark */
  background-image: 
    radial-gradient(ellipse 60% 30% at 50% -10%, rgba(220, 180, 120, 0.06) 0%, transparent 60%),
    linear-gradient(180deg, var(--color-bg, #f5f5f0) 0%, var(--color-bg, #f5f5f0) 100%);
  background-attachment: fixed;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
}

/* Accent color: warm amber - visible in both light and dark mode */
:root {
  --color-accent-1: #f5e6c8;
  --color-accent-2: #d4b896;
  --color-accent-3: #c4a77d;
  --color-accent-4: #a08050;
  --color-accent-5: #6b552f;
}

img, video, svg {
  max-width: 100%;
  height: auto;
  display: block;
}

a {
  color: var(--color-link, #5c5c4a);
  text-decoration: underline;
  text-underline-offset: 2px;
}

a:hover {
  color: var(--color-link-hover, #3d3d3d);
}

/* ─── Layout ───────────────────────────────────────── */
.main {
  flex: 1;
  width: 100%;
  max-width: var(--max-width, 720px);
  margin-inline: auto;
  padding: var(--space-4, 1rem);
  padding-bottom: var(--space-16, 4rem);
}

/* ─── Nav ───────────────────────────────────────────── */
.nav {
  position: sticky;
  top: 0;
  z-index: var(--z-nav, 6);
  background: color-mix(in srgb, var(--color-bg) 85%, transparent);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--color-border, #d8d8d0);
  padding: var(--space-3, 0.75rem) var(--space-4, 1rem);
}

.nav-inner {
  max-width: var(--max-width, 720px);
  margin-inline: auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4, 1rem);
}

.nav-brand {
  font-family: var(--font-brand, serif);
  font-size: var(--text-lg, 1.125rem);
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--color-text) !important;
  text-decoration: none;
}

.nav-github {
  margin-left: auto;
  color: var(--color-text) !important;
  text-decoration: none;
  display: flex;
  align-items: center;
  transition: opacity 0.2s ease;
}

.nav-github:hover {
  opacity: 0.7;
}

/* Mobile nav */
@media (max-width: 640px) {
  .nav-github {
    display: none;
  }
}

.nav-links {
  display: flex;
  list-style: none;
  gap: var(--space-4, 1rem);
  flex-wrap: wrap;
}

.nav-link {
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-muted, #7a7a6e) !important;
  text-decoration: none;
  transition: color 0.15s ease;
  position: relative;
}

.nav-link:hover {
  color: var(--color-text) !important;
}

.nav-link.active,
.nav-link[aria-current="page"] {
  color: var(--color-accent-4, #a08050) !important;
}

/* Active indicator: small dot below the link */
.nav-link.active::after,
.nav-link[aria-current="page"]::after {
  content: "";
  position: absolute;
  bottom: -4px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--color-accent-3, #c4a77d);
}

/* Mobile menu toggle button (hidden on desktop) */
.nav-toggle {
  display: none;
  position: relative;
  z-index: calc(var(--z-nav, 6) + 1);
  width: 32px;
  height: 32px;
  padding: 4px;
  background: none;
  border: none;
  cursor: pointer;
  appearance: none;
  -webkit-appearance: none;
}

.nav-toggle span,
.nav-toggle span::before,
.nav-toggle span::after {
  display: block;
  width: 24px;
  height: 2px;
  background: var(--color-text, #3d3d3d);
  border-radius: 1px;
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.nav-toggle span {
  position: relative;
}

.nav-toggle span::before,
.nav-toggle span::after {
  content: '';
  position: absolute;
  left: 0;
}

.nav-toggle span::before {
  top: -7px;
}

.nav-toggle span::after {
  top: 7px;
}

/* Mobile menu - checkbox hack for JS-free toggle */
.nav-menu-checkbox {
  position: absolute;
  opacity: 0;
  pointer-events: none;
}

/* Mobile menu container */
.nav-menu {
  display: none;
}

/* ─── Footer ────────────────────────────────────────── */
.footer {
  padding: var(--space-8, 2rem) var(--space-4, 1rem);
  text-align: center;
  font-size: var(--text-xs, 0.75rem);
  color: var(--color-muted, #7a7a6e);
  border-top: 1px solid var(--color-border, #d8d8d0);
}

/* Agent link - subtle but discoverable */
.agent-notice {
  margin-top: var(--space-3, 0.75rem);
}

.footer-copyright,
.footer-license {
  font-size: var(--text-xs, 0.75rem);
  color: var(--color-muted, #7a7a6e);
  margin: var(--space-1, 0.25rem) 0;
}

.agent-link {
  color: var(--color-muted, #7a7a6e);
  font-size: var(--text-xs, 0.75rem);
  text-decoration: none;
  opacity: 0.7;
  transition: opacity 0.15s ease;
}

.agent-link:hover {
  opacity: 1;
  text-decoration: underline;
}

/* ─── Sections ──────────────────────────────────────── */
.section {
  padding: var(--space-8, 2rem) 0;
}

.section-title {
  font-family: var(--font-heading, serif);
  font-size: var(--text-2xl, 1.625rem);
  font-weight: 700;
  letter-spacing: -0.01em;
  color: var(--color-heading, #2a2a25);
  margin-bottom: var(--space-3, 0.75rem);
  line-height: 1.15;
}

.section-content {
  font-size: var(--text-base, 1.0625rem);
  line-height: 1.55;
  color: var(--color-text, #3d3d3d);
}

/* ─── Section Break ──────────────────────────────────── */
.section-break {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-8, 2rem) 0;
  margin: var(--space-8, 2rem) 0;
}

.section-break::before,
.section-break::after {
  content: "";
  flex: 1;
  height: 1px;
  background: var(--color-border, #d8d8d0);
}

.section-break-icon {
  padding: 0 var(--space-4, 1rem);
  color: var(--color-muted, #7a7a6e);
  font-size: var(--text-lg, 1.125rem);
}

.section-content p {
  margin-bottom: var(--space-4, 1rem);
}

.section-content h2 {
  font-family: var(--font-heading, serif);
  font-size: var(--text-xl, 1.375rem);
  font-weight: 700;
  margin-top: var(--space-8, 2rem);
  margin-bottom: var(--space-3, 0.75rem);
  letter-spacing: -0.01em;
  color: var(--color-heading, #2a2a25);
}

.section-content h3 {
  font-family: var(--font-heading, serif);
  font-size: var(--text-lg, 1.1875rem);
  font-weight: 600;
  margin-top: var(--space-6, 1.5rem);
  margin-bottom: var(--space-2, 0.5rem);
  letter-spacing: -0.01em;
  color: var(--color-heading, #2a2a25);
}

.section-content ul, .section-content ol {
  margin-bottom: var(--space-4, 1rem);
  padding-left: var(--space-6, 1.5rem);
}

.section-content li {
  margin-bottom: var(--space-2, 0.5rem);
}

/* ─── Tables ────────────────────────────────────────── */
.section-content table {
  width: 100%;
  border-collapse: collapse;
  margin: var(--space-4, 1rem) 0;
  font-size: var(--text-sm, 0.875rem);
}

.section-content th,
.section-content td {
  padding: var(--space-3, 0.75rem) var(--space-4, 1rem);
  text-align: left;
  border-bottom: 1px solid var(--color-border, #d8d8d0);
}

.section-content th {
  font-family: var(--font-mono, monospace);
  font-weight: 600;
  background: var(--color-surface, #ecece6);
  color: var(--color-heading, #2a2a25);
}

.section-content td {
  background: var(--color-bg, #f5f5f0);
}

.section-content tr:last-child td {
  border-bottom: none;
}

.section-content code {
  font-family: var(--font-mono, monospace);
  font-size: 0.875em;
  background: var(--color-surface, #ecece6);
  padding: 0.125em 0.375em;
  border-radius: 3px;
}

.section-content pre {
  background: var(--color-surface, #ecece6);
  padding: var(--space-4, 1rem);
  border-radius: 4px;
  overflow-x: auto;
  margin-bottom: var(--space-4, 1rem);
}

.section-content pre code {
  background: none;
  padding: 0;
}

.section-content blockquote {
  border-left: 3px solid var(--color-border, #d8d8d0);
  padding-left: var(--space-4, 1rem);
  color: var(--color-muted, #7a7a6e);
  font-style: italic;
  margin-bottom: var(--space-4, 1rem);
}

/* ─── Cards ─────────────────────────────────────────── */
.cards {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--space-4, 1rem);
  margin-top: var(--space-6, 1.5rem);
}

@media (min-width: 640px) {
  .cards {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .cards {
    grid-template-columns: repeat(3, 1fr);
  }
}

.card {
  position: relative;
  z-index: var(--z-content, 1);
  background: var(--color-surface, #ecece6);
  border: 1px solid var(--color-border, #d8d8d0);
  padding: var(--space-5, 1.25rem);
  border-radius: 4px;
  box-shadow:
    2px 2px 0 0 var(--color-border, #d8d8d0),
    4px 4px 0 0 color-mix(in srgb, var(--color-border) 30%, transparent);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.card:hover {
  transform: translate(-1px, -1px);
  box-shadow:
    3px 3px 0 0 var(--color-border, #d8d8d0),
    6px 6px 0 0 color-mix(in srgb, var(--color-border) 30%, transparent);
}

.card-title {
  font-family: var(--font-heading, monospace);
  font-size: var(--text-base, 1rem);
  font-weight: 700;
  margin-bottom: var(--space-2, 0.5rem);
  color: var(--color-heading, #2a2a25);
}

.card-body {
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-text, #3d3d3d);
  line-height: 1.5;
}

/* ─── Buttons ───────────────────────────────────────── */
.btn {
  display: inline-block;
  position: relative;
  z-index: var(--z-floating, 3);
  padding: var(--space-2, 0.5rem) var(--space-4, 1rem);
  font-family: var(--font-mono, monospace);
  font-size: var(--text-sm, 0.875rem);
  font-weight: 600;
  letter-spacing: 0.02em;
  text-decoration: none !important;
  border: 2px solid var(--color-border, #3d3d3d);
  border-radius: 2px;
  background: var(--color-bg, #f5f5f0);
  color: var(--color-text, #3d3d3d);
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease, transform 0.1s ease;
}

.btn:hover {
  background: var(--color-text, #3d3d3d);
  color: var(--color-bg, #f5f5f0);
}

.btn:active {
  transform: translate(1px, 1px);
}

.btn-primary {
  background: var(--color-text, #3d3d3d);
  color: var(--color-bg, #f5f5f0);
}

.btn-primary:hover {
  background: var(--color-heading, #2a2a25);
  color: var(--color-bg, #f5f5f0);
}

.btn-outline {
  background: transparent;
  border-color: var(--color-text, #3d3d3d);
}

.btn-outline:hover {
  background: var(--color-text, #3d3d3d);
  color: var(--color-bg, #f5f5f0);
}

/* ─── CTA Group ─────────────────────────────────────── */
.cta-group {
  display: flex;
  gap: var(--space-3, 0.75rem);
  flex-wrap: wrap;
  margin-top: var(--space-6, 1.5rem);
}

/* ─── Blog Post ─────────────────────────────────────── */
.blog-post {
  padding: var(--space-8, 2rem) 0;
}

.post-title {
  font-family: var(--font-heading, monospace);
  font-size: clamp(var(--text-2xl, 1.5rem), 5vw, var(--text-4xl, 2.25rem));
  font-weight: 800;
  letter-spacing: -0.03em;
  line-height: 1.1;
  color: var(--color-heading, #2a2a25);
  margin-bottom: var(--space-4, 1rem);
}

.post-meta {
  font-family: var(--font-mono, monospace);
  font-size: var(--text-xs, 0.75rem);
  color: var(--color-muted, #7a7a6e);
  margin-bottom: var(--space-6, 1.5rem);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.post-description {
  font-size: var(--text-lg, 1.125rem);
  color: var(--color-text, #3d3d3d);
  line-height: var(--leading-relaxed, 1.625);
  margin-bottom: var(--space-8, 2rem);
  padding-bottom: var(--space-6, 1.5rem);
  border-bottom: 1px solid var(--color-border, #d8d8d0);
}

.post-body {
  font-size: var(--text-base, 1rem);
  line-height: var(--leading-relaxed, 1.625);
}

.post-body p {
  margin-bottom: var(--space-4, 1rem);
}

.post-body h2 {
  font-family: var(--font-heading, monospace);
  font-size: var(--text-xl, 1.25rem);
  font-weight: 700;
  margin-top: var(--space-8, 2rem);
  margin-bottom: var(--space-3, 0.75rem);
}

.post-body h3 {
  font-family: var(--font-heading, monospace);
  font-size: var(--text-lg, 1.125rem);
  font-weight: 600;
  margin-top: var(--space-6, 1.5rem);
  margin-bottom: var(--space-2, 0.5rem);
}

.post-body ul, .post-body ol {
  margin-bottom: var(--space-4, 1rem);
  padding-left: var(--space-6, 1.5rem);
}

.post-body li {
  margin-bottom: var(--space-2, 0.5rem);
}

.post-body a {
  color: var(--color-link, #5c5c4a);
  text-decoration: underline;
  text-underline-offset: 2px;
}

.post-body blockquote {
  border-left: 3px solid var(--color-border, #d8d8d0);
  padding-left: var(--space-4, 1rem);
  color: var(--color-muted, #7a7a6e);
  font-style: italic;
  margin-bottom: var(--space-4, 1rem);
}

.post-body code {
  font-family: var(--font-mono, monospace);
  font-size: 0.875em;
  background: var(--color-surface, #ecece6);
  padding: 0.125em 0.375em;
  border-radius: 3px;
}

.post-body pre {
  background: var(--color-surface, #ecece6);
  padding: var(--space-4, 1rem);
  border-radius: 4px;
  overflow-x: auto;
  margin-bottom: var(--space-4, 1rem);
}

.post-body pre code {
  background: none;
  padding: 0;
}

.post-body img {
  border-radius: 4px;
  margin: var(--space-6, 1.5rem) 0;
}

/* ─── Hero ──────────────────────────────────────────── */
.hero {
  padding: var(--space-16, 4rem) 0 var(--space-12, 3rem);
}

.hero-title {
  font-family: var(--font-heading, serif);
  font-size: clamp(var(--text-3xl, 2rem), 8vw, var(--text-6xl, 3.75rem));
  font-weight: 700;
  letter-spacing: -0.02em;
  line-height: 1.1;
  color: var(--color-heading, #2a2a25);
  margin-bottom: var(--space-3, 0.75rem);
}

.hero-subtitle {
  font-family: var(--font-body, sans-serif);
  font-size: var(--text-base, 1.0625rem);
  font-weight: 400;
  letter-spacing: 0.02em;
  color: var(--color-muted, #7a7a6e);
  margin-bottom: var(--space-5, 1.25rem);
}

.hero-body {
  font-size: var(--text-lg, 1.1875rem);
  line-height: 1.55;
  color: var(--color-text, #3d3d3d);
  max-width: 55ch;
  margin-bottom: var(--space-5, 1.25rem);
}

.hero-buttons {
  display: flex;
  gap: var(--space-3, 0.75rem);
  flex-wrap: wrap;
  margin-top: var(--space-6, 1.5rem);
}

/* ─── Features Grid ───────────────────────────────────── */
.features-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--space-4, 1rem);
  margin-top: var(--space-6, 1.5rem);
}

@media (min-width: 640px) {
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1024px) {
  .features-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

.feature-card {
  position: relative;
  z-index: var(--z-content, 1);
  background: var(--color-surface, #ecece6);
  border: 1px solid var(--color-border, #d8d8d0);
  padding: var(--space-5, 1.25rem);
  border-radius: 4px;
  box-shadow:
    2px 2px 0 0 var(--color-border, #d8d8d0),
    4px 4px 0 0 color-mix(in srgb, var(--color-border) 30%, transparent);
}

.feature-icon {
  font-size: var(--text-2xl, 1.5rem);
  margin-bottom: var(--space-3, 0.75rem);
}

.feature-title {
  font-family: var(--font-heading, monospace);
  font-size: var(--text-base, 1rem);
  font-weight: 700;
  margin-bottom: var(--space-2, 0.5rem);
  color: var(--color-heading, #2a2a25);
}

.feature-body {
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-text, #3d3d3d);
  line-height: 1.5;
}

/* ─── Callout ────────────────────────────────────────── */
.callout {
  position: relative;
  z-index: var(--z-content, 1);
  padding: var(--space-4, 1rem);
  margin: var(--space-6, 1.5rem) 0;
  border-radius: 4px;
  border-left: 4px solid var(--color-muted, #7a7a6e);
  background: var(--color-surface, #ecece6);
}

.callout-info {
  border-color: var(--color-muted, #7a7a6e);
}

.callout-tip {
  border-color: var(--color-accent-cool, #8090a0);
  background: color-mix(in srgb, var(--color-accent-cool, #8090a0) 10%, var(--color-surface, #ecece6));
}

.callout-note {
  border-color: var(--color-accent-warm, #a89880);
  background: color-mix(in srgb, var(--color-accent-warm, #a89880) 10%, var(--color-surface, #ecece6));
}

.callout-title {
  display: block;
  font-family: var(--font-heading, monospace);
  font-size: var(--text-sm, 0.875rem);
  font-weight: 700;
  margin-bottom: var(--space-2, 0.5rem);
  color: var(--color-heading, #2a2a25);
}

.callout-body {
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-text, #3d3d3d);
  line-height: 1.5;
}

.callout-body p:last-child {
  margin-bottom: 0;
}

/* ─── Pull Quote ──────────────────────────────────────── */
.pull-quote {
  position: relative;
  z-index: var(--z-content, 1);
  margin: var(--space-8, 2rem) 0;
  padding: var(--space-6, 1.5rem);
  padding-left: var(--space-8, 2rem);
  border-left: 4px solid var(--color-accent-warm, #a89880);
  background: var(--color-surface, #ecece6);
  border-radius: 0 4px 4px 0;
}

.pull-quote p {
  font-family: var(--font-body, system-ui, sans-serif);
  font-size: var(--text-xl, 1.25rem);
  font-style: italic;
  line-height: var(--leading-relaxed, 1.625);
  color: var(--color-heading, #2a2a25);
  margin-bottom: var(--space-4, 1rem);
}

.pull-quote-attribution {
  display: block;
  font-family: var(--font-mono, monospace);
  font-size: var(--text-xs, 0.75rem);
  color: var(--color-muted, #7a7a6e);
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.pull-quote-attribution cite {
  font-style: normal;
  font-weight: 600;
}

/* ─── Citation ───────────────────────────────────────── */
.citation {
  display: block;
  position: relative;
  z-index: var(--z-content, 1);
  margin: var(--space-6, 1.5rem) 0;
  padding: var(--space-4, 1rem);
  background: var(--color-surface, #ecece6);
  border-left: 3px solid var(--color-border, #d8d8d0);
  border-radius: 0 4px 4px 0;
  font-style: normal;
}

.citation-link {
  font-weight: 600;
  color: var(--color-link, #5c5c4a);
  text-decoration: none;
}

.citation-link:hover {
  color: var(--color-link-hover, #3d3d3d);
  text-decoration: underline;
}

.citation-details {
  display: block;
  margin-top: var(--space-2, 0.5rem);
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-muted, #7a7a6e);
}

/* ─── CTA Section ────────────────────────────────────── */
.section-cta {
  position: relative;
  z-index: var(--z-floating, 3);
  text-align: center;
  padding: var(--space-12, 3rem) var(--space-6, 1.5rem);
  margin: var(--space-8, 2rem) 0;
  background: var(--color-surface, #ecece6);
  border: 1px solid var(--color-border, #d8d8d0);
  border-radius: 8px;
  box-shadow:
    4px 4px 0 0 var(--color-border, #d8d8d0),
    8px 8px 0 0 color-mix(in srgb, var(--color-border) 20%, transparent);
}

.cta-title {
  font-family: var(--font-heading, monospace);
  font-size: var(--text-2xl, 1.5rem);
  font-weight: 800;
  letter-spacing: -0.02em;
  color: var(--color-heading, #2a2a25);
  margin-bottom: var(--space-4, 1rem);
}

.cta-body {
  font-size: var(--text-base, 1rem);
  color: var(--color-text, #3d3d3d);
  max-width: 45ch;
  margin: 0 auto var(--space-6, 1.5rem);
  line-height: var(--leading-relaxed, 1.625);
}

/* ─── Warning Box ────────────────────────────────────── */
.warning-box {
  position: relative;
  z-index: var(--z-content, 1);
  padding: var(--space-4, 1rem);
  margin: var(--space-6, 1.5rem) 0;
  border-radius: 4px;
  border-left: 4px solid var(--color-system-warning, #8b6a30);
  background: color-mix(in srgb, var(--color-system-warning, #8b6a30) 10%, var(--color-bg, #f5f5f0));
}

.warning-warning {
  border-color: var(--color-system-warning, #8b6a30);
  background: color-mix(in srgb, var(--color-system-warning, #8b6a30) 10%, var(--color-bg, #f5f5f0));
  box-shadow:
    2px 2px 0 0 var(--color-system-warning, #8b6a30),
    4px 4px 0 0 color-mix(in srgb, var(--color-system-warning, #8b6a30) 30%, transparent);
}

.warning-error {
  border-color: var(--color-system-error, #8b3a3a);
  background: color-mix(in srgb, var(--color-system-error, #8b3a3a) 10%, var(--color-bg, #f5f5f0));
  box-shadow:
    2px 2px 0 0 var(--color-system-error, #8b3a3a),
    4px 4px 0 0 color-mix(in srgb, var(--color-system-error, #8b3a3a) 30%, transparent);
}

.warning-info {
  border-color: var(--color-accent-cool, #8090a0);
  background: color-mix(in srgb, var(--color-accent-cool, #8090a0) 10%, var(--color-bg, #f5f5f0));
  box-shadow:
    2px 2px 0 0 var(--color-accent-cool, #8090a0),
    4px 4px 0 0 color-mix(in srgb, var(--color-accent-cool, #8090a0) 30%, transparent);
}

.warning-success {
  border-color: var(--color-system-success, #3a6b3a);
  background: color-mix(in srgb, var(--color-system-success, #3a6b3a) 10%, var(--color-bg, #f5f5f0));
  box-shadow:
    2px 2px 0 0 var(--color-system-success, #3a6b3a),
    4px 4px 0 0 color-mix(in srgb, var(--color-system-success, #3a6b3a) 30%, transparent);
}

.warning-title {
  display: block;
  font-family: var(--font-heading, monospace);
  font-size: var(--text-sm, 0.875rem);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: var(--space-2, 0.5rem);
  color: var(--color-heading, #2a2a25);
}

.warning-body {
  font-size: var(--text-sm, 0.875rem);
  color: var(--color-text, #3d3d3d);
  line-height: 1.5;
}

.warning-body p:last-child {
  margin-bottom: 0;
}

/* ─── Code Example ─────────────────────────────────────── */
.code-example {
  position: relative;
  z-index: var(--z-content, 1);
  margin: var(--space-6, 1.5rem) 0;
  border-radius: 4px;
  overflow: hidden;
  border: 1px solid var(--color-border, #d8d8d0);
  background: var(--color-surface, #ecece6);
}

.code-example-header {
  display: flex;
  align-items: center;
  gap: var(--space-3, 0.75rem);
  padding: var(--space-2, 0.5rem) var(--space-4, 1rem);
  background: var(--color-border, #d8d8d0);
  border-bottom: 1px solid var(--color-border, #d8d8d0);
}

.code-example-title {
  font-family: var(--font-mono, monospace);
  font-size: var(--text-xs, 0.75rem);
  color: var(--color-text, #3d3d3d);
}

.code-example-lang {
  font-family: var(--font-mono, monospace);
  font-size: var(--text-xs, 0.75rem);
  font-weight: 600;
  color: var(--color-muted, #7a7a6e);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.code-block {
  margin: 0;
  padding: var(--space-4, 1rem);
  overflow-x: auto;
}

.code-block code {
  font-family: var(--font-mono, monospace);
  font-size: var(--text-sm, 0.875rem);
  line-height: 1.6;
  color: var(--color-text, #3d3d3d);
  background: none;
  padding: 0;
}

/* ─── Utility ───────────────────────────────────────── */
.error {
  color: var(--color-error, #8b3a3a);
  padding: var(--space-4, 1rem);
  background: color-mix(in srgb, var(--color-error) 10%, transparent);
  border: 1px solid var(--color-error, #8b3a3a);
  border-radius: 2px;
}

/* ─── Responsive ────────────────────────────────────── */
@media (max-width: 639px) {
  /* Show hamburger menu on mobile */
  .nav-toggle {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  /* Mobile menu - hidden by default */
  .nav-menu {
    display: none;
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: var(--color-bg, #f5f5f0);
    border-bottom: 1px solid var(--color-border, #d8d8d0);
    padding: var(--space-4, 1rem);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  /* When checkbox is checked, show menu */
  .nav-menu-checkbox:checked ~ .nav-menu {
    display: block;
    animation: slideDown 0.2s ease;
  }

  /* Animate hamburger to X */
  .nav-menu-checkbox:checked ~ .nav-toggle span {
    background: transparent;
  }

  .nav-menu-checkbox:checked ~ .nav-toggle span::before {
    transform: translateY(7px) rotate(45deg);
  }

  .nav-menu-checkbox:checked ~ .nav-toggle span::after {
    transform: translateY(-7px) rotate(-45deg);
  }

  /* Hide desktop links on mobile */
  .nav-links {
    display: none;
  }

  /* Mobile links - shown in menu */
  .nav-mobile-links {
    display: flex;
    flex-direction: column;
    gap: var(--space-3, 0.75rem);
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .nav-mobile-link {
    display: block;
    padding: var(--space-2, 0.5rem) var(--space-3, 0.75rem);
    font-size: var(--text-base, 1rem);
    color: var(--color-text, #3d3d3d) !important;
    text-decoration: none;
    border-radius: 4px;
    transition: background 0.15s ease;
  }

  .nav-mobile-link:hover {
    background: var(--color-surface, #ecece6);
    color: var(--color-text, #3d3d3d) !important;
  }

  .nav-mobile-link.active,
  .nav-mobile-link[aria-current="page"] {
    background: var(--color-surface, #ecece6);
    font-weight: 600;
  }

  .hero-title {
    font-size: var(--text-4xl, 2.25rem);
  }

  @keyframes slideDown {
    from {
      opacity: 0;
      transform: translateY(-8px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
}

@media (min-width: 640px) {
  .main {
    padding: var(--space-8, 2rem);
  }

  .hero {
    padding: var(--space-20, 5rem) 0 var(--space-16, 4rem);
  }
}

@media (min-width: 1024px) {
  .main {
    padding: var(--space-12, 3rem);
  }
}

/* ─── Dark Mode ─────────────────────────────────────── */
@media (prefers-color-scheme: dark) {
  :root {
    /* Dark mode: Warm night sky with fireflies */
    --color-bg: #111215;
    --color-surface: #16181f;
    --color-surface-elevated: #1e2028;
    --color-border: #2a2e38;

    /* Universal gradient works in both modes */
    
    --color-text: #e4dfd4;
    --color-text-secondary: #9a9488;
    --color-muted: #5e5b54;
    --color-heading: #ece6d8;
    
    --color-link: #d4b896;
    --color-link-hover: #f0dec2;
    
    /* Accent: warm amber - firefly glow */
    --color-accent-1: #2d2418;
    --color-accent-2: #453828;
    --color-accent-3: #6b5230;
    --color-accent-4: #a08050;
    --color-accent-5: #e4c896;
    
    --color-system-warning: #d4b896;
    --color-system-error: #d98a8a;
    --color-system-success: #8ab87a;
    
    --color-error: #d98a8a;
  }
  
  /* Adjustments for dark mode */
  .nav {
    background: color-mix(in srgb, var(--color-bg) 90%, transparent);
  }
  
  .nav-toggle span,
  .nav-toggle span::before,
  .nav-toggle span::after {
    background: var(--color-text, #e8e4dc);
  }
  
  .card,
  .feature-card,
  .callout,
  .pull-quote,
  .citation,
  .section-cta,
  .code-example {
    border-color: var(--color-border, #262e3a);
    box-shadow:
      2px 2px 0 0 #1a1c22,
      4px 4px 12px 0 rgba(0, 0, 0, 0.4);
  }
  
  .card:hover {
    box-shadow:
      3px 3px 0 0 #1a1c22,
      6px 6px 16px 0 rgba(0, 0, 0, 0.35);
  }
  
  .code-example-header {
    background: var(--color-border, #262e3a);
  }
  
  .section-cta {
    box-shadow:
      4px 4px 0 0 #1a1c22,
      8px 8px 20px 0 rgba(0, 0, 0, 0.35);
  }
  
  .nav-mobile-link:hover,
  .nav-mobile-link.active {
    background: var(--color-surface, #131920);
  }
  
  /* Warning box dark mode - consistent stacked shadows */
  .warning-warning {
    background: color-mix(in srgb, var(--color-system-warning) 12%, var(--color-bg));
    border-color: var(--color-system-warning);
    box-shadow:
      2px 2px 0 0 var(--color-system-warning),
      4px 4px 0 0 color-mix(in srgb, var(--color-system-warning) 30%, transparent);
  }
  
  .warning-error {
    background: color-mix(in srgb, var(--color-system-error) 12%, var(--color-bg));
    border-color: var(--color-system-error);
    box-shadow:
      2px 2px 0 0 var(--color-system-error),
      4px 4px 0 0 color-mix(in srgb, var(--color-system-error) 30%, transparent);
  }
  
  .warning-info {
    background: color-mix(in srgb, var(--color-accent-cool) 12%, var(--color-bg));
    border-color: var(--color-accent-cool);
    box-shadow:
      2px 2px 0 0 var(--color-accent-cool),
      4px 4px 0 0 color-mix(in srgb, var(--color-accent-cool) 30%, transparent);
  }
  
  .warning-success {
    background: color-mix(in srgb, var(--color-system-success) 12%, var(--color-bg));
    border-color: var(--color-system-success);
    box-shadow:
      2px 2px 0 0 var(--color-system-success),
      4px 4px 0 0 color-mix(in srgb, var(--color-system-success) 30%, transparent);
  }
}
`

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Force initialization even if directory is not empty")
}

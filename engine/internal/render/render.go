package render

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/b7r-dev/lyt/engine/internal/content"
	"github.com/b7r-dev/lyt/engine/internal/markdown"
)

// Renderer produces HTML from content
type Renderer struct {
	Collection *content.Collection
	contentDir string
	verbose    bool
}

// PageTemplate is the base HTML template
const PageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="description" content="{{.Description}}">
  <meta name="generator" content="lyt">
  <meta name="theme-color" content="#f5f5f0">
  
  <title>{{.Title}} | {{.SiteTitle}}</title>
  
  <!-- Open Graph -->
  <meta property="og:title" content="{{.Title}}">
  <meta property="og:description" content="{{.Description}}">
  <meta property="og:type" content="website">
  <meta property="og:url" content="{{.OGURL}}">
  
  <!-- Twitter Card -->
  <meta name="twitter:card" content="summary">
  <meta name="twitter:title" content="{{.Title}}">
  <meta name="twitter:description" content="{{.Description}}">
  
  <!-- Performance: preload critical CSS -->
  <link rel="preload" href="/base.css" as="style">
  <link rel="preload" href="/tokens.css" as="style">
  
  <!-- Stylesheets with defer -->
  <link rel="stylesheet" href="/base.css" media="print" onload="this.media='all'">
  <link rel="stylesheet" href="/tokens.css" media="print" onload="this.media='all'">
  <noscript>
    <link rel="stylesheet" href="/base.css">
    <link rel="stylesheet" href="/tokens.css">
  </noscript>
  
  <!-- Favicon -->
  <link rel="icon" href="/favicon.svg" type="image/svg+xml">
  
  <!-- SEO: canonical URL -->
  {{.Canonical}}
</head>
<body class="{{.BodyClass}}">
  {{.Nav}}
  <main class="main">
    {{.Content}}
  </main>
  <footer class="footer">
    <p>{{.SiteTitle}}</p>
    <p class="footer-copyright">{{.Copyright}}</p>
    <p class="footer-license">Licensed under {{.License}}</p>
    {{if .ShowAgentLink}}
    <p class="agent-notice">
      <a href="{{.AgentLink}}" class="agent-link">Agents Read This First →</a>
    </p>
    {{end}}
  </footer>
</body>
</html>`

// AgentPageTemplate is a minimal HTML template for LLM agents
// No nav, no hero, no footer - just pure content
const AgentPageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="generator" content="lyt">
  <title>{{.Title}} | {{.SiteTitle}} (Agent)</title>
</head>
<body>
  <main>
    {{.Content}}
  </main>
</body>
</html>`

func NewRenderer(c *content.Collection, contentDir string, verbose bool) *Renderer {
	return &Renderer{Collection: c, contentDir: contentDir, verbose: verbose}
}

type PageData struct {
	Title         string
	Description   string
	SiteTitle     string
	Content       htmltemplate.HTML
	Nav           htmltemplate.HTML
	BodyClass     string
	Canonical     htmltemplate.HTML
	OGURL         string
	ShowAgentLink bool
	AgentLink     string
	Copyright     string
	License       string
}

func (r *Renderer) RenderPage(cf content.ContentFile) (string, error) {
	meta := content.GetMeta(cf)

	title := getString(meta, "title", "Page")
	desc := getString(meta, "description", "")
	slug := getString(meta, "slug", "/")

	bodyClass := "page page-" + slugToClass(slug)

	// Build content from sections
	contentHTML, err := r.renderSections(cf)
	if err != nil {
		contentHTML = fmt.Sprintf(`<p class="error">render error: %v</p>`, err)
	}

	// Nav
	navHTML := r.renderNav(slug)

	// Canonical URL
	canonical := r.renderCanonical(slug)
	ogURL := "https://lyt.b7r.dev" + slug

	// Assemble
	data := PageData{
		Title:         title,
		Description:   desc,
		SiteTitle:     r.Collection.GetSiteTitle(),
		Content:       htmltemplate.HTML(contentHTML),
		Nav:           htmltemplate.HTML(navHTML),
		BodyClass:     bodyClass,
		Canonical:     htmltemplate.HTML(canonical),
		OGURL:         ogURL,
		ShowAgentLink: r.Collection.ShowAgentHubLink(cf),
		AgentLink:     r.getAgentLink(cf, slug),
		Copyright:     r.Collection.GetCopyright(),
		License:       r.Collection.GetLicense(),
	}

	return r.executeTemplate(PageTemplate, data)
}

// getAgentLink returns the appropriate agent link for a page
// If page has its own agent_content, link to /agents/{slug}
// Otherwise, link to the agent hub /agents
func (r *Renderer) getAgentLink(cf content.ContentFile, slug string) string {
	if r.Collection.HasAgentPage(cf) {
		return r.Collection.GetAgentPath() + slug
	}
	// Link to agent hub
	return r.Collection.GetAgentPath()
}

func (r *Renderer) RenderBlogPost(cf content.ContentFile) (string, error) {
	meta := content.GetMeta(cf)

	title := getString(meta, "title", "Post")
	desc := getString(meta, "description", "")
	slug := getString(meta, "slug", "")

	bodyClass := "page page-blog"

	// Check for sections first (new format)
	sectionsContent, _ := r.renderSections(cf)

	// Fallback to body if no sections
	bodyContent := ""
	if sectionsContent == "" {
		if cf.Markdown != nil {
			if ref, ok := cf.Markdown["@body.md"]; ok {
				html, _ := markdown.Render(ref)
				bodyContent = html
			} else if ref, ok := cf.Markdown["body.md"]; ok {
				html, _ := markdown.Render(ref)
				bodyContent = html
			}
		}

		// Fallback inline body
		if bodyContent == "" {
			if body, ok := cf.Data["body"].(string); ok {
				html, _ := markdown.Render(body)
				bodyContent = html
			}
		}
	}

	contentHTML := `<article class="blog-post">`
	contentHTML += fmt.Sprintf(`<h1 class="post-title">%s</h1>`, title)
	contentHTML += fmt.Sprintf(`<p class="post-meta">%s</p>`, getString(meta, "date", ""))
	if desc != "" {
		contentHTML += fmt.Sprintf(`<p class="post-description">%s</p>`, desc)
	}

	// Use sections if available, otherwise use body
	if sectionsContent != "" {
		contentHTML += sectionsContent
	} else {
		contentHTML += `<div class="post-body">` + bodyContent + `</div>`
	}
	contentHTML += `</article>`

	navHTML := r.renderNav("/blog/" + slug)

	// Canonical URL for blog posts
	canonical := r.renderCanonical("/blog/" + slug)
	ogURL := "https://lyt.b7r.dev/blog/" + slug

	data := PageData{
		Title:         title,
		Description:   desc,
		SiteTitle:     r.Collection.GetSiteTitle(),
		Content:       htmltemplate.HTML(contentHTML),
		Nav:           htmltemplate.HTML(navHTML),
		BodyClass:     bodyClass,
		Canonical:     htmltemplate.HTML(canonical),
		OGURL:         ogURL,
		ShowAgentLink: r.Collection.ShowAgentHubLink(cf),
		AgentLink:     r.getAgentLink(cf, "/blog/"+slug),
		Copyright:     r.Collection.GetCopyright(),
		License:       r.Collection.GetLicense(),
	}

	return r.executeTemplate(PageTemplate, data)
}

// RenderAgentPage renders a page for LLM agents with custom agent content
func (r *Renderer) RenderAgentPage(cf content.ContentFile) (string, error) {
	meta := content.GetMeta(cf)

	title := getString(meta, "title", "Page")
	slug := getString(meta, "slug", "/")

	var contentHTML string

	// Check for custom agent_content section
	agentContent := content.GetAgentContent(cf)
	if agentContent != nil {
		// Render the custom agent content
		contentHTML = r.renderAgentContent(agentContent)
	} else {
		// Fall back to sections for blog posts
		if cf.Type == "blog" {
			blogHTML, _ := r.renderBlogContent(cf)
			contentHTML = blogHTML
		} else {
			// No agent content defined - render stripped human content
			sectionsHTML, _ := r.renderSections(cf)
			contentHTML = sectionsHTML
		}
	}

	// Add links to other agent pages
	agentNavHTML := r.renderAgentNav(slug)

	data := PageData{
		Title:         title,
		SiteTitle:     r.Collection.GetSiteTitle(),
		Content:       htmltemplate.HTML(contentHTML),
		Nav:           htmltemplate.HTML(agentNavHTML),
		BodyClass:     "agent",
		Canonical:     htmltemplate.HTML(""),
		OGURL:         "",
		ShowAgentLink: false,
		AgentLink:     "",
	}

	return r.executeTemplate(AgentPageTemplate, data)
}

// renderAgentContent renders custom agent_content section from YAML
func (r *Renderer) renderAgentContent(agentContent map[string]interface{}) string {
	var sb strings.Builder

	// Title
	if title, ok := agentContent["title"].(string); ok {
		sb.WriteString(fmt.Sprintf(`<h1>%s</h1>`, title))
	}

	// Description
	if desc, ok := agentContent["description"].(string); ok {
		sb.WriteString(fmt.Sprintf(`<p class="agent-description">%s</p>`, desc))
	}

	// Sections - like regular sections but for agents
	if sections, ok := agentContent["sections"].([]interface{}); ok {
		for _, s := range sections {
			if sec, ok := s.(map[string]interface{}); ok {
				sb.WriteString(r.renderAgentSection(sec))
			}
		}
	}

	return sb.String()
}

// renderAgentSection renders a section in agent content
func (r *Renderer) renderAgentSection(sec map[string]interface{}) string {
	var sb strings.Builder

	secType := getString(sec, "type", "default")

	switch secType {
	case "cli":
		sb.WriteString(r.renderAgentCLISection(sec))
	case "schema":
		sb.WriteString(r.renderAgentSchemaSection(sec))
	case "example":
		sb.WriteString(r.renderAgentExampleSection(sec))
	case "link":
		sb.WriteString(r.renderAgentLinkSection(sec))
	default:
		sb.WriteString(r.renderAgentDefaultSection(sec))
	}

	return sb.String()
}

// renderAgentCLISection renders CLI reference for agents
func (r *Renderer) renderAgentCLISection(sec map[string]interface{}) string {
	var sb strings.Builder
	title := getString(sec, "title", "CLI Reference")
	commands := getString(sec, "commands", "")

	sb.WriteString(fmt.Sprintf(`<section class="agent-section agent-cli"><h2>%s</h2>`, title))
	sb.WriteString(fmt.Sprintf(`<pre><code>%s</code></pre>`, escapeHTML(commands)))
	sb.WriteString(`</section>`)
	return sb.String()
}

// renderAgentSchemaSection renders a schema for agents
func (r *Renderer) renderAgentSchemaSection(sec map[string]interface{}) string {
	var sb strings.Builder
	title := getString(sec, "title", "Schema")
	schema := getString(sec, "schema", "")

	sb.WriteString(fmt.Sprintf(`<section class="agent-section agent-schema"><h2>%s</h2>`, title))
	sb.WriteString(fmt.Sprintf(`<pre><code>%s</code></pre>`, escapeHTML(schema)))
	sb.WriteString(`</section>`)
	return sb.String()
}

// renderAgentExampleSection renders an example for agents
func (r *Renderer) renderAgentExampleSection(sec map[string]interface{}) string {
	var sb strings.Builder
	title := getString(sec, "title", "Example")
	example := getString(sec, "example", "")
	language := getString(sec, "language", "yaml")

	sb.WriteString(fmt.Sprintf(`<section class="agent-section agent-example"><h2>%s</h2>`, title))
	sb.WriteString(fmt.Sprintf(`<pre><code class="language-%s">%s</code></pre>`, language, escapeHTML(example)))
	sb.WriteString(`</section>`)
	return sb.String()
}

// renderAgentLinkSection renders a link to another agent page
func (r *Renderer) renderAgentLinkSection(sec map[string]interface{}) string {
	var sb strings.Builder
	text := getString(sec, "text", "")
	href := getString(sec, "href", "#")

	sb.WriteString(fmt.Sprintf(`<p class="agent-link"><a href="%s">%s</a></p>`, href, text))
	return sb.String()
}

// renderAgentDefaultSection renders a default agent section
func (r *Renderer) renderAgentDefaultSection(sec map[string]interface{}) string {
	var sb strings.Builder
	title := getString(sec, "title", "")
	body := getString(sec, "body", "")

	if title != "" {
		sb.WriteString(fmt.Sprintf(`<section class="agent-section"><h2>%s</h2>`, title))
	} else {
		sb.WriteString(`<section class="agent-section">`)
	}

	if body != "" {
		html, _ := markdown.Render(body)
		sb.WriteString(html)
	}

	sb.WriteString(`</section>`)
	return sb.String()
}

// renderAgentNav renders navigation links between agent pages
func (r *Renderer) renderAgentNav(currentSlug string) string {
	agentPath := r.Collection.GetAgentPath()

	var sb strings.Builder
	sb.WriteString(`<nav class="agent-nav">`)

	// Find all pages with agent content
	for _, page := range r.Collection.Pages {
		if !r.Collection.HasAgentPage(page) {
			continue
		}
		meta := content.GetMeta(page)
		pageTitle := getString(meta, "title", "Page")
		pageSlug := getString(meta, "slug", "/")
		pageHref := agentPath + pageSlug

		active := ""
		if pageSlug == currentSlug || strings.HasPrefix(currentSlug, pageSlug) {
			active = ` class="active"`
		}

		sb.WriteString(fmt.Sprintf(`<a href="%s"%s>%s</a>`, pageHref, active, pageTitle))
	}

	sb.WriteString(`</nav>`)
	return sb.String()
}

// renderBlogContent renders just the blog post content without wrapper
func (r *Renderer) renderBlogContent(cf content.ContentFile) (string, error) {
	meta := content.GetMeta(cf)
	title := getString(meta, "title", "Post")
	desc := getString(meta, "description", "")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<h1>%s</h1>`, title))

	if desc != "" {
		sb.WriteString(fmt.Sprintf(`<p class="description">%s</p>`, desc))
	}

	// Check for sections first (new format)
	sectionsContent, _ := r.renderSections(cf)

	// Fallback to body if no sections
	bodyContent := ""
	if sectionsContent == "" {
		if cf.Markdown != nil {
			if ref, ok := cf.Markdown["@body.md"]; ok {
				html, _ := markdown.Render(ref)
				bodyContent = html
			} else if ref, ok := cf.Markdown["body.md"]; ok {
				html, _ := markdown.Render(ref)
				bodyContent = html
			}
		}

		// Fallback inline body
		if bodyContent == "" {
			if body, ok := cf.Data["body"].(string); ok {
				html, _ := markdown.Render(body)
				bodyContent = html
			}
		}
	}

	if sectionsContent != "" {
		sb.WriteString(sectionsContent)
	} else {
		sb.WriteString(bodyContent)
	}

	return sb.String(), nil
}

func (r *Renderer) renderSections(cf content.ContentFile) (string, error) {
	var sb strings.Builder
	sections, ok := cf.Data["sections"].([]interface{})
	if !ok {
		return "", nil
	}

	for _, s := range sections {
		sec, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		html := r.renderSection(sec)
		sb.WriteString(html)
	}

	return sb.String(), nil
}

func (r *Renderer) renderSection(sec map[string]interface{}) string {
	secType := getString(sec, "type", "default")
	id := getString(sec, "id", "")

	var sb strings.Builder

	// Handle special section types with custom rendering
	switch secType {
	case "hero":
		sb.WriteString(r.renderHeroSection(sec))
	case "features":
		sb.WriteString(r.renderFeaturesSection(sec))
	case "callout":
		sb.WriteString(r.renderCalloutSection(sec))
	case "pull-quote":
		sb.WriteString(r.renderPullQuote(sec))
	case "citation":
		sb.WriteString(r.renderCitation(sec))
	case "cta":
		sb.WriteString(r.renderCTASection(sec))
	case "warning":
		sb.WriteString(r.renderWarningSection(sec))
	case "code-example":
		sb.WriteString(r.renderCodeExample(sec))
	case "divider":
		// Section break/divider - a visual separator
		icon := getString(sec, "icon", "•••")
		sb.WriteString(fmt.Sprintf(`<div class="section-break" id="%s"><span class="section-break-icon">%s</span></div>`, id, icon))
	default:
		// Default section rendering
		sb.WriteString(fmt.Sprintf(`<section class="section section-%s" id="%s">`, secType, id))

		// Title
		if title := getString(sec, "title", ""); title != "" {
			level := getString(sec, "title_level", "h2")
			sb.WriteString(fmt.Sprintf(`<%s class="section-title">%s</%s>`, level, title, level))
		}

		// Content: @file.md ref or inline
		contentHTML := r.renderSectionContent(sec)
		if contentHTML != "" {
			sb.WriteString(`<div class="section-content">` + contentHTML + `</div>`)
		}

		// Cards
		if cards, ok := sec["cards"].([]interface{}); ok {
			sb.WriteString(`<div class="cards">`)
			for _, c := range cards {
				if card, ok := c.(map[string]interface{}); ok {
					sb.WriteString(r.renderCard(card))
				}
			}
			sb.WriteString(`</div>`)
		}

		// CTA buttons
		if buttons, ok := sec["buttons"].([]interface{}); ok {
			sb.WriteString(`<div class="cta-group">`)
			for _, b := range buttons {
				if btn, ok := b.(map[string]interface{}); ok {
					sb.WriteString(r.renderButton(btn))
				}
			}
			sb.WriteString(`</div>`)
		}

		sb.WriteString(`</section>`)
	}
	return sb.String()
}

func (r *Renderer) renderSectionContent(sec map[string]interface{}) string {
	// Support both "content" and "body" fields for section content
	contentRef := getString(sec, "content", "")
	if contentRef == "" {
		contentRef = getString(sec, "body", "")
	}
	if strings.HasPrefix(contentRef, "@") {
		mdFile := strings.TrimPrefix(contentRef, "@")
		if md, ok := sec["_md"].(map[string]string)[contentRef]; ok {
			html, _ := markdown.Render(md)
			return html
		} else if data, err := os.ReadFile(filepath.Join(r.contentDir, "pages", mdFile)); err == nil {
			html, _ := markdown.Render(string(data))
			return html
		}
	} else if contentRef != "" {
		html, _ := markdown.Render(contentRef)
		return html
	}
	return ""
}

// renderHeroSection - Hero with subtitle and CTA buttons
func (r *Renderer) renderHeroSection(sec map[string]interface{}) string {
	id := getString(sec, "id", "")
	title := getString(sec, "title", "")
	subtitle := getString(sec, "subtitle", "")
	body := getString(sec, "body", "")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<section class="section section-hero" id="%s">`, id))

	if title != "" {
		sb.WriteString(fmt.Sprintf(`<h1 class="hero-title">%s</h1>`, title))
	}

	if subtitle != "" {
		sb.WriteString(fmt.Sprintf(`<p class="hero-subtitle">%s</p>`, subtitle))
	}

	if body != "" {
		html, _ := markdown.Render(body)
		sb.WriteString(`<div class="hero-body">` + html + `</div>`)
	}

	// CTA buttons
	if buttons, ok := sec["buttons"].([]interface{}); ok && len(buttons) > 0 {
		sb.WriteString(`<div class="hero-buttons">`)
		for _, b := range buttons {
			if btn, ok := b.(map[string]interface{}); ok {
				sb.WriteString(r.renderButton(btn))
			}
		}
		sb.WriteString(`</div>`)
	}

	sb.WriteString(`</section>`)
	return sb.String()
}

// renderFeaturesSection - Features rendered as card grid
func (r *Renderer) renderFeaturesSection(sec map[string]interface{}) string {
	id := getString(sec, "id", "")
	title := getString(sec, "title", "")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<section class="section section-features" id="%s">`, id))

	if title != "" {
		sb.WriteString(fmt.Sprintf(`<h2 class="section-title">%s</h2>`, title))
	}

	// Features always render cards
	if cards, ok := sec["cards"].([]interface{}); ok && len(cards) > 0 {
		sb.WriteString(`<div class="features-grid">`)
		for _, c := range cards {
			if card, ok := c.(map[string]interface{}); ok {
				sb.WriteString(r.renderFeatureCard(card))
			}
		}
		sb.WriteString(`</div>`)
	}

	sb.WriteString(`</section>`)
	return sb.String()
}

func (r *Renderer) renderFeatureCard(card map[string]interface{}) string {
	title := getString(card, "title", "")
	body := getString(card, "body", "")
	icon := getString(card, "icon", "")

	var sb strings.Builder
	sb.WriteString(`<div class="feature-card">`)

	if icon != "" {
		sb.WriteString(fmt.Sprintf(`<div class="feature-icon">%s</div>`, icon))
	}

	if title != "" {
		sb.WriteString(fmt.Sprintf(`<h3 class="feature-title">%s</h3>`, title))
	}

	if body != "" {
		sb.WriteString(fmt.Sprintf(`<p class="feature-body">%s</p>`, body))
	}

	sb.WriteString(`</div>`)
	return sb.String()
}

// renderCalloutSection - Styled callout/info box
func (r *Renderer) renderCalloutSection(sec map[string]interface{}) string {
	id := getString(sec, "id", "")
	title := getString(sec, "title", "")
	body := getString(sec, "body", "")
	variant := getString(sec, "variant", "info") // info, tip, note

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<aside class="callout callout-%s" id="%s">`, variant, id))

	if title != "" {
		sb.WriteString(fmt.Sprintf(`<strong class="callout-title">%s</strong>`, title))
	}

	if body != "" {
		html, _ := markdown.Render(body)
		sb.WriteString(`<div class="callout-body">` + html + `</div>`)
	}

	sb.WriteString(`</aside>`)
	return sb.String()
}

// renderPullQuote - Blockquote with attribution
func (r *Renderer) renderPullQuote(sec map[string]interface{}) string {
	quote := getString(sec, "quote", "")
	attribution := getString(sec, "attribution", "")
	cite := getString(sec, "cite", "")

	var sb strings.Builder
	sb.WriteString(`<blockquote class="pull-quote">`)

	if quote != "" {
		html, _ := markdown.Render(quote)
		sb.WriteString(html)
	}

	if attribution != "" || cite != "" {
		sb.WriteString(`<footer class="pull-quote-attribution">`)
		if cite != "" {
			sb.WriteString(fmt.Sprintf(`<cite>%s</cite>`, cite))
		}
		if attribution != "" {
			sb.WriteString(fmt.Sprintf(`<span class="attribution">%s</span>`, attribution))
		}
		sb.WriteString(`</footer>`)
	}

	sb.WriteString(`</blockquote>`)
	return sb.String()
}

// renderCitation - Book/article reference
func (r *Renderer) renderCitation(sec map[string]interface{}) string {
	text := getString(sec, "text", "")
	url := getString(sec, "url", "")
	publisher := getString(sec, "publisher", "")
	author := getString(sec, "author", "")
	year := getString(sec, "year", "")

	var sb strings.Builder
	sb.WriteString(`<cite class="citation">`)

	if url != "" {
		sb.WriteString(fmt.Sprintf(`<a href="%s" class="citation-link">`, url))
	}

	if text != "" {
		sb.WriteString(text)
	}

	if url != "" {
		sb.WriteString(`</a>`)
	}

	// Citation details
	details := []string{}
	if author != "" {
		details = append(details, author)
	}
	if publisher != "" {
		details = append(details, publisher)
	}
	if year != "" {
		details = append(details, year)
	}

	if len(details) > 0 {
		sb.WriteString(fmt.Sprintf(`<span class="citation-details">%s</span>`, strings.Join(details, ", ")))
	}

	sb.WriteString(`</cite>`)
	return sb.String()
}

// renderCTASection - Call-to-action banner
func (r *Renderer) renderCTASection(sec map[string]interface{}) string {
	id := getString(sec, "id", "")
	title := getString(sec, "title", "")
	body := getString(sec, "body", "")
	buttonText := getString(sec, "button_text", "Learn More")
	buttonHref := getString(sec, "button_href", "#")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<section class="section section-cta" id="%s">`, id))

	if title != "" {
		sb.WriteString(fmt.Sprintf(`<h2 class="cta-title">%s</h2>`, title))
	}

	if body != "" {
		html, _ := markdown.Render(body)
		sb.WriteString(`<div class="cta-body">` + html + `</div>`)
	}

	if buttonText != "" {
		sb.WriteString(fmt.Sprintf(`<a href="%s" class="button button-primary">%s</a>`, buttonHref, buttonText))
	}

	sb.WriteString(`</section>`)
	return sb.String()
}

// renderWarningSection - Styled warning/alert box
func (r *Renderer) renderWarningSection(sec map[string]interface{}) string {
	id := getString(sec, "id", "")
	title := getString(sec, "title", "")
	body := getString(sec, "body", "")
	variant := getString(sec, "variant", "warning") // warning, error, info, success

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<div class="warning-box warning-%s" id="%s">`, variant, id))

	if title != "" {
		sb.WriteString(fmt.Sprintf(`<strong class="warning-title">%s</strong>`, title))
	}

	if body != "" {
		html, _ := markdown.Render(body)
		sb.WriteString(`<div class="warning-body">` + html + `</div>`)
	}

	sb.WriteString(`</div>`)
	return sb.String()
}

// renderCodeExample - Code block with syntax highlighting
func (r *Renderer) renderCodeExample(sec map[string]interface{}) string {
	id := getString(sec, "id", "")
	title := getString(sec, "title", "")
	language := getString(sec, "language", "")
	code := getString(sec, "code", "")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<div class="code-example" id="%s">`, id))

	if title != "" {
		sb.WriteString(fmt.Sprintf(`<div class="code-example-header">`))
		if language != "" {
			sb.WriteString(fmt.Sprintf(`<span class="code-example-lang">%s</span>`, language))
		}
		sb.WriteString(fmt.Sprintf(`<span class="code-example-title">%s</span>`, title))
		sb.WriteString(`</div>`)
	}

	// Render code as pre/code block
	if code != "" {
		escapedCode := escapeHTML(code)
		if language != "" {
			sb.WriteString(fmt.Sprintf(`<pre class="code-block"><code class="language-%s">%s</code></pre>`, language, escapedCode))
		} else {
			sb.WriteString(fmt.Sprintf(`<pre class="code-block"><code>%s</code></pre>`, escapedCode))
		}
	}

	sb.WriteString(`</div>`)
	return sb.String()
}

// escapeHTML escapes HTML special characters for code blocks
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func (r *Renderer) renderCard(card map[string]interface{}) string {
	title := getString(card, "title", "")
	body := getString(card, "body", "")
	variant := getString(card, "variant", "default")
	href := getString(card, "href", "")

	html := `<div class="card card-` + variant + `">`
	if title != "" {
		if href != "" {
			html += `<h3 class="card-title"><a href="` + href + `">` + title + `</a></h3>`
		} else {
			html += `<h3 class="card-title">` + title + `</h3>`
		}
	}
	if body != "" {
		html += `<p class="card-body">` + body + `</p>`
	}
	html += `</div>`
	return html
}

func (r *Renderer) renderButton(btn map[string]interface{}) string {
	text := getString(btn, "text", "Button")
	href := getString(btn, "href", "#")
	variant := getString(btn, "variant", "primary")
	return fmt.Sprintf(`<a href="%s" class="btn btn-%s">%s</a>`, href, variant, text)
}

func (r *Renderer) renderNav(currentSlug string) string {
	nav := r.Collection.GetNav()
	if len(nav) == 0 {
		return `<nav class="nav">
		<div class="nav-inner">
			<a href="/" class="nav-brand">lyt</a>
		</div>
	</nav>`
	}

	var sb strings.Builder
	sb.WriteString(`<nav class="nav">`)
	sb.WriteString(`<div class="nav-inner">`)
	sb.WriteString(`<a href="/" class="nav-brand">lyt</a>`)

	// Desktop links
	sb.WriteString(`<ul class="nav-links">`)
	for _, item := range nav {
		if m, ok := item.(map[string]interface{}); ok {
			label := getString(m, "label", "")
			href := getString(m, "href", "#")
			active := r.isActiveNavLink(href, currentSlug)
			sb.WriteString(fmt.Sprintf(`<li><a href="%s" class="nav-link"%s>%s</a></li>`, href, active, label))
		}
	}
	sb.WriteString(`</ul>`)

	// Mobile hamburger + menu (checkbox hack for JS-free toggle)
	sb.WriteString(`<input type="checkbox" id="nav-menu" class="nav-menu-checkbox">`)
	sb.WriteString(`<label for="nav-menu" class="nav-toggle" aria-label="Toggle menu"><span></span></label>`)
	sb.WriteString(`<div class="nav-menu">`)
	sb.WriteString(`<ul class="nav-mobile-links">`)
	for _, item := range nav {
		if m, ok := item.(map[string]interface{}); ok {
			label := getString(m, "label", "")
			href := getString(m, "href", "#")
			active := r.isActiveNavLink(href, currentSlug)
			sb.WriteString(fmt.Sprintf(`<li><a href="%s" class="nav-mobile-link"%s>%s</a></li>`, href, active, label))
		}
	}
	sb.WriteString(`</ul>`)
	sb.WriteString(`</div>`)

	sb.WriteString(`</div>`)
	sb.WriteString(`</nav>`)
	return sb.String()
}

// isActiveNavLink returns aria-current="page" if the href matches the current slug
func (r *Renderer) isActiveNavLink(href, currentSlug string) string {
	// Normalize both for comparison
	href = strings.TrimSuffix(href, "/")
	slug := strings.TrimSuffix(currentSlug, "/")

	if href == slug {
		return ` aria-current="page"`
	}
	// Also check if currentSlug starts with href (for section pages like /docs/)
	if strings.HasPrefix(slug, href) && href != "" && href != "/" {
		return ` aria-current="page"`
	}
	return ""
}

// renderCanonical generates the canonical URL meta tag
func (r *Renderer) renderCanonical(slug string) string {
	baseURL := "https://lyt.b7r.dev" // Could be configurable
	if slug == "" || slug == "/" {
		return fmt.Sprintf(`<link rel="canonical" href="%s/">`, baseURL)
	}
	return fmt.Sprintf(`<link rel="canonical" href="%s%s">`, baseURL, slug)
}

func (r *Renderer) executeTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := htmltemplate.New("page").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getString(m map[string]interface{}, key, fallback string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return fallback
}

func slugToClass(slug string) string {
	slug = strings.Trim(slug, "/")
	slug = strings.ReplaceAll(slug, "/", "-")
	return slug
}

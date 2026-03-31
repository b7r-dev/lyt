# Changelog

All notable changes to lyt will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.1] - 2026-03-31

### Fixed
- **Init templates** — Refreshed all embedded templates to match current codebase:
  - Design tokens now use nested structure (base, accent, system colors)
  - CSS updated with full component styles (1182 lines)
  - Site config includes agent_section and correct nav
  - Index and about pages include agent_content for AI assistance
  - Added blog.yaml and docs.yaml templates

## [1.1.0] - 2026-03-30

### Added
- **Schema validation** (`lyt validate --schema`):
  - Validates YAML content structure against schema
  - Checks required fields (title, slug)
  - Validates section types (hero, features, callout, etc.)
  - Validates CTA fields (button_text + button_href required together)
  - Validates code-example has code field
  - Validates agent sections (cli, example, link, schema types)
  - Validates config nav items and agent sections
- **Link validation** (`lyt validate --links`):
  - Reads sitemap.xml to get all pages
  - Checks all internal links in HTML files
  - Resolves relative URLs correctly (same-dir, parent, subdir)
  - Skips external links, anchors, mailto, tel
- **`lyt validate` command**:
  - Default validates both schema and links
  - `--schema` flag for schema-only validation
  - `--links` flag for link-only validation
  - `--dir` flag to specify custom dist directory
  - `--strict` to treat warnings as errors
- **Comprehensive unit tests** for all validation logic
- **Agent documentation** at multiple levels (CLI help, human docs, agent_content)

### Fixed
- YAML indentation bug in getting-started page causing 404
- Missing agent_content sections in documentation parity

### Changed
- Project restructured: Go module moved to root for `go install` compatibility
- CI workflow updated for new module structure

## [1.0.0] - 2026-03-29

### Added
- **Initial public release**
- YAML content parsing (pages, blog posts)
- Markdown rendering (goldmark, GFM)
- HTML template rendering
- Design tokens → CSS custom properties
- Dev server with hot reload
- Static asset copying
- Sitemap generation
- Blog index page with post cards
- Deployment documentation (Vercel, Netlify, Cloudflare, rsync)
- CI/CD with GitHub Actions

### Fixed
- Blog index now renders with post listing
- Card links now render as clickable anchors
- Docs path routing
- Nav duplication issues
- Dark mode palette (beige/fireflies theme)
- Consistent stacked shadows on warning boxes

---

## Version History

| Version | Date | Notes |
|---------|------|-------|
| 1.1.1 | 2026-03-31 | Fix init templates |
| 1.1.0 | 2026-03-30 | Schema + link validation |
| 1.0.0 | 2026-03-29 | Initial public release |
| 0.2.0 | 2026-03-29 | Composable content + rendering quality |
| 0.1.0 | 2026-03-29 | MVP - Core engine |

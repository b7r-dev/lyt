# Changelog

All notable changes to lyt will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2026-03-29

### Added
- **Composable components**:
  - `pull-quote` - Blockquote with attribution and cite
  - `citation` - Book/article reference with URL, author, publisher, year
  - `cta` - Call-to-action banner with title, body, button
  - `warning` - Styled alert box (warning, error, info, success variants)
  - `callout` - Styled info/tip/note boxes
- **Section type rendering**:
  - Hero now renders subtitle and CTA buttons
  - Features renders as card grid with icons
  - Callout section type with styled variants
- **Incremental builds** - Cache-based rebuild detection with `lyt build`
- `--force` / `-f` flag to force rebuild
- Build time displayed after build

### Changed
- Improved section type rendering for all types
- Build performance optimizations (<1s for typical sites)
- Tests for build cache

## [0.1.0] - 2026-03-29

### Added
- YAML content parsing (pages, blog posts)
- Markdown rendering (goldmark, GFM)
- HTML template rendering
- Design tokens → CSS custom properties
- Dev server with hot reload
- Static asset copying
- Sitemap generation
- Blog index page with post cards
- Card href support for linking

### Fixed
- Blog index now renders with post listing
- Card links now render as clickable anchors

---

## Version History

| Version | Date | Notes |
|---------|------|-------|
| 0.1.0 | 2026-03-29 | MVP - Core engine |
| 0.2.0 | - | Composable content + rendering quality |

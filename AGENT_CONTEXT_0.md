# lyt — Agent Context

## Project Overview

`lyt` is a minimal static site generator. It is a fork of `fyt` (a React/Node.js SSG) stripped of all JS ecosystem dependencies.

**Core principle**: Engine, content, markup, and styling are kept strictly separate. They only meet at build time.

**No runtime JavaScript.** Output is pure static HTML. Zero JS in the browser.

**Publish workflow**: `git push` → CI runs `lyt build` → deploy `dist/` to any static host.

## Repository

```
/Users/bweller/src/b7r/lyt/
```

## Directory Structure

```
lyt/
├── engine/           # Go source code (the tool)
│   ├── main.go       # Entry point
│   ├── go.mod       # Go module (github.com/b7r-dev/lyt/engine)
│   ├── go.sum       # Dependencies
│   ├── cmd/
│   │   ├── root.go      # CLI root, cobra initialization
│   │   ├── build.go     # `lyt build` command
│   │   └── serve.go     # `lyt serve` command (dev server)
│   └── internal/
│       ├── content/
│       │   └── scanner.go   # YAML scanner, @file.md resolution
│       ├── markdown/
│       │   └── markdown.go   # goldmark processor (CommonMark + GFM)
│       ├── render/
│       │   └── render.go    # HTML generator (go html/template)
│       └── tokens/
│           └── processor.go  # Design tokens → CSS custom properties
├── content/          # Site content (YAML + Markdown)
│   ├── tokens.yaml      # Design tokens (colors, spacing, typography, z-planes)
│   ├── config/
│   │   └── site.yaml    # Site config (nav, meta)
│   ├── pages/           # Page YAML files
│   ├── blog/            # Blog post YAML files (YAML + inline Markdown body)
│   └── docs/            # Docs page (placeholder)
├── templates/
│   └── base.css         # Base stylesheet (mobile-first, 3-plane depth)
├── public/
│   └── favicon.svg      # SVG favicon
└── dist/                # Built output (generated, not committed)
```

## Key Technologies

| Layer | Technology |
|-------|------------|
| CLI framework | spf13/cobra |
| YAML parsing | gopkg.in/yaml.v3 |
| Markdown | yuin/goldmark (CommonMark + GFM) |
| HTML templating | Go html/template (stdlib) |
| File watching | fsnotify/fsnotify |
| WebSocket (hot reload) | gorilla/websocket |

**Templ**: Originally included but removed — version `v0.3.866` doesn't exist. Re-add with correct version once verified. Current approach uses Go `html/template` with inline template strings in `render.go`.

## Content Schema

### Pages (`content/pages/*.yaml`)

```yaml
meta:
  title: "Page Title"
  slug: "/page-slug"        # URL path (e.g., "/about")
  description: "SEO desc"

sections:
  - id: "section-id"
    type: "default"        # or "hero", "features", etc.
    title: "Section Title"
    title_level: "h2"      # optional, defaults to h2
    body: |                 # inline Markdown
        Section **content** here.
    # OR
    content: "@filename.md" # Markdown file reference (resolved by scanner)
    cards:                 # optional card grid
      - title: "Card"
        body: "Card content"
        variant: "default"
    buttons:               # optional CTA buttons
      - text: "Click me"
        href: "/link"
        variant: "primary"
```

### Blog posts (`content/blog/*.yaml`)

```yaml
meta:
  title: "Post Title"
  slug: "post-slug"       # becomes /blog/post-slug/index.html
  description: "..."
  date: "2026-03-29"
  author: "Name"
  tags: ["tag1", "tag2"]
  published: true

body: |                    # inline Markdown body
    ## Heading

    Prose here.
# OR
body: "@body.md"          # external Markdown file reference
```

### Markdown linking

Fields prefixed with `@` or referencing `.md` files are resolved by the scanner:
- `@filename.md` → load file relative to the YAML file's directory
- Resolved content stored in `ContentFile.Markdown` map

### Design tokens (`content/tokens.yaml`)

```yaml
colors:
  base:
    bg: "#f5f5f0"
    surface: "#ecece6"
    border: "#d8d8d0"
    muted: "#7a7a6e"
    text: "#3d3d3d"
    heading: "#2a2a25"
    link: "#5c5c4a"
    link-hover: "#3d3d3d"

spacing:
  0: "0"
  1: "0.25rem"
  ...

typography:
  font_family:
    body: "Georgia, 'Times New Roman', serif"
    heading: "'Courier New', Courier, monospace"
    ...
  font_size:
    xs: "0.75rem"
    ...
  font_weight:
    normal: "400"
    ...
  line_height:
    normal: "1.5"
    ...

z:
  base: "0"       # Prose, body text
  middle: "10"    # Cards, elevated elements
  top: "100"      # Nav, floating elements
  modal: "1000"   # Modal overlay
```

Output: CSS custom properties (`--color-base-bg`, `--z-base`, etc.)

## Design System

- **Palette**: warm beige (`#f5f5f0` bg), desaturated grays, softened blacks
- **Typography**: Georgia serif (body), Courier mono (headings/brand/code)
- **Depth**: 3 z-planes conveyed via stacked shadow offsets + `backdrop-filter: blur()`
- **Mobile**: `clamp()` for fluid type, grid collapses at breakpoints, nav hides on mobile
- **Shadows**: Two-layer stacked offsets (solid + faded), no pseudo-elements, no extra resources

Key CSS features used:
- `backdrop-filter: blur(12px)` on nav
- Stacked `box-shadow` for depth simulation
- `color-mix()` for shadow fade
- CSS custom properties from tokens

## CLI Commands

```bash
cd engine
go build -o lyt

./lyt build              # Build to ../dist
./lyt build -o /path     # Build to custom output dir
./lyt serve              # Dev server on :5173 with hot reload
./lyt serve --port 8080  # Custom port
./lyt serve --static     # Serve pre-built dist without watching
```

## Build Pipeline

1. Scan `content/` for YAML files
2. Resolve `@file.md` references in content
3. Process `tokens.yaml` → `tokens.css`
4. Copy `templates/base.css` → `dist/`
5. Render each page/blog post to HTML (using inline template in `render.go`)
6. Copy `public/` assets → `dist/`
7. Generate `sitemap.xml`

## Render Pipeline (per page)

1. Parse YAML → `ContentFile` with `Data` map and `Markdown` map
2. Extract `meta` (title, slug, description)
3. For each section:
   - Inline Markdown body → HTML (goldmark)
   - `@file.md` reference → load file → Markdown → HTML
   - Render cards, buttons
4. Render nav from `config/site.yaml`
5. Execute HTML template (go `html/template`)
6. Write to `dist/[slug].html`

## Current State

- **Working**: `lyt build` and `lyt serve` both functional
- **Server running**: `*:57837` (bound to `0.0.0.0`)
- **Design**: Functional beige 3-plane system
- **TODOs**: See below

## Known Issues / TODOs

1. **Templ not integrated**: Researched and decided on Templ, but version issue prevented integration. Currently using `html/template` inline strings in `render.go`. Re-add with correct Templ version.

2. **No component templating**: Components defined in YAML, rendered as inline strings. True Templ `.templ` files in `templates/components/` would enable reusable, version-controlled components.

3. **Section types not rendered**: `type: "hero"` in YAML produces `<section class="section section-hero">` but no special rendering. Need hero-specific HTML generation.

4. **Blog index page**: `content/pages/blog.yaml` has no sections — the blog listing page renders empty.

5. **No tag/category pages**: Blog posts have tags but no tag pages generated.

6. **No RSS feed**: No `feed.xml` generation.

7. **Markdown frontmatter**: Blog posts use YAML `meta` block; true frontmatter in `.md` files (with `---` delimiters) not supported.

8. **Serve hot reload**: Works but only rebuilds content. Doesn't refresh WebSocket clients reliably.

9. **Build paths**: Build/serve use `../content`, `../dist`, `../templates` relative to `engine/`. Must run from `engine/` directory. Absolute paths or config-based paths would be more robust.

10. **`@file.md` reference resolution**: Currently only works for YAML fields; not integrated into section `content` field resolution.

## Style Conventions

- Go: standard formatting (`go fmt`), error handling with `fmt.Errorf("...: %w", err)`
- CSS: 2-space indent, BEM-ish class naming, CSS custom properties for all values
- YAML: lowercase keys, no trailing whitespace, `|` for multiline strings
- No emoji in code or commit messages

## Git

- Branch: `main`
- Binary (`engine/lyt`) and `dist/` are in `.gitignore`
- `go.sum` is committed (engine subdirectory)

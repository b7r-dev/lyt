# lyt

**[Live Site](https://lyt.b7r.dev)** · **[GitHub](https://github.com/b7r-dev/lyt)** · **[Introduction for Agents](https://lyt.b7r.dev/agents)**

lyt is a minimal static site generator. No runtime JavaScript. No build step for templates. Just YAML, Markdown, and Go.

## Structure

```
lyt/
├── cmd/           ← CLI commands (Go)
├── internal/      ← Core libraries (Go)
├── content/       ← Site content (YAML, Markdown)
├── templates/     ← Styling (CSS)
├── public/        ← Static assets (fonts, images)
└── dist/          ← Built output (generated)
```

Code, content, markup, and styling are kept strictly separate. They only meet at build time.

## Quick start

```bash
# Install
go install github.com/b7r-dev/lyt@latest

# Or build from source
git clone https://github.com/b7r-dev/lyt.git
cd lyt
go build -o lyt .

# Build the site
./lyt build

# Serve locally (with hot reload)
./lyt serve
```

## Content

Pages are YAML files in `content/pages/`:

```yaml
meta:
  title: "Page Title"
  slug: "/page-slug"
  description: "SEO description"

sections:
  - id: "intro"
    type: "default"
    title: "Section Title"
    body: |
      Section content in Markdown.
```

Long-form content lives in separate Markdown files, referenced via `@filename.md`:

```yaml
sections:
  - id: "article"
    content: "@blog/my-article.md"
```

Blog posts live in `content/blog/`. Same format as pages, plus a `body` field (inline Markdown) or `body.md` file reference.

## Design tokens

Tokens live in `content/tokens.yaml`. They're processed into CSS custom properties:

```yaml
colors:
  base:
    bg: "#f5f5f0"
    text: "#3d3d3d"
```

Output:

```css
:root {
  --color-base-bg: #f5f5f0;
  --color-base-text: #3d3d3d;
}
```

## Design system

The design system is built on three principles:

1. **3-plane z-axis** — `z: 0` (prose), `z: 10` (cards), `z: 100` (nav)
2. **Beige palette** — warm, desaturated, content-forward
3. **Mobile-first** — works without JavaScript on any device

Depth is conveyed through stacked shadow offsets and `backdrop-filter: blur()` — no images, no extra markup.

## Publishing

Publish is `git push`. Run `lyt build` in CI/CD and deploy `dist/` to any static host.

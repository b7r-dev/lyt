# lyt

**[Live Demo](https://lyt.b7r.dev)** · **[GitHub](https://github.com/b7r-dev/lyt)** · **[Agent Docs](https://lyt.b7r.dev/agents)**

lyt is a minimal static site generator written in Go. It transforms YAML content and Markdown into pure HTML—zero runtime JavaScript, no client-side framework, no template compilation step.

## Why lyt?

Modern static site generators often ship megabytes of JavaScript to the browser. lyt takes a different approach:

- **Zero JS output** — The built site is pure HTML and CSS. It works on any device, including text-only browsers.
- **Content/Style/Engine separation** — Your content lives in YAML, styling in CSS, and the engine is just... Go code. They only meet at build time.
- **Design tokens** — Define colors, typography, and spacing once in YAML. lyt transforms them into CSS custom properties.
- **Built-in components** — Hero sections, feature grids, callouts, CTAs, warnings, pull quotes—declared in YAML, rendered as semantic HTML.
- **Schema validation** — Catch missing fields, invalid types, and broken links before you deploy.
- **AI-ready** — Includes structured help for AI agents (`lyt help agent`) so assistants can work with your content intelligently.

## Quick Start

```bash
# Install
go install github.com/b7r-dev/lyt@latest

# Or build from source
git clone https://github.com/b7r-dev/lyt
cd lyt
go build -o lyt .

# Create a project
lyt init my-site
cd my-site

# Build
./lyt build

# Develop with hot reload
./lyt serve
```

## How It Works

### Content: YAML + Markdown

Pages live in `content/pages/`:

```yaml
meta:
  title: "Hello World"
  slug: "/hello"
  description: "A simple page"

sections:
  - id: "intro"
    type: "hero"
    title: "Hello World"
    subtitle: "Welcome to my site"

  - id: "about"
    type: "default"
    title: "About"
    body: |
      This is **Markdown** content.
```

Blog posts go in `content/blog/` with an added `date` field. Long-form content can reference external Markdown files with `@filename.md`.

### Design Tokens

Define your visual theme in `content/tokens.yaml`:

```yaml
colors:
  base:
    bg: "#faf9f6"
    text: "#2d2d2d"
  accent:
    primary: "#e07a5f"
    
typography:
  fonts:
    body: "system-ui, sans-serif"
```

These become CSS custom properties in your build output—no runtime processing needed.

### Section Components

lyt renders these section types from your YAML:

| Type | Output |
|------|--------|
| `hero` | Full-width header with title, subtitle, body, buttons |
| `features` | Card grid with icons and descriptions |
| `default` | Standard prose section |
| `cta` | Call-to-action banner |
| `callout` | Styled info/tip/warning boxes |
| `pull-quote` | Blockquote with attribution |
| `citation` | Book/article reference |
| `warning` | Alert box (error, warning, info, success) |
| `code-example` | Syntax-highlighted code block |

## Commands

| Command | Description |
|---------|-------------|
| `lyt build` | Build the site to `./dist` |
| `lyt serve` | Dev server with hot reload |
| `lyt init` | Initialize a new project |
| `lyt validate` | Check content schema and links |
| `lyt help agent` | AI agent documentation |

Flags:
- `-o, --output` — Output directory (default: `./dist`)
- `-f, --force` — Force rebuild
- `-v, --verbose` — Detailed output

## Project Structure

```
my-site/
├── content/
│   ├── pages/         # Page YAML files
│   ├── blog/          # Blog post YAML files
│   ├── config/        # Site configuration
│   ├── tokens.yaml   # Design tokens
│   └── schema.yaml   # Content schema
├── templates/        # HTML templates + CSS
├── public/           # Static assets
└── dist/             # Built output
```

## Deployment

Build output is static files in `./dist`. Deploy to anything:

- **Vercel** — Connect repo, runs `lyt build`
- **Netlify** — Connect repo, runs `lyt build`  
- **Cloudflare Pages** — Connect repo, runs `lyt build`
- **GitHub Pages** — Use GitHub Actions
- **VPS/rsync** — `lyt build -o ./dist && rsync -av dist/ user@server:`

## Validation

lyt validates your content before building:

```bash
# Validate schema and links
lyt validate

# Schema only
lyt validate --schema

# Links only  
lyt validate --links

# Custom dist directory
lyt validate --links --dir /path/to/dist
```

Checks include:
- Required fields (title, slug)
- Section type validity
- CTA field pairs (button_text + button_href)
- Internal link integrity
- Agent section structures

## Philosophy

1. **Separation of concerns** — Content, style, and engine are independent. Swap any without touching the others.
2. **No runtime dependencies** — The built site needs no JavaScript, no CDN, no special hosting.
3. **Fast builds** — Incremental builds with change detection. Typical rebuilds complete in milliseconds.
4. **Agent-friendly** — AI assistants can read, modify, and validate your content because it's just structured YAML.

## License

MIT

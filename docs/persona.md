# User Persona: The CLI-Native Content Creator

## Persona: Martin

**Age**: 44  
**Role**: Technical writer, indie hacker, maintainer of small open-source projects  
**Tech comfort level**: High — lives in the terminal, skeptical of web-based GUIs  

### Background

Martin has been building websites since 2003. He's used every generation of static site generator: first Movable Type's Perl scripts, then Jekyll when it launched in 2008, then Hugo in 2013. He has 15 years of content across four different site generators, scattered across repositories, exported and re-imported more times than he can count.

His current site — a personal blog and project docs — runs on a 2019 Hugo version that's so outdated he's afraid to run `hugo update`. The theme is a fork of a fork of a theme he doesn't remember choosing. Half the shortcodes don't work. The build times are 40 seconds for 80 pages. He keeps meaning to migrate but the thought of touching the config file makes him tired.

Martin doesn't want a new hobby. He wants to write. He wants the site to just work.

### Goals

1. **Migrate existing content with minimal friction** — export from old formats, import cleanly
2. **Control everything from the terminal** — no web dashboards, no browser-based editors
3. **Fast build times** — seconds, not tens of seconds
4. **Clean, predictable output** — he needs to understand exactly what HTML is being generated
5. **Easy deployment** — push to any host, no locked-in platform
6. **Global installation** — install once, run from any directory
7. **Project-relative paths** — output `./dist` in current working directory
8. **Instant project creation** — `lyt init` creates a ready-to-use project in seconds
9. **LLM-friendly help** — `lyt help agent` gives AI agents the full instruction tree

### Pain Points

- Configuration drift between theme versions
- Mysterious build failures when upgrading dependencies
- Bloated JavaScript bundles in output he didn't ask for
- Needing to learn a templating language he'll only use for this one site
- Theme authors who abandon projects, leaving him with security vulnerabilities
- Slow iteration cycle: write → preview → wait → adjust
- **Hardcoded theme components** — he can't add a pull-quote or citation block without editing theme templates
- **Shortcode chaos** — each theme defines its own shortcodes with different syntax; migrating means rewriting content
- **No reusability** — when he writes a good component in one post, there's no way to use it in another without copypasting HTML
- **Project-specific binaries** — having to keep a binary in each project directory, or remembering complex paths
- **Output path confusion** — build output going to sibling directories instead of current working directory
- **Manual project setup** — having to create directory structure and copy template files manually
- **LLM ignorance** — AI agents don't know lyt's commands, workflows, or best practices

---

## User Story: Migrating a Decade of Content

### Scenario

Martin sits down on a Saturday morning with a fresh install of `lyt`. He's going to migrate his blog from that old Hugo version.

### Story

**Given** Martin installs lyt via `go install` or a package manager,  
**When** he runs `lyt --version` from any terminal,  
**Then** lyt should respond with its version number,  
**And** the binary should be in his PATH — no need to manage project-specific binaries.

---

**Given** Martin has a lyt project in any directory,  
**When** he runs `cd /path/to/my-blog && lyt build`,  
**Then** lyt should detect this is a lyt project and build from the current directory,  
**And** the output should go to `./dist` (in the current working directory), not a sibling directory.

---

**Given** Martin has a directory of markdown content exported from his old site,  
**And** the content includes frontmatter with dates, tags, and draft status,  
**When** he runs `lyt import --from hugo ./content`,  
**Then** lyt should parse the frontmatter and produce valid YAML files in the expected directory structure,  
**And** any unsupported frontmatter fields should be preserved as comments or silently ignored.

---

**Given** Martin places his YAML content files in `content/pages/` and `content/blog/`,  
**And** the YAML files reference markdown files via `@filename.md`,  
**When** he runs `lyt build`,  
**Then** lyt should resolve those references and render complete HTML,  
**And** the build should complete in under 2 seconds for 100 content files.

---

**Given** Martin runs `lyt serve` to preview his site,  
**When** he edits a content file and saves it,  
**Then** the server should rebuild and refresh the page automatically,  
**And** the rebuild should complete in under 500ms,  
**And** the browser should update without a full page reload if possible.

---

**Given** Martin has configured his deployment target (Cloudflare Pages, Netlify, or just rsync to a VPS),  
**When** he runs `lyt build -o ./dist`,  
**Then** the output directory should contain only static files — no build artifacts, no source, no config,  
**And** the HTML should be readable: indented, no minified inline scripts, semantic markup he can verify by hand if needed.

---

**Given** Martin wants to deploy with a single command,  
**When** he runs `lyt build && lyc deploy` (or `lyt build && cp -r dist/* server:/var/www/`),  
**Then** the site should be live with no additional processing required on the server,  
**And** the output should work on any basic HTTP server — Nginx, Caddy, GitHub Pages.

---

**Given** Martin needs to verify his rendered output is correct,  
**When** he inspects `dist/blog/my-post/index.html`,  
**Then** the HTML should use semantic tags (`<article>`, `<main>`, `<nav>`),  
**And** styles should be loaded via `<link>` tags pointing to static CSS files,  
**And** there should be zero `<script>` tags unless he explicitly added them in his markdown.

---

---

**Given** Martin has a fresh idea for a new blog project,  
**When** he runs `mkdir my-new-blog && cd my-new-blog && lyt init`,  
**Then** lyt should create a complete project structure with:
- `content/pages/index.yaml` — a home page with hero and features sections
- `content/pages/about.yaml` — an about page template
- `content/blog/hello-world.yaml` — a sample blog post with lorem ipsum
- `content/tokens.yaml` — default design tokens
- `templates/base.css` — default stylesheet
- `public/` — empty directory for assets

**And** the project should be immediately buildable with `lyt build`,  
**And** the output should demonstrate all available components (hero, features, cta, pull-quote, etc.).

---

**Given** Martin runs `lyt init` in a non-empty directory,  
**When** lyt detects existing files,  
**Then** it should prompt for confirmation before overwriting any conflicting files,  
**Or** use a `--force` flag to skip prompts.

---

**Given** Martin uses an LLM to help build his site,  
**When** the LLM runs `lyt help agent`,  
**Then** lyt should output a structured instruction tree containing:
- All available commands with descriptions
- All flags for each command
- Project directory structure
- Content file schema (YAML frontmatter fields)
- Available component types
- Recommended workflow for common tasks
- Deployment options

**And** the output should be parseable by an LLM (JSON or markdown with clear headings),  
**So that** the LLM can provide accurate, context-aware assistance without guessing.

---

**Given** Martin has defined a pull-quote component in his content schema,  
**And** his blog post YAML includes:

```yaml
sections:
  - id: "key-point"
    type: "pull-quote"
    quote: "The tools we use shape the work we produce."
    attribution: "Martin, 2024"
```

**When** lyt builds the page,  
**Then** it should render semantic `<blockquote>` HTML with proper citation markup,  
**And** the styling should come from `base.css` — no additional CSS required in the post,  
**And** the component should be reusable across any post by simply changing the YAML fields.

---

**Given** Martin wants to include a footnote or citation at the bottom of a blog post,  
**And** his YAML includes:

```yaml
sections:
  - type: "citation"
    text: "Design Patterns of Successful Software"
    url: "https://example.com/book"
    publisher: "O'Reilly, 2023"
```

**When** the page renders,  
**Then** the citation should appear as a properly linked reference in the footer or at the section's position,  
**And** multiple citations should be numbered or bulleted consistently.

---

**Given** Martin needs a call-to-action block at the end of certain blog posts,  
**And** his YAML includes:

```yaml
sections:
  - type: "cta"
    title: "Want more?"
    body: "Subscribe to the newsletter for weekly posts on tools and craft."
    button_text: "Subscribe"
    button_href: "/subscribe"
```

**When** lyt builds the page,  
**Then** it should render a styled CTA section with the button,  
**And** the CTA should be visually distinct from prose (via CSS classes),  
**And** he should not need to touch any template files to add or modify CTA blocks.

---

**Given** Martin discovers he needs a new component type — say, a "warning" banner for deprecation notices,  
**When** he adds a `type: "warning"` section to his YAML,  
**Then** lyt should render it as a styled warning box with appropriate semantic markup,  
**And** future posts can use it immediately without code changes.

---

**Given** Martin is tired of copying HTML snippets between posts when he wants the same layout,  
**When** he defines reusable component templates in a central location,  
**Then** he can reference those components by name in any YAML file,  
**And** the rendered output should be consistent every time.

---

### Success Criteria

| Criterion | Target |
|-----------|--------|
| Import existing markdown/frontmatter content | Clean conversion with no data loss |
| Build time (100 pages) | < 2 seconds |
| Serve rebuild time | < 500ms |
| Output HTML | Readable, semantic, JS-free by default |
| Deployment | One command to static host |
| CLI ergonomics | Flags are predictable, `--help` is complete, no surprise behavior |
| Composable templates | Pull-quotes, hero sections, citations, CTAs defined in YAML, rendered consistently |
| Component reusability | Same component used in multiple posts renders identically |
| Global command | `lyt` works from any directory after install |
| Project-relative output | `lyt build` outputs to `./dist` in current directory |
| Instant project creation | `lyt init` scaffolds a ready-to-build project in seconds |
| LLM-native help | `lyt help agent` provides full instruction tree for AI agents |

---

## Design Principles Derived

1. **CLI first** — Every feature accessible from command line. No web UI, no optional dashboard.
2. **Fast iteration** — Sub-second rebuilds enable flow state while writing.
3. **Transparent output** — What you write is what you get. No hidden transformations.
4. **Zero JS in output** — Unless the author explicitly includes it. The browser receives HTML + CSS only.
5. **Portable deployment** — Output works on any static host. No platform lock-in.
6. **Schema over configuration** — Content structure is defined by YAML schema, not by tweaking theme templates.
7. **Composable templates** — Rich components (pull-quotes, citations, CTAs, warnings) defined once, used anywhere.
8. **Component-driven content** — Content files declare structure via YAML; rendering is handled by reusable components.

> "Give me a binary, a content folder, and a template. That's it. Everything else is overhead."
> — Martin

# ADR-005: Portable Static Output

**Status**: Accepted

## Context

Martin's goals include:
- Easy deployment - push to any host, no locked-in platform
- Output should work on any basic HTTP server (Nginx, Caddy, GitHub Pages)

His pain points include:
- Platform lock-in from some static site generators

## Decision

The `dist/` output directory will contain **only static files** that work on any HTTP server:

```
dist/
├── index.html
├── about.html
├── blog/
│   ├── index.html
│   └── my-post/
│       └── index.html
├── base.css
├── tokens.css
├── favicon.svg
└── (other assets)
```

No:
- Build artifacts
- Source files
- Config files
- Server-side processing

The output is portable to any static hosting:
- GitHub Pages
- Cloudflare Pages
- Netlify
- Vercel
- rsync to a VPS
- AWS S3 + CloudFront

## Consequences

**Positive:**
- No platform lock-in
- Works on any HTTP server
- Simple deployment - just copy files
- Easy to understand what's being deployed
- CDN-friendly (immutable files)

**Negative:**
- No server-side processing (no redirects via config, etc.)
- Must rebuild to deploy changes

## Related

- Derived from: [Persona - Portable deployment goal](../persona.md#goals)

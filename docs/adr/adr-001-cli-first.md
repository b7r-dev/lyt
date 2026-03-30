# ADR-001: CLI-first Architecture

**Status**: Accepted

## Context

Martin, our target user, lives in the terminal. He's skeptical of web-based GUIs and wants every feature accessible from the command line. His pain points include:
- Needing to learn a templating language he'll only use for this one site
- Theme authors who abandon projects
- Mysterious build failures when upgrading dependencies

## Decision

lyt will be a CLI-only tool. All features are accessible via command-line flags and subcommands.

- `lyt build` - Build the site
- `lyt serve` - Dev server with hot reload
- `lyt import` - Import content from other formats (future)
- `lyt deploy` - Deploy to static host (future)

No web dashboard, no browser-based editor, no optional GUI. The binary is the product.

## Consequences

**Positive:**
- Martin gets exactly what he wants: a binary and a content folder
- No learning curve for a web UI
- Easy automation in CI/CD pipelines
- Predictable behavior - flags are the interface

**Negative:**
- Users expecting a web UI will be disappointed
- No visual preview without running a browser

## Related

- Derived from: [Persona - CLI ergonomics goal](../persona.md#goals)

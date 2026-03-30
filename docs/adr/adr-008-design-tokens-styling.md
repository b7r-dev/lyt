# ADR-008: Design Tokens as Styling Source

**Status**: Accepted

## Context

From the [AGENT_CONTEXT_0.md](../../AGENT_CONTEXT_0.md):
- Design tokens in `content/tokens.yaml` define colors, spacing, typography, z-planes
- Output is CSS custom properties (`--color-base-bg`, `--z-base`, etc.)

The design system uses:
- Warm beige palette
- 3 z-planes for depth
- CSS custom properties for all values
- Mobile-first, fluid typography

## Decision

**Design tokens are the single source of truth for styling.**

All visual design decisions are defined in `content/tokens.yaml`:
- Colors (base, accent, system)
- Spacing (0-32 scale)
- Typography (font families, sizes, weights, line heights)
- Z-indices (base, middle, top, modal)

The lyt engine processes these tokens into CSS custom properties that both:
1. `templates/base.css` uses directly
2. Any component rendering uses as reference

This ensures:
- Consistent styling across all pages
- Single place to change colors/spacing
- No hardcoded values in CSS or components

## Consequences

**Positive:**
- Single source of truth for design
- Easy to change themes globally
- Consistent component styling
- Tokens can be extended without code changes
- Enables future theme support

**Negative:**
- Must maintain tokens.yaml format
- Less granular control than inline styles
- Need to regenerate CSS when tokens change

## Related

- Implementation: [tokens/processor.go](../../internal/tokens/processor.go)
- Input: [content/tokens.yaml](../../content/tokens.yaml)
- Output: [dist/tokens.css](../../dist/tokens.css)

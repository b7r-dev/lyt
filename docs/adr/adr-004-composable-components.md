# ADR-004: Composable Component System

**Status**: Proposed

## Context

Martin's pain points include:
- Hardcoded theme components - can't add pull-quote or citation without editing theme templates
- Shortcode chaos - each theme defines its own shortcodes with different syntax
- No reusability - when he writes a good component, there's no way to use it in another post

His user stories include:
- Pull-quote component (`type: "pull-quote"`)
- Citation component (`type: "citation"`)
- Call-to-action component (`type: "cta"`)
- Warning banner component (`type: "warning"`)

## Decision

lyt will implement a **composable component system** where rich content elements are defined in YAML and rendered to consistent HTML.

### Core Components

| Component | YAML Type | Description |
|-----------|-----------|-------------|
| Pull-quote | `pull-quote` | Blockquote with attribution |
| Citation | `citation` | Footnote/reference with URL |
| CTA | `cta` | Call-to-action with button |
| Warning | `warning` | Styled alert/warning box |
| Hero | `hero` | Large title section |

### Implementation

1. Components are identified by `type` field in section YAML
2. Each type maps to a rendering function in the engine
3. Styles are provided by `base.css` - no per-component CSS needed
4. New components can be added by extending the renderer

```yaml
sections:
  - type: "pull-quote"
    quote: "The tools we use shape the work we produce."
    attribution: "Martin, 2024"
```

## Consequences

**Positive:**
- Consistent rendering across all posts
- No HTML copypasting between posts
- Easy to add new component types
- No theme modifications needed
- Content stays in YAML - readable and portable

**Negative:**
- Initial implementation effort
- Need to document all available components
- May need to balance flexibility vs. simplicity

## Related

- Derived from: [Persona - Composable templates stories](../persona.md#user-story-migrating-a-decade-of-content)
- See also: [ADR-003: YAML Schema](./adr-003-yaml-schema-content.md)

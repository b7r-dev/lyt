# ADR-002: Zero Runtime JavaScript

**Status**: Accepted

## Context

Martin's pain points include:
- Bloated JavaScript bundles in output he didn't ask for
- Zero `<script>` tags unless he explicitly added them in his markdown

The core principle of lyt is: "Engine, content, markup, and styling are kept strictly separate. They only meet at build time."

## Decision

lyt produces **zero runtime JavaScript** in the output. The browser receives only:
- HTML (semantic, readable)
- CSS (static files + generated design tokens)

Authors may add `<script>` tags manually in markdown content if needed, but lyt will never inject them.

## Consequences

**Positive:**
- Fast page loads, especially on slow connections
- No JavaScript security vulnerabilities in the output
- Works on any browser, including text-only browsers
- No hydration mismatch issues
- Simplified deployment - just static files

**Negative:**
- No client-side interactivity (tabs, accordions, etc.) without manual script tags
- Some modern UX patterns require JS

## Related

- Core principle: [No runtime JavaScript](../../AGENT_CONTEXT_0.md#project-overview)
- Derived from: [Persona - Zero JS goal](../persona.md#user-story-migrating-a-decade-of-content)

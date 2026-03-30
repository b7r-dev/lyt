# ADR-003: YAML Schema over Templating

**Status**: Accepted

## Context

Martin's pain points include:
- Configuration drift between theme versions
- Needing to learn a templating language he'll only use for this one site
- Theme authors who abandon projects, leaving him with security vulnerabilities

He wants: "Content structure is defined by YAML schema, not by tweaking theme templates."

## Decision

Content is defined declaratively via YAML. The schema is fixed (or extendable via composable components), not via arbitrary template logic.

```yaml
meta:
  title: "My Post"
  slug: "my-post"
  
sections:
  - id: "intro"
    type: "default"
    title: "Introduction"
    body: |
      Markdown content here.
```

The rendering is handled by lyt's engine, not by user-defined templates. This ensures:
- Consistent output across all pages
- No template version drift
- Security - no arbitrary template execution

## Consequences

**Positive:**
- No templating language to learn
- Consistent rendering everywhere
- Fast build - no template evaluation overhead
- Secure - no user code execution

**Negative:**
- Less flexible than arbitrary templating
- Users can't customize every aspect of HTML output
- May feel restrictive for power users

## Related

- Derived from: [Persona - Schema over configuration principle](../persona.md#design-principles-derived)
- See also: [ADR-004: Composable Components](./adr-004-composable-components.md)

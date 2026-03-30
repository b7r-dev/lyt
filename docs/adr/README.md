# Architecture Decision Records

> "A document that captures an important architectural decision made along with its context and consequences."

This directory contains ADRs for the lyt project, informed by the [user persona and story](../persona.md).

## Index

| ADR | Title | Status |
|-----|-------|--------|
| [ADR-001](./adr-001-cli-first.md) | CLI-first Architecture | Accepted |
| [ADR-002](./adr-002-zero-runtime-js.md) | Zero Runtime JavaScript | Accepted |
| [ADR-003](./adr-003-yaml-schema-content.md) | YAML Schema over Templating | Accepted |
| [ADR-004](./adr-004-composable-components.md) | Composable Component System | Proposed |
| [ADR-005](./adr-005-portable-static-output.md) | Portable Static Output | Accepted |
| [ADR-006](./adr-006-fast-incremental-builds.md) | Fast Incremental Builds | Proposed |
| [ADR-007](./adr-007-readable-html-output.md) | Readable HTML Output | Accepted |
| [ADR-008](./adr-008-design-tokens-styling.md) | Design Tokens as Styling Source | Accepted |

## Format

Each ADR follows the [Michael Nygard format](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions):

1. **Title** - Brief description
2. **Status** - Proposed, Accepted, Deprecated, or Superseded
3. **Context** - The situation forcing this decision
4. **Decision** - What we're doing
5. **Consequences** - Outcomes, both positive and negative

## Creating New ADRs

When making significant architectural choices:

1. Create a new file `adr-NNN-title.md`
2. Fill in all sections
3. Update this index
4. Commit with message "ADR-NNN: title"

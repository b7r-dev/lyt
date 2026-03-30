# ADR-006: Fast Incremental Builds

**Status**: Proposed

## Context

Martin's success criteria include:
- Build time (100 pages): < 2 seconds
- Serve rebuild time: < 500ms

His pain point is: "Slow iteration cycle: write → preview → wait → adjust"

## Decision

lyt will prioritize **fast build times** through:

1. **Incremental compilation** - Only rebuild changed content files
2. **Parallel processing** - Process independent pages concurrently  
3. **Efficient dependency tracking** - Track which content affects which output
4. **In-memory rendering** - Minimize I/O during development

### Targets

| Operation | Target |
|-----------|--------|
| Full build (100 pages) | < 2 seconds |
| Incremental rebuild | < 500ms |
| Cold start | < 3 seconds |

### Implementation Considerations

- Use file modification times to detect changes
- Cache parsed YAML and rendered markdown
- Parallelize content scanning and rendering
- Minimize allocations in hot paths

## Consequences

**Positive:**
- Enables flow state while writing
- Fast feedback loop
- Better developer experience
- Competitive advantage over slower SSGs

**Negative:**
- Complexity in incremental logic
- Potential for stale caches
- More code to maintain

## Related

- Derived from: [Persona - Fast iteration principle](../persona.md#design-principles-derived)

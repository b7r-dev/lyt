# ADR-007: Readable HTML Output

**Status**: Accepted

## Context

Martin's goals include:
- Clean, predictable output - he needs to understand exactly what HTML is being generated
- The ability to verify output by hand if needed

His story includes: "The HTML should be readable: indented, no minified inline scripts, semantic markup he can verify by hand if needed."

## Decision

Output HTML will be:
- **Indented** - Human-readable formatting
- **Semantic** - Uses `<article>`, `<main>`, `<nav>`, `<section>`, etc.
- **Non-minified** - No inline minified scripts or styles
- **Traceable** - Clear relationship between input YAML and output HTML

Example output:
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>My Post | lyt</title>
  <link rel="stylesheet" href="/base.css">
</head>
<body>
  <nav class="nav">...</nav>
  <main class="main">
    <article class="blog-post">
      <h1 class="post-title">My Post</h1>
      <div class="post-body">
        ...
      </div>
    </article>
  </main>
</body>
</html>
```

## Consequences

**Positive:**
- Easy to verify output is correct
- Useful for debugging
- Can teach HTML structure by example
- Makes git diffs meaningful
- No "black box" output

**Negative:**
- Slightly larger file sizes (indentation)
- Requires proper templating (not inline string concatenation)

## Related

- Derived from: [Persona - Transparent output principle](../persona.md#design-principles-derived)

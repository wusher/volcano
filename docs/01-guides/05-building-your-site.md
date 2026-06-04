# Building Your Site

Generate static HTML from your markdown files.

`--url` is required for every build (it's used for canonical and Open Graph
tags). Either pass it on the CLI or set it once in `volcano.json` — the
shortened examples below assume one or the other.

## Basic Build

```bash
volcano ./docs --url="https://docs.example.com"
```

This writes the site to `./output`.

## Custom Output Directory

```bash
volcano ./docs -o ./public --url="https://docs.example.com"
```

## With Metadata

```bash
volcano ./docs \
  --url="https://docs.example.com" \
  --title="My Documentation" \
  --author="Your Name"
```

## Build Options

The flag examples below show only the option being demonstrated. Add `--url` (or a config file) when you actually run them.

### Themes

```bash
# Use blog theme
volcano ./docs --theme blog

# Use custom CSS
volcano ./docs --css ./my-theme.css
```

### Navigation Features

```bash
# Hover prefetching for fast clicks
volcano ./docs --instant-nav

# Cmd+K command palette
volcano ./docs --search

# Previous/next links at bottom of pages
volcano ./docs --page-nav

# Horizontal bar with root pages
volcano ./docs --top-nav

# Hierarchy trail above each page
volcano ./docs --breadcrumbs
```

### Advanced Features

```bash
# Installable site with offline support
volcano ./docs --pwa

# Custom accent color
volcano ./docs --accent-color="#0066cc"
```

## Complete Example

```bash
volcano ./docs \
  -o ./public \
  --title="My Documentation" \
  --url="https://docs.example.com" \
  --theme docs \
  --instant-nav \
  --search \
  --pwa \
  --favicon="./icon.png"
```

## Build Output

After building, your output directory contains:

```
public/
├── index.html           # Homepage
├── getting-started/
│   └── index.html       # Clean URLs (no .html in path)
├── guides/
│   ├── index.html
│   └── advanced/
│       └── index.html
├── styles.css           # Combined CSS
├── search-index.json    # If --search enabled
├── manifest.json        # If --pwa enabled
├── sw.js               # If --pwa enabled
└── 404.html            # Custom 404 page
```

## What Gets Processed

### Included

- `.md` files → HTML pages
- Images referenced in markdown
- Favicon (if specified)
- Automatic navigation tree
- Auto-generated 404 page

### Skipped

- Hidden files (starting with `.` or `_`)
- Dotfiles (`.gitignore`, `.DS_Store`, etc.)
- Non-markdown files not referenced

## Link Validation

Volcano automatically validates all internal links during the build. If broken links are found, the build fails:

```
✗ Found 2 broken internal links:
  Page /guides/intro/: broken link /setup/ (not found)
  Page /reference/api/: broken link /deprecated/ (not found)
```

Fix the links and rebuild.

### Bypass Link Validation

For work-in-progress sites:

```bash
volcano ./docs --url="https://example.com" --allow-broken-links
```

This warns about broken links but doesn't fail the build.

## Quiet and Verbose Modes

### Quiet Mode

Suppress output (useful for scripts):

```bash
volcano ./docs --url="$URL" -q
```

### Verbose Mode

See detailed processing information:

```bash
volcano ./docs --url="$URL" --verbose
```

Shows:
- Files being processed
- Navigation tree structure
- Link resolution details

## Configuration File

Create `volcano.json` in your input directory:

```json
{
  "title": "My Documentation",
  "output": "./public",
  "theme": "docs",
  "instantNav": true,
  "search": true,
  "pwa": true,
  "url": "https://docs.example.com"
}
```

Then build with:

```bash
volcano ./docs
```

CLI flags override config file values.

## When to Rebuild

You need to rebuild when:
- You change markdown content
- You add or remove pages
- You update site metadata (title, theme, etc.)
- You modify custom CSS

## Next Steps

- [[serving-your-site]] — Preview and test your build
- [[deploying-your-site]] — Deploy to production
- [[reference/cli]] — All build options

# Command Line Reference

Complete reference for the Volcano CLI.

## Commands

### Generate

Generate a static site from markdown files:

```bash
volcano <input-folder> [flags]
volcano build <input-folder> [flags]
```

Both commands do the same thing.

**Example:**

```bash
volcano ./docs -o ./public --title="My Documentation"
```

### Serve

Serve your markdown source directory for local preview:

```bash
volcano serve <folder> [flags]
```

**Example:**

```bash
volcano serve ./docs -p 8080
```

Pages are rendered on each request. Changes appear on browser refresh.

### CSS Export

Export the vanilla CSS skeleton for customization:

```bash
volcano css [-o file]
```

**To stdout:**

```bash
volcano css
```

**To file:**

```bash
volcano css -o my-theme.css
```

## Flags

### Output

| Flag | Default | Description |
|------|---------|-------------|
| `-o, --output` | `./output` | Output directory |

### Server

Use `volcano serve <folder>` instead of the build command:

| Flag | Default | Description |
|------|---------|-------------|
| `-p, --port` | `1776` | Server port |

### Site Metadata

| Flag | Default | Description |
|------|---------|-------------|
| `--title` | `My Site` | Site title for header and meta tags |
| `--url` | (none) | Base URL for canonical links and SEO |
| `--author` | (none) | Author for meta tags |
| `--og-image` | (none) | Default Open Graph image URL |
| `--favicon` | (none) | Path to favicon file (.ico, .png, .svg) |

For subpath deployments, include the path in `--url`:

```bash
volcano ./docs --url="https://username.github.io/my-repo"
```

### Navigation

| Flag | Default | Description |
|------|---------|-------------|
| `--top-nav` | `false` | Show root files in top navigation bar |
| `--page-nav` | `false` | Show previous/next page links |
| `--breadcrumbs` | `true` | Show breadcrumb trail |
| `--instant-nav` | `false` | Enable hover prefetching |

### Theming

| Flag | Default | Description |
|------|---------|-------------|
| `--theme` | `docs` | Built-in theme: `docs`, `blog`, `vanilla` |
| `--css` | (none) | Custom CSS file (overrides theme) |
| `--accent-color` | (none) | Custom accent color (hex, e.g., `#ff6600`) |

### Advanced

| Flag | Default | Description |
|------|---------|-------------|
| `--search` | `false` | Enable site search with Cmd+K command palette (see [[features/navigation#search\|Navigation]]) |
| `--pwa` | `false` | Enable PWA manifest and service worker (see [[features/progressive-web-app\|PWA]]) |
| `--inline-assets` | `false` | Embed CSS/JS inline instead of external files |
| `--allow-broken-links` | `false` | Don't fail build on broken internal links |

### Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-c, --config` | (none) | Path to config file |

Place a `volcano.json` in your input directory to set defaults:

```json
{
  "title": "My Docs",
  "output": "./public",
  "theme": "docs",
  "search": true,
  "pwa": true
}
```

CLI flags override config file values.

### Output Control

| Flag | Default | Description |
|------|---------|-------------|
| `-q, --quiet` | `false` | Suppress non-error output |
| `--verbose` | `false` | Enable debug output |

### Information

| Flag | Description |
|------|-------------|
| `-v, --version` | Show version |
| `-h, --help` | Show help |

## Examples

### Basic Generation

```bash
volcano ./docs
```

Generates to `./output` with default settings.

### Production Build

```bash
volcano ./docs \
  -o ./public \
  --title="My Project" \
  --url="https://docs.example.com" \
  --author="My Team" \
  --og-image="https://docs.example.com/og.png"
```

### Development Server

```bash
volcano serve -p 3000 ./docs
```

Serves with dynamic regeneration at `http://localhost:3000`.

### Blog Setup

```bash
volcano ./posts \
  -o ./blog \
  --theme blog \
  --title="My Blog" \
  --page-nav
```

### Custom Theme

```bash
# Export vanilla CSS
volcano css -o my-theme.css

# Edit my-theme.css...

# Use custom CSS
volcano ./docs --css ./my-theme.css
```

### Full Navigation

```bash
volcano ./docs \
  --top-nav \
  --page-nav \
  --breadcrumbs \
  --instant-nav \
  --search
```

### CI/CD Build

```bash
volcano ./docs -o ./public -q --url="$SITE_URL"
```

Quiet mode for cleaner CI logs.

## Link Validation

Volcano automatically validates all internal links during generation. The build will **fail** if any broken links are found.

### What's Validated

1. **Navigation links** — All sidebar navigation links are verified
2. **Content links** — Internal links within page content are checked
3. **Wiki links** — `[[Page Name]]` links are resolved and validated

### Error Output

When broken links are found:

```
✗ Found 2 broken internal links:
  Page /guides/intro/: broken link /setup/ (not found)
  Page /reference/api/: broken link /deprecated/ (not found)
```

### Fixing Broken Links

1. Check the link target exists as a markdown file
2. Verify the path is correct (paths are case-insensitive)
3. For wiki links, ensure the page name matches an existing file

### Dynamic Server

When using `volcano serve ./docs` (dynamic serving), broken links are shown inline on the page with detailed error messages instead of failing silently.

## Generated Files

In addition to your pages, Volcano automatically generates:

### 404 Page

A `404.html` file is created in the output directory for custom error handling. Most web servers and hosting platforms (GitHub Pages, Netlify, Vercel) will automatically serve this page for missing URLs.

The 404 page includes:
- Full site navigation (so users can find their way)
- Consistent styling with your theme
- A simple "Page Not Found" message

### Styles

A `styles.css` file containing the combined and minified theme CSS.

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Error (invalid args, missing files, generation failure, broken links) |

## Environment

Volcano detects terminal capabilities automatically:
- Colored output when stderr is a TTY
- Progress indicators in interactive mode

## Front Matter Support

Volcano strips YAML front matter from markdown files during rendering:

```markdown
---
title: My Page
date: 2024-01-15
---

# Page Content
```

Front matter is removed so it doesn't appear in the output. While front matter values aren't currently used, this ensures compatibility with Obsidian, Hugo, Jekyll, and other tools that use front matter.

**Note:** Page titles come from H1 headings, not front matter `title` fields.

## Related

- [[url-structure]] — URL generation rules

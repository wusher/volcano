# Command Line Reference

Complete reference for the Volcano CLI.

## Commands

### Generate

Generate a static site from markdown files:

```bash
volcano <input-folder> [flags]
```

**Example:**

```bash
volcano ./docs -o ./public --title="My Documentation"
```

### Serve

Serve a generated site or regenerate on request:

```bash
volcano -s <folder> [flags]
```

**Static serving (pre-generated):**

```bash
volcano -s ./public -p 8080
```

**Dynamic regeneration (source folder):**

```bash
volcano -s ./docs -p 8080
```

When serving a source folder, pages are regenerated on each request.

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

### Output Options

#### `-o, --output <dir>`

Output directory for generated HTML files.

- **Default:** `./output`
- **Type:** String (directory path)

```bash
volcano ./docs -o ./public
volcano ./docs --output ./dist
```

### Server Options

#### `-s, --serve`

Run in serve mode instead of generating.

- **Default:** `false`
- **Type:** Boolean

```bash
volcano -s ./public
volcano --serve ./docs
```

#### `-p, --port <port>`

Port for the HTTP server.

- **Default:** `1776`
- **Type:** Integer

```bash
volcano -s -p 8080 ./public
volcano -s --port 3000 ./docs
```

### Site Metadata

#### `--title <title>`

Site title displayed in the header and used in meta tags.

- **Default:** `My Site`
- **Type:** String

```bash
volcano ./docs --title="My Project Documentation"
```

#### `--url <url>`

Base URL for the site. Used for:
- Canonical link tags
- Open Graph URLs
- Subpath prefix for all internal links
- Sitemap generation

- **Default:** (none)
- **Type:** URL string

```bash
volcano ./docs --url="https://docs.example.com"
```

**Subpath Deployments:**

When deploying to a subpath (e.g., GitHub Pages project sites), include the path:

```bash
volcano ./docs --url="https://username.github.io/my-repo"
```

All internal links will be prefixed with `/my-repo/` automatically.

#### `--author <name>`

Site author for meta tags.

- **Default:** (none)
- **Type:** String

```bash
volcano ./docs --author="Jane Smith"
```

#### `--og-image <url>`

Default Open Graph image URL for social sharing.

- **Default:** (none)
- **Type:** URL string

```bash
volcano ./docs --og-image="https://example.com/og-image.png"
```

#### `--favicon <path>`

Path to favicon file. Supports `.ico`, `.png`, and `.svg` formats.

- **Default:** (none)
- **Type:** File path

```bash
volcano ./docs --favicon="./favicon.ico"
```

### Navigation Options

#### `--top-nav`

Display root-level files in a top navigation bar instead of the sidebar.

- **Default:** `false`
- **Type:** Boolean

```bash
volcano ./docs --top-nav
```

#### `--page-nav`

Show previous/next page navigation at the bottom of each page.

- **Default:** `false`
- **Type:** Boolean

```bash
volcano ./docs --page-nav
```

#### `--breadcrumbs`

Show breadcrumb trail navigation at the top of each page. Breadcrumbs help users understand their location in the site hierarchy and provide quick navigation to parent pages.

- **Default:** `true`
- **Type:** Boolean

```bash
# Enable breadcrumbs (default)
volcano ./docs --breadcrumbs

# Disable breadcrumbs
volcano ./docs --breadcrumbs=false
```

#### `--instant-nav`

Enable instant navigation with hover prefetching. When enabled, pages are prefetched when users hover over links, making navigation feel instantaneous.

- **Default:** `false`
- **Type:** Boolean

```bash
volcano ./docs --instant-nav
```

#### `--last-modified`

Show last modified date on pages.

- **Default:** `false`
- **Type:** Boolean

```bash
volcano ./docs --last-modified
```

### Theming

#### `--theme <name>`

Built-in theme to use.

- **Default:** `docs`
- **Options:** `docs`, `blog`, `vanilla`

```bash
volcano ./docs --theme blog
volcano ./docs --theme vanilla
```

#### `--css <path>`

Path to custom CSS file. When specified, overrides the `--theme` flag.

- **Default:** (none)
- **Type:** File path

```bash
volcano ./docs --css ./my-theme.css
```

#### `--accent-color <hex>`

Custom accent color in hex format. This overrides the theme's default accent color for links, buttons, and other UI elements.

- **Default:** (none, uses theme default)
- **Type:** Hex color string

```bash
volcano ./docs --accent-color="#ff6600"
volcano ./docs --accent-color="#3b82f6"
```

### Output Control

#### `-q, --quiet`

Suppress non-error output.

- **Default:** `false`
- **Type:** Boolean

```bash
volcano ./docs -q
volcano ./docs --quiet
```

#### `--verbose`

Enable debug output for troubleshooting.

- **Default:** `false`
- **Type:** Boolean

```bash
volcano ./docs --verbose
```

### Information

#### `-v, --version`

Show version information.

```bash
volcano -v
volcano --version
```

Output: `volcano version 0.1.0`

#### `-h, --help`

Show help message.

```bash
volcano -h
volcano --help
```

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
volcano -s -p 3000 ./docs
```

Serves with dynamic regeneration at `http://localhost:3000`.

### Blog Setup

```bash
volcano ./posts \
  -o ./blog \
  --theme blog \
  --title="My Blog" \
  --last-modified
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
  --last-modified
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

When using `volcano -s ./docs` (dynamic serving), broken links are shown inline on the page with detailed error messages instead of failing silently.

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

## Related

- [[url-structure]] — URL generation rules
- [[front-matter]] — YAML front matter support
- [[guides/installation]] — Installation methods

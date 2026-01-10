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
- Sitemap generation

- **Default:** (none)
- **Type:** URL string

```bash
volcano ./docs --url="https://docs.example.com"
```

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
  --last-modified
```

### CI/CD Build

```bash
volcano ./docs -o ./public -q --url="$SITE_URL"
```

Quiet mode for cleaner CI logs.

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Error (invalid args, missing files, generation failure) |

## Environment

Volcano detects terminal capabilities automatically:
- Colored output when stderr is a TTY
- Progress indicators in interactive mode

## Related

- [[url-structure]] — URL generation rules
- [[front-matter]] — YAML front matter support
- [[guides/installation]] — Installation methods

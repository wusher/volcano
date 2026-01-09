# Volcano

[![CI](https://github.com/wusher/volcano/actions/workflows/ci.yml/badge.svg)](https://github.com/wusher/volcano/actions/workflows/ci.yml)

A Go CLI static site generator that converts a folder of markdown files into a styled static website with a tree navigation layout.

## Features

### Core Features
- **No frontmatter** - Simple markdown files, no YAML required
- **Alphabetical ordering** by filename
- **Clean URLs** - `guides/intro.md` → `/guides/intro/`
- **Tree navigation** - Collapsible folder structure in sidebar
- **Light/dark mode** - With browser preference detection
- **Responsive design** - Desktop sidebar, mobile drawer
- **Embedded styling** - No external dependencies
- **Fast builds** - Simple, efficient Go implementation

### Navigation & UX (Stories 14-20)
- **Table of Contents** - Auto-generated from h2-h4 headings with scroll spy highlighting
- **Breadcrumb Navigation** - With schema.org structured data for SEO
- **Previous/Next Navigation** - Page-level navigation in depth-first order
- **Heading Anchor Links** - Clickable anchors with unique IDs for all headings
- **External Link Indicators** - Visual icon and `target="_blank"` for external links
- **Code Block Copy Button** - One-click copy to clipboard for code blocks
- **Keyboard Shortcuts** - Press `?` to see all shortcuts (/, t, n, p, h)

### Display Features (Stories 21-25)
- **Print Stylesheet** - Optimized print layout hiding navigation
- **Reading Time** - Estimated reading time based on word count
- **Last Modified Date** - Shows when content was last updated (git or filesystem)
- **Scroll Progress** - Visual indicator of page scroll position
- **Back to Top Button** - Smooth scroll back to page top

### SEO & Meta (Stories 26-28)
- **SEO Meta Tags** - Description, robots, author, canonical URL
- **Open Graph Support** - Full og:title, og:description, og:image, etc.
- **Custom Favicon** - Support for .ico, .png, and .svg favicons

### Content Features (Stories 29-31)
- **Admonition Blocks** - Note, tip, warning, danger, info callout blocks
- **Code Line Highlighting** - Highlight specific lines in code blocks
- **Smooth Scroll** - Smooth scrolling with reduced-motion support

### Advanced Navigation (Stories 32-37)
- **Clickable Folders** - Folders with index files are clickable in navigation
- **Navigation Search** - Filter navigation tree by typing
- **Top Navigation Bar** - Optional horizontal nav for root-level files
- **Auto-Generated Index** - Folders without index.md get automatic listings
- **H1-Based Labels** - Navigation uses H1 titles instead of filenames
- **Filename Prefixes** - Date (2024-01-01-) and number (01-) prefixes stripped from URLs

## Installation

```bash
go install github.com/wusher/volcano@latest
```

Or build from source:

```bash
git clone https://github.com/wusher/volcano.git
cd volcano
go build -o volcano .
```

## Usage

### Generate a static site

```bash
volcano <input-folder> -o <output-folder> --title="My Site"
```

### Serve a static site

```bash
volcano -s -p 1776 <folder>
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Output directory | `./output` |
| `-s, --serve` | Run in serve mode | `false` |
| `-p, --port` | Server port | `1776` |
| `--title` | Site title | `My Site` |
| `--url` | Site base URL (for SEO) | |
| `--author` | Site author | |
| `--og-image` | Default Open Graph image URL | |
| `--favicon` | Path to favicon file (.ico, .png, .svg) | |
| `--last-modified` | Show last modified date on pages | `false` |
| `--top-nav` | Display root files in top navigation bar | `false` |
| `-q, --quiet` | Suppress non-error output | `false` |
| `--verbose` | Enable debug output | `false` |
| `-v, --version` | Show version | |
| `-h, --help` | Show help | |

## Markdown Features

### Admonition Blocks

Use fenced blocks to create callouts:

```markdown
:::note
This is a note callout.
:::

:::tip
This is a tip callout.
:::

:::warning
This is a warning callout.
:::

:::danger
This is a danger callout.
:::

:::info
This is an info callout.
:::
```

### Code Blocks with Line Highlighting

Highlight specific lines in code blocks using the `{lines}` annotation:

````markdown
```go {3-5}
func main() {
    // These lines
    // will be
    // highlighted
    fmt.Println("Hello")
}
```
````

### Filename Prefixes

Files with date or number prefixes are automatically cleaned:

- `2024-01-15-hello-world.md` → `/hello-world/` (date stripped)
- `01-introduction.md` → `/introduction/` (number stripped)
- `_draft.md` → skipped (hidden file)

## Examples

Generate site from `./docs` to `./public`:

```bash
volcano ./docs -o ./public --title="My Documentation"
```

Generate with SEO options:

```bash
volcano ./docs -o ./public \
  --title="My Docs" \
  --url="https://mydocs.example.com" \
  --author="Your Name" \
  --og-image="https://mydocs.example.com/og-image.png" \
  --favicon="./favicon.ico"
```

Generate with top navigation bar:

```bash
volcano ./docs -o ./public --title="My Docs" --top-nav
```

Serve the generated site:

```bash
volcano -s -p 8080 ./public
```

## Example Site

The `example/` folder contains sample markdown content demonstrating all features:

```bash
# Generate the example site
volcano ./example -o ./output --title="Volcano Docs"

# Serve it locally
volcano -s ./output

# Open http://localhost:1776 in your browser
```

## Development

### Prerequisites

- Go 1.21 or later
- golangci-lint (for linting)

### Running tests

```bash
# Run all tests
go test -race ./...

# Run with verbose output
go test -v ./...

# Run integration tests only
go test -v -run TestIntegration ./...

# Run story-specific tests
go test -v -run "TestIntegrationStory" ./...
```

### Check coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
```

### Running lint

```bash
golangci-lint run
```

### Full quality check

```bash
# Lint, test, and verify coverage
golangci-lint run && \
go test -race ./... && \
go test -coverprofile=coverage.out ./... && \
go tool cover -func=coverage.out | grep total
```

## Project Structure

```
volcano/
├── main.go                  # CLI entry point
├── cmd/                     # Command implementations
├── internal/
│   ├── assets/              # Favicon handling
│   ├── content/             # Reading time, last modified
│   ├── generator/           # Site generation engine
│   ├── markdown/            # Markdown parsing, admonitions, headings
│   ├── navigation/          # Breadcrumbs, pagination
│   ├── output/              # Colored logging
│   ├── seo/                 # Meta tags, Open Graph
│   ├── server/              # HTTP server
│   ├── styles/              # Embedded CSS
│   ├── templates/           # HTML templates
│   ├── toc/                 # Table of contents
│   └── tree/                # File tree building, scanning
├── example/                 # Example content
└── integration_test.go      # Integration tests (24 story tests)
```

## License

MIT

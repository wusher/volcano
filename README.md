# Volcano

![Volcano Logo](docs/logo.png)

**Turn your markdown folders into websites. No config required.**

[![CI](https://github.com/wusher/volcano/actions/workflows/ci.yml/badge.svg)](https://github.com/wusher/volcano/actions/workflows/ci.yml)

---

Volcano is an opinionated static site generator for people who just want their markdown files on the web. Point it at a folder—your Obsidian vault, your notes directory, your documentation—and get a beautiful, navigable website. No configuration files. No frontmatter. No build pipelines. Just markdown in, website out.

## Why Volcano?

**Your folder structure is your site structure.** Volcano reads your directories and creates matching navigation automatically. No need to define routes, menus, or page hierarchies in config files.

**Zero configuration by design.** Other generators ask you to learn their templating language, configure plugins, and maintain YAML files. Volcano has one command: point it at markdown, get a website.

**Works with your existing notes.** Have an Obsidian vault? A folder of documentation? Meeting notes organized by date? Volcano handles them all without requiring you to restructure anything.

**Single binary, no dependencies.** No Node.js, Python, or Ruby. Download one file and run it. Works offline, builds fast, deploys anywhere.

## Features

- **Tree navigation** — Collapsible sidebar mirrors your folder structure
- **Instant navigation** — Hover prefetching and smooth page transitions
- **Search** — Command palette (Cmd+K) searches pages and headings
- **Table of contents** — Auto-generated from headings with scroll tracking
- **Dark mode** — Automatic detection with manual toggle
- **Wiki links** — Obsidian-style `[[Page Name]]` linking
- **Admonitions** — Note, tip, warning, and info callout blocks
- **Code highlighting** — Syntax highlighting with copy button
- **SEO ready** — Meta tags, Open Graph, and automatic sitemaps
- **Keyboard shortcuts** — Press `?` to see all navigation shortcuts
- **PWA support** — Progressive Web App for offline access
- **Reading time** — Estimated time displayed on each page
- **Breadcrumbs** — Always know where you are in the hierarchy
- **Mobile responsive** — Drawer navigation on small screens

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

## Releases

Tag a release to publish GitHub release binaries:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

The release workflow builds cross-platform binaries and attaches them to the GitHub release.

## Usage

### Generate a static site

```bash
volcano <input-folder> -o <output-folder> --title="My Site"
```

### Preview locally

```bash
volcano serve ./docs
```

### Common Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Output directory | `./output` |
| `-s, --serve` | Run in serve mode | `false` |
| `-p, --port` | Server port | `1776` |
| `--title` | Site title | `My Site` |
| `--url` | Site base URL (for SEO) | |
| `--theme` | Built-in theme: docs, blog, vanilla | `docs` |
| `--accent-color` | Custom accent color (hex) | |
| `--instant-nav` | Enable hover prefetching | `false` |
| `--search` | Enable search (Cmd+K) | `false` |
| `--pwa` | Enable PWA support | `false` |
| `--top-nav` | Display root files in top nav bar | `false` |
| `--page-nav` | Show previous/next page links | `false` |
| `-q, --quiet` | Suppress non-error output | `false` |
| `-v, --version` | Show version | |
| `-h, --help` | Show help | |

Run `volcano --help` or `volcano serve --help` for all options.

## Key Features

### Wiki Links

Link between pages using `[[Page Name]]` syntax:

```markdown
See the [[Installation]] guide for setup instructions.
Check out [[Advanced/Configuration]] for more options.
```

### Instant Navigation

Enable hover prefetching for near-instant page loads:

```bash
volcano ./docs --instant-nav --search
```

Pages prefetch on hover and use smooth view transitions for seamless navigation.

### Search

Add a command palette with `--search`:

```bash
volcano ./docs --search
```

Press Cmd+K (or Ctrl+K) to search pages and headings across your entire site.

### Filename Organization

Files with date or number prefixes are automatically cleaned:

- `2024-01-15-hello-world.md` → `/hello-world/` (date stripped)
- `01-introduction.md` → `/introduction/` (number stripped)
- `_draft.md` → skipped (hidden file)

## Examples

### Basic Documentation

```bash
volcano ./docs --title="My Documentation"
```

### Full-Featured Site

```bash
volcano ./docs \
  --title="My Docs" \
  --theme docs \
  --instant-nav \
  --search \
  --url="https://docs.example.com"
```

### Blog Setup

```bash
volcano ./posts \
  --theme blog \
  --page-nav \
  --title="My Blog"
```

### Development Server

```bash
volcano serve ./docs
```

Changes appear on browser refresh.

## Example Site

The `example/` folder contains sample markdown content demonstrating all features:

```bash
# Preview the example site
volcano serve ./example --title="Volcano Docs"

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
│   ├── content/             # Reading time calculation
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

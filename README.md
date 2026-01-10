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
- **Clean URLs** — `guides/intro.md` → `/guides/intro/`
- **Table of contents** — Auto-generated from headings with scroll tracking
- **Dark mode** — Automatic detection with manual toggle
- **Admonitions** — Note, tip, warning, and info callout blocks
- **Code highlighting** — Syntax highlighting with copy button
- **SEO ready** — Meta tags, Open Graph, and sitemaps
- **Keyboard shortcuts** — Press `?` to see all navigation shortcuts
- **Reading time** — Estimated time displayed on each page
- **Breadcrumbs** — Always know where you are in the hierarchy
- **Search** — Filter navigation by typing
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

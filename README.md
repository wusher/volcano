# Volcano

[![CI](https://github.com/wusher/volcano/actions/workflows/ci.yml/badge.svg)](https://github.com/wusher/volcano/actions/workflows/ci.yml)

A Go CLI static site generator that converts a folder of markdown files into a styled static website with a tree navigation layout.

## Features

- **No frontmatter** - Simple markdown files, no YAML required
- **Alphabetical ordering** by filename
- **Clean URLs** - `guides/intro.md` → `/guides/intro/`
- **Tree navigation** - Collapsible folder structure in sidebar
- **Light/dark mode** - With browser preference detection
- **Responsive design** - Desktop sidebar, mobile drawer
- **Embedded styling** - No external dependencies
- **Fast builds** - Simple, efficient Go implementation

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
| `-q, --quiet` | Suppress non-error output | `false` |
| `--verbose` | Enable debug output | `false` |
| `-v, --version` | Show version | |
| `-h, --help` | Show help | |

## Examples

Generate site from `./docs` to `./public`:

```bash
volcano ./docs -o ./public --title="My Documentation"
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
│   ├── generator/           # Site generation engine
│   ├── markdown/            # Markdown parsing
│   ├── templates/           # HTML templates
│   ├── tree/                # File tree building
│   ├── styles/              # Embedded CSS
│   ├── server/              # HTTP server
│   └── output/              # Colored logging
├── example/                 # Example content
└── integration_test.go      # Integration tests
```

## License

MIT

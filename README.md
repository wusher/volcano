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

## Development

### Prerequisites

- Go 1.21 or later
- golangci-lint (for linting)

### Running tests

```bash
go test -race ./...
```

### Running lint

```bash
golangci-lint run
```

### Check all (lint, test, coverage)

```bash
./scripts/check.sh
```

## License

MIT

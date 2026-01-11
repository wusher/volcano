# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
# Build
go build -o volcano .

# Run all tests with race detection
go test -race ./...

# Run a single test
go test -v -run TestFunctionName ./path/to/package

# Run integration tests only
go test -v -run TestIntegration ./...

# Run story-specific tests (e2e tests for features)
go test -v -run "TestIntegrationStory" ./...

# Check coverage
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep total

# Lint
golangci-lint run

# Full quality check (lint + test + coverage)
golangci-lint run && go test -race ./... && go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep total
```

## CLI Usage

```bash
# Generate static site from markdown folder
volcano ./docs -o ./public --title="My Docs"

# Serve a static site
volcano -s -p 8080 ./public

# Use example/ folder for testing
volcano ./example -o ./output --title="Volcano Docs"
```

## Architecture

Volcano is a zero-dependency static site generator that converts markdown folders to styled HTML sites with tree navigation.

### Core Data Flow

1. **tree.Scan()** walks input directory → builds `*tree.Node` tree + collects `AllPages` slice
2. **generator.Generate()** iterates pages:
   - Reads markdown, applies transforms (admonitions, anchors, external links, code blocks)
   - Builds navigation, breadcrumbs, TOC, SEO meta
   - Renders via templates.Renderer with embedded layout.html
3. Output: clean URLs (`file.md` → `file/index.html`)

### Key Packages

- **cmd/**: Config struct and command handlers (Generate, Serve)
- **internal/tree/**: File scanning, Node tree structure, URL/path generation, filename metadata extraction (date/number prefix stripping, slugification)
- **internal/generator/**: Orchestrates the build pipeline, auto-index generation
- **internal/templates/**: HTML layout (embedded), PageData struct, navigation rendering
- **internal/markdown/**: Goldmark parser, content transforms (admonitions, headings, links, code blocks)
- **internal/styles/**: Embedded CSS themes (docs, blog, vanilla)
- **internal/server/**: Static file server with dynamic regeneration mode

### Tree Node Structure

```go
type Node struct {
    Name       string   // Display name (H1 title or cleaned filename)
    Path       string   // Relative path from input root
    SourcePath string   // Absolute filesystem path
    IsFolder   bool
    HasIndex   bool     // Folder has index.md
    Children   []*Node
}
```

### URL Generation

- `tree.GetOutputPath(node)` → filesystem output path
- `tree.GetURLPath(node)` → URL for links
- Date prefixes (`2024-01-15-`) and number prefixes (`01-`) are stripped from slugs
- Directory paths are slugified (spaces/special chars → lowercase-hyphenated)

### Embedded Assets

Templates and CSS use Go's `//go:embed` directive:
- `internal/templates/layout.html` - main HTML structure
- `internal/styles/themes/*.css` - theme stylesheets

### Test Fixtures

When fixing wiki link or markdown processing bugs, always store test fixtures in the `testdata/` folder:

1. Create a folder under `testdata/` with markdown files that reproduce the issue
2. Add an integration test that references the testdata folder
3. Update `testdata/README.md` to document the test scenario

This keeps test fixtures as real files (easier to inspect/debug) rather than inline strings in test code.

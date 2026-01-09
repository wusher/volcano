# Volcano - Static Site Generator Implementation Plan

## Overview

Volcano is a Go CLI static site generator that converts a folder of markdown files into a styled static website with a tree navigation layout.

### Core Commands

```bash
# Generate static site
volcano FOLDER_NAME -o OUTPUT_FOLDER --title="My Site"

# Serve static site
volcano -s -p 1776 FOLDER_NAME
```

### Key Design Decisions

- **No frontmatter** (future enhancement)
- **Alphabetical ordering** by filename
- **No live reload** - simple static file serving
- **Ignore non-markdown assets** - only process `.md` files
- **Embedded Tailwind CSS** - no external dependencies
- **Clean nav labels** - `getting-started.md` → "Getting Started"
- **Collapsible folders** in navigation tree
- **Light/dark mode** with browser preference default
- **Colors**: shades of black and white only

---

## Development Workflow Rules

**These rules MUST be followed for every story:**

### Rule 1: Commit & Push After Each Story

After completing each story:
1. Stage all changes: `git add .`
2. Commit with story reference: `git commit -m "Story N: <brief description>"`
3. Push to remote: `git push`

Example:
```bash
git add .
git commit -m "Story 2: Add CI pipeline with lint, test, coverage, and e2e"
git push
```

### Rule 2: Pass All Checks Before Starting Next Story

Before starting the next story, **ALL** of the following must pass:

```bash
# Run locally before moving on
golangci-lint run          # Linting must pass with zero errors
go test -race ./...        # All tests must pass
go test -coverprofile=coverage.out ./... && \
  go tool cover -func=coverage.out | grep total  # Coverage must be ≥ 70%
```

**If any check fails:**
1. **STOP** - Do not start the next story
2. Fix all lint errors
3. Fix all failing tests
4. Add tests if coverage dropped below 70%
5. Commit the fixes: `git commit -m "Story N: Fix lint/test/coverage issues"`
6. Push: `git push`
7. Re-run all checks to confirm they pass
8. Only then proceed to the next story

### Quick Check Script

Create this as `scripts/check.sh` for convenience:

```bash
#!/bin/bash
set -e

echo "Running lint..."
golangci-lint run

echo "Running tests with race detector..."
go test -race ./...

echo "Checking coverage..."
go test -coverprofile=coverage.out ./...
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "Total coverage: $COVERAGE%"

if (( $(echo "$COVERAGE < 70" | bc -l) )); then
  echo "FAIL: Coverage $COVERAGE% is below 70% threshold"
  exit 1
fi

echo ""
echo "All checks passed! Ready for next story."
```

Run with: `./scripts/check.sh`

---

## Story 1: Project Scaffolding & CLI Foundation

**Goal**: Set up the Go project structure with basic CLI argument parsing.

### Acceptance Criteria

- [ ] Initialize Go module as `volcano`
- [ ] Create main.go entry point
- [ ] Parse CLI arguments without external dependencies (use `flag` package)
- [ ] Support two modes:
  - **Generate mode** (default): `volcano <input-folder> -o <output-folder>`
  - **Serve mode**: `volcano -s -p <port> <folder>`
- [ ] Implement flags:
  - `-o, --output` - output directory (default: `./output`)
  - `-s, --serve` - enable serve mode
  - `-p, --port` - server port (default: `1776`)
  - `--title` - site title (default: `"My Site"`)
- [ ] Display help text with `-h` or `--help`
- [ ] Display version with `-v` or `--version`
- [ ] Validate input folder exists and is a directory
- [ ] Exit with appropriate error codes (0=success, 1=error)

### File Structure

```
volcano/
├── main.go           # Entry point, CLI parsing
├── cmd/
│   ├── generate.go   # Generate command logic
│   └── serve.go      # Serve command logic
├── go.mod
└── go.sum
```

### Implementation Notes

```go
// Flag structure
type Config struct {
    InputDir   string
    OutputDir  string
    ServeMode  bool
    Port       int
    Title      string
}
```

---

## Story 2: CI Pipeline Setup (Lint, Test, Coverage, E2E)

**Goal**: Set up GitHub Actions CI to enforce code quality, run tests with coverage, and execute end-to-end tests on every push and PR.

### Acceptance Criteria

- [ ] Create `.github/workflows/ci.yml` workflow file
- [ ] Trigger on:
  - Push to `main` branch
  - All pull requests
- [ ] **Linting job**:
  - Use `golangci-lint` with strict configuration
  - Create `.golangci.yml` config file
  - Enable linters: `gofmt`, `govet`, `errcheck`, `staticcheck`, `unused`, `gosimple`, `ineffassign`, `misspell`
  - Fail CI if any lint errors
- [ ] **Test job**:
  - Run `go test ./...` with race detector
  - Generate coverage report (`-coverprofile=coverage.out`)
  - Upload coverage to workflow artifacts
  - Set minimum coverage threshold: 70%
  - Fail CI if coverage drops below threshold
- [ ] **E2E job**:
  - Build the binary
  - Run volcano against `example/` folder
  - Verify output files exist
  - Verify generated HTML contains expected content
  - Test serve mode starts and responds to HTTP requests
- [ ] **Build matrix**:
  - Test on `ubuntu-latest` and `macos-latest`
  - Test with Go 1.21 and 1.22
- [ ] Add status badges to README

### Workflow Structure

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - uses: golangci/golangci-lint-action@v4
        with:
          version: latest

  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest]
        go: ['1.21', '1.22']
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      - name: Run tests with coverage
        run: go test -race -coverprofile=coverage.out -covermode=atomic ./...
      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$COVERAGE < 70" | bc -l) )); then
            echo "Coverage $COVERAGE% is below 70% threshold"
            exit 1
          fi
      - uses: actions/upload-artifact@v4
        with:
          name: coverage-${{ matrix.os }}-${{ matrix.go }}
          path: coverage.out

  e2e:
    runs-on: ubuntu-latest
    needs: [lint, test]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Build binary
        run: go build -o volcano .
      - name: Generate site
        run: ./volcano ./example -o ./test-output --title="E2E Test"
      - name: Verify output structure
        run: |
          test -f ./test-output/index.html
          test -f ./test-output/404.html
          test -d ./test-output/guides
      - name: Verify HTML content
        run: |
          grep -q "<title>" ./test-output/index.html
          grep -q "nav" ./test-output/index.html
      - name: Test serve mode
        run: |
          ./volcano -s -p 8888 ./test-output &
          sleep 2
          curl -f http://localhost:8888/ || exit 1
          curl -f http://localhost:8888/guides/ || exit 1
          kill %1
```

### Golangci-lint Configuration

```yaml
# .golangci.yml
run:
  timeout: 5m

linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - misspell
    - gocyclo
    - gocritic
    - revive

linters-settings:
  gocyclo:
    min-complexity: 15
  gocritic:
    enabled-tags:
      - diagnostic
      - style
      - performance

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

### File Location

```
volcano/
├── .github/
│   └── workflows/
│       └── ci.yml
├── .golangci.yml
```

### Makefile Targets (Optional)

```makefile
.PHONY: lint test coverage e2e

lint:
	golangci-lint run

test:
	go test -race ./...

coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

e2e: build
	./volcano ./example -o ./test-output
	./scripts/e2e-verify.sh
```

---

## Story 3: Markdown File Discovery & Tree Building

**Goal**: Scan input directory and build an in-memory tree structure of markdown files.

### Acceptance Criteria

- [ ] Recursively walk input directory
- [ ] Identify all `.md` files (case insensitive: `.md`, `.MD`, `.Md`)
- [ ] Skip hidden files/folders (starting with `.`)
- [ ] Skip empty folders (folders with no `.md` files at any depth)
- [ ] Build hierarchical tree structure representing folder/file relationships
- [ ] Sort entries alphabetically (folders first, then files)
- [ ] Generate clean labels from filenames:
  - Remove `.md` extension
  - Replace `-` and `_` with spaces
  - Title case each word
  - Examples:
    - `getting-started.md` → `"Getting Started"`
    - `api_reference.md` → `"Api Reference"`
    - `FAQ.md` → `"FAQ"`
- [ ] Identify `index.md` files as folder landing pages

### Data Structures

```go
// Represents a node in the content tree
type TreeNode struct {
    Name       string      // Clean display label
    Path       string      // Relative path from input root
    SourcePath string      // Full path to source .md file
    IsFolder   bool
    HasIndex   bool        // True if folder contains index.md
    IndexPath  string      // Path to index.md if exists
    Children   []*TreeNode // Sorted alphabetically
    Parent     *TreeNode
}

// Represents the full site structure
type SiteTree struct {
    Root     *TreeNode
    AllPages []*TreeNode // Flat list for easy iteration
}
```

### File Location

```
volcano/
├── internal/
│   └── tree/
│       ├── tree.go      # TreeNode, SiteTree types
│       ├── scanner.go   # Directory walking logic
│       └── labels.go    # Filename to label conversion
```

---

## Story 4: Markdown to HTML Conversion

**Goal**: Parse markdown files and convert to HTML content.

### Acceptance Criteria

- [ ] Use `github.com/gomarkdown/markdown` or `github.com/yuin/goldmark` for parsing
- [ ] Support standard markdown features:
  - Headings (h1-h6)
  - Paragraphs
  - Bold, italic, strikethrough
  - Links (internal and external)
  - Images (with alt text)
  - Code blocks with syntax highlighting
  - Inline code
  - Ordered and unordered lists
  - Blockquotes
  - Horizontal rules
  - Tables
- [ ] Enable syntax highlighting for code blocks using `github.com/alecthomas/chroma`
- [ ] Generate safe HTML (escape dangerous content)
- [ ] Extract first H1 as page title (fallback to clean filename)

### Data Structures

```go
type Page struct {
    Title      string // First H1 or clean filename
    Content    string // Rendered HTML
    SourcePath string // Path to original .md file
    OutputPath string // Path for output .html file
    URLPath    string // URL path for navigation links
}
```

### File Location

```
volcano/
├── internal/
│   └── markdown/
│       ├── parser.go      # Markdown parsing setup
│       ├── renderer.go    # HTML rendering with syntax highlighting
│       └── extractor.go   # Title extraction from content
```

---

## Story 5: HTML Template System

**Goal**: Create embedded HTML templates for the generated site.

### Acceptance Criteria

- [ ] Use Go's `html/template` package
- [ ] Embed templates using `//go:embed` directive
- [ ] Create base layout template with:
  - HTML5 doctype
  - Responsive viewport meta tag
  - Site title in `<title>` tag
  - Embedded Tailwind CSS (typography plugin styles)
  - Light/dark mode CSS variables
  - Navigation sidebar (desktop: left side, mobile: drawer)
  - Main content area (right side)
  - Dark mode toggle button
  - Theme detection script (browser preference)
- [ ] Template receives:
  - Site title
  - Navigation tree (for sidebar)
  - Current page content
  - Current page path (for active state highlighting)

### Template Structure

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.PageTitle}} - {{.SiteTitle}}</title>
    <style>/* Embedded Tailwind CSS */</style>
</head>
<body class="bg-white dark:bg-black text-black dark:text-white">
    <!-- Mobile menu button -->
    <!-- Sidebar navigation (tree) -->
    <!-- Main content area -->
    <!-- Dark mode toggle -->
    <!-- Theme + mobile nav JavaScript -->
</body>
</html>
```

### File Location

```
volcano/
├── internal/
│   └── templates/
│       ├── embed.go       # //go:embed directives
│       ├── layout.html    # Base page template
│       ├── nav.html       # Navigation tree partial
│       └── renderer.go    # Template execution logic
```

---

## Story 6: Navigation Tree Component

**Goal**: Implement the collapsible tree navigation in the sidebar.

### Acceptance Criteria

- [ ] Render tree structure as nested `<ul>/<li>` elements
- [ ] Folders are collapsible with expand/collapse toggle
- [ ] Default state: expand path to current page, collapse others
- [ ] Folder click behavior:
  - If folder has `index.md`: clicking folder name navigates to index
  - Expand/collapse toggle is separate clickable element (chevron icon)
- [ ] Files link directly to their page
- [ ] Current page highlighted with distinct background color
- [ ] Smooth expand/collapse animation (CSS transitions)
- [ ] Accessible: proper ARIA attributes for tree navigation

### HTML Structure

```html
<nav class="tree-nav" aria-label="Site navigation">
  <ul role="tree">
    <li role="treeitem" aria-expanded="true">
      <div class="folder">
        <button class="toggle" aria-label="Toggle section">
          <svg><!-- Chevron icon --></svg>
        </button>
        <a href="/guides/">Guides</a>
      </div>
      <ul role="group">
        <li role="treeitem">
          <a href="/guides/getting-started/" class="active">Getting Started</a>
        </li>
      </ul>
    </li>
  </ul>
</nav>
```

### JavaScript Behavior

```javascript
// Toggle folder expansion
// Remember expanded state in localStorage
// Collapse all others when navigating (optional)
```

---

## Story 7: Embedded CSS Styling

**Goal**: Create embedded CSS with Tailwind-style typography and theming.

### Acceptance Criteria

- [ ] Generate minimal CSS bundle (no build step, embedded in Go binary)
- [ ] Include Tailwind Typography styles for markdown content:
  - Proper heading sizes and margins
  - Paragraph spacing
  - List styling (bullets, numbers)
  - Code block styling (background, padding, font)
  - Inline code styling
  - Blockquote styling
  - Table styling
  - Link colors and hover states
- [ ] Color scheme: black and white shades only
  - Light mode: white background (#ffffff), black text (#000000)
  - Dark mode: black background (#000000), white text (#ffffff)
  - Gray variations for subtle elements (#f5f5f5, #e5e5e5, #333, #666)
- [ ] CSS custom properties for theme colors
- [ ] Responsive breakpoints:
  - Mobile: < 768px (single column, drawer nav)
  - Desktop: >= 768px (sidebar + content)
- [ ] Sidebar width: 280px on desktop
- [ ] Max content width: 800px

### CSS Custom Properties

```css
:root {
  --bg-primary: #ffffff;
  --bg-secondary: #f5f5f5;
  --text-primary: #000000;
  --text-secondary: #333333;
  --border-color: #e5e5e5;
}

[data-theme="dark"] {
  --bg-primary: #000000;
  --bg-secondary: #111111;
  --text-primary: #ffffff;
  --text-secondary: #cccccc;
  --border-color: #333333;
}
```

### File Location

```
volcano/
├── internal/
│   └── styles/
│       ├── embed.go      # //go:embed for CSS
│       ├── base.css      # Reset, variables, layout
│       ├── typography.css # Prose/content styling
│       └── nav.css       # Navigation tree styling
```

---

## Story 8: Mobile Navigation Drawer

**Goal**: Implement responsive drawer navigation for mobile devices.

### Acceptance Criteria

- [ ] On mobile (< 768px):
  - Navigation hidden by default
  - Hamburger menu button fixed in top-left corner
  - Clicking hamburger opens drawer from left
  - Drawer overlays content with semi-transparent backdrop
  - Clicking backdrop closes drawer
  - Close button inside drawer
  - Drawer contains full tree navigation
- [ ] On desktop (>= 768px):
  - Hamburger button hidden
  - Sidebar always visible
  - No overlay/backdrop
- [ ] Smooth slide-in/slide-out animation
- [ ] Body scroll locked when drawer is open
- [ ] Accessible: focus trap when drawer is open, ESC key closes

### HTML Structure

```html
<!-- Mobile hamburger -->
<button class="mobile-menu-btn md:hidden" aria-label="Open menu">
  <svg><!-- Hamburger icon --></svg>
</button>

<!-- Backdrop -->
<div class="drawer-backdrop hidden" aria-hidden="true"></div>

<!-- Drawer -->
<aside class="drawer" aria-label="Navigation">
  <button class="close-btn" aria-label="Close menu">
    <svg><!-- X icon --></svg>
  </button>
  <!-- Tree nav here -->
</aside>
```

---

## Story 9: Dark Mode Toggle

**Goal**: Implement light/dark mode switching with browser preference detection.

### Acceptance Criteria

- [ ] Toggle button visible on all pages (top-right corner)
- [ ] Icon changes based on current mode (sun/moon)
- [ ] On first visit: detect `prefers-color-scheme` media query
- [ ] Store user preference in `localStorage`
- [ ] Preference persists across pages and sessions
- [ ] Apply theme before page paint (prevent flash of wrong theme)
- [ ] Toggle smoothly transitions colors (CSS transition on custom properties)

### JavaScript Implementation

```javascript
// Run before body renders to prevent flash
(function() {
  const stored = localStorage.getItem('theme');
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  const theme = stored || (prefersDark ? 'dark' : 'light');
  document.documentElement.setAttribute('data-theme', theme);
})();

// Toggle function
function toggleTheme() {
  const current = document.documentElement.getAttribute('data-theme');
  const next = current === 'dark' ? 'light' : 'dark';
  document.documentElement.setAttribute('data-theme', next);
  localStorage.setItem('theme', next);
}
```

---

## Story 10: Static Site Generation Engine

**Goal**: Orchestrate the full build process from input folder to output folder.

### Acceptance Criteria

- [ ] Create output directory if it doesn't exist
- [ ] Clean output directory before generation (optional `--clean` flag)
- [ ] For each markdown file:
  1. Read file contents
  2. Convert to HTML
  3. Apply template with navigation tree
  4. Write to output path preserving folder structure
- [ ] URL structure:
  - `input/guides/intro.md` → `output/guides/intro/index.html` (clean URLs)
  - `input/index.md` → `output/index.html` (root index)
- [ ] Generate 404.html page
- [ ] Print summary on completion:
  - Number of pages generated
  - Output directory path
  - Any warnings (empty folders skipped, etc.)
- [ ] Return appropriate exit code

### Build Pipeline

```
1. Parse CLI arguments
2. Validate input directory
3. Scan for markdown files
4. Build tree structure
5. For each page:
   a. Read markdown
   b. Convert to HTML
   c. Execute template
   d. Write output file
6. Print summary
```

### File Location

```
volcano/
├── internal/
│   └── generator/
│       ├── generator.go  # Main orchestration
│       ├── writer.go     # File writing utilities
│       └── urls.go       # URL/path generation
```

---

## Story 11: Static File Server

**Goal**: Implement simple HTTP server for previewing generated sites.

### Acceptance Criteria

- [ ] Serve files from specified directory
- [ ] Default port: 1776 (configurable via `-p` flag)
- [ ] Clean URL support: `/guides/intro/` serves `/guides/intro/index.html`
- [ ] Serve 404.html for missing pages
- [ ] Print server URL on startup
- [ ] Graceful shutdown on SIGINT/SIGTERM
- [ ] Log requests to stdout (method, path, status, duration)
- [ ] MIME type detection for CSS, JS, images
- [ ] Cache-Control headers for development (no-cache)

### Implementation

```go
func Serve(dir string, port int) error {
    fs := http.FileServer(http.Dir(dir))

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Try clean URL resolution
        // Log request
        // Handle 404
        fs.ServeHTTP(w, r)
    })

    addr := fmt.Sprintf(":%d", port)
    fmt.Printf("Serving at http://localhost%s\n", addr)
    return http.ListenAndServe(addr, handler)
}
```

### File Location

```
volcano/
├── internal/
│   └── server/
│       └── server.go
```

---

## Story 12: Error Handling & User Feedback

**Goal**: Provide clear error messages and progress feedback.

### Acceptance Criteria

- [ ] Descriptive error messages for common issues:
  - Input directory doesn't exist
  - Input directory is empty (no markdown files)
  - Output directory not writable
  - Port already in use
  - Invalid markdown file (parse error)
- [ ] Colored terminal output (if TTY detected):
  - Errors in red
  - Warnings in yellow
  - Success in green
- [ ] Progress output during generation:
  ```
  Scanning input directory...
  Found 12 markdown files in 4 folders
  Generating pages...
    ✓ index.md
    ✓ guides/getting-started.md
    ✓ guides/advanced/configuration.md
    ...

  Generated 12 pages in output/
  ```
- [ ] Quiet mode flag (`-q`) to suppress non-error output
- [ ] Verbose mode flag (`-v`) for debug output

### File Location

```
volcano/
├── internal/
│   └── output/
│       ├── logger.go    # Logging utilities
│       └── colors.go    # Terminal color support
```

---

## Story 13: Testing & Example Site

**Goal**: Set up testing infrastructure and create example content.

### Acceptance Criteria

- [ ] Unit tests for:
  - Label generation from filenames
  - Tree building from directory structure
  - Markdown parsing
  - URL generation
- [ ] Integration test:
  - Generate site from example folder
  - Verify output file structure
  - Verify HTML contains expected elements
- [ ] Create `example/` folder with test content:
  ```
  example/
  ├── index.md              # Welcome page
  ├── getting-started.md    # Quick start guide
  ├── guides/
  │   ├── index.md          # Guides overview
  │   ├── installation.md
  │   └── configuration.md
  ├── api/
  │   ├── index.md
  │   ├── endpoints.md
  │   └── authentication.md
  └── faq.md
  ```
- [ ] Example content demonstrates:
  - All markdown features (headers, code, lists, tables)
  - Nested folder structure
  - Index files for folders
- [ ] Test commands documented in README

### File Location

```
volcano/
├── example/           # Test input folder
├── internal/
│   └── tree/
│       └── tree_test.go
│   └── markdown/
│       └── parser_test.go
```

---

## Future Enhancements

These items are out of scope for initial implementation but documented for future development:

### Frontmatter Support
- YAML frontmatter parsing
- Fields: title, description, order, draft, date
- Override auto-generated titles
- Custom ordering within folders
- Draft pages excluded from production builds

### Search Functionality
- Client-side search using generated JSON index
- Search index includes: title, content preview, URL
- Search UI in navigation header
- Keyboard shortcut (Cmd/Ctrl + K)

### Asset Handling
- Copy non-markdown files to output
- Image optimization
- CSS/JS minification

### Live Reload
- WebSocket-based live reload in serve mode
- File watcher for markdown changes
- Auto-refresh browser on change

### Additional Features
- RSS/Atom feed generation
- Sitemap.xml generation
- Custom CSS injection via flag
- Multiple theme options

---

## Dependencies

```go
// go.mod
module volcano

go 1.21

require (
    github.com/yuin/goldmark v1.6.0
    github.com/alecthomas/chroma/v2 v2.12.0
)
```

Minimal dependencies - only markdown parsing and syntax highlighting.

---

## Build & Run

```bash
# Build
go build -o volcano .

# Generate site
./volcano ./example -o ./output --title="My Docs"

# Serve site
./volcano -s -p 1776 ./output

# Combined: generate and serve
./volcano ./example -o ./output && ./volcano -s ./output
```

# Volcano Reorganization

Two implementation stories for restructuring the CLI and CSS.

---

## Story 1: CLI Restructure - Separate Build and Serve Commands

### Goal
Change CLI from flag-based mode selection (`-s`) to explicit subcommands like the existing `css` command.

### New CLI Structure
```bash
volcano build [flags] <folder>     # Generate static site
volcano serve [flags] <folder>     # Start dev server
volcano server [flags] <folder>    # Alias for serve
volcano css [-o file]              # Output vanilla CSS (unchanged)
volcano <folder>                   # Shorthand for build (backward compat)
```

### Files to Modify

#### 1. `main.go`
- Add switch statement for subcommand dispatch (build, serve, server, css)
- Fall through to build for backward compatibility (`volcano ./docs` still works)
- Remove flag definitions (move to cmd/build.go and cmd/serve.go)
- Remove `runWithConfig()` function (logic moves to subcommands)
- Keep `printUsage()` but update for new command structure
- Keep `reorderArgs()` helper (used by both commands)

**New dispatch logic:**
```go
if len(args) > 0 {
    switch args[0] {
    case "css":
        return cmd.CSS(args[1:], stdout)
    case "build":
        return cmd.Build(args[1:], stdout, stderr)
    case "serve", "server":
        return cmd.Serve(args[1:], stdout, stderr)
    }
}
// Fall through: treat as shorthand for build
return cmd.Build(args, stdout, stderr)
```

#### 2. `cmd/build.go` (rename from generate.go)
Create new `Build(args []string, stdout, stderr io.Writer) error` function that:
- Creates its own FlagSet named "build"
- Defines all build-specific flags
- Parses args with `reorderArgs()` for flexible arg order
- Validates input directory
- Calls existing `Generate(cfg, stdout)` internally

**Build flags:**
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | ./output | Output directory |
| `--title` | | My Site | Site title |
| `--theme` | | docs | Theme (docs, blog, vanilla) |
| `--css` | | | Custom CSS file path |
| `--url` | | | Base URL for SEO |
| `--author` | | | Site author |
| `--og-image` | | | Default Open Graph image |
| `--favicon` | | | Favicon file path |
| `--last-modified` | | false | Show last modified dates |
| `--top-nav` | | false | Display root files in top nav |
| `--page-nav` | | false | Show prev/next navigation |
| `--quiet` | `-q` | false | Suppress output |
| `--verbose` | | false | Debug output |

#### 3. `cmd/serve.go`
Update existing file to export `Serve(args []string, stdout, stderr io.Writer) error` that:
- Creates its own FlagSet named "serve"
- Defines serve-specific flags
- Parses args with `reorderArgs()`
- Calls existing serve logic internally

**Serve flags:**
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--port` | `-p` | 1776 | Server port |
| `--title` | | My Site | Site title (dynamic mode) |
| `--theme` | | docs | Theme (dynamic mode) |
| `--css` | | | Custom CSS (dynamic mode) |
| `--quiet` | `-q` | false | Suppress output |
| `--verbose` | | false | Debug output |

#### 4. `cmd/config.go`
- Keep `Config` struct unchanged
- Keep `DefaultConfig()` unchanged
- Both Build and Serve use this internally

#### 5. Update help text
```
Volcano - Static site generator

Usage:
  volcano build [flags] <input>    Generate static site
  volcano serve [flags] <input>    Start development server
  volcano server [flags] <input>   Alias for serve
  volcano css [-o file]            Output vanilla CSS
  volcano <input>                  Shorthand for build

Run 'volcano <command> --help' for command-specific help.
```

### Verification
```bash
./volcano build ./docs -o ./output --title="Test"
./volcano ./docs -o ./output  # Shorthand still works
./volcano build --help
./volcano serve ./docs -p 8080
./volcano server ./docs  # Alias works
./volcano serve --help
./volcano css
go test -race ./...
```

---

## Story 2: CSS Split - Layout vs Styling

### Goal
Split CSS into a shared layout file + per-theme styling files. All themes share the same layout structure but have different colors, fonts, and decorations.

### New File Structure
```
internal/styles/themes/
├── layout.css          # Shared structural CSS
├── docs.css            # Docs theme styling only
├── blog.css            # Blog theme styling only
└── vanilla.css         # Minimal/empty styling
```

### What Goes in `layout.css` (Shared Base)
- CSS reset
- Flexbox/grid layout rules
- Position properties (fixed, absolute, relative)
- Layout dimension variables (--sidebar-width, --content-max-width, --toc-width)
- Z-index stacking
- Display properties (flex, block, none for mobile)
- Overflow handling
- Mobile responsiveness (@media queries for layout changes)
- Print stylesheet (hiding elements)
- Transform animations (translateX for drawer)
- All structural selectors with ONLY layout properties

### What Goes in Theme Files (Styling Only)
- CSS color variables (--bg-*, --text-*, --border-*, --accent-*)
- Font family declarations
- Font sizes, weights, line-heights, letter-spacing
- Background colors
- Border colors and styles
- Box shadows
- Border-radius
- Link colors and hover states
- Syntax highlighting colors (.chroma classes)
- Admonition styling (colors, icons)
- Decorative transitions

### Files to Modify

#### 1. Create `internal/styles/themes/layout.css`
Extract structural CSS from docs.css:
- All position/display/flex properties
- Width/height/margin/padding dimensions
- Grid/flexbox definitions
- Mobile breakpoints (structural changes only)

#### 2. Refactor `internal/styles/themes/docs.css`
- Remove layout properties (now in layout.css)
- Keep only styling: colors, fonts, borders, shadows
- Continue using same CSS variable naming convention

#### 3. Refactor `internal/styles/themes/blog.css`
- Remove layout properties
- Keep blog-specific styling (serif fonts, warmer colors, larger text)

#### 4. Refactor `internal/styles/themes/vanilla.css`
- Remove layout properties (now shared)
- Keep as minimal file with comments showing customization hints

#### 5. Update `internal/styles/embed.go`
```go
//go:embed themes/layout.css
var LayoutCSS string

//go:embed themes/docs.css
var DocsCSS string

//go:embed themes/blog.css
var BlogCSS string

//go:embed themes/vanilla.css
var VanillaCSS string

func GetCSS(theme string) string {
    var themeCSS string
    switch theme {
    case "blog":
        themeCSS = BlogCSS
    case "vanilla":
        themeCSS = VanillaCSS
    default:
        themeCSS = DocsCSS
    }
    return LayoutCSS + "\n" + themeCSS  // Combine layout + theme
}
```

#### 6. Update `cmd/css.go`
No changes needed - `GetCSS("vanilla")` already includes layout.

### Verification
```bash
./volcano build ./docs -o ./out1 --theme=docs
./volcano build ./docs -o ./out2 --theme=blog
./volcano build ./docs -o ./out3 --theme=vanilla
# Verify each output looks correct in browser
./volcano css | head -50  # Should show layout first
go test -race ./...
```

---

## Story 3: OG Image as Local File Path

### Goal
Change `--og-image` from accepting a URL to accepting a local file path. The image is copied to the output directory and the og:image meta tag references the built asset.

### Current Behavior
```bash
volcano build ./docs --og-image="https://example.com/image.png"
# og:image meta tag uses the URL directly
```

### New Behavior
```bash
volcano build ./docs --og-image="./assets/social.png" --url="https://mysite.com"
# Image copied to output/og-image.png
# og:image meta tag: <meta property="og:image" content="https://mysite.com/og-image.png">
```

### Implementation Pattern
Follow the existing favicon pattern in `internal/assets/favicon.go`:
1. Validate source file exists
2. Copy to output directory with standardized name
3. Generate meta tag with full URL

### Files to Modify

#### 1. Create `internal/assets/ogimage.go`
```go
package assets

// OGImageConfig holds configuration for OG image handling
type OGImageConfig struct {
    ImagePath string // Local path to OG image file
    BaseURL   string // Site base URL for absolute URL generation
}

// ProcessOGImage copies the OG image to output and returns the URL
func ProcessOGImage(config OGImageConfig, outputDir string) (string, error) {
    if config.ImagePath == "" {
        return "", nil
    }

    // Validate file exists
    if _, err := os.Stat(config.ImagePath); os.IsNotExist(err) {
        return "", fmt.Errorf("og-image file not found: %s", config.ImagePath)
    }

    // Get extension and validate format
    ext := strings.ToLower(filepath.Ext(config.ImagePath))
    if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" && ext != ".webp" {
        return "", fmt.Errorf("unsupported og-image format: %s (use png, jpg, gif, or webp)", ext)
    }

    // Copy to output as og-image.{ext}
    destFilename := "og-image" + ext
    destPath := filepath.Join(outputDir, destFilename)
    if err := copyFile(config.ImagePath, destPath); err != nil {
        return "", err
    }

    // Build absolute URL
    ogImageURL := "/" + destFilename
    if config.BaseURL != "" {
        ogImageURL = strings.TrimSuffix(config.BaseURL, "/") + "/" + destFilename
    }

    return ogImageURL, nil
}
```

#### 2. Update `internal/seo/meta.go`
Modify to accept pre-computed og:image URL instead of raw path:
- Change `OGImage string` field to accept the processed URL
- The URL is computed by the generator before calling SEO functions

#### 3. Update `internal/generator/generator.go`
In the generation flow:
```go
// Process OG image (similar to favicon processing)
var ogImageURL string
if g.config.OGImage != "" {
    url, err := assets.ProcessOGImage(assets.OGImageConfig{
        ImagePath: g.config.OGImage,
        BaseURL:   g.config.SiteURL,
    }, g.config.OutputDir)
    if err != nil {
        return fmt.Errorf("failed to process og-image: %w", err)
    }
    ogImageURL = url
}
// Pass ogImageURL to SEO meta generation
```

#### 4. Update `cmd/config.go`
- Rename `OGImage` field comment to clarify it's a path, not URL
- Validation: check file exists if path provided

#### 5. Update `internal/server/dynamic.go`
For dynamic serve mode, also process og-image similarly.

### Supported Formats
- `.png` (recommended for OG images)
- `.jpg` / `.jpeg`
- `.gif`
- `.webp`

### Output
```
output/
├── index.html
├── og-image.png      # Copied from --og-image path
├── favicon.ico       # Existing favicon handling
└── ...
```

### Meta Tag Output
```html
<meta property="og:image" content="https://mysite.com/og-image.png">
```

### Edge Cases
- No `--url` provided: Use relative path `/og-image.png`
- No `--og-image` provided: No og:image meta tag (existing behavior)
- Invalid file format: Error with supported formats list
- File not found: Clear error message

### Verification
```bash
# With base URL
./volcano build ./docs --og-image="./social.png" --url="https://example.com"
# Check output/og-image.png exists
# Check index.html contains: <meta property="og:image" content="https://example.com/og-image.png">

# Without base URL
./volcano build ./docs --og-image="./social.png"
# Check og:image uses relative path

# Error cases
./volcano build ./docs --og-image="./missing.png"  # Should error
./volcano build ./docs --og-image="./file.txt"    # Should error on format

go test -race ./...
```

---

---

## Story 4: Fix Favicon Support in Serve Mode

### Goal
Enable favicon support in serve mode (`volcano -s`) to match build mode functionality. Currently, the `--favicon` flag is accepted but ignored during serve mode, resulting in no favicon links in the HTML and no file being served.

### Current Behavior
```bash
./volcano -s -p 4242 --favicon="docs/logo.png" docs
# ✗ No <link rel="icon"> in HTML
# ✗ favicon file not accessible
# ✗ Browser shows no favicon
```

### Expected Behavior
```bash
./volcano -s -p 4242 --favicon="docs/logo.png" docs
# ✓ <link rel="icon" type="image/png" href="/logo.png"> in HTML
# ✓ Favicon file accessible at http://localhost:4242/logo.png
# ✓ Browser displays favicon
```

### Root Causes

**Build mode works correctly:**
- `internal/generator/generator.go:120-129` - Processes favicon, copies file to output, stores HTML tags
- `internal/generator/generator.go:361` - Passes `FaviconLinks` to PageData
- Result: Static files include favicon, HTML has correct tags

**Serve mode is broken:**
1. **Missing config field** - `DynamicConfig` (internal/server/dynamic.go:29-39) lacks `FaviconPath`
2. **Config not passed** - `cmd/serve.go:20-30` doesn't pass `cfg.FaviconPath` to DynamicConfig
3. **No processing logic** - `DynamicServer` never calls `assets.ProcessFavicon()`
4. **No HTML tags** - PageData in dynamic.go never sets `FaviconLinks` field (lines 340-353, 587-593, 626-632, 737-746)
5. **File not served** - Static file serving works for existing files in source dir, but favicon isn't accessible at root level

### Implementation Strategy

Serve mode needs a **hybrid approach** because it doesn't write to an output directory:
1. Store favicon in memory after processing (no file copy needed)
2. Generate HTML link tags (reuse existing `assets.RenderFaviconLinks`)
3. Serve the favicon file from memory via special HTTP handler
4. Pass favicon HTML to all PageData instances

### Files to Modify

#### 1. Update `internal/server/dynamic.go`

**Add FaviconPath to DynamicConfig:**
```go
type DynamicConfig struct {
	SourceDir   string
	Title       string
	Port        int
	Quiet       bool
	Verbose     bool
	TopNav      bool
	ShowPageNav bool
	Theme       string
	CSSPath     string
	FaviconPath string // Path to favicon file
}
```

**Add favicon state to DynamicServer:**
```go
type DynamicServer struct {
	config       DynamicConfig
	renderer     *templates.Renderer
	transformer  *markdown.ContentTransformer
	writer       io.Writer
	server       *http.Server
	fs           FileSystem
	scanner      TreeScanner
	cssLoader    styles.CSSLoader
	faviconLinks template.HTML // Favicon HTML tags
	faviconData  []byte        // Favicon file content (in memory)
	faviconMime  string        // Favicon MIME type
	faviconName  string        // Favicon filename (e.g., "logo.png")
}
```

**Process favicon in NewDynamicServer:**
```go
func NewDynamicServer(config DynamicConfig, writer io.Writer) (*DynamicServer, error) {
	// ... existing CSS loading code ...

	srv := &DynamicServer{
		config:      config,
		renderer:    renderer,
		transformer: markdown.NewContentTransformer(""),
		writer:      writer,
		fs:          osFileSystem{},
		scanner:     defaultScanner{},
		cssLoader:   cssLoader,
	}

	// Process favicon if configured
	if config.FaviconPath != "" {
		if err := srv.loadFavicon(); err != nil {
			// Log warning but don't fail server startup
			_, _ = fmt.Fprintf(writer, "Warning: Failed to load favicon: %v\n", err)
		}
	}

	return srv, nil
}

// loadFavicon loads the favicon file into memory and generates HTML tags
func (s *DynamicServer) loadFavicon() error {
	// Validate file exists
	if _, err := os.Stat(s.config.FaviconPath); os.IsNotExist(err) {
		return fmt.Errorf("favicon file not found: %s", s.config.FaviconPath)
	}

	// Read file into memory
	data, err := os.ReadFile(s.config.FaviconPath)
	if err != nil {
		return fmt.Errorf("failed to read favicon: %w", err)
	}

	// Get filename and MIME type
	filename := filepath.Base(s.config.FaviconPath)
	mimeType := getMimeType(filename)
	if mimeType == "" {
		return fmt.Errorf("unsupported favicon format: %s", filename)
	}

	// Store in memory
	s.faviconData = data
	s.faviconMime = mimeType
	s.faviconName = filename

	// Generate HTML link tags
	links := []assets.FaviconLink{{
		Rel:  "icon",
		Type: mimeType,
		Href: "/" + filename,
	}}
	s.faviconLinks = assets.RenderFaviconLinks(links)

	return nil
}

// getMimeType returns MIME type for favicon file (helper function)
func getMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".ico":
		return "image/x-icon"
	case ".png":
		return "image/png"
	case ".svg":
		return "image/svg+xml"
	case ".gif":
		return "image/gif"
	default:
		return ""
	}
}
```

**Update handleRequest to serve favicon:**
```go
func (s *DynamicServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

	// Set cache control headers
	rec.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	rec.Header().Set("Pragma", "no-cache")
	rec.Header().Set("Expires", "0")

	urlPath := r.URL.Path

	// Serve favicon from memory if requested
	if s.serveFavicon(rec, urlPath) {
		s.logRequest(r.Method, urlPath, rec.statusCode, time.Since(start))
		return
	}

	// ... rest of existing handlers (static files, pages, auto-index, 404) ...
}

// serveFavicon serves the favicon from memory
func (s *DynamicServer) serveFavicon(w http.ResponseWriter, urlPath string) bool {
	if s.faviconData == nil || s.faviconName == "" {
		return false
	}

	// Check if the request is for the favicon
	requestedFile := strings.TrimPrefix(urlPath, "/")
	if requestedFile != s.faviconName {
		return false
	}

	// Serve from memory
	w.Header().Set("Content-Type", s.faviconMime)
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	_, _ = w.Write(s.faviconData)
	return true
}
```

**Update all PageData instances to include FaviconLinks:**

In `renderPage()` (line ~340-353):
```go
data := templates.PageData{
	SiteTitle:    s.config.Title,
	PageTitle:    page.Title,
	Content:      template.HTML(htmlContent),
	Navigation:   nav,
	CurrentPath:  nodeURLPath,
	Breadcrumbs:  breadcrumbsHTML,
	PageNav:      pageNavHTML,
	TOC:          tocHTML,
	FaviconLinks: s.faviconLinks, // Add this line
	ReadingTime:  readingTime,
	HasTOC:       hasTOC,
	ShowSearch:   true,
	TopNavItems:  topNavItems,
}
```

In `serveBrokenLinksError()` (line ~587-593):
```go
data := templates.PageData{
	SiteTitle:    s.config.Title,
	PageTitle:    "Build Error",
	Content:      template.HTML(sb.String()),
	Navigation:   nav,
	CurrentPath:  "",
	FaviconLinks: s.faviconLinks, // Add this line
}
```

In `serve404()` (line ~626-632):
```go
data := templates.PageData{
	SiteTitle:    s.config.Title,
	PageTitle:    "Page Not Found",
	Content:      template.HTML(content),
	Navigation:   nav,
	CurrentPath:  "",
	FaviconLinks: s.faviconLinks, // Add this line
}
```

In `renderAutoIndex()` (line ~737-746):
```go
data := templates.PageData{
	SiteTitle:    s.config.Title,
	PageTitle:    node.Name,
	Content:      htmlContent,
	Navigation:   nav,
	CurrentPath:  urlPath,
	Breadcrumbs:  breadcrumbsHTML,
	FaviconLinks: s.faviconLinks, // Add this line
	ShowSearch:   true,
	TopNavItems:  topNavItems,
}
```

#### 2. Update `cmd/serve.go`

Pass `FaviconPath` to DynamicConfig:
```go
func Serve(cfg *Config, w io.Writer) error {
	if isSourceDirectory(cfg.InputDir) {
		dynamicCfg := server.DynamicConfig{
			SourceDir:   cfg.InputDir,
			Title:       cfg.Title,
			Port:        cfg.Port,
			Quiet:       cfg.Quiet,
			Verbose:     cfg.Verbose,
			TopNav:      cfg.TopNav,
			ShowPageNav: cfg.ShowPageNav,
			Theme:       cfg.Theme,
			CSSPath:     cfg.CSSPath,
			FaviconPath: cfg.FaviconPath, // Add this line
		}
		// ... rest of serve logic ...
	}
	// ... static file serving ...
}
```

### Design Decisions

**Why serve from memory instead of copying to output?**
- Serve mode is for live development - no output directory exists
- In-memory serving is fast and matches the dynamic nature of serve mode
- Avoids filesystem pollution and cleanup complexity

**Why not reuse the full `assets.ProcessFavicon()`?**
- That function copies files to disk, which serve mode doesn't need
- We only need the MIME type detection and HTML generation parts
- The `getMimeType()` helper provides what we need without file I/O

**Cache headers strategy:**
- Dynamic content: `no-cache` (for live reload during development)
- Favicon: `max-age=3600` (static asset, safe to cache for 1 hour)

### Edge Cases

1. **Invalid file format** - Server logs warning, continues without favicon
2. **File not found** - Server logs warning, continues without favicon
3. **File deleted after server starts** - Served from memory, still works
4. **File modified after server starts** - Restart required (acceptable for development)
5. **No `--favicon` flag** - No favicon served, no HTML tags (existing behavior)

### Verification

```bash
# Test with PNG favicon
./volcano -s -p 4242 --favicon="docs/logo.png" docs
# Open http://localhost:4242/
# Check HTML source for <link rel="icon" type="image/png" href="/logo.png">
# Check http://localhost:4242/logo.png directly
# Verify browser shows favicon in tab

# Test with ICO favicon
./volcano -s -p 4242 --favicon="./favicon.ico" docs
# Check HTML and browser

# Test with SVG favicon
./volcano -s -p 4242 --favicon="./icon.svg" docs
# Check HTML and browser

# Test without favicon flag (should work as before)
./volcano -s -p 4242 docs
# Check no favicon links in HTML

# Test invalid file (should log warning but serve normally)
./volcano -s -p 4242 --favicon="./missing.png" docs
# Check server logs for warning
# Check site still works without favicon

# Test invalid format (should log warning but serve normally)
./volcano -s -p 4242 --favicon="./file.txt" docs
# Check server logs for warning

# Run tests
go test -race ./...

# Test that build mode still works
./volcano build ./docs --favicon="docs/logo.png"
# Check output/logo.png exists
# Check HTML includes favicon link
```

### Testing Strategy

**Unit Tests (internal/server/dynamic_test.go):**
- `TestDynamicServer_LoadFavicon_Success` - Valid PNG/ICO/SVG files
- `TestDynamicServer_LoadFavicon_FileNotFound` - Missing file handling
- `TestDynamicServer_LoadFavicon_InvalidFormat` - Unsupported format handling
- `TestDynamicServer_ServeFavicon` - HTTP serving from memory
- `TestDynamicServer_PageData_IncludesFaviconLinks` - HTML tags in output

**Integration Tests (integration_test.go):**
- `TestIntegrationServe_WithFavicon` - End-to-end serve with favicon
- Verify HTML includes favicon link tags
- Verify favicon file is accessible via HTTP
- Verify behavior without --favicon flag

---

## Implementation Order

1. **Story 1 first** - CLI changes are independent and foundational
2. **Story 2 second** - CSS split can be done after CLI is stable
3. **Story 3 third** - OG image can be done independently
4. **Story 4 fourth** - Favicon fix for serve mode

Each story should be committed separately.

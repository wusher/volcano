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

## Implementation Order

1. **Story 1 first** - CLI changes are independent and foundational
2. **Story 2 second** - CSS split can be done after CLI is stable

Each story should be committed separately.

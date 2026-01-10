# PLAN_3: Themes, CSS, and Navigation Enhancements

Stories 38-44 for the Volcano static site generator.

---

## Story 38: Theme System Architecture

### Description
Refactor CSS to support multiple themes: `docs` (current), `blog`, and `vanilla`. The current embedded CSS becomes the `docs` theme.

### Implementation Details

**1. Create theme CSS files structure:**
```
internal/styles/
├── themes/
│   ├── docs.css      # Current styles.css content
│   ├── blog.css      # New blog theme
│   └── vanilla.css   # Minimal structural CSS
├── embed.go          # Updated to embed all themes
└── themes.go         # Theme selection logic
```

**2. Add `--theme` flag:**
- Flag: `--theme <name>` (default: `docs`)
- Valid values: `docs`, `blog`, `vanilla`
- Add to `Config` struct in `cmd/config.go`
- Add flag in `main.go`
- Pass to generator and dynamic server

**3. Update `internal/styles/embed.go`:**
```go
//go:embed themes/docs.css
var DocsCSS string

//go:embed themes/blog.css
var BlogCSS string

//go:embed themes/vanilla.css
var VanillaCSS string

func GetCSS(theme string) string {
    switch theme {
    case "blog":
        return BlogCSS
    case "vanilla":
        return VanillaCSS
    default:
        return DocsCSS
    }
}
```

**4. Files to modify:**
- `cmd/config.go`: Add `Theme string` field
- `main.go`: Add `--theme` flag, update `isValueFlag()`
- `cmd/generate.go`: Pass theme to generator
- `cmd/serve.go`: Pass theme to dynamic server
- `internal/generator/generator.go`: Use `styles.GetCSS(theme)`
- `internal/server/dynamic.go`: Use `styles.GetCSS(theme)`

### Acceptance Criteria
- [ ] `volcano docs --theme=docs` uses current styling
- [ ] `volcano docs --theme=blog` uses serif/blog styling
- [ ] `volcano docs --theme=vanilla` uses minimal styling
- [ ] Default is `docs` when flag not specified
- [ ] Invalid theme name produces error message

---

## Story 39: Blog Theme

### Description
Create a clean, readable blog theme with serif typography and a less prominent sidebar.

### CSS Specifications

**Typography:**
```css
body {
  font-family: Georgia, "Times New Roman", Times, serif;
  font-size: 18px;
  line-height: 1.8;
}

.prose h1, .prose h2, .prose h3 {
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
}
```

**Sidebar changes:**
```css
:root {
  --sidebar-width: 220px;  /* Reduced from 280px */
}

.sidebar {
  border-right: none;  /* No visual separator */
  background-color: var(--bg-primary);  /* Same as content */
}

.tree-nav {
  font-size: 13px;  /* Smaller than docs theme */
}

.tree-nav a {
  padding: 0.25rem 0.375rem;  /* Tighter spacing */
}
```

**Content area:**
```css
.prose {
  max-width: 680px;  /* Narrower for readability */
}
```

### Files to Create
- `internal/styles/themes/blog.css` - Full theme CSS

### Acceptance Criteria
- [ ] Serif body font with sans-serif headings
- [ ] Narrower sidebar (220px vs 280px)
- [ ] No border between sidebar and content
- [ ] Smaller navigation text
- [ ] Narrower content column for optimal reading

---

## Story 40: Vanilla Theme

### Description
Create a minimal CSS theme with only structural/positioning styles. No colors, no decorative styling. This serves as a base for users who want to provide their own CSS.

### CSS Specifications

**Include only:**
- Box model reset
- Layout positioning (sidebar, content, TOC)
- Flexbox/grid for structure
- Button/toggle positioning
- Code block copy button positioning
- Mobile breakpoints for responsive layout
- Print stylesheet basics

**Exclude:**
- All color values (use `inherit` or browser defaults)
- Borders (except where structurally needed)
- Border-radius
- Box-shadow
- Background colors
- Font families (use system defaults)
- Font weights (except structural)
- Transitions/animations
- Hover effects

**Example structure:**
```css
/* Reset */
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

/* Layout */
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: 280px;
  height: 100vh;
  overflow-y: auto;
}

.main-wrapper { margin-left: 280px; }
.content { padding: 2rem; min-height: 100vh; }
.prose { max-width: 800px; margin: 0 auto; }

/* Copy button positioning */
.code-block { position: relative; }
.copy-button { position: absolute; top: 8px; right: 8px; }

/* Theme toggle positioning */
.theme-toggle { position: fixed; top: 1rem; right: 1rem; }

/* Mobile */
@media (max-width: 768px) {
  .sidebar { transform: translateX(-100%); }
  .main-wrapper { margin-left: 0; }
}
```

### Files to Create
- `internal/styles/themes/vanilla.css`

### Acceptance Criteria
- [ ] No color definitions in vanilla.css
- [ ] All layout/positioning preserved
- [ ] Site remains functional and navigable
- [ ] Copy buttons and toggles are positioned correctly
- [ ] Mobile responsive behavior works

---

## Story 41: CSS Export Command

### Description
Add a new command to output the vanilla CSS skeleton to stdout or a file. This allows users to use it as a starting point for custom styling.

### CLI Design
```bash
# Output to stdout
volcano css

# Output to file
volcano css > custom.css
volcano css -o custom.css
```

### Implementation Details

**1. Add command detection in `main.go`:**
```go
// Before flag parsing, check for subcommand
if len(args) > 0 && args[0] == "css" {
    return runCSSCommand(args[1:], stdout)
}
```

**2. Create `cmd/css.go`:**
```go
package cmd

import (
    "flag"
    "io"
    "os"
    "volcano/internal/styles"
)

func CSS(args []string, w io.Writer) error {
    fs := flag.NewFlagSet("css", flag.ContinueOnError)
    var outputFile string
    fs.StringVar(&outputFile, "o", "", "Output file path")
    fs.StringVar(&outputFile, "output", "", "Output file path")

    if err := fs.Parse(args); err != nil {
        return err
    }

    css := styles.GetCSS("vanilla")

    if outputFile != "" {
        return os.WriteFile(outputFile, []byte(css), 0644)
    }

    _, err := w.Write([]byte(css))
    return err
}
```

**3. Update help/usage:**
```
Commands:
  volcano <input> [flags]    Generate static site
  volcano css [-o file]      Output vanilla CSS skeleton
```

### Files to Modify
- `main.go`: Add command detection and routing
- Create `cmd/css.go`: CSS export logic

### Acceptance Criteria
- [ ] `volcano css` outputs vanilla CSS to stdout
- [ ] `volcano css -o file.css` writes to file
- [ ] CSS output is complete vanilla theme
- [ ] Help shows css command

---

## Story 42: Custom CSS Path Flag

### Description
Add `--css <path>` flag for serve and build commands to load custom CSS from a file instead of embedded themes.

### CLI Design
```bash
# Build with custom CSS
volcano docs --css ./custom.css

# Serve with custom CSS (reloaded on each request)
volcano -s docs --css ./custom.css
```

### Implementation Details

**1. Add to Config:**
```go
type Config struct {
    // ... existing fields
    CSSPath string  // Path to custom CSS file
}
```

**2. Add flag in `main.go`:**
```go
fs.StringVar(&cfg.CSSPath, "css", "", "Path to custom CSS file")
```

**3. Update `isValueFlag()`:**
```go
valueFlags := map[string]bool{
    // ... existing
    "css": true,
}
```

**4. CSS loading logic in generator:**
```go
func (g *Generator) getCSS() (string, error) {
    if g.config.CSSPath != "" {
        content, err := os.ReadFile(g.config.CSSPath)
        if err != nil {
            return "", fmt.Errorf("failed to load CSS: %w", err)
        }
        return string(content), nil
    }
    return styles.GetCSS(g.config.Theme), nil
}
```

**5. Dynamic server live reload:**
In `internal/server/dynamic.go`, read CSS file on each request:
```go
func (s *DynamicServer) getCSS() string {
    if s.config.CSSPath != "" {
        content, err := os.ReadFile(s.config.CSSPath)
        if err != nil {
            // Log error, fall back to theme
            return styles.GetCSS(s.config.Theme)
        }
        return string(content)
    }
    return styles.GetCSS(s.config.Theme)
}
```

### Files to Modify
- `cmd/config.go`: Add `CSSPath` field
- `main.go`: Add `--css` flag
- `cmd/serve.go`: Pass CSSPath to dynamic server
- `internal/server/dynamic.go`: Add CSSPath to config, implement live reload
- `internal/generator/generator.go`: Load custom CSS in build

### Acceptance Criteria
- [ ] `--css path` loads CSS from file for build
- [ ] `--css path` loads CSS from file for serve
- [ ] In serve mode, CSS changes reflect immediately (no restart)
- [ ] Missing CSS file produces clear error
- [ ] `--css` and `--theme` are mutually exclusive (error if both)

---

## Story 43: Page Navigation Behind Flag

### Description
Move the prev/next page navigation links to the bottom of the page and make them opt-in via a flag.

### Current State
- `PageNav` is always rendered in layout.html
- Placed after `</article>` in the main content area

### Implementation Details

**1. Add flag:**
```go
// config.go
ShowPageNav bool  // Show prev/next navigation

// main.go
fs.BoolVar(&cfg.ShowPageNav, "page-nav", false, "Show previous/next page navigation")
```

**2. Conditional rendering:**
In generator and dynamic server, only populate `PageNav` field when flag is enabled:
```go
if g.config.ShowPageNav {
    pageNav := navigation.BuildPageNavigation(node, allPages)
    data.PageNav = navigation.RenderPageNavigation(pageNav)
}
```

**3. Layout already has correct placement:**
```html
{{.PageNav}}  <!-- After article, at bottom of content -->
```

### Files to Modify
- `cmd/config.go`: Add `ShowPageNav bool`
- `main.go`: Add `--page-nav` flag
- `cmd/serve.go`: Pass to dynamic server config
- `internal/server/dynamic.go`: Add to DynamicConfig, conditionally render
- `internal/generator/generator.go`: Conditionally render

### Acceptance Criteria
- [ ] By default, no prev/next navigation shown
- [ ] `--page-nav` flag enables prev/next links
- [ ] Links appear at bottom of page content
- [ ] Works in both build and serve modes

---

## Story 44: TOC Smooth Scroll

### Description
Make TOC links scroll smoothly to headings, and ensure headers scroll to the top of viewport even when near the bottom of the page.

### Current State
- `html { scroll-behavior: smooth; }` is already set
- TOC links use `href="#heading-id"` anchors
- `h1-h6 { scroll-margin-top: 80px; }` for fixed header offset

### Problem
When clicking a TOC link to a header near the bottom of the page, the header doesn't scroll to the top because there's not enough content after it.

### Solution
Add JavaScript to handle TOC link clicks with custom scroll behavior:

```javascript
// In layout.html script section
(function() {
    const tocLinks = document.querySelectorAll('.toc a[href^="#"]');

    tocLinks.forEach(function(link) {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const targetId = this.getAttribute('href').slice(1);
            const target = document.getElementById(targetId);
            if (!target) return;

            // Calculate scroll position to put heading near top
            const headerOffset = 80;
            const targetPosition = target.getBoundingClientRect().top + window.scrollY - headerOffset;

            window.scrollTo({
                top: targetPosition,
                behavior: 'smooth'
            });

            // Update URL hash without jumping
            history.pushState(null, null, '#' + targetId);
        });
    });
})();
```

### Additional CSS for visibility
Ensure there's enough scroll room:
```css
/* Add padding at bottom of content to allow last heading to scroll to top */
.prose {
    padding-bottom: calc(100vh - 200px);
}

/* Or use min-height on article */
article.prose {
    min-height: calc(100vh - 100px);
}
```

### Files to Modify
- `internal/templates/layout.html`: Add click handler for TOC links
- `internal/styles/styles.css` (and theme variants): Add bottom padding

### Acceptance Criteria
- [ ] Clicking TOC link smoothly scrolls to heading
- [ ] Heading ends up near top of viewport (with offset for fixed elements)
- [ ] Works even for headings near bottom of page
- [ ] URL hash updates correctly
- [ ] Works with `prefers-reduced-motion` (instant scroll)
- [ ] Works in all themes

---

## Implementation Order

1. **Story 38** (Theme Architecture) - Foundation for all theme work
2. **Story 40** (Vanilla Theme) - Needed for CSS export
3. **Story 41** (CSS Export Command) - Depends on vanilla theme
4. **Story 39** (Blog Theme) - Independent, can be done anytime after 38
5. **Story 42** (Custom CSS Path) - Independent after 38
6. **Story 43** (Page Nav Flag) - Independent
7. **Story 44** (TOC Smooth Scroll) - Independent

Stories 38-42 are related (themes/CSS). Stories 43-44 are independent.

---

## Verification

After implementation, test:
1. Build with each theme: `--theme=docs`, `--theme=blog`, `--theme=vanilla`
2. Export CSS: `volcano css > test.css`
3. Build with custom CSS: `volcano docs --css test.css`
4. Serve with custom CSS, modify file, verify changes appear
5. Enable page nav: `volcano docs --page-nav`
6. Click TOC links, verify smooth scroll to all headings

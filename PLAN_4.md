# PLAN_4.md - Bug Fixes and Improvements

## Summary of Issues Found

1. **Front matter rendering** - YAML front matter is rendered as content instead of being stripped
2. **Blog theme styling** - Many elements use sans-serif instead of serif, divider lines present
3. **Tree links broken** - Output paths not slugified but URLs are, causing 404s
4. **TOC smooth scroll** - Needs investigation (code present but may have issues)

---

## Story 45: Strip YAML Front Matter

### Description
YAML front matter (delimited by `---`) at the beginning of markdown files is being rendered as content instead of being stripped. This causes ugly output when files have front matter metadata.

### Current Behavior
A file like:
```markdown
---
title: My Page
date: 2024-01-01
---

# Heading

Content here
```

Renders the `---` and YAML as visible content.

### Implementation Details

**Option A: Use goldmark-meta extension**
```go
import (
    meta "github.com/yuin/goldmark-meta"
)

md := goldmark.New(
    goldmark.WithExtensions(
        meta.Meta,  // Parses and strips front matter
        // ... other extensions
    ),
)
```

**Option B: Manual stripping before parsing**
Add a function in `internal/markdown/frontmatter.go`:
```go
package markdown

import (
    "bytes"
    "regexp"
)

var frontMatterRegex = regexp.MustCompile(`(?s)^---\n(.+?)\n---\n?`)

// StripFrontMatter removes YAML front matter from markdown content
func StripFrontMatter(content []byte) []byte {
    return frontMatterRegex.ReplaceAll(content, nil)
}
```

Then call it in `ParseContent` before parsing:
```go
func ParseContent(content []byte, ...) (*Page, error) {
    // Strip front matter before processing
    content = StripFrontMatter(content)

    // Extract title from markdown
    title := ExtractTitle(content)
    // ...
}
```

**Recommendation**: Option B (manual stripping) - simpler, no new dependency, we don't need to parse the YAML values.

### Files to Modify
- Create `internal/markdown/frontmatter.go`
- Modify `internal/markdown/page.go` - call StripFrontMatter in ParseContent
- Add tests in `internal/markdown/frontmatter_test.go`

### Acceptance Criteria
- [ ] YAML front matter delimited by `---` is stripped from output
- [ ] Content after front matter renders correctly
- [ ] Files without front matter work unchanged
- [ ] Front matter at non-start of file is not stripped (only beginning)

---

## Story 46: Blog Theme - Full Serif Styling

### Description
The blog theme should use serif fonts everywhere except for code elements. Currently many UI elements incorrectly use sans-serif fonts.

### Current Issues
Elements incorrectly using sans-serif:
- `.site-title` (line 119)
- `.search-input` (line 163)
- `.breadcrumbs` (line 194)
- `.toc-title` (line 303)
- Headings h1-h6 (line 387)
- `.page-nav` (line 490)
- `.page-meta` (line 534, 556)
- `.back-to-top` (line 615)
- `kbd` (line 748)
- `blockquote` (line 804)
- `.admonition` (line 856)
- `table` (line 1005)
- Code elements (line 965) - should stay monospace

### Implementation Details
Remove ALL explicit `font-family` declarations except for code elements, allowing everything to inherit the serif body font. Keep only:
- Code/pre: monospace (functional requirement)
- Everything else including headings: inherit from body (serif)

### Files to Modify
- `internal/styles/themes/blog.css`

### Acceptance Criteria
- [ ] Body text uses serif font (Georgia)
- [ ] Headings use sans-serif font (intentional contrast)
- [ ] Code, pre, kbd use monospace font
- [ ] All other UI elements (nav, breadcrumbs, TOC, buttons) use serif
- [ ] No visual regressions in functionality

---

## Story 47: Blog Theme - Remove Divider Lines

### Description
The blog theme should have a cleaner, less cluttered appearance without divider lines separating sections.

### Current Issues
Divider lines present on:
- `.sidebar-header` border-bottom (line 114)
- `.nav-search` border-bottom (line 136)
- `.breadcrumbs` border-bottom (line 286)
- `.page-nav` border-top (line 613)
- `hr` elements (line 996-997)

### Implementation Details
Remove or make transparent the border properties on these elements:
```css
.sidebar-header {
    border-bottom: none;
}

.nav-search {
    border-bottom: none;
}

.breadcrumbs {
    border-bottom: none;
}

.page-nav {
    border-top: none;
}
```

Keep `hr` element styling minimal but visible (it's intentional content).

### Files to Modify
- `internal/styles/themes/blog.css`

### Acceptance Criteria
- [ ] No border between sidebar header and content
- [ ] No border below search input
- [ ] No border below breadcrumbs
- [ ] No border above page navigation
- [ ] `hr` elements remain visible but subtle

---

## Story 48: Fix Tree Link Resolution (Critical Bug)

### Description
Navigation links in the tree are broken because `GetURLPath` slugifies paths but `GetOutputPath` does not. This causes all links to folders/files with special characters, numbers, or spaces to 404.

### Root Cause
In `internal/tree/scanner.go`:
- `GetURLPath` calls `SlugifyPath(dir)` to create URL-friendly paths
- `GetOutputPath` uses `dir` directly without slugification

Result:
- URL: `/archive/2023-easter-break/` (slugified)
- File: `/8. Archive/2023 Easter Break/index.html` (original)

### Implementation Details

**Fix GetOutputPath to slugify directory paths:**

```go
func GetOutputPath(node *Node) string {
    if node.IsFolder {
        return ""
    }

    // Get directory and filename
    dir := filepath.Dir(node.Path)
    filename := filepath.Base(node.Path)

    // Slugify the directory path to match URLs
    slugDir := SlugifyPath(dir)

    // Handle index files - they stay as index.html
    stem := strings.TrimSuffix(filename, filepath.Ext(filename))
    if strings.ToLower(stem) == "index" || strings.ToLower(stem) == "readme" {
        if slugDir == "" || slugDir == "." {
            return "index.html"
        }
        return filepath.Join(slugDir, "index.html")
    }

    // Extract metadata to get slug (strips date/number prefixes)
    meta := ExtractFileMetadata(filename, node.ModTime())
    slug := meta.Slug

    // For non-index files, create clean URLs: file.md → file/index.html
    if slugDir == "" || slugDir == "." {
        return filepath.Join(slug, "index.html")
    }
    return filepath.Join(slugDir, slug, "index.html")
}
```

**Add link verification step:**
Add a verification pass after generation that checks all navigation links resolve to actual files:

```go
// In generator.go, after page generation
func (g *Generator) verifyLinks(site *tree.Site) []string {
    var broken []string
    for _, node := range site.AllPages {
        urlPath := tree.GetURLPath(node)
        outputPath := tree.GetOutputPath(node)
        fullPath := filepath.Join(g.config.OutputDir, outputPath)

        if _, err := os.Stat(fullPath); os.IsNotExist(err) {
            broken = append(broken, fmt.Sprintf("%s -> %s (file: %s)",
                node.Path, urlPath, outputPath))
        }
    }
    return broken
}
```

### Files to Modify
- `internal/tree/scanner.go` - Fix GetOutputPath
- `internal/generator/generator.go` - Add verification step
- Update tests in `internal/tree/scanner_test.go`

### Acceptance Criteria
- [ ] Output paths are slugified to match URL paths
- [ ] All navigation links resolve correctly
- [ ] Verification step warns/fails on broken links
- [ ] Works with Obsidian vault (~/Obsidian/Primary)
- [ ] Existing tests updated/passing

---

## Story 49: TOC Smooth Scroll Fix

### Description
TOC links scroll instantly instead of smoothly. The JavaScript intercepts clicks but `window.scrollTo` with `behavior: 'smooth'` isn't working as expected.

### Root Cause Analysis
The current implementation:
1. Intercepts click with `e.preventDefault()` - blocks native behavior
2. Calls `window.scrollTo({ behavior: 'smooth' })` - should animate
3. But scroll happens instantly

Likely issue: The JavaScript `scrollTo` behavior option may not be working due to:
- CSS `scroll-behavior: smooth` on `html` conflicting with JS
- Browser quirks with programmatic smooth scroll

### Solution: Remove JavaScript, Use Pure CSS

The simplest and most reliable approach is to NOT intercept TOC clicks and let the browser handle it natively with CSS:

```css
html {
    scroll-behavior: smooth;
}

@media (prefers-reduced-motion: reduce) {
    html {
        scroll-behavior: auto;
    }
}

h1, h2, h3, h4, h5, h6 {
    scroll-margin-top: 80px;
}
```

Remove the JavaScript TOC smooth scroll handler entirely. The native anchor links (`<a href="#heading-id">`) will:
1. Scroll smoothly (CSS handles animation)
2. Respect reduced motion preference
3. Position correctly (scroll-margin-top handles offset)
4. Update URL hash automatically

### Files to Modify
- `internal/templates/layout.html` - Remove TOC smooth scroll JavaScript block
- `internal/styles/themes/*.css` - Ensure scroll-margin-top is set on headings

### Acceptance Criteria
- [ ] Clicking TOC links scrolls smoothly to heading
- [ ] Heading ends up visible below any fixed headers
- [ ] URL hash updates when clicking TOC links
- [ ] Respects prefers-reduced-motion
- [ ] Works in Chrome, Firefox, Safari

---

## Implementation Order

1. **Story 48** (Critical) - Fix tree links, this is a blocking bug
2. **Story 45** - Strip front matter, common issue with Obsidian
3. **Story 46** - Blog theme serif styling
4. **Story 47** - Blog theme divider removal
5. **Story 49** - TOC smooth scroll (needs user feedback)

---

## Testing

After implementation:
1. Build with Obsidian vault: `volcano ~/Obsidian/Primary -o /tmp/test`
2. Verify no 404s when clicking navigation links
3. Check front matter is stripped from rendered pages
4. Test blog theme: `volcano ~/Obsidian/Primary --theme=blog`
5. Verify TOC scrolling works in browser

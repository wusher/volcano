# Plan 7 - Brainstorming Session

**Date**: 2026-01-10
**Branch**: claude/brainstorm-stories-m8YXR

---

## Ideas & Notes

### 1. Accent Color Flag

**Concept**: Add CLI flag like `--accent-color="#ff6600"` to customize theme colors

**Current Theme Architecture**:
- Themes use CSS custom properties (CSS variables) extensively
- Dark/light mode handled via `[data-theme="dark"]` selectors
- `layout.css` = structure only, theme files = colors/styling
- Already have variables like `--text-primary`, `--link-color`, `--border-color`

**Implementation Approach**:
- Add `--accent-color` flag (maybe also `--accent-color-dark` for dark mode)
- Inject as CSS variable in template: `:root { --accent-color: #value; }`
- Update themes to use `var(--accent-color, fallback)` for:
  - Active nav items
  - TOC active border
  - Link hover states
  - Scroll progress bar
  - Possibly breadcrumbs, admonitions

**Dark/Light Mode Challenge**:
- Single color may not work well in both modes
- Options:
  1. Same color in both (simple, might have contrast issues)
  2. Auto-derive lighter/darker variant (complex, results may vary)
  3. Two separate flags: `--accent-color-light` and `--accent-color-dark`
  4. CSS filters to adjust (hacky)

**Pros**:
- Easy customization without writing full custom themes
- Match brand colors
- Non-breaking (can use defaults)
- Relatively simple to implement

**Cons**:
- Dark/light mode handling is tricky
- Potential accessibility/contrast issues if poorly chosen
- Adds CLI complexity
- Might lead to requests for more customization (fonts, spacing, etc.)
- The monochrome aesthetic might be intentional design choice

**Alternative Approaches**:
- Pre-made accent variants ("blue", "green", "orange") instead of custom colors
- Document custom CSS approach instead (already possible)
- Skip this in favor of other features

**Decision Direction**: Skip presets, allow custom brand colors

**Two viable approaches**:

**Option A: HSL Auto-adjustment**
- Convert hex input to HSL (Hue, Saturation, Lightness)
- For dark mode: automatically lighten the color (increase L value)
- For light mode: use as-is or slightly darken if needed
- Pros: One flag, smart adaptation
- Cons: Need HSL conversion logic in Go, results may vary by color

**Option B: Keep It Simple**
- Single `--accent-color` flag
- Use same color in both light and dark modes
- Let user provide a color that works reasonably in both
- Pros: Zero complexity, user has full control
- Cons: User needs to test in both modes, some colors won't work well

**Option C: Generate Tints & Shades (Proper Color Space)**
- User provides base color: `--accent-color="#ff6600"`
- Convert to better color space (HSL, LAB, or LCH - NOT RGB)
- Compute variants:
  - **Dark shade**: Set lightness to ~10% (keep hue/saturation)
  - **Light tint**: Set lightness to ~95% (keep hue/saturation)
  - **Base**: Use original or normalized (50% lightness)
- Creates 3 CSS variables: `--accent`, `--accent-dark`, `--accent-light`
- Use cases:
  - `--accent` for main interactive elements (links, active states, scroll bar)
  - `--accent-light` for subtle backgrounds in light mode (code blocks, admonitions)
  - `--accent-dark` for subtle backgrounds in dark mode
- **Color space options**:
  - **HSL**: Simple, widely understood, good enough for tints/shades
  - **Oklab/LCH**: Perceptually uniform, better color mixing, but more complex math
- Pros: One input creates mini palette, subtle branded backgrounds, flexible, proper color mixing
- Cons: Need color space conversion, more variables to manage, ~100 lines of code

**When flag is NOT provided**:
- **Option 1**: No accent color variables injected at all (themes use their defaults)
- **Option 2**: Inject a default accent color value
- **Leaning towards**: Option 1 (no default) - keeps themes working exactly as they do now, purely opt-in feature

**Verdict**: **OPTION C** - Generate tints & shades using HSL

**Decision**:
- Single `--accent-color="#ff6600"` flag
- Convert to HSL and generate 3 variants:
  - `--accent`: Original color (or normalized to 50% lightness)
  - `--accent-dark`: ~10% lightness (for dark mode backgrounds/accents)
  - `--accent-light`: ~95% lightness (for light mode backgrounds/accents)
- Themes can use these variables where appropriate
- No variables injected if flag not provided (purely opt-in)

**Next Steps**:
1. Implement HSL conversion utilities (hex → HSL → hex)
2. Add `--accent-color` flag to Config struct
3. Generate 3 CSS variables in template
4. Update themes to use accent variables
5. Test with various brand colors in both modes

---

### 2. Icon-Only Copy Buttons

**Concept**: Remove "Copy" and "Copied!" text from code block copy buttons, show only icons

**Current Behavior**:
- Copy button shows text + icon: "Copy 📋" or "Copied! ✓"
- Button toggles between these states on click

**Proposed Change**:
- Show only icons: 📋 → ✓
- No text labels
- Cleaner, more minimal appearance
- Icons are already clear/universal

**Pros**:
- Cleaner visual design
- Less visual clutter
- Icons are self-explanatory
- More modern pattern (GitHub, VS Code do this)
- Smaller button footprint

**Cons**:
- Slightly less accessible (no explicit text label)
- Could add aria-label for screen readers

**Implementation**:
- Remove text nodes from `.copy-button` in template
- Keep icon swap logic (copy-icon ↔ check-icon)
- Add `aria-label="Copy code"` / `aria-label="Copied!"` for accessibility

**Complexity**: Very simple, probably 5 lines changed

---

### 3. Missing Meta Tags - Browser Theme Color

**Concept**: Add missing meta tags, particularly `theme-color` for mobile browser chrome

**Current Meta Tags** (from `internal/seo/meta.go`):
- ✅ `charset="UTF-8"`
- ✅ `name="viewport"`
- ✅ `name="description"`
- ✅ `name="robots"`
- ✅ `name="author"`
- ✅ `rel="canonical"`
- ✅ Open Graph tags (og:title, og:description, og:type, og:url, og:site_name, og:image)
- ✅ Twitter Card tags (twitter:card, twitter:title, twitter:description, twitter:image)

**Missing/Recommended Meta Tags**:

1. **`theme-color`** (HIGH PRIORITY)
   - Colors mobile browser's address bar/UI chrome
   - Example: `<meta name="theme-color" content="#ffffff">`
   - Should match theme background colors
   - Can use `media` attribute for dark mode: `<meta name="theme-color" content="#1a1a1a" media="(prefers-color-scheme: dark)">`

2. **`color-scheme`** (MEDIUM PRIORITY)
   - Tells browser the page supports light/dark mode
   - Example: `<meta name="color-scheme" content="light dark">`
   - Improves native browser UI consistency

3. **Apple iOS Meta Tags** (NICE TO HAVE)
   - `apple-mobile-web-app-capable` - enables iOS standalone mode
   - `apple-mobile-web-app-status-bar-style` - iOS status bar styling
   - Only useful if users want PWA-like behavior

4. **Windows Tile Color** (NICE TO HAVE)
   - `msapplication-TileColor` - Windows Start menu tile color
   - Less critical nowadays

**Proposed Implementation**:
- Add `theme-color` meta tags to `RenderMetaTags()` in `internal/seo/meta.go`
- Use light/dark variants:
  - Light: `#ffffff` or match theme background
  - Dark: `#1a1a1a` or match dark theme background
- Add `color-scheme: light dark` meta tag
- Make configurable via Config if needed (or use hardcoded theme colors)

**Interaction with Accent Color Feature**:
- If accent color flag is provided, could optionally use accent for theme-color
- Or stick with neutral background colors (safer default)

**Pros**:
- Better mobile browser experience
- Professional polish
- Colors match site theme automatically

**Cons**:
- Need to maintain color values
- Should match CSS theme colors

**Complexity**: Simple, ~10-20 lines of code

---

### 4. Full-Text Search (Client-Side)

**Concept**: Add full-text search across all page content, not just navigation titles

**Current Search**:
- Only searches navigation tree by page titles (`data-search-text` attribute)
- Filters/hides nav items that don't match
- No content search capability

**Proposed Client-Side Search**:
- Generate a search index at build time (JSON file)
- Include page titles, content, URLs
- **Lazy-load index**: Only fetch `search-index.json` when user opens search (not on every page load)
- Search in browser using JavaScript
- Display results in modal or dedicated UI

**Lazy-Loading Strategy**:
- Don't include search index in initial page load
- When user presses `/` or clicks search button, fetch the index
- Cache in memory for rest of session
- First search has ~200ms delay (download), subsequent searches are instant
- Only users who search pay the bandwidth cost

**Implementation Approaches**:

**Option A: Custom Lightweight Search**
- Build simple JSON index: `[{title, url, content, excerpt}, ...]`
- ~100 lines of custom JS for searching (string matching, ranking)
- Pros: Zero dependencies (aligns with Volcano philosophy), small, full control
- Cons: Basic ranking, slower on huge sites, need to write search logic

**Option B: Lunr.js (Popular Static Site Library)**
- Pre-build Lunr index at generation time
- Bundle serialized index with site
- Load and search with Lunr.js (~40KB minified)
- Pros: Good ranking, stemming, fast, well-tested
- Cons: Adds external dependency, larger index files

**Option C: Fuse.js (Fuzzy Search)**
- Generate simple JSON index
- Include Fuse.js (~15KB minified) for fuzzy matching
- Pros: Typo-tolerant, smaller than Lunr, simple API
- Cons: Still a dependency, can be slow on huge datasets (1000+ pages)

**Option D: FlexSearch**
- Lightweight (~6KB minified), extremely fast
- Claims to be faster and more memory-efficient than Fuse/Lunr
- Supports multi-language indexing and custom scoring
- Generate simple JSON index
- Pros: Very fast even on large datasets, small size, good balance
- Cons: Still a dependency, slightly more complex API than Fuse

**Option E: MiniSearch**
- Tiny (~8KB minified), full-text search with minimal resources
- Can add/remove documents from index dynamically
- Built-in support for prefix search, fuzzy matching, boosting
- Pros: Small size, flexible, good documentation, pure JS
- Cons: Still a dependency, newer/less battle-tested than Lunr

**Option F: ElasticLunr**
- Fork of Lunr.js with more features (query-time boosting, field search)
- Faster than original Lunr
- Pros: More flexible than Lunr, familiar API for Lunr users
- Cons: Larger than other options, still heavyweight

**Recommendation**:
- **For zero-dependency purists**: Option A (Custom) - keeps Volcano truly dependency-free
- **For best balance**: Option D (FlexSearch) or Option E (MiniSearch) - small, fast, modern
- **For proven stability**: Option B (Lunr) - battle-tested but heavier
- **If fuzzy matching is priority**: Option C (Fuse) - but watch performance on large sites

**Build-Time Index Generation**:
1. During `generator.Generate()`, collect all page data:
   - Title, URL, full text content (stripped HTML)
   - Maybe first 200 chars as excerpt
2. Write to `search-index.json` in output dir
3. Estimate size: ~1-2KB per page (depends on content length)
   - 100 pages = 100-200KB index (reasonable)
   - 1000 pages = 1-2MB (might need optimization)

**UI Considerations**:
- Add search modal/overlay (triggered by existing `/` shortcut or new button)
- Show results with title, excerpt, URL
- Highlight search terms in results
- Navigate to page on click
- Keep existing nav search OR replace it with full-text search

**Complexity**:
- **Option A (Custom)**: ~200-300 lines Go (indexing) + ~100-150 lines JS (search UI)
- **Option B/C (Library)**: ~150 lines Go + ~50-100 lines JS + external lib

**Zero-Dependency Consideration**:
- Volcano is marketed as "zero-dependency"
- Adding Lunr/Fuse doesn't add a *build* dependency (just a runtime JS file)
- Could bundle the JS directly (like CSS is embedded)
- Or recommend CDN link
- Custom implementation keeps it truly zero-dependency

**Pros**:
- Major feature for documentation sites
- Makes large sites much more usable
- Common user request for static site generators

**Cons**:
- Index file size grows with content
- Adds complexity to build process
- Slower builds (need to process all content)
- Need to maintain search UI/logic

**Similar SSGs that have this**:
- Hugo (has built-in search index generation)
- MkDocs (search plugin)
- Docusaurus (Algolia or local search)
- VuePress (built-in search)

**Verdict**: **NO** - Not pursuing full-text search at this time due to complexity

---

### 5. Instant Navigation with DOM Morphing

**Concept**: Make navigation instant by pre-fetching on hover and morphing DOM (not full replacement)

**What This Does**:
1. **Cache current page** in memory
2. **Pre-fetch links** when user hovers (200-300ms before click)
3. **Fetch in background** when clicked
4. **Morph DOM** - intelligently update only what changed, preserve state
5. **Update URL** with History API
6. **Result**: Instant navigation, no white flash, preserves JS state, smooth updates

**Key Technologies**:

### DOM Morphing Libraries

**1. idiomorph** (3.3KB, zero dependencies)
- Created by htmx team, now used by Turbo 8
- Morphs one DOM tree to another while preserving state
- Smart ID-based matching
- Works great for updating pages without losing state
- Pros: Tiny, standalone, well-tested, no dependencies
- Cons: Need to combine with pre-fetch solution

**2. Turbo 8 with Morph Mode** (~45KB with Stimulus)
- Uses idiomorph under the hood
- Pre-fetches on hover automatically
- Morphs DOM instead of full replacement (only updates what changed)
- Can ignore certain elements (keep popovers open, etc.)
- Pros: Complete solution, battle-tested, does everything
- Cons: Larger bundle, opinionated

**3. Swup with Morph Plugin** (~10-15KB total)
- Swup handles transitions + pre-fetching
- Morph plugin adds DOM morphing (uses morphdom)
- Good for multi-language sites, persistent headers/menus
- Pros: Smaller than Turbo, focused, flexible
- Cons: Two libraries to integrate

**4. Alpine Morph** (part of Alpine.js)
- Morphs elements while preserving Alpine/browser state
- Good if already using Alpine
- Pros: Integrates with Alpine ecosystem
- Cons: Requires Alpine.js

### Pre-fetching/Hover Libraries

**1. instant.page** (~1KB)
- Pre-fetches links on mouse hover (uses `<link rel=prefetch>`)
- ~300ms from hover to click = free time for fetching
- Pros: Tiny, simple, just prefetch
- Cons: Just prefetch, need to add morphing separately

**2. quicklink** (<1KB, by Google)
- Pre-fetches links in viewport (as soon as visible)
- Respects user preferences (data-saver mode, slow connection)
- Uses requestIdleCallback to be responsible
- Pros: Tiniest, smart about bandwidth, Google-backed
- Cons: Viewport-based (not hover), just prefetch

**3. Flying Pages**
- Combines both: viewport preloading + hover preloading
- Rate limiting (3 requests/sec default)
- More control than quicklink/instant.page
- Pros: Best of both worlds, rate limiting
- Cons: Slightly larger, more WordPress-focused

**4. InstantClick** (Proof of concept)
- Pre-fetch on hover + AJAX navigation
- Makes site into SPA
- Pros: Does both prefetch and navigation
- Cons: Mostly abandoned, poor docs, high GitHub issues

### Recommended Approaches for Volcano

**Option A: Turbo 8 with Morph Mode** (Complete, ~45KB)
- Does everything: pre-fetch on hover, DOM morphing, caching, state preservation
- Battle-tested, used in production by many sites
- Simple to integrate: just include the script
- Pros: Complete solution, just works, well-documented
- Cons: Larger bundle (45KB), most opinionated

**Option B: Swup + Morph Plugin** (Smaller, ~10-15KB)
- Swup handles pre-fetching and page transitions
- Morph plugin adds DOM morphing
- More lightweight than Turbo
- Pros: Smaller, flexible, good docs
- Cons: Two pieces to integrate, smaller community

**Option C: Custom with idiomorph + instant.page** (Smallest, ~4-5KB)
- idiomorph (3.3KB) for DOM morphing
- instant.page (1KB) for hover pre-fetching
- Write ~50-100 lines of glue code to connect them
- Pros: Tiniest bundle, full control, truly minimal
- Cons: Need to write integration code, less battle-tested combo

**Option D: Custom with idiomorph + quicklink** (Smallest, ~4KB)
- idiomorph (3.3KB) for DOM morphing
- quicklink (<1KB) for viewport-based prefetching
- More responsible with bandwidth (respects data-saver, slow connections)
- Pros: Tiniest bundle, Google-backed prefetch, smart about resources
- Cons: Viewport-based (not hover), need integration code

**Best Recommendation for Volcano**: **Option C (idiomorph + instant.page)**
- Smallest bundle (~4-5KB total)
- Hover pre-fetching matches Turbo's UX
- idiomorph is battle-tested (used by Turbo 8)
- Clean, minimal approach that fits Volcano's philosophy
- ~100 lines of integration code needed

**Implementation Considerations**:
- Add `--instant-nav` or `--morph-navigation` flag (opt-in)
- Embed minified JS directly (like CSS) to keep single-binary philosophy
- Integration code:
  1. Listen for hover events on links
  2. Pre-fetch with instant.page
  3. On click: intercept, fetch, use idiomorph to morph DOM
  4. Update URL with History API
  5. Handle edge cases: external links, downloads, anchors
- May need special handling for theme toggle, search state
- Test that existing JS (theme toggle, search, etc.) survives morphing

**Zero-Dependency Trade-off**:
- Adds ~4-5KB runtime dependency (but not build dependency)
- Can embed JS directly to keep "single binary" philosophy
- Small enough to align with Volcano's minimal approach
- idiomorph is from htmx team (same philosophy as Volcano)

**Pros**:
- **Instant navigation** - feels incredibly fast
- No white flash between pages
- Preserves JavaScript state (theme preference, search, etc.)
- Smooth, modern UX
- Small bundle size (4-5KB vs 45KB for Turbo)
- Pre-fetching on hover uses the 200-300ms "dead time" perfectly

**Cons**:
- Adds JavaScript dependency (~4-5KB)
- Need to write ~100 lines of integration code
- Slight complexity in testing
- Need to handle edge cases (external links, downloads, anchor links)
- Scripts run differently (need to re-init on morph)
- Breaking change for users with custom JS (need docs on how to adapt)

**Verdict**: **YES** - Pursue Option C (idiomorph + instant.page) behind CLI flag

**Implementation Plan**:
- CLI flag: `--instant-nav` (opt-in, disabled by default)
- Embed idiomorph (3.3KB) + instant.page (1KB) minified into binary
- Write ~100 lines of integration glue code
- Only inject JS when flag is enabled
- Document usage and edge cases

---

## Stories

### Story 1: Accent Color Customization

**Feature**: Add `--accent-color` flag to customize theme colors with HSL-based tint/shade generation

**Technical Specification**:

**1. Add CLI Flag**
- **File**: `cmd/config.go`
- **Change**: Add field to Config struct:
  ```go
  AccentColor string // Custom accent color in hex format (e.g., "#ff6600")
  ```
- **File**: `cmd/root.go` (or wherever flags are defined)
- **Change**: Add CLI flag:
  ```go
  rootCmd.Flags().StringVar(&cfg.AccentColor, "accent-color", "", "Custom accent color (hex format, e.g., '#ff6600')")
  ```

**2. Create Color Utilities Package**
- **New File**: `internal/color/hsl.go`
- **Functions to implement**:
  ```go
  package color

  // HSL represents a color in HSL color space
  type HSL struct {
      H float64 // Hue: 0-360
      S float64 // Saturation: 0-100
      L float64 // Lightness: 0-100
  }

  // HexToRGB converts hex string to RGB values (0-255)
  func HexToRGB(hex string) (r, g, b uint8, err error)

  // RGBToHSL converts RGB (0-255) to HSL
  func RGBToHSL(r, g, b uint8) HSL

  // HSLToRGB converts HSL to RGB (0-255)
  func HSLToRGB(hsl HSL) (r, g, b uint8)

  // RGBToHex converts RGB to hex string
  func RGBToHex(r, g, b uint8) string

  // GenerateAccentVariants creates accent, accent-dark, accent-light
  // Returns: (accent at 50% L, dark at 10% L, light at 95% L)
  func GenerateAccentVariants(hexColor string) (accent, dark, light string, err error)
  ```

**3. HSL Conversion Algorithm**
- **RGB to HSL** (reference implementation):
  ```go
  func RGBToHSL(r, g, b uint8) HSL {
      rf := float64(r) / 255.0
      gf := float64(g) / 255.0
      bf := float64(b) / 255.0

      max := math.Max(math.Max(rf, gf), bf)
      min := math.Min(math.Min(rf, gf), bf)
      delta := max - min

      var h, s, l float64
      l = (max + min) / 2.0

      if delta == 0 {
          h, s = 0, 0
      } else {
          if l < 0.5 {
              s = delta / (max + min)
          } else {
              s = delta / (2.0 - max - min)
          }

          switch max {
          case rf:
              h = ((gf - bf) / delta)
              if gf < bf {
                  h += 6
              }
          case gf:
              h = ((bf - rf) / delta) + 2
          case bf:
              h = ((rf - gf) / delta) + 4
          }
          h /= 6
      }

      return HSL{H: h * 360, S: s * 100, L: l * 100}
  }
  ```

- **HSL to RGB** (reference implementation):
  ```go
  func HSLToRGB(hsl HSL) (r, g, b uint8) {
      h := hsl.H / 360.0
      s := hsl.S / 100.0
      l := hsl.L / 100.0

      if s == 0 {
          gray := uint8(l * 255)
          return gray, gray, gray
      }

      var q float64
      if l < 0.5 {
          q = l * (1 + s)
      } else {
          q = l + s - l*s
      }
      p := 2*l - q

      r = uint8(hueToRGB(p, q, h+1.0/3.0) * 255)
      g = uint8(hueToRGB(p, q, h) * 255)
      b = uint8(hueToRGB(p, q, h-1.0/3.0) * 255)
      return
  }

  func hueToRGB(p, q, t float64) float64 {
      if t < 0 {
          t += 1
      }
      if t > 1 {
          t -= 1
      }
      if t < 1.0/6.0 {
          return p + (q-p)*6*t
      }
      if t < 1.0/2.0 {
          return q
      }
      if t < 2.0/3.0 {
          return p + (q-p)*(2.0/3.0-t)*6
      }
      return p
  }
  ```

**4. Inject CSS Variables**
- **File**: `internal/templates/renderer.go`
- **Modify**: `PageData` struct to include:
  ```go
  AccentColorCSS template.CSS // CSS variables for accent colors
  ```
- **Modify**: `NewRenderer()` or add new function:
  ```go
  func GenerateAccentColorCSS(accentColor string) (template.CSS, error) {
      if accentColor == "" {
          return "", nil
      }

      accent, dark, light, err := color.GenerateAccentVariants(accentColor)
      if err != nil {
          return "", err
      }

      css := fmt.Sprintf(`
  :root {
    --accent: %s;
    --accent-dark: %s;
    --accent-light: %s;
  }`, accent, dark, light)

      return template.CSS(css), nil
  }
  ```
- **File**: `internal/templates/layout.html`
- **Change**: Add accent color CSS in `<style>` block after main CSS:
  ```html
  <style>
  {{.CSS}}
  {{if .AccentColorCSS}}{{.AccentColorCSS}}{{end}}
  </style>
  ```

**5. Update Theme CSS Files**
- **Files**: All theme CSS files in `internal/styles/themes/`
  - `docs.css`
  - `blog.css`
  - `vanilla.css`

- **Changes**: Replace hardcoded colors with `var(--accent, fallback)`:

  **Active Navigation**:
  ```css
  /* Before */
  .tree-nav a.active {
      color: #0066cc;
      font-weight: 500;
  }

  /* After */
  .tree-nav a.active {
      color: var(--accent, #0066cc);
      font-weight: 500;
  }
  ```

  **TOC Active Border**:
  ```css
  /* Before */
  .toc a.active {
      border-left: 2px solid #0066cc;
  }

  /* After */
  .toc a.active {
      border-left: 2px solid var(--accent, #0066cc);
  }
  ```

  **Link Hover States**:
  ```css
  /* Before */
  .prose a:hover {
      color: #0066cc;
  }

  /* After */
  .prose a:hover {
      color: var(--accent, #0066cc);
  }
  ```

  **Scroll Progress Bar**:
  ```css
  /* Before */
  .scroll-progress-bar {
      background: #0066cc;
  }

  /* After */
  .scroll-progress-bar {
      background: var(--accent, #0066cc);
  }
  ```

**6. Update Generator**
- **File**: `internal/generator/generator.go`
- **Modify**: Where `PageData` is populated:
  ```go
  accentColorCSS, err := templates.GenerateAccentColorCSS(g.config.AccentColor)
  if err != nil {
      log.Printf("Warning: failed to generate accent color CSS: %v", err)
      accentColorCSS = ""
  }

  pageData := templates.PageData{
      // ... existing fields ...
      AccentColorCSS: accentColorCSS,
  }
  ```

**Testing**:

**Unit Tests** (`internal/color/hsl_test.go`):
```go
func TestHexToRGB(t *testing.T) {
    tests := []struct {
        hex string
        r, g, b uint8
        wantErr bool
    }{
        {"#ff6600", 255, 102, 0, false},
        {"#000000", 0, 0, 0, false},
        {"#ffffff", 255, 255, 255, false},
        {"invalid", 0, 0, 0, true},
    }
    // ... test implementation
}

func TestGenerateAccentVariants(t *testing.T) {
    accent, dark, light, err := GenerateAccentVariants("#ff6600")
    assert.NoError(t, err)
    assert.NotEmpty(t, accent)
    assert.NotEmpty(t, dark)
    assert.NotEmpty(t, light)
    // Verify lightness values are correct
}
```

**Integration Tests**:
1. Generate site with `--accent-color="#ff6600"`
2. Verify generated HTML contains CSS variables
3. Verify all themes display correctly with accent color
4. Test with various colors: `#0066cc`, `#ff0000`, `#00ff00`, `#8b00ff`
5. Verify no accent variables present when flag not provided

**Manual Testing**:
1. Generate site with accent color
2. Open in browser, verify:
   - Active nav items use accent color
   - TOC active border uses accent color
   - Link hover uses accent color
   - Scroll progress bar uses accent color
3. Toggle light/dark mode, verify colors work in both
4. Test with extreme colors (very dark, very light, saturated)

**Acceptance Criteria**:
- ✅ Works with various brand colors (#ff6600, #0066cc, #8b00ff, etc.)
- ✅ Colors look good in both light/dark modes
- ✅ No visual change when flag not provided (backward compatible)
- ✅ Generated CSS variables have correct lightness values (50%, 10%, 95%)
- ✅ All themes (docs, blog, vanilla) use accent color when provided
- ✅ Invalid hex colors show error message
- ✅ Color conversion is accurate (HSL ↔ RGB ↔ Hex round-trip)

**Documentation Updates**:
- Update CLAUDE.md with `--accent-color` flag usage
- Add examples with common brand colors
- Document that colors should have reasonable contrast in both modes

**Estimated Effort**: Medium (200-300 lines: ~100 color utils, ~50 integration, ~50 CSS updates, 100 tests)

---

### Story 2: Icon-Only Copy Buttons

**Feature**: Simplify code block copy buttons to show only icons (no text)

**Technical Specification**:

**1. Update HTML Template**
- **File**: `internal/markdown/codeblock.go` (or wherever copy button HTML is generated)
- **Current Code** (approximate):
  ```html
  <button class="copy-button" aria-label="Copy code">
    <svg class="copy-icon">...</svg>
    <span class="copy-text">Copy</span>
    <svg class="check-icon">...</svg>
  </button>
  ```
- **New Code**:
  ```html
  <button class="copy-button" aria-label="Copy code">
    <svg class="copy-icon" aria-hidden="true">...</svg>
    <svg class="check-icon" aria-hidden="true">...</svg>
  </button>
  ```
- **Changes**:
  1. Remove `<span class="copy-text">Copy</span>` element
  2. Add `aria-label="Copy code"` to button
  3. Add `aria-hidden="true"` to SVG icons (accessibility best practice)

**2. Update JavaScript**
- **File**: `internal/templates/layout.html` (in the `<script>` section)
- **Current JavaScript** (lines 349-364):
  ```javascript
  // Copy code button
  document.querySelectorAll('.copy-button').forEach(function(button) {
      button.addEventListener('click', async function() {
          const code = this.parentElement.querySelector('code').textContent;
          try {
              await navigator.clipboard.writeText(code);
              this.classList.add('copied');
              this.querySelector('.copy-text').textContent = 'Copied!';  // REMOVE THIS LINE
              setTimeout(function() {
                  button.classList.remove('copied');
                  button.querySelector('.copy-text').textContent = 'Copy';  // REMOVE THIS LINE
              }, 2000);
          } catch (err) {
              console.error('Failed to copy:', err);
          }
      });
  });
  ```
- **New JavaScript**:
  ```javascript
  // Copy code button
  document.querySelectorAll('.copy-button').forEach(function(button) {
      button.addEventListener('click', async function() {
          const code = this.parentElement.querySelector('code').textContent;
          try {
              await navigator.clipboard.writeText(code);
              this.classList.add('copied');
              this.setAttribute('aria-label', 'Copied!');
              setTimeout(function() {
                  button.classList.remove('copied');
                  button.setAttribute('aria-label', 'Copy code');
              }, 2000);
          } catch (err) {
              console.error('Failed to copy:', err);
          }
      });
  });
  ```
- **Changes**:
  1. Replace `.querySelector('.copy-text').textContent = 'Copied!'` with `.setAttribute('aria-label', 'Copied!')`
  2. Replace `.querySelector('.copy-text').textContent = 'Copy'` with `.setAttribute('aria-label', 'Copy code')`

**3. Update CSS (Optional - Cleanup)**
- **File**: `internal/styles/themes/*.css` (all theme files if needed)
- **Remove** (if present):
  ```css
  .copy-text {
      font-size: 14px;
      margin-left: 4px;
  }
  ```
- **Note**: CSS for `.copy-button` icon styling should remain unchanged

**Testing**:

**Manual Testing**:
1. Generate a site with code blocks
2. Open in browser
3. Verify copy buttons:
   - ✅ Show only icon (clipboard icon)
   - ✅ No text "Copy" visible
   - ✅ Hover shows button is clickable
   - ✅ Click button, icon changes to checkmark
   - ✅ No text "Copied!" visible
   - ✅ After 2 seconds, icon reverts to clipboard
4. Test with screen reader (or browser's accessibility inspector):
   - ✅ Initial state announces "Copy code"
   - ✅ After click announces "Copied!"
   - ✅ After 2 seconds announces "Copy code" again
5. Test in multiple browsers:
   - Chrome/Edge
   - Firefox
   - Safari
6. Test on mobile devices:
   - Android Chrome
   - iOS Safari

**Accessibility Testing**:
- Use browser dev tools accessibility inspector
- Verify `aria-label` is present and changes on click
- Verify SVG icons have `aria-hidden="true"`
- Test with NVDA (Windows) or VoiceOver (Mac/iOS)

**Visual Regression**:
- Compare before/after screenshots
- Verify icon size and positioning unchanged
- Verify button alignment in code blocks
- Verify hover states work correctly

**Acceptance Criteria**:
- ✅ Buttons show only icons (no text)
- ✅ Screen readers announce "Copy code" initially
- ✅ Screen readers announce "Copied!" after click
- ✅ Icon swaps from clipboard to checkmark on click
- ✅ Icon reverts after 2 seconds
- ✅ Visual appearance is clean and modern
- ✅ Works in all major browsers
- ✅ Works on mobile devices
- ✅ Accessibility audit passes (no warnings)

**Documentation Updates**:
- No user-facing documentation needed (visual change only)
- Update internal comments if needed

**Estimated Effort**: Trivial (10-15 lines changed, 30 min testing)

---

### Story 3: Browser Theme Color Meta Tags

**Feature**: Add `theme-color` and `color-scheme` meta tags for better mobile browser integration

**Technical Specification**:

**1. Update SEO Meta Package**
- **File**: `internal/seo/meta.go`
- **Modify**: `RenderMetaTags()` function

**2. Add Theme Color Meta Tags**
- **Location**: In `RenderMetaTags()` function, after existing meta tags
- **Code to Add**:
  ```go
  // Theme color meta tags (after line ~167, after canonical link)
  sb.WriteString("\n")
  sb.WriteString(`  <!-- Browser Theme Colors -->`)
  sb.WriteString("\n")

  // Light mode theme color
  sb.WriteString(`  <meta name="theme-color" content="#ffffff">`)
  sb.WriteString("\n")

  // Dark mode theme color
  sb.WriteString(`  <meta name="theme-color" content="#1a1a1a" media="(prefers-color-scheme: dark)">`)
  sb.WriteString("\n")

  // Color scheme support
  sb.WriteString(`  <meta name="color-scheme" content="light dark">`)
  sb.WriteString("\n")
  ```

**3. Theme Color Values**
- **Light Mode**: `#ffffff` (white background)
  - Matches light mode background in all themes
- **Dark Mode**: `#1a1a1a` (dark gray background)
  - Matches dark mode background in docs/blog themes
  - Alternative: `#0d1117` (GitHub dark), `#1e1e1e` (VS Code dark)
- **Recommendation**: Use `#1a1a1a` for consistency with existing themes

**4. Verify Theme Background Colors**
- **Check Files**: `internal/styles/themes/docs.css`, `blog.css`, `vanilla.css`
- **Light Mode Background**: Should be `#ffffff` or close
- **Dark Mode Background**: Check `[data-theme="dark"] body` selector
- **Adjust if needed**: Update theme-color values to match actual theme backgrounds

**Example**: If `docs.css` has:
```css
body {
    background: #ffffff;
}

[data-theme="dark"] body {
    background: #1a1a1a;
}
```
Then use those exact values in theme-color meta tags.

**5. Optional: Make Theme Colors Configurable**
- **If needed**: Add to `Config` struct:
  ```go
  ThemeColorLight string // Default: "#ffffff"
  ThemeColorDark  string // Default: "#1a1a1a"
  ```
- **Use in meta generation**:
  ```go
  lightColor := "#ffffff"
  if meta.Config.ThemeColorLight != "" {
      lightColor = meta.Config.ThemeColorLight
  }
  // similar for dark
  ```
- **Note**: Probably not needed initially, can add later if requested

**Testing**:

**Manual Testing**:
1. Generate a site
2. View generated HTML source
3. Verify meta tags present in `<head>`:
   ```html
   <meta name="theme-color" content="#ffffff">
   <meta name="theme-color" content="#1a1a1a" media="(prefers-color-scheme: dark)">
   <meta name="color-scheme" content="light dark">
   ```

**Mobile Browser Testing**:
1. **Android Chrome**:
   - Open site on Android phone
   - Verify address bar is white in light mode
   - Switch system to dark mode
   - Verify address bar is dark gray (`#1a1a1a`)
2. **iOS Safari**:
   - Open site on iPhone
   - Verify status bar/UI chrome matches theme
   - Toggle Appearance (Settings > Display & Brightness)
   - Verify colors update
3. **Samsung Internet**:
   - Test on Samsung device
   - Verify theme color works

**Desktop Browser Testing**:
1. Some desktop browsers use theme-color for UI elements
2. Test in Chrome/Edge with custom themes
3. Verify no negative impact

**Acceptance Criteria**:
- ✅ Meta tags present in all generated HTML pages
- ✅ Light mode theme-color is `#ffffff`
- ✅ Dark mode theme-color is `#1a1a1a`
- ✅ color-scheme meta tag set to `light dark`
- ✅ Mobile browser chrome colors match site theme on:
  - Android Chrome
  - iOS Safari
  - Samsung Internet
- ✅ Colors update when user switches system light/dark mode
- ✅ No visual regressions on desktop browsers

**Documentation Updates**:
- Update CLAUDE.md to mention theme-color support
- Note that mobile browsers will show branded colors

**Estimated Effort**: Trivial (20 lines of code, 1 hour testing on mobile devices)

---

### Story 3: Browser Theme Color Meta Tags

**Feature**: Add `theme-color` and `color-scheme` meta tags for better mobile browser integration

**Requirements**:
- Add `theme-color` meta tags to `RenderMetaTags()`:
  - Light mode: `<meta name="theme-color" content="#ffffff">`
  - Dark mode: `<meta name="theme-color" content="#1a1a1a" media="(prefers-color-scheme: dark)">`
- Add color-scheme meta tag: `<meta name="color-scheme" content="light dark">`
- Match colors to current theme CSS variables
- Consider interaction with accent color feature (use neutral colors by default)

**Acceptance Criteria**:
- Mobile browser chrome colors match site theme
- Works on iOS Safari and Chrome Android
- Colors update properly for light/dark mode

**Estimated Effort**: Trivial (~20 lines in `internal/seo/meta.go`)

---

### Story 4: Instant Navigation with DOM Morphing

**Feature**: Add opt-in instant navigation using hover prefetching + DOM morphing

**Requirements**:
- Add `--instant-nav` CLI flag (disabled by default)
- Embed minified libraries when flag enabled:
  - idiomorph (3.3KB) - DOM morphing
  - instant.page (1KB) - hover prefetching
- Write integration glue code (~100 lines):
  - Hook into instant.page prefetch events
  - Intercept link clicks
  - Fetch new page HTML
  - Use idiomorph to morph DOM (preserve state)
  - Update URL with History API
  - Handle edge cases:
    - External links (skip morphing)
    - Downloads (skip morphing)
    - Anchor links (smooth scroll, no fetch)
    - Hash changes
- Ensure existing JS survives morphing:
  - Theme toggle state preserved
  - Search state preserved
  - Event listeners re-attached if needed
- Add data attributes for control:
  - `data-no-instant` to exclude links from instant nav
- Documentation:
  - Usage guide
  - How to adapt custom JavaScript
  - Performance benefits
  - Edge cases and limitations

**Acceptance Criteria**:
- Navigation feels instant (no white flash)
- Pre-fetches on hover (~200-300ms before click)
- DOM morphs smoothly (only updates changed elements)
- Theme toggle works across navigations
- Search state preserved
- External links open normally
- Downloads work normally
- Anchor links smooth scroll
- Total JS bundle: ~4-5KB added
- No impact when flag not enabled

**Estimated Effort**: Large (~100-150 lines integration code + testing + documentation)

---

### 6. JavaScript Minification

**Concept**: Minify inline JavaScript in generated HTML to reduce page size and improve performance

**Current Situation**:
- Current inline JS in `layout.html`: ~10.9KB unminified
  - Theme detection script: ~415 bytes
  - Main script block: ~10.5KB
- All JavaScript is inline (embedded in HTML template)
- No minification currently applied

**Benefits of Minification**:

**Size Reduction**:
- Typical savings: 70-80% of JavaScript file size
- For Volcano: ~11KB → ~2-3KB (saves ~8KB per page)
- Combined with Brotli/gzip: up to 90% total size reduction
- With instant nav feature: additional ~4-5KB → ~1-1.5KB minified

**Performance Impact**:
- Faster page load times (smaller download)
- Improved Core Web Vitals (LCP, FID, CLS)
- Better mobile performance (critical on slow connections)
- Reduced bandwidth costs for users

**SEO Benefits**:
- Page speed is a Google ranking factor
- Better Core Web Vitals scores improve SEO
- Faster sites = better user experience = better engagement metrics

**Available Tools (Go-based)**:

**Option A: tdewolff/minify** (Pure Go)
- GitHub: `github.com/tdewolff/minify`
- Pure Go library, no external dependencies
- Supports JS, CSS, HTML, SVG, XML, JSON
- Very fast: 8.63 KB output in 3ms (fastest in benchmarks)
- Used by many Go projects
- Active maintenance
- Pros: Pure Go, zero build dependencies, fast, comprehensive
- Cons: Adds Go module dependency (~100KB)

**Option B: esbuild** (Go + external binary)
- Written in Go but needs external esbuild binary
- Extremely fast (10x+ faster than alternatives)
- Best-in-class minification quality
- Modern ES6+ support
- Pros: Fastest, best compression, modern syntax support
- Cons: Requires external binary, more setup complexity

**Option C: Call external tools** (Node/terser)
- Use exec.Command to call terser/uglify
- Best compression ratios
- Pros: Industry-standard tools, proven quality
- Cons: Requires Node.js installed, slow, adds build dependency

**Recommendation**: **Option A (tdewolff/minify)**
- Pure Go, fits Volcano's zero-dependency philosophy (Go module only)
- Very fast (3ms for typical JS)
- Good compression (70-80% reduction)
- No external tools needed
- Simple integration

**Implementation Approach**:

**Option 1: Always-On Minification**
- Minify all JS during template preparation
- Embed minified JS in binary
- Always serve minified JS
- Pros: Best performance for all users, no flags needed
- Cons: Harder to debug generated HTML (but users rarely need to)

**Option 2: CLI Flag** (`--minify-js` or `--production`)
- Only minify when flag is enabled
- Keep readable JS by default
- Pros: Easier debugging during development
- Cons: Users might forget to enable it, two code paths to maintain

**Option 3: Automatic based on output**
- Minify when generating to filesystem
- Don't minify in serve mode (for debugging)
- Pros: Smart default, best of both worlds
- Cons: Inconsistent behavior between modes

**Recommended Approach**: **Option 1 (Always-On)**
- Users rarely inspect generated HTML source
- Browser dev tools show formatted code anyway
- Best performance by default
- Simpler implementation (one code path)
- Can add `--no-minify` flag if debugging needed

**Integration Points**:
1. Add `github.com/tdewolff/minify/v2` as Go module dependency
2. Minify JS in `NewRenderer()` when loading template:
   - Extract `<script>` blocks
   - Minify each block with `js.Minify()`
   - Replace in template
3. For instant-nav feature: minify idiomorph + instant.page before embedding
4. Update CLAUDE.md with new dependency

**Edge Cases & Considerations**:
- Preserve template variables (e.g., `{{.BaseURL}}`) - minify around them
- Test that minified JS works correctly (no syntax errors)
- Consider source maps (probably overkill for inline scripts)
- Error handling: if minification fails, fall back to unminified

**Estimated Savings**:
- Current JS: ~11KB → ~2-3KB minified (**saves ~8KB**)
- With instant-nav: ~15KB → ~4-5KB minified (**saves ~10KB**)
- Per-page savings scales with every page view
- For site with 1000 views/day: saves ~10MB/day bandwidth

**Complexity**: Low-Medium (~50-100 lines Go code + dependency)

**Trade-offs**:
- **Pro**: Significant performance improvement (70-80% JS size reduction)
- **Pro**: Better SEO and Core Web Vitals scores
- **Pro**: Minimal implementation effort with tdewolff/minify
- **Con**: Adds Go module dependency (goes against "zero dependency" slightly)
- **Con**: Minified output harder to read (but users rarely inspect)
- **Con**: Adds ~50ms to build time (negligible)

**Zero-Dependency Philosophy**:
- This is a **build-time** dependency (Go module), not runtime
- No external binaries required (pure Go)
- Still ships as single binary to users
- Alternative: skip minification, accept larger page sizes
- Middle ground: make it opt-in with flag

**Verdict**: **YES** - Always-on minification with tdewolff/minify

**Final Decision**:
- Use tdewolff/minify (pure Go)
- Always minify (no flag, enabled by default)
- Minify at build time (embed minified JS in binary)
- Add as Story 5

---

### Story 5: JavaScript Minification

**Feature**: Minify all inline JavaScript to reduce page size and improve performance

**Requirements**:
- Add `github.com/tdewolff/minify/v2` as Go module dependency
- Minify JavaScript during template loading in `NewRenderer()`:
  - Extract `<script>` blocks from `layout.html`
  - Minify each block with `js.Minify()` from tdewolff/minify
  - Replace original JS with minified version in template
  - Preserve template variables (e.g., `{{.BaseURL}}`) during minification
- Apply to both embedded JS blocks:
  - Theme detection script (~415 bytes)
  - Main script block (~10.5KB)
- When instant-nav feature is added: minify idiomorph + instant.page before embedding
- Error handling: if minification fails, fall back to unminified JS with warning
- Update CLAUDE.md to document new dependency
- Always enabled (no flag needed)
- Optional: Add `--no-minify` flag for debugging if needed later

**Acceptance Criteria**:
- JavaScript size reduced by 70-80% (~11KB → ~2-3KB)
- All JavaScript functionality works correctly (theme toggle, search, navigation, etc.)
- No runtime errors from minification
- Build completes successfully with minified JS
- Template variables preserved and working
- Generated HTML contains minified inline JS

**Testing**:
- Verify theme toggle works
- Verify navigation search works
- Verify copy buttons work
- Verify keyboard shortcuts work
- Verify TOC scroll spy works
- Verify all interactive features function correctly
- Test in multiple browsers (Chrome, Firefox, Safari)

**Estimated Effort**: Small-Medium (~50-100 lines Go code + testing)

**Performance Impact**:
- Saves ~8KB per page load
- Faster page loads (especially on mobile/slow connections)
- Better Core Web Vitals scores
- Improved SEO from page speed

---

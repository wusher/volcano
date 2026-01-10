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

**Requirements**:
- Add `--accent-color="#hexcolor"` CLI flag
- Implement hex → HSL → hex conversion utilities
- Generate 3 CSS variables from input color:
  - `--accent`: Original color (normalized to 50% lightness)
  - `--accent-dark`: 10% lightness (dark mode backgrounds)
  - `--accent-light`: 95% lightness (light mode backgrounds)
- Inject CSS variables into template when flag provided
- Update themes to use accent variables for:
  - Active navigation items
  - TOC active border
  - Link hover states
  - Scroll progress bar
  - Breadcrumbs (optional)
  - Admonitions (optional)
- No variables injected if flag not provided (backward compatible)

**Acceptance Criteria**:
- Works with various brand colors
- Colors look good in both light/dark modes
- No visual change when flag not provided
- Documentation with examples

**Estimated Effort**: Medium (~200-300 lines Go code for color conversion + CSS variable injection + theme updates)

---

### Story 2: Icon-Only Copy Buttons

**Feature**: Simplify code block copy buttons to show only icons (no text)

**Requirements**:
- Remove "Copy" and "Copied!" text from copy buttons
- Keep icon swap functionality (📋 → ✓)
- Add proper `aria-label` attributes for accessibility:
  - `aria-label="Copy code"` (initial state)
  - `aria-label="Copied!"` (success state)
- Update JavaScript to change aria-label on state change

**Acceptance Criteria**:
- Buttons show only icons
- Screen readers announce actions properly
- Visual appearance matches modern patterns (GitHub, VS Code)

**Estimated Effort**: Trivial (~10 lines changed in template + JS)

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

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

### 5. SPA-like Page Transitions (Turbo/PJAX alternatives)

**Concept**: Make static site feel like a SPA by intercepting navigation and doing client-side transitions

**What This Does**:
- Intercepts clicks on links
- Fetches next page via AJAX/fetch
- Swaps content in browser without full page reload
- Updates URL with History API
- Results: Faster navigation, no white flash between pages, can add custom transitions
- Keeps JavaScript state between pages

**Popular Libraries**:

**1. Turbo (from Hotwire)** (~45KB with Stimulus)
- Most popular for this pattern
- Intercepts all links/forms automatically
- Battle-tested (evolved from Turbolinks/PJAX)
- Opinionated but smooth out-of-the-box experience
- Pros: Comprehensive, well-documented, large community
- Cons: Larger bundle, more opinionated, ties into Hotwire ecosystem

**2. htmx** (~14KB)
- Most popular alternative (90K+ weekly npm downloads, 46K+ GitHub stars)
- Attribute-based: `<a href="/page" hx-boost="true">`
- More flexible and transparent than Turbo
- Can do partial updates, not just full page swaps
- Pros: Smaller, more flexible, broader feature set (AJAX, WebSockets, SSE)
- Cons: More manual setup, different paradigm

**3. Unpoly** (Similar size to Turbo)
- Like htmx but more Rails-like (opinionated, convention-based)
- Has migration scripts for upgrading
- Comprehensive like Turbo but more flexible
- Pros: Good balance of opinions and flexibility
- Cons: Smaller community than Turbo/htmx

**4. Barba.js** (~9KB)
- Focused specifically on page transitions with animations
- Lightweight, flexible for custom transitions
- Great for creative/portfolio sites
- Pros: Small, focused, excellent animation support
- Cons: Less comprehensive (just transitions, not forms/etc)

**5. Swup** (Small)
- Similar to Barba.js, focused on smooth animated transitions
- Easy setup
- Pros: Simple, works well with various frameworks
- Cons: Less feature-rich than Turbo/htmx

**6. InstantClick** (Tiny)
- Preloads pages on hover, then swaps instantly on click
- Extremely fast perceived performance
- Pros: Tiny, clever approach, very fast
- Cons: Limited to link navigation, no forms

**Recommendation for Volcano**:
- **htmx**: Best balance of size (~14KB), flexibility, and features. Attribute-based fits well with HTML-first philosophy
- **Turbo**: If you want comprehensive, opinionated, battle-tested solution and don't mind the size
- **Barba.js/Swup**: If you mainly want smooth transitions and don't need form handling
- **InstantClick**: If you want the absolute smallest/simplest solution

**Implementation Considerations**:
- Would need to inject library into generated HTML
- Could embed minified JS directly (like CSS) or link to CDN
- May need special handling for theme toggle, search, and other JS state
- Breaking change: users would need to opt-in via flag (e.g., `--spa-navigation`)

**Zero-Dependency Trade-off**:
- Adds runtime dependency (but not build dependency)
- Could embed JS directly to keep "single binary" philosophy
- This is a significant departure from current minimal JS approach

**Pros**:
- Much faster perceived performance
- Modern, app-like feel
- Better user experience for navigation-heavy sites

**Cons**:
- Adds JavaScript dependency (~14-45KB)
- Can break assumptions (scripts run differently)
- More complexity in testing
- Need to handle edge cases (external links, downloads, etc.)

**Verdict**: TBD - is this valuable enough to add runtime dependency?

---

## Stories

<!-- Stories will be formalized here once we've discussed the ideas -->

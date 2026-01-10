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

## Stories

<!-- Stories will be formalized here once we've discussed the ideas -->

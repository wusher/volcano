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

**Verdict**: TBD - HSL computation vs simple approach

---

## Stories

<!-- Stories will be formalized here once we've discussed the ideas -->

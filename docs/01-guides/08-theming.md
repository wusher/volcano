# Theming

Customize your site's appearance with built-in themes or custom CSS.

## Built-in Themes

Volcano includes four themes:

### docs (Default)

A full-featured documentation theme:

```bash
volcano ./docs --theme docs
```

**Features:**
- Sidebar tree navigation
- Table of contents (auto, when a page has 3+ headings)
- Dark mode toggle
- Typography tuned for long-form reading
- Responsive layout

Optional add-ons like search, breadcrumbs, top nav, and page nav are flag-controlled — see [[customizing-appearance]].

**Best for:** Technical documentation, API references, user guides.

### blog

A reading-focused theme:

```bash
volcano ./docs --theme blog
```

**Features:**
- Content-centered layout
- Optimized typography for reading
- Clean, minimal design
- Dark mode support

**Best for:** Blogs, articles, announcements, changelogs.

### presentation

A slide-deck inspired theme for talks and demos:

```bash
volcano ./talk --theme presentation
```

**Features:**
- Oversized, fluid display typography (H1 reads like a title slide)
- Section dividers above every H2 — pages feel like slide decks
- Pull-quote style blockquotes and banner admonitions
- High-contrast palette tuned for projector legibility
- Dark mode optimized for stage lighting

**Best for:** Conference talks, product demos, internal showcases, narrative writeups.

### vanilla

A structural skeleton with no visual styling:

```bash
volcano ./docs --theme vanilla
```

**Features:**
- All layout structure
- No colors, fonts, or decorations
- Extensively commented CSS
- Perfect starting point for custom themes

**Best for:** Custom designs, unique branding, learning the CSS structure.

## Selecting a Theme

Use the `--theme` flag:

```bash
volcano ./docs --theme blog
```

## Accent Color

Every built-in theme honors an accent color. Set it with `--accent-color` using
a Tailwind color name (uses the `500` shade), a hex value, or a two-color
gradient. The color is exposed to CSS as the `--accent` custom property. The
default is `sky`:

```bash
# Tailwind names (case-insensitive)
volcano ./docs --accent-color sky
volcano ./docs --accent-color rose
volcano ./docs --accent-color emerald

# Hex values
volcano ./docs --accent-color "#0ea5e9"

# Two-color gradient (each side can be name or hex)
volcano ./docs --accent-color lime-sky
volcano ./docs --accent-color "#444444-#555555"
volcano ./docs --accent-color "lime-#0ea5e9"
```

Supported names: `slate`, `gray`, `zinc`, `neutral`, `stone`, `red`,
`orange`, `amber`, `yellow`, `lime`, `green`, `emerald`, `teal`, `cyan`,
`sky`, `blue`, `indigo`, `violet`, `purple`, `fuchsia`, `pink`, `rose`.

When a gradient is set, three custom properties become available alongside
`--accent`:

- `--accent-end` — the second color (single value)
- `--accent-gradient` — `linear-gradient(to right, start, end)`, used for wide horizontal accents (page H1, links, page-nav, scroll progress bar)
- `--accent-gradient-vertical` — `linear-gradient(to bottom, start, end)`, used for narrow vertical accents (admonition + blockquote left borders)

Two directions ship because a single 135° diagonal looked first-color-heavy on
wide headings and didn't read as a clean top-to-bottom sweep on tall vertical
bars. Themes and custom CSS can reference either variable wherever a gradient
is a valid value (backgrounds, `border-image`, etc.).

## Extracting CSS

Export the vanilla CSS to create your own theme:

```bash
volcano css -o my-theme.css
```

This creates a file with all structural CSS and detailed comments:

```css
/* =============================================================================
   SIDEBAR (Left Navigation Panel)
   =============================================================================
   Fixed panel containing site title, search, and tree navigation.

   Structure:
   .sidebar
   ├── .sidebar-header (site title + close button for mobile)
   ├── .nav-search (search input)
   └── .tree-nav (folder/file navigation tree)

   Styling suggestions:
   - Add background-color
   - Add border-right for separation from content
   ============================================================================= */

.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: var(--sidebar-width);
  height: 100vh;
  overflow-y: auto;
  z-index: 100;
  display: flex;
  flex-direction: column;
}
```

## Using Custom CSS

Apply your custom CSS with `--css`:

```bash
volcano ./docs --css ./my-theme.css
```

When using `--css`:
- The `--theme` flag is ignored
- Your CSS completely replaces the theme
- You have full control over styling

## Creating a Custom Theme

### 1. Export the Skeleton

```bash
volcano css -o my-theme.css
```

### 2. Add CSS Variables

Define your color scheme with custom properties:

```css
:root {
  /* Layout */
  --sidebar-width: 280px;
  --toc-width: 220px;
  --content-max-width: 800px;

  /* Colors */
  --color-bg: #ffffff;
  --color-text: #1a1a1a;
  --color-text-muted: #666666;
  --color-link: #2563eb;
  --color-border: #e5e5e5;

  /* Typography */
  --font-sans: system-ui, -apple-system, sans-serif;
  --font-mono: "SF Mono", Consolas, monospace;
}
```

### 3. Style Core Elements

```css
body {
  font-family: var(--font-sans);
  background-color: var(--color-bg);
  color: var(--color-text);
  line-height: 1.6;
}

.sidebar {
  background-color: #f8f9fa;
  border-right: 1px solid var(--color-border);
}

.prose a {
  color: var(--color-link);
  text-decoration: none;
}

.prose a:hover {
  text-decoration: underline;
}
```

### 4. Add Dark Mode

```css
[data-theme="dark"] {
  --color-bg: #1a1a1a;
  --color-text: #e5e5e5;
  --color-text-muted: #999999;
  --color-link: #60a5fa;
  --color-border: #333333;
}

[data-theme="dark"] .sidebar {
  background-color: #0d0d0d;
}
```

### 5. Use Your Theme

```bash
volcano ./docs --css ./my-theme.css
```

## CSS Architecture

### Layout Classes

| Class | Element |
|-------|---------|
| `.sidebar` | Left navigation panel |
| `.main-wrapper` | Main content container |
| `.content` | Inner content area |
| `.prose` | Article content |
| `.toc-sidebar` | Table of contents panel |

### Navigation Classes

| Class | Element |
|-------|---------|
| `.sidebar-header` | Site title area |
| `.nav-search` | Search input container |
| `.tree-nav` | Navigation tree |
| `.folder-header` | Folder row |
| `.folder-toggle` | Expand/collapse button |
| `.breadcrumbs` | Breadcrumb navigation |
| `.page-nav` | Previous/next links |
| `.top-nav` | Top navigation bar |

### Content Classes

| Class | Element |
|-------|---------|
| `.prose h1` - `.prose h6` | Headings |
| `.prose p` | Paragraphs |
| `.prose a` | Links |
| `.prose code` | Inline code |
| `.prose pre` | Code blocks |
| `.prose blockquote` | Blockquotes |
| `.prose table` | Tables |

### Component Classes

| Class | Element |
|-------|---------|
| `.admonition` | Callout boxes |
| `.admonition-note` | Note type |
| `.admonition-tip` | Tip type |
| `.admonition-warning` | Warning type |
| `.admonition-danger` | Danger type |
| `.code-block` | Code block wrapper |
| `.copy-button` | Copy code button |
| `.theme-toggle` | Dark mode button |
| `.back-to-top` | Scroll to top button |

## Dark Mode

### How It Works

Dark mode is controlled by the `data-theme` attribute on `<html>`:

```html
<html data-theme="light">  <!-- Light mode -->
<html data-theme="dark">   <!-- Dark mode -->
```

The theme toggle button switches between these and saves the preference to localStorage.

### Implementing Dark Mode

Define styles for both modes:

```css
/* Light mode (default) */
:root {
  --color-bg: #ffffff;
  --color-text: #1a1a1a;
}

/* Dark mode */
[data-theme="dark"] {
  --color-bg: #1a1a1a;
  --color-text: #e5e5e5;
}

/* Use variables */
body {
  background: var(--color-bg);
  color: var(--color-text);
}
```

## Responsive Design

Built-in themes include responsive breakpoints:

```css
/* Tablet - Hide TOC */
@media (max-width: 1280px) {
  .toc-sidebar {
    display: none;
  }
}

/* Mobile - Collapsible sidebar */
@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
  }

  body.drawer-open .sidebar {
    transform: translateX(0);
  }
}
```

## Syntax Highlighting

Code blocks use Chroma CSS classes:

```css
.chroma .k  { color: #d73a49; }  /* Keywords */
.chroma .s  { color: #032f62; }  /* Strings */
.chroma .c  { color: #6a737d; }  /* Comments */
.chroma .nf { color: #6f42c1; }  /* Functions */
.chroma .m  { color: #005cc5; }  /* Numbers */
```

See the vanilla CSS for the complete list of token classes.

## Examples

### Minimal Light Theme

A serif body, soft sidebar, classic blue links — about a dozen lines of CSS.

![Minimal Light custom theme](/images/custom-themes/minimal-light.png)

```css
:root {
  --sidebar-width: 260px;
  --content-max-width: 720px;
}

body {
  font-family: Georgia, serif;
  background: #fefefe;
  color: #333;
}

.sidebar {
  background: #f5f5f5;
  border-right: 1px solid #ddd;
}

.prose a {
  color: #0066cc;
}
```

### Dark Neon Theme

Deep black background, gradient sidebar, cyan accents, and a magenta-to-cyan gradient on the page title.

![Dark Neon custom theme](/images/custom-themes/dark-neon.png)

```css
:root {
  --sidebar-width: 280px;
}

body {
  font-family: 'Inter', sans-serif;
  background: #0a0a0f;
  color: #e0e0e0;
}

.sidebar {
  background: linear-gradient(180deg, #1a1a2e 0%, #0a0a0f 100%);
  border-right: 1px solid #2a2a4a;
}

.prose a {
  color: #00ffff;
}

.prose h1 {
  background: linear-gradient(90deg, #ff00ff, #00ffff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
```

## Print Styles

Built-in themes include print-optimized stylesheets that:

- Hide navigation elements (sidebar, TOC, breadcrumbs)
- Remove interactive elements (theme toggle, back-to-top button)
- Optimize typography for paper
- Ensure proper page breaks

To customize print styles:

```css
@media print {
  /* Hide elements you don't want printed */
  .sidebar,
  .toc-sidebar,
  .page-nav {
    display: none !important;
  }

  /* Adjust content width */
  .content {
    max-width: 100%;
  }

  /* Better link visibility */
  .prose a::after {
    content: " (" attr(href) ")";
    font-size: 0.8em;
  }
}
```

## Related

- [[customizing-appearance]] — Quick customization guide
- [[building-your-site]] — Build with your theme

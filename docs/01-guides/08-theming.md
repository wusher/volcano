# Theming

Customize your site's appearance with built-in themes or custom CSS.

## Built-in Themes

Volcano includes three themes:

### docs (Default)

A full-featured documentation theme:

```bash
volcano ./docs --theme docs
```

**Features:**
- Sidebar navigation with search
- Breadcrumb navigation
- Table of contents
- Dark mode toggle
- Professional typography
- Responsive design

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

# Custom CSS

When a theme + accent color aren't enough.

## Export the Skeleton

```bash
volcano css -o my-theme.css
```

This dumps the `vanilla` theme — every layout class, fully commented, no decoration. It's the cleanest starting point.

## Use Your CSS

```bash
volcano ./docs --css ./my-theme.css --url="https://example.com"
```

`--css` replaces the theme entirely. `--theme` is ignored when `--css` is set.

## What to Customize

Override CSS custom properties at the top of your file — most theming happens through these:

```css
:root {
  /* Layout */
  --sidebar-width: 280px;
  --toc-width: 220px;
  --content-max-width: 800px;

  /* Colors */
  --bg-primary: #ffffff;
  --bg-secondary: #f8f9fa;
  --text-primary: #1a1a1a;
  --text-muted: #666666;
  --border-color: #e5e5e5;

  /* Typography */
  --font-sans: system-ui, -apple-system, sans-serif;
  --font-mono: "SF Mono", Consolas, monospace;
}

[data-theme="dark"] {
  --bg-primary: #1a1a1a;
  --text-primary: #e5e5e5;
  --border-color: #333333;
}
```

The theme toggle button sets `data-theme="light"` or `data-theme="dark"` on `<html>` and persists the choice in `localStorage`.

## Key Classes

| Class | What it is |
|-------|------------|
| `.sidebar`, `.tree-nav` | Left navigation panel + tree |
| `.main-wrapper`, `.content`, `.prose` | Page chrome → article wrapper → article content |
| `.toc-sidebar` | Right-side table of contents |
| `.breadcrumbs`, `.top-nav`, `.page-nav` | Optional nav elements |
| `.command-palette` | Search modal (Cmd+K) |
| `.admonition`, `.admonition-tip`, `.admonition-warning`, `.admonition-danger`, `.admonition-note` | Callout boxes |
| `.code-block`, `.copy-button` | Code blocks + copy button |
| `.scroll-progress`, `.back-to-top`, `.theme-toggle` | UI affordances |

Open the exported `my-theme.css` for the full list with structural comments.

## Two Worked Examples

### Minimal serif

Soft sidebar, classic blue links, Georgia body.

![Minimal Light theme](/images/custom-themes/minimal-light.png)

```css
:root {
  --content-max-width: 720px;
}
body {
  font-family: Georgia, serif;
  background: #fefefe;
  color: #333;
}
.sidebar { background: #f5f5f5; border-right: 1px solid #ddd; }
.prose a { color: #0066cc; }
```

### Dark neon

Black background, cyan accents, gradient H1.

![Dark Neon theme](/images/custom-themes/dark-neon.png)

```css
body {
  font-family: 'Inter', sans-serif;
  background: #0a0a0f;
  color: #e0e0e0;
}
.sidebar {
  background: linear-gradient(180deg, #1a1a2e, #0a0a0f);
  border-right: 1px solid #2a2a4a;
}
.prose a { color: #00ffff; }
.prose h1 {
  background: linear-gradient(90deg, #ff00ff, #00ffff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
```

## AI-Assisted Themes

The vanilla CSS export is short, well-commented, and self-contained — paste it into Claude or ChatGPT, describe the look you want, get a usable theme back. See [[advanced/ai-prompts|AI prompts]] for a tested starter prompt.

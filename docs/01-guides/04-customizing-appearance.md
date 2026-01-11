# Customizing Appearance

Change your site's look with themes, custom CSS, and branding options.

## Choosing a Theme

Volcano includes three built-in themes:

### docs (Default)

A full-featured documentation theme with sidebar navigation, breadcrumbs, and table of contents.

```bash
volcano ./docs --theme docs
```

Best for: Technical documentation, API references, user guides.

### blog

A blog-optimized layout with emphasis on content readability.

```bash
volcano ./docs --theme blog
```

Best for: Blogs, articles, news sites, changelogs.

### vanilla

A minimal structural skeleton with no colors or visual styling. Use this as a starting point for completely custom designs.

```bash
volcano ./docs --theme vanilla
```

Best for: Custom designs, sites requiring unique branding.

## Extracting CSS for Customization

Export the vanilla CSS skeleton to customize:

```bash
volcano css -o custom.css
```

This creates a CSS file with all structural styles and extensive comments explaining each section:

```css
/* =============================================================================
   SIDEBAR (Left Navigation Panel)
   =============================================================================
   Fixed panel containing site title, search, and tree navigation.
   ... */

.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: var(--sidebar-width);
  /* Add your styles here */
}
```

## Using Custom CSS

Apply your custom CSS with the `--css` flag:

```bash
volcano ./docs --css ./custom.css --title="My Site"
```

When using `--css`, the theme flag is ignored. Your CSS file completely replaces the built-in theme.

## Creating a Custom Theme

1. **Export the skeleton:**
   ```bash
   volcano css -o my-theme.css
   ```

2. **Add your colors and styles:**
   ```css
   :root {
     --sidebar-width: 280px;
     --content-max-width: 800px;

     /* Add your colors */
     --color-bg: #ffffff;
     --color-text: #1a1a1a;
     --color-link: #2563eb;
     --color-border: #e5e5e5;
   }

   body {
     background-color: var(--color-bg);
     color: var(--color-text);
     font-family: system-ui, sans-serif;
   }

   .sidebar {
     background-color: #f8f8f8;
     border-right: 1px solid var(--color-border);
   }
   ```

3. **Use your theme:**
   ```bash
   volcano ./docs --css ./my-theme.css
   ```

## Adding a Favicon

Add a favicon to your site:

```bash
volcano ./docs --favicon ./favicon.ico
```

Supported formats:
- `.ico` — Traditional favicon format
- `.png` — PNG image
- `.svg` — SVG vector image
- `.gif` — Animated GIF

The favicon file is copied to your output directory and proper `<link>` tags are added to all pages.

## Setting Site Metadata

### Site Title

The site title appears in the sidebar header and page titles:

```bash
volcano ./docs --title="My Project Documentation"
```

Page titles are formatted as: `Page Title - Site Title`

### Author

Set the author for meta tags:

```bash
volcano ./docs --author="Jane Smith"
```

This adds `<meta name="author" content="Jane Smith">` to all pages.

### Base URL

Set the base URL for canonical links and SEO:

```bash
volcano ./docs --url="https://docs.example.com"
```

This enables:
- Canonical URL meta tags
- Absolute URLs in Open Graph tags

## Display Options

### Last Modified Date

Show when each page was last modified:

```bash
volcano ./docs --last-modified
```

Volcano checks Git history first, falling back to file modification time.

### Top Navigation Bar

Display root-level items in a top navigation bar:

```bash
volcano ./docs --top-nav
```

The top nav shows files and folders at the root level (1-8 items maximum).

### Page Navigation

Show previous/next page links at the bottom of each page:

```bash
volcano ./docs --page-nav
```

## SEO Options

### Open Graph Image

Set a default Open Graph image for social sharing:

```bash
volcano ./docs --og-image="https://example.com/og-image.png"
```

This image appears when your pages are shared on social media.

## Combining Options

Options can be combined:

```bash
volcano ./docs \
  -o ./public \
  --title="My Project" \
  --url="https://docs.myproject.com" \
  --author="My Team" \
  --favicon="./assets/favicon.png" \
  --theme docs \
  --top-nav \
  --page-nav \
  --last-modified
```

## Dark Mode

The built-in themes support dark mode through the `data-theme` attribute:

```html
<html data-theme="dark">
```

The theme toggle button is included in the generated pages and persists the user's preference.

When creating custom CSS, define dark mode styles:

```css
[data-theme="dark"] {
  --color-bg: #1a1a1a;
  --color-text: #e5e5e5;
  --color-link: #60a5fa;
}

[data-theme="dark"] body {
  background-color: var(--color-bg);
  color: var(--color-text);
}
```

## CSS Class Reference

Key classes you can style:

| Class | Element |
|-------|---------|
| `.sidebar` | Left navigation panel |
| `.sidebar-header` | Site title area |
| `.tree-nav` | Navigation tree |
| `.folder-header` | Folder toggle + name |
| `.main-wrapper` | Main content container |
| `.content` | Content area |
| `.prose` | Article content |
| `.breadcrumbs` | Breadcrumb navigation |
| `.toc-sidebar` | Table of contents |
| `.page-nav` | Previous/next links |
| `.admonition` | Callout boxes |
| `.code-block` | Code block wrapper |

See the exported vanilla CSS for the complete list with documentation.

## Next Steps

- [[development-workflow]] — Preview changes with the dev server
- [[features/theming]] — Deep dive into theming
- [[deploying-your-site]] — Publish your customized site

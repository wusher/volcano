# Navigation

Volcano automatically generates comprehensive navigation from your folder structure.

## Sidebar Navigation Tree

The sidebar displays your content as a hierarchical tree:

```
My Site
├── Getting Started
├── Guides
│   ├── Installation
│   ├── Configuration
│   └── Deployment
└── Reference
    ├── CLI
    └── API
```

### Features

- **Expand/collapse folders** — Click the chevron or folder name
- **Active page highlighting** — Current page is visually distinct
- **Auto-expand to current page** — Navigating to a nested page expands its parents
- **Keyboard accessible** — Full ARIA tree role support

### Customizing

The sidebar reflects your folder structure. To change navigation:
1. Rename or move files and folders
2. Use [[guides/organizing-content|number prefixes]] to control order
3. Create `index.md` files for folder landing pages

## Search

Enable site search with `--search`:

```bash
volcano ./docs --search
```

This adds a command palette (Cmd+K / Ctrl+K) that searches page titles and content.

## Breadcrumbs

Breadcrumbs show the current page's location in the hierarchy:

```
My Site > Guides > Installation
```

Breadcrumbs are **enabled by default**. To disable them:

```bash
volcano ./docs --breadcrumbs=false
```

### Features

- Links to each ancestor level
- Current page shown without link
- Schema.org BreadcrumbList markup for SEO
- Hidden on homepage (single item)

### Structure

For a page at `/guides/advanced/topics/`:

```
Home > Guides > Advanced > Topics
```

Each item (except the last) is a clickable link.

## Table of Contents

Pages with 3 or more headings get a table of contents sidebar:

```
On this page
├── Introduction
├── Installation
│   ├── Requirements
│   └── Steps
├── Configuration
└── Troubleshooting
```

### Features

- Extracts h2, h3, and h4 headings
- Nested structure matches heading hierarchy
- Click to scroll to section
- Active heading highlighted on scroll

### Behavior

- **Minimum headings**: TOC only appears if 3+ headings exist
- **Desktop only**: Hidden on mobile (viewport < 1280px)
- **Scroll tracking**: Active heading updates as you scroll

### Heading Requirements

For the TOC to work well:
- Use h2 (`##`) for main sections
- Use h3 (`###`) for subsections
- Keep heading text concise
- Use unique heading text (duplicates get numbered IDs)

## Page Navigation

Enable previous/next links at the bottom of each page:

```bash
volcano ./docs --page-nav
```

Shows:

```
← Previous: Installation    Next: Configuration →
```

### Ordering

Pages are ordered by:
1. Position in sidebar navigation
2. Sequential through the tree (depth-first)

The first page has no "Previous", the last has no "Next".

## Top Navigation Bar

Display root-level items in a horizontal bar:

```bash
volcano ./docs --top-nav
```

### Features

- Shows files and folders at root level
- Limited to 1-8 items (for usability)
- Includes theme toggle button
- Responsive on mobile

### Item Order

Items are sorted:
1. Files first, then folders
2. By date prefix (newest first)
3. By number prefix
4. Alphabetically

### When to Use

Top nav works well when:
- You have a small number of main sections
- Root-level organization is meaningful
- You want quick access to key pages

## Keyboard Shortcuts

Press `?` to see all shortcuts.

| Key | Action |
|-----|--------|
| `Cmd+K` / `Ctrl+K` | Open search (requires `--search`) |
| `t` | Toggle theme |
| `z` | Toggle zen mode (hides sidebar) |
| `n` | Next page (requires `--page-nav`) |
| `p` | Previous page (requires `--page-nav`) |
| `h` | Go to homepage |
| `?` | Show shortcuts |
| `Esc` | Close modal |

## Instant Navigation

Enable instant navigation with hover prefetching for near-instant page loads:

```bash
volcano ./docs --instant-nav
```

### How It Works

When enabled, Volcano prefetches pages when you hover over links. By the time you click, the page is often already loaded, making navigation feel instantaneous.

### Features

- Prefetches pages on hover with a 65ms delay
- Smart throttling prevents excessive requests
- Works with all internal links (sidebar, breadcrumbs, content)
- Respects browser cache
- No configuration needed

### Performance

Instant navigation uses minimal bandwidth because:
- Only prefetches when hovering (not all links at once)
- Respects the browser's cache
- Throttles multiple hovers on the same link

## Scroll Progress Indicator

A thin progress bar at the top of the page shows how far you've scrolled through the current page. This helps readers track their position in longer documents.

The progress bar:
- Appears at the very top of the viewport
- Fills from left to right as you scroll
- Updates smoothly in real-time

## Back to Top Button

A floating button appears in the bottom-right corner after scrolling down 300 pixels. Click it to smoothly scroll back to the top of the page.

The button:
- Hidden when near the top
- Fades in when scrolling down
- Uses smooth scrolling animation
- Styled to match the current theme

## Navigation Elements Summary

| Element | Location | Purpose |
|---------|----------|---------|
| Sidebar tree | Left side | Browse all pages |
| Search | Top of sidebar | Filter pages |
| Breadcrumbs | Above content | Show current location |
| TOC | Right side | Jump within page |
| Page nav | Bottom of content | Sequential browsing |
| Top nav | Top bar | Quick section access |
| Progress bar | Top of viewport | Show scroll position |
| Back to top | Bottom-right | Quick return to top |

## Accessibility

Volcano's navigation is built with accessibility in mind:

- **ARIA roles**: Tree navigation uses proper `tree` and `treeitem` roles
- **Keyboard support**: Navigate with Tab, Enter, Arrow keys
- **Screen reader labels**: Buttons have descriptive labels
- **Focus indicators**: Visible focus states throughout
- **Skip links**: Hidden link to skip to main content

## Related

- [[guides/organizing-content]] — Structure your content
- [[reference/url-structure]] — How URLs are generated

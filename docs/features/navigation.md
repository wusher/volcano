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

The sidebar includes a search box that filters the navigation tree:

- Type to filter pages by name
- Matching pages remain visible
- Non-matching pages are hidden
- Clear the search to show all pages

Search matches against page titles, not content.

## Breadcrumbs

Breadcrumbs show the current page's location in the hierarchy:

```
My Site > Guides > Installation
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

## Navigation Elements Summary

| Element | Location | Purpose |
|---------|----------|---------|
| Sidebar tree | Left side | Browse all pages |
| Search | Top of sidebar | Filter pages |
| Breadcrumbs | Above content | Show current location |
| TOC | Right side | Jump within page |
| Page nav | Bottom of content | Sequential browsing |
| Top nav | Top bar | Quick section access |

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

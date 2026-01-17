# Top Navigation

Display root-level items in a horizontal navigation bar at the top.

## Enabling Top Navigation

```bash
volcano ./docs --top-nav
```

## What You Get

A horizontal bar at the top showing root-level pages and folders:

```
[Home]  [Guides]  [Reference]  [About]  🌙
```

The theme toggle (🌙) appears on the right.

## What Gets Shown

Only items at the root level of your content directory:

```
docs/
├── index.md           → "Home" in top nav
├── guides/            → "Guides" in top nav
├── reference/         → "Reference" in top nav
└── about.md          → "About" in top nav
└── nested/
    └── deep.md       → NOT in top nav (too deep)
```

## Item Limit

Top nav displays 1-8 items for usability. If you have more than 8 root items, only the first 8 appear.

## Ordering

Items are sorted:
1. **Files first**, then folders
2. By **date prefix** (newest first)
3. By **number prefix**
4. **Alphabetically**

Example:

```
docs/
├── 01-home.md
├── 02-guides/
├── 03-reference/
└── about.md
```

Shows: Home | Guides | Reference | About

## Active State

The current page's top-level section is highlighted:

```
Home  [Guides]  Reference  About
      ^^^^^^^^ highlighted when on any /guides/* page
```

## Mobile Behavior

On mobile (< 768px width):
- Top nav becomes a collapsible menu
- Hamburger icon opens the full sidebar
- Theme toggle remains accessible

## Styling

Top navigation:
- Matches your theme colors
- Includes hover effects
- Shows active state
- Responsive layout

## When to Use

Top nav works well when:
- You have a small number of main sections (2-8)
- Root-level organization is intuitive
- You want quick access to major areas

Example sites:
- **Documentation**: Docs | API | Guides | Blog
- **Product site**: Features | Pricing | Docs | Support
- **Blog**: Posts | About | Contact

## When to Skip

Consider skipping top nav when:
- You have many root items (>8)
- Deep hierarchy is more important
- Sidebar navigation is sufficient

## Combining with Sidebar

Top nav and sidebar work together:

- **Top nav**: Quick access to main sections
- **Sidebar**: Detailed navigation within sections

```bash
volcano ./docs \
  --top-nav \
  --search \
  --breadcrumbs
```

This gives users multiple navigation options.

## Accessibility

- Semantic HTML (`<nav>`, `<ul>`, `<li>`)
- `aria-label="Main navigation"`
- `aria-current="page"` on active items
- Keyboard navigable (Tab, Enter)

## Layout Impact

With top nav enabled:
- Content starts below the nav bar
- Sidebar begins under the nav bar
- Vertical space is slightly reduced

## Theme Toggle Location

The theme toggle button (light/dark mode) appears:
- **Without top nav**: In the sidebar or floating toolbar
- **With top nav**: On the right side of the top bar

## Example

For a documentation site with clear sections:

```bash
volcano ./docs \
  --top-nav \
  --theme docs \
  --title="My Product Docs"
```

Structure:

```
docs/
├── index.md           # Overview
├── getting-started/   # Guides
├── api-reference/     # API docs
└── examples/          # Code examples
```

Top nav shows: Overview | Getting Started | API Reference | Examples

## Related

- [[navigation]] — Navigation overview
- [[breadcrumbs]] — Hierarchical breadcrumbs
- [[search]] — Search all pages

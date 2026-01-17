# Page Navigation

Add previous/next links at the bottom of each page for sequential reading.

## Enabling Page Navigation

```bash
volcano ./docs --page-nav
```

## What You Get

Links at the bottom of each page:

```
← Previous: Installation    Next: Configuration →
```

## How It Works

Pages are ordered by:
1. Their position in the sidebar navigation tree
2. Sequential traversal (depth-first)

The system automatically:
- Shows previous link (except on first page)
- Shows next link (except on last page)
- Uses actual page titles

## Navigation Order

Given this structure:

```
docs/
├── index.md              # 1. First page (no "Previous")
├── getting-started.md    # 2.
├── guides/
│   ├── index.md         # 3.
│   ├── basic.md         # 4.
│   └── advanced.md      # 5.
└── reference.md         # 6. Last page (no "Next")
```

Page navigation follows this order: 1 → 2 → 3 → 4 → 5 → 6

## Controlling Order

Use number prefixes to control order:

```
docs/
├── 01-introduction.md
├── 02-installation.md
├── 03-configuration.md
└── 04-deployment.md
```

Navigation will follow: Introduction → Installation → Configuration → Deployment

## Keyboard Shortcuts

When page navigation is enabled:

| Key | Action |
|-----|--------|
| `n` | Go to next page |
| `p` | Go to previous page |

## Styling

Page navigation:
- Appears at the bottom of content
- Uses your theme's colors
- Shows arrows for direction (← →)
- Is responsive (stacks on mobile)

## Best For

Page navigation works well for:

- **Tutorials** - Step-by-step guides
- **Books** - Linear reading order
- **Courses** - Sequential lessons
- **Documentation** - Progressive learning

## Not Ideal For

Skip page navigation when:
- Content is reference material (no reading order)
- Site structure is non-linear
- Users typically search for specific topics

## Accessibility

- Proper semantic HTML (`<nav>`)
- `aria-label="Page navigation"`
- Clear link text (not just "Next"/"Previous")
- Keyboard accessible

## Combining with Other Navigation

Works well with:

```bash
volcano ./docs \
  --page-nav \
  --breadcrumbs \
  --search
```

- **Breadcrumbs** - Show where you are
- **Page nav** - Move sequentially
- **Search** - Jump to any page

## Example

For a tutorial site:

```bash
volcano ./tutorials \
  --page-nav \
  --theme docs \
  --title="My Tutorial"
```

Readers can click "Next" to follow the tutorial in order, or use breadcrumbs/search to jump around.

## Related

- [[navigation]] — Navigation overview
- [[breadcrumbs]] — Hierarchical navigation
- [[keyboard-shortcuts]] — All keyboard shortcuts

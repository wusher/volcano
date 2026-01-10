# Auto-Generated Index Pages

Volcano automatically creates index pages for folders that don't have one.

## How It Works

When you have a folder without an `index.md` or `readme.md` file, Volcano generates an index page automatically. This ensures every folder in your site is accessible.

### Before

```
docs/
├── guides/
│   ├── installation.md
│   └── configuration.md
└── index.md
```

The `guides/` folder has no index file.

### After Generation

Visiting `/guides/` shows an auto-generated page listing:

- Configuration
- Installation

## Generated Content

Auto-index pages include:

1. **Heading** — The folder name as the page title
2. **Link list** — All files and subfolders in the folder

### Example Output

For a folder structure like:

```
reference/
├── cli.md
├── api.md
└── advanced/
    └── plugins.md
```

The auto-generated `/reference/` page shows:

```html
<article class="auto-index-page">
  <h1>Reference</h1>
  <ul class="folder-index">
    <li class="page-item"><a href="/reference/api/">API</a></li>
    <li class="page-item"><a href="/reference/cli/">CLI</a></li>
    <li class="folder-item"><a href="/reference/advanced/">Advanced</a></li>
  </ul>
</article>
```

## Sort Order

Items are sorted:

1. **Pages first** — Markdown files appear before folders
2. **Alphabetically** — Within each group, sorted A-Z (case-insensitive)

This matches the sidebar navigation order.

## When Auto-Index is NOT Generated

Auto-index pages are skipped when:

1. **index.md exists** — Folder has an `index.md` file
2. **readme.md exists** — Folder has a `readme.md` file (case-insensitive)
3. **Root folder** — The input root directory (use index.md for homepage)

### Example

```
docs/
├── guides/
│   ├── index.md        ← Uses this file
│   └── setup.md
├── reference/
│   ├── readme.md       ← Uses this file
│   └── cli.md
└── examples/           ← Auto-index generated
    ├── basic.md
    └── advanced.md
```

## Styling Auto-Index Pages

Auto-index pages use specific CSS classes:

```css
/* The entire auto-index page */
.auto-index-page { }

/* The list of links */
.folder-index { }

/* Individual page links */
.folder-index .page-item { }

/* Individual folder links */
.folder-index .folder-item { }

/* Empty folder message */
.empty-folder { }
```

### Example Custom Styling

```css
.folder-index {
  list-style: none;
  padding: 0;
}

.folder-index li {
  padding: 8px 12px;
  border-bottom: 1px solid var(--color-border);
}

.folder-index .folder-item::before {
  content: "📁 ";
}

.folder-index .page-item::before {
  content: "📄 ";
}
```

## Best Practices

### Write Custom Index Pages

Auto-index is a fallback. For better user experience, write custom index pages:

```markdown
# Guides

Learn how to use Volcano with these step-by-step guides.

## Getting Started

- [[installation]] — Install Volcano on your system
- [[first-site]] — Create your first documentation site

## Advanced Topics

- [[customization]] — Customize themes and styling
- [[deployment]] — Deploy to production
```

Custom pages let you:
- Add descriptions to links
- Group items logically
- Include introductory content
- Control the exact order

### Use Auto-Index for Large Folders

Auto-index works well for:
- Reference sections with many pages
- Auto-generated API documentation
- Large file collections

### Check Generated Pages

Review auto-generated pages during development:

```bash
volcano serve ./docs
```

Visit each folder URL to ensure the generated index is acceptable.

## Navigation Integration

Auto-index pages are fully integrated:

- **Sidebar** — Folder appears in navigation tree
- **Breadcrumbs** — Shows path to folder
- **SEO** — Includes proper meta tags

## Related

- [[navigation]] — Sidebar and navigation features
- [[guides/organizing-content]] — File structure best practices
- [[theming]] — Customize auto-index styling

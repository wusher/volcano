# docs Theme

The default theme. Full-featured documentation chrome — sidebar tree, optional TOC, dark mode, comfortable reading typography.

```bash
volcano ./docs --theme docs --url="https://example.com"
```

## Light

![docs theme, light mode](/images/themes/docs-light.png)

## Dark

![docs theme, dark mode](/images/themes/docs-dark.png)

## Features

- Sidebar tree navigation with active-page highlight
- Auto-generated TOC sidebar on pages with 3+ headings
- Dark mode toggle (`t` key)
- Typography tuned for long-form technical reading
- Responsive — sidebar becomes a drawer below 768px
- Code blocks with syntax highlighting + copy buttons
- Accent gradient paints page H1, link colors, and H2 underlines

## Best For

Technical documentation, API references, user guides, software handbooks — anywhere readers need to scan a tree and jump around.

## Try It With Optional Features

```bash
volcano ./docs \
  --theme docs \
  --search \
  --breadcrumbs \
  --top-nav \
  --page-nav \
  --url="https://example.com"
```

See [Features](/features/) for what each flag adds.

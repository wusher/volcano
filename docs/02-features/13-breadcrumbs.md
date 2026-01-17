# Breadcrumbs

Breadcrumbs show the current page's location in your site hierarchy.

## Default Behavior

**Breadcrumbs are enabled by default.** You don't need to do anything to use them.

```
Home > Guides > Getting Started
```

## Disabling Breadcrumbs

To turn them off:

```bash
volcano ./docs --breadcrumbs=false
```

## What They Show

For a page at `/guides/advanced/configuration/`:

```
Home > Guides > Advanced > Configuration
```

Each item (except the last) is clickable.

## Features

### Interactive Links

Every breadcrumb except the current page is a link:
- Click to jump to any parent level
- Navigate up the hierarchy quickly

### SEO Markup

Breadcrumbs include Schema.org structured data:

```html
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "BreadcrumbList",
  "itemListElement": [...]
}
</script>
```

This helps search engines understand your site structure.

### Smart Display

- **Homepage**: No breadcrumbs (only one level)
- **Top-level pages**: Shows "Home > Page"
- **Nested pages**: Shows full path

### Responsive

On mobile, breadcrumbs automatically:
- Stack vertically if needed
- Maintain touch-friendly hit targets
- Stay accessible

## Styling

Breadcrumbs use your theme's colors and fonts. Location is above the page content, below the header.

## Accessibility

- Proper semantic HTML (`nav`, `ol`, `li`)
- `aria-label="Breadcrumb"` for screen readers
- `aria-current="page"` on current item
- Keyboard navigable

## How They're Generated

Breadcrumbs are built from:
1. Your folder structure
2. File/folder names (cleaned and title-cased)
3. H1 headings (if available)

Example:

```
Folder: guides/getting-started/
File: installation.md
H1: Installing Volcano

Breadcrumb: Home > Guides > Getting Started > Installing Volcano
```

## Common Patterns

### Documentation Sites

```
Docs Home > API Reference > Authentication
Docs Home > Guides > Installation
```

### Knowledge Base

```
Help Center > Account Settings > Password Reset
Help Center > Billing > Invoices
```

### Blog

```
Home > Posts > 2024 > My Article
```

## When to Disable

Consider disabling breadcrumbs when:
- Your site is very flat (few nested levels)
- You're using top navigation that shows structure
- You want a minimal, distraction-free layout

## Related

- [[navigation]] — Navigation overview
- [[top-navigation]] — Horizontal top nav bar
- [[page-navigation]] — Previous/next links

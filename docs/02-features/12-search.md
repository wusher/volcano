# Search

Add a command palette that searches your entire site.

## Enabling Search

```bash
volcano ./docs --search
```

## What You Get

A keyboard-activated command palette (Cmd+K or Ctrl+K) that searches:
- Page titles
- Heading text (h2-h4)
- URL paths

## How It Works

When enabled, Volcano generates:

1. **Search Index** (`search-index.json`) - Contains all searchable content
2. **Search UI** - Command palette with keyboard navigation
3. **Keyboard Shortcut** - Cmd+K (Mac) or Ctrl+K (Windows/Linux)

## Using Search

**Open search:**
- Press Cmd+K (or Ctrl+K)
- Click the search icon in the header (mobile)

**Navigate results:**
- Type to filter
- ↑/↓ arrow keys to select
- Enter to navigate
- Esc to close

## Search Features

### Multi-term Queries

```
"getting started"     → Finds pages/headings with both words
"api authentication"  → Searches titles and headings for both terms
```

### Result Types

Search shows different result types:

- **📄 Pages** - Matches page titles or URL paths
- **🔗 Headings** - Matches section headings (h2-h4)

### Ranking

Results appear in order:
1. Page title and path matches
2. Heading matches within pages

### Context

Results show:
- Page titles for page matches
- Parent page title for heading matches

## Performance

The search index is:
- Generated at build time
- Loaded on-demand (only when search is opened)
- Cached by the browser
- Typically 50-200KB for documentation sites

## Limitations

**Titles and headings only:** Search indexes page titles, headings, and URL paths. It does not search page body content.

**Client-side only:** Search runs in the browser, not on a server. Very large sites (1000+ pages) may have slower search loading.

**No fuzzy search:** Typos won't match. "autentication" won't find "authentication".

**Substring matching only:** Uses simple substring matching without stemming or natural language processing.

## Customization

Search uses your site's theme colors automatically. The command palette adapts to dark mode.

## Keyboard Shortcuts

With search enabled:

| Key | Action |
|-----|--------|
| Cmd+K / Ctrl+K | Open search |
| ↑ / ↓ | Navigate results |
| Enter | Go to selected result |
| Esc | Close search |

## Mobile Support

On mobile devices:
- Tap the search icon in the header
- Type in the search box
- Tap results to navigate

## Deployment

When deploying with search enabled, ensure `search-index.json` is uploaded along with your HTML files.

## Example

```bash
volcano ./docs \
  --search \
  --instant-nav \
  --title="My Documentation"
```

Search works great with instant navigation for a fast, app-like experience.

## Related

- [[navigation]] — Other navigation features
- [[instant-navigation]] — Faster page loads
- [[keyboard-shortcuts]] — All keyboard shortcuts

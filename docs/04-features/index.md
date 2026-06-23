# Features

Optional behaviors you turn on. Each is independent.

Every feature here has both a CLI flag (one-off) and a `volcano.json` key (durable). See **[Configuration](/cli/configuration/)** for the full lookup table and precedence rules.

## Search

> **Configure:** `--search` · `"search": true`

```bash
volcano ./docs --search --url="https://example.com"
```

Adds a Cmd+K (Mac) / Ctrl+K (Windows/Linux) command palette that searches page titles, headings, and URL paths.

![Search palette open](/images/ui/search-palette.png)

What it indexes: H1 page titles, H2–H4 headings, URL paths. It does **not** index page body content. No fuzzy matching (typos won't match). The index is generated at build time and lazy-loaded on first open. Typically 50-200 KB for a docs site.

## Breadcrumbs

> **Configure:** `--breadcrumbs` · `"breadcrumbs": true`

```bash
volcano ./docs --breadcrumbs --url="https://example.com"
```

Shows the page's path in the site hierarchy, above the content:

![Breadcrumb trail](/images/ui/breadcrumbs.png)

Hidden on the homepage (no useful path). Each segment is linked. Includes schema.org `BreadcrumbList` JSON-LD for search engines.

## Top Navigation

> **Configure:** `--top-nav` · `"topNav": true`

```bash
volcano ./docs --top-nav --url="https://example.com"
```

A horizontal bar above the content showing your root-level pages and folders:

![Top navigation bar](/images/ui/top-nav.png)

Caps at 8 items — anything beyond that is dropped. Active section is highlighted. On mobile it collapses into the hamburger menu.

## Previous / Next Links

> **Configure:** `--page-nav` · `"pageNav": true`

```bash
volcano ./docs --page-nav --url="https://example.com"
```

Adds "← Previous" / "Next →" links at the bottom of each page. Order follows the sidebar tree (files first, then folders, depth-first). Good for tutorials and books, less useful for reference sites where readers jump around.

When enabled, `n` and `p` keys navigate. Pages without a previous/next page hide the corresponding link.

## Table of Contents

Auto-generated, no flag needed. Pages with 3+ headings get a right-side TOC of `##`, `###`, `####` headings. Click to jump, scroll to update the active highlight, URL anchor stays in sync.

- Visible on screens ≥ 1280px wide
- Hidden 768–1279px (limited horizontal space)
- Mobile gets a TOC toggle button in the header

## Instant Navigation

> **Configure:** `--instant-nav` · `"instantNav": true`

```bash
volcano ./docs --instant-nav --url="https://example.com"
```

Prefetches the linked page when you hover for 65ms+, then swaps content via AJAX on click. Pages appear in under 10ms once cached. Works with View Transitions for smooth animation. See [the advanced page](/advanced/instant-navigation/) for the full mechanism.

## Progressive Web App

> **Configure:** `--pwa` · `"pwa": true`

```bash
volcano ./docs --pwa --url="https://example.com"
```

Generates a `manifest.json` and service worker. Users can install your site as an app; cached pages load offline. See [the advanced page](/advanced/pwa/) for the full setup.

## Keyboard Shortcuts

Press `?` anywhere to see the full list — it adapts to which features you have enabled.

Always available:

| Key | Action |
|-----|--------|
| `t` | Toggle theme (light/dark) |
| `z` | Toggle zen mode (hide sidebar) |
| `h` | Go to homepage |
| `?` | Show shortcuts |
| `Esc` | Close modal |

Feature-gated:

| Key | Action | Requires |
|-----|--------|----------|
| `⌘K` / `Ctrl+K` | Open search | `--search` |
| `n` / `p` | Next / previous page | `--page-nav` |
| `←` / `→` | Previous / next heading | `--theme presentation` |

## Combining Features

These flags compose freely. A typical "give me everything" build:

```bash
volcano ./docs \
  --url="https://docs.example.com" \
  --search \
  --breadcrumbs \
  --top-nav \
  --page-nav \
  --instant-nav
```

Or put them in `volcano.json` once and stop typing flags — see [Configuration](/cli/configuration/) for the equivalent JSON, precedence rules, and a worked example.

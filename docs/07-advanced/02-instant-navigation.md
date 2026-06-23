# Instant Navigation

Sub-10ms page loads via hover prefetching.

## Enable It

```bash
volcano ./docs --instant-nav --url="https://example.com"
```

## How It Works

1. Hover a link for 65ms+ → Volcano fetches the page in the background
2. Browser caches it
3. Click → page loads from cache (typically < 10ms)
4. Content swaps via AJAX (no full reload). Scroll restored, history updated.
5. View Transitions API animates the swap on browsers that support it

## What Gets Prefetched

All internal links — sidebar, breadcrumbs, in-content links, page nav, wiki links. External links and `data-no-instant` opt out:

```html
<a href="/slow-page/" data-no-instant>Skip prefetch for this one</a>
```

Each URL is only fetched once per session.

## Why 65ms?

Long enough to filter out cursor-passing-by. Short enough to feel instant on intentional hovers. Empirically, intentional hovers cluster above 80ms; accidental ones cluster below 40ms.

## When to Skip It

- **Very large pages (>1MB).** Prefetching them on hover wastes bandwidth.
- **Sites with few internal links.** No payoff.
- **Users on metered connections.** No control yet to gate prefetch on `navigator.connection`.

## Verifying It Works

Open DevTools → Network tab. Hover a sidebar link. You should see a request fire after ~65ms. Click — the request shows "from disk cache" or "from memory cache" instead of a fresh fetch.

## Technical Details

- `fetch()` API for the prefetch and the AJAX page load
- `history.pushState()` for URL updates
- View Transitions API for the animation (graceful degradation when unsupported)
- Browser HTTP cache (not custom storage)

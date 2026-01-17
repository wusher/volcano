# Instant Navigation

Enable hover prefetching for near-instant page loads.

## Enabling Instant Navigation

```bash
volcano ./docs --instant-nav
```

## How It Works

When you hover over a link for 65ms:
1. Volcano prefetches the page in the background
2. The page is cached by your browser
3. When you click, the page loads instantly (often <10ms)

With smooth view transitions, navigation feels like a single-page app.

## Features

### Hover Prefetching

Links are prefetched when you hover over them:
- Sidebar navigation links
- Breadcrumb links
- In-content links
- Page navigation links

### Smart Throttling

- Only prefetches after 65ms hover (prevents false positives)
- Each URL is only prefetched once
- Respects browser cache
- Minimal bandwidth usage

### AJAX Page Loading

After prefetch, clicks load pages via AJAX:
- No full page reload
- Content swaps smoothly
- Preserves scroll state
- Updates browser history

### Works Everywhere

Instant navigation works with:
- All internal links
- Sidebar navigation
- Breadcrumbs
- Wiki links in content
- Page prev/next navigation

## Performance

Instant navigation is efficient:

**Bandwidth:** Only prefetches pages you're likely to visit (hover = interest)

**Cache-friendly:** Uses browser's HTTP cache, not custom storage

**Fast:** Prefetching happens in the background, doesn't block anything

## User Experience

With instant navigation enabled:
- Hovering a link → Prefetch starts
- Clicking the link → Page appears instantly
- Smooth transition (with View Transitions)
- Feels like a native app

## Browser Support

Works in all modern browsers:
- Chrome/Edge ✓
- Firefox ✓
- Safari ✓

Graceful degradation: If prefetch fails, pages load normally.

## Combining with Other Features

Great with:

```bash
volcano ./docs \
  --instant-nav \
  --search \
  --pwa
```

- **Search**: Jump to pages instantly
- **PWA**: Offline + instant loading
- **View Transitions**: Smooth animations (enabled by default)

## Technical Details

Instant nav uses:
- HTML `<link rel="prefetch">` for resource hints
- `fetch()` API for AJAX page loads
- `history.pushState()` for URL updates
- View Transitions API for animations

## Disabling for Specific Links

Add `data-no-instant` to any link to skip instant nav:

```html
<a href="/slow-page/" data-no-instant>Not instant</a>
```

Or use external links (different domain) - they always bypass instant nav.

## Accessibility

Instant navigation maintains:
- Browser back/forward buttons
- Bookmarking
- Sharing links
- Keyboard navigation
- Screen reader compatibility

URLs update normally, so assistive tech works as expected.

## Performance Monitoring

To see instant nav in action:
1. Open browser DevTools → Network tab
2. Hover over links
3. Watch prefetch requests appear
4. Click the link
5. See the instant load (from cache)

## When to Use

Instant navigation is great for:
- Documentation sites
- Knowledge bases
- Blogs
- Any site with lots of internal links

## When to Skip

Consider skipping instant nav when:
- Your pages are very large (>1MB)
- Your site has few pages
- Users rarely click internal links

## Example

For a documentation site with many pages:

```bash
volcano ./docs \
  -o ./public \
  --instant-nav \
  --search \
  --title="My Docs"
```

Users get:
- Instant page loads on click
- Fast search results
- Smooth animations
- App-like experience

## Related

- [[search]] — Fast site-wide search
- [[navigation]] — Navigation overview

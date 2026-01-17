# Progressive Web App (PWA)

Enable PWA support to make your site installable and work offline.

## Enabling PWA

Add the `--pwa` flag when generating your site:

```bash
volcano ./docs --pwa --title="My Docs"
```

## What Gets Generated

When PWA is enabled, Volcano automatically creates:

### 1. Web App Manifest (`manifest.json`)

A JSON file that tells browsers how to install your site:

```json
{
  "name": "My Docs",
  "short_name": "My Docs",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#ffffff",
  "theme_color": "#ffffff"
}
```

### 2. Service Worker (`sw.js`)

A script that caches your site's assets for offline access. The service worker:
- Caches HTML pages, CSS, and JavaScript
- Serves cached content when offline
- Updates the cache in the background

### 3. Manifest Link in HTML

```html
<link rel="manifest" href="/manifest.json">
```

## Installation Experience

When PWA is enabled, users can install your site:

**Desktop:**
- Chrome/Edge: Shows an install icon in the address bar
- Safari: File → Add to Dock

**Mobile:**
- iOS Safari: Share → Add to Home Screen
- Android Chrome: Menu → Install app

Once installed, your site:
- Appears as an app icon on the home screen
- Opens in a standalone window (no browser UI)
- Works offline (cached pages load without internet)

## Offline Behavior

The service worker implements a "network first, cache fallback" strategy:

1. **Online**: Fetches fresh content from the network
2. **Offline**: Serves cached content if available
3. **Background Updates**: Updates cache when online

## Configuration Options

PWA settings automatically use your site metadata:

| Setting | Source | Default |
|---------|--------|---------|
| `name` | `--title` flag | "My Site" |
| `start_url` | `--url` flag | "/" |
| `theme_color` | `--accent-color` or theme | "#ffffff" |
| `background_color` | Theme default | "#ffffff" |

## Complete Example

```bash
volcano ./docs \
  --pwa \
  --title="My Documentation" \
  --url="https://docs.example.com" \
  --accent-color="#0066cc" \
  --favicon="./icon.png"
```

This generates a fully installable PWA with:
- Custom branding (title, icon, colors)
- Offline support for all pages
- Proper canonical URLs

## Testing PWA

### Lighthouse

Use Chrome DevTools Lighthouse to audit PWA features:

1. Open your site in Chrome
2. Open DevTools (F12)
3. Go to Lighthouse tab
4. Run PWA audit

Look for:
- ✓ Installable
- ✓ Works offline
- ✓ Has a manifest

### Manual Testing

**Test offline mode:**
1. Load your site
2. Open DevTools → Network tab
3. Set throttling to "Offline"
4. Refresh the page
5. Site should still load (from cache)

**Test installation:**
1. Look for install prompt in address bar
2. Click install
3. Verify app opens standalone
4. Check home screen icon appears

## Limitations

**For deployed sites:** PWA features are for production builds. Dynamic features (forms, real-time updates) require additional server setup beyond Volcano's scope.

**Cache invalidation:** The service worker caches aggressively. When you update your site:
1. Users get the old version initially
2. Service worker fetches updates in background
3. Updates apply on next page load

**Browser support:**
- Chrome/Edge: Full support
- Firefox: Partial (no install prompt)
- Safari iOS: Limited (no background sync)

## Removing PWA

To disable PWA features:

1. Regenerate without `--pwa` flag
2. Remove `manifest.json` and `sw.js` from output
3. Users who installed will keep cached version until they uninstall

## Related

- [[guides/theming]] — Theme colors affect PWA appearance
- [[reference/cli]] — All CLI flags

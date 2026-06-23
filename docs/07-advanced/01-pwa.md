# Progressive Web App

Make your site installable and offline-capable.

## Enable It

```bash
volcano ./docs --pwa --url="https://docs.example.com" --title="My Docs"
```

Volcano generates a manifest, a service worker, and the `<link rel="manifest">` tag.

## What Gets Created

`manifest.json`:

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

`sw.js` — caches HTML, CSS, JS on first visit. Strategy: network first, cache fallback. Background refresh keeps cache current.

## What Users See

**Desktop:**
- Chrome / Edge: install icon in the address bar
- Safari: File → Add to Dock

**Mobile:**
- iOS Safari: Share → Add to Home Screen
- Android Chrome: Menu → Install app

Installed sites get a home-screen icon, open in a standalone window (no browser chrome), and load cached pages when offline.

## Settings Source

| Manifest field | Comes from |
|----------------|------------|
| `name` / `short_name` | `--title` |
| `start_url` | `--url` |
| `theme_color` | `--accent-color` or theme default |
| `background_color` | Theme default |

## Limits

- **Cache invalidation is aggressive.** Users see the old version once, the service worker refreshes in the background, the new version appears on the next page load. Acceptable for docs, not for fast-changing apps.
- **No background sync** on iOS Safari.
- **Firefox** doesn't show an install prompt.

## Removing It

Rebuild without `--pwa` and remove `manifest.json` + `sw.js` from your deployment. Already-installed users keep their cached copy until they uninstall.

## Testing

Chrome DevTools → Lighthouse → PWA audit covers the basics (installable, has manifest, works offline). For offline testing: DevTools → Network → throttling = Offline, then reload.

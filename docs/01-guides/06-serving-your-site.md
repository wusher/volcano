# Serving Your Site

Preview your site locally with Volcano's built-in server.

:::note Local Development Only
This guide covers Volcano's built-in server for **local development and testing**. For production deployments, see [[deploying-your-site]].
:::

## Basic Usage

Just point serve at your markdown source directory:

```bash
volcano serve ./docs
```

That's it. No build step needed.

**How it works:**
- Renders pages fresh on each request
- Changes visible on browser refresh
- See updates immediately as you edit

## With Options

You can pass the same flags as build:

```bash
volcano serve ./docs \
  --title="My Docs" \
  --theme blog \
  --instant-nav \
  --search
```

All settings apply during rendering.

## Custom Port

The default port is 1776. Change it with `-p`:

```bash
volcano serve -p 8080 ./docs
```

## Server Output

```
⚡ Serving at http://localhost:1776

Press Ctrl+C to stop
```

## Development Workflow

Start the server once:

```bash
volcano serve ./docs
```

Then:
1. Edit your markdown files
2. Save
3. Refresh browser
4. See changes immediately

No build or regeneration needed.

## How It Works

- Generates navigation on startup
- Renders each page on request
- Shows broken link warnings inline
- All flags (title, theme, etc.) apply during rendering

## Auto-Refresh

Volcano doesn't have built-in auto-refresh. For automatic browser reloading, use a browser extension:

- **Chrome/Edge**: LiveReload, Auto Refresh Plus
- **Firefox**: Auto Reload
- **Safari**: LiveReload

Or use an external tool like `browser-sync`:

```bash
browser-sync start --proxy localhost:1776 --files "docs/**/*.md"
```

## Common Issues

### Port Already in Use

```
Error: listen tcp :1776: bind: address already in use
```

**Solution 1:** Use a different port
```bash
volcano serve -p 8080 ./docs
```

**Solution 2:** Kill the process using the port
```bash
# Find the process
lsof -i :1776

# Kill it
kill <PID>
```

### Changes Not Appearing

- Try hard refresh: Ctrl+Shift+R (or Cmd+Shift+R)
- Clear browser cache
- Check you're editing the source files (not generated output)

## Production Serving

`volcano serve` is for development only. For production:

- Use a real web server (nginx, Apache, Caddy)
- Or deploy to a static host (Netlify, Vercel, GitHub Pages)
- Or use a CDN

See [[deploying-your-site]] for deployment guides.

## Next Steps

- [[building-your-site]] — Build options and configuration
- [[deploying-your-site]] — Deploy to production
- [[reference/cli]] — All server options

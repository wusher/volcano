# Deploying

Volcano produces plain static HTML, CSS, and JS — any host that serves files will do. The only build-time requirement is `--url` (used for canonical and Open Graph tags).

```bash
volcano ./docs -o ./public --url="https://docs.example.com" --title="My Site"
```

After that, `./public` is your deployable artifact.

## GitHub Pages

The cleanest setup uses GitHub Actions. Save this as `.github/workflows/deploy.yml`:

```yaml
name: Deploy
on:
  push:
    branches: [main]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.24' }
      - run: go install github.com/wusher/volcano@latest
      - run: volcano ./docs -o ./public --url="https://${{ github.repository_owner }}.github.io/${{ github.event.repository.name }}"
      - uses: peaceiris/actions-gh-pages@v4
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

Then enable Pages on the `gh-pages` branch in repo settings.

**Subpath note:** if your site lives at `https://you.github.io/repo-name/` (not a custom domain), put the path in `--url`. Volcano handles the prefix everywhere.

## Netlify

`netlify.toml` at your repo root:

```toml
[build]
  command = "go install github.com/wusher/volcano@latest && volcano ./docs -o ./public --url='https://your-site.netlify.app'"
  publish = "public"

[build.environment]
  GO_VERSION = "1.24"
```

Or drag and drop: build locally with `volcano ./docs -o ./public --url=...`, then drag `./public` into the Netlify deploy zone.

## Vercel

`vercel.json`:

```json
{
  "buildCommand": "go install github.com/wusher/volcano@latest && volcano ./docs -o ./public --url='https://your-site.vercel.app'",
  "outputDirectory": "public",
  "installCommand": "echo skip"
}
```

## Cloudflare Pages

In the Cloudflare dashboard, set:

- **Build command:** `go install github.com/wusher/volcano@latest && volcano ./docs -o ./public --url="https://my-docs.pages.dev"`
- **Output directory:** `public`

## Self-Hosting

### Nginx

```nginx
server {
    listen 80;
    server_name docs.example.com;
    root /var/www/docs;

    location / {
        try_files $uri $uri/ $uri/index.html =404;
    }
    error_page 404 /404.html;
}
```

### Caddy

```caddyfile
docs.example.com {
    root * /var/www/docs
    file_server
    try_files {path} {path}/ {path}/index.html
    handle_errors {
        @404 expression {http.error.status_code} == 404
        rewrite @404 /404.html
        file_server
    }
}
```

### Apache (.htaccess)

```apache
RewriteEngine On
RewriteCond %{REQUEST_FILENAME} !-f
RewriteCond %{REQUEST_FILENAME} !-d
RewriteRule ^(.*)$ $1/index.html [L]
ErrorDocument 404 /404.html
```

## What Volcano Generates

Every build produces a `404.html` with full site navigation — most static hosts serve this for missing URLs automatically.

If you enabled features, you'll also see:

- `search-index.json` + `search.js` — with `--search`
- `manifest.json` + `sw.js` — with `--pwa`
- `styles.css` — always

Make sure your host uploads everything in `./public`, not just `*.html`.

## Pre-Flight Checklist

- [ ] `--url` matches your production domain
- [ ] `volcano ./docs --url=...` exits cleanly (broken-link validation passes)
- [ ] Custom favicon set with `--favicon` if you have one
- [ ] `--og-image` set if you want rich social previews

## Next

- **[CLI reference](/cli/)** — every flag
- **[Features](/features/)** — opt-in behaviors worth enabling for production

# Deploying Your Site

Publish your Volcano site to the web.

## Generating for Production

Before deploying, generate your site with production settings:

```bash
volcano ./docs \
  -o ./public \
  --title="My Documentation" \
  --url="https://docs.example.com"
```

The `--url` flag is important for:
- Canonical URL meta tags
- Correct Open Graph URLs
- SEO best practices

## Output Structure

Volcano generates static HTML files that can be served by any web server:

```
public/
├── index.html
├── getting-started/
│   └── index.html
├── guides/
│   ├── index.html
│   └── installation/
│       └── index.html
└── styles.css
```

No special server configuration is needed — just serve the files.

## GitHub Pages

### Using GitHub Actions

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install Volcano
        run: go install github.com/wusher/volcano@latest

      - name: Build site
        run: |
          volcano ./docs \
            -o ./public \
            --title="My Project" \
            --url="https://username.github.io/repo"

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

### Manual Deployment

1. Generate your site:
   ```bash
   volcano ./docs -o ./public --url="https://username.github.io/repo"
   ```

2. Create a `gh-pages` branch with the output:
   ```bash
   git checkout --orphan gh-pages
   git rm -rf .
   cp -r public/* .
   git add .
   git commit -m "Deploy site"
   git push origin gh-pages
   ```

3. Enable GitHub Pages in your repository settings, selecting the `gh-pages` branch.

## Netlify

### Using netlify.toml

Create `netlify.toml` in your repository root:

```toml
[build]
  command = "go install github.com/wusher/volcano@latest && volcano ./docs -o ./public --title='My Site' --url='https://your-site.netlify.app'"
  publish = "public"

[build.environment]
  GO_VERSION = "1.21"
```

### Using Netlify CLI

```bash
# Install Netlify CLI
npm install -g netlify-cli

# Build locally
volcano ./docs -o ./public

# Deploy
netlify deploy --prod --dir=public
```

### Drag and Drop

1. Generate your site: `volcano ./docs -o ./public`
2. Go to [app.netlify.com](https://app.netlify.com)
3. Drag your `public` folder to the deploy zone

## Vercel

### Using vercel.json

Create `vercel.json`:

```json
{
  "buildCommand": "go install github.com/wusher/volcano@latest && volcano ./docs -o ./public --title='My Site'",
  "outputDirectory": "public",
  "installCommand": "echo 'No install needed'"
}
```

### Using Vercel CLI

```bash
# Install Vercel CLI
npm install -g vercel

# Build locally
volcano ./docs -o ./public

# Deploy
vercel --prod
```

## Cloudflare Pages

### Using wrangler.toml

Create `wrangler.toml`:

```toml
name = "my-docs"
compatibility_date = "2024-01-01"

[site]
bucket = "./public"
```

Build command in Cloudflare dashboard:
```bash
go install github.com/wusher/volcano@latest && volcano ./docs -o ./public
```

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
}
```

### Apache

Create `.htaccess` in your document root:

```apache
RewriteEngine On
RewriteCond %{REQUEST_FILENAME} !-f
RewriteCond %{REQUEST_FILENAME} !-d
RewriteRule ^(.*)$ $1/index.html [L]
```

### Caddy

```caddyfile
docs.example.com {
    root * /var/www/docs
    file_server
    try_files {path} {path}/ {path}/index.html
}
```

## CI/CD Automation

### Generic CI Script

```bash
#!/bin/bash
set -e

# Install Volcano
go install github.com/wusher/volcano@latest

# Generate site
volcano ./docs \
  -o ./public \
  --title="$SITE_TITLE" \
  --url="$SITE_URL" \
  --author="$SITE_AUTHOR"

echo "Site generated in ./public"
```

### Environment Variables

Use environment variables for sensitive or environment-specific values:

```bash
export SITE_TITLE="My Documentation"
export SITE_URL="https://docs.example.com"
export SITE_AUTHOR="My Team"

volcano ./docs \
  -o ./public \
  --title="$SITE_TITLE" \
  --url="$SITE_URL" \
  --author="$SITE_AUTHOR"
```

## Custom Domains

Most hosting platforms support custom domains:

1. **Add domain in platform settings** — Point your custom domain to the hosting platform
2. **Configure DNS** — Add CNAME or A records as instructed
3. **Update Volcano URL** — Generate with the custom domain:
   ```bash
   volcano ./docs --url="https://docs.mycompany.com"
   ```
4. **Enable HTTPS** — Most platforms provide free SSL certificates

## Deployment Checklist

Before deploying:

- [ ] Set the correct `--url` for your production domain
- [ ] Verify all links work locally
- [ ] Check page titles and meta descriptions
- [ ] Test on mobile devices
- [ ] Validate the site structure matches expectations
- [ ] Review the generated HTML for any issues

## Next Steps

- [[reference/cli]] — See all generation options
- [[features/seo-and-meta]] — Optimize for search engines

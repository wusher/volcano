# Documentation Site

A typical documentation site structure:

```
docs/
├── index.md                 # Homepage
├── getting-started.md
├── guides/
│   ├── index.md
│   ├── 01-installation.md
│   └── 02-configuration.md
├── reference/
│   ├── index.md
│   └── cli.md
└── faq.md
```

## Build Command

```bash
volcano ./docs \
  -o ./public \
  --title="My Project" \
  --url="https://docs.example.com" \
  --page-nav \
  --search
```

## GitHub Actions

```yaml
name: Deploy Docs

on:
  push:
    branches: [main]
    paths: ['docs/**']

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Volcano
        run: go install github.com/wusher/volcano@latest

      - name: Build
        run: volcano ./docs -o ./public --title="My Project" --url="${{ vars.SITE_URL }}"

      - name: Deploy
        uses: peaceiris/actions-gh-pages@v4
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

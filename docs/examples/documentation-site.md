# Documentation Site Example

How to create a comprehensive documentation site like this one.

## Structure

```
docs/
├── index.md                    # Homepage/introduction
├── getting-started.md          # Quick start guide
│
├── guides/                     # How-to tutorials
│   ├── index.md               # Guides overview
│   ├── 01-installation.md     # Numbered for order
│   ├── 02-configuration.md
│   ├── 03-customization.md
│   └── 04-deployment.md
│
├── features/                   # Feature documentation
│   ├── index.md
│   ├── feature-a.md
│   ├── feature-b.md
│   └── feature-c.md
│
├── reference/                  # Technical reference
│   ├── index.md
│   ├── api.md
│   ├── cli.md
│   └── configuration.md
│
└── faq.md                      # Frequently asked questions
```

## Homepage (index.md)

```markdown
# My Project

Welcome to My Project documentation.

## Quick Links

- [[getting-started]] — Get up and running in 5 minutes
- [[guides/installation]] — Detailed installation guide
- [[reference/api]] — API reference

## Features

- **Fast** — Built for performance
- **Simple** — Easy to use
- **Flexible** — Highly customizable

## Getting Help

- [GitHub Issues](https://github.com/example/project/issues)
- [Discord Community](https://discord.gg/example)
```

## Guides Section

### guides/index.md

```markdown
# Guides

Step-by-step tutorials for common tasks.

## Getting Started

- [[01-installation|Installation]] — Install on your system
- [[02-configuration|Configuration]] — Configure your setup

## Advanced

- [[03-customization|Customization]] — Customize behavior
- [[04-deployment|Deployment]] — Deploy to production
```

### guides/01-installation.md

```markdown
# Installation

Install My Project on your system.

## Requirements

- Node.js 18 or later
- npm or yarn

## Install via npm

```bash
npm install -g my-project
```

## Verify Installation

```bash
my-project --version
```

:::tip
If you encounter permission issues, see [[faq#permissions]].
:::

## Next Steps

Continue to [[02-configuration|Configuration]].
```

## Reference Section

### reference/cli.md

```markdown
# CLI Reference

Complete command-line interface documentation.

## Commands

### `init`

Initialize a new project:

```bash
my-project init [name]
```

**Arguments:**

| Argument | Description | Default |
|----------|-------------|---------|
| `name` | Project name | `my-project` |

**Options:**

| Option | Description |
|--------|-------------|
| `--template` | Template to use |
| `--force` | Overwrite existing files |

**Examples:**

```bash
my-project init my-app
my-project init my-app --template minimal
```

### `build`

Build the project:

```bash
my-project build [options]
```

**Options:**

| Option | Description | Default |
|--------|-------------|---------|
| `-o, --output` | Output directory | `./dist` |
| `--minify` | Minify output | `true` |
| `--sourcemaps` | Generate sourcemaps | `false` |
```

## Build Command

```bash
volcano ./docs \
  -o ./public \
  --title="My Project Documentation" \
  --url="https://docs.myproject.com" \
  --page-nav \
  --favicon="./favicon.ico"
```

### Recommended Flags

| Flag | Purpose |
|------|---------|
| `--page-nav` | Adds previous/next navigation |
| `--url` | Required for SEO and canonical URLs |
| `--favicon` | Brand your documentation |

## CI/CD Integration

### GitHub Actions

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
        run: go install github.com/example/volcano@latest

      - name: Build Documentation
        run: |
          volcano ./docs \
            -o ./public \
            --title="My Project" \
            --url="https://docs.myproject.com"

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v4
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

## Styling Tips

### Custom CSS for Docs

```css
/* Wider content for documentation */
:root {
  --content-max-width: 900px;
}

/* Prominent code blocks */
.prose pre {
  border: 1px solid var(--color-border);
  border-radius: 8px;
}

/* Better tables */
.prose table {
  width: 100%;
}

.prose th {
  background: var(--color-bg-muted);
}
```

## Best Practices

### Organization

1. **Clear hierarchy** — Use folders for major sections
2. **Numbered guides** — Use prefixes for sequential content
3. **Descriptive names** — `01-installation.md` not `01-install.md`

### Content

1. **Start with overview** — Each section has an index.md
2. **Link liberally** — Use wiki links to connect related pages
3. **Include examples** — Real code examples, not just descriptions

### Navigation

1. **Logical order** — Most important content first
2. **Previous/next** — Enable `--page-nav` for guides
3. **Cross-references** — Link between sections

## Related

- [[guides/organizing-content]] — Content structure guide
- [[features/navigation]] — Navigation features
- [[reference/cli]] — CLI reference

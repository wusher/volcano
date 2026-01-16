# Creating Your First Site

## Structure

```
docs/
├── index.md           # Homepage
├── getting-started.md
├── guides/
│   ├── index.md       # Section landing page
│   ├── installation.md
│   └── configuration.md
└── reference/
    ├── index.md
    └── cli.md
```

## Generate

```bash
volcano ./docs -o ./public --title="My Site"
```

## Preview

```bash
volcano -s ./public
```

Open http://localhost:1776

## Development Mode

Serve source files directly (regenerates on each request):

```bash
volcano -s ./docs
```

## Output Structure

```
public/
├── index.html              # /
├── getting-started/
│   └── index.html          # /getting-started/
├── guides/
│   ├── index.html          # /guides/
│   └── installation/
│       └── index.html      # /guides/installation/
```

Each page includes sidebar navigation, breadcrumbs, and table of contents.

## Next

- [[organizing-content]] — Naming conventions
- [[customizing-appearance]] — Themes and CSS
- [[deploying-your-site]] — Publish online

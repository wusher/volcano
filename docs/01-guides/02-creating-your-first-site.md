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

## Preview While You Write

The fastest loop — point `serve` at your source folder and Volcano renders each request on the fly:

```bash
volcano serve ./docs
```

Open http://localhost:1776. Edit a markdown file, refresh the browser, see the change.

## Generate Static Files

When you're ready to ship, generate the static site. `--url` is required so canonical and Open Graph tags resolve correctly:

```bash
volcano ./docs -o ./public --title="My Site" --url="https://example.com"
```

You can also serve the built output (handy for testing pre-deploy):

```bash
volcano serve ./public
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

Each page includes sidebar navigation and an auto-generated table of contents (when the page has 3+ headings).

## Next

- [[organizing-content]] — Naming conventions
- [[customizing-appearance]] — Themes and CSS
- [[deploying-your-site]] — Publish online

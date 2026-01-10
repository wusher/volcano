# Examples

Real-world examples and templates for common use cases.

## Example Projects

### [[documentation-site|Documentation Site]]

A technical documentation site like this one:

- Organized sections with nested folders
- API references and guides
- Code examples with syntax highlighting
- Search and navigation

**Best for:** Software documentation, API references, user guides.

### [[blog|Blog]]

A blog with date-organized posts:

- Chronological posts with date prefixes
- Blog-focused theme
- Reading time estimates
- Clean article layouts

**Best for:** Personal blogs, company news, changelogs.

### [[knowledge-base|Knowledge Base]]

An interconnected wiki-style knowledge base:

- Wiki links between pages
- Obsidian vault compatibility
- Flat or nested structure
- Cross-referencing notes

**Best for:** Personal wikis, team knowledge bases, note publishing.

## Quick Start Templates

### Minimal Documentation

```
docs/
├── index.md           # Homepage
├── getting-started.md # Quick start
└── reference.md       # API/CLI reference
```

```bash
volcano ./docs -o ./public --title="My Docs"
```

### Structured Guide

```
docs/
├── index.md
├── guides/
│   ├── index.md
│   ├── 01-installation.md
│   └── 02-configuration.md
├── reference/
│   ├── index.md
│   └── cli.md
└── faq.md
```

```bash
volcano ./docs -o ./public --title="My Project" --page-nav
```

### Blog

```
posts/
├── index.md
├── 2024-01-15-hello-world.md
├── 2024-02-01-getting-started.md
└── 2024-03-10-advanced-tips.md
```

```bash
volcano ./posts -o ./blog --theme blog --title="My Blog"
```

### Wiki

```
wiki/
├── index.md
├── projects/
│   ├── project-a.md
│   └── project-b.md
├── notes/
│   ├── meeting-notes.md
│   └── ideas.md
└── references/
    └── resources.md
```

```bash
volcano ./wiki -o ./public --title="My Wiki"
```

## Configuration Examples

### Production Build

```bash
volcano ./docs \
  -o ./public \
  --title="My Project" \
  --url="https://docs.example.com" \
  --author="My Team" \
  --og-image="https://docs.example.com/og.png" \
  --favicon="./favicon.ico"
```

### Development

```bash
volcano -s -p 3000 ./docs
```

### Custom Theme

```bash
# Export CSS skeleton
volcano css -o theme.css

# Use custom theme
volcano ./docs --css ./theme.css
```

## Related

- [[guides/creating-your-first-site]] — Step-by-step tutorial
- [[guides/organizing-content]] — Content structure tips
- [[reference/cli]] — All CLI options

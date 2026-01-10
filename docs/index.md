# Volcano

A zero-dependency static site generator that transforms markdown folders into beautiful, navigable documentation sites.

## Why Volcano?

Volcano is built for developers and writers who want to turn a folder of markdown files into a polished website without configuration files, complex setups, or external dependencies.

**Zero Configuration** — Just point Volcano at your markdown folder and get a site. No YAML files, no plugins to install, no build pipelines to configure.

**Instant Navigation** — Automatically generates a hierarchical sidebar, breadcrumbs, and table of contents from your folder structure.

**Beautiful by Default** — Ships with carefully designed themes that look professional out of the box, with full dark mode support.

**Single Binary** — One executable file. No Node.js, no Python, no Ruby. Download and run.

## Key Features

- **GitHub Flavored Markdown** — Tables, task lists, strikethrough, footnotes, and more
- **Wiki-Style Links** — Use `[[Page Name]]` syntax for easy cross-referencing
- **Admonitions** — Create note, tip, warning, and info callout blocks
- **Syntax Highlighting** — Automatic code highlighting for 200+ languages
- **SEO Ready** — Auto-generated meta tags, Open Graph, and Twitter Cards
- **Fast Dev Server** — Preview your site locally with live rendering
- **Clean URLs** — `docs/guide.md` becomes `/guide/` automatically

## Quick Example

```bash
# Install
go install github.com/wusher/volcano@latest

# Generate a site
volcano ./my-docs --title="My Documentation"

# Preview it
volcano -s ./output
```

Your markdown files:

```
my-docs/
├── index.md
├── getting-started.md
└── guides/
    ├── installation.md
    └── configuration.md
```

Become a fully navigable site with sidebar, breadcrumbs, and search.

## Get Started

Ready to build your first site? Head to the [[Getting Started]] guide.

## Documentation

- **[[guides/index|Guides]]** — Step-by-step tutorials for common tasks
- **[[features/index|Features]]** — Deep dives into Volcano's capabilities
- **[[reference/index|Reference]]** — CLI flags, URL structure, and technical details
- **[[examples/index|Examples]]** — Real-world site configurations

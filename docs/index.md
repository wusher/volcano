# Volcano

**Turn your markdown folders into websites. No config required.**

<style>
.logo-float {
  max-width: 260px;
}

@media (min-width: 768px) {
  .logo-float {
    float: right;
  }
}
</style>

<img src="logo.png" alt="Volcano Logo" class="logo-float">


---

Volcano is an opinionated static site generator for people who just want their markdown files on the web. Point it at a folder—your Obsidian vault, your notes directory, your documentation—and get a beautiful, navigable website.

No configuration files. No frontmatter. No build pipelines. Just markdown in, website out.

## Why Volcano?

**Your folder structure is your site structure.** Volcano reads your directories and creates matching navigation automatically. No need to define routes, menus, or page hierarchies in config files.

**Works with your existing notes.** Have an Obsidian vault? A folder of documentation? Meeting notes organized by date? Volcano handles them all without requiring you to restructure anything or add metadata.

**Single binary, no dependencies.** No Node.js, Python, or Ruby. Download one file and run it. Works offline, builds fast, deploys anywhere.

## Quick Start

```bash
# Install
go install github.com/wusher/volcano@latest

# Generate a site from your markdown folder
volcano ./my-notes --title="My Site"

# Preview it locally
volcano -s ./output
```

That's it. Your folder:

```
my-notes/
├── index.md
├── getting-started.md
└── guides/
    ├── installation.md
    └── configuration.md
```

Becomes a website with tree navigation, breadcrumbs, table of contents, and search—all generated from your folder structure.

## Features

- **Tree navigation** — Collapsible sidebar mirrors your folder structure
- **Clean URLs** — `guides/intro.md` → `/guides/intro/`
- **Table of contents** — Auto-generated from headings
- **Dark mode** — Automatic detection with toggle
- **Admonitions** — Note, tip, warning, and info callouts
- **Code highlighting** — Syntax highlighting with copy button
- **SEO ready** — Meta tags, Open Graph support
- **Keyboard shortcuts** — Press `?` for all shortcuts
- **Search** — Filter navigation by typing
- **Mobile responsive** — Works on any screen size

## Get Started

Ready to build your first site? Head to the [[Getting Started]] guide.

## Documentation

- **[[guides/index|Guides]]** — Step-by-step tutorials
- **[[features/index|Features]]** — Deep dives into capabilities
- **[[reference/index|Reference]]** — CLI flags and technical details
- **[[examples/index|Examples]]** — Real-world configurations

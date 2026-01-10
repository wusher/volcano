# Volcano Documentation Outline

A comprehensive documentation structure for the Volcano static site generator.

---

## Documentation Structure

```
docs/
├── index.md
├── getting-started.md
├── guides/
│   ├── index.md
│   ├── installation.md
│   ├── creating-your-first-site.md
│   ├── organizing-content.md
│   ├── customizing-appearance.md
│   ├── development-workflow.md
│   └── deploying-your-site.md
├── features/
│   ├── index.md
│   ├── markdown-syntax.md
│   ├── navigation.md
│   ├── wiki-links.md
│   ├── admonitions.md
│   ├── code-blocks.md
│   ├── seo-and-meta.md
│   ├── theming.md
│   ├── auto-index.md
│   └── reading-time.md
├── reference/
│   ├── index.md
│   ├── cli.md
│   ├── url-structure.md
│   └── front-matter.md
└── examples/
    ├── index.md
    ├── documentation-site.md
    ├── blog.md
    └── knowledge-base.md
```

---

## Detailed Outline

### Root

- **index.md** — Home / Overview
  - What is Volcano?
  - Key features at a glance
  - Why choose Volcano?
  - Quick feature comparison with other tools
  - Links to getting started

- **getting-started.md** — Quick Start (5 minutes)
  - Prerequisites (Go installed)
  - Install Volcano (one command)
  - Create a docs folder
  - Add your first markdown file
  - Generate the site
  - View in browser
  - Next steps

---

### Guides (How-To Tutorials)

- **guides/index.md** — Guides Overview
  - What you'll learn
  - Guide listing with descriptions

- **guides/installation.md** — Installation
  - Using `go install` (recommended)
  - Downloading pre-built binaries
  - Building from source
  - Verifying installation
  - Updating to latest version

- **guides/creating-your-first-site.md** — Creating Your First Site
  - Planning your content structure
  - Creating the source folder
  - Writing your first pages
  - Adding an index.md
  - Running the generator
  - Understanding the output
  - Previewing with the dev server

- **guides/organizing-content.md** — Organizing Content
  - Folder structure best practices
  - Using index.md files for sections
  - File naming conventions
    - Date prefixes (2024-01-15-)
    - Number prefixes (01-, 001-)
    - Readable filenames
  - How titles are extracted
  - Controlling sort order
  - Hidden files and drafts (_ prefix)

- **guides/customizing-appearance.md** — Customizing Appearance
  - Choosing a built-in theme
    - `docs` — Documentation sites
    - `blog` — Blog layouts
    - `vanilla` — Minimal starting point
  - Extracting CSS for customization
  - Creating a custom stylesheet
  - Adding a favicon
  - Setting site title and metadata

- **guides/development-workflow.md** — Development Workflow
  - Static server mode (pre-generated)
  - Dynamic server mode (live rendering)
  - When to use each mode
  - Watching for changes
  - Browser refresh workflow
  - Debugging with verbose mode

- **guides/deploying-your-site.md** — Deploying Your Site
  - Generating for production
  - Setting the base URL
  - GitHub Pages deployment
  - Netlify deployment
  - Vercel deployment
  - Self-hosting options
  - CI/CD automation examples

---

### Features (Capability Deep Dives)

- **features/index.md** — Features Overview
  - Feature categories
  - Quick feature matrix
  - Links to detailed docs

- **features/markdown-syntax.md** — Markdown Syntax
  - GitHub Flavored Markdown support
  - Headings and paragraphs
  - Emphasis (bold, italic, strikethrough)
  - Lists (ordered, unordered, task lists)
  - Links and images
  - Tables
  - Blockquotes
  - Footnotes
  - Definition lists
  - Smart typography (quotes, dashes)
  - Raw HTML support
  - Escaping special characters

- **features/navigation.md** — Navigation
  - Sidebar navigation tree
    - Hierarchical structure
    - Expand/collapse folders
    - Active page highlighting
    - Search filtering
  - Breadcrumbs
    - Schema.org markup
    - Automatic generation
  - Table of contents (TOC)
    - Heading extraction
    - Nested structure
    - Sidebar placement
    - Minimum heading threshold
  - Page navigation (prev/next)
    - Enabling with --page-nav
    - Sequential ordering
  - Top navigation bar
    - Enabling with --top-nav
    - Root-level items only
    - Item limits (1-8)

- **features/wiki-links.md** — Wiki-Style Links
  - Basic syntax: `[[Page Name]]`
  - Custom display text: `[[Page Name|Display Text]]`
  - Linking to paths: `[[folder/Page Name]]`
  - How links resolve
    - Sibling page lookup
    - Path-based lookup
  - Embed syntax: `![[Page Name]]`
  - Compatibility with Obsidian
  - Best practices

- **features/admonitions.md** — Admonitions (Callout Blocks)
  - Syntax: `:::type` ... `:::`
  - Available types
    - `:::note` — General information
    - `:::tip` — Helpful suggestions
    - `:::warning` — Important cautions
    - `:::danger` — Critical warnings
    - `:::info` — Informational callouts
  - Custom titles: `:::note Custom Title`
  - Styling and icons
  - Nesting content inside admonitions

- **features/code-blocks.md** — Code Blocks
  - Fenced code blocks
  - Language specification
  - Supported languages (Chroma)
  - Syntax highlighting
  - Line highlighting: ```lang {3,5-7}
  - Copy button functionality
  - Inline code styling
  - Code in admonitions

- **features/seo-and-meta.md** — SEO & Meta Tags
  - Automatic meta tags
    - Title and description
    - Author
    - Robots directives
  - Open Graph tags
    - og:title, og:description
    - og:image configuration
    - og:url and site_name
  - Twitter Card tags
    - Card types
    - Image handling
  - Canonical URLs
    - Setting base URL
    - Preventing duplicates
  - Description extraction
    - Automatic from content
    - Length and formatting

- **features/theming.md** — Theming
  - Built-in themes
    - `docs` — Full-featured documentation
    - `blog` — Blog-optimized layout
    - `vanilla` — Structural skeleton only
  - Theme selection: `--theme`
  - Extracting vanilla CSS: `volcano css`
  - Custom CSS: `--css path/to/style.css`
  - CSS class reference
    - Layout classes
    - Navigation classes
    - Content classes
    - Component classes
  - Dark mode support
  - Creating a custom theme

- **features/auto-index.md** — Auto-Generated Index Pages
  - How it works
  - When indexes are generated
  - Index page content
    - Child listing
    - Folder vs file distinction
  - Styling auto-indexes
  - Overriding with custom index.md

- **features/reading-time.md** — Reading Time
  - Enabling reading time display
  - Calculation method
    - 225 WPM for prose
    - 100 WPM for code blocks
  - Display format
  - Customizing appearance

---

### Reference (Technical Documentation)

- **reference/index.md** — Reference Overview
  - Available reference docs
  - Quick links

- **reference/cli.md** — CLI Reference
  - Commands
    - `volcano <input>` — Generate site
    - `volcano -s <folder>` — Serve site
    - `volcano css` — Output CSS
  - Generation flags
    - `-o, --output` — Output directory
    - `--title` — Site title
    - `--url` — Base URL
    - `--author` — Site author
    - `--og-image` — Default OG image
    - `--favicon` — Favicon path
    - `--theme` — Theme selection
    - `--css` — Custom CSS path
    - `--last-modified` — Show dates
    - `--top-nav` — Top navigation
    - `--page-nav` — Prev/next links
  - Server flags
    - `-s, --serve` — Enable server
    - `-p, --port` — Server port
  - Utility flags
    - `-q, --quiet` — Suppress output
    - `--verbose` — Debug output
    - `-v, --version` — Show version
    - `-h, --help` — Show help
  - Exit codes

- **reference/url-structure.md** — URL Structure
  - Clean URL generation
    - `file.md` → `/file/`
    - `folder/file.md` → `/folder/file/`
  - Index file handling
    - `index.md` → `/folder/`
    - `readme.md` → `/folder/`
  - Prefix stripping
    - Date prefixes: `2024-01-15-post.md` → `/post/`
    - Number prefixes: `01-intro.md` → `/intro/`
  - Slugification rules
    - Lowercase conversion
    - Space to hyphen
    - Special character removal
  - Trailing slashes

- **reference/front-matter.md** — Front Matter
  - YAML front matter syntax
  - Supported fields (if any)
  - How front matter is processed
  - Examples

---

### Examples (Real-World Use Cases)

- **examples/index.md** — Examples Overview
  - Available examples
  - Choosing the right setup

- **examples/documentation-site.md** — Documentation Site
  - Recommended structure
  - Example folder layout
  - Configuration flags
  - Sample content
  - Deployment tips

- **examples/blog.md** — Blog
  - Date-based organization
  - Using date prefixes
  - Blog theme setup
  - Post structure
  - Archive organization
  - RSS considerations

- **examples/knowledge-base.md** — Knowledge Base
  - Wiki-style organization
  - Using wiki links extensively
  - Cross-referencing pages
  - Flat vs hierarchical structure
  - Obsidian vault compatibility
  - Search optimization

---

## Content Guidelines

### Writing Style
- Clear, concise language
- Active voice
- Present tense
- Second person ("you")

### Page Structure
- Clear H1 title
- Brief introduction
- Logical sections with H2/H3
- Code examples where relevant
- Links to related pages

### Code Examples
- Use realistic examples
- Show both input and output
- Highlight key parts
- Include copy-friendly blocks

---

## Priority Order for Writing

1. **index.md** — First impression
2. **getting-started.md** — Critical for adoption
3. **guides/installation.md** — Unblocks users
4. **reference/cli.md** — Most referenced
5. **features/markdown-syntax.md** — Core functionality
6. **features/navigation.md** — Key differentiator
7. **guides/organizing-content.md** — Common questions
8. **features/theming.md** — Customization demand
9. Remaining guides
10. Remaining features
11. Examples
12. Reference docs

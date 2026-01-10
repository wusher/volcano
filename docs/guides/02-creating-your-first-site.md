# Creating Your First Site

Build a complete documentation site from scratch.

## Planning Your Content

Before writing, plan your site structure. A typical documentation site might look like:

```
docs/
├── index.md           # Homepage
├── getting-started.md # Quick start guide
├── guides/
│   ├── index.md      # Guides overview
│   ├── installation.md
│   └── configuration.md
└── api/
    ├── index.md      # API overview
    └── endpoints.md
```

## Creating the Source Folder

Create your documentation directory:

```bash
mkdir -p docs/guides docs/api
```

## Writing Your Homepage

Create `docs/index.md`:

```markdown
# My Project

Welcome to My Project documentation.

## Quick Links

- [[Getting Started]] — Get up and running
- [[guides/index|Guides]] — Step-by-step tutorials
- [[api/index|API Reference]] — Technical documentation

## Features

- Easy to use
- Well documented
- Actively maintained
```

The homepage introduces your project and links to key sections using [[features/wiki-links|wiki-style links]].

## Writing Content Pages

Create `docs/getting-started.md`:

```markdown
# Getting Started

Get My Project running in 5 minutes.

## Prerequisites

- Node.js 18+
- npm or yarn

## Installation

Install via npm:

\`\`\`bash
npm install my-project
\`\`\`

:::tip
Use the `--save` flag to add it to your package.json.
:::

## Next Steps

See the [[guides/installation|full installation guide]] for more options.
```

## Creating Section Index Pages

For each folder, create an `index.md` to serve as the section landing page.

Create `docs/guides/index.md`:

```markdown
# Guides

Step-by-step tutorials for common tasks.

- [[installation]] — Install My Project
- [[configuration]] — Configure settings
```

## Running the Generator

Generate your site:

```bash
volcano ./docs -o ./public --title="My Project"
```

Output:

```
Generating site...
  Input:  ./docs
  Output: ./public
  Title:  My Project

Scanning input directory...
Found 6 markdown files in 3 folders

Generating pages...
  ✓ api/endpoints.md
  ✓ api/index.md
  ✓ getting-started.md
  ✓ guides/configuration.md
  ✓ guides/index.md
  ✓ guides/installation.md
  ✓ index.md

Generated 7 pages in ./public
```

## Understanding the Output

Volcano creates a clean URL structure:

```
public/
├── index.html           # /
├── getting-started/
│   └── index.html       # /getting-started/
├── guides/
│   ├── index.html       # /guides/
│   ├── installation/
│   │   └── index.html   # /guides/installation/
│   └── configuration/
│       └── index.html   # /guides/configuration/
└── api/
    ├── index.html       # /api/
    └── endpoints/
        └── index.html   # /api/endpoints/
```

Each page gets:
- Sidebar navigation with your folder structure
- Breadcrumbs showing the current location
- Table of contents (if the page has 3+ headings)
- Clean URLs without `.html` extensions

## Previewing Your Site

Start the development server:

```bash
volcano -s ./public
```

Open [http://localhost:1776](http://localhost:1776) in your browser.

Navigate around to see:
- The sidebar expands to show your current section
- Breadcrumbs update as you navigate
- Wiki links connect your pages together

## Making Changes

Edit your markdown files and regenerate:

```bash
volcano ./docs -o ./public --title="My Project"
```

Or use dynamic mode to see changes instantly:

```bash
volcano -s ./docs
```

In dynamic mode, Volcano renders pages on-the-fly from your source files.

## Next Steps

- [[organizing-content]] — Learn naming conventions and best practices
- [[customizing-appearance]] — Add your own branding
- [[deploying-your-site]] — Publish your site online

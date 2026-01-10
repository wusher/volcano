# Organizing Content

Structure your documentation for clarity and easy navigation.

## Folder Structure

Volcano mirrors your folder structure in the navigation sidebar. Organize your content hierarchically:

```
docs/
├── index.md              # Site homepage (/)
├── getting-started.md    # Top-level page (/getting-started/)
├── guides/
│   ├── index.md         # Section index (/guides/)
│   ├── basics.md        # Section page (/guides/basics/)
│   └── advanced/
│       ├── index.md     # Subsection index (/guides/advanced/)
│       └── topics.md    # Nested page (/guides/advanced/topics/)
└── reference/
    └── api.md           # Another section (/reference/api/)
```

**Best practices:**
- Keep nesting to 2-3 levels maximum
- Use descriptive folder names
- Always include an `index.md` in folders with multiple pages

## Index Files

An `index.md` (or `readme.md`) in a folder serves as that section's landing page.

Without an index file, Volcano auto-generates a simple listing of the folder's contents.

With an index file, you control what users see:

```markdown
# API Reference

Complete API documentation for My Project.

## Endpoints

- [[authentication]] — Auth flows and tokens
- [[users]] — User management
- [[resources]] — Resource CRUD operations
```

## File Naming Conventions

### Basic Names

Use lowercase with hyphens:

```
getting-started.md    # Good
Getting Started.md    # Works but creates messier URLs
GettingStarted.md     # Works but less readable
```

### Date Prefixes

Prefix files with dates for chronological ordering (useful for blogs):

```
2024-01-15-my-post.md
2024-01-20-another-post.md
2024-02-01-latest-post.md
```

The date prefix is stripped from URLs:
- `2024-01-15-my-post.md` → `/my-post/`

But files are sorted by the date, newest first.

### Number Prefixes

Prefix with numbers to control sort order:

```
01-introduction.md
02-installation.md
03-configuration.md
10-advanced-topics.md
```

Number prefixes are stripped from URLs:
- `01-introduction.md` → `/introduction/`

Files are sorted numerically:
- `01-` comes before `02-`
- `02-` comes before `10-`

### Combining Prefixes

You can use both numbers and dates:

```
guides/
├── 01-getting-started.md
├── 02-basics.md
posts/
├── 2024-01-15-welcome.md
├── 2024-01-20-update.md
```

## How Titles Are Extracted

Volcano determines page titles in this order:

1. **First H1 heading** in the markdown file:
   ```markdown
   # My Page Title

   Content here...
   ```

2. **Cleaned filename** if no H1 is found:
   - `getting-started.md` → "Getting Started"
   - `01-introduction.md` → "Introduction"
   - `2024-01-15-my-post.md` → "My Post"

**Best practice:** Always include an H1 heading in your files for explicit control over titles.

## Controlling Sort Order

Pages in the sidebar are sorted as follows:

1. **Files first**, then folders
2. **Within each group**: by date prefix (newest first), then number prefix, then alphabetically

Example:

```
guides/
├── 2024-01-20-update.md    # 1st (has date, newest first)
├── 2024-01-15-news.md      # 2nd
├── 01-intro.md             # 3rd (number prefix)
├── 02-basics.md            # 4th
├── advanced.md             # 5th (alphabetical)
└── zz-appendix.md          # 6th
```

## Hidden Files and Drafts

Files and folders starting with `_` or `.` are ignored:

```
docs/
├── _drafts/              # Ignored folder
│   └── upcoming.md       # Not published
├── .hidden/              # Ignored folder
├── index.md              # Published
└── _work-in-progress.md  # Ignored file
```

Use this for:
- Draft content not ready for publication
- Template files
- Notes and scratch files

## Linking Between Pages

Use wiki links for easy cross-referencing:

```markdown
See the [[Installation]] guide.
Check out [[guides/advanced/topics|Advanced Topics]].
```

Or standard markdown links:

```markdown
See the [Installation](/guides/installation/) guide.
```

Wiki links are simpler and automatically resolve to the correct URL.

## Example: Documentation Site

```
docs/
├── index.md
├── getting-started.md
├── guides/
│   ├── index.md
│   ├── 01-installation.md
│   ├── 02-configuration.md
│   └── 03-deployment.md
├── reference/
│   ├── index.md
│   ├── cli.md
│   └── api.md
└── about.md
```

## Example: Blog

```
blog/
├── index.md
├── posts/
│   ├── 2024-01-15-welcome.md
│   ├── 2024-02-01-new-feature.md
│   └── 2024-03-10-update.md
└── about.md
```

## Next Steps

- [[customizing-appearance]] — Style your site
- [[features/wiki-links]] — Learn wiki link syntax
- [[reference/url-structure]] — Understand URL generation

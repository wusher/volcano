# Front Matter

Volcano supports YAML front matter in markdown files.

## Format

Front matter is YAML content at the beginning of a file, delimited by `---`:

```markdown
---
title: My Page Title
author: Jane Smith
---

# Page Content

Your markdown content here...
```

## How It's Processed

Currently, Volcano **strips front matter** from content before rendering. The front matter block is removed so it doesn't appear in the output HTML.

### Before Processing

```markdown
---
title: Getting Started
date: 2024-01-15
tags: [tutorial, beginner]
---

# Getting Started

Welcome to the tutorial...
```

### After Processing

The markdown rendered is:

```markdown
# Getting Started

Welcome to the tutorial...
```

## Current Behavior

Volcano's front matter support is focused on compatibility:

1. **Stripping** — Front matter is cleanly removed from content
2. **Obsidian Compatible** — Works with Obsidian-style front matter
3. **No Errors** — Front matter doesn't cause rendering issues

### What Works

- Standard YAML front matter syntax
- Any valid YAML content
- Unix (`\n`) and Windows (`\r\n`) line endings
- Front matter at the very start of the file

### Requirements

For front matter to be recognized:

1. File must start with `---` (no leading whitespace)
2. Must have a closing `---` on its own line
3. Must be valid YAML (though content is not parsed)

## Examples

### Minimal

```markdown
---
---

Content starts here.
```

### Typical Usage

```markdown
---
title: Installation Guide
description: How to install the software
author: Documentation Team
---

# Installation Guide

Follow these steps...
```

### Complex YAML

```markdown
---
title: API Reference
metadata:
  version: 2.0
  status: stable
tags:
  - api
  - reference
  - v2
---

# API Reference

Documentation content...
```

## Page Titles

Volcano determines page titles from:

1. **H1 Heading** — First `# Heading` in the content (preferred)
2. **Filename** — Cleaned and title-cased if no H1

Front matter `title` fields are not currently used for the page title.

### Recommendation

Use an H1 heading for your page title:

```markdown
---
description: Optional metadata
---

# Your Page Title

This H1 becomes the page title in navigation and browser tabs.
```

## Compatibility

### Obsidian

Volcano's front matter handling is compatible with Obsidian:

```markdown
---
aliases: [home, index]
tags: [documentation]
cssclass: wide-page
---

# Welcome

Content here...
```

The front matter is stripped, and the content renders normally.

### Other Tools

Front matter from other tools works:

**Hugo:**
```markdown
---
title: "My Post"
date: 2024-01-15T10:00:00Z
draft: false
---
```

**Jekyll:**
```markdown
---
layout: post
title: My Post
categories: [blog]
---
```

**Docusaurus:**
```markdown
---
id: intro
title: Introduction
sidebar_position: 1
---
```

All stripped cleanly, content renders correctly.

## Limitations

Current limitations of front matter support:

1. **Not Parsed** — YAML values aren't extracted or used
2. **No Custom Titles** — Front matter `title` doesn't set page title
3. **No Layout Control** — Can't specify templates
4. **No Custom Metadata** — Values aren't available in templates

### Workarounds

**For page titles:** Use an H1 heading
```markdown
# My Custom Title
```

**For descriptions:** Write good opening paragraphs (auto-extracted for SEO)

**For dates:** Use filename prefixes
```
2024-01-15-my-post.md
```

**For ordering:** Use number prefixes
```
01-introduction.md
02-getting-started.md
```

## Future Considerations

Potential future front matter features:

- Title override from front matter
- Custom descriptions for SEO
- Draft status (`draft: true`)
- Custom sidebar labels
- Page-specific CSS classes

## Related

- [[url-structure]] — How filenames affect URLs
- [[features/seo-and-meta]] — Automatic meta tag generation
- [[guides/organizing-content]] — File organization patterns

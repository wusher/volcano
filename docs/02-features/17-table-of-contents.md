# Table of Contents

Automatic table of contents sidebar for easy within-page navigation.

## Automatic Generation

**No configuration needed.** Tables of contents are automatically generated for pages with 3 or more headings.

## What You Get

A sidebar showing the page's heading structure:

```
On this page
├── Introduction
├── Installation
│   ├── Requirements
│   └── Steps
├── Configuration
│   ├── Basic Setup
│   └── Advanced Options
└── Troubleshooting
```

## Features

### Auto-Generated

TOC is built from your markdown headings:
- **H2** (`##`) - Main sections
- **H3** (`###`) - Subsections
- **H4** (`####`) - Sub-subsections

H1 is not included (it's the page title).

### Click to Jump

Click any heading in the TOC to scroll to that section. Smooth scrolling animation brings it into view.

### Active Tracking

As you scroll, the TOC highlights which section you're currently reading:

```
On this page
├── Introduction
├── Installation         ← highlighted (you're here)
│   ├── Requirements
│   └── Steps
└── Configuration
```

### Nested Structure

The TOC preserves your heading hierarchy:

```markdown
## Getting Started        ← Top level
### Install              ← Nested under "Getting Started"
### Configure            ← Nested under "Getting Started"
## Advanced Usage         ← Top level again
### Optimization          ← Nested under "Advanced Usage"
```

## Minimum Headings

TOC only appears if a page has **3 or more headings**. This prevents clutter on short pages.

Single-section pages (0-2 headings) don't show a TOC.

## Desktop Only

TOC appears on the right side on larger screens (≥1280px width).

On smaller screens:
- Desktop (1280px+): TOC sidebar visible
- Tablet (768-1279px): TOC hidden (limited space)
- Mobile (<768px): Mobile TOC toggle button in header

## Mobile TOC

On mobile devices, tap the TOC icon in the header to open a slide-out TOC panel.

## Heading Requirements

For the TOC to work well:

### Use Proper Hierarchy

```markdown
## Main Section
### Subsection
#### Detail

## Another Main Section
```

### Keep Text Concise

```markdown
## Installation          ✓ Good
## How to Install This Software on Your Computer  ✗ Too long
```

### Use Unique Headings

Duplicate headings get numbered IDs:

```markdown
## Configuration        → #configuration
## Configuration        → #configuration-1
## Configuration        → #configuration-2
```

## Styling

The TOC:
- Matches your theme
- Shows active heading with accent color
- Indents nested levels
- Has hover effects

## Scroll Behavior

When clicking a TOC link:
1. Page scrolls smoothly to the heading
2. Heading appears near the top (with offset for header)
3. URL updates with the heading anchor (`#section-name`)
4. Active state updates in the TOC

## URL Anchors

Every heading gets an anchor:

```markdown
## Getting Started
```

Creates anchor: `#getting-started`

You can link directly:
```markdown
See [[page#getting-started|Getting Started section]]
```

## Accessibility

The TOC is fully accessible:
- Semantic `<nav>` element
- Proper heading hierarchy
- Keyboard navigable (Tab, Enter)
- `aria-label="Table of contents"`
- `aria-current` on active item

## Layout

With TOC:
```
┌─────────┬───────────────┬──────┐
│ Sidebar │ Page Content  │ TOC  │
└─────────┴───────────────┴──────┘
```

Without TOC:
```
┌─────────┬─────────────────────┐
│ Sidebar │ Page Content        │
└─────────┴─────────────────────┘
```

## When TOC Appears

```markdown
# Page Title

## Section One           ← Heading 1
## Section Two           ← Heading 2
## Section Three         ← Heading 3

✓ TOC appears (3 headings)
```

```markdown
# Page Title

## Only Section

Just one section here.

✗ No TOC (only 1 heading)
```

## Best Practices

### Structure Your Content

Use headings to organize:
```markdown
## Overview
## Installation
### Requirements
### Steps
## Configuration
### Basic
### Advanced
```

### Balance Depth

Too shallow:
```markdown
## Installation
## Configuration
## Deployment
## Troubleshooting
## FAQ
## Support
```

Too deep:
```markdown
## Installation
### Step 1
#### Substep A
##### Detail 1
###### Note
```

Ideal:
```markdown
## Installation
### Requirements
### Installation Steps
## Configuration
### Basic Setup
### Advanced Options
```

## Related

- [[navigation]] — Navigation overview
- [[breadcrumbs]] — Page hierarchy
- [[search]] — Find content across pages

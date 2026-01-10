# Wiki-Style Links

Volcano supports Obsidian-style wiki links for easy cross-referencing between pages.

## Basic Syntax

Link to another page using double brackets:

```markdown
See the [[Installation]] guide.
Check out [[Getting Started]].
```

The link text is the page name, and Volcano resolves it to the correct URL.

## Custom Display Text

Use a pipe to set different display text:

```markdown
Read the [[Installation|install guide]] first.
See [[Getting Started|how to get started]].
```

Renders as:
- Read the [install guide](/guides/installation/) first.
- See [how to get started](/getting-started/).

## Linking to Paths

For pages in subdirectories, include the path:

```markdown
See [[guides/installation]].
Check [[reference/cli|the CLI reference]].
```

Paths are relative to the site root (not the current page).

## How Links Resolve

### Simple Links

A simple link like `[[Page Name]]` is resolved relative to the current page's directory:

| Current Page | Link | Resolves To |
|--------------|------|-------------|
| `/guides/` | `[[Installation]]` | `/guides/installation/` |
| `/reference/` | `[[CLI]]` | `/reference/cli/` |
| `/` | `[[About]]` | `/about/` |

### Path Links

A link with a path like `[[folder/page]]` resolves from the site root:

| Current Page | Link | Resolves To |
|--------------|------|-------------|
| `/` | `[[guides/installation]]` | `/guides/installation/` |
| `/reference/` | `[[guides/installation]]` | `/guides/installation/` |

### Index Files

Links to `index` or `readme` resolve to the folder:

```markdown
[[guides/index]]  → /guides/
[[readme]]        → /  (if at root)
```

## URL Generation

Link targets are converted to URL-friendly slugs:

| Wiki Link | URL |
|-----------|-----|
| `[[Getting Started]]` | `/getting-started/` |
| `[[API Reference]]` | `/api-reference/` |
| `[[guides/My Guide]]` | `/guides/my-guide/` |

Transformations:
- Spaces → hyphens
- Uppercase → lowercase
- `.md` extension removed if present

## Embeds

Obsidian embed syntax is converted to regular links:

```markdown
![[Page Name]]
```

Becomes a regular link to that page. Full embed functionality (showing page content inline) is not supported.

## Examples

### Linking Within a Section

In `/guides/installation.md`:

```markdown
After installing, see [[Configuration]] for setup options.
For troubleshooting, check [[FAQ]].
```

Both links resolve relative to `/guides/`.

### Linking Across Sections

In `/guides/installation.md`:

```markdown
See the [[reference/cli|CLI reference]] for all flags.
Return to the [[index|homepage]].
```

### Building a Navigation Page

```markdown
# User Guide

Welcome to the user guide.

## Getting Started

- [[Getting Started]] — Quick start tutorial
- [[Installation]] — Install the software

## Advanced Topics

- [[guides/advanced/performance|Performance Tuning]]
- [[guides/advanced/security|Security Best Practices]]

## Reference

- [[reference/cli|Command Line Reference]]
- [[reference/api|API Documentation]]
```

## Compatibility with Obsidian

Volcano's wiki links are designed to be compatible with Obsidian:

**Supported:**
- `[[Page Name]]`
- `[[Page Name|Display Text]]`
- `[[folder/Page Name]]`
- `![[Page Name]]` (converted to link)

**Not Supported:**
- `[[Page Name#Heading]]` (heading anchors)
- `[[Page Name#^block]]` (block references)
- Transclusion (embedding page content)

## Best Practices

### Use Descriptive Names

```markdown
# Good
See [[Installation Guide]] for setup instructions.

# Less Clear
See [[here]] for setup instructions.
```

### Be Consistent

Pick a naming convention and stick with it:

```markdown
# Consistent
[[Getting Started]]
[[Installation Guide]]
[[Configuration Options]]

# Inconsistent
[[Getting Started]]
[[install]]
[[ConfigurationOptions]]
```

### Use Paths for Clarity

When linking across sections, include the path:

```markdown
# Clear
See [[reference/cli]] for all options.

# Might Be Ambiguous (if multiple "cli" pages exist)
See [[cli]] for all options.
```

## Troubleshooting

### Link Not Working

If a wiki link doesn't resolve:

1. **Check the page exists** — The target page must exist
2. **Check the path** — Use paths for cross-section links
3. **Check spelling** — Names are case-insensitive but must match

### Wrong Page Linked

If linking to the wrong page:

1. **Use explicit paths** — `[[guides/setup]]` instead of `[[setup]]`
2. **Check for duplicates** — Same filename in different folders

## Related

- [[markdown-syntax]] — All markdown features
- [[guides/organizing-content]] — File and folder naming
- [[reference/url-structure]] — URL generation rules

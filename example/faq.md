# Frequently Asked Questions

## General

### What is Volcano?

Volcano is a static site generator written in Go. It converts markdown files into a styled, navigable website with automatic tree navigation.

### Why should I use Volcano?

- **Simple** - No configuration files needed
- **Fast** - Generates sites in milliseconds
- **Portable** - Single binary, no dependencies
- **Beautiful** - Built-in styling with dark mode

## Technical

### What markdown features are supported?

Volcano supports standard markdown plus:

- **Tables**
- **Syntax highlighting** for code blocks
- **Strikethrough** text
- Smart quotes and typography

### How do I customize the appearance?

Currently, Volcano uses a built-in stylesheet. Custom themes are planned for a future release.

> **Note:** You can use CSS overrides by including a `<style>` tag in your markdown.

### Can I use frontmatter?

Frontmatter support is planned for a future release. For now, the page title is extracted from the first `# Heading` in your document.

## Troubleshooting

### My pages aren't showing up

Make sure your markdown files have the `.md` extension. Volcano ignores files with other extensions.

### The navigation order is wrong

Volcano sorts pages alphabetically by filename. Use numbered prefixes to control order:

```
01-introduction.md
02-getting-started.md
03-advanced-topics.md
```

The numbers will be stripped from the display labels.

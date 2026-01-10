# Reference

Technical reference documentation for Volcano.

## Quick Links

### [[cli|Command Line Interface]]

Complete reference for all CLI flags and commands:

- Generate and serve commands
- Output and theming options
- SEO and metadata flags
- The `css` subcommand

### [[url-structure|URL Structure]]

How Volcano generates URLs:

- Clean URL format
- Slugification rules
- Date and number prefix stripping
- Directory path handling

### [[front-matter|Front Matter]]

YAML front matter handling:

- Supported format
- How it's processed
- Current limitations

## Command Quick Reference

### Generate Site

```bash
volcano <input-folder> [flags]
```

### Serve Site

```bash
volcano -s <folder> [flags]
```

### Export CSS

```bash
volcano css [-o file]
```

## Flag Summary

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output` | Output directory | `./output` |
| `-s, --serve` | Run in serve mode | `false` |
| `-p, --port` | Server port | `1776` |
| `--title` | Site title | `My Site` |
| `--url` | Site base URL | (none) |
| `--theme` | Theme name | `docs` |
| `--css` | Custom CSS path | (none) |

See [[cli]] for the complete list.

## URL Examples

| Input File | Output URL |
|------------|------------|
| `docs/intro.md` | `/intro/` |
| `guides/setup.md` | `/guides/setup/` |
| `2024-01-15-post.md` | `/post/` |
| `01-getting-started.md` | `/getting-started/` |
| `index.md` | `/` |

See [[url-structure]] for complete rules.

# SEO & Meta Tags

Volcano automatically generates SEO-friendly meta tags for all pages.

## Automatic Meta Tags

Every page includes essential meta tags:

```html
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="robots" content="index, follow">
```

## Title Tag

Page titles follow the format: `Page Title - Site Title`

```html
<title>Installation - My Documentation</title>
```

Set the site title with `--title`:

```bash
volcano ./docs --title="My Documentation"
```

## Description

Volcano extracts a description from the first ~160 characters of page content:

```html
<meta name="description" content="Learn how to install and configure...">
```

The description is automatically:
- Extracted from the beginning of content
- Stripped of HTML tags
- Truncated at word boundaries
- Limited to ~160 characters

## Author

Add an author meta tag with `--author`:

```bash
volcano ./docs --author="Jane Smith"
```

```html
<meta name="author" content="Jane Smith">
```

## Canonical URLs

Set the base URL with `--url` for canonical link tags:

```bash
volcano ./docs --url="https://docs.example.com"
```

Each page gets a canonical URL:

```html
<link rel="canonical" href="https://docs.example.com/guides/installation/">
```

Canonical URLs:
- Help search engines identify the primary URL
- Prevent duplicate content issues
- Are essential for SEO

## Open Graph Tags

Open Graph tags enable rich previews when shared on social media:

```html
<meta property="og:title" content="Installation">
<meta property="og:description" content="Learn how to install...">
<meta property="og:type" content="website">
<meta property="og:url" content="https://docs.example.com/guides/installation/">
<meta property="og:site_name" content="My Documentation">
```

### Open Graph Image

Set a default OG image with `--og-image`:

```bash
volcano ./docs --og-image="https://example.com/og-image.png"
```

```html
<meta property="og:image" content="https://example.com/og-image.png">
```

**Image recommendations:**
- Minimum size: 1200 x 630 pixels
- Aspect ratio: 1.91:1
- Format: PNG or JPEG
- Host on a reliable CDN

## Twitter Card Tags

Twitter Cards enable rich previews on Twitter/X:

```html
<meta name="twitter:card" content="summary_large_image">
<meta name="twitter:title" content="Installation">
<meta name="twitter:description" content="Learn how to install...">
<meta name="twitter:image" content="https://example.com/og-image.png">
```

The card type is automatically set:
- `summary_large_image` — When an OG image is set
- `summary` — When no image is provided

## Complete Example

Generate a fully SEO-optimized site:

```bash
volcano ./docs \
  -o ./public \
  --title="My Project Documentation" \
  --url="https://docs.myproject.com" \
  --author="My Team" \
  --og-image="https://docs.myproject.com/og-image.png"
```

Generated meta tags:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>Installation - My Project Documentation</title>

    <meta name="description" content="Step-by-step guide to installing My Project on your system...">
    <meta name="author" content="My Team">
    <meta name="robots" content="index, follow">

    <link rel="canonical" href="https://docs.myproject.com/guides/installation/">

    <meta property="og:title" content="Installation">
    <meta property="og:description" content="Step-by-step guide to installing My Project...">
    <meta property="og:type" content="website">
    <meta property="og:url" content="https://docs.myproject.com/guides/installation/">
    <meta property="og:site_name" content="My Project Documentation">
    <meta property="og:image" content="https://docs.myproject.com/og-image.png">

    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:title" content="Installation">
    <meta name="twitter:description" content="Step-by-step guide to installing My Project...">
    <meta name="twitter:image" content="https://docs.myproject.com/og-image.png">
</head>
```

## Best Practices

### Always Set Base URL

The `--url` flag is crucial for production sites:

```bash
# Development (no URL needed)
volcano ./docs -o ./output

# Production (always set URL)
volcano ./docs -o ./output --url="https://docs.example.com"
```

### Write Good Introductions

Since descriptions are extracted from content, make your page introductions count:

```markdown
# Installation

Get My Project running on your system in under 5 minutes. This guide covers
installation on Windows, macOS, and Linux.

## Prerequisites
...
```

The first paragraph becomes the meta description.

### Use Meaningful Page Titles

Page titles come from H1 headings:

```markdown
<!-- Good: Specific and descriptive -->
# Installing My Project on Ubuntu

<!-- Less Good: Vague -->
# Installation
```

### Provide a Good OG Image

A compelling image increases click-through rates:

- Include your project logo or name
- Use readable text (it will be shown small)
- Keep important content in the center
- Test with [Facebook Sharing Debugger](https://developers.facebook.com/tools/debug/)

## Testing SEO

### Validate Meta Tags

Use these tools to check your meta tags:

- [Google Rich Results Test](https://search.google.com/test/rich-results)
- [Facebook Sharing Debugger](https://developers.facebook.com/tools/debug/)
- [Twitter Card Validator](https://cards-dev.twitter.com/validator)

### Check Generated HTML

Inspect the generated HTML directly:

```bash
# Generate and check
volcano ./docs -o ./output
cat ./output/guides/installation/index.html | head -50
```

## Related

- [[guides/deploying-your-site]] — Deploy with proper URLs
- [[reference/cli]] — All SEO-related flags

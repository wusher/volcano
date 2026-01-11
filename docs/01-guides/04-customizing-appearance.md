# Customizing Appearance

Customize your site with themes, CSS, and branding options.

## Themes

Choose from three built-in themes:

```bash
volcano ./docs --theme docs    # Default: documentation
volcano ./docs --theme blog    # Reading-focused
volcano ./docs --theme vanilla # Unstyled skeleton
```

For detailed theming information, see [[features/theming]].

## Accent Color

Set a custom accent color:

```bash
volcano ./docs --accent-color="#0284c7"
```

## Custom CSS

Export and customize the base theme:

```bash
# Export base CSS
volcano css -o custom.css

# Use your custom CSS
volcano ./docs --css ./custom.css
```

## Branding

```bash
# Add favicon
volcano ./docs --favicon ./favicon.png

# Set site title
volcano ./docs --title="My Project"

# Set author for meta tags
volcano ./docs --author="Jane Smith"

# Base URL for SEO
volcano ./docs --url="https://docs.example.com"
```

## Display Options

```bash
# Show last modified dates
volcano ./docs --last-modified

# Top navigation bar
volcano ./docs --top-nav

# Previous/next page links
volcano ./docs --page-nav
```

## Example Build

Combine multiple options:

```bash
volcano ./docs \
  -o ./public \
  --title="My Project" \
  --accent-color="#0284c7" \
  --theme docs \
  --favicon="./favicon.png" \
  --top-nav \
  --page-nav
```

## Next Steps

- [[features/theming]] — Deep dive into theming and CSS
- [[development-workflow]] — Preview changes with dev server
- [[deploying-your-site]] — Publish your site

# Customizing Appearance

Customize your site with themes, CSS, and branding options.

## Themes

Choose from four built-in themes:

```bash
volcano ./docs --theme docs         # Default: documentation
volcano ./docs --theme blog         # Reading-focused
volcano ./docs --theme presentation # Slide-deck / talk style
volcano ./docs --theme vanilla      # Unstyled skeleton
```

### docs

| Light | Dark |
|-------|------|
| ![docs theme, light mode](/images/themes/docs-light.png) | ![docs theme, dark mode](/images/themes/docs-dark.png) |

### blog

| Light | Dark |
|-------|------|
| ![blog theme, light mode](/images/themes/blog-light.png) | ![blog theme, dark mode](/images/themes/blog-dark.png) |

### presentation

| Light | Dark |
|-------|------|
| ![presentation theme, light mode](/images/themes/presentation-light.png) | ![presentation theme, dark mode](/images/themes/presentation-dark.png) |

### vanilla

| Light | Dark |
|-------|------|
| ![vanilla theme, light mode](/images/themes/vanilla-light.png) | ![vanilla theme, dark mode](/images/themes/vanilla-dark.png) |

For detailed theming information, see [[theming]].

## Accent Color

Set an accent color by Tailwind name (uses the `500` shade) or any hex value.
The default is `sky`:

```bash
# Tailwind color names
volcano ./docs --accent-color sky
volcano ./docs --accent-color rose
volcano ./docs --accent-color emerald
volcano ./docs --accent-color teal

# Hex values still work
volcano ./docs --accent-color="#0284c7"
```

Supported Tailwind names: `slate`, `gray`, `zinc`, `neutral`, `stone`, `red`,
`orange`, `amber`, `yellow`, `lime`, `green`, `emerald`, `teal`, `cyan`, `sky`,
`blue`, `indigo`, `violet`, `purple`, `fuchsia`, `pink`, `rose`.

## Custom CSS

Export and customize the base theme:

```bash
# Export base CSS
volcano css -o custom.css

# Use your custom CSS
volcano ./docs --css ./custom.css
```

:::tip Generate with AI
You can use AI tools like Claude to generate custom CSS themes. Just describe your desired look and get production-ready CSS that works with Volcano's vanilla theme.
:::

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
  --accent-color sky \
  --theme docs \
  --favicon="./favicon.png" \
  --top-nav \
  --page-nav
```

## Next Steps

- [[theming]] — Deep dive into theming and custom CSS
- [[building-your-site]] — Build with your custom theme
- [[serving-your-site]] — Preview your site locally

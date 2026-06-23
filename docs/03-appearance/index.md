# Appearance

Themes, accent colors, and branding.

## Themes

Five themes ship in. Pick one with `--theme`:

```bash
volcano ./docs --theme docs --url="https://example.com"
```

| Theme | Best for |
|-------|----------|
| **[[docs|docs]]** (default) | Documentation, API references, technical guides |
| **[[blog|blog]]** | Long-form articles, posts, changelogs |
| **[[presentation|presentation]]** | Talks, demos, narrative writeups |
| **[[readable|readable]]** | Dyslexia-friendly reading (OpenDyslexic font) |
| **[[vanilla|vanilla]]** | Starting point for fully custom themes |

Each theme page has full-size light + dark screenshots — click any screenshot to zoom.

## Accent Color

Every theme exposes one accent color. Set it with `--accent-color`:

```bash
# A Tailwind palette name (uses the 500 shade)
volcano ./docs --accent-color sky
volcano ./docs --accent-color rose
volcano ./docs --accent-color emerald

# A hex value
volcano ./docs --accent-color "#0ea5e9"
```

Supported names: `slate`, `gray`, `zinc`, `neutral`, `stone`, `red`, `orange`, `amber`, `yellow`, `lime`, `green`, `emerald`, `teal`, `cyan`, `sky` (default), `blue`, `indigo`, `violet`, `purple`, `fuchsia`, `pink`, `rose`.

### Gradient Accents

Pass two colors separated by `-` to get a gradient:

```bash
volcano ./docs --accent-color lime-sky
volcano ./docs --accent-color "#444444-#555555"
volcano ./docs --accent-color "lime-#0ea5e9"
```

Volcano emits four CSS custom properties when a gradient is set:

```css
:root {
  --accent: <start color>;
  --accent-end: <end color>;
  --accent-gradient: linear-gradient(to right, <start>, <end>);
  --accent-gradient-vertical: linear-gradient(to bottom, <start>, <end>);
}
```

The horizontal gradient automatically paints the scroll progress bar, page H1, in-content links, page-nav, and (in the docs theme) H2 underlines. The vertical variant paints admonition and blockquote left borders.

## Branding

```bash
volcano ./docs \
  --title="My Project" \
  --favicon=./icon.png \
  --author="Your Name" \
  --og-image="https://docs.example.com/og.png" \
  --url="https://docs.example.com"
```

## Next

- **[[custom-css|Custom CSS]]** — go beyond themes
- **[CLI reference](/cli/)** — every appearance flag

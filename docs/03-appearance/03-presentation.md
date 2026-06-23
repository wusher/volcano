# presentation Theme

Slide-deck typography for talks and demos. Oversized fluid headings, generous whitespace, chrome that recedes.

```bash
volcano ./talk --theme presentation --url="https://example.com"
```

## Light

![presentation theme, light mode](/images/themes/presentation-light.png)

## Dark

![presentation theme, dark mode](/images/themes/presentation-dark.png)

## Features

- Oversized fluid H1 — reads like a title slide
- Each `##` heading feels like a new slide (huge top margin)
- Pull-quote style blockquotes
- High-contrast palette tuned for projector legibility
- Dark mode optimized for stage lighting
- Sidebar fades to transparent until you hover — chrome disappears, content dominates

## Keyboard Shortcuts

When this theme is active:

| Key | Action |
|-----|--------|
| `←` | Previous heading |
| `→` | Next heading |

Press `?` to see the full shortcuts list — the heading-nav row only appears on this theme.

## Best For

Conference talks, product demos, internal showcases, narrative writeups. Anything you'd advance with arrow keys in front of an audience.

## Tip: Pair With Top Nav

```bash
volcano ./talk \
  --theme presentation \
  --top-nav \
  --url="https://talks.example.com"
```

Top nav gives a stable "deck index" along the top while the sidebar stays out of the way.

# Configuration

Every setting Volcano accepts, and where it can live.

## Where Settings Come From

Volcano resolves each setting from three sources, in order:

1. **CLI flag** — `--title="My Site"` on the command line
2. **`volcano.json`** — config file in the input directory (or `-c path`)
3. **Default** — the built-in default

CLI flags always win. When a CLI flag overrides a config-file value, Volcano prints a notice:

```
CLI --title overrides config file
```

This means: keep durable settings (theme, title, URL) in `volcano.json` and reach for the CLI only when you need to override for a one-off build.

## Creating a Config File

`volcano init` writes a `volcano.json` with every option at its default:

```bash
volcano init              # creates ./volcano.json
volcano init -o ./docs    # creates ./docs/volcano.json
```

If a config file already exists, `init` adds any missing keys but preserves your existing values.

After that, every command uses it automatically:

```bash
volcano serve ./docs           # picks up ./docs/volcano.json
volcano ./docs                 # same — config provides --url and the rest
```

Or point at a specific file:

```bash
volcano ./docs -c ./prod.json
```

## All Settings

Every setting Volcano accepts. Click a "Feature" link to learn what each one does in context.

### Build basics

| CLI flag | JSON key | Default | What it does |
|----------|----------|---------|--------------|
| `--url` | `"url"` | `""` | **Required** for builds. Base URL for canonical / Open Graph tags. |
| `-o`, `--output` | `"output"` | `"./output"` | Where to write the static site |
| `-p`, `--port` | `"port"` | `1776` | Dev server port (`serve` only) |
| `-c`, `--config` | — | — | Path to a config file (otherwise auto-discovers `volcano.json` in input dir) |

### Site metadata

| CLI flag | JSON key | Default | What it does |
|----------|----------|---------|--------------|
| `--title` | `"title"` | `"My Site"` | Site title in header and `<title>` tag |
| `--author` | `"author"` | `""` | Author meta tag |
| `--og-image` | `"ogImage"` | `""` | Default Open Graph image URL |
| `--favicon` | `"favicon"` | `""` | Path to `.ico`, `.png`, or `.svg` favicon |

### Appearance

See [Appearance](/appearance/) for theme screenshots and accent-color examples.

| CLI flag | JSON key | Default | What it does |
|----------|----------|---------|--------------|
| `--theme` | `"theme"` | `"docs"` | One of `docs`, `blog`, `presentation`, `readable`, `vanilla` |
| `--css` | `"css"` | `""` | Custom CSS file (overrides `--theme`) — see [Custom CSS](/appearance/custom-css/) |
| `--accent-color` | `"accentColor"` | `"sky"` | Tailwind name, hex, or gradient (`lime-sky`, `#444-#555`) |

### Navigation features

All opt-in. See [Features](/features/) for what each one looks and feels like.

| CLI flag | JSON key | Default | Feature |
|----------|----------|---------|---------|
| `--breadcrumbs` | `"breadcrumbs"` | `false` | [Breadcrumbs](/features/#breadcrumbs) |
| `--top-nav` | `"topNav"` | `false` | [Top Navigation](/features/#top-navigation) |
| `--page-nav` | `"pageNav"` | `false` | [Previous / Next Links](/features/#previous--next-links) |
| `--instant-nav` | `"instantNav"` | `false` | [Instant Navigation](/features/#instant-navigation) |
| `--search` | `"search"` | `false` | [Search](/features/#search) |

### Advanced features

| CLI flag | JSON key | Default | Feature |
|----------|----------|---------|---------|
| `--pwa` | `"pwa"` | `false` | [PWA](/advanced/pwa/) — installable + offline |
| `--inline-assets` | `"inlineAssets"` | `false` | Embed CSS/JS in each HTML file instead of separate files |
| `--allow-broken-links` | `"allowBrokenLinks"` | `false` | Warn instead of failing the build on broken links |

### Output control

CLI-only — not in the config file.

| CLI flag | Default | What it does |
|----------|---------|--------------|
| `-q`, `--quiet` | `false` | Suppress non-error output |
| `--verbose` | `false` | Print debug info |

## Complete `volcano.json`

The file `volcano init` generates:

```json
{
  "output": "./output",
  "port": 1776,
  "title": "My Site",
  "url": "",
  "author": "",
  "theme": "docs",
  "css": "",
  "accentColor": "sky",
  "favicon": "",
  "topNav": false,
  "breadcrumbs": false,
  "pageNav": false,
  "instantNav": false,
  "inlineAssets": false,
  "pwa": false,
  "search": false,
  "ogImage": "",
  "allowBrokenLinks": false
}
```

## Worked Example

Say you want a documentation site at `https://docs.example.com`, with search and breadcrumbs, plus a custom accent color. Put the durable bits in `volcano.json`:

```json
{
  "url": "https://docs.example.com",
  "title": "My Docs",
  "theme": "docs",
  "accentColor": "emerald",
  "search": true,
  "breadcrumbs": true,
  "topNav": true,
  "pageNav": true
}
```

Then everyday commands are short:

```bash
volcano serve ./docs            # preview
volcano ./docs -o ./public      # build
```

Override per-build when needed:

```bash
volcano ./docs --title="My Docs (Beta)" -o ./public-beta
```

## Related

- **[CLI Reference](/cli/)** — every command, plus URL generation rules
- **[Features](/features/)** — what each navigation flag does
- **[Appearance](/appearance/)** — themes and accent color in context

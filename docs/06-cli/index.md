# CLI Reference

Every command and flag.

For setting things permanently via `volcano.json` — and a unified lookup table that ties flags ↔ JSON keys ↔ feature pages — see **[[01-configuration|Configuration]]**.

## Commands

| Command | Purpose |
|---------|---------|
| `volcano <folder>` or `volcano build <folder>` | Generate static site. `--url` required. |
| `volcano serve <folder>` | Dev server. Renders each request fresh. |
| `volcano` (no args) | Shortcut for `volcano serve .` |
| `volcano init [-o path]` | Create or update `volcano.json` with all options + defaults |
| `volcano css [-o file]` | Export the `vanilla` theme CSS (skeleton for custom themes) |
| `volcano --version` / `-v` | Print version |
| `volcano --help` / `-h` | Print help |

## Flags

### Required for builds

| Flag | Default | Description |
|------|---------|-------------|
| `--url` | — | Base URL for canonical / Open Graph tags. **Required** for `build`; optional for `serve`. |

### Output

| Flag | Default | Description |
|------|---------|-------------|
| `-o`, `--output` | `./output` | Output directory |
| `-p`, `--port` | `1776` | Dev server port (serve only) |

### Site metadata

| Flag | Default | Description |
|------|---------|-------------|
| `--title` | `My Site` | Site title in header and `<title>` |
| `--author` | — | Author meta tag |
| `--og-image` | — | Default Open Graph image URL |
| `--favicon` | — | Path to `.ico`, `.png`, or `.svg` favicon |

### Appearance

| Flag | Default | Description |
|------|---------|-------------|
| `--theme` | `docs` | One of `docs`, `blog`, `presentation`, `readable`, `vanilla` |
| `--css` | — | Path to custom CSS file (overrides `--theme`) |
| `--accent-color` | `sky` | Tailwind color name, hex, or two-color gradient (`lime-sky`, `#444444-#555555`) |

### Navigation (all opt-in)

| Flag | Default | Description |
|------|---------|-------------|
| `--breadcrumbs` | `false` | Show breadcrumb trail above content |
| `--top-nav` | `false` | Horizontal nav bar with root-level pages |
| `--page-nav` | `false` | Previous/next links at page bottom |
| `--instant-nav` | `false` | Hover prefetching for fast clicks |

### Advanced features

| Flag | Default | Description |
|------|---------|-------------|
| `--search` | `false` | Cmd+K command palette + search index |
| `--pwa` | `false` | Generate `manifest.json` and service worker |
| `--inline-assets` | `false` | Embed CSS/JS inline instead of separate files |
| `--allow-broken-links` | `false` | Warn instead of failing the build |

### Output control

| Flag | Default | Description |
|------|---------|-------------|
| `-c`, `--config` | — | Path to config file (otherwise looks for `volcano.json` in input directory) |
| `-q`, `--quiet` | `false` | Suppress non-error output |
| `--verbose` | `false` | Print debug info |

## Config File

Run `volcano init` to scaffold a `volcano.json` with every option at its default. Volcano then auto-discovers it in your input directory.

CLI flags always win over the config file. See **[[01-configuration|Configuration]]** for the full settings reference and worked examples.

## URL Generation

How `.md` files turn into URLs.

### Clean URLs

Each page becomes a folder with an `index.html`:

| Input | Output URL |
|-------|------------|
| `intro.md` | `/intro/` |
| `guides/setup.md` | `/guides/setup/` |
| `index.md` or `readme.md` | `/` |
| `guides/index.md` | `/guides/` |

Both `index.md` and `readme.md` are treated as folder landing pages (case-insensitive).

### Slugs

Filenames and folder names are lowercased and slugified:

| Input | Slug |
|-------|------|
| `Hello World.md` | `hello-world` |
| `API_Reference.md` | `api-reference` |
| `FAQ.md` | `faq` |

Rules: lowercase, spaces and `_` → `-`, non-alphanumeric stripped, multiple hyphens collapsed, leading/trailing hyphens trimmed.

### Prefix Stripping

Date prefixes (`YYYY-MM-DD-`) and number prefixes (`NN-`) are stripped from URLs but preserved for sort order:

| Input | URL | Sort key |
|-------|-----|----------|
| `2024-03-15-launch.md` | `/launch/` | 2024-03-15 (newest first) |
| `01-intro.md` | `/intro/` | 1 (ascending) |
| `0. Inbox/notes.md` | `/inbox/notes/` | 0 |

Separators that work: `-`, `_`, `.`, space. Both prefix styles can combine — `2024-01-15-01-featured.md` → `/featured/`.

### Display Names

Sidebar labels come from, in order:

1. The first `# H1` heading in the file
2. The filename, cleaned (prefix stripped, separators → spaces, title-cased, all-uppercase words preserved like `FAQ` / `API`)

### Hidden / Drafts

Files and folders starting with `_` or `.` are skipped:

```
_drafts/         ← ignored
.work/           ← ignored
_template.md     ← ignored
```

## Sort Order

In the sidebar, within each folder:

1. Files before folders
2. Date-prefixed files (newest first)
3. Number-prefixed files (ascending)
4. Everything else, alphabetical

## Link Validation

Builds fail if any internal link is broken — wiki links, markdown links, and navigation entries are all checked. Output looks like:

```
Found 2 broken internal links:
  Page /guides/intro/: broken link /setup/ (not found)
  Page /reference/api/: broken link /deprecated/ (not found)
```

Pass `--allow-broken-links` to turn this into a warning. In `serve` mode, broken links are shown inline on the page instead of failing.

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Anything went wrong — invalid args, missing files, broken links, generation failure |

## Generated Files

Every build writes:

- `index.html` and the page tree
- `styles.css` (combined + minified)
- `404.html` (always — most hosts serve this on missing URLs automatically)

Plus, conditionally:

- `search-index.json` + `search.js` — with `--search`
- `manifest.json` + `sw.js` — with `--pwa`

## Front Matter

YAML front matter is stripped before rendering — Volcano doesn't read any fields from it. This means existing Obsidian / Hugo / Jekyll files work without modification:

```markdown
---
date: 2024-01-15
tags: [draft]
---

# Page Title
```

Page titles come from the H1, not the front matter `title:` field.

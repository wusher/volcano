# Organizing Files

Your folder structure is your site structure. No config to manage that.

## The Tree Becomes the Sidebar

```
my-site/
├── index.md                    →  /
├── getting-started.md          →  /getting-started/
├── guides/
│   ├── index.md                →  /guides/
│   ├── installation.md         →  /guides/installation/
│   └── advanced/
│       ├── index.md            →  /guides/advanced/
│       └── tuning.md           →  /guides/advanced/tuning/
└── reference.md                →  /reference/
```

What you see in the sidebar mirrors that:

![Sidebar tree](/images/ui/sidebar-tree.png)

## Index Files

`index.md` (or `readme.md`, case-insensitive) is the landing page for a folder. Without one, Volcano auto-generates a simple listing.

## Sort Order: Prefix Tricks

### Number prefixes — force an order

```
01-introduction.md   →  /introduction/
02-installation.md   →  /installation/
10-advanced.md       →  /advanced/
```

URLs strip the `NN-` prefix. Sort is numeric (so `02-` comes before `10-`).

### Date prefixes — newest first

```
2024-03-15-launch.md     →  /launch/
2024-02-01-roadmap.md    →  /roadmap/
```

Useful for blogs. URLs strip the date. Sort is by date descending.

Separators `-`, `_`, `.`, and space all work (`01-foo`, `01_foo`, `01.foo`, `01 foo`).

## Titles

Display names come from, in order:

1. The first `# H1` heading in the file
2. The filename, cleaned up (prefix stripped, hyphens → spaces, title-cased)

```markdown
# Welcome to My Project    ← used as sidebar label + page <title>
```

## Hidden and Draft Files

Files or folders starting with `_` or `.` are skipped:

```
_drafts/               ← ignored
.work-in-progress.md   ← ignored
_template.md           ← ignored
```

Use this for in-progress content you don't want published yet.

## Linking

Prefer wiki links over hand-rolled paths — they survive renames and reorganizations:

```markdown
[[installation]]                       ← finds installation.md anywhere
[[guides/advanced/tuning|Tuning]]      ← explicit path + custom text
[[setup#requirements]]                 ← jump to a heading
```

See [[index|Writing]] for syntax.

## Next

- **[[appearance/index|Appearance]]** — themes and colors
- **[Features](/features/)** — search, breadcrumbs, page navigation

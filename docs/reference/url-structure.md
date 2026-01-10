# URL Structure

How Volcano generates clean URLs from your file structure.

## Clean URLs

Volcano generates clean URLs without file extensions:

| Input File | Output URL |
|------------|------------|
| `intro.md` | `/intro/` |
| `guides/setup.md` | `/guides/setup/` |
| `index.md` | `/` |

### Output Files

Each page becomes a directory with an `index.html`:

```
Input:                    Output:
docs/                     public/
├── intro.md              ├── intro/
├── guides/               │   └── index.html
│   └── setup.md          ├── guides/
└── index.md              │   └── setup/
                          │       └── index.html
                          └── index.html
```

This enables clean URLs like `/guides/setup/` instead of `/guides/setup.html`.

## Slugification

File and folder names are converted to URL-safe slugs.

### Rules

1. **Lowercase** — All characters converted to lowercase
2. **Spaces to hyphens** — Spaces become `-`
3. **Underscores to hyphens** — `_` becomes `-`
4. **Remove special characters** — Only letters, numbers, and hyphens kept
5. **Clean up hyphens** — Multiple hyphens collapsed to one
6. **Trim hyphens** — Leading/trailing hyphens removed

### Examples

| Input | Slug |
|-------|------|
| `Hello World.md` | `hello-world` |
| `API Reference.md` | `api-reference` |
| `Setup_Guide.md` | `setup-guide` |
| `FAQ.md` | `faq` |

## Prefix Stripping

Volcano strips date and number prefixes from URLs while preserving them for sorting.

### Date Prefixes

Files starting with `YYYY-MM-DD-` have the date stripped:

| Input | URL | Sorted By |
|-------|-----|-----------|
| `2024-01-15-hello.md` | `/hello/` | 2024-01-15 |
| `2024-02-20-world.md` | `/world/` | 2024-02-20 |

Supported separators: `-`, `_`, space

```
2024-01-15-post.md    ✓
2024_01_15_post.md    ✓
2024-01-15 post.md    ✓
```

### Number Prefixes

Files starting with numbers followed by a separator:

| Input | URL | Sorted By |
|-------|-----|-----------|
| `01-introduction.md` | `/introduction/` | 1 |
| `02-setup.md` | `/setup/` | 2 |
| `10-advanced.md` | `/advanced/` | 10 |

Supported separators: `-`, `_`, `.`, space

```
01-intro.md           ✓
01_intro.md           ✓
01 intro.md           ✓
0. Intro.md           ✓ (Obsidian style)
```

### Combined Prefixes

Date and number can be combined:

| Input | URL |
|-------|-----|
| `2024-01-15-01-featured.md` | `/featured/` |

## Folder Paths

Folder names are also slugified:

| Input Path | URL Path |
|------------|----------|
| `My Guides/setup.md` | `/my-guides/setup/` |
| `0. Inbox/notes.md` | `/inbox/notes/` |
| `API Reference/auth.md` | `/api-reference/auth/` |

Each segment is slugified independently:

```
"0. Inbox/1. Health/notes.md" → /inbox/health/notes/
```

## Index Files

Special handling for index files:

| Input | URL |
|-------|-----|
| `index.md` | `/` |
| `readme.md` | `/` |
| `guides/index.md` | `/guides/` |
| `guides/readme.md` | `/guides/` |

Both `index.md` and `readme.md` (case-insensitive) are treated as folder index files.

## Display Names

While URLs are slugified, display names are cleaned differently for the sidebar:

| Filename | Display Name | URL |
|----------|--------------|-----|
| `01-getting-started.md` | "Getting Started" | `/getting-started/` |
| `2024-01-15-hello.md` | "Hello" | `/hello/` |
| `FAQ.md` | "FAQ" | `/faq/` |
| `api_reference.md` | "Api Reference" | `/api-reference/` |

### Display Name Rules

1. Remove `.md` extension
2. Strip date/number prefixes
3. Replace `-` and `_` with spaces
4. Title case each word
5. Preserve all-uppercase words (FAQ, API)

## H1 Title Override

If a markdown file has an H1 heading, it overrides the display name:

```markdown
# Welcome to My Project

Content here...
```

Filename: `01-introduction.md`
- **Display Name:** "Welcome to My Project" (from H1)
- **URL:** `/introduction/` (from filename)

## Sorting Order

Files and folders are sorted in the sidebar:

1. **Files before folders**
2. **By date** — Files with date prefixes sorted newest first
3. **By number** — Files with number prefixes sorted ascending
4. **Alphabetically** — Remaining items sorted A-Z

### Example Sort Order

Input files:
```
2024-02-01-post.md      (has date)
2024-01-15-older.md     (has date)
01-intro.md             (has number)
02-setup.md             (has number)
faq.md                  (no prefix)
about.md                (no prefix)
```

Sorted order:
```
2024-02-01-post.md      (newest date first)
2024-01-15-older.md
01-intro.md             (lowest number)
02-setup.md
about.md                (alphabetical)
faq.md
```

## Hidden Files

Files and folders starting with `.` are ignored:

```
.git/           ← ignored
.DS_Store       ← ignored
.draft.md       ← ignored
```

## Draft Files

Files starting with `_` are treated as drafts and excluded:

```
_draft-post.md  ← excluded from build
_wip.md         ← excluded from build
```

## Related

- [[cli]] — Command line options
- [[guides/organizing-content]] — File organization best practices
- [[features/wiki-links]] — Wiki link resolution

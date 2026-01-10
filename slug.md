# Slug Rules

How Volcano converts filenames and wiki links to URL slugs.

## Filename to Slug Conversion

### Basic Rules

| Filename | Slug | URL |
|----------|------|-----|
| `Homepage.md` | `homepage` | `/homepage/` |
| `Getting Started.md` | `getting-started` | `/getting-started/` |
| `API Reference.md` | `api-reference` | `/api-reference/` |
| `FAQ.md` | `faq` | `/faq/` |

### Prefix Stripping

Date and number prefixes are stripped from slugs:

| Filename | Slug | URL |
|----------|------|-----|
| `2024-01-15-hello.md` | `hello` | `/hello/` |
| `01-introduction.md` | `introduction` | `/introduction/` |
| `001_setup.md` | `setup` | `/setup/` |

### Folder Prefixes

Folder names with number prefixes are also slugified:

| Folder | Slug | URL |
|--------|------|-----|
| `0. Inbox` | `inbox` | `/inbox/` |
| `1. Projects` | `projects` | `/projects/` |
| `01-guides` | `guides` | `/guides/` |

### Transformations Applied

1. **Remove `.md` extension**
2. **Strip date prefix** (`YYYY-MM-DD-` pattern)
3. **Strip number prefix** (`01-`, `001_`, `0. ` patterns)
4. **Convert to lowercase**
5. **Replace spaces and underscores with hyphens**
6. **Remove non-URL-safe characters** (dots, special chars)
7. **Collapse multiple hyphens**
8. **Trim leading/trailing hyphens**

## Wiki Link Resolution

### Simple Links (No Path)

`[[Page Name]]` resolves relative to the **source file's directory**:

| Current File | Wiki Link | Resolves To |
|--------------|-----------|-------------|
| `guides/intro.md` | `[[Setup]]` | `/guides/setup/` |
| `guides/intro.md` | `[[FAQ]]` | `/guides/faq/` |
| `index.md` (root) | `[[About]]` | `/about/` |

### Explicit Path Links

`[[folder/Page]]` resolves from the **site root**:

| Current File | Wiki Link | Resolves To |
|--------------|-----------|-------------|
| `guides/intro.md` | `[[reference/cli]]` | `/reference/cli/` |
| `guides/intro.md` | `[[guides/setup]]` | `/guides/setup/` |
| `index.md` | `[[guides/intro]]` | `/guides/intro/` |

### Index/Readme Resolution

Links to `index` or `readme` resolve to the parent directory:

| Wiki Link | Resolves To |
|-----------|-------------|
| `[[guides/index]]` | `/guides/` |
| `[[readme]]` | `/` (if in root) |
| `[[folder/readme]]` | `/folder/` |

### Display Text

`[[Target|Display Text]]` uses custom display text:

| Wiki Link | Display | URL |
|-----------|---------|-----|
| `[[setup|Get Started]]` | Get Started | `/setup/` |
| `[[guides/api|API Docs]]` | API Docs | `/guides/api/` |

## Examples

### File: `guides/customizing-appearance.md`

```markdown
See [[development-workflow]] for dev server info.
```

Resolves to: `/guides/development-workflow/`

### File: `guides/customizing-appearance.md`

```markdown
Check [[features/theming]] for theme details.
```

Resolves to: `/features/theming/` (explicit path = from root)

### File: `guides/index.md`

```markdown
Start with [[installation]].
```

Resolves to: `/guides/installation/`

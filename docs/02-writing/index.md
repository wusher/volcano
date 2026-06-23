# Writing

How content turns into pages.

## Markdown

Volcano renders standard CommonMark plus GitHub-flavored extensions: tables, fenced code blocks, strikethrough, task lists. The first `#` heading in each file becomes the page title.

```markdown
# Page Title

A paragraph with **bold**, *italic*, `inline code`, and a [link](https://example.com).

- Bullet
- Bullet

| Column | Column |
|--------|--------|
| Cell   | Cell   |
```

## Wiki Links

Link to other pages by name — Volcano resolves the path:

```markdown
See [[installation]] for setup.
Read about [[guides/deploying|deploying]] (custom display text).
Jump to [[cli#flags|the flags section]].
```

Resolution is case-insensitive and matches against filenames, slugs, and folder paths. Broken wiki links fail the build by default (pass `--allow-broken-links` to warn instead).

## Admonitions

Highlighted callouts using triple-colon fences. Four types ship in:

````markdown
:::note
For your information.
:::

:::tip
A helpful suggestion.
:::

:::warning
Pay attention.
:::

:::danger
Stop and read.
:::
````

![A tip admonition rendered in the docs theme](/images/ui/admonition.png)

## Code Blocks

Triple-backtick blocks with a language tag get syntax highlighting (via Chroma) and a copy button in the rendered output:

````markdown
```go
func main() {
    fmt.Println("hello")
}
```
````

## Headings → Anchors → TOC

Every `##` / `###` / `####` heading gets an auto-generated anchor (`#section-name`) and shows up in the table-of-contents sidebar on pages with 3+ headings.

```markdown
## Installation

### Requirements

### Steps
```

Link to a specific heading from anywhere:

```markdown
See [[setup#requirements]].
```

## Images

Standard markdown syntax. Paths are relative to the markdown file:

```markdown
![Alt text](images/diagram.png)
```

Images referenced from markdown are copied to the output directory.

## Front Matter

YAML front matter is stripped from rendered output (kept for Obsidian / Hugo compatibility — Volcano doesn't read any fields from it; titles come from the first H1):

```markdown
---
date: 2024-01-15
tags: [draft, notes]
---

# Page Title
```

## Next

- **[[organizing|Organizing files]]** — folders, sort order, drafts
- **[[appearance/index|Appearance]]** — make it look the way you want

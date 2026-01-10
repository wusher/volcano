# Markdown Syntax

Volcano supports GitHub Flavored Markdown (GFM) plus additional extensions.

## Basic Formatting

### Headings

```markdown
# Heading 1
## Heading 2
### Heading 3
#### Heading 4
##### Heading 5
###### Heading 6
```

The first H1 (`#`) in a file becomes the page title.

### Emphasis

```markdown
*italic* or _italic_
**bold** or __bold__
***bold italic*** or ___bold italic___
~~strikethrough~~
```

Renders as: *italic*, **bold**, ***bold italic***, ~~strikethrough~~

### Paragraphs

Separate paragraphs with a blank line:

```markdown
This is the first paragraph.

This is the second paragraph.
```

Line breaks within a paragraph are preserved (hard wraps).

## Links and Images

### Links

```markdown
[Link text](https://example.com)
[Link with title](https://example.com "Title")
<https://example.com>
```

External links automatically:
- Open in a new tab (`target="_blank"`)
- Include `rel="noopener noreferrer"` for security
- Show an external link icon

### Images

```markdown
![Alt text](image.png)
![Alt text](image.png "Title")
```

Images are displayed with `max-width: 100%` for responsiveness.

## Lists

### Unordered Lists

```markdown
- Item one
- Item two
  - Nested item
  - Another nested item
- Item three
```

### Ordered Lists

```markdown
1. First item
2. Second item
   1. Nested item
   2. Another nested item
3. Third item
```

### Task Lists

```markdown
- [x] Completed task
- [ ] Incomplete task
- [ ] Another task
```

Renders as interactive checkboxes (display only):

- [x] Completed task
- [ ] Incomplete task
- [ ] Another task

## Code

### Inline Code

```markdown
Use `code` for inline code.
```

Renders as: Use `code` for inline code.

### Code Blocks

````markdown
```javascript
function hello() {
  console.log("Hello, world!");
}
```
````

See [[code-blocks]] for syntax highlighting and advanced features.

## Tables

```markdown
| Header 1 | Header 2 | Header 3 |
|----------|----------|----------|
| Cell 1   | Cell 2   | Cell 3   |
| Cell 4   | Cell 5   | Cell 6   |
```

Alignment:

```markdown
| Left | Center | Right |
|:-----|:------:|------:|
| L    | C      | R     |
```

| Left | Center | Right |
|:-----|:------:|------:|
| L    | C      | R     |

## Blockquotes

```markdown
> This is a blockquote.
>
> It can span multiple paragraphs.

> Nested quotes:
>> Are also supported.
```

> This is a blockquote.
>
> It can span multiple paragraphs.

## Horizontal Rules

```markdown
---
```

or

```markdown
***
```

---

## Footnotes

```markdown
Here's a sentence with a footnote[^1].

[^1]: This is the footnote content.
```

Here's a sentence with a footnote[^1].

[^1]: This is the footnote content.

Footnotes are collected at the bottom of the page.

## Definition Lists

```markdown
Term 1
: Definition for term 1

Term 2
: Definition for term 2
: Another definition for term 2
```

Term 1
: Definition for term 1

Term 2
: Definition for term 2
: Another definition for term 2

## Smart Typography

Volcano automatically converts:

| Input | Output |
|-------|--------|
| `"quotes"` | "curly quotes" |
| `'quotes'` | 'curly quotes' |
| `--` | en-dash – |
| `---` | em-dash — |
| `...` | ellipsis … |

## Raw HTML

HTML is allowed in markdown:

```markdown
<div class="custom">
  Custom HTML content
</div>
```

Use sparingly — markdown is preferred for maintainability.

## Escaping Special Characters

Use backslash to escape markdown characters:

```markdown
\*not italic\*
\# not a heading
\[not a link\]
```

Renders as: \*not italic\*, \# not a heading, \[not a link\]

## Volcano Extensions

Beyond standard GFM, Volcano adds:

- **[[wiki-links]]** — `[[Page Name]]` syntax
- **[[admonitions]]** — `:::note` callout blocks
- **[[code-blocks|Line highlighting]]** — `` ```js {1,3-5} ``

## Related

- [[wiki-links]] — Link between pages easily
- [[admonitions]] — Create callout boxes
- [[code-blocks]] — Syntax highlighting details

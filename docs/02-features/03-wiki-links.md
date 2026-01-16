# Wiki Links

Link to other pages using double brackets:

```markdown
See the [[Installation]] guide.
```

## Custom Display Text

```markdown
Read the [[Installation|install guide]] first.
```

## Paths

```markdown
See [[guides/installation]].
Check [[reference/cli|the CLI docs]].
```

Paths are relative to the site root.

## Heading Anchors

```markdown
See [[Installation#requirements]].
Jump to [[#section]] on the same page.
```

## How Links Resolve

| Current Page | Link | Resolves To |
|--------------|------|-------------|
| `/guides/` | `[[Installation]]` | `/guides/installation/` |
| `/` | `[[guides/installation]]` | `/guides/installation/` |
| `/guides/` | `[[Installation#setup]]` | `/guides/installation/#setup` |

## Obsidian Compatibility

**Supported:**
- `[[Page Name]]`
- `[[Page Name|Display Text]]`
- `[[folder/Page Name]]`
- `[[Page Name#Heading]]`
- `[[#Heading]]` (same page)
- `![[Page Name]]` (converted to link)

**Not Supported:**
- `[[Page Name#^block]]` (block references)
- Transclusion (embedding content)

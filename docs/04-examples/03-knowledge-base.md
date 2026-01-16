# Knowledge Base

A wiki-style knowledge base with interconnected notes:

```
wiki/
├── index.md
├── projects/
│   ├── index.md
│   └── project-alpha.md
├── processes/
│   ├── index.md
│   └── code-review.md
└── references/
    └── tools.md
```

## Wiki Links

```markdown
See [[coding-standards]] for details.
Check the [[code-review|code review process]].
Follow [[processes/deployment|deployment steps]].
Jump to [[page#section]].
```

## Build Command

```bash
volcano ./wiki \
  -o ./public \
  --title="Team Wiki" \
  --url="https://wiki.example.com" \
  --search
```

## Obsidian Compatibility

**Supported:**
- `[[Page Name]]` — Wiki links
- `[[Page|Text]]` — Custom link text
- `[[folder/Page]]` — Path links
- `[[Page#Heading]]` — Heading anchors
- `![[Page]]` — Converted to regular link
- Numbered folders — `01-Inbox` → `/inbox/`
- Front matter — Stripped from output

**Not supported:**
- Block references (`[[Page#^block]]`)
- Transclusion (embedding)
- Tags (`#tag`)
- Auto-generated backlinks

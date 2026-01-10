# Knowledge Base Example

How to create a wiki-style knowledge base with Volcano.

## Overview

A knowledge base is a collection of interconnected notes, perfect for:

- Personal wikis
- Team documentation
- Published Obsidian vaults
- Note collections

## Structure

### Flat Organization

```
wiki/
├── index.md
├── project-alpha.md
├── meeting-notes.md
├── coding-standards.md
├── team-contacts.md
└── resources.md
```

Best for smaller knowledge bases with free-form linking.

### Topic-Based

```
wiki/
├── index.md
├── projects/
│   ├── index.md
│   ├── project-alpha.md
│   └── project-beta.md
├── processes/
│   ├── index.md
│   ├── code-review.md
│   └── deployment.md
├── references/
│   ├── index.md
│   ├── tools.md
│   └── resources.md
└── people/
    └── team.md
```

Best for larger knowledge bases with clear categories.

### Obsidian Vault

```
vault/
├── 0. Inbox/
│   └── quick-notes.md
├── 1. Projects/
│   ├── project-alpha.md
│   └── project-beta.md
├── 2. Areas/
│   ├── health.md
│   └── finance.md
├── 3. Resources/
│   └── books.md
└── 4. Archive/
    └── old-notes.md
```

Volcano handles Obsidian's numbered folders — prefixes are stripped from URLs.

## Wiki Links

The power of a knowledge base is interconnection:

### Basic Links

```markdown
# Project Alpha

This project uses our [[coding-standards]].

See also:
- [[project-beta]] — Related project
- [[team-contacts]] — Team members
```

### Links with Labels

```markdown
Check the [[code-review|code review process]] before submitting.
Contact [[team-contacts|the team]] for questions.
```

### Cross-Folder Links

```markdown
# Project Alpha

Follow our [[processes/code-review|code review process]].
See [[references/tools|recommended tools]].
```

## Example Pages

### index.md (Homepage)

```markdown
# Knowledge Base

Welcome to our team knowledge base.

## Quick Links

- [[projects/index|Projects]] — Active and past projects
- [[processes/index|Processes]] — How we work
- [[references/index|References]] — Tools and resources

## Recent Updates

- [[projects/project-alpha]] — Updated deployment guide
- [[processes/code-review]] — New review checklist

## Need Help?

Contact [[people/team|the team]] or check [[references/resources]].
```

### projects/project-alpha.md

```markdown
# Project Alpha

Our main product project.

## Overview

Project Alpha is our flagship product. Development started in Q1 2024.

## Team

- Lead: [[people/team#alice|Alice]]
- Backend: [[people/team#bob|Bob]]
- Frontend: [[people/team#carol|Carol]]

## Documentation

- [[processes/deployment|Deployment Process]]
- [[references/tools#monitoring|Monitoring Setup]]

## Links

- Repository: [GitHub](https://github.com/example/alpha)
- Staging: [staging.example.com](https://staging.example.com)

## Notes

:::note
This project follows our [[coding-standards]].
:::
```

### processes/code-review.md

```markdown
# Code Review Process

How we review code before merging.

## Checklist

- [ ] Code follows [[coding-standards]]
- [ ] Tests pass locally
- [ ] Documentation updated
- [ ] No security issues

## Process

1. Create pull request
2. Request review from [[people/team|team member]]
3. Address feedback
4. Get approval
5. Merge and deploy per [[deployment|deployment process]]

## Tips

:::tip
Use the PR template for consistent descriptions.
:::

## Related

- [[coding-standards]] — Code style guide
- [[deployment]] — How to deploy
```

## Build Command

```bash
volcano ./wiki \
  -o ./public \
  --title="Team Wiki" \
  --url="https://wiki.example.com"
```

### For Obsidian Vaults

```bash
volcano ./vault \
  -o ./public \
  --title="My Notes" \
  --url="https://notes.example.com"
```

## Wiki-Style Features

### Backlinks (Manual)

Volcano doesn't auto-generate backlinks. Create them manually:

```markdown
# Coding Standards

...

## Pages That Reference This

- [[project-alpha]] — Uses these standards
- [[code-review]] — Checks against these standards
```

### Table of Contents

Each page automatically gets a table of contents in the sidebar.

### Search

The sidebar includes search to find pages quickly.

## Obsidian Compatibility

### Supported

- `[[Page Name]]` — Wiki links
- `[[Page Name|Display Text]]` — Aliased links
- `[[folder/Page Name]]` — Path links
- `![[Page Name]]` — Converted to regular link
- Front matter — Stripped cleanly
- Numbered folders — `0. Inbox` → `/inbox/`

### Not Supported

- `[[Page#Heading]]` — Heading anchors
- `[[Page#^block]]` — Block references
- Transclusion — Embedding content
- Tags (`#tag`) — Not processed
- Backlinks — Not auto-generated

## Custom Styling

### Wiki-Friendly CSS

```css
/* Denser content for reference */
.prose {
  max-width: 800px;
}

/* Clear wiki link styling */
.prose a {
  color: var(--color-link);
  text-decoration: none;
  border-bottom: 1px dotted var(--color-link);
}

.prose a:hover {
  border-bottom-style: solid;
}

/* Prominent internal links */
.prose a[href^="/"]:not([href^="//"]) {
  font-weight: 500;
}

/* Callout boxes for notes */
.admonition {
  margin: 1.5rem 0;
  padding: 1rem;
  border-radius: 4px;
}
```

## Tips for Knowledge Bases

### Linking Strategy

1. **Link generously** — When mentioning a concept, link it
2. **Use aliases** — `[[coding-standards|our standards]]` reads better
3. **Bidirectional** — Link both directions between related pages

### Organization

1. **Start flat** — Add folders only when needed
2. **Index pages** — Each folder should have an overview
3. **Consistent naming** — Use same terms across pages

### Maintenance

1. **Review regularly** — Update outdated information
2. **Fix broken links** — Check for pages that moved
3. **Archive old content** — Move to archive folder

## Related

- [[features/wiki-links]] — Wiki link syntax
- [[reference/url-structure]] — How URLs are generated
- [[guides/organizing-content]] — File organization

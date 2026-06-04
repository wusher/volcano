# example-nested

A regression fixture for the deep-nested-directory bug.

## Why this exists

When the input directory has **no `.md` files at the top level** — only subdirectories — Volcano used to fall back to a static file server, and every URL returned 404 even though the markdown files existed deeper in the tree.

This fixture reproduces that layout:

```
example-nested/
├── 1/
│   ├── 1/   research-design-doc.md, interview-summary.md
│   ├── 2/   api-surface-proposal.md, migration-checklist.md
│   └── 3/   timing-experiments.md, observations.md
├── 2/
│   ├── 1/   launch-readiness-review.md, rollback-runbook.md
│   ├── 2/   incident-postmortem.md, lessons-learned.md
│   └── 3/   security-considerations.md, threat-model-draft.md
└── 3/
    ├── 1/   onboarding-guide.md, architecture-overview.md
    ├── 2/   api-reference.md, cli-cheatsheet.md
    └── 3/   glossary.md, changelog.md
```

Key properties matched to the real-world report:

- **Numeric-only folder names** at every directory level.
- **Kebab-case filenames** at the leaves.
- **No `index.md` anywhere** — not at root, not in any folder.
- **YAML frontmatter** on every leaf file with `status`, `priority`, `issue_id`, `tags`, `dependencies: []`.

## How to test

```bash
go build -o volcano .
./volcano serve example-nested
```

Then visit:

- <http://localhost:1776/> — should show an auto-index of `1/`, `2/`, `3/`
- <http://localhost:1776/1/> — auto-index of `1/1`, `1/2`, `1/3`
- <http://localhost:1776/1/1/research-design-doc/> — leaf page renders with frontmatter stripped

All three should return HTTP 200. If any return 404, the fix in `cmd/serve.go` (`hasMarkdownDescendant`) has regressed.

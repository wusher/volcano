# AI Prompts

Tested prompts for generating Volcano content with Claude, ChatGPT, or other LLMs.

## Custom Theme

For an LLM that can run shell commands (Claude Code, Cursor, etc.):

````
I'm building a static site with Volcano and need help creating a custom CSS theme.

## Step 1: Ask me a few questions

Ask one at a time:
- What type of site? (blog, docs, portfolio, etc.)
- Overall feel? (professional, playful, minimal, bold)
- Brand colors, or should you suggest some?
- Light, dark, or both?
- Typography style? (modern, classic, technical)
- Any websites whose design I like?

## Step 2: Export Volcano's CSS skeleton

```
go install github.com/wusher/volcano@latest
volcano css -o skeleton.css
```

Read skeleton.css to understand the CSS variables, components, and selectors available.

## Step 3: Generate the theme

Based on my answers AND the skeleton structure, produce `custom.css` that:
1. Preserves the skeleton structure (same selectors)
2. Applies my preferences (colors, typography, spacing)
3. Includes both light (`:root`) and dark (`[data-theme="dark"]`) modes
4. Maintains WCAG AA contrast
5. Keeps responsive breakpoints

## Step 4: Show me how to use it

```
volcano ./docs --css ./custom.css -o ./public --url="https://example.com"
volcano serve ./public
```

Ready when you are — ask your first question.
````

## Generate a Docs Site From Source Code

Hand an LLM a repo and ask for a `docs/` folder:

````
Generate documentation for my project that will be built with Volcano.

About my project:
- Name: [PROJECT NAME]
- Description: [BRIEF DESCRIPTION]
- Language: [Go / Python / TypeScript / etc.]
- Source: [github URL or path]
- Audience: [developers / end users / both]

Instructions:

1. Analyze the source to identify the public API and main concepts.

2. Create a `docs/` folder with:
   - `index.md` — overview, value prop, quick links
   - `quickstart.md` — install + minimal example
   - `guides/` — practical tutorials
   - `reference/` — one file per major module / class
     - Description, constructor, methods (signature + params + return + example),
       properties, types/enums

3. Each markdown file uses:
   - Clear H1 title
   - H2/H3 section structure
   - Tables for parameters and return values
   - Fenced code blocks with language tags
   - Admonitions (`:::note`, `:::tip`, `:::warning`) sparingly
   - Wiki links (`[[page-name]]`) for cross-references

4. To preview:
   ```
   go install github.com/wusher/volcano@latest
   volcano serve ./docs
   ```

5. To deploy:
   ```
   volcano ./docs -o ./public --url="https://[PROJECT].example.com" --title="[PROJECT]"
   ```

Generate the docs now.
````

## Wire Up GitHub Pages

````
Set up GitHub Pages for my Volcano docs.

Fill in:
- REPO:        [owner/name]
- BRANCH:      [main]
- DOCS_FOLDER: [docs]
- ACCENT:      [#0284c7 or a Tailwind name like "sky"]

The Pages URL for a repo `owner/name` is `https://owner.github.io/name/`.

Create `.github/workflows/deploy-docs.yml` with these top-level permissions:

```yaml
permissions:
  contents: read
  pages: write
  id-token: write
```

Use this workflow as the template, substituting my values:

```yaml
name: Deploy Documentation
on:
  push:
    branches: [BRANCH]

jobs:
  docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
          cache: true
      - run: |
          go install github.com/wusher/volcano@latest
          echo "$HOME/go/bin" >> $GITHUB_PATH
      - run: |
          cd DOCS_FOLDER
          volcano build . --output ../public --url "https://owner.github.io/name/" --accent-color "ACCENT"
      - uses: actions/upload-pages-artifact@v3
        with:
          path: ./public
      - uses: actions/deploy-pages@v4
```

After the file is created, tell me to enable Pages with Source = "GitHub Actions" in repo settings.
````

## Tips

- **Be specific about your stack.** "Volcano docs theme" → the LLM may invent flags. Paste the [CLI reference](/cli/) in or link to it.
- **Iterate.** Ask for one section at a time. Easier to spot mistakes.
- **Verify links.** LLMs hallucinate file paths. Run a build (`volcano ./docs --url=...`) — Volcano's link validator catches the obvious mistakes.

# Add GitHub Pages to Your Repo

Set up GitHub Pages deployment for your project documentation using Volcano.

## About

This prompt helps you configure GitHub Actions to automatically build and deploy your Volcano documentation site to GitHub Pages whenever you push changes.

## Prompt

```
I need help setting up GitHub Pages for my project documentation using Volcano static site generator.

Fill in these values:
- REPO: [e.g., wusher/volcano]
- BRANCH: [main or master]
- DOCS_FOLDER: [e.g., docs]
- ACCENT_COLOR: [e.g., #0284c7]

Using the example workflow below, create a GitHub Actions workflow at .github/workflows/deploy-docs.yml with my values substituted.

The workflow needs these permissions at the top level:
  permissions:
    contents: read
    pages: write
    id-token: write

After creating the workflow, tell me to:
1. Go to repository Settings > Pages
2. Set Source to "GitHub Actions"

## Example Workflow

Note: GITHUB_PAGES_URL for REPO "owner/name" is https://owner.github.io/name/

  docs:
    name: Deploy Documentation
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/BRANCH'
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true

      - name: Install Volcano
        run: |
          go install github.com/wusher/volcano@latest
          echo "$HOME/go/bin" >> $GITHUB_PATH

      - name: Build documentation
        run: |
          cd DOCS_FOLDER
          volcano build . --output ../public --url GITHUB_PAGES_URL --view-transitions --accent-color="ACCENT_COLOR"

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./public

      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4

Please create the workflow file now.
```

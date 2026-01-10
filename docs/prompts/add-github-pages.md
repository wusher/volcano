# Add GitHub Pages to Your Repo Using Volcano

Use this prompt with an LLM to set up GitHub Pages deployment for your Volcano documentation.

## Prompt

```
I need help setting up GitHub Pages for my project documentation using Volcano static site generator.

About my repository:
- Repository URL: [YOUR GITHUB REPO URL]
- Repository structure: [e.g., docs/ folder at root, separate docs branch, etc.]
- Documentation location: [e.g., ./docs, ./documentation, etc.]
- Preferred deployment method: [GitHub Actions, manual, branch-based]

Please help me set up GitHub Pages with Volcano by providing:

1. **GitHub Actions Workflow** (.github/workflows/deploy-docs.yml):
   - Trigger on push to main branch (and/or docs changes)
   - Install Volcano (via go install or binary download)
   - Build the documentation site
   - Deploy to GitHub Pages
   - Include caching for faster builds

2. **Configuration Steps**:
   - Repository settings needed for GitHub Pages
   - Branch and folder configuration
   - Custom domain setup (if applicable)
   - HTTPS configuration

3. **Build Script** (optional):
   - Shell script or Makefile for local testing
   - Command to build docs locally
   - Command to serve docs locally for preview

4. **Documentation Updates**:
   - README section explaining how the docs are built and deployed
   - Instructions for contributors on how to preview docs locally
   - Badge/link to the live documentation site

5. **Volcano Configuration**:
   - Recommended build command
   - Output directory setup
   - Theme selection (docs, blog, or custom)
   - Title and other options

Additional requirements:
- [ ] Custom domain: [DOMAIN NAME if applicable]
- [ ] Deploy only on docs changes (not all commits)
- [ ] Include build status badge
- [ ] Support for multiple doc versions
- [ ] Deploy to a subdirectory [e.g., /docs]
- [ ] Other: [SPECIFY ANY OTHER REQUIREMENTS]
```

## Example GitHub Actions Workflow

Here's a basic template you can customize:

```yaml
name: Deploy Documentation

on:
  push:
    branches: [main]
    paths:
      - 'docs/**'
      - '.github/workflows/deploy-docs.yml'

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install Volcano
        run: go install github.com/wusher/volcano@latest

      - name: Build Documentation
        run: volcano ./docs -o ./public --title="My Project Docs"

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./public

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## Tips

- Test the build locally first: `volcano ./docs -o ./public`
- Preview locally: `volcano -s -p 8080 ./public`
- Enable GitHub Pages in repository Settings > Pages
- Choose "GitHub Actions" as the source (not branch-based if using Actions)
- Add a CNAME file in your docs folder if using a custom domain
- Consider adding a `.nojekyll` file if you have underscore-prefixed files
- Use branch protection rules to prevent accidental deployments
- Check the Actions tab for deployment logs and troubleshooting

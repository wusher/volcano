# Volcano

**Turn any markdown folder into a website. Zero config.**

<style>
.logo-float { max-width: 220px; }
@media (min-width: 768px) {
  .logo-float { float: right; margin-left: 2rem; }
}
</style>

<img src="logo.png" alt="" class="logo-float">

Point Volcano at a folder of `.md` files. Get a styled, navigable site with sidebar tree, search, dark mode, wiki links, and clean URLs — no `_config.yml`, no plugins, no build pipeline. Single Go binary.

It works on whatever you already have: an Obsidian vault, a `docs/` folder, a pile of notes.

## 30-Second Start

```bash
# Install
go install github.com/wusher/volcano@latest

# Preview a folder in the browser
volcano serve ./my-notes
```

That's it. Open [http://localhost:1776](http://localhost:1776). Edit files, refresh the page, see changes.

![Volcano running on its own docs](/images/ui/homepage.png)

## What You Get By Default

- **Sidebar tree** — your folder structure is the navigation
- **Wiki links** — `[[Page Name]]` resolves automatically
- **Admonitions** — `:::tip`, `:::note`, `:::warning`, `:::danger` callout boxes
- **Code highlighting** with copy buttons
- **Image lightbox** — click any image in the content area to view it full-size
- **Dark mode** toggle (press `t`)
- **Clean URLs** — `setup.md` becomes `/setup/`
- **SEO meta tags** — Open Graph, canonical, schema.org
- **Mobile responsive**

Optional with one flag each: `--search` (Cmd+K palette), `--breadcrumbs`, `--top-nav`, `--page-nav`, `--instant-nav`, `--pwa`.

## Where Next

- **[Quickstart](/quickstart/)** — install → write → preview → deploy in 15 minutes
- **[Writing](/writing/)** — markdown, wiki links, admonitions, file organization
- **[Appearance](/appearance/)** — themes, accent colors, custom CSS
- **[Features](/features/)** — search, navigation, keyboard shortcuts
- **[Deploying](/deploying/)** — GitHub Pages, Netlify, Vercel, self-hosting
- **[CLI Reference](/cli/)** — every flag and config option

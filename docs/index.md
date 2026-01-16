# Volcano

**Turn markdown folders into websites. Zero config.**

<style>
.logo-float {
  max-width: 260px;
}

@media (min-width: 768px) {
  .logo-float {
    float: right;
  }
}
</style>

<img src="logo.png" alt="Volcano Logo" class="logo-float">


---

Point Volcano at any markdown folder and get a beautiful, navigable website. Works with Obsidian vaults, documentation folders, or any collection of markdown files.

**Single binary. No dependencies. No config files.**

## Quick Start

```bash
# Install
go install github.com/wusher/volcano@latest

# Generate a site
volcano ./my-notes --title="My Site"

# Preview it
volcano serve ./output
```

Your folder structure becomes the site navigation. That's it.

## Features

- Tree navigation and search
- Clean URLs and SEO tags
- Dark mode
- Wiki links and admonitions
- Code highlighting with copy button
- Mobile responsive

## Documentation

- [[Getting Started]] — Build your first site in 5 minutes
- [[guides/index|Guides]] — Learn Volcano features
- [[reference/index|Reference]] — CLI flags and options
- [[examples/index|Examples]] — Real-world setups

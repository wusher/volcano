# Getting Started

Build your first Volcano site in 5 minutes.

## Install

Requires [Go](https://go.dev/dl/) 1.24 or later.

```bash
go install github.com/wusher/volcano@latest
volcano --version
```

## Create Content

```bash
mkdir my-docs && cd my-docs

cat > index.md << 'EOF'
# Welcome

This is my documentation site.

Check out the [[Installation]] guide.
EOF

cat > installation.md << 'EOF'
# Installation

1. Download the installer
2. Run the setup wizard
3. Follow the prompts

:::tip
Use the default installation path.
:::
EOF
```

## Generate Site

```bash
volcano . -o ./public --title="My Docs"
```

## Preview

```bash
volcano serve .
```

Open [http://localhost:1776](http://localhost:1776)

You'll see sidebar navigation, breadcrumbs, wiki links, and styled admonitions.

## Next Steps

- [[guides/organizing-content]] — Structure your content
- [[guides/theming]] — Themes and styles
- [[features/navigation]] — Navigation and search features
- [[reference/cli]] — All CLI options

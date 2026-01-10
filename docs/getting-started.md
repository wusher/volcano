# Getting Started

Get your first Volcano site running in under 5 minutes.

## Prerequisites

You need [Go](https://go.dev/dl/) 1.21 or later installed on your system.

Verify your Go installation:

```bash
go version
```

## Install Volcano

Install Volcano using `go install`:

```bash
go install github.com/wusher/volcano@latest
```

Verify the installation:

```bash
volcano --version
```

## Create Your Content

Create a folder for your documentation:

```bash
mkdir my-docs
cd my-docs
```

Create your first page. This will be your site's homepage:

```bash
cat > index.md << 'EOF'
# Welcome to My Site

This is my documentation site built with Volcano.

## Getting Started

Check out the [[Installation]] guide to get started.
EOF
```

Create another page:

```bash
cat > installation.md << 'EOF'
# Installation

Follow these steps to install the software.

## Requirements

- Operating system: Windows, macOS, or Linux
- Disk space: 100MB

## Steps

1. Download the installer
2. Run the setup wizard
3. Follow the prompts

:::tip
Use the default installation path for best compatibility.
:::
EOF
```

## Generate Your Site

Run Volcano to generate your site:

```bash
volcano . -o ./public --title="My Documentation"
```

You'll see output like:

```
Generating site...
  Input:  .
  Output: ./public
  Title:  My Documentation

Scanning input directory...
Found 2 markdown files in 1 folder

Generating pages...
  ✓ index.md
  ✓ installation.md

Generated 2 pages in ./public
```

## Preview Your Site

Start the built-in server to preview your site:

```bash
volcano -s ./public
```

Open your browser to [http://localhost:1776](http://localhost:1776).

You'll see your site with:
- A sidebar navigation showing your pages
- Breadcrumb navigation at the top
- The wiki link `[[Installation]]` converted to a working link
- The tip admonition styled as a callout box

## Next Steps

You now have a working Volcano site. Here's where to go next:

- **[[guides/organizing-content]]** — Learn how to structure your content
- **[[guides/customizing-appearance]]** — Change themes and add custom styles
- **[[features/markdown-syntax]]** — Explore all supported markdown features
- **[[reference/cli]]** — See all available command-line options

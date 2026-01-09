# Getting Started

This guide will help you get up and running with Volcano in just a few minutes.

## Installation

Download the latest release from GitHub:

```bash
go install github.com/example/volcano@latest
```

Or build from source:

```bash
git clone https://github.com/example/volcano.git
cd volcano
go build -o volcano .
```

## Your First Site

1. Create a folder with your markdown files:

```
docs/
├── index.md
├── getting-started.md
└── guides/
    ├── index.md
    └── configuration.md
```

2. Generate the site:

```bash
volcano ./docs -o ./public --title="My Documentation"
```

3. Serve it locally:

```bash
volcano -s ./public
```

4. Open `http://localhost:1776` in your browser!

## Folder Structure

Volcano automatically creates navigation from your folder structure:

| Source File | Output URL |
|------------|------------|
| `index.md` | `/` |
| `about.md` | `/about/` |
| `guides/index.md` | `/guides/` |
| `guides/setup.md` | `/guides/setup/` |

## Next Steps

- Learn about [Configuration](guides/configuration/) options
- Explore the [API Reference](api/)
- Check out [Installation](guides/installation/) details

# Quickstart

From zero to a deployed site in 15 minutes.

## 1. Install

Requires [Go](https://go.dev/dl/) 1.24 or later.

```bash
go install github.com/wusher/volcano@latest
volcano --version
```

Make sure `$(go env GOPATH)/bin` is on your `PATH`. To build from source:

```bash
git clone https://github.com/wusher/volcano.git
cd volcano
go build -o volcano .
```

## 2. Write a Page

```bash
mkdir my-site && cd my-site

cat > index.md << 'EOF'
# My Site

Welcome. Read about [[installation]].

:::tip
You don't need to declare links — `[[installation]]` finds `installation.md`
automatically.
:::
EOF

cat > installation.md << 'EOF'
# Installation

1. Download the installer
2. Run it
3. Done
EOF
```

## 3. Preview It

```bash
volcano serve .
```

Open [http://localhost:1776](http://localhost:1776). Edit either file, hit refresh, see changes — no build step.

![Volcano running locally](/images/ui/homepage.png)

## 4. Build for Deployment

A static build needs a base URL (Volcano uses it for canonical and Open Graph tags):

```bash
volcano . -o ./public --url="https://example.com" --title="My Site"
```

You now have a `./public` folder containing static HTML, CSS, and JS. Drop it on any web host.

## 5. Deploy

Pick a platform — most need zero configuration beyond pointing at `./public`:

- **GitHub Pages** — commit `./public` to a `gh-pages` branch, or use the Actions snippet in [Deploying](/deploying/)
- **Netlify / Vercel / Cloudflare Pages** — set build command to `volcano . -o ./public --url=...` and publish directory to `public`
- **Self-host** — `nginx`, `Caddy`, or any static file server

See [Deploying](/deploying/) for copy-paste configs.

## What's Next

- **Make it look how you want** — [Themes & accent colors](/appearance/)
- **Organize a bigger site** — [Folders, prefixes, drafts](/writing/organizing/)
- **Turn on search / breadcrumbs / etc.** — [Features](/features/)
- **Every flag** — [CLI reference](/cli/)

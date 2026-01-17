# Blog

A blog with date-ordered posts:

```
blog/
├── index.md
├── 2024-03-15-latest-post.md
├── 2024-02-20-another-post.md
└── 2024-01-10-first-post.md
```

Date prefixes (`YYYY-MM-DD-`) control sort order (newest first) and are stripped from URLs.

## Build Command

```bash
volcano ./blog \
  -o ./public \
  --theme blog \
  --title="My Blog" \
  --url="https://blog.example.com" \
  --page-nav
```

## Organization Options

**By year:**
```
blog/2024/2024-03-15-post.md  →  /2024/post/
```

**By category:**
```
blog/tech/2024-03-15-post.md  →  /tech/post/
```

**Flat (simplest):**
```
blog/2024-03-15-post.md  →  /post/
```

## Deployment

**Netlify:**
```toml
[build]
  command = "volcano ./blog -o ./public --theme blog --title='My Blog'"
  publish = "public"
```

**Vercel:**
```json
{
  "buildCommand": "volcano ./blog -o ./public --theme blog",
  "outputDirectory": "public"
}
```

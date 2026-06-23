# blog Theme

A reading-focused theme. Content-centered, narrower column, optimized for long-form prose.

```bash
volcano ./posts --theme blog --url="https://example.com"
```

## Light

![blog theme, light mode](/images/themes/blog-light.png)

## Dark

![blog theme, dark mode](/images/themes/blog-dark.png)

## Features

- Narrower content column tuned for comfortable line length (~70 characters)
- Reading-optimized typography — generous line height, refined hierarchy
- Dark mode support
- Minimal chrome — the words are the point
- Works well with date-prefixed filenames (`2024-03-15-post.md` → newest first)

## Best For

Blogs, newsletters, articles, release announcements, changelogs, essays — anywhere a reader scrolls top-to-bottom on a single page.

## Try It With Page Nav

```bash
volcano ./posts \
  --theme blog \
  --page-nav \
  --url="https://blog.example.com"
```

Page nav adds previous/next links at the bottom of each post — useful for sequential reading. See [Features](/features/).

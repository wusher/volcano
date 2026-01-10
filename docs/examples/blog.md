# Blog Example

How to create a blog with Volcano.

## Structure

```
blog/
├── index.md                        # Blog homepage
├── 2024-03-15-latest-post.md       # Newest first
├── 2024-02-20-another-post.md
├── 2024-01-10-first-post.md
└── about.md                        # Static page
```

### Date Prefix Format

Posts use `YYYY-MM-DD-` prefix:

```
2024-03-15-my-post-title.md
```

- **Date extracted:** March 15, 2024
- **URL generated:** `/my-post-title/`
- **Sorting:** Newest posts appear first

## Blog Homepage

### index.md

```markdown
# My Blog

Welcome to my blog where I write about technology and life.

## Recent Posts

- [[2024-03-15-latest-post|Latest Post]] — March 15, 2024
- [[2024-02-20-another-post|Another Post]] — February 20, 2024
- [[2024-01-10-first-post|First Post]] — January 10, 2024

## About

I'm a software developer writing about my experiences.
[More about me →](about/)
```

### Alternative: Auto-generated Index

Don't create an `index.md` — Volcano will auto-generate a list of posts.

## Blog Post Template

### 2024-03-15-latest-post.md

```markdown
# My Latest Post Title

A brief introduction that explains what this post is about.
This becomes the meta description for SEO.

## Main Section

Your content here. Use headings to structure the post.

### Subsection

More detailed content.

## Code Examples

Include code when relevant:

```javascript
function example() {
  return "Hello, World!";
}
```

## Conclusion

Wrap up your thoughts and invite discussion.

---

*Thanks for reading! Follow me on [Twitter](https://twitter.com/example).*
```

## Build Command

```bash
volcano ./blog \
  -o ./public \
  --theme blog \
  --title="My Blog" \
  --url="https://blog.example.com" \
  --author="Your Name" \
  --og-image="https://blog.example.com/og.png" \
  --last-modified
```

### Recommended Flags

| Flag | Purpose |
|------|---------|
| `--theme blog` | Reading-focused layout |
| `--last-modified` | Shows post dates |
| `--author` | Adds author meta tag |
| `--og-image` | Social media preview image |

## Blog Theme Features

The `blog` theme is optimized for reading:

- **Centered content** — Focused reading experience
- **Optimized typography** — Comfortable line length and spacing
- **Minimal distractions** — Clean, uncluttered design
- **Dark mode** — Easy on the eyes

## Organizing Posts

### By Year

```
blog/
├── index.md
├── 2024/
│   ├── 2024-03-15-post.md
│   └── 2024-02-20-post.md
└── 2023/
    └── 2023-12-01-post.md
```

URLs: `/2024/post/`, `/2023/post/`

### By Category

```
blog/
├── index.md
├── tech/
│   ├── 2024-03-15-coding-tips.md
│   └── 2024-02-20-new-framework.md
└── life/
    └── 2024-01-10-travel-story.md
```

URLs: `/tech/coding-tips/`, `/life/travel-story/`

### Flat (Recommended for Small Blogs)

```
blog/
├── index.md
├── 2024-03-15-post-a.md
├── 2024-02-20-post-b.md
└── 2024-01-10-post-c.md
```

URLs: `/post-a/`, `/post-b/`, `/post-c/`

## Custom Blog Styling

### Export and Customize

```bash
# Start with blog theme CSS
volcano css -o blog-custom.css

# Use your customized version
volcano ./blog --css ./blog-custom.css
```

### Example Customizations

```css
/* Wider reading width */
:root {
  --content-max-width: 720px;
}

/* Serif fonts for posts */
.prose {
  font-family: Georgia, "Times New Roman", serif;
}

/* Larger body text */
.prose p {
  font-size: 1.125rem;
  line-height: 1.8;
}

/* Stylish blockquotes */
.prose blockquote {
  border-left: 4px solid var(--color-accent);
  padding-left: 1.5rem;
  font-style: italic;
}

/* Featured post styling */
.prose h1 {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}
```

## RSS Feed

Volcano doesn't generate RSS feeds automatically. Options:

1. **Manual RSS** — Create an `rss.xml` file
2. **Build script** — Generate RSS during CI/CD
3. **Third-party** — Use a service like Feedburner

## Deployment

### Netlify

```bash
# netlify.toml
[build]
  command = "volcano ./blog -o ./public --theme blog --title='My Blog'"
  publish = "public"
```

### Vercel

```json
{
  "buildCommand": "volcano ./blog -o ./public --theme blog",
  "outputDirectory": "public"
}
```

### GitHub Pages

```yaml
# .github/workflows/deploy.yml
- name: Build Blog
  run: |
    volcano ./blog \
      -o ./public \
      --theme blog \
      --title="My Blog" \
      --url="https://username.github.io/blog"
```

## Writing Tips

### Good Post Structure

1. **Hook** — Compelling opening paragraph
2. **Context** — Background information
3. **Main content** — Organized with headings
4. **Conclusion** — Summary and next steps

### SEO Tips

- Write descriptive titles (H1)
- Strong opening paragraph (becomes meta description)
- Use headings for structure
- Include relevant keywords naturally

### Readability

- Keep paragraphs short (3-4 sentences)
- Use lists for multiple items
- Include images where helpful
- Break up long posts with headings

## Related

- [[features/theming]] — Theme customization
- [[features/reading-time]] — Reading time calculation
- [[guides/deploying-your-site]] — Deployment options

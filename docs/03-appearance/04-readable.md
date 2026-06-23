# readable Theme

Dyslexia-friendly typography. Uses the OpenDyslexic font, increased letter spacing, generous line height.

```bash
volcano ./docs --theme readable --url="https://example.com"
```

## Light

![readable theme, light mode](/images/themes/readable-light.png)

## Dark

![readable theme, dark mode](/images/themes/readable-dark.png)

## Features

- **OpenDyslexic** font — letterforms designed to reduce mirroring/swapping confusion
- Wider letter spacing
- Larger base font size
- Higher line height
- Higher contrast palette
- Dark mode tuned for low-light reading

## Best For

Readers with dyslexia, ADHD, or anyone who finds standard documentation typography fatiguing. Drop-in alternative to the `docs` theme when accessibility matters.

## Combine With Other Features

The readable theme respects the same flags as every other theme — you can opt into search, breadcrumbs, top nav, etc.:

```bash
volcano ./docs \
  --theme readable \
  --search \
  --breadcrumbs \
  --url="https://docs.example.com"
```

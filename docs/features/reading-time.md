# Reading Time

Volcano calculates and displays estimated reading time for each page.

## How It Works

Reading time appears in the page header, showing how long the content takes to read:

```
5 min read
```

### Calculation Method

Volcano uses different reading speeds for different content types:

| Content Type | Words Per Minute |
|--------------|------------------|
| Regular text | 225 WPM |
| Code blocks | 100 WPM |

Code blocks are read slower because readers often study them carefully.

### Formula

```
Reading Time = (Regular Words / 225) + (Code Words / 100)
```

The result is rounded to the nearest minute, with a minimum of 1 minute.

## Examples

### Text-Heavy Page

A 1,125-word article with no code:

```
1125 words / 225 WPM = 5 minutes
```

**Display:** `5 min read`

### Code-Heavy Page

A page with 450 words of text and 200 words of code:

```
(450 / 225) + (200 / 100) = 2 + 2 = 4 minutes
```

**Display:** `4 min read`

### Short Page

A 50-word page:

```
50 / 225 = 0.22 minutes → rounds to 1 minute (minimum)
```

**Display:** `1 min read`

## Display Location

Reading time appears in the page header area, typically near the title or metadata section.

The exact position depends on your theme:

- **docs theme** — Below the page title
- **blog theme** — In the article header with date
- **vanilla theme** — Structure only, style as needed

## Styling

Reading time uses the `.reading-time` CSS class:

```css
.reading-time {
  color: var(--color-text-muted);
  font-size: 0.875rem;
}
```

### Custom Styling Examples

#### Minimal

```css
.reading-time {
  opacity: 0.6;
  font-size: 0.8em;
}
```

#### With Icon

```css
.reading-time::before {
  content: "⏱ ";
}
```

#### Badge Style

```css
.reading-time {
  display: inline-block;
  background: var(--color-bg-muted);
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
}
```

## Word Counting

### What's Counted

- Regular text content
- Headings
- List items
- Table content
- Code block content (at slower rate)
- Inline code

### What's Not Counted

- HTML tags
- Markdown syntax characters
- Image alt text
- Link URLs

## Accuracy Considerations

Reading time is an estimate. Actual reading time varies based on:

- Reader's familiarity with the topic
- Content complexity
- Code comprehension speed
- Reading environment

The 225 WPM baseline is based on average adult reading speed for technical content.

## Best Practices

### Write Scannable Content

Help readers match or beat the estimate:

- Use clear headings
- Keep paragraphs short
- Use bullet points
- Add code comments

### Consider Your Audience

Technical documentation readers often:
- Skim familiar sections
- Study code examples carefully
- Jump between sections

The reading time represents a thorough read, not a skim.

### Don't Pad Content

Reading time reflects content length. Focus on:
- Clear, concise writing
- Relevant examples
- Necessary detail only

## Disabling Reading Time

Reading time is automatically calculated. To hide it, use CSS:

```css
.reading-time {
  display: none;
}
```

Or target specific pages:

```css
/* Hide on index pages */
.auto-index-page .reading-time {
  display: none;
}
```

## Related

- [[theming]] — Customize reading time appearance
- [[markdown-syntax]] — Content formatting
- [[guides/customizing-appearance]] — CSS customization

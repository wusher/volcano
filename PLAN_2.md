# Volcano - Extended Feature Stories (Phase 2)

## Overview

This document extends the original PLAN.md with 24 additional feature stories that align with Volcano's core philosophy:

- **Simplicity first**: Minimal dependencies, embedded assets
- **Documentation-focused**: Features that enhance technical documentation sites
- **Black and white aesthetic**: Consistent with the monochrome color scheme
- **No external runtime dependencies**: Everything embedded in the binary
- **Progressive enhancement**: Features work without JavaScript where possible

These stories assume completion of Stories 1-13 from PLAN.md.

---

## Story 14: Table of Contents Component

**Goal**: Auto-generate a floating table of contents from page headings.

### Acceptance Criteria

- [ ] Parse rendered HTML to extract all h2-h4 headings
- [ ] Generate nested list structure reflecting heading hierarchy
- [ ] Display TOC in right sidebar on desktop (>= 1024px)
- [ ] Hide TOC on mobile and tablet (< 1024px)
- [ ] Highlight current section based on scroll position
- [ ] Smooth scroll to heading when TOC item clicked
- [ ] Skip TOC generation for pages with fewer than 3 headings
- [ ] Add `--no-toc` flag to disable TOC globally
- [ ] TOC container fixed position, scrollable if too long

### Data Structures

```go
type TOCItem struct {
    ID       string     // Heading ID for anchor link
    Text     string     // Heading text content
    Level    int        // 2, 3, or 4
    Children []*TOCItem // Nested headings
}

type PageTOC struct {
    Items    []*TOCItem
    MinItems int // Minimum headings to show TOC (default: 3)
}
```

### HTML Structure

```html
<aside class="toc-sidebar" aria-label="Table of contents">
  <nav class="toc">
    <h2 class="toc-title">On this page</h2>
    <ul>
      <li>
        <a href="#installation" class="active">Installation</a>
        <ul>
          <li><a href="#requirements">Requirements</a></li>
          <li><a href="#quick-start">Quick Start</a></li>
        </ul>
      </li>
      <li><a href="#configuration">Configuration</a></li>
    </ul>
  </nav>
</aside>
```

### JavaScript Behavior

```javascript
// Intersection Observer for scroll spy
const observer = new IntersectionObserver((entries) => {
  entries.forEach(entry => {
    if (entry.isIntersecting) {
      // Update active TOC item
      document.querySelectorAll('.toc a').forEach(a => a.classList.remove('active'));
      document.querySelector(`.toc a[href="#${entry.target.id}"]`)?.classList.add('active');
    }
  });
}, { rootMargin: '-80px 0px -80% 0px' });

document.querySelectorAll('h2, h3, h4').forEach(h => observer.observe(h));
```

### File Location

```
volcano/
├── internal/
│   └── toc/
│       ├── extractor.go   # Heading extraction from HTML
│       ├── builder.go     # TOC tree building
│       └── toc_test.go
```

### Styling Notes

- TOC width: 220px
- Font size: 14px
- Active item: bold, with left border indicator
- Nested items indented 16px per level
- Max height with overflow-y: auto for long TOCs

---

## Story 15: Breadcrumb Navigation

**Goal**: Display hierarchical path from root to current page above content.

### Acceptance Criteria

- [ ] Show breadcrumb trail above page title
- [ ] Each breadcrumb segment is a clickable link (except current page)
- [ ] Use clean labels (same transformation as nav tree)
- [ ] Separator: `/` or `›` character between segments
- [ ] Root shows site title or "Home"
- [ ] Current page shown but not linked (text only)
- [ ] Hide breadcrumbs on root index page
- [ ] Accessible with proper `nav` and `aria-label` attributes
- [ ] Structured data (JSON-LD) for SEO

### HTML Structure

```html
<nav class="breadcrumbs" aria-label="Breadcrumb">
  <ol itemscope itemtype="https://schema.org/BreadcrumbList">
    <li itemprop="itemListElement" itemscope itemtype="https://schema.org/ListItem">
      <a itemprop="item" href="/"><span itemprop="name">Home</span></a>
      <meta itemprop="position" content="1" />
    </li>
    <li itemprop="itemListElement" itemscope itemtype="https://schema.org/ListItem">
      <a itemprop="item" href="/guides/"><span itemprop="name">Guides</span></a>
      <meta itemprop="position" content="2" />
    </li>
    <li itemprop="itemListElement" itemscope itemtype="https://schema.org/ListItem">
      <span itemprop="name" aria-current="page">Getting Started</span>
      <meta itemprop="position" content="3" />
    </li>
  </ol>
</nav>
```

### Data Structures

```go
type Breadcrumb struct {
    Label   string
    URL     string
    Current bool
}

func BuildBreadcrumbs(page *Page, tree *SiteTree) []Breadcrumb
```

### File Location

```
volcano/
├── internal/
│   └── navigation/
│       ├── breadcrumbs.go
│       └── breadcrumbs_test.go
```

### Styling Notes

- Positioned above page title with margin-bottom
- Muted text color for non-current items
- Items separated by `›` with padding
- Font size: 14px
- No wrapping; truncate middle segments on mobile if needed

---

## Story 16: Previous/Next Page Navigation

**Goal**: Add sequential navigation links at the bottom of each page.

### Acceptance Criteria

- [ ] Display "Previous" and "Next" links at page bottom
- [ ] Order determined by tree structure (depth-first, alphabetical)
- [ ] Show page title as link text
- [ ] Include parent folder name for context if different from current
- [ ] Arrow icons indicating direction (← Previous, Next →)
- [ ] First page has no "Previous" link
- [ ] Last page has no "Next" link
- [ ] Keyboard accessible (part of normal tab order)
- [ ] `--no-pagination` flag to disable

### HTML Structure

```html
<nav class="page-nav" aria-label="Page navigation">
  <a href="/guides/installation/" class="page-nav-prev">
    <span class="page-nav-label">Previous</span>
    <span class="page-nav-title">← Installation</span>
  </a>
  <a href="/guides/advanced/" class="page-nav-next">
    <span class="page-nav-label">Next</span>
    <span class="page-nav-title">Advanced Usage →</span>
  </a>
</nav>
```

### Data Structures

```go
type PageNavigation struct {
    Previous *NavLink
    Next     *NavLink
}

type NavLink struct {
    Title   string // Page title
    URL     string
    Section string // Parent folder name (optional)
}

func BuildPageNavigation(page *Page, allPages []*Page) PageNavigation
```

### File Location

```
volcano/
├── internal/
│   └── navigation/
│       ├── pagination.go
│       └── pagination_test.go
```

### Styling Notes

- Flexbox layout: previous left-aligned, next right-aligned
- Border-top separator from content
- Padding: 24px top, 0 bottom
- Hover state: background color change
- Full width on mobile (stacked vertically)

---

## Story 17: Heading Anchor Links

**Goal**: Add clickable anchor links to all headings for easy deep linking.

### Acceptance Criteria

- [ ] Generate unique IDs for all h1-h6 headings
- [ ] ID derived from heading text (slugified, lowercase, hyphenated)
- [ ] Handle duplicate headings by appending `-1`, `-2`, etc.
- [ ] Show link icon on heading hover
- [ ] Click copies URL with anchor to clipboard (optional)
- [ ] Anchor link appears before or after heading text
- [ ] Links are accessible (proper aria-label)
- [ ] Smooth scroll offset accounts for fixed header (if any)

### HTML Output

```html
<h2 id="installation">
  <a href="#installation" class="heading-anchor" aria-label="Link to Installation section">
    <svg class="anchor-icon"><!-- Link icon --></svg>
  </a>
  Installation
</h2>
```

### Slug Generation

```go
func Slugify(text string) string {
    // "Getting Started with Go" → "getting-started-with-go"
    // Handle unicode, special characters
    // Ensure uniqueness within page
}

type HeadingID struct {
    Original string
    Slug     string
    Count    int // For deduplication
}
```

### File Location

```
volcano/
├── internal/
│   └── markdown/
│       ├── headings.go     # ID generation, anchor injection
│       └── slugify.go      # Text to slug conversion
```

### Styling Notes

- Anchor icon: 16x16px, muted color
- Opacity: 0 by default, 1 on heading hover
- Transition: opacity 0.2s
- Icon positioned to left of heading with negative margin
- Focus visible outline for keyboard users

---

## Story 18: External Link Indicators

**Goal**: Visually distinguish external links and configure their behavior.

### Acceptance Criteria

- [ ] Detect external links (different domain from site)
- [ ] Add small external link icon after link text
- [ ] Add `target="_blank"` to external links
- [ ] Add `rel="noopener noreferrer"` for security
- [ ] Internal links unchanged (same domain)
- [ ] Configurable via `--external-links-new-tab` flag (default: true)
- [ ] Icon does not appear on image links
- [ ] Screen reader text: "opens in new tab"

### Implementation

```go
func ProcessExternalLinks(html string, siteURL string) string {
    // Parse HTML
    // Find all <a> tags
    // Check if href is external
    // Add attributes and icon
    // Return modified HTML
}

func IsExternalURL(href string, siteURL string) bool {
    // Handle relative URLs, protocol-relative, etc.
}
```

### HTML Output

```html
<!-- Internal link (unchanged) -->
<a href="/guides/intro/">Introduction</a>

<!-- External link (modified) -->
<a href="https://golang.org" target="_blank" rel="noopener noreferrer">
  Go Programming Language
  <svg class="external-icon" aria-hidden="true"><!-- icon --></svg>
  <span class="sr-only">(opens in new tab)</span>
</a>
```

### File Location

```
volcano/
├── internal/
│   └── markdown/
│       ├── links.go        # Link processing
│       └── links_test.go
```

### Styling Notes

- Icon: 12x12px, inline with text baseline
- Margin-left: 4px
- Color: matches link color
- Icon SVG: simple arrow pointing up-right

---

## Story 19: Code Block Copy Button

**Goal**: Add one-click copy-to-clipboard button on code blocks.

### Acceptance Criteria

- [ ] Copy button appears on all fenced code blocks
- [ ] Button positioned in top-right corner of code block
- [ ] Shows "Copy" text or copy icon
- [ ] On click: copies code content to clipboard
- [ ] Visual feedback: button changes to "Copied!" for 2 seconds
- [ ] Works with syntax-highlighted code blocks
- [ ] Excludes line numbers if present (copy only code)
- [ ] Graceful fallback if clipboard API unavailable
- [ ] Keyboard accessible

### HTML Structure

```html
<div class="code-block">
  <button class="copy-button" aria-label="Copy code to clipboard">
    <svg class="copy-icon"><!-- Copy icon --></svg>
    <span class="copy-text">Copy</span>
  </button>
  <pre><code class="language-go">package main

func main() {
    fmt.Println("Hello")
}</code></pre>
</div>
```

### JavaScript Implementation

```javascript
document.querySelectorAll('.copy-button').forEach(button => {
  button.addEventListener('click', async () => {
    const code = button.parentElement.querySelector('code').textContent;
    try {
      await navigator.clipboard.writeText(code);
      button.querySelector('.copy-text').textContent = 'Copied!';
      button.classList.add('copied');
      setTimeout(() => {
        button.querySelector('.copy-text').textContent = 'Copy';
        button.classList.remove('copied');
      }, 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  });
});
```

### File Location

```
volcano/
├── internal/
│   └── markdown/
│       └── codeblock.go    # Code block wrapper injection
├── internal/
│   └── templates/
│       └── scripts.js      # Copy functionality
```

### Styling Notes

- Button background: semi-transparent gray
- Position: absolute, top-right with padding
- Hover: darker background
- Copied state: green checkmark icon, green text
- Button only visible on code block hover (desktop)
- Always visible on touch devices

---

## Story 20: Keyboard Navigation Shortcuts

**Goal**: Enable power users to navigate the site using keyboard shortcuts.

### Acceptance Criteria

- [ ] Global shortcuts (work anywhere on page):
  - `/` or `Ctrl+K`: Focus search (if search exists) or nav
  - `t`: Toggle dark/light theme
  - `n`: Next page
  - `p`: Previous page
  - `h`: Go to home page
  - `?`: Show keyboard shortcuts modal
- [ ] Navigation shortcuts (when nav focused):
  - `↑/↓`: Move through nav items
  - `Enter`: Navigate to selected item
  - `→`: Expand folder
  - `←`: Collapse folder
- [ ] Content shortcuts:
  - `Escape`: Close any open modal/drawer
  - `g` then `t`: Scroll to top
  - `g` then `b`: Scroll to bottom
- [ ] Disable shortcuts when typing in input/textarea
- [ ] `--no-shortcuts` flag to disable

### HTML Structure

```html
<!-- Shortcuts help modal -->
<dialog id="shortcuts-modal" class="shortcuts-modal">
  <h2>Keyboard Shortcuts</h2>
  <dl class="shortcuts-list">
    <div class="shortcut-group">
      <dt><kbd>t</kbd></dt>
      <dd>Toggle theme</dd>
    </div>
    <div class="shortcut-group">
      <dt><kbd>n</kbd></dt>
      <dd>Next page</dd>
    </div>
    <!-- ... -->
  </dl>
  <button class="close-modal">Close</button>
</dialog>
```

### JavaScript Implementation

```javascript
const shortcuts = {
  't': () => toggleTheme(),
  'n': () => navigateToNext(),
  'p': () => navigateToPrevious(),
  'h': () => window.location.href = '/',
  '?': () => showShortcutsModal(),
};

document.addEventListener('keydown', (e) => {
  // Skip if in input/textarea
  if (['INPUT', 'TEXTAREA'].includes(e.target.tagName)) return;

  const handler = shortcuts[e.key];
  if (handler) {
    e.preventDefault();
    handler();
  }
});
```

### File Location

```
volcano/
├── internal/
│   └── templates/
│       ├── shortcuts.html   # Shortcuts modal
│       └── shortcuts.js     # Keyboard handling
```

### Styling Notes

- Modal: centered, max-width 500px
- Backdrop: semi-transparent black
- `<kbd>` elements: monospace, bordered, subtle background
- Two-column layout for shortcuts list
- Close on Escape key or backdrop click

---

## Story 21: Print Stylesheet

**Goal**: Provide clean, printer-friendly output for documentation pages.

### Acceptance Criteria

- [ ] Print styles applied via `@media print`
- [ ] Hide elements not needed for print:
  - Navigation sidebar
  - TOC sidebar
  - Dark mode toggle
  - Copy buttons
  - Mobile menu
  - Previous/Next navigation
- [ ] Show elements useful for print:
  - Full URLs for links (via `content: attr(href)`)
  - Page title as header
  - Site name in footer
- [ ] Typography optimized for print:
  - Serif font for body text (optional)
  - Black text on white background
  - Appropriate margins
  - Page break control (avoid orphans/widows)
- [ ] Code blocks: preserve formatting, prevent page breaks inside
- [ ] Images: reasonable max-width, page break handling
- [ ] Optional: include page URL as footer

### CSS Implementation

```css
@media print {
  /* Hide UI elements */
  .sidebar,
  .toc-sidebar,
  .theme-toggle,
  .mobile-menu-btn,
  .copy-button,
  .page-nav,
  .heading-anchor {
    display: none !important;
  }

  /* Reset colors */
  body {
    background: white !important;
    color: black !important;
  }

  /* Content fills page */
  .content {
    width: 100% !important;
    max-width: none !important;
    margin: 0 !important;
    padding: 0 !important;
  }

  /* Show URLs for links */
  a[href^="http"]::after {
    content: " (" attr(href) ")";
    font-size: 0.8em;
    color: #666;
  }

  /* Prevent page breaks */
  pre, code, img, table {
    page-break-inside: avoid;
  }

  h1, h2, h3, h4 {
    page-break-after: avoid;
  }

  /* Page margins */
  @page {
    margin: 2cm;
  }
}
```

### File Location

```
volcano/
├── internal/
│   └── styles/
│       └── print.css
```

### Testing Notes

- Test with browser print preview
- Verify code blocks don't break mid-block
- Check URL display for external links
- Confirm all interactive elements hidden

---

## Story 22: Reading Time Indicator

**Goal**: Display estimated reading time for each page.

### Acceptance Criteria

- [ ] Calculate reading time from content word count
- [ ] Use average reading speed: 200-250 words per minute
- [ ] Display format: "X min read"
- [ ] Show below page title, near metadata area
- [ ] Account for code blocks (reduced reading speed)
- [ ] Minimum display: "1 min read"
- [ ] Round to nearest minute
- [ ] `--no-reading-time` flag to disable

### Implementation

```go
type ReadingTime struct {
    Minutes int
    Words   int
}

func CalculateReadingTime(content string) ReadingTime {
    // Strip HTML tags
    // Count words
    // Account for code blocks (slower reading)
    // words / 225 = minutes
    // Round and ensure minimum of 1
}

const wordsPerMinute = 225
const codeWordsPerMinute = 100 // Code is read slower
```

### HTML Output

```html
<article>
  <header class="page-header">
    <h1>Getting Started</h1>
    <div class="page-meta">
      <span class="reading-time">
        <svg class="clock-icon"><!-- clock --></svg>
        5 min read
      </span>
    </div>
  </header>
  <!-- content -->
</article>
```

### File Location

```
volcano/
├── internal/
│   └── content/
│       ├── readingtime.go
│       └── readingtime_test.go
```

### Styling Notes

- Font size: 14px
- Color: muted/secondary text color
- Clock icon: 14x14px
- Positioned below title, left-aligned
- Subtle, not distracting from content

---

## Story 23: Last Modified Display

**Goal**: Show when each page was last updated.

### Acceptance Criteria

- [ ] Detect last modified date from:
  1. Git commit date (preferred)
  2. File system modification time (fallback)
- [ ] Display format: "Last updated: Month Day, Year"
- [ ] Relative time option: "Last updated: 3 days ago"
- [ ] Show in page header/metadata area
- [ ] `--last-modified` flag to enable (off by default)
- [ ] `--last-modified-format` to choose format (date|relative)
- [ ] Handle files not in git gracefully
- [ ] Generate structured data for SEO (dateModified)

### Implementation

```go
type ModifiedDate struct {
    Time     time.Time
    Source   string // "git" or "filesystem"
    Relative string // "3 days ago"
    Absolute string // "January 5, 2025"
}

func GetLastModified(filePath string) ModifiedDate {
    // Try git first
    gitDate, err := getGitModifiedDate(filePath)
    if err == nil {
        return gitDate
    }
    // Fallback to filesystem
    return getFileModifiedDate(filePath)
}

func getGitModifiedDate(path string) (ModifiedDate, error) {
    // git log -1 --format=%cI <path>
}
```

### HTML Output

```html
<div class="page-meta">
  <span class="last-modified">
    <svg class="calendar-icon"><!-- calendar --></svg>
    Last updated: January 5, 2025
  </span>
</div>

<!-- Or with relative time -->
<span class="last-modified" title="January 5, 2025">
  Last updated 3 days ago
</span>
```

### File Location

```
volcano/
├── internal/
│   └── content/
│       ├── modified.go
│       └── modified_test.go
```

### Styling Notes

- Same styling as reading time
- Calendar icon: 14x14px
- Full date shown on hover (via title) for relative format
- Grouped with other metadata items

---

## Story 24: Scroll Progress Indicator

**Goal**: Visual indicator showing reading progress through the current page.

### Acceptance Criteria

- [ ] Thin progress bar at top of viewport
- [ ] Width represents scroll progress (0% at top, 100% at bottom)
- [ ] Fixed position, always visible
- [ ] Smooth animation as user scrolls
- [ ] Color: accent color or primary text color
- [ ] Height: 2-3px
- [ ] `--no-progress` flag to disable
- [ ] Accessible: decorative only, hidden from screen readers

### HTML Structure

```html
<div class="scroll-progress" aria-hidden="true">
  <div class="scroll-progress-bar" style="width: 0%"></div>
</div>
```

### JavaScript Implementation

```javascript
const progressBar = document.querySelector('.scroll-progress-bar');

function updateProgress() {
  const scrollTop = window.scrollY;
  const docHeight = document.documentElement.scrollHeight - window.innerHeight;
  const progress = docHeight > 0 ? (scrollTop / docHeight) * 100 : 0;
  progressBar.style.width = `${progress}%`;
}

// Use passive listener for performance
window.addEventListener('scroll', updateProgress, { passive: true });
updateProgress(); // Initial state
```

### File Location

```
volcano/
├── internal/
│   └── templates/
│       └── progress.js
├── internal/
│   └── styles/
│       └── progress.css
```

### Styling Notes

```css
.scroll-progress {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: var(--bg-secondary);
  z-index: 1000;
}

.scroll-progress-bar {
  height: 100%;
  background: var(--text-primary);
  width: 0;
  transition: width 0.1s linear;
}
```

---

## Story 25: Back to Top Button

**Goal**: Floating button to quickly scroll back to the top of the page.

### Acceptance Criteria

- [ ] Floating button appears after scrolling down (> 300px)
- [ ] Fixed position in bottom-right corner
- [ ] Smooth scroll to top when clicked
- [ ] Fade in/out animation
- [ ] Shows arrow-up icon
- [ ] Accessible: proper aria-label, keyboard focusable
- [ ] Hidden when near top of page
- [ ] `--no-back-to-top` flag to disable

### HTML Structure

```html
<button
  class="back-to-top"
  aria-label="Scroll to top"
  hidden
>
  <svg class="arrow-up-icon"><!-- Up arrow --></svg>
</button>
```

### JavaScript Implementation

```javascript
const backToTop = document.querySelector('.back-to-top');
const showThreshold = 300;

function toggleBackToTop() {
  if (window.scrollY > showThreshold) {
    backToTop.hidden = false;
    backToTop.classList.add('visible');
  } else {
    backToTop.classList.remove('visible');
  }
}

backToTop.addEventListener('click', () => {
  window.scrollTo({ top: 0, behavior: 'smooth' });
});

window.addEventListener('scroll', toggleBackToTop, { passive: true });
```

### File Location

```
volcano/
├── internal/
│   └── templates/
│       └── backtotop.js
├── internal/
│   └── styles/
│       └── backtotop.css
```

### Styling Notes

```css
.back-to-top {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  cursor: pointer;
  opacity: 0;
  transform: translateY(10px);
  transition: opacity 0.3s, transform 0.3s;
  z-index: 100;
}

.back-to-top.visible {
  opacity: 1;
  transform: translateY(0);
}

.back-to-top:hover {
  background: var(--bg-primary);
  border-color: var(--text-secondary);
}
```

---

## Story 26: SEO Meta Tags Generation

**Goal**: Generate comprehensive meta tags for search engine optimization.

### Acceptance Criteria

- [ ] Generate for each page:
  - `<title>` with page title and site name
  - `<meta name="description">` from first paragraph or custom
  - `<meta name="robots">` (index, follow)
  - `<link rel="canonical">` with full URL
  - `<meta name="author">` if configured
- [ ] Global configuration:
  - `--site-url` base URL for canonical links
  - `--description` default site description
  - `--author` site author name
- [ ] Auto-generate description from first 160 characters of content
- [ ] Strip markdown/HTML from generated descriptions
- [ ] Handle special characters properly (HTML entities)

### Implementation

```go
type SEOMeta struct {
    Title       string
    Description string
    Canonical   string
    Robots      string
    Author      string
}

func GenerateSEOMeta(page *Page, config *Config) SEOMeta {
    description := page.Description
    if description == "" {
        description = extractFirstParagraph(page.Content, 160)
    }
    return SEOMeta{
        Title:       fmt.Sprintf("%s - %s", page.Title, config.SiteTitle),
        Description: description,
        Canonical:   config.SiteURL + page.URLPath,
        Robots:      "index, follow",
        Author:      config.Author,
    }
}
```

### HTML Output

```html
<head>
  <title>Getting Started - My Documentation</title>
  <meta name="description" content="Learn how to install and configure Volcano for your documentation site.">
  <meta name="robots" content="index, follow">
  <meta name="author" content="Documentation Team">
  <link rel="canonical" href="https://docs.example.com/guides/getting-started/">
</head>
```

### File Location

```
volcano/
├── internal/
│   └── seo/
│       ├── meta.go
│       ├── description.go
│       └── seo_test.go
```

---

## Story 27: Open Graph Support

**Goal**: Add Open Graph and Twitter Card meta tags for rich social sharing.

### Acceptance Criteria

- [ ] Generate Open Graph tags:
  - `og:title` - page title
  - `og:description` - page description
  - `og:type` - "article" for pages
  - `og:url` - canonical URL
  - `og:site_name` - site title
  - `og:image` - default or page-specific image
- [ ] Generate Twitter Card tags:
  - `twitter:card` - "summary_large_image" or "summary"
  - `twitter:title`, `twitter:description`
  - `twitter:image`
- [ ] Configuration flags:
  - `--og-image` default social sharing image URL
  - `--twitter-handle` Twitter username
- [ ] Support first image in content as og:image fallback
- [ ] Image dimensions recommendations in docs

### HTML Output

```html
<head>
  <!-- Open Graph -->
  <meta property="og:title" content="Getting Started">
  <meta property="og:description" content="Learn how to install and configure Volcano.">
  <meta property="og:type" content="article">
  <meta property="og:url" content="https://docs.example.com/guides/getting-started/">
  <meta property="og:site_name" content="My Documentation">
  <meta property="og:image" content="https://docs.example.com/og-image.png">

  <!-- Twitter Card -->
  <meta name="twitter:card" content="summary_large_image">
  <meta name="twitter:title" content="Getting Started">
  <meta name="twitter:description" content="Learn how to install and configure Volcano.">
  <meta name="twitter:image" content="https://docs.example.com/og-image.png">
</head>
```

### File Location

```
volcano/
├── internal/
│   └── seo/
│       ├── opengraph.go
│       └── twitter.go
```

---

## Story 28: Custom Favicon Support

**Goal**: Allow custom favicon configuration for site branding.

### Acceptance Criteria

- [ ] `--favicon` flag to specify favicon file path
- [ ] Support formats: .ico, .png, .svg
- [ ] Copy favicon to output root directory
- [ ] Generate appropriate `<link>` tags:
  - `<link rel="icon" type="image/x-icon" href="/favicon.ico">`
  - `<link rel="icon" type="image/png" href="/favicon.png">`
  - `<link rel="icon" type="image/svg+xml" href="/favicon.svg">`
- [ ] Support Apple touch icon via `--apple-touch-icon` flag
- [ ] Default: no favicon (no broken links)
- [ ] Validate favicon file exists during build

### Implementation

```go
type FaviconConfig struct {
    IconPath       string // Path to favicon file
    AppleTouchIcon string // Path to Apple touch icon
}

func ProcessFavicon(config *FaviconConfig, outputDir string) ([]FaviconLink, error) {
    // Validate file exists
    // Determine MIME type from extension
    // Copy to output directory
    // Return link tags to include
}

type FaviconLink struct {
    Rel  string // "icon", "apple-touch-icon"
    Type string // "image/png", etc.
    Href string // URL path
}
```

### HTML Output

```html
<head>
  <link rel="icon" type="image/svg+xml" href="/favicon.svg">
  <link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
  <link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
</head>
```

### File Location

```
volcano/
├── internal/
│   └── assets/
│       ├── favicon.go
│       └── favicon_test.go
```

---

## Story 29: Admonition/Callout Blocks

**Goal**: Support note, warning, tip, and other callout blocks in markdown.

### Acceptance Criteria

- [ ] Custom markdown syntax for callouts:
  ```markdown
  :::note
  This is a note callout.
  :::

  :::warning
  This is a warning.
  :::
  ```
- [ ] Supported types:
  - `note` - blue, info icon
  - `tip` - green, lightbulb icon
  - `warning` - yellow, warning icon
  - `danger` - red, error icon
  - `info` - gray, info icon
- [ ] Custom title support: `:::note Custom Title`
- [ ] Render with appropriate styling and icon
- [ ] Accessible: proper ARIA role="note" or role="alert"
- [ ] Nested content support (code blocks, lists inside)
- [ ] Graceful degradation if syntax not recognized

### Markdown Syntax

```markdown
:::note
This is a basic note.
:::

:::warning Be Careful
This action cannot be undone. Make sure to backup your data first.
:::

:::tip Performance Tip
You can improve build speed by enabling caching.
:::

:::danger
Do not expose your API keys in client-side code.
:::
```

### HTML Output

```html
<div class="admonition admonition-warning" role="note">
  <div class="admonition-heading">
    <svg class="admonition-icon"><!-- warning icon --></svg>
    <span class="admonition-title">Be Careful</span>
  </div>
  <div class="admonition-content">
    <p>This action cannot be undone. Make sure to backup your data first.</p>
  </div>
</div>
```

### Implementation

```go
// Custom Goldmark extension
type AdmonitionExtension struct{}

func (e *AdmonitionExtension) Extend(m goldmark.Markdown) {
    // Register custom block parser
    // Register custom renderer
}

type Admonition struct {
    Type    string // note, warning, tip, danger, info
    Title   string // Custom or default from type
    Content []byte
}
```

### File Location

```
volcano/
├── internal/
│   └── markdown/
│       ├── admonition.go      # Parser and renderer
│       └── admonition_test.go
├── internal/
│   └── styles/
│       └── admonition.css
```

### Styling Notes

```css
.admonition {
  padding: 16px;
  margin: 16px 0;
  border-left: 4px solid;
  border-radius: 4px;
}

.admonition-note {
  background: #f0f7ff;
  border-color: #0066cc;
}

.admonition-warning {
  background: #fff8e6;
  border-color: #cc8800;
}

/* Dark mode variants */
[data-theme="dark"] .admonition-note {
  background: #1a2a3a;
  border-color: #4499ff;
}
```

**Note**: Colors break the black/white-only rule but can be made monochrome:
- note: light gray background
- warning: medium gray with thicker border
- danger: dark gray background, white text

---

## Story 30: Code Line Highlighting

**Goal**: Allow highlighting specific lines in code blocks for emphasis.

### Acceptance Criteria

- [ ] Markdown syntax to specify highlighted lines:
  ```markdown
  ```go {3,5-7}
  package main

  import "fmt"  // highlighted

  func main() {  // highlighted
      fmt.Println("Hello")  // highlighted
  }  // highlighted
  ```
  ```
- [ ] Support:
  - Single lines: `{3}`
  - Multiple lines: `{3,5,9}`
  - Ranges: `{3-7}`
  - Combined: `{1,3-5,10}`
- [ ] Visual distinction: background highlight on specified lines
- [ ] Works with syntax highlighting
- [ ] Line numbers displayed correctly alongside highlights

### Implementation

```go
type CodeBlockMeta struct {
    Language       string
    HighlightLines []int      // Parsed from {1,3-5}
    ShowLineNumbers bool      // Optional
}

func ParseCodeBlockMeta(info string) CodeBlockMeta {
    // Parse "go {3,5-7}" into language and highlight lines
}

func ParseLineSpec(spec string) []int {
    // "3,5-7,10" → [3, 5, 6, 7, 10]
}
```

### HTML Output

```html
<pre class="code-block"><code class="language-go"><span class="line">package main</span>
<span class="line">  </span>
<span class="line highlight">import "fmt"</span>
<span class="line">  </span>
<span class="line highlight">func main() {</span>
<span class="line highlight">    fmt.Println("Hello")</span>
<span class="line highlight">}</span>
</code></pre>
```

### File Location

```
volcano/
├── internal/
│   └── markdown/
│       ├── codeblock.go       # Extended for line highlighting
│       └── codeblock_test.go
```

### Styling Notes

```css
.code-block .line.highlight {
  background: rgba(255, 255, 0, 0.15); /* Light mode */
  display: inline-block;
  width: 100%;
  margin: 0 -16px;
  padding: 0 16px;
}

[data-theme="dark"] .code-block .line.highlight {
  background: rgba(255, 255, 255, 0.1);
}
```

---

## Story 31: Smooth Scroll Behavior

**Goal**: Enable smooth scrolling throughout the site for polished UX.

### Acceptance Criteria

- [ ] CSS `scroll-behavior: smooth` on `html` element
- [ ] Smooth scrolling for:
  - Anchor links (heading navigation)
  - TOC navigation
  - Back to top button
  - Breadcrumb navigation (within page)
- [ ] Scroll offset for fixed headers (if applicable)
- [ ] Respect `prefers-reduced-motion` media query
- [ ] Keyboard navigation respects smooth scroll
- [ ] `--no-smooth-scroll` flag to disable

### CSS Implementation

```css
html {
  scroll-behavior: smooth;
}

/* Respect user preference for reduced motion */
@media (prefers-reduced-motion: reduce) {
  html {
    scroll-behavior: auto;
  }
}

/* Offset for fixed header (if using one) */
:target {
  scroll-margin-top: 80px;
}

/* Alternative: CSS scroll-margin on headings */
h1, h2, h3, h4, h5, h6 {
  scroll-margin-top: 80px;
}
```

### JavaScript Enhancement

```javascript
// For browsers with poor CSS smooth scroll support
// or for programmatic scrolling with offset
function smoothScrollTo(targetId, offset = 80) {
  const target = document.getElementById(targetId);
  if (!target) return;

  const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  const top = target.getBoundingClientRect().top + window.scrollY - offset;

  window.scrollTo({
    top,
    behavior: prefersReducedMotion ? 'auto' : 'smooth'
  });
}
```

### File Location

```
volcano/
├── internal/
│   └── styles/
│       └── base.css   # Add scroll-behavior rules
├── internal/
│   └── templates/
│       └── scroll.js  # Optional JS enhancement
```

---

## Story 32: Clickable Folder Navigation

**Goal**: Make folder names in the navigation tree clickable when they have associated content.

### Acceptance Criteria

- [ ] Folder names are clickable links when content is available
- [ ] Two patterns for folder content detection:
  1. **Index file pattern**: `posts/index.md` makes `posts/` folder clickable
  2. **Sibling file pattern**: `posts.md` next to `posts/` folder makes folder clickable
- [ ] Priority order when both exist:
  1. Index file takes precedence (`posts/index.md`)
  2. Sibling file as fallback (`posts.md`)
- [ ] **Hide content file from tree**: The associated file should NOT appear as a separate item in the navigation
  - `posts/index.md` is hidden; only "Posts" folder shows (clickable)
  - `posts.md` (sibling) is hidden; only "Posts" folder shows (clickable)
- [ ] URL structure:
  - `posts/index.md` → `/posts/`
  - `posts.md` (with `posts/` folder) → `/posts/`
- [ ] Visual distinction:
  - Clickable folders: folder name is a link (underline on hover)
  - Non-clickable folders: plain text, only chevron is interactive
- [ ] Expand/collapse toggle remains separate from folder link
- [ ] Keyboard accessible: folder link and toggle are separate tab stops
- [ ] Screen reader announces folder as link when clickable

### Directory Examples

```
# Pattern 1: Index file
docs/
├── guides/
│   ├── index.md        ← Hidden from nav; makes "Guides" folder clickable → /guides/
│   ├── getting-started.md
│   └── advanced.md

Nav tree shows:
  ▼ Guides          ← clickable link to /guides/
      Getting Started
      Advanced

# Pattern 2: Sibling file
docs/
├── api.md              ← Hidden from nav; makes "Api" folder clickable → /api/
├── api/
│   ├── endpoints.md
│   └── authentication.md

Nav tree shows:
  ▼ Api             ← clickable link to /api/
      Endpoints
      Authentication

# Pattern 3: Both exist (index takes precedence)
docs/
├── tutorials.md        ← Hidden (sibling consumed by folder)
├── tutorials/
│   ├── index.md        ← Hidden; makes "Tutorials" folder clickable → /tutorials/
│   └── basics.md

Nav tree shows:
  ▼ Tutorials       ← clickable link to /tutorials/ (from index.md)
      Basics

# Pattern 4: No associated content
docs/
├── reference/          ← "Reference" folder NOT clickable
│   ├── commands.md
│   └── config.md

Nav tree shows:
  ▼ Reference       ← plain text, not a link
      Commands
      Config
```

### Data Structures

```go
type TreeNode struct {
    Name        string
    Path        string
    SourcePath  string
    IsFolder    bool
    IsClickable bool        // True if folder has associated content
    IsHidden    bool        // True if file should not appear in nav tree
    ContentURL  string      // URL to navigate to when clicked
    ContentType string      // "index" or "sibling"
    Children    []*TreeNode
    Parent      *TreeNode
}

func (n *TreeNode) ResolveClickableFolder(siblings []*TreeNode) {
    if !n.IsFolder {
        return
    }

    // Check for index.md inside folder
    for _, child := range n.Children {
        if isIndexFile(child.Name) {
            n.IsClickable = true
            n.ContentURL = child.URLPath
            n.ContentType = "index"
            child.IsHidden = true  // Hide index from nav tree
            return
        }
    }

    // Check for sibling file with same name
    folderName := strings.TrimSuffix(n.Name, "/")
    for _, sibling := range siblings {
        if !sibling.IsFolder && cleanName(sibling.Name) == folderName {
            n.IsClickable = true
            n.ContentURL = sibling.URLPath
            n.ContentType = "sibling"
            sibling.IsHidden = true  // Hide sibling file from nav tree
            return
        }
    }
}

func isIndexFile(name string) bool {
    lower := strings.ToLower(name)
    return lower == "index.md" || lower == "readme.md"
}

// VisibleChildren returns only non-hidden children for nav rendering
func (n *TreeNode) VisibleChildren() []*TreeNode {
    var visible []*TreeNode
    for _, child := range n.Children {
        if !child.IsHidden {
            visible = append(visible, child)
        }
    }
    return visible
}
```

### HTML Structure

```html
<!-- Clickable folder (has index.md or sibling file) -->
<li role="treeitem" aria-expanded="true">
  <div class="folder clickable">
    <button class="toggle" aria-label="Collapse Guides section">
      <svg class="chevron"><!-- chevron --></svg>
    </button>
    <a href="/guides/" class="folder-link">Guides</a>
  </div>
  <ul role="group">
    <!-- children -->
  </ul>
</li>

<!-- Non-clickable folder (no associated content) -->
<li role="treeitem" aria-expanded="false">
  <div class="folder">
    <button class="toggle" aria-label="Expand Reference section">
      <svg class="chevron"><!-- chevron --></svg>
    </button>
    <span class="folder-name">Reference</span>
  </div>
  <ul role="group">
    <!-- children -->
  </ul>
</li>
```

### File Location

```
volcano/
├── internal/
│   └── tree/
│       ├── scanner.go      # Update to detect clickable folders
│       ├── clickable.go    # Folder content resolution logic
│       └── clickable_test.go
├── internal/
│   └── templates/
│       └── nav.html        # Update folder rendering
```

### Styling Notes

```css
.folder-link {
  color: var(--text-primary);
  text-decoration: none;
}

.folder-link:hover {
  text-decoration: underline;
}

.folder-name {
  color: var(--text-primary);
  cursor: default;
}

/* Differentiate clickable vs non-clickable folders */
.folder.clickable .folder-link {
  cursor: pointer;
}
```

### Edge Cases

- **Empty index.md**: Still makes folder clickable (content may be intentionally minimal)
- **Case sensitivity**: Match `index.md`, `Index.md`, `INDEX.MD` (case insensitive)
- **README.md**: Treat as equivalent to `index.md` for folder content
- **Multiple sibling matches**: Only exact name match counts (`api.md` for `api/`, not `api-v2.md`)
- **Nested folders**: Each folder independently checked for clickability
- **Root folder**: Always clickable (links to site home)
- **Hidden files still generate pages**: Files marked `IsHidden` are excluded from nav tree but still generate HTML output (accessible via folder click or direct URL)
- **Standalone file without folder**: `posts.md` without a `posts/` folder appears normally in tree (not hidden)

---

## Story 33: Navigation Tree Search

**Goal**: Add a search input that filters the navigation tree in real-time as the user types.

### Acceptance Criteria

- [ ] Search input field at top of navigation sidebar
- [ ] Filters tree nodes as user types (client-side, no server)
- [ ] Matching behavior:
  - Match against page/folder display names (not file paths)
  - Case-insensitive matching
  - Partial/substring matching (e.g., "start" matches "Getting Started")
- [ ] When filtering:
  - Show only nodes that match the query
  - Auto-expand parent folders of matching items
  - Highlight matching portion of text
  - Hide non-matching nodes entirely
- [ ] Clear button (×) to reset filter and restore full tree
- [ ] Empty state: show "No results" message when nothing matches
- [ ] Keyboard shortcuts:
  - `/` or `Ctrl+K` focuses search input
  - `Escape` clears search and returns focus to content
  - `Enter` navigates to first visible result
  - `↓` moves focus to first result in tree
- [ ] Debounce input (150ms) for performance on large trees
- [ ] Preserve search query across page navigation (optional)
- [ ] Accessible: proper labels, live region for result count

### HTML Structure

```html
<div class="nav-search">
  <div class="search-input-wrapper">
    <svg class="search-icon" aria-hidden="true"><!-- magnifying glass --></svg>
    <input
      type="search"
      class="search-input"
      placeholder="Search..."
      aria-label="Search navigation"
      aria-controls="nav-tree"
      aria-describedby="search-results-status"
    />
    <button class="search-clear" aria-label="Clear search" hidden>
      <svg><!-- × icon --></svg>
    </button>
  </div>
  <div id="search-results-status" class="sr-only" aria-live="polite">
    <!-- "5 results found" or "No results" -->
  </div>
</div>

<nav id="nav-tree" class="tree-nav" aria-label="Site navigation">
  <!-- Tree nodes with data attributes for filtering -->
  <ul role="tree">
    <li role="treeitem" data-search-text="getting started">
      <a href="/getting-started/">
        <span class="match-highlight">Getting</span> Started
      </a>
    </li>
    <!-- ... -->
  </ul>
</nav>

<!-- Empty state -->
<div class="search-empty" hidden>
  <p>No pages found for "<span class="search-query"></span>"</p>
</div>
```

### JavaScript Implementation

```javascript
class NavSearch {
  constructor() {
    this.input = document.querySelector('.search-input');
    this.clearBtn = document.querySelector('.search-clear');
    this.tree = document.querySelector('.tree-nav');
    this.status = document.getElementById('search-results-status');
    this.emptyState = document.querySelector('.search-empty');
    this.allNodes = this.tree.querySelectorAll('[data-search-text]');

    this.init();
  }

  init() {
    this.input.addEventListener('input', this.debounce(() => this.filter(), 150));
    this.clearBtn.addEventListener('click', () => this.clear());
    document.addEventListener('keydown', (e) => this.handleGlobalKeys(e));
  }

  filter() {
    const query = this.input.value.trim().toLowerCase();

    if (!query) {
      this.clear();
      return;
    }

    this.clearBtn.hidden = false;
    let matchCount = 0;

    this.allNodes.forEach(node => {
      const text = node.dataset.searchText;
      const matches = text.includes(query);

      if (matches) {
        matchCount++;
        node.hidden = false;
        this.highlightMatch(node, query);
        this.expandParents(node);
      } else {
        node.hidden = true;
      }
    });

    // Update status for screen readers
    this.status.textContent = matchCount
      ? `${matchCount} result${matchCount === 1 ? '' : 's'} found`
      : 'No results found';

    // Show/hide empty state
    this.emptyState.hidden = matchCount > 0;
    if (!this.emptyState.hidden) {
      this.emptyState.querySelector('.search-query').textContent = this.input.value;
    }
  }

  highlightMatch(node, query) {
    const link = node.querySelector('a, .folder-link');
    const originalText = node.dataset.searchText;
    const regex = new RegExp(`(${this.escapeRegex(query)})`, 'gi');
    link.innerHTML = originalText.replace(regex, '<mark class="match-highlight">$1</mark>');
  }

  expandParents(node) {
    let parent = node.parentElement.closest('[role="treeitem"]');
    while (parent) {
      parent.setAttribute('aria-expanded', 'true');
      parent.hidden = false;
      parent = parent.parentElement.closest('[role="treeitem"]');
    }
  }

  clear() {
    this.input.value = '';
    this.clearBtn.hidden = true;
    this.emptyState.hidden = true;
    this.status.textContent = '';

    // Restore all nodes
    this.allNodes.forEach(node => {
      node.hidden = false;
      const link = node.querySelector('a, .folder-link');
      link.textContent = node.dataset.searchText;
    });

    // Collapse folders to default state
    this.restoreDefaultExpansion();
  }

  handleGlobalKeys(e) {
    // Focus search on / or Ctrl+K
    if ((e.key === '/' || (e.ctrlKey && e.key === 'k')) &&
        !['INPUT', 'TEXTAREA'].includes(e.target.tagName)) {
      e.preventDefault();
      this.input.focus();
    }

    // Clear on Escape when input focused
    if (e.key === 'Escape' && document.activeElement === this.input) {
      this.clear();
      this.input.blur();
    }
  }

  debounce(fn, delay) {
    let timeout;
    return (...args) => {
      clearTimeout(timeout);
      timeout = setTimeout(() => fn.apply(this, args), delay);
    };
  }

  escapeRegex(str) {
    return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  }
}

// Initialize
document.addEventListener('DOMContentLoaded', () => new NavSearch());
```

### File Location

```
volcano/
├── internal/
│   └── templates/
│       ├── search.html      # Search input markup
│       └── search.js        # Filter logic
├── internal/
│   └── styles/
│       └── search.css       # Search styling
```

### Styling Notes

```css
.nav-search {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  position: sticky;
  top: 0;
  background: var(--bg-primary);
  z-index: 10;
}

.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 10px;
  width: 16px;
  height: 16px;
  color: var(--text-secondary);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 8px 32px 8px 36px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 14px;
}

.search-input:focus {
  outline: none;
  border-color: var(--text-secondary);
}

.search-input::placeholder {
  color: var(--text-secondary);
}

.search-clear {
  position: absolute;
  right: 8px;
  padding: 4px;
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-secondary);
}

.search-clear:hover {
  color: var(--text-primary);
}

/* Match highlighting */
.match-highlight {
  background: rgba(255, 220, 0, 0.3);
  border-radius: 2px;
}

[data-theme="dark"] .match-highlight {
  background: rgba(255, 220, 0, 0.2);
}

/* Empty state */
.search-empty {
  padding: 24px 16px;
  text-align: center;
  color: var(--text-secondary);
  font-size: 14px;
}

/* Hidden nodes */
[role="treeitem"][hidden] {
  display: none;
}

/* Screen reader only */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  border: 0;
}
```

### Template Data

```go
// Each tree node needs search text as data attribute
type TreeNode struct {
    // ... existing fields
    SearchText string // Lowercase display name for filtering
}

func (n *TreeNode) GetSearchText() string {
    return strings.ToLower(n.Name)
}
```

### Mobile Considerations

- Search input visible in mobile drawer
- On mobile, tapping a result closes drawer and navigates
- Virtual keyboard doesn't obscure results (scroll into view)
- Touch-friendly clear button (min 44px tap target)

### Performance Notes

- Tree is filtered client-side (no network requests)
- `data-search-text` attributes avoid DOM text extraction on each keystroke
- Debouncing prevents excessive DOM updates
- For very large trees (500+ nodes), consider virtual scrolling

---

## Story 34: Top Navigation Bar (Root Files)

**Goal**: Optionally move root-level markdown files from the tree to a horizontal navigation bar above page content.

### Acceptance Criteria

- [ ] Enabled via `--top-nav` flag (disabled by default)
- [ ] Only applies when root folder has ≤5 markdown files
- [ ] Root files removed from sidebar tree when enabled
- [ ] Root files displayed as horizontal nav links above content
- [ ] Navigation bar appears on all pages
- [ ] Current page highlighted in nav bar
- [ ] Folders remain in sidebar tree (only files move)
- [ ] Root `index.md` becomes site home, not shown in nav bar
- [ ] Mobile: nav bar scrolls horizontally if needed
- [ ] If >5 root files, flag has no effect (files stay in tree)

### Directory Example

```
docs/
├── index.md           ← Home page (not in nav bar)
├── about.md           ← Moves to nav bar
├── contact.md         ← Moves to nav bar
├── faq.md             ← Moves to nav bar
├── guides/
│   ├── index.md
│   └── getting-started.md
└── api/
    └── endpoints.md

With --top-nav:
┌─────────────────────────────────────────────────────┐
│  About   Contact   FAQ                              │  ← Top nav bar
├─────────────────────────────────────────────────────┤
│ ▼ Guides        │                                   │
│     Getting...  │   [Page Content]                  │
│ ▼ Api           │                                   │
│     Endpoints   │                                   │
└─────────────────────────────────────────────────────┘

Without --top-nav (default):
┌─────────────────────────────────────────────────────┐
│   About         │                                   │
│   Contact       │                                   │
│   Faq           │   [Page Content]                  │
│ ▼ Guides        │                                   │
│     Getting...  │                                   │
│ ▼ Api           │                                   │
│     Endpoints   │                                   │
└─────────────────────────────────────────────────────┘
```

### Configuration

```go
type Config struct {
    // ... existing fields
    TopNav bool // --top-nav flag
}
```

### Data Structures

```go
type SiteLayout struct {
    TopNavItems []*TreeNode // Root files for top nav (when enabled)
    TreeRoot    *TreeNode   // Sidebar tree (folders only when top-nav enabled)
}

func BuildSiteLayout(tree *SiteTree, config *Config) SiteLayout {
    if !config.TopNav {
        return SiteLayout{TreeRoot: tree.Root}
    }

    rootFiles := getRootFiles(tree.Root)
    if len(rootFiles) > 5 {
        // Too many files, keep default behavior
        return SiteLayout{TreeRoot: tree.Root}
    }

    // Separate files from folders
    topNav := filterRootFiles(tree.Root)      // Excludes index.md
    treeRoot := filterRootFolders(tree.Root)  // Folders only

    return SiteLayout{
        TopNavItems: topNav,
        TreeRoot:    treeRoot,
    }
}

func getRootFiles(root *TreeNode) []*TreeNode {
    var files []*TreeNode
    for _, child := range root.Children {
        if !child.IsFolder && !isIndexFile(child.Name) {
            files = append(files, child)
        }
    }
    return files
}
```

### HTML Structure

```html
<!-- Top navigation bar (when --top-nav enabled) -->
{{if .TopNavItems}}
<nav class="top-nav" aria-label="Main navigation">
  <ul class="top-nav-list">
    {{range .TopNavItems}}
    <li>
      <a href="{{.URL}}" {{if eq .URL $.CurrentPath}}aria-current="page" class="active"{{end}}>
        {{.Name}}
      </a>
    </li>
    {{end}}
  </ul>
</nav>
{{end}}

<!-- Page content -->
<main class="content">
  <!-- ... -->
</main>
```

### File Location

```
volcano/
├── cmd/
│   └── root.go          # Add --top-nav flag
├── internal/
│   └── tree/
│       └── layout.go    # SiteLayout logic
├── internal/
│   └── templates/
│       └── topnav.html  # Top nav partial
├── internal/
│   └── styles/
│       └── topnav.css
```

### Styling Notes

```css
.top-nav {
  border-bottom: 1px solid var(--border-color);
  padding: 0 24px;
  background: var(--bg-primary);
}

.top-nav-list {
  display: flex;
  gap: 8px;
  list-style: none;
  margin: 0;
  padding: 0;
  overflow-x: auto;
}

.top-nav-list li a {
  display: block;
  padding: 12px 16px;
  color: var(--text-secondary);
  text-decoration: none;
  white-space: nowrap;
  border-bottom: 2px solid transparent;
}

.top-nav-list li a:hover {
  color: var(--text-primary);
}

.top-nav-list li a.active {
  color: var(--text-primary);
  border-bottom-color: var(--text-primary);
}

/* Mobile: horizontal scroll */
@media (max-width: 768px) {
  .top-nav {
    padding: 0 16px;
  }

  .top-nav-list {
    -webkit-overflow-scrolling: touch;
  }
}
```

---

## Story 35: Auto-Generated Folder Index

**Goal**: Automatically generate an index page for folders that lack an `index.md` or matching sibling file.

### Acceptance Criteria

- [ ] Applies to folders without `index.md` and without matching sibling (e.g., `posts/` without `posts.md`)
- [ ] Generates real HTML file at `/folder-name/index.html`
- [ ] Page content: list of child pages as links
- [ ] List shows page titles (from H1 or filename)
- [ ] Links sorted alphabetically (folders first, then files)
- [ ] Folder remains clickable in navigation tree
- [ ] Page uses same template/layout as regular pages
- [ ] Page title derived from folder name (e.g., "Posts")
- [ ] Nested folders also get auto-generated indexes if needed
- [ ] Breadcrumbs work correctly on auto-generated pages

### Directory Example

```
docs/
├── guides/
│   ├── index.md           ← Has index, no auto-generation
│   └── getting-started.md
├── posts/                  ← No index.md, no posts.md
│   ├── hello-world.md
│   ├── second-post.md
│   └── drafts/            ← Nested, also no index
│       └── wip.md
└── api/
    └── endpoints.md       ← No index, but only one file

Generated pages:
- /posts/index.html        ← Auto-generated listing
- /posts/drafts/index.html ← Auto-generated listing
- /api/index.html          ← Auto-generated (even for 1 file)
```

### Generated Page Content

```html
<!-- Auto-generated /posts/index.html -->
<article>
  <h1>Posts</h1>
  <ul class="folder-index">
    <li><a href="/posts/drafts/">Drafts</a></li>
    <li><a href="/posts/hello-world/">Hello World</a></li>
    <li><a href="/posts/second-post/">Second Post</a></li>
  </ul>
</article>
```

### Data Structures

```go
type AutoIndex struct {
    FolderNode *TreeNode
    Title      string      // Derived from folder name
    Children   []IndexItem // Sorted list of children
    OutputPath string      // e.g., "posts/index.html"
    URLPath    string      // e.g., "/posts/"
}

type IndexItem struct {
    Title    string // Page title or folder name
    URL      string
    IsFolder bool
}

func NeedsAutoIndex(node *TreeNode) bool {
    if !node.IsFolder {
        return false
    }
    // Already has index.md
    for _, child := range node.Children {
        if isIndexFile(child.Name) {
            return false
        }
    }
    // Has matching sibling file (handled by Story 32)
    if node.IsClickable && node.ContentType == "sibling" {
        return false
    }
    return true
}

func GenerateAutoIndex(node *TreeNode) AutoIndex {
    var items []IndexItem

    for _, child := range node.VisibleChildren() {
        items = append(items, IndexItem{
            Title:    child.Name,
            URL:      child.URLPath,
            IsFolder: child.IsFolder,
        })
    }

    // Sort: folders first, then alphabetically
    sort.Slice(items, func(i, j int) bool {
        if items[i].IsFolder != items[j].IsFolder {
            return items[i].IsFolder
        }
        return items[i].Title < items[j].Title
    })

    return AutoIndex{
        FolderNode: node,
        Title:      node.Name,
        Children:   items,
        OutputPath: filepath.Join(node.Path, "index.html"),
        URLPath:    "/" + node.Path + "/",
    }
}
```

### Template

```html
{{define "auto-index"}}
<article class="auto-index-page">
  <h1>{{.Title}}</h1>

  {{if .Children}}
  <ul class="folder-index">
    {{range .Children}}
    <li class="{{if .IsFolder}}folder-item{{else}}page-item{{end}}">
      <a href="{{.URL}}">
        {{if .IsFolder}}<span class="folder-icon">📁</span>{{end}}
        {{.Title}}
      </a>
    </li>
    {{end}}
  </ul>
  {{else}}
  <p class="empty-folder">This folder is empty.</p>
  {{end}}
</article>
{{end}}
```

### File Location

```
volcano/
├── internal/
│   └── generator/
│       ├── autoindex.go      # Auto-index detection and generation
│       └── autoindex_test.go
├── internal/
│   └── templates/
│       └── autoindex.html    # Listing template
├── internal/
│   └── styles/
│       └── autoindex.css
```

### Styling Notes

```css
.folder-index {
  list-style: none;
  padding: 0;
  margin: 24px 0;
}

.folder-index li {
  margin: 0;
  padding: 0;
}

.folder-index a {
  display: block;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
  text-decoration: none;
}

.folder-index a:hover {
  background: var(--bg-secondary);
}

.folder-index li:last-child a {
  border-bottom: none;
}

.folder-item a {
  font-weight: 500;
}

.folder-icon {
  margin-right: 8px;
}

.empty-folder {
  color: var(--text-secondary);
  font-style: italic;
}
```

### Integration with Story 32

- Story 32 makes folders clickable when they have content
- Story 35 ensures ALL folders become clickable by generating content
- After Story 35, `NeedsAutoIndex` folders get:
  1. `IsClickable = true`
  2. `ContentURL = auto-generated index URL`
  3. `ContentType = "auto"`

```go
func (n *TreeNode) ResolveClickableFolder(siblings []*TreeNode) {
    // ... existing logic for index.md and sibling file ...

    // Fallback: auto-generated index (Story 35)
    if NeedsAutoIndex(n) {
        n.IsClickable = true
        n.ContentURL = "/" + n.Path + "/"
        n.ContentType = "auto"
    }
}
```

---

## Story 36: H1-Based Tree Labels

**Goal**: Use the first H1 heading from markdown files as the display name in the navigation tree, falling back to the filename-derived label if no H1 exists.

### Acceptance Criteria

- [ ] Extract first H1 (`# Heading`) from each markdown file during scanning
- [ ] Use H1 text as the display name in navigation tree
- [ ] Fallback to filename-derived label if no H1 found:
  - `getting-started.md` → "Getting Started"
- [ ] H1 extraction happens at build time, not runtime
- [ ] Handle edge cases:
  - Multiple H1s: use only the first one
  - H1 with inline formatting: strip markdown (`**bold**` → "bold")
  - H1 with links: extract text only
  - Empty H1: fallback to filename
- [ ] Also use H1 for:
  - Page `<title>` tag
  - Breadcrumb labels
  - Previous/Next navigation labels
  - Search text in nav search (Story 33)

### Examples

```
# File: docs/getting-started.md
# Quick Start Guide        ← H1 found

Tree shows: "Quick Start Guide" (not "Getting Started")

# File: docs/api-reference.md
(no H1 in file)

Tree shows: "Api Reference" (fallback to filename)

# File: docs/faq.md
# Frequently Asked Questions

Tree shows: "Frequently Asked Questions" (not "Faq")
```

### Data Structures

```go
type TreeNode struct {
    Name        string // Display name (from H1 or filename)
    FileName    string // Original filename for reference
    H1Title     string // Extracted H1 (empty if none)
    // ... other fields
}

// Extract H1 from markdown content
func ExtractH1(content []byte) string {
    // Match first line starting with "# "
    // Strip inline formatting
    // Return empty string if not found
}

// Build display name with H1 priority
func (n *TreeNode) DisplayName() string {
    if n.H1Title != "" {
        return n.H1Title
    }
    return n.Name // Fallback to filename-derived name
}
```

### Implementation

```go
import (
    "bufio"
    "regexp"
    "strings"
)

var h1Regex = regexp.MustCompile(`^#\s+(.+)$`)
var inlineMarkdown = regexp.MustCompile(`[\*_~\[\]` + "`" + `]`)

func ExtractH1(content []byte) string {
    scanner := bufio.NewScanner(bytes.NewReader(content))

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())

        // Skip empty lines and frontmatter
        if line == "" || line == "---" {
            continue
        }

        // Check for H1
        if matches := h1Regex.FindStringSubmatch(line); len(matches) > 1 {
            title := matches[1]
            // Strip inline markdown formatting
            title = stripInlineMarkdown(title)
            title = strings.TrimSpace(title)
            if title != "" {
                return title
            }
        }

        // If first non-empty line isn't H1, stop looking
        // (H1 should be at the top)
        break
    }

    return ""
}

func stripInlineMarkdown(text string) string {
    // Remove **bold**, *italic*, `code`, [links](url), etc.
    // Simple approach: remove common markdown characters
    text = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`).ReplaceAllString(text, "$1")
    text = regexp.MustCompile(`[*_~` + "`" + `]`).ReplaceAllString(text, "")
    return text
}
```

### Integration Points

```go
// During tree building (scanner.go)
func buildTreeNode(path string, info os.FileInfo) *TreeNode {
    node := &TreeNode{
        FileName: info.Name(),
        Name:     cleanFilename(info.Name()), // Default from filename
    }

    if !info.IsDir() && isMarkdownFile(info.Name()) {
        content, err := os.ReadFile(path)
        if err == nil {
            if h1 := ExtractH1(content); h1 != "" {
                node.H1Title = h1
                node.Name = h1 // Override display name
            }
        }
    }

    return node
}
```

### File Location

```
volcano/
├── internal/
│   └── tree/
│       ├── scanner.go      # Update to extract H1 during scan
│       ├── h1.go           # H1 extraction logic
│       └── h1_test.go      # Test various H1 formats
```

### Test Cases

```go
func TestExtractH1(t *testing.T) {
    tests := []struct {
        name     string
        content  string
        expected string
    }{
        {
            name:     "simple h1",
            content:  "# Hello World\n\nContent here",
            expected: "Hello World",
        },
        {
            name:     "h1 with bold",
            content:  "# Getting **Started**\n\nContent",
            expected: "Getting Started",
        },
        {
            name:     "h1 with link",
            content:  "# [Home](/) Page\n\nContent",
            expected: "Home Page",
        },
        {
            name:     "no h1",
            content:  "## Second Level\n\nNo h1 here",
            expected: "",
        },
        {
            name:     "h1 after content",
            content:  "Some intro text\n\n# Late H1",
            expected: "", // H1 must be first
        },
        {
            name:     "h1 with code",
            content:  "# Using `volcano` CLI\n\nContent",
            expected: "Using volcano CLI",
        },
        {
            name:     "empty h1",
            content:  "# \n\nContent",
            expected: "",
        },
        {
            name:     "frontmatter then h1",
            content:  "---\ntitle: ignored\n---\n# Real Title",
            expected: "Real Title",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ExtractH1([]byte(tt.content))
            if result != tt.expected {
                t.Errorf("got %q, want %q", result, tt.expected)
            }
        })
    }
}
```

### Performance Note

- H1 extraction reads only the first few lines of each file
- Uses early termination (stops after first non-empty, non-H1 line)
- File content is read once during tree building, then cached

---

## Story 37: Filename Date & Number Prefixes (Medusa-Compatible)

**Goal**: Support date and number prefixes in filenames for sorting and organization, matching medusa-ssg's behavior.

### Acceptance Criteria

- [ ] Parse date prefix from filenames: `YYYY-MM-DD-` format
- [ ] Parse number prefix from filenames: leading digits like `01-`
- [ ] Support combined date + number: `2024-01-15-01-section.md`
- [ ] Strip prefixes from display names and URLs
- [ ] Sort using three-tier priority system
- [ ] Support draft files with `_` prefix
- [ ] Fallback to file modification time if no date prefix

### Date Prefix Format

```
YYYY-MM-DD-rest-of-filename.md

Examples:
2024-01-15-hello-world.md    → date: 2024-01-15, slug: "hello-world"
2024-12-31-new-years.md      → date: 2024-12-31, slug: "new-years"
no-date-here.md              → date: file mtime, slug: "no-date-here"
```

### Number Prefix Format

```
# Simple number prefix (no date)
01-introduction.md           → number: 1, slug: "introduction"
02-getting-started.md        → number: 2, slug: "getting-started"
10-conclusion.md             → number: 10, slug: "conclusion"

# Date + number combined
2024-01-15-01-part-one.md    → date: 2024-01-15, number: 1, slug: "part-one"
2024-01-15-02-part-two.md    → date: 2024-01-15, number: 2, slug: "part-two"

# Date without number
2024-01-15-my-post.md        → date: 2024-01-15, number: nil, slug: "my-post"
```

### Sorting Rules (Three-Tier Priority)

**Default sort order (newest/highest first):**

1. **Primary: Date** - Newer dates first
2. **Secondary: Number** - Higher numbers first (when dates equal)
3. **Tertiary: Filename** - Alphabetical (when dates and numbers equal)

```
# Example sort order (newest first):
2024-03-01-post.md           ← newest date
2024-02-15-02-part-two.md    ← same date, higher number
2024-02-15-01-part-one.md    ← same date, lower number
2024-01-01-zebra.md          ← older date, 'z' comes after 'a'
2024-01-01-alpha.md          ← older date, 'a' comes first alphabetically
01-intro.md                  ← no date (uses mtime), has number
about.md                     ← no date, no number
```

### Draft Files

Files and folders prefixed with `_` are drafts:

```
posts/
├── 2024-01-15-published.md     ← included in build
├── _2024-02-01-work-in-progress.md  ← DRAFT, excluded
└── _drafts/                    ← DRAFT folder, all contents excluded
    └── idea.md

# Draft behavior:
- Excluded from navigation tree
- Excluded from generated output
- Future: --drafts flag to include them
```

### URL/Slug Generation

Date and number prefixes stripped from URLs:

```
Source                          → URL
───────────────────────────────────────────────────
2024-01-15-hello-world.md       → /hello-world/
01-introduction.md              → /introduction/
2024-01-15-01-part-one.md       → /part-one/
posts/2024-02-20-news.md        → /posts/news/
docs/01-getting-started.md      → /docs/getting-started/
```

### Display Name Generation

For nav tree (before H1 override from Story 36):

```
Filename                        → Display Name
───────────────────────────────────────────────────
2024-01-15-hello-world.md       → "Hello World"
01-introduction.md              → "Introduction"
2024-01-15-01-part-one.md       → "Part One"
getting-started.md              → "Getting Started"
```

### Data Structures

```go
type FileMetadata struct {
    OriginalName string     // "2024-01-15-01-hello-world.md"
    Date         time.Time  // Parsed or file mtime
    DateSource   string     // "filename" or "mtime"
    Number       *int       // nil if no number prefix
    Slug         string     // "hello-world"
    DisplayName  string     // "Hello World"
    IsDraft      bool       // starts with "_"
}

type TreeNode struct {
    // ... existing fields
    Metadata FileMetadata
}
```

### Implementation

```go
import (
    "regexp"
    "strconv"
    "strings"
    "time"
)

var datePrefix = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})-(.+)$`)
var numberPrefix = regexp.MustCompile(`^(\d+)-(.+)$`)

// ExtractFileMetadata parses date/number prefixes from filename
func ExtractFileMetadata(filename string, modTime time.Time) FileMetadata {
    stem := strings.TrimSuffix(filename, filepath.Ext(filename))
    meta := FileMetadata{
        OriginalName: filename,
        IsDraft:      strings.HasPrefix(stem, "_"),
    }

    // Remove draft prefix for further processing
    if meta.IsDraft {
        stem = strings.TrimPrefix(stem, "_")
    }

    remaining := stem

    // Try to extract date prefix (YYYY-MM-DD-)
    if matches := datePrefix.FindStringSubmatch(stem); len(matches) == 5 {
        year, _ := strconv.Atoi(matches[1])
        month, _ := strconv.Atoi(matches[2])
        day, _ := strconv.Atoi(matches[3])

        if date, err := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC); err == nil {
            meta.Date = date
            meta.DateSource = "filename"
            remaining = matches[4]
        }
    }

    // If no date from filename, use file mtime
    if meta.DateSource == "" {
        meta.Date = modTime
        meta.DateSource = "mtime"
    }

    // Try to extract number prefix from remaining
    if matches := numberPrefix.FindStringSubmatch(remaining); len(matches) == 3 {
        num, _ := strconv.Atoi(matches[1])
        meta.Number = &num
        remaining = matches[2]
    }

    meta.Slug = remaining
    meta.DisplayName = titleize(remaining)

    return meta
}

// titleize converts "hello-world" to "Hello World"
func titleize(slug string) string {
    parts := strings.FieldsFunc(slug, func(r rune) bool {
        return r == '-' || r == '_'
    })
    for i, part := range parts {
        parts[i] = strings.Title(part)
    }
    return strings.Join(parts, " ")
}
```

### Sorting Implementation

```go
// SortNodes sorts tree nodes by date, number, then name
func SortNodes(nodes []*TreeNode, newestFirst bool) {
    sort.Slice(nodes, func(i, j int) bool {
        a, b := nodes[i], nodes[j]

        // Folders always before files
        if a.IsFolder != b.IsFolder {
            return a.IsFolder
        }

        // Primary: Date
        if !a.Metadata.Date.Equal(b.Metadata.Date) {
            if newestFirst {
                return a.Metadata.Date.After(b.Metadata.Date)
            }
            return a.Metadata.Date.Before(b.Metadata.Date)
        }

        // Secondary: Number (nil treated as infinity when newestFirst)
        aNum := getNumberForSort(a.Metadata.Number, newestFirst)
        bNum := getNumberForSort(b.Metadata.Number, newestFirst)
        if aNum != bNum {
            if newestFirst {
                return aNum > bNum
            }
            return aNum < bNum
        }

        // Tertiary: Slug (alphabetical)
        if newestFirst {
            return a.Metadata.Slug > b.Metadata.Slug
        }
        return a.Metadata.Slug < b.Metadata.Slug
    })
}

func getNumberForSort(n *int, newestFirst bool) int {
    if n == nil {
        if newestFirst {
            return -1 // Sort after numbered items
        }
        return 999999 // Sort after numbered items
    }
    return *n
}
```

### File Location

```
volcano/
├── internal/
│   └── tree/
│       ├── metadata.go      # FileMetadata extraction
│       ├── metadata_test.go
│       ├── sort.go          # Sorting logic
│       └── sort_test.go
```

### Test Cases

```go
func TestExtractFileMetadata(t *testing.T) {
    mtime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

    tests := []struct {
        filename    string
        wantDate    string
        wantNumber  *int
        wantSlug    string
        wantDisplay string
        wantDraft   bool
    }{
        // Date prefix
        {"2024-01-15-hello-world.md", "2024-01-15", nil, "hello-world", "Hello World", false},

        // Number prefix
        {"01-introduction.md", "2020-01-01", intPtr(1), "introduction", "Introduction", false},

        // Date + number
        {"2024-01-15-01-part-one.md", "2024-01-15", intPtr(1), "part-one", "Part One", false},

        // No prefix
        {"about.md", "2020-01-01", nil, "about", "About", false},

        // Draft
        {"_2024-01-15-draft.md", "2024-01-15", nil, "draft", "Draft", true},

        // Draft folder marker
        {"_work-in-progress.md", "2020-01-01", nil, "work-in-progress", "Work In Progress", true},
    }

    for _, tt := range tests {
        t.Run(tt.filename, func(t *testing.T) {
            meta := ExtractFileMetadata(tt.filename, mtime)

            wantDate, _ := time.Parse("2006-01-02", tt.wantDate)
            if !meta.Date.Equal(wantDate) {
                t.Errorf("Date = %v, want %v", meta.Date, wantDate)
            }
            // ... more assertions
        })
    }
}

func TestSortNodes(t *testing.T) {
    // Test sorting priority: date > number > name
    nodes := []*TreeNode{
        {Metadata: FileMetadata{Date: date("2024-01-01"), Slug: "alpha"}},
        {Metadata: FileMetadata{Date: date("2024-03-01"), Slug: "post"}},
        {Metadata: FileMetadata{Date: date("2024-02-15"), Number: intPtr(1), Slug: "part-one"}},
        {Metadata: FileMetadata{Date: date("2024-02-15"), Number: intPtr(2), Slug: "part-two"}},
    }

    SortNodes(nodes, true) // newest first

    // Expected order: 2024-03-01, 2024-02-15-02, 2024-02-15-01, 2024-01-01
    expected := []string{"post", "part-two", "part-one", "alpha"}
    for i, node := range nodes {
        if node.Metadata.Slug != expected[i] {
            t.Errorf("Position %d: got %s, want %s", i, node.Metadata.Slug, expected[i])
        }
    }
}
```

### Integration with Existing Stories

- **Story 3 (Tree Building)**: Use `ExtractFileMetadata` during scan
- **Story 32 (Clickable Folders)**: Slugs used for folder URLs
- **Story 35 (Auto-Index)**: Sort children using new sorting rules
- **Story 36 (H1 Labels)**: H1 overrides `DisplayName` from metadata

### Configuration Options (Future)

```go
type Config struct {
    // ... existing
    SortNewestFirst bool   // --sort-newest-first (default: true)
    IncludeDrafts   bool   // --drafts (include _ prefixed files)
}
```

---

## Implementation Priority

### Phase 2A: Essential Navigation Enhancements
- Story 15: Breadcrumb Navigation
- Story 16: Previous/Next Page Navigation
- Story 17: Heading Anchor Links
- Story 31: Smooth Scroll Behavior
- Story 32: Clickable Folder Navigation
- Story 33: Navigation Tree Search
- Story 34: Top Navigation Bar (Root Files)
- Story 35: Auto-Generated Folder Index
- Story 36: H1-Based Tree Labels
- Story 37: Filename Date & Number Prefixes

### Phase 2B: Content Improvements
- Story 14: Table of Contents Component
- Story 19: Code Block Copy Button
- Story 29: Admonition/Callout Blocks
- Story 30: Code Line Highlighting

### Phase 2C: User Experience Polish
- Story 18: External Link Indicators
- Story 21: Print Stylesheet
- Story 22: Reading Time Indicator
- Story 25: Back to Top Button

### Phase 2D: SEO & Sharing
- Story 26: SEO Meta Tags Generation
- Story 27: Open Graph Support
- Story 28: Custom Favicon Support

### Phase 2E: Advanced Features
- Story 20: Keyboard Navigation Shortcuts
- Story 23: Last Modified Display
- Story 24: Scroll Progress Indicator

---

## Dependencies

Additional dependencies for new features:

```go
// go.mod additions
require (
    golang.org/x/net v0.20.0  // HTML parsing for external links
)
```

Most features use only the standard library and existing dependencies (Goldmark, Chroma).

---

## Testing Strategy

Each story should include:

1. **Unit tests** for core logic (parsing, generation, transformations)
2. **Integration tests** for template rendering with new components
3. **E2E tests** for visual verification (optional, via screenshots)
4. **Accessibility tests** for ARIA compliance

Example test structure:

```
volcano/
├── internal/
│   └── toc/
│       ├── extractor_test.go
│       └── testdata/
│           ├── simple.html
│           └── nested.html
```

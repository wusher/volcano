# Navigation

Volcano provides multiple navigation features to help users find content.

## Sidebar Tree

The left sidebar shows your folder structure as a collapsible tree. Always visible on desktop, drawer on mobile.

**Features:**
- Expand/collapse folders
- Active page highlighting
- Auto-expands to current page
- Keyboard accessible

Your folder structure becomes the navigation automatically.

## Navigation Features

Each feature can be enabled independently:

### [[search|Search]]
Command palette (Cmd+K) searches pages and headings.
```bash
volcano ./docs --search
```

### [[breadcrumbs|Breadcrumbs]]
Shows current page location in hierarchy. Enabled by default.
```bash
volcano ./docs --breadcrumbs=false  # to disable
```

### [[table-of-contents|Table of Contents]]
Auto-generated sidebar for pages with 3+ headings. Always enabled.

### [[page-navigation|Page Navigation]]
Previous/next links at bottom of pages.
```bash
volcano ./docs --page-nav
```

### [[top-navigation|Top Navigation]]
Horizontal bar showing root-level pages.
```bash
volcano ./docs --top-nav
```

### [[instant-navigation|Instant Navigation]]
Hover prefetching for near-instant page loads.
```bash
volcano ./docs --instant-nav
```

## UI Elements

### Scroll Progress Bar
Thin bar at top showing scroll position through the page.

### Back to Top Button
Floating button (bottom-right) appears after scrolling 300px.

### Theme Toggle
Switch between light and dark mode. Located in sidebar or top nav.

## Keyboard Shortcuts

Press `?` to see all shortcuts.

| Key | Action |
|-----|--------|
| `Cmd+K` / `Ctrl+K` | Open search (requires [[search\|--search]]) |
| `t` | Toggle theme |
| `z` | Toggle zen mode (hides sidebar) |
| `n` | Next page (requires [[page-navigation\|--page-nav]]) |
| `p` | Previous page (requires [[page-navigation\|--page-nav]]) |
| `h` | Go to homepage |
| `?` | Show shortcuts |
| `Esc` | Close modal |

## Combining Features

Mix and match navigation features:

```bash
# Full navigation suite
volcano ./docs \
  --search \
  --instant-nav \
  --page-nav \
  --top-nav

# Minimal setup
volcano ./docs \
  --breadcrumbs=false

# Documentation focus
volcano ./docs \
  --search \
  --instant-nav
```

## Accessibility

All navigation features include:
- Proper ARIA roles and labels
- Keyboard navigation support
- Screen reader compatibility
- Focus indicators
- Skip links

## Related

- [[guides/organizing-content]] — Structure your content
- [[keyboard-shortcuts]] — Complete shortcut reference

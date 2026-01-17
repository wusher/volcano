# Keyboard Shortcuts

Navigate your site quickly with keyboard shortcuts.

## Default Shortcuts

These shortcuts are always available:

| Key | Action |
|-----|--------|
| `t` | Toggle theme (light/dark) |
| `z` | Toggle zen mode (hide sidebar) |
| `h` | Go to homepage |
| `?` | Show shortcuts modal |
| `Esc` | Close modals, exit zen mode |

## Feature-Specific Shortcuts

Additional shortcuts when features are enabled:

### Search

**Requires:** `--search` flag

| Key | Action |
|-----|--------|
| `Cmd+K` (Mac)<br>`Ctrl+K` (Win/Linux) | Open search |
| `↑` / `↓` | Navigate results |
| `Enter` | Go to selected result |
| `Esc` | Close search |

### Page Navigation

**Requires:** `--page-nav` flag

| Key | Action |
|-----|--------|
| `n` | Next page |
| `p` | Previous page |

## Zen Mode

Press `z` to hide the sidebar for distraction-free reading:
- Sidebar hidden
- Content takes full width
- Press `z` or `Esc` to exit

Perfect for focused reading of long pages.

## Shortcuts Modal

Press `?` to see the full shortcuts reference:
- Lists all available shortcuts
- Adapts based on enabled features
- Press `Esc` or click outside to close

## Accessibility

Shortcuts are designed to:
- Work with screen readers
- Not conflict with browser shortcuts
- Be easy to remember (single keys)
- Follow common conventions (Cmd+K for search)

## Disabled in Input Fields

Shortcuts are automatically disabled when typing in:
- Search input
- Any text field
- Text areas

This prevents accidental triggers while editing.

## Browser Conflicts

Some browsers use overlapping shortcuts:

**Cmd+K / Ctrl+K:**
- Chrome: Address bar search
- Firefox: Address bar search
- **Solution:** Volcano's search captures it first when focused on page

**Cmd+H / Ctrl+H:**
- Most browsers: History
- **Solution:** Volcano uses `h` alone (no modifier)

## Customization

Keyboard shortcuts are built-in and cannot be customized. They're chosen to:
- Avoid browser conflicts
- Be memorable
- Follow common conventions

## Related

- [[search]] — Search feature and shortcuts
- [[page-navigation]] — Previous/next navigation
- [[navigation]] — Navigation overview

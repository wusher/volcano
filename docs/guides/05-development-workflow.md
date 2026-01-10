# Development Workflow

Use Volcano's server modes to preview and develop your site efficiently.

## Two Server Modes

Volcano offers two ways to preview your site:

| Mode | Command | Best For |
|------|---------|----------|
| Static | `volcano -s ./output` | Testing final output |
| Dynamic | `volcano -s ./docs` | Active development |

## Static Server Mode

Serve pre-generated HTML files:

```bash
# First, generate your site
volcano ./docs -o ./output --title="My Site"

# Then serve the output
volcano -s ./output
```

Open [http://localhost:1776](http://localhost:1776)

**Characteristics:**
- Serves files exactly as they'll appear in production
- Fast — no rendering on each request
- Requires regeneration to see changes

**Use when:**
- Testing the final build before deployment
- Checking that all links work
- Verifying the production output

## Dynamic Server Mode

Render pages on-the-fly from markdown source:

```bash
volcano -s ./docs
```

Volcano detects that `./docs` contains markdown files (not generated HTML) and automatically uses dynamic rendering.

**Characteristics:**
- Renders pages fresh on each request
- See changes by refreshing your browser
- No need to regenerate after edits

**Use when:**
- Actively writing content
- Iterating on site structure
- Testing new pages quickly

## Custom Port

Change the server port with `-p`:

```bash
# Static mode on port 8080
volcano -s -p 8080 ./output

# Dynamic mode on port 3000
volcano -s -p 3000 ./docs
```

## Development Cycle

### For Content Changes

Using dynamic mode (fastest iteration):

1. Start the server: `volcano -s ./docs`
2. Edit your markdown files
3. Refresh your browser to see changes
4. Repeat

### For Theme/CSS Changes

If using custom CSS:

1. Start the server: `volcano -s ./output`
2. Edit your CSS file
3. Regenerate: `volcano ./docs --css ./custom.css`
4. Refresh your browser
5. Repeat

:::tip
Keep two terminal windows open — one running the server, one for regenerating.
:::

## Debugging

### Verbose Mode

See detailed information about what Volcano is doing:

```bash
volcano ./docs --verbose
```

Output includes:
- File processing details
- Navigation tree building
- Link resolution

### Quiet Mode

Suppress all output except errors:

```bash
volcano ./docs -q
```

Useful for scripts and CI pipelines.

## Watching for Changes

Volcano doesn't have built-in file watching, but you can use external tools:

### Using entr (Linux/macOS)

```bash
# Install entr
# macOS: brew install entr
# Linux: apt install entr

# Watch and regenerate
find docs -name '*.md' | entr volcano ./docs -o ./output
```

### Using watchexec

```bash
# Install watchexec
# cargo install watchexec-cli

# Watch and regenerate
watchexec -e md -w docs -- volcano ./docs -o ./output
```

### Using fswatch (macOS)

```bash
fswatch -o docs | xargs -n1 -I{} volcano ./docs -o ./output
```

## Link Validation

Volcano automatically validates all internal links during generation. If any broken links are found, the build fails with detailed error messages:

```
✗ Found 2 broken internal links:
  Page /guides/intro/: broken link /setup/ (not found)
  Page /reference/api/: broken link /deprecated/ (not found)
```

This catches broken wiki links, markdown links, and navigation references before deployment.

In dynamic serve mode (`volcano -s ./docs`), broken links are displayed inline on the page with helpful error messages.

## Manual Testing

After generating, you can also manually verify your site:

```bash
# Generate the site
volcano ./docs -o ./output

# Start the server
volcano -s ./output
```

Navigate through your site checking:
- All sidebar links work
- Wiki links resolve correctly
- Breadcrumbs link to the right pages
- Table of contents links scroll to headings

## Performance Tips

### Large Sites

For sites with many pages:

1. Use static mode for testing (faster serving)
2. Only regenerate after significant changes
3. Use `--quiet` to speed up generation slightly

### Quick Previews

For the fastest preview cycle:

1. Start dynamic server: `volcano -s ./docs`
2. Make your changes
3. Refresh — no regeneration needed

## Common Issues

### Port Already in Use

If you see "address already in use":

```bash
# Use a different port
volcano -s -p 8081 ./output

# Or find and stop the process using the port
lsof -i :1776
kill <PID>
```

### Changes Not Appearing

In static mode:
- Make sure you regenerated after changes
- Check you're editing the source files, not output

In dynamic mode:
- Try a hard refresh (Ctrl+Shift+R / Cmd+Shift+R)
- Clear your browser cache

### Wrong Site Title

The title is set at generation time:

```bash
# This sets the title
volcano ./docs --title="Correct Title"

# Serving doesn't change the title
volcano -s ./output
```

## Next Steps

- [[deploying-your-site]] — Publish your finished site
- [[reference/cli]] — See all server options

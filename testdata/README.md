# Test Data

This folder contains markdown fixtures for integration tests.

## Folders

### `number-prefixes/`

Tests wiki link resolution for files with various number prefix formats:

- `6 - Name.md` - Space-dash-space prefix (Obsidian Kanban style)
- `01-name.md` - Dash prefix
- `02 name.md` - Space prefix with leading zeros (ordering)
- `2023 Goals.md` - Year in filename (should be preserved)
- `0. Inbox.md` - Dot prefix (folder style)

Used by: `TestIntegrationMarkdown_WikiLinksWithNumberPrefixes`

### `wikilinks-md-anchor/`

Tests wiki links that have both `.md` extension AND an anchor:

- `[[file.md#section]]` - Extension should be stripped, anchor preserved
- `[[folder/file.md#section]]` - Nested paths with extension and anchor

Used by: `TestIntegrationMarkdown_WikiLinksWithMdAnchor`

### `wikilinks-attachments/`

Tests wiki links to attachments (images, PDFs, videos, etc.):

- `[[image.png]]` - Extension should be preserved (not slugified)
- `[[attachments/photo.jpg]]` - Nested attachment paths
- `[[My Image.png]]` - Spaces in filename become dashes, but extension preserved
- `![[image.png]]` - Embed links to attachments

Used by: `TestIntegrationMarkdown_WikiLinksWithAttachments`

## Adding New Test Fixtures

1. Create a new folder under `testdata/` for your test scenario
2. Add markdown files that reproduce the issue
3. Reference the folder in your integration test using a relative path:
   ```go
   inputDir := "testdata/your-scenario"
   ```

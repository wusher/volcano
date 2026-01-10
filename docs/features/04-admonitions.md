# Admonitions

Create callout boxes to highlight important information.

## Syntax

Admonitions use a fenced syntax with `:::type`:

```markdown
:::note
This is a note admonition.
:::
```

:::note
This is a note admonition.
:::

## Available Types

### Note

General information or clarification:

```markdown
:::note
Notes provide additional context or background information.
:::
```

:::note
Notes provide additional context or background information.
:::

### Tip

Helpful suggestions or best practices:

```markdown
:::tip
Tips offer useful suggestions to improve your workflow.
:::
```

:::tip
Tips offer useful suggestions to improve your workflow.
:::

### Warning

Important cautions to be aware of:

```markdown
:::warning
Warnings highlight potential issues or things to watch out for.
:::
```

:::warning
Warnings highlight potential issues or things to watch out for.
:::

### Danger

Critical warnings that could cause problems:

```markdown
:::danger
Danger callouts warn about actions that could cause data loss or security issues.
:::
```

:::danger
Danger callouts warn about actions that could cause data loss or security issues.
:::

### Info

Informational callouts:

```markdown
:::info
Info boxes provide additional information that may be useful.
:::
```

:::info
Info boxes provide additional information that may be useful.
:::

## Custom Titles

Add a custom title after the type:

```markdown
:::tip Pro Tip
Use custom titles to make your admonitions more specific.
:::
```

:::tip Pro Tip
Use custom titles to make your admonitions more specific.
:::

```markdown
:::warning Breaking Change
This version introduces breaking changes to the API.
:::
```

:::warning Breaking Change
This version introduces breaking changes to the API.
:::

## Content Inside Admonitions

Admonitions support full markdown content:

```markdown
:::note About Configuration
You can configure multiple options:

- **Option A**: Does something
- **Option B**: Does something else

See the [customization guide](/guides/customizing-appearance/) for details.
:::
```

:::note About Configuration
You can configure multiple options:

- **Option A**: Does something
- **Option B**: Does something else

See the [customization guide](/guides/customizing-appearance/) for details.
:::

### Code in Admonitions

````markdown
:::tip
Use environment variables for sensitive data:

```bash
export API_KEY="your-key-here"
```
:::
````

:::tip
Use environment variables for sensitive data:

```bash
export API_KEY="your-key-here"
```
:::

## Styling

Each admonition type has distinct styling:

| Type | Default Title | Icon | Color |
|------|---------------|------|-------|
| `note` | Note | Info circle | Blue |
| `tip` | Tip | Lightbulb | Green |
| `warning` | Warning | Triangle | Yellow/Orange |
| `danger` | Danger | Octagon | Red |
| `info` | Info | Info circle | Blue |

## CSS Classes

Admonitions use these classes for custom styling:

```css
.admonition { }                    /* Base styles */
.admonition-note { }               /* Note type */
.admonition-tip { }                /* Tip type */
.admonition-warning { }            /* Warning type */
.admonition-danger { }             /* Danger type */
.admonition-info { }               /* Info type */
.admonition-heading { }            /* Title container */
.admonition-icon { }               /* SVG icon */
.admonition-title { }              /* Title text */
.admonition-content { }            /* Content area */
```

## Best Practices

### Use Sparingly

Too many admonitions reduce their impact:

```markdown
<!-- Good: Strategic use -->
Regular content explaining the concept.

:::warning
This action cannot be undone.
:::

More content here.

<!-- Bad: Overuse -->
:::note
Here's some info.
:::

:::tip
Here's a tip.
:::

:::warning
Here's a warning.
:::
```

### Choose the Right Type

| Situation | Type |
|-----------|------|
| Additional context | `note` |
| Helpful suggestion | `tip` |
| Potential issue | `warning` |
| Critical danger | `danger` |
| Supplementary info | `info` |

### Keep Titles Concise

```markdown
<!-- Good -->
:::warning Security Notice
...
:::

<!-- Too Long -->
:::warning Please Be Aware Of This Important Security Consideration
...
:::
```

### Write Actionable Content

```markdown
<!-- Good: Actionable -->
:::warning
Always backup your data before upgrading.
Run `backup.sh` before proceeding.
:::

<!-- Less Useful: Vague -->
:::warning
Be careful with upgrades.
:::
```

## Examples

### Installation Warning

```markdown
:::warning Prerequisites
Before installing, ensure you have:
- Node.js 18 or later
- At least 1GB of free disk space
:::
```

### Pro Tip

```markdown
:::tip Keyboard Shortcut
Press `Ctrl+S` (or `Cmd+S` on Mac) to save your work quickly.
:::
```

### Breaking Change Notice

```markdown
:::danger Breaking Change in v2.0
The `oldFunction()` has been removed.
Use `newFunction()` instead.

Migration guide: [[guides/migration]]
:::
```

### Cross-Platform Note

```markdown
:::note Platform Differences
- **Windows**: Use `\` for path separators
- **macOS/Linux**: Use `/` for path separators
:::
```

## Related

- [[markdown-syntax]] — All markdown features
- [[theming]] — Customize admonition styling

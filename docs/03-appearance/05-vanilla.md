# vanilla Theme

A structural skeleton. All layout, no decoration. Extensively commented CSS. The starting point if you want a fully custom look.

```bash
volcano ./docs --theme vanilla --url="https://example.com"
```

## Light

![vanilla theme, light mode](/images/themes/vanilla-light.png)

## Dark

![vanilla theme, dark mode](/images/themes/vanilla-dark.png)

## Features

- All layout structure preserved (sidebar, content, TOC, etc.)
- No colors, fonts, or decorations beyond browser defaults
- Every class is documented with inline comments
- Drop-in scaffold — replace with your own design

## Best For

Starting a custom theme. The vanilla CSS is what `volcano css -o my-theme.css` exports — you get the structure for free, then style it however you want.

## Next Step

Export the CSS skeleton and start customizing:

```bash
volcano css -o my-theme.css
volcano ./docs --css ./my-theme.css --url="https://example.com"
```

See **[[custom-css|Custom CSS]]** for a walkthrough of CSS variables, key classes, and worked examples.

# Generate a Custom Theme

Create a custom CSS theme for your Volcano site through an interactive process.

## About

This prompt guides an LLM through exporting Volcano's CSS skeleton, asking about your design preferences, and generating a complete custom theme with light and dark modes.

## Prompt

```
I'm building a static site with Volcano and need help creating a custom CSS theme.

## Step 1: Ask Me Questions

Ask me 5-7 questions about my design preferences one at a time:
- What type of site? (blog, docs, portfolio, etc.)
- Overall theme/styling? (professional, playful, minimal, bold, etc.)
- Brand colors or should you suggest?
- Light mode, dark mode, or both?
- Typography style? (modern, classic, technical, etc.)
- Any websites whose design I like?

## Step 2: Export the CSS Skeleton

After I answer, run this command to export Volcano's CSS skeleton:

go install github.com/wusher/volcano@latest
volcano css -o skeleton.css

Read skeleton.css to understand:
- What CSS variables are defined
- What components exist (sidebar, content, code blocks, etc.)
- What selectors are used

## Step 3: Generate the Theme

Based on my preferences AND the skeleton structure, generate a complete theme that:
1. Preserves the skeleton structure (same selectors, components)
2. Applies my design preferences (colors, typography, spacing)
3. Includes both light (:root) and dark ([data-theme="dark"]) modes
4. Maintains WCAG AA accessibility contrast
5. Keeps responsive breakpoints

Write the CSS to custom.css with clear comments.

## Step 4: Tell Me How to Use It

After saving custom.css, show me these commands:

volcano ./docs --css ./custom.css -o ./public
volcano -s -p 8080 ./public

Ready to start! Ask me your first question.
```

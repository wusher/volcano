# Generate Documentation

Create a complete documentation site for your project using Volcano.

## About

This prompt generates a well-organized `docs/` folder with guides, tutorials, examples, and API reference documentation by analyzing your project's source code.

## Prompt

```
I need help creating documentation for my project that will be built with Volcano static site generator.

About my project:
- Project name: [YOUR PROJECT NAME]
- Description: [BRIEF PROJECT DESCRIPTION]
- Programming language/framework: [e.g., Go, Python, JavaScript, React]
- Target audience: [e.g., developers, end users, both]
- Source code: [GITHUB URL OR FILE PATH TO ANALYZE]

## Instructions

1. Analyze the source code to understand the project structure and public APIs

2. Create a docs/ folder with this structure:
   - index.md: Landing page with project overview, key benefits, quick links
   - getting-started.md: Installation, prerequisites, first example
   - guides/: Practical tutorials (basic usage, common use cases, best practices)
   - examples/: Real-world examples and code samples
   - api/: API reference with one file per major class/module containing:
     - Description and purpose
     - Constructor/initialization
     - Methods with signature, parameters, return values, examples
     - Properties/attributes
     - Types, interfaces, enums (for typed languages)

3. Format each markdown file with:
   - Clear H1 title
   - Well-organized H2/H3 sections
   - Markdown tables for parameters and return values
   - Code examples with syntax highlighting
   - Admonitions (:::note, :::warning, :::tip) for callouts
   - Cross-links to related pages

4. To build and preview the documentation:
   go install github.com/wusher/volcano@latest
   volcano ./docs -o ./public --title="[PROJECT NAME]"
   volcano -s -p 8080 ./public

Please analyze the source code and generate the documentation now.
```

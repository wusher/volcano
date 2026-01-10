# Generate API Documentation for Your Library

Use this prompt with an LLM to generate comprehensive API documentation that can be built with Volcano.

## Prompt

```
I need help creating API documentation for my library that will be built with Volcano static site generator.

About my library:
- Library name: [YOUR LIBRARY NAME]
- Programming language: [e.g., JavaScript, Python, Go, Ruby, etc.]
- Purpose: [WHAT THE LIBRARY DOES]
- Package manager: [e.g., npm, pip, go modules, rubygems, etc.]

Source code location:
[PROVIDE GITHUB URL, FILE PATH, OR PASTE KEY CODE SNIPPETS]

Please analyze my library and create a docs/api/ folder with comprehensive API documentation including:

1. **index.md** - API overview with:
   - Introduction to the API
   - Installation and setup
   - Quick start example
   - Authentication/initialization (if applicable)
   - Links to main API sections

2. **Classes or Modules** - One markdown file per major class/module:
   - Class/module description and purpose
   - Constructor/initialization
   - Methods with:
     - Method signature
     - Parameters (name, type, description, required/optional)
     - Return value (type and description)
     - Code examples
     - Common use cases
   - Properties/attributes
   - Events (if applicable)

3. **Functions** - For standalone functions:
   - Function signature
   - Parameters
   - Return values
   - Examples

4. **Types/Interfaces** - For typed languages:
   - Type definitions
   - Interface contracts
   - Enums/constants

5. **Error Handling** - Document:
   - Error types
   - Error codes
   - Exception handling patterns

Format requirements:
- Use clear markdown tables for parameters and return values
- Include realistic code examples with syntax highlighting
- Use admonitions (:::note, :::warning, :::tip) for important notes
- Cross-link related API methods
- Show both basic and advanced usage examples
- Include TypeScript definitions if applicable

Additional requirements:
[SPECIFY ANY SPECIAL FORMATTING, EXAMPLES, OR SECTIONS YOU NEED]
```

## Example Structure

For a JavaScript library, the structure might look like:

```
docs/api/
├── index.md              # API overview
├── client.md             # Client class documentation
├── authentication.md     # Auth methods
├── database.md           # Database operations
├── types.md              # TypeScript types
└── errors.md             # Error reference
```

## Tips

- Provide the actual source code or a link to your repository for accurate documentation
- Specify your preferred code example style (ES6, TypeScript, async/await, etc.)
- Mention any unique patterns or conventions your library uses
- Include version information if documenting a specific release
- Specify if you want migration guides from previous versions
- Indicate whether you want examples in multiple languages (if applicable)

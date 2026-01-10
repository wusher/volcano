# Generate a Docs Folder for Your Project

Use this prompt with an LLM to generate comprehensive documentation for your project using Volcano.

## Prompt

```
I need help creating a documentation folder for my project that will be built with Volcano static site generator.

About my project:
- Project name: [YOUR PROJECT NAME]
- Description: [BRIEF PROJECT DESCRIPTION]
- Programming language/framework: [e.g., Go, Python, JavaScript, React, etc.]
- Target audience: [e.g., developers, end users, both]

Please create a well-structured docs folder with the following:

1. **index.md** - A landing page with:
   - Project overview and key benefits
   - Quick links to important sections
   - Getting started call-to-action

2. **getting-started.md** - Installation and setup guide including:
   - Prerequisites
   - Installation steps
   - First example/hello world
   - Common troubleshooting

3. **guides/** folder with practical tutorials:
   - Basic usage guide
   - Common use cases
   - Best practices

4. **reference/** folder with technical documentation:
   - Configuration options
   - CLI commands (if applicable)
   - API reference structure

5. **examples/** folder with:
   - Real-world example scenarios
   - Code samples
   - Tips and tricks

Please use clear, concise markdown formatting. Include code blocks where appropriate, and use admonitions (:::note, :::warning, :::tip) for important callouts.

Each file should have:
- A clear H1 title
- Well-organized sections with H2/H3 headings
- Practical examples
- Links to related documentation

Additional context about my project:
[ADD ANY SPECIFIC REQUIREMENTS OR DETAILS HERE]
```

## Tips

- Be specific about your project's unique features and requirements
- Mention any existing documentation you want to migrate or reference
- Specify your preferred documentation style (tutorial-heavy, reference-heavy, etc.)
- Include information about your project's complexity level
- Mention if you need multilingual support or specific technical terms

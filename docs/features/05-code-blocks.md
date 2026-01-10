# Code Blocks

Display code with syntax highlighting and interactive features.

## Basic Code Blocks

Create a code block with triple backticks:

````markdown
```
Plain code block
No syntax highlighting
```
````

## Syntax Highlighting

Specify the language after the opening backticks:

````markdown
```javascript
function greet(name) {
  return `Hello, ${name}!`;
}
```
````

```javascript
function greet(name) {
  return `Hello, ${name}!`;
}
```

## Supported Languages

Volcano uses [Chroma](https://github.com/alecthomas/chroma) for highlighting, supporting 200+ languages including:

| Language | Identifier |
|----------|------------|
| JavaScript | `javascript`, `js` |
| TypeScript | `typescript`, `ts` |
| Python | `python`, `py` |
| Go | `go`, `golang` |
| Rust | `rust`, `rs` |
| Ruby | `ruby`, `rb` |
| Java | `java` |
| C | `c` |
| C++ | `cpp`, `c++` |
| C# | `csharp`, `c#` |
| PHP | `php` |
| Swift | `swift` |
| Kotlin | `kotlin` |
| HTML | `html` |
| CSS | `css` |
| SCSS | `scss` |
| JSON | `json` |
| YAML | `yaml`, `yml` |
| Markdown | `markdown`, `md` |
| Bash | `bash`, `sh`, `shell` |
| SQL | `sql` |
| Docker | `dockerfile` |

See the [full list](https://github.com/alecthomas/chroma#supported-languages) for all supported languages.

## Line Highlighting

Highlight specific lines by adding line numbers in curly braces:

````markdown
```go {3,5-7}
package main

import "fmt"  // Highlighted

func main() {         // Highlighted
    fmt.Println("Hello")  // Highlighted
    fmt.Println("World")  // Highlighted
}
```
````

```go {3,5-7}
package main

import "fmt"  // Highlighted

func main() {         // Highlighted
    fmt.Println("Hello")  // Highlighted
    fmt.Println("World")  // Highlighted
}
```

### Highlight Syntax

- Single line: `{3}`
- Multiple lines: `{1,3,5}`
- Line range: `{5-10}`
- Combined: `{1,3,5-7,10}`

## Copy Button

All code blocks include a copy button that appears on hover. Click to copy the code to your clipboard.

The button shows "Copied" briefly after clicking.

## Inline Code

Use single backticks for inline code:

```markdown
Use the `console.log()` function for debugging.
```

Use the `console.log()` function for debugging.

Inline code is styled differently from code blocks — typically with a background color and monospace font.

## Examples by Language

### JavaScript

```javascript
// ES6 arrow function
const add = (a, b) => a + b;

// Async/await
async function fetchData(url) {
  const response = await fetch(url);
  return response.json();
}
```

### Python

```python
# List comprehension
squares = [x**2 for x in range(10)]

# Context manager
with open('file.txt', 'r') as f:
    content = f.read()
```

### Go

```go
package main

import "fmt"

func main() {
    messages := make(chan string)

    go func() {
        messages <- "Hello"
    }()

    msg := <-messages
    fmt.Println(msg)
}
```

### Bash

```bash
#!/bin/bash

# Variables
NAME="World"

# Function
greet() {
    echo "Hello, $1!"
}

greet "$NAME"
```

### HTML

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Example</title>
</head>
<body>
    <h1>Hello, World!</h1>
</body>
</html>
```

### CSS

```css
.container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
}

.button {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    border: none;
    padding: 12px 24px;
    border-radius: 8px;
}
```

### JSON

```json
{
  "name": "my-project",
  "version": "1.0.0",
  "dependencies": {
    "lodash": "^4.17.21"
  }
}
```

### YAML

```yaml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: npm install
      - run: npm run build
```

## Code in Admonitions

Code blocks work inside admonitions:

````markdown
:::tip
Install dependencies first:

```bash
npm install
```
:::
````

:::tip
Install dependencies first:

```bash
npm install
```
:::

## Styling

Code blocks use CSS classes from Chroma. Key classes:

```css
.code-block { }          /* Wrapper with copy button */
.copy-button { }         /* Copy button */
.chroma { }              /* Chroma syntax container */
pre { }                  /* Pre-formatted block */
code { }                 /* Code element */

/* Syntax token classes */
.chroma .k { }           /* Keyword */
.chroma .s { }           /* String */
.chroma .c { }           /* Comment */
.chroma .nf { }          /* Function name */
.chroma .m { }           /* Number */
```

## Best Practices

### Specify the Language

Always specify a language for syntax highlighting:

````markdown
<!-- Good -->
```python
print("Hello")
```

<!-- Less useful -->
```
print("Hello")
```
````

### Keep Examples Concise

Show just enough code to illustrate the point:

````markdown
<!-- Good: Focused -->
```javascript
const result = items.filter(item => item.active);
```

<!-- Too much: Includes irrelevant setup -->
```javascript
const express = require('express');
const app = express();
// ... 50 lines of setup ...
const result = items.filter(item => item.active);
```
````

### Use Line Highlighting for Focus

When explaining specific lines:

````markdown
The key part is the filter callback:

```javascript {2}
const activeItems = items
  .filter(item => item.active)  // This line filters
  .map(item => item.name);
```
````

## Related

- [[markdown-syntax]] — All markdown features
- [[admonitions]] — Callout boxes for code examples
- [[theming]] — Customize code block appearance

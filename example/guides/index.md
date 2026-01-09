# Guides

Welcome to the Volcano guides section. Here you'll find detailed documentation on all aspects of using Volcano.

## Available Guides

### Installation

Learn how to install Volcano on your system.

- [Installation Guide](installation/)

### Configuration

Understand all the command-line options and how to customize your site.

- [Configuration Guide](configuration/)

## Code Examples

Here's a simple Go program:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Volcano!")
}
```

And some JavaScript:

```javascript
// Theme toggle
document.addEventListener('DOMContentLoaded', () => {
    const toggle = document.querySelector('.theme-toggle');
    toggle.addEventListener('click', () => {
        const current = document.documentElement.getAttribute('data-theme');
        const next = current === 'dark' ? 'light' : 'dark';
        document.documentElement.setAttribute('data-theme', next);
    });
});
```

# API Reference

This section documents Volcano's internal APIs for developers who want to extend or embed Volcano in their own projects.

## Package Overview

Volcano is organized into several internal packages:

| Package | Description |
|---------|-------------|
| `cmd` | Command implementations |
| `internal/generator` | Site generation engine |
| `internal/markdown` | Markdown parsing |
| `internal/templates` | HTML templates |
| `internal/tree` | File tree building |
| `internal/styles` | Embedded CSS |
| `internal/server` | HTTP server |
| `internal/output` | Colored logging |

## Usage as a Library

While Volcano is primarily a CLI tool, you can use its packages in your own Go programs:

```go
package main

import (
    "os"
    "volcano/internal/generator"
)

func main() {
    config := generator.Config{
        InputDir:  "./docs",
        OutputDir: "./public",
        Title:     "My Site",
    }

    gen, err := generator.New(config, os.Stdout)
    if err != nil {
        panic(err)
    }

    result, err := gen.Generate()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Generated %d pages\n", result.PagesGenerated)
}
```

## Learn More

- [Endpoints](endpoints/) - HTTP server endpoints
- [Authentication](authentication/) - (Future) Auth features

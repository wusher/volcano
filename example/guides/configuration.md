# Configuration

Volcano is configured entirely through command-line flags. No configuration files are needed!

## Command Overview

### Generate Mode (default)

```bash
volcano <input-folder> [flags]
```

### Serve Mode

```bash
volcano -s <folder> [flags]
```

## Available Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output directory | `./output` |
| `--serve` | `-s` | Run in serve mode | `false` |
| `--port` | `-p` | Server port | `1776` |
| `--title` | | Site title | `My Site` |
| `--quiet` | `-q` | Suppress output | `false` |
| `--verbose` | | Debug output | `false` |
| `--version` | `-v` | Show version | |
| `--help` | `-h` | Show help | |

## Examples

### Basic Generation

```bash
volcano ./docs
```

Generates site from `./docs` to `./output`.

### Custom Output Directory

```bash
volcano ./docs -o ./public
```

### With Custom Title

```bash
volcano ./docs --title="My Project Docs"
```

### Serve Generated Site

```bash
volcano -s ./output
```

### Serve on Custom Port

```bash
volcano -s -p 8080 ./output
```

### Quiet Mode

```bash
volcano -q ./docs
```

Suppresses all output except errors.

### Verbose Mode

```bash
volcano --verbose ./docs
```

Shows detailed debug information.

## Environment Variables

Currently, Volcano does not use any environment variables. All configuration is done through flags.

## Tips

1. **Clean URLs**: Volcano automatically generates clean URLs. No trailing `.html` needed!

2. **Index Files**: Create `index.md` in folders to set the folder's landing page.

3. **Alphabetical Order**: Pages are sorted alphabetically. Use numbered prefixes to control order.

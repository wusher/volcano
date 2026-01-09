# Installation

This guide covers all the ways to install Volcano on your system.

## Prerequisites

- Go 1.21 or later (for building from source)
- A terminal/command line

## Option 1: Go Install

The easiest way to install Volcano is using `go install`:

```bash
go install github.com/example/volcano@latest
```

This will download, compile, and install the `volcano` binary to your `$GOPATH/bin` directory.

## Option 2: Download Binary

Pre-built binaries are available for:

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

Download from the [releases page](https://github.com/example/volcano/releases).

## Option 3: Build from Source

Clone the repository and build:

```bash
git clone https://github.com/example/volcano.git
cd volcano
go build -o volcano .
```

Move the binary to your PATH:

```bash
sudo mv volcano /usr/local/bin/
```

## Verify Installation

Check that Volcano is installed correctly:

```bash
volcano --version
```

You should see output like:

```
volcano version 0.1.0
```

## Next Steps

Now that you have Volcano installed, check out the [Getting Started](../getting-started/) guide.

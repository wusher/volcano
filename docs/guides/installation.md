# Installation

Install Volcano on your system using one of these methods.

## Using Go Install (Recommended)

If you have Go 1.21+ installed, this is the simplest method:

```bash
go install github.com/wusher/volcano@latest
```

This installs the `volcano` binary to your `$GOPATH/bin` directory (usually `~/go/bin`).

Make sure your Go bin directory is in your PATH:

```bash
# Add to ~/.bashrc, ~/.zshrc, or equivalent
export PATH="$PATH:$(go env GOPATH)/bin"
```

Verify the installation:

```bash
volcano --version
```

## Building from Source

Clone the repository and build:

```bash
git clone https://github.com/wusher/volcano.git
cd volcano
go build -o volcano .
```

Move the binary to a location in your PATH:

```bash
# Linux/macOS
sudo mv volcano /usr/local/bin/

# Or add to your local bin
mv volcano ~/bin/
```

## Verifying Installation

After installation, verify Volcano is working:

```bash
volcano --version
```

You should see output like:

```
volcano version 0.1.0
```

View available commands and options:

```bash
volcano --help
```

## Updating Volcano

To update to the latest version using Go:

```bash
go install github.com/wusher/volcano@latest
```

If you built from source:

```bash
cd volcano
git pull
go build -o volcano .
```

## Troubleshooting

### Command not found

If you get "command not found" after installing with `go install`:

1. Check that Go's bin directory is in your PATH:
   ```bash
   echo $PATH | grep -o "$(go env GOPATH)/bin"
   ```

2. If it's not there, add it to your shell profile and reload:
   ```bash
   echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
   source ~/.bashrc
   ```

### Permission denied

If you get permission errors when moving the binary:

```bash
# Use sudo for system directories
sudo mv volcano /usr/local/bin/

# Or use a user-owned directory
mkdir -p ~/bin
mv volcano ~/bin/
```

## Next Steps

Now that Volcano is installed, [[creating-your-first-site|create your first site]].

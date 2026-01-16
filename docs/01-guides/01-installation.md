# Installation

Requires Go 1.24+.

## Install

```bash
go install github.com/wusher/volcano@latest
```

Ensure `$(go env GOPATH)/bin` is in your PATH.

## Verify

```bash
volcano --version
```

## Update

```bash
go install github.com/wusher/volcano@latest
```

## Build from Source

```bash
git clone https://github.com/wusher/volcano.git
cd volcano
go build -o volcano .
sudo mv volcano /usr/local/bin/
```

## Next

[[creating-your-first-site|Create your first site]].

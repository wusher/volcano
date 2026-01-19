// Package main provides the entry point for the volcano CLI.
package main

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/wusher/volcano/cmd"
	"github.com/wusher/volcano/internal/output"
)

// Version is the CLI version (overridden at build time for releases)
var Version = "dev"

func init() {
	// If Version wasn't set via ldflags, try to get it from build info
	// This works when installed via `go install github.com/wusher/volcano@v1.0.0`
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

// Run executes the CLI with the given arguments and writers
func Run(args []string, stdout, stderr io.Writer) int {
	// Handle no arguments
	if len(args) == 0 {
		errLogger := output.NewLogger(stderr, output.IsStderrTTY(), false, false)
		errLogger.Error("input folder is required")
		_, _ = fmt.Fprintln(stderr, "")
		printUsage(stderr)
		return 1
	}

	// Handle version flag at top level
	if args[0] == "-v" || args[0] == "--version" {
		_, _ = fmt.Fprintf(stdout, "volcano version %s\n", Version)
		return 0
	}

	// Handle help flag at top level
	if args[0] == "-h" || args[0] == "--help" {
		printUsage(stdout)
		return 0
	}

	// Check for subcommands
	var err error
	switch args[0] {
	case "css":
		err = cmd.CSS(args[1:], stdout)
	case "build":
		err = cmd.Build(args[1:], stdout, stderr)
	case "serve", "server":
		err = cmd.ServeCommand(args[1:], stdout, stderr)
	case "init":
		err = cmd.Init(args[1:], stdout, stderr)
	default:
		// Fall through: treat as shorthand for build (backward compatibility)
		// This allows `volcano ./docs` to work like `volcano build ./docs`
		err = cmd.Build(args, stdout, stderr)
	}

	if err != nil {
		// Error already logged by the command
		return 1
	}

	return 0
}

func printUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Volcano - Static site generator")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Usage:")
	_, _ = fmt.Fprintln(w, "  volcano build [flags] <input>    Generate static site")
	_, _ = fmt.Fprintln(w, "  volcano serve [flags] <input>    Start development server")
	_, _ = fmt.Fprintln(w, "  volcano server [flags] <input>   Alias for serve")
	_, _ = fmt.Fprintln(w, "  volcano init [flags]             Create/update volcano.json config")
	_, _ = fmt.Fprintln(w, "  volcano css [-o file]            Output vanilla CSS")
	_, _ = fmt.Fprintln(w, "  volcano <input>                  Shorthand for build")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Run 'volcano <command> --help' for command-specific help.")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Flags:")
	_, _ = fmt.Fprintln(w, "  -v, --version        Show version")
	_, _ = fmt.Fprintln(w, "  -h, --help           Show help")
}

// ValidateInputDir checks if the given path is a valid directory
// Kept for backward compatibility with tests
func ValidateInputDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("input directory does not exist: %s", path)
		}
		return fmt.Errorf("cannot access input directory: %v", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("input path is not a directory: %s", path)
	}

	return nil
}

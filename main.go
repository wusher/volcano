// Package main provides the entry point for the volcano CLI.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"volcano/cmd"
)

// Version is the CLI version
const Version = "0.1.0"

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}

// Run executes the CLI with the given arguments and writers
func Run(args []string, stdout, stderr io.Writer) int {
	cfg := cmd.DefaultConfig()
	exitCode, _ := runWithConfig(args, cfg, stdout, stderr)
	return exitCode
}

// runWithConfig is the internal implementation that returns both exit code and error
func runWithConfig(args []string, cfg *cmd.Config, stdout, stderr io.Writer) (int, error) {
	// Create a new FlagSet to avoid issues with flag.Parse being called multiple times
	fs := flag.NewFlagSet("volcano", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// Define flags
	var showVersion bool
	var showHelp bool

	fs.StringVar(&cfg.OutputDir, "o", cfg.OutputDir, "Output directory for generated site")
	fs.StringVar(&cfg.OutputDir, "output", cfg.OutputDir, "Output directory for generated site")
	fs.BoolVar(&cfg.ServeMode, "s", false, "Run in serve mode")
	fs.BoolVar(&cfg.ServeMode, "serve", false, "Run in serve mode")
	fs.IntVar(&cfg.Port, "p", cfg.Port, "Port for the HTTP server (serve mode)")
	fs.IntVar(&cfg.Port, "port", cfg.Port, "Port for the HTTP server (serve mode)")
	fs.StringVar(&cfg.Title, "title", cfg.Title, "Site title")
	fs.BoolVar(&cfg.Quiet, "q", false, "Suppress non-error output")
	fs.BoolVar(&cfg.Quiet, "quiet", false, "Suppress non-error output")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Enable debug output")
	fs.BoolVar(&showVersion, "v", false, "Show version")
	fs.BoolVar(&showVersion, "version", false, "Show version")
	fs.BoolVar(&showHelp, "h", false, "Show help")
	fs.BoolVar(&showHelp, "help", false, "Show help")

	// Custom usage message
	fs.Usage = func() {
		printUsage(stdout)
	}

	if err := fs.Parse(args); err != nil {
		return 1, err
	}

	// Handle version flag
	if showVersion {
		_, _ = fmt.Fprintf(stdout, "volcano version %s\n", Version)
		return 0, nil
	}

	// Handle help flag
	if showHelp {
		printUsage(stdout)
		return 0, nil
	}

	// Get input directory from positional arguments
	remainingArgs := fs.Args()
	if len(remainingArgs) < 1 {
		_, _ = fmt.Fprintln(stderr, "Error: input folder is required")
		_, _ = fmt.Fprintln(stderr, "")
		printUsage(stderr)
		return 1, fmt.Errorf("input folder is required")
	}

	cfg.InputDir = remainingArgs[0]

	// Validate input directory exists and is a directory
	if err := ValidateInputDir(cfg.InputDir); err != nil {
		_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1, err
	}

	// Execute the appropriate command
	var err error
	if cfg.ServeMode {
		err = cmd.Serve(cfg, stdout)
	} else {
		err = cmd.Generate(cfg, stdout)
	}

	if err != nil {
		_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1, err
	}

	return 0, nil
}

func printUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "volcano - A static site generator")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Usage:")
	_, _ = fmt.Fprintln(w, "  volcano <input-folder> [flags]     Generate static site")
	_, _ = fmt.Fprintln(w, "  volcano -s <folder> [flags]        Serve static site")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Flags:")
	_, _ = fmt.Fprintln(w, "  -o, --output <dir>   Output directory (default: ./output)")
	_, _ = fmt.Fprintln(w, "  -s, --serve          Run in serve mode")
	_, _ = fmt.Fprintln(w, "  -p, --port <port>    Server port (default: 1776)")
	_, _ = fmt.Fprintln(w, "  --title <title>      Site title (default: My Site)")
	_, _ = fmt.Fprintln(w, "  -q, --quiet          Suppress non-error output")
	_, _ = fmt.Fprintln(w, "  --verbose            Enable debug output")
	_, _ = fmt.Fprintln(w, "  -v, --version        Show version")
	_, _ = fmt.Fprintln(w, "  -h, --help           Show help")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Examples:")
	_, _ = fmt.Fprintln(w, "  volcano ./docs -o ./public --title=\"My Docs\"")
	_, _ = fmt.Fprintln(w, "  volcano -s -p 8080 ./public")
}

// ValidateInputDir checks if the given path is a valid directory
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

// Package main provides the entry point for the volcano CLI.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/wusher/volcano/cmd"
	"github.com/wusher/volcano/internal/output"
	"github.com/wusher/volcano/internal/styles"
)

// Version is the CLI version (overridden at build time for releases)
var Version = "dev"

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
	// Check for subcommands before flag parsing
	if len(args) > 0 && args[0] == "css" {
		if err := cmd.CSS(args[1:], stdout); err != nil {
			errLogger := output.NewLogger(stderr, output.IsStderrTTY(), false, false)
			errLogger.Error("%v", err)
			return 1, err
		}
		return 0, nil
	}

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
	fs.StringVar(&cfg.SiteURL, "url", "", "Site base URL for SEO")
	fs.StringVar(&cfg.Author, "author", "", "Site author")
	fs.StringVar(&cfg.OGImage, "og-image", "", "Default Open Graph image URL")
	fs.StringVar(&cfg.FaviconPath, "favicon", "", "Path to favicon file")
	fs.BoolVar(&cfg.ShowLastMod, "last-modified", false, "Show last modified date")
	fs.BoolVar(&cfg.TopNav, "top-nav", false, "Display root files in top navigation bar")
	fs.BoolVar(&cfg.ShowPageNav, "page-nav", false, "Show previous/next page navigation")
	fs.StringVar(&cfg.Theme, "theme", "docs", "Theme name (docs, blog, vanilla)")
	fs.StringVar(&cfg.CSSPath, "css", "", "Path to custom CSS file")
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

	// Reorder args to put flags first (Go's flag package stops at first non-flag)
	args = reorderArgs(args)

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

	// Detect if stderr is a terminal for colored output
	errLogger := output.NewLogger(stderr, output.IsStderrTTY(), false, false)

	// Get input directory from positional arguments
	remainingArgs := fs.Args()
	if len(remainingArgs) < 1 {
		errLogger.Error("input folder is required")
		_, _ = fmt.Fprintln(stderr, "")
		printUsage(stderr)
		return 1, fmt.Errorf("input folder is required")
	}

	cfg.InputDir = remainingArgs[0]

	// Validate input directory exists and is a directory
	if err := ValidateInputDir(cfg.InputDir); err != nil {
		errLogger.Error("%v", err)
		return 1, err
	}

	// Validate theme
	if err := styles.ValidateTheme(cfg.Theme); err != nil {
		errLogger.Error("%v", err)
		return 1, err
	}

	// Validate custom CSS path if provided
	if cfg.CSSPath != "" {
		if _, err := os.Stat(cfg.CSSPath); err != nil {
			if os.IsNotExist(err) {
				errLogger.Error("CSS file not found: %s", cfg.CSSPath)
				return 1, fmt.Errorf("CSS file not found: %s", cfg.CSSPath)
			}
			errLogger.Error("cannot access CSS file: %v", err)
			return 1, err
		}
		// When using custom CSS, ignore the theme flag
		cfg.Theme = ""
	}

	// Set colored output based on stdout TTY detection
	cfg.Colored = output.IsStdoutTTY()

	// Execute the appropriate command
	var err error
	if cfg.ServeMode {
		err = cmd.Serve(cfg, stdout)
	} else {
		err = cmd.Generate(cfg, stdout)
	}

	if err != nil {
		errLogger.Error("%v", err)
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
	_, _ = fmt.Fprintln(w, "  volcano css [-o file]              Output vanilla CSS skeleton")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Flags:")
	_, _ = fmt.Fprintln(w, "  -o, --output <dir>   Output directory (default: ./output)")
	_, _ = fmt.Fprintln(w, "  -s, --serve          Run in serve mode")
	_, _ = fmt.Fprintln(w, "  -p, --port <port>    Server port (default: 1776)")
	_, _ = fmt.Fprintln(w, "  --title <title>      Site title (default: My Site)")
	_, _ = fmt.Fprintln(w, "  --url <url>          Site base URL for SEO")
	_, _ = fmt.Fprintln(w, "  --author <name>      Site author")
	_, _ = fmt.Fprintln(w, "  --og-image <url>     Default Open Graph image URL")
	_, _ = fmt.Fprintln(w, "  --favicon <path>     Path to favicon file (.ico, .png, .svg)")
	_, _ = fmt.Fprintln(w, "  --last-modified      Show last modified date on pages")
	_, _ = fmt.Fprintln(w, "  --top-nav            Display root files in top navigation bar")
	_, _ = fmt.Fprintln(w, "  --page-nav           Show previous/next page navigation")
	_, _ = fmt.Fprintln(w, "  --theme <name>       Theme (docs, blog, vanilla; default: docs)")
	_, _ = fmt.Fprintln(w, "  --css <path>         Path to custom CSS file (overrides theme)")
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

// reorderArgs moves flags before positional arguments
// This is needed because Go's flag package stops at the first non-flag argument
func reorderArgs(args []string) []string {
	var flags []string
	var positional []string

	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			// This is a flag
			flags = append(flags, arg)
			// Check if this flag takes a value (not a boolean flag)
			// Flags with = don't need special handling
			if !strings.Contains(arg, "=") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				// Check if the next arg looks like a flag value (not a path)
				nextArg := args[i+1]
				// Only treat as value if flag is a known value-taking flag
				if isValueFlag(arg) {
					i++
					flags = append(flags, nextArg)
				}
			}
		} else {
			positional = append(positional, arg)
		}
		i++
	}

	return append(flags, positional...)
}

// isValueFlag returns true if the flag takes a value argument
func isValueFlag(flag string) bool {
	// Strip leading dashes
	name := strings.TrimLeft(flag, "-")
	// List of flags that take values
	valueFlags := map[string]bool{
		"o": true, "output": true,
		"p": true, "port": true,
		"title": true, "url": true, "author": true,
		"og-image": true, "favicon": true,
		"theme": true, "css": true,
	}
	return valueFlags[name]
}

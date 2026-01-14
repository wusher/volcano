package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/wusher/volcano/internal/output"
	"github.com/wusher/volcano/internal/styles"
)

// Build handles the build subcommand for generating static sites
func Build(args []string, stdout, stderr io.Writer) error {
	cfg := DefaultConfig()

	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// Define build-specific flags
	var showHelp bool
	var viewTransitionsFlag bool // Deprecated flag, kept for backwards compatibility

	fs.StringVar(&cfg.OutputDir, "o", cfg.OutputDir, "Output directory for generated site")
	fs.StringVar(&cfg.OutputDir, "output", cfg.OutputDir, "Output directory for generated site")
	fs.StringVar(&cfg.Title, "title", cfg.Title, "Site title")
	fs.StringVar(&cfg.SiteURL, "url", "", "Site base URL for SEO")
	fs.StringVar(&cfg.Author, "author", "", "Site author")
	fs.StringVar(&cfg.OGImage, "og-image", "", "Default Open Graph image URL")
	fs.StringVar(&cfg.FaviconPath, "favicon", "", "Path to favicon file")
	fs.BoolVar(&cfg.ShowLastMod, "last-modified", false, "Show last modified date")
	fs.BoolVar(&cfg.TopNav, "top-nav", false, "Display root files in top navigation bar")
	fs.BoolVar(&cfg.ShowPageNav, "page-nav", false, "Show previous/next page navigation")
	fs.BoolVar(&cfg.ShowBreadcrumbs, "breadcrumbs", cfg.ShowBreadcrumbs, "Show breadcrumb navigation")
	fs.StringVar(&cfg.Theme, "theme", "docs", "Theme name (docs, blog, vanilla)")
	fs.StringVar(&cfg.CSSPath, "css", "", "Path to custom CSS file")
	fs.StringVar(&cfg.AccentColor, "accent-color", "", "Custom accent color (hex format, e.g., '#ff6600')")
	fs.BoolVar(&cfg.InstantNav, "instant-nav", false, "Enable instant navigation with hover prefetching")
	fs.BoolVar(&cfg.InlineAssets, "inline-assets", false, "Embed CSS/JS inline instead of external files")
	fs.BoolVar(&cfg.PWA, "pwa", false, "Enable PWA manifest and service worker for offline support")
	fs.BoolVar(&viewTransitionsFlag, "view-transitions", false, "Deprecated: view transitions are now enabled by default")
	fs.BoolVar(&cfg.Quiet, "q", false, "Suppress non-error output")
	fs.BoolVar(&cfg.Quiet, "quiet", false, "Suppress non-error output")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Enable debug output")
	fs.BoolVar(&showHelp, "h", false, "Show help")
	fs.BoolVar(&showHelp, "help", false, "Show help")

	fs.Usage = func() {
		printBuildUsage(stdout)
	}

	// Reorder args to put flags first (Go's flag package stops at first non-flag)
	args = reorderArgs(args, buildValueFlags)

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Handle help flag
	if showHelp {
		printBuildUsage(stdout)
		return nil
	}

	errLogger := output.NewLogger(stderr, output.IsStderrTTY(), false, false)

	// Check if deprecated --view-transitions flag was used
	if viewTransitionsFlag {
		errLogger.Warning("--view-transitions is deprecated: view transitions are now enabled by default")
	}

	// Get input directory from positional arguments
	remainingArgs := fs.Args()
	if len(remainingArgs) < 1 {
		errLogger.Error("input folder is required")
		_, _ = fmt.Fprintln(stderr, "")
		printBuildUsage(stderr)
		return fmt.Errorf("input folder is required")
	}

	cfg.InputDir = remainingArgs[0]

	// Validate input directory exists and is a directory
	if err := validateInputDir(cfg.InputDir); err != nil {
		errLogger.Error("%v", err)
		return err
	}

	// Validate theme
	if err := styles.ValidateTheme(cfg.Theme); err != nil {
		errLogger.Error("%v", err)
		return err
	}

	// Validate custom CSS path if provided
	if cfg.CSSPath != "" {
		if _, err := os.Stat(cfg.CSSPath); err != nil {
			if os.IsNotExist(err) {
				errLogger.Error("CSS file not found: %s", cfg.CSSPath)
				return fmt.Errorf("CSS file not found: %s", cfg.CSSPath)
			}
			errLogger.Error("cannot access CSS file: %v", err)
			return err
		}
		// When using custom CSS, ignore the theme flag
		cfg.Theme = ""
	}

	// Set colored output based on stdout TTY detection
	cfg.Colored = output.IsStdoutTTY()

	if err := Generate(cfg, stdout); err != nil {
		errLogger.Error("%v", err)
		return err
	}
	return nil
}

// validateInputDir checks if the given path is a valid directory
func validateInputDir(path string) error {
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

func printBuildUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Generate a static site from markdown files")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Usage:")
	_, _ = fmt.Fprintln(w, "  volcano build [flags] <input>")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Output:")
	_, _ = fmt.Fprintln(w, "  -o, --output <dir>   Output directory (default: ./output)")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Site Configuration:")
	_, _ = fmt.Fprintln(w, "  --title <title>      Site title (default: My Site)")
	_, _ = fmt.Fprintln(w, "  --url <url>          Base URL for canonical links and SEO")
	_, _ = fmt.Fprintln(w, "  --author <name>      Site author for meta tags")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Appearance:")
	_, _ = fmt.Fprintln(w, "  --theme <name>       Theme: docs, blog, vanilla (default: docs)")
	_, _ = fmt.Fprintln(w, "  --css <path>         Custom CSS file (overrides theme)")
	_, _ = fmt.Fprintln(w, "  --accent-color <hex> Custom accent color (e.g., '#ff6600')")
	_, _ = fmt.Fprintln(w, "  --favicon <path>     Favicon file (.ico, .png, .svg)")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Navigation:")
	_, _ = fmt.Fprintln(w, "  --top-nav            Show root files in top navigation bar")
	_, _ = fmt.Fprintln(w, "  --breadcrumbs        Show breadcrumb trail (default: true)")
	_, _ = fmt.Fprintln(w, "  --page-nav           Show previous/next page links")
	_, _ = fmt.Fprintln(w, "  --instant-nav        Enable hover prefetching for faster navigation")
	_, _ = fmt.Fprintln(w, "  --inline-assets      Embed CSS/JS inline instead of external files")
	_, _ = fmt.Fprintln(w, "  --pwa                Enable PWA manifest and service worker for offline support")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Content:")
	_, _ = fmt.Fprintln(w, "  --last-modified      Show last modified date on pages")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "SEO:")
	_, _ = fmt.Fprintln(w, "  --og-image <path>    Default Open Graph image")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Logging:")
	_, _ = fmt.Fprintln(w, "  -q, --quiet          Suppress non-error output")
	_, _ = fmt.Fprintln(w, "  --verbose            Show detailed build information")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Other:")
	_, _ = fmt.Fprintln(w, "  -h, --help           Show this help message")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Examples:")
	_, _ = fmt.Fprintln(w, "  volcano build ./docs -o ./public --title=\"My Docs\"")
	_, _ = fmt.Fprintln(w, "  volcano build --theme=blog --accent-color='#ff6600' ./posts")
	_, _ = fmt.Fprintln(w, "  volcano build --top-nav --page-nav --instant-nav ./docs")
}

// buildValueFlags is the set of flags that take values for the build command
var buildValueFlags = map[string]bool{
	"o": true, "output": true,
	"title": true, "url": true, "author": true,
	"og-image": true, "favicon": true,
	"theme": true, "css": true, "accent-color": true,
}

// reorderArgs moves flags before positional arguments
// This is needed because Go's flag package stops at the first non-flag argument
func reorderArgs(args []string, valueFlags map[string]bool) []string {
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
				if isValueFlagInSet(arg, valueFlags) {
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

// isValueFlagInSet returns true if the flag takes a value argument
func isValueFlagInSet(flag string, valueFlags map[string]bool) bool {
	// Strip leading dashes
	name := strings.TrimLeft(flag, "-")
	return valueFlags[name]
}

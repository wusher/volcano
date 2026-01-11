package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wusher/volcano/internal/output"
	"github.com/wusher/volcano/internal/server"
	"github.com/wusher/volcano/internal/styles"
)

// ServeCommand handles the serve subcommand for starting the development server
func ServeCommand(args []string, stdout, stderr io.Writer) error {
	cfg := DefaultConfig()

	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// Define serve-specific flags
	var showHelp bool

	fs.IntVar(&cfg.Port, "p", cfg.Port, "Server port")
	fs.IntVar(&cfg.Port, "port", cfg.Port, "Server port")
	fs.StringVar(&cfg.Title, "title", cfg.Title, "Site title")
	fs.StringVar(&cfg.Theme, "theme", "docs", "Theme name (docs, blog, vanilla)")
	fs.StringVar(&cfg.CSSPath, "css", "", "Path to custom CSS file")
	fs.StringVar(&cfg.AccentColor, "accent-color", "", "Custom accent color (hex format, e.g., '#ff6600')")
	fs.StringVar(&cfg.FaviconPath, "favicon", "", "Path to favicon file")
	fs.BoolVar(&cfg.TopNav, "top-nav", false, "Display root files in top navigation bar")
	fs.BoolVar(&cfg.ShowPageNav, "page-nav", false, "Show previous/next page navigation")
	fs.BoolVar(&cfg.ShowBreadcrumbs, "breadcrumbs", cfg.ShowBreadcrumbs, "Show breadcrumb navigation")
	fs.BoolVar(&cfg.InstantNav, "instant-nav", false, "Enable instant navigation with hover prefetching")
	fs.BoolVar(&cfg.ViewTransitions, "view-transitions", false, "Enable browser view transitions API")
	fs.BoolVar(&cfg.Quiet, "q", false, "Suppress non-error output")
	fs.BoolVar(&cfg.Quiet, "quiet", false, "Suppress non-error output")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Enable debug output")
	fs.BoolVar(&showHelp, "h", false, "Show help")
	fs.BoolVar(&showHelp, "help", false, "Show help")

	fs.Usage = func() {
		printServeUsage(stdout)
	}

	// Reorder args to put flags first (Go's flag package stops at first non-flag)
	args = reorderArgs(args, serveValueFlags)

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Handle help flag
	if showHelp {
		printServeUsage(stdout)
		return nil
	}

	errLogger := output.NewLogger(stderr, output.IsStderrTTY(), false, false)

	// Get input directory from positional arguments
	remainingArgs := fs.Args()
	if len(remainingArgs) < 1 {
		errLogger.Error("input folder is required")
		_, _ = fmt.Fprintln(stderr, "")
		printServeUsage(stderr)
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

	// Validate favicon path if provided
	if cfg.FaviconPath != "" {
		if _, err := os.Stat(cfg.FaviconPath); err != nil {
			if os.IsNotExist(err) {
				errLogger.Error("favicon file not found: %s", cfg.FaviconPath)
				return fmt.Errorf("favicon file not found: %s", cfg.FaviconPath)
			}
			errLogger.Error("cannot access favicon file: %v", err)
			return err
		}
	}

	if err := Serve(cfg, stdout); err != nil {
		errLogger.Error("%v", err)
		return err
	}
	return nil
}

// Serve starts the HTTP server to preview the site
// If the input directory is a source directory (contains .md files), it uses
// dynamic rendering so changes are reflected immediately without restart.
// Otherwise, it serves static files from the directory.
func Serve(cfg *Config, w io.Writer) error {
	// Check if this is a source directory (contains .md files but no index.html)
	if isSourceDirectory(cfg.InputDir) {
		// Use dynamic server for live rendering
		dynamicCfg := server.DynamicConfig{
			SourceDir:       cfg.InputDir,
			Title:           cfg.Title,
			Port:            cfg.Port,
			Quiet:           cfg.Quiet,
			Verbose:         cfg.Verbose,
			TopNav:          cfg.TopNav,
			ShowPageNav:     cfg.ShowPageNav,
			ShowBreadcrumbs: cfg.ShowBreadcrumbs,
			Theme:           cfg.Theme,
			CSSPath:         cfg.CSSPath,
			AccentColor:     cfg.AccentColor,
			InstantNav:      cfg.InstantNav,
			ViewTransitions: cfg.ViewTransitions,
			FaviconPath:     cfg.FaviconPath,
		}

		srv, err := server.NewDynamicServer(dynamicCfg, w)
		if err != nil {
			return err
		}

		return srv.Start()
	}

	// Serve static files from the directory
	srvConfig := server.Config{
		Dir:     cfg.InputDir,
		Port:    cfg.Port,
		Quiet:   cfg.Quiet,
		Verbose: cfg.Verbose,
	}

	srv := server.New(srvConfig, w)
	return srv.Start()
}

// isSourceDirectory checks if a directory is a source directory
// (contains .md files but no index.html)
func isSourceDirectory(dir string) bool {
	// Check if index.html exists
	indexHTML := filepath.Join(dir, "index.html")
	if _, err := os.Stat(indexHTML); err == nil {
		return false // Has index.html, it's generated output
	}

	// Check if any .md files exist
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			return true
		}
	}

	return false
}

func printServeUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Start development server for live preview")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Usage:")
	_, _ = fmt.Fprintln(w, "  volcano serve [flags] <input>")
	_, _ = fmt.Fprintln(w, "  volcano server [flags] <input>   (alias)")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Server:")
	_, _ = fmt.Fprintln(w, "  -p, --port <port>    Server port (default: 1776)")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Site Configuration:")
	_, _ = fmt.Fprintln(w, "  --title <title>      Site title (default: My Site)")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Appearance:")
	_, _ = fmt.Fprintln(w, "  --theme <name>       Theme: docs, blog, vanilla (default: docs)")
	_, _ = fmt.Fprintln(w, "  --css <path>         Custom CSS file (overrides theme)")
	_, _ = fmt.Fprintln(w, "  --accent-color <hex> Custom accent color (e.g., '#ff6600')")
	_, _ = fmt.Fprintln(w, "  --favicon <path>     Favicon file (ico, png, svg, gif)")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Navigation:")
	_, _ = fmt.Fprintln(w, "  --top-nav            Show root files in top navigation bar")
	_, _ = fmt.Fprintln(w, "  --breadcrumbs        Show breadcrumb trail (default: true)")
	_, _ = fmt.Fprintln(w, "  --page-nav           Show previous/next page links")
	_, _ = fmt.Fprintln(w, "  --instant-nav        Enable hover prefetching for faster navigation")
	_, _ = fmt.Fprintln(w, "  --view-transitions   Enable browser view transitions API")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Logging:")
	_, _ = fmt.Fprintln(w, "  -q, --quiet          Suppress non-error output")
	_, _ = fmt.Fprintln(w, "  --verbose            Show detailed server logs")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Other:")
	_, _ = fmt.Fprintln(w, "  -h, --help           Show this help message")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Examples:")
	_, _ = fmt.Fprintln(w, "  volcano serve ./docs -p 8080")
	_, _ = fmt.Fprintln(w, "  volcano serve --theme=blog --instant-nav ./posts")
}

// serveValueFlags is the set of flags that take values for the serve command
var serveValueFlags = map[string]bool{
	"p": true, "port": true,
	"title": true, "theme": true, "css": true, "accent-color": true, "favicon": true,
}

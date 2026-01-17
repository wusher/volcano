package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wusher/volcano/internal/config"
	"github.com/wusher/volcano/internal/output"
	"github.com/wusher/volcano/internal/server"
	"github.com/wusher/volcano/internal/styles"
)

// ServeCommand handles the serve subcommand for starting the development server
func ServeCommand(args []string, stdout, stderr io.Writer) error {
	cfg := DefaultConfig()
	errLogger := output.NewLogger(stderr, output.IsStderrTTY(), false, false)
	tracker := newConfigTracker()

	// Initialize tracker with defaults
	initServeDefaultTracker(tracker, cfg)

	// Pre-scan args to find input directory and config file path
	inputDir, configPath := prescanServeArgs(args)

	// Load config file if specified or discovered
	var fileCfg *config.FileConfig
	if inputDir != "" {
		if err := validateInputDir(inputDir); err == nil {
			var cfgPath string
			var err error
			fileCfg, cfgPath, err = config.LoadOrDiscover(configPath, inputDir)
			if err != nil {
				errLogger.Error("%v", err)
				return err
			}
			if fileCfg != nil {
				applyServeFileConfig(cfg, fileCfg, tracker)
				cfg.configFilePath = cfgPath
			}
		}
	}

	// Store pre-CLI values for override detection
	preCLIValues := copyServeConfigValues(cfg)

	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// Define serve-specific flags
	var showHelp bool
	var configFlag string
	var viewTransitionsFlag bool // Deprecated flag, kept for backwards compatibility

	fs.IntVar(&cfg.Port, "p", cfg.Port, "Server port")
	fs.IntVar(&cfg.Port, "port", cfg.Port, "Server port")
	fs.StringVar(&cfg.Title, "title", cfg.Title, "Site title")
	fs.StringVar(&cfg.SiteURL, "url", cfg.SiteURL, "Site base URL")
	fs.StringVar(&cfg.Author, "author", cfg.Author, "Site author")
	fs.StringVar(&cfg.Theme, "theme", cfg.Theme, "Theme name (docs, blog, vanilla)")
	fs.StringVar(&cfg.CSSPath, "css", cfg.CSSPath, "Path to custom CSS file")
	fs.StringVar(&cfg.AccentColor, "accent-color", cfg.AccentColor, "Custom accent color (hex format, e.g., '#ff6600')")
	fs.StringVar(&cfg.FaviconPath, "favicon", cfg.FaviconPath, "Path to favicon file")
	fs.BoolVar(&cfg.TopNav, "top-nav", cfg.TopNav, "Display root files in top navigation bar")
	fs.BoolVar(&cfg.ShowPageNav, "page-nav", cfg.ShowPageNav, "Show previous/next page navigation")
	fs.BoolVar(&cfg.ShowBreadcrumbs, "breadcrumbs", cfg.ShowBreadcrumbs, "Show breadcrumb navigation")
	fs.BoolVar(&cfg.InstantNav, "instant-nav", cfg.InstantNav, "Enable instant navigation with hover prefetching")
	fs.BoolVar(&viewTransitionsFlag, "view-transitions", false, "Deprecated: view transitions are now enabled by default")
	fs.BoolVar(&cfg.Quiet, "q", cfg.Quiet, "Suppress non-error output")
	fs.BoolVar(&cfg.Quiet, "quiet", cfg.Quiet, "Suppress non-error output")
	fs.BoolVar(&cfg.Verbose, "verbose", cfg.Verbose, "Enable debug output")
	fs.BoolVar(&cfg.PWA, "pwa", cfg.PWA, "Enable PWA manifest and service worker for offline support")
	fs.BoolVar(&cfg.Search, "search", cfg.Search, "Enable site search with Cmd+K command palette")
	fs.StringVar(&configFlag, "config", "", "Path to config file (default: volcano.json in input directory)")
	fs.StringVar(&configFlag, "c", "", "Path to config file (default: volcano.json in input directory)")
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

	// Check if deprecated --view-transitions flag was used
	if viewTransitionsFlag {
		errLogger.Warning("--view-transitions is deprecated: view transitions are now enabled by default")
	}

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

	// Detect CLI overrides and update tracker
	detectServeCLIOverrides(cfg, preCLIValues, tracker)

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

	// Create logger for printing config
	stdLogger := output.NewLogger(stdout, output.IsStdoutTTY(), cfg.Quiet, cfg.Verbose)

	// Print config file info if loaded
	if cfg.configFilePath != "" {
		stdLogger.Println("Using config: %s", cfg.configFilePath)
	}

	// Print CLI override messages
	printServeCLIOverrides(stdLogger, tracker, fileCfg != nil)

	// Print configuration values
	if !cfg.Quiet {
		printServeConfig(stdLogger, cfg)
	}

	if err := Serve(cfg, stdout); err != nil {
		errLogger.Error("%v", err)
		return err
	}
	return nil
}

// initServeDefaultTracker initializes tracker with default values for serve command
func initServeDefaultTracker(tracker *configTracker, cfg *Config) {
	tracker.set("port", cfg.Port, sourceDefault)
	tracker.set("title", cfg.Title, sourceDefault)
	tracker.set("url", cfg.SiteURL, sourceDefault)
	tracker.set("author", cfg.Author, sourceDefault)
	tracker.set("theme", cfg.Theme, sourceDefault)
	tracker.set("css", cfg.CSSPath, sourceDefault)
	tracker.set("accentColor", cfg.AccentColor, sourceDefault)
	tracker.set("favicon", cfg.FaviconPath, sourceDefault)
	tracker.set("topNav", cfg.TopNav, sourceDefault)
	tracker.set("breadcrumbs", cfg.ShowBreadcrumbs, sourceDefault)
	tracker.set("pageNav", cfg.ShowPageNav, sourceDefault)
	tracker.set("instantNav", cfg.InstantNav, sourceDefault)
	tracker.set("pwa", cfg.PWA, sourceDefault)
	tracker.set("search", cfg.Search, sourceDefault)
}

// copyServeConfigValues creates a copy of config values for override detection
func copyServeConfigValues(cfg *Config) map[string]interface{} {
	return map[string]interface{}{
		"port":        cfg.Port,
		"title":       cfg.Title,
		"url":         cfg.SiteURL,
		"author":      cfg.Author,
		"theme":       cfg.Theme,
		"css":         cfg.CSSPath,
		"accentColor": cfg.AccentColor,
		"favicon":     cfg.FaviconPath,
		"topNav":      cfg.TopNav,
		"breadcrumbs": cfg.ShowBreadcrumbs,
		"pageNav":     cfg.ShowPageNav,
		"instantNav":  cfg.InstantNav,
		"pwa":         cfg.PWA,
		"search":      cfg.Search,
	}
}

// detectServeCLIOverrides detects which values were changed by CLI flags
func detectServeCLIOverrides(cfg *Config, preCLI map[string]interface{}, tracker *configTracker) {
	checkOverride := func(name string, oldVal, newVal interface{}) {
		if oldVal != newVal {
			tracker.set(name, newVal, sourceCLI)
		}
	}

	checkOverride("port", preCLI["port"], cfg.Port)
	checkOverride("title", preCLI["title"], cfg.Title)
	checkOverride("url", preCLI["url"], cfg.SiteURL)
	checkOverride("author", preCLI["author"], cfg.Author)
	checkOverride("theme", preCLI["theme"], cfg.Theme)
	checkOverride("css", preCLI["css"], cfg.CSSPath)
	checkOverride("accentColor", preCLI["accentColor"], cfg.AccentColor)
	checkOverride("favicon", preCLI["favicon"], cfg.FaviconPath)
	checkOverride("topNav", preCLI["topNav"], cfg.TopNav)
	checkOverride("breadcrumbs", preCLI["breadcrumbs"], cfg.ShowBreadcrumbs)
	checkOverride("pageNav", preCLI["pageNav"], cfg.ShowPageNav)
	checkOverride("instantNav", preCLI["instantNav"], cfg.InstantNav)
	checkOverride("pwa", preCLI["pwa"], cfg.PWA)
	checkOverride("search", preCLI["search"], cfg.Search)
}

// printServeCLIOverrides prints messages for CLI flags that override config file values
func printServeCLIOverrides(logger *output.Logger, tracker *configTracker, hasConfigFile bool) {
	if !hasConfigFile {
		return
	}

	// Map of option names to their CLI flag names
	flagNames := map[string]string{
		"port":        "--port",
		"title":       "--title",
		"url":         "--url",
		"author":      "--author",
		"theme":       "--theme",
		"css":         "--css",
		"accentColor": "--accent-color",
		"favicon":     "--favicon",
		"topNav":      "--top-nav",
		"breadcrumbs": "--breadcrumbs",
		"pageNav":     "--page-nav",
		"instantNav":  "--instant-nav",
		"pwa":         "--pwa",
		"search":      "--search",
	}

	for name, flagName := range flagNames {
		if tracker.getSource(name) == sourceCLI {
			logger.Println("CLI %s overrides config file", flagName)
		}
	}
}

// printServeConfig prints the configuration values being used for serve
func printServeConfig(logger *output.Logger, cfg *Config) {
	logger.Println("Configuration:")
	logger.Println("  input:       %s", cfg.InputDir)
	logger.Println("  port:        %d", cfg.Port)
	logger.Println("  title:       %s", cfg.Title)

	if cfg.SiteURL != "" {
		logger.Println("  url:         %s", cfg.SiteURL)
	}
	if cfg.Author != "" {
		logger.Println("  author:      %s", cfg.Author)
	}
	if cfg.Theme != "" {
		logger.Println("  theme:       %s", cfg.Theme)
	}
	if cfg.CSSPath != "" {
		logger.Println("  css:         %s", cfg.CSSPath)
	}
	if cfg.AccentColor != "" {
		logger.Println("  accentColor: %s", cfg.AccentColor)
	}
	if cfg.FaviconPath != "" {
		logger.Println("  favicon:     %s", cfg.FaviconPath)
	}

	// Print feature flags that are enabled
	var features []string
	if cfg.TopNav {
		features = append(features, "topNav")
	}
	if cfg.ShowBreadcrumbs {
		features = append(features, "breadcrumbs")
	}
	if cfg.ShowPageNav {
		features = append(features, "pageNav")
	}
	if cfg.InstantNav {
		features = append(features, "instantNav")
	}
	if cfg.PWA {
		features = append(features, "pwa")
	}
	if cfg.Search {
		features = append(features, "search")
	}

	if len(features) > 0 {
		logger.Println("  features:    %s", strings.Join(features, ", "))
	}
}

// applyServeFileConfig applies values from a config file to the Config struct for serve.
func applyServeFileConfig(cfg *Config, fileCfg *config.FileConfig, tracker *configTracker) {
	// Port value
	if fileCfg.Port != nil {
		cfg.Port = *fileCfg.Port
		tracker.set("port", *fileCfg.Port, sourceFile)
	}

	// String values - only apply if not empty
	if fileCfg.Title != "" {
		cfg.Title = fileCfg.Title
		tracker.set("title", fileCfg.Title, sourceFile)
	}
	if fileCfg.URL != "" {
		cfg.SiteURL = fileCfg.URL
		tracker.set("url", fileCfg.URL, sourceFile)
	}
	if fileCfg.Author != "" {
		cfg.Author = fileCfg.Author
		tracker.set("author", fileCfg.Author, sourceFile)
	}
	if fileCfg.Theme != "" {
		cfg.Theme = fileCfg.Theme
		tracker.set("theme", fileCfg.Theme, sourceFile)
	}
	if fileCfg.CSS != "" {
		cfg.CSSPath = fileCfg.CSS
		tracker.set("css", fileCfg.CSS, sourceFile)
	}
	if fileCfg.AccentColor != "" {
		cfg.AccentColor = fileCfg.AccentColor
		tracker.set("accentColor", fileCfg.AccentColor, sourceFile)
	}
	if fileCfg.Favicon != "" {
		cfg.FaviconPath = fileCfg.Favicon
		tracker.set("favicon", fileCfg.Favicon, sourceFile)
	}

	// Boolean values - only apply if explicitly set (non-nil)
	if fileCfg.TopNav != nil {
		cfg.TopNav = *fileCfg.TopNav
		tracker.set("topNav", *fileCfg.TopNav, sourceFile)
	}
	if fileCfg.Breadcrumbs != nil {
		cfg.ShowBreadcrumbs = *fileCfg.Breadcrumbs
		tracker.set("breadcrumbs", *fileCfg.Breadcrumbs, sourceFile)
	}
	if fileCfg.PageNav != nil {
		cfg.ShowPageNav = *fileCfg.PageNav
		tracker.set("pageNav", *fileCfg.PageNav, sourceFile)
	}
	if fileCfg.InstantNav != nil {
		cfg.InstantNav = *fileCfg.InstantNav
		tracker.set("instantNav", *fileCfg.InstantNav, sourceFile)
	}
	if fileCfg.PWA != nil {
		cfg.PWA = *fileCfg.PWA
		tracker.set("pwa", *fileCfg.PWA, sourceFile)
	}
	if fileCfg.Search != nil {
		cfg.Search = *fileCfg.Search
		tracker.set("search", *fileCfg.Search, sourceFile)
	}
}

// prescanServeArgs extracts the input directory and config path from args
func prescanServeArgs(args []string) (inputDir, configPath string) {
	i := 0
	for i < len(args) {
		arg := args[i]

		// Check for --config or -c flag
		if arg == "--config" || arg == "-c" {
			if i+1 < len(args) {
				configPath = args[i+1]
				i += 2
				continue
			}
		} else if strings.HasPrefix(arg, "--config=") {
			configPath = strings.TrimPrefix(arg, "--config=")
			i++
			continue
		} else if strings.HasPrefix(arg, "-c=") {
			configPath = strings.TrimPrefix(arg, "-c=")
			i++
			continue
		}

		// Skip other flags and their values
		if strings.HasPrefix(arg, "-") {
			name := strings.TrimLeft(arg, "-")
			if strings.Contains(name, "=") {
				i++
				continue
			}
			if serveValueFlags[name] && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i += 2
				continue
			}
			i++
			continue
		}

		// First non-flag argument is the input directory
		if inputDir == "" {
			inputDir = arg
		}
		i++
	}

	return inputDir, configPath
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
			PWA:             cfg.PWA,
			Search:          cfg.Search,
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
	_, _ = fmt.Fprintln(w, "  --url <url>          Site base URL")
	_, _ = fmt.Fprintln(w, "  --author <name>      Site author")
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
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Logging:")
	_, _ = fmt.Fprintln(w, "  -q, --quiet          Suppress non-error output")
	_, _ = fmt.Fprintln(w, "  --verbose            Show detailed server logs")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "PWA:")
	_, _ = fmt.Fprintln(w, "  --pwa                Enable PWA manifest and service worker")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Search:")
	_, _ = fmt.Fprintln(w, "  --search             Enable site search with Cmd+K command palette")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Configuration:")
	_, _ = fmt.Fprintln(w, "  -c, --config <path>  Config file (default: volcano.json in input dir)")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Other:")
	_, _ = fmt.Fprintln(w, "  -h, --help           Show this help message")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Config File:")
	_, _ = fmt.Fprintln(w, "  Create with 'volcano init' or place volcano.json in your input directory.")
	_, _ = fmt.Fprintln(w, "  CLI flags override config file values.")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Examples:")
	_, _ = fmt.Fprintln(w, "  volcano serve ./docs -p 8080")
	_, _ = fmt.Fprintln(w, "  volcano serve --theme=blog --instant-nav ./posts")
	_, _ = fmt.Fprintln(w, "  volcano serve --config=./my-config.json ./docs")
}

// serveValueFlags is the set of flags that take values for the serve command
var serveValueFlags = map[string]bool{
	"p": true, "port": true,
	"title": true, "url": true, "author": true,
	"theme": true, "css": true, "accent-color": true, "favicon": true,
	"config": true, "c": true,
}

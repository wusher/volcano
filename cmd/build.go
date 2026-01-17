package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/wusher/volcano/internal/config"
	"github.com/wusher/volcano/internal/output"
	"github.com/wusher/volcano/internal/styles"
)

// configSource tracks where a config value came from
type configSource int

const (
	sourceDefault configSource = iota
	sourceFile
	sourceCLI
)

// configTracker tracks config values and their sources for override messages
type configTracker struct {
	values  map[string]interface{}
	sources map[string]configSource
}

func newConfigTracker() *configTracker {
	return &configTracker{
		values:  make(map[string]interface{}),
		sources: make(map[string]configSource),
	}
}

func (ct *configTracker) set(name string, value interface{}, source configSource) {
	ct.values[name] = value
	ct.sources[name] = source
}

func (ct *configTracker) getSource(name string) configSource {
	return ct.sources[name]
}

// Build handles the build subcommand for generating static sites
func Build(args []string, stdout, stderr io.Writer) error {
	cfg := DefaultConfig()
	errLogger := output.NewLogger(stderr, output.IsStderrTTY(), false, false)
	tracker := newConfigTracker()

	// Initialize tracker with defaults
	initDefaultTracker(tracker, cfg)

	// Pre-scan args to find input directory and config file path
	// This allows us to load config file before parsing flags
	inputDir, configPath := prescanArgs(args)

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
				applyFileConfig(cfg, fileCfg, tracker)
				cfg.configFilePath = cfgPath
			}
		}
	}

	// Store pre-CLI values for override detection
	preCLIValues := copyConfigValues(cfg)

	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// Define build-specific flags
	var showHelp bool
	var configFlag string
	var viewTransitionsFlag bool // Deprecated flag, kept for backwards compatibility

	fs.StringVar(&cfg.OutputDir, "o", cfg.OutputDir, "Output directory for generated site")
	fs.StringVar(&cfg.OutputDir, "output", cfg.OutputDir, "Output directory for generated site")
	fs.StringVar(&cfg.Title, "title", cfg.Title, "Site title")
	fs.StringVar(&cfg.SiteURL, "url", cfg.SiteURL, "Site base URL for SEO")
	fs.StringVar(&cfg.Author, "author", cfg.Author, "Site author")
	fs.StringVar(&cfg.OGImage, "og-image", cfg.OGImage, "Default Open Graph image URL")
	fs.StringVar(&cfg.FaviconPath, "favicon", cfg.FaviconPath, "Path to favicon file")
	fs.BoolVar(&cfg.TopNav, "top-nav", cfg.TopNav, "Display root files in top navigation bar")
	fs.BoolVar(&cfg.ShowPageNav, "page-nav", cfg.ShowPageNav, "Show previous/next page navigation")
	fs.BoolVar(&cfg.ShowBreadcrumbs, "breadcrumbs", cfg.ShowBreadcrumbs, "Show breadcrumb navigation")
	fs.StringVar(&cfg.Theme, "theme", cfg.Theme, "Theme name (docs, blog, vanilla)")
	fs.StringVar(&cfg.CSSPath, "css", cfg.CSSPath, "Path to custom CSS file")
	fs.StringVar(&cfg.AccentColor, "accent-color", cfg.AccentColor, "Custom accent color (hex format, e.g., '#ff6600')")
	fs.BoolVar(&cfg.InstantNav, "instant-nav", cfg.InstantNav, "Enable instant navigation with hover prefetching")
	fs.BoolVar(&cfg.InlineAssets, "inline-assets", cfg.InlineAssets, "Embed CSS/JS inline instead of external files")
	fs.BoolVar(&cfg.PWA, "pwa", cfg.PWA, "Enable PWA manifest and service worker for offline support")
	fs.BoolVar(&cfg.Search, "search", cfg.Search, "Enable site search with Cmd+K command palette")
	fs.BoolVar(&cfg.AllowBrokenLinks, "allow-broken-links", cfg.AllowBrokenLinks, "Don't fail build on broken internal links")
	fs.BoolVar(&viewTransitionsFlag, "view-transitions", false, "Deprecated: view transitions are now enabled by default")
	fs.BoolVar(&cfg.Quiet, "q", cfg.Quiet, "Suppress non-error output")
	fs.BoolVar(&cfg.Quiet, "quiet", cfg.Quiet, "Suppress non-error output")
	fs.BoolVar(&cfg.Verbose, "verbose", cfg.Verbose, "Enable debug output")
	fs.StringVar(&configFlag, "config", "", "Path to config file (default: volcano.json in input directory)")
	fs.StringVar(&configFlag, "c", "", "Path to config file (default: volcano.json in input directory)")
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

	// Check for --url requirement
	if cfg.SiteURL == "" {
		errLogger.Error("--url is required for build (set site base URL for SEO)")
		_, _ = fmt.Fprintln(stderr, "")
		_, _ = fmt.Fprintln(stderr, "Example: volcano build ./docs --url=\"https://example.com\"")
		_, _ = fmt.Fprintln(stderr, "")
		_, _ = fmt.Fprintln(stderr, "Or set it in volcano.json:")
		_, _ = fmt.Fprintln(stderr, "  { \"url\": \"https://example.com\" }")
		return fmt.Errorf("--url is required")
	}

	// Detect CLI overrides and update tracker
	detectCLIOverrides(cfg, preCLIValues, tracker)

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

	// Set colored output based on stdout TTY detection
	cfg.Colored = output.IsStdoutTTY()

	// Create logger for printing config
	stdLogger := output.NewLogger(stdout, output.IsStdoutTTY(), cfg.Quiet, cfg.Verbose)

	// Print config file info if loaded
	if cfg.configFilePath != "" {
		stdLogger.Println("Using config: %s", cfg.configFilePath)
	}

	// Print CLI override messages
	printCLIOverrides(stdLogger, tracker, fileCfg != nil)

	// Print configuration values
	if !cfg.Quiet {
		printBuildConfig(stdLogger, cfg)
	}

	if err := Generate(cfg, stdout); err != nil {
		errLogger.Error("%v", err)
		return err
	}
	return nil
}

// initDefaultTracker initializes tracker with default values
func initDefaultTracker(tracker *configTracker, cfg *Config) {
	tracker.set("output", cfg.OutputDir, sourceDefault)
	tracker.set("title", cfg.Title, sourceDefault)
	tracker.set("url", cfg.SiteURL, sourceDefault)
	tracker.set("author", cfg.Author, sourceDefault)
	tracker.set("theme", cfg.Theme, sourceDefault)
	tracker.set("css", cfg.CSSPath, sourceDefault)
	tracker.set("accentColor", cfg.AccentColor, sourceDefault)
	tracker.set("favicon", cfg.FaviconPath, sourceDefault)
	tracker.set("ogImage", cfg.OGImage, sourceDefault)
	tracker.set("topNav", cfg.TopNav, sourceDefault)
	tracker.set("breadcrumbs", cfg.ShowBreadcrumbs, sourceDefault)
	tracker.set("pageNav", cfg.ShowPageNav, sourceDefault)
	tracker.set("instantNav", cfg.InstantNav, sourceDefault)
	tracker.set("inlineAssets", cfg.InlineAssets, sourceDefault)
	tracker.set("pwa", cfg.PWA, sourceDefault)
	tracker.set("search", cfg.Search, sourceDefault)
	tracker.set("allowBrokenLinks", cfg.AllowBrokenLinks, sourceDefault)
}

// copyConfigValues creates a copy of config values for override detection
func copyConfigValues(cfg *Config) map[string]interface{} {
	return map[string]interface{}{
		"output":           cfg.OutputDir,
		"title":            cfg.Title,
		"url":              cfg.SiteURL,
		"author":           cfg.Author,
		"theme":            cfg.Theme,
		"css":              cfg.CSSPath,
		"accentColor":      cfg.AccentColor,
		"favicon":          cfg.FaviconPath,
		"ogImage":          cfg.OGImage,
		"topNav":           cfg.TopNav,
		"breadcrumbs":      cfg.ShowBreadcrumbs,
		"pageNav":          cfg.ShowPageNav,
		"instantNav":       cfg.InstantNav,
		"inlineAssets":     cfg.InlineAssets,
		"pwa":              cfg.PWA,
		"search":           cfg.Search,
		"allowBrokenLinks": cfg.AllowBrokenLinks,
	}
}

// detectCLIOverrides detects which values were changed by CLI flags
func detectCLIOverrides(cfg *Config, preCLI map[string]interface{}, tracker *configTracker) {
	checkOverride := func(name string, oldVal, newVal interface{}) {
		if oldVal != newVal {
			tracker.set(name, newVal, sourceCLI)
		}
	}

	checkOverride("output", preCLI["output"], cfg.OutputDir)
	checkOverride("title", preCLI["title"], cfg.Title)
	checkOverride("url", preCLI["url"], cfg.SiteURL)
	checkOverride("author", preCLI["author"], cfg.Author)
	checkOverride("theme", preCLI["theme"], cfg.Theme)
	checkOverride("css", preCLI["css"], cfg.CSSPath)
	checkOverride("accentColor", preCLI["accentColor"], cfg.AccentColor)
	checkOverride("favicon", preCLI["favicon"], cfg.FaviconPath)
	checkOverride("ogImage", preCLI["ogImage"], cfg.OGImage)
	checkOverride("topNav", preCLI["topNav"], cfg.TopNav)
	checkOverride("breadcrumbs", preCLI["breadcrumbs"], cfg.ShowBreadcrumbs)
	checkOverride("pageNav", preCLI["pageNav"], cfg.ShowPageNav)
	checkOverride("instantNav", preCLI["instantNav"], cfg.InstantNav)
	checkOverride("inlineAssets", preCLI["inlineAssets"], cfg.InlineAssets)
	checkOverride("pwa", preCLI["pwa"], cfg.PWA)
	checkOverride("search", preCLI["search"], cfg.Search)
	checkOverride("allowBrokenLinks", preCLI["allowBrokenLinks"], cfg.AllowBrokenLinks)
}

// printCLIOverrides prints messages for CLI flags that override config file values
func printCLIOverrides(logger *output.Logger, tracker *configTracker, hasConfigFile bool) {
	if !hasConfigFile {
		return
	}

	// Map of option names to their CLI flag names
	flagNames := map[string]string{
		"output":           "--output",
		"title":            "--title",
		"url":              "--url",
		"author":           "--author",
		"theme":            "--theme",
		"css":              "--css",
		"accentColor":      "--accent-color",
		"favicon":          "--favicon",
		"ogImage":          "--og-image",
		"topNav":           "--top-nav",
		"breadcrumbs":      "--breadcrumbs",
		"pageNav":          "--page-nav",
		"instantNav":       "--instant-nav",
		"inlineAssets":     "--inline-assets",
		"pwa":              "--pwa",
		"search":           "--search",
		"allowBrokenLinks": "--allow-broken-links",
	}

	for name, flagName := range flagNames {
		if tracker.getSource(name) == sourceCLI {
			logger.Println("CLI %s overrides config file", flagName)
		}
	}
}

// printBuildConfig prints the configuration values being used
func printBuildConfig(logger *output.Logger, cfg *Config) {
	logger.Println("Configuration:")
	logger.Println("  input:       %s", cfg.InputDir)
	logger.Println("  output:      %s", cfg.OutputDir)
	logger.Println("  url:         %s", cfg.SiteURL)
	logger.Println("  title:       %s", cfg.Title)

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
	if cfg.OGImage != "" {
		logger.Println("  ogImage:     %s", cfg.OGImage)
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
	if cfg.InlineAssets {
		features = append(features, "inlineAssets")
	}
	if cfg.PWA {
		features = append(features, "pwa")
	}
	if cfg.Search {
		features = append(features, "search")
	}
	if cfg.AllowBrokenLinks {
		features = append(features, "allowBrokenLinks")
	}

	if len(features) > 0 {
		logger.Println("  features:    %s", strings.Join(features, ", "))
	}
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
	_, _ = fmt.Fprintln(w, "Site Configuration (required):")
	_, _ = fmt.Fprintln(w, "  --url <url>          Base URL for canonical links and SEO (REQUIRED)")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Site Configuration (optional):")
	_, _ = fmt.Fprintln(w, "  --title <title>      Site title (default: My Site)")
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
	_, _ = fmt.Fprintln(w, "  --search             Enable site search with Cmd+K command palette")
	_, _ = fmt.Fprintln(w, "  --allow-broken-links Don't fail build on broken internal links")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "SEO:")
	_, _ = fmt.Fprintln(w, "  --og-image <path>    Default Open Graph image")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Logging:")
	_, _ = fmt.Fprintln(w, "  -q, --quiet          Suppress non-error output")
	_, _ = fmt.Fprintln(w, "  --verbose            Show detailed build information")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Configuration:")
	_, _ = fmt.Fprintln(w, "  -c, --config <path>  Config file (default: volcano.json in input dir)")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Other:")
	_, _ = fmt.Fprintln(w, "  -h, --help           Show this help message")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Config File:")
	_, _ = fmt.Fprintln(w, "  Create with 'volcano init' or place volcano.json in your input directory.")
	_, _ = fmt.Fprintln(w, "  CLI flags override config file values. Example volcano.json:")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "    {")
	_, _ = fmt.Fprintln(w, "      \"url\": \"https://example.com\",")
	_, _ = fmt.Fprintln(w, "      \"title\": \"My Docs\",")
	_, _ = fmt.Fprintln(w, "      \"output\": \"./public\",")
	_, _ = fmt.Fprintln(w, "      \"theme\": \"docs\",")
	_, _ = fmt.Fprintln(w, "      \"pwa\": true")
	_, _ = fmt.Fprintln(w, "    }")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Examples:")
	_, _ = fmt.Fprintln(w, "  volcano build ./docs --url=\"https://docs.example.com\"")
	_, _ = fmt.Fprintln(w, "  volcano build ./docs -o ./public --url=\"https://docs.example.com\" --title=\"My Docs\"")
	_, _ = fmt.Fprintln(w, "  volcano build --theme=blog --accent-color='#ff6600' --url=\"https://blog.example.com\" ./posts")
	_, _ = fmt.Fprintln(w, "  volcano build --config=./my-config.json ./docs")
}

// buildValueFlags is the set of flags that take values for the build command
var buildValueFlags = map[string]bool{
	"o": true, "output": true,
	"title": true, "url": true, "author": true,
	"og-image": true, "favicon": true,
	"theme": true, "css": true, "accent-color": true,
	"config": true, "c": true,
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

// prescanArgs extracts the input directory and config path from args
// before the main flag parsing. This allows loading the config file early.
func prescanArgs(args []string) (inputDir, configPath string) {
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
			// Skip value flags
			name := strings.TrimLeft(arg, "-")
			// Handle --flag=value format
			if strings.Contains(name, "=") {
				i++
				continue
			}
			// Check if this flag takes a value
			if buildValueFlags[name] && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
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

// applyFileConfig applies values from a config file to the Config struct.
// Only non-empty/non-nil values from the file are applied.
func applyFileConfig(cfg *Config, fileCfg *config.FileConfig, tracker *configTracker) {
	// String values - only apply if not empty
	if fileCfg.Output != "" {
		cfg.OutputDir = fileCfg.Output
		tracker.set("output", fileCfg.Output, sourceFile)
	}
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
	if fileCfg.OGImage != "" {
		cfg.OGImage = fileCfg.OGImage
		tracker.set("ogImage", fileCfg.OGImage, sourceFile)
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
	if fileCfg.InlineAssets != nil {
		cfg.InlineAssets = *fileCfg.InlineAssets
		tracker.set("inlineAssets", *fileCfg.InlineAssets, sourceFile)
	}
	if fileCfg.PWA != nil {
		cfg.PWA = *fileCfg.PWA
		tracker.set("pwa", *fileCfg.PWA, sourceFile)
	}
	if fileCfg.Search != nil {
		cfg.Search = *fileCfg.Search
		tracker.set("search", *fileCfg.Search, sourceFile)
	}
	if fileCfg.AllowBrokenLinks != nil {
		cfg.AllowBrokenLinks = *fileCfg.AllowBrokenLinks
		tracker.set("allowBrokenLinks", *fileCfg.AllowBrokenLinks, sourceFile)
	}
}

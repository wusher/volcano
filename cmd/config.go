// Package cmd provides the command implementations for the volcano CLI.
package cmd

// Config holds all configuration options for the volcano CLI
type Config struct {
	InputDir        string // Input directory containing markdown files
	OutputDir       string // Output directory for generated HTML
	ServeMode       bool   // Whether to run in serve mode
	Port            int    // Port for the HTTP server
	Title           string // Site title
	Quiet           bool   // Suppress non-error output
	Verbose         bool   // Enable debug output
	Colored         bool   // Enable colored output (auto-detected from TTY)
	SiteURL         string // Base URL for canonical links and SEO
	Author          string // Site author
	OGImage         string // Default Open Graph image
	FaviconPath     string // Path to favicon file
	ShowLastMod     bool   // Show last modified date
	TopNav          bool   // Display root files in top navigation bar
	ShowPageNav     bool   // Show previous/next page navigation
	ShowBreadcrumbs bool   // Show breadcrumb navigation
	Theme           string // Theme name (docs, blog, vanilla)
	CSSPath         string // Path to custom CSS file
	AccentColor     string // Custom accent color in hex format (e.g., "#ff6600")
	InstantNav      bool   // Enable instant navigation with hover prefetching
	ViewTransitions bool   // Enable browser view transitions API
	InlineAssets    bool   // Embed CSS/JS inline instead of external files
	PWA             bool   // Enable PWA manifest and service worker generation

	// Internal fields (not settable via CLI)
	configFilePath string // Path to loaded config file (for verbose logging)
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		OutputDir:       "./output",
		Port:            1776,
		Title:           "My Site",
		ShowBreadcrumbs: true, // Show breadcrumbs by default
		ViewTransitions: true, // View transitions enabled by default
	}
}

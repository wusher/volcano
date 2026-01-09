// Package cmd provides the command implementations for the volcano CLI.
package cmd

// Config holds all configuration options for the volcano CLI
type Config struct {
	InputDir    string // Input directory containing markdown files
	OutputDir   string // Output directory for generated HTML
	ServeMode   bool   // Whether to run in serve mode
	Port        int    // Port for the HTTP server
	Title       string // Site title
	Quiet       bool   // Suppress non-error output
	Verbose     bool   // Enable debug output
	Colored     bool   // Enable colored output (auto-detected from TTY)
	SiteURL     string // Base URL for canonical links and SEO
	Author      string // Site author
	OGImage     string // Default Open Graph image
	FaviconPath string // Path to favicon file
	ShowLastMod bool   // Show last modified date
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		OutputDir: "./output",
		Port:      1776,
		Title:     "My Site",
	}
}

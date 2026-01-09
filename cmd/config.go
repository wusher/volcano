// Package cmd provides the command implementations for the volcano CLI.
package cmd

// Config holds all configuration options for the volcano CLI
type Config struct {
	InputDir  string // Input directory containing markdown files
	OutputDir string // Output directory for generated HTML
	ServeMode bool   // Whether to run in serve mode
	Port      int    // Port for the HTTP server
	Title     string // Site title
	Quiet     bool   // Suppress non-error output
	Verbose   bool   // Enable debug output
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		OutputDir: "./output",
		Port:      1776,
		Title:     "My Site",
	}
}

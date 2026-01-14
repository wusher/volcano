// Package config provides configuration file loading for Volcano.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ConfigFileName is the default config file name to look for.
const ConfigFileName = "volcano.json"

// FileConfig represents the configuration options that can be set in a config file.
// All fields use pointers so we can distinguish between "not set" and "set to zero/false".
type FileConfig struct {
	// Output settings
	Output string `json:"output,omitempty"` // Output directory

	// Site configuration
	Title  string `json:"title,omitempty"`  // Site title
	URL    string `json:"url,omitempty"`    // Base URL for SEO
	Author string `json:"author,omitempty"` // Site author

	// Appearance
	Theme       string `json:"theme,omitempty"`       // Theme name (docs, blog, vanilla)
	CSS         string `json:"css,omitempty"`         // Path to custom CSS file
	AccentColor string `json:"accentColor,omitempty"` // Custom accent color
	Favicon     string `json:"favicon,omitempty"`     // Path to favicon file

	// Navigation
	TopNav       *bool `json:"topNav,omitempty"`       // Show top navigation bar
	Breadcrumbs  *bool `json:"breadcrumbs,omitempty"`  // Show breadcrumbs
	PageNav      *bool `json:"pageNav,omitempty"`      // Show prev/next navigation
	InstantNav   *bool `json:"instantNav,omitempty"`   // Enable instant navigation
	InlineAssets *bool `json:"inlineAssets,omitempty"` // Embed CSS/JS inline
	PWA          *bool `json:"pwa,omitempty"`          // Enable PWA support

	// Content
	LastModified *bool `json:"lastModified,omitempty"` // Show last modified date

	// SEO
	OGImage string `json:"ogImage,omitempty"` // Default Open Graph image
}

// Load reads a config file from the given path and returns the parsed configuration.
func Load(path string) (*FileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg FileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Discover looks for a config file in the given directory.
// Returns the path to the config file if found, or empty string if not found.
func Discover(dir string) string {
	configPath := filepath.Join(dir, ConfigFileName)
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}
	return ""
}

// LoadOrDiscover loads a config file from the explicit path if provided,
// or discovers one in the input directory. Returns nil if no config is found.
func LoadOrDiscover(explicitPath, inputDir string) (*FileConfig, string, error) {
	var configPath string

	if explicitPath != "" {
		// Explicit path provided - it must exist
		if _, err := os.Stat(explicitPath); err != nil {
			if os.IsNotExist(err) {
				return nil, "", fmt.Errorf("config file not found: %s", explicitPath)
			}
			return nil, "", fmt.Errorf("cannot access config file: %w", err)
		}
		configPath = explicitPath
	} else {
		// Try to discover config file in input directory
		configPath = Discover(inputDir)
		if configPath == "" {
			// No config file found - this is fine
			return nil, "", nil
		}
	}

	cfg, err := Load(configPath)
	if err != nil {
		return nil, configPath, err
	}

	return cfg, configPath, nil
}

// BoolPtr is a helper to create a pointer to a bool value.
func BoolPtr(b bool) *bool {
	return &b
}

// GetBool returns the value of a bool pointer, or the default if nil.
func GetBool(ptr *bool, defaultVal bool) bool {
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

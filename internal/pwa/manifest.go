// Package pwa provides Progressive Web App support for generated sites.
package pwa

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// DefaultThemeColor is the default theme color when no accent color is provided.
const DefaultThemeColor = "#3b82f6"

// ManifestConfig holds configuration for manifest generation.
type ManifestConfig struct {
	SiteTitle   string // Full site title
	Description string // Site description (from first page or empty)
	ThemeColor  string // From --accent-color or default #3b82f6
	BaseURL     string // Base URL path prefix
	HasIcons    bool   // Whether PWA icons were generated
}

// Manifest represents the web app manifest.
type Manifest struct {
	Name            string         `json:"name"`
	ShortName       string         `json:"short_name"`
	Description     string         `json:"description,omitempty"`
	StartURL        string         `json:"start_url"`
	Scope           string         `json:"scope"`
	Display         string         `json:"display"`
	BackgroundColor string         `json:"background_color"`
	ThemeColor      string         `json:"theme_color"`
	Icons           []ManifestIcon `json:"icons,omitempty"`
}

// ManifestIcon represents an icon entry in the manifest.
type ManifestIcon struct {
	Src   string `json:"src"`
	Sizes string `json:"sizes"`
	Type  string `json:"type"`
}

// GenerateManifest creates and writes the manifest.json file.
func GenerateManifest(outputDir string, config ManifestConfig) error {
	// Build start URL and scope
	startURL := "/"
	scope := "/"
	if config.BaseURL != "" {
		startURL = config.BaseURL + "/"
		scope = config.BaseURL + "/"
	}

	// Use provided theme color or default
	themeColor := config.ThemeColor
	if themeColor == "" {
		themeColor = DefaultThemeColor
	}

	// Truncate short name to 12 characters
	shortName := config.SiteTitle
	if len(shortName) > 12 {
		shortName = shortName[:12]
	}

	manifest := Manifest{
		Name:            config.SiteTitle,
		ShortName:       shortName,
		Description:     config.Description,
		StartURL:        startURL,
		Scope:           scope,
		Display:         "standalone",
		BackgroundColor: "#ffffff",
		ThemeColor:      themeColor,
	}

	// Add icons if generated
	if config.HasIcons {
		iconBase := "/"
		if config.BaseURL != "" {
			iconBase = config.BaseURL + "/"
		}
		manifest.Icons = []ManifestIcon{
			{
				Src:   iconBase + "icon-192.png",
				Sizes: "192x192",
				Type:  "image/png",
			},
			{
				Src:   iconBase + "icon-512.png",
				Sizes: "512x512",
				Type:  "image/png",
			},
		}
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filepath.Join(outputDir, "manifest.json"), data, 0644)
}

// GetManifestLinkTag returns the HTML link tag for the manifest.
func GetManifestLinkTag(baseURL string) string {
	href := "/manifest.json"
	if baseURL != "" {
		href = baseURL + "/manifest.json"
	}
	return `<link rel="manifest" href="` + href + `">`
}

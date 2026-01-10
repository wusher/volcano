// Package styles provides embedded CSS styles for the generated site.
package styles

import (
	_ "embed"
	"fmt"
)

// Theme CSS files embedded from themes directory
//
//go:embed themes/docs.css
var DocsCSS string

//go:embed themes/blog.css
var BlogCSS string

//go:embed themes/vanilla.css
var VanillaCSS string

// CSS is kept for backward compatibility, points to docs theme
var CSS = DocsCSS

// ValidThemes lists all available theme names
var ValidThemes = []string{"docs", "blog", "vanilla"}

// GetCSS returns the CSS for the specified theme
// If theme is empty, returns the docs theme (default)
func GetCSS(theme string) string {
	switch theme {
	case "blog":
		return BlogCSS
	case "vanilla":
		return VanillaCSS
	case "docs", "":
		return DocsCSS
	default:
		return DocsCSS
	}
}

// ValidateTheme checks if the given theme name is valid
func ValidateTheme(theme string) error {
	if theme == "" {
		return nil // Empty means default (docs)
	}
	for _, valid := range ValidThemes {
		if theme == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid theme %q, valid themes are: docs, blog, vanilla", theme)
}

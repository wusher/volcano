// Package styles provides embedded CSS styles for the generated site.
package styles

import (
	_ "embed"
	"fmt"

	"github.com/wusher/volcano/internal/color"
)

// LayoutCSS is the embedded shared layout stylesheet.
//
//go:embed themes/layout.css
var LayoutCSS string

// DocsCSS is the embedded docs theme stylesheet.
//
//go:embed themes/docs.css
var DocsCSS string

// BlogCSS is the embedded blog theme stylesheet.
//
//go:embed themes/blog.css
var BlogCSS string

// VanillaCSS is the embedded vanilla theme stylesheet.
//
//go:embed themes/vanilla.css
var VanillaCSS string

// CSS is kept for backward compatibility, points to docs theme
var CSS = LayoutCSS + "\n" + DocsCSS

// ValidThemes lists all available theme names
var ValidThemes = []string{"docs", "blog", "vanilla"}

// GetCSS returns the CSS for the specified theme
// If theme is empty, returns the docs theme (default)
// CSS is combined: layout.css + theme.css
func GetCSS(theme string) string {
	var themeCSS string
	switch theme {
	case "blog":
		themeCSS = BlogCSS
	case "vanilla":
		themeCSS = VanillaCSS
	case "docs", "":
		themeCSS = DocsCSS
	default:
		themeCSS = DocsCSS
	}
	return LayoutCSS + "\n" + themeCSS
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

// CSSLoader provides CSS content loading functionality.
// This interface enables dependency injection for testing.
type CSSLoader interface {
	LoadCSS() (string, error)
}

// CSSConfig holds configuration for loading CSS
type CSSConfig struct {
	Theme       string // Theme name (docs, blog, vanilla)
	CSSPath     string // Path to custom CSS file (takes precedence over Theme)
	AccentColor string // Custom accent color in hex format (e.g., "#ff6600")
}

// cssLoader implements CSSLoader
type cssLoader struct {
	config   CSSConfig
	readFile func(string) ([]byte, error)
}

// NewCSSLoader creates a new CSSLoader with the given configuration
func NewCSSLoader(config CSSConfig, readFile func(string) ([]byte, error)) CSSLoader {
	return &cssLoader{
		config:   config,
		readFile: readFile,
	}
}

// LoadCSS returns minified CSS content from custom file or embedded theme
func (l *cssLoader) LoadCSS() (string, error) {
	var css string
	if l.config.CSSPath != "" {
		content, err := l.readFile(l.config.CSSPath)
		if err != nil {
			return "", err
		}
		css = string(content)
	} else {
		css = GetCSS(l.config.Theme)
	}

	// Append accent color CSS if configured
	if l.config.AccentColor != "" {
		accentCSS, err := color.GenerateAccentCSS(l.config.AccentColor)
		if err != nil {
			return "", fmt.Errorf("invalid accent color: %w", err)
		}
		css = css + "\n" + accentCSS
	}

	return MinifyCSS(css)
}

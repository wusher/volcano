// Package styles provides embedded CSS styles for the generated site.
package styles

import (
	_ "embed"
)

// CSS contains the embedded stylesheet for the generated site.
//
//go:embed styles.css
var CSS string

// GetCSS returns the embedded CSS styles
func GetCSS() string {
	return CSS
}

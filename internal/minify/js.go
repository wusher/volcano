// Package minify provides JavaScript minification utilities.
package minify

import (
	"bytes"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

// JS minifies JavaScript code.
// Returns minified JS on success, or original JS if minification fails.
func JS(input string) string {
	if input == "" {
		return ""
	}

	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)

	var buf bytes.Buffer
	if err := m.Minify("text/javascript", &buf, bytes.NewReader([]byte(input))); err != nil {
		// Return original input if minification fails
		return input
	}

	return buf.String()
}

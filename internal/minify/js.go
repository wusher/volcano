// Package minify provides JavaScript minification utilities.
package minify

import (
	"bytes"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
)

// jsMinifier is a cached minifier instance for JavaScript.
// Creating a minifier once and reusing it is more efficient than
// creating a new one on each call.
var jsMinifier *minify.M

func init() {
	jsMinifier = minify.New()
	jsMinifier.AddFunc("text/javascript", js.Minify)
}

// JS minifies JavaScript code.
// Returns minified JS on success, or original JS if minification fails.
func JS(input string) string {
	if input == "" {
		return ""
	}

	var buf bytes.Buffer
	if err := jsMinifier.Minify("text/javascript", &buf, bytes.NewReader([]byte(input))); err != nil {
		// Return original input if minification fails
		return input
	}

	return buf.String()
}

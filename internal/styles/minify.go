package styles

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

// MinifyCSS minifies CSS content using tdewolff/minify
func MinifyCSS(input string) (string, error) {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	return m.String("text/css", input)
}

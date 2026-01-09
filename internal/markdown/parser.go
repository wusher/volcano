// Package markdown provides markdown parsing and HTML rendering functionality.
package markdown

import (
	"bytes"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Parser handles markdown parsing and HTML rendering
type Parser struct {
	md goldmark.Markdown
}

// NewParser creates a new markdown parser with all features enabled
func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,            // GitHub Flavored Markdown: tables, strikethrough, autolinks, task lists
			extension.Typographer,    // Smart quotes and dashes
			extension.Footnote,       // Footnotes
			extension.DefinitionList, // Definition lists
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true), // Use CSS classes instead of inline styles
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // Automatically generate heading IDs
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // Convert soft line breaks to <br>
			html.WithXHTML(),     // Use XHTML-style self-closing tags
			html.WithUnsafe(),    // Allow raw HTML (needed for some content)
		),
	)

	return &Parser{md: md}
}

// Parse converts markdown content to HTML
func (p *Parser) Parse(source []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := p.md.Convert(source, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ParseString converts markdown string to HTML string
func (p *Parser) ParseString(source string) (string, error) {
	result, err := p.Parse([]byte(source))
	if err != nil {
		return "", err
	}
	return string(result), nil
}

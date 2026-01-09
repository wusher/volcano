package markdown

import (
	"regexp"
	"strings"
)

var (
	// h1Pattern matches a markdown H1 heading
	h1Pattern = regexp.MustCompile(`(?m)^#\s+(.+)$`)

	// htmlH1Pattern matches an HTML h1 tag
	htmlH1Pattern = regexp.MustCompile(`(?i)<h1[^>]*>([^<]+)</h1>`)
)

// ExtractTitle extracts the first H1 heading from markdown content
// Returns empty string if no H1 is found
func ExtractTitle(markdownContent []byte) string {
	matches := h1Pattern.FindSubmatch(markdownContent)
	if len(matches) >= 2 {
		return strings.TrimSpace(string(matches[1]))
	}
	return ""
}

// ExtractTitleFromHTML extracts the first H1 heading from HTML content
// Returns empty string if no H1 is found
func ExtractTitleFromHTML(htmlContent []byte) string {
	matches := htmlH1Pattern.FindSubmatch(htmlContent)
	if len(matches) >= 2 {
		return strings.TrimSpace(string(matches[1]))
	}
	return ""
}

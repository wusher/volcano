package search

import (
	"html"
	"regexp"
	"strings"
)

// headingRegex matches h2-h4 tags with id attribute.
var headingRegex = regexp.MustCompile(`(?i)<h([2-4])[^>]*\s+id="([^"]+)"[^>]*>(.*?)</h[2-4]>`)

// stripTagsRegex removes HTML tags from text.
var stripTagsRegex = regexp.MustCompile(`<[^>]*>`)

// ExtractHeadings extracts heading entries from HTML content.
func ExtractHeadings(htmlContent string) []HeadingEntry {
	matches := headingRegex.FindAllStringSubmatch(htmlContent, -1)
	if len(matches) == 0 {
		return nil
	}

	entries := make([]HeadingEntry, 0, len(matches))
	for _, match := range matches {
		level := int(match[1][0] - '0')
		anchor := match[2]
		text := stripTagsRegex.ReplaceAllString(match[3], "")
		text = strings.TrimSpace(html.UnescapeString(text))

		if text != "" {
			entries = append(entries, HeadingEntry{
				Text:   text,
				Anchor: anchor,
				Level:  level,
			})
		}
	}
	return entries
}

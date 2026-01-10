package markdown

import (
	"regexp"
	"strings"
)

// internalLinkRegex matches href attributes with internal links (starting with /)
var internalLinkRegex = regexp.MustCompile(`href="(/[^"]*)"`)

// BrokenLink represents a broken internal link
type BrokenLink struct {
	SourcePage string // The page containing the link
	LinkURL    string // The broken link URL
}

// ExtractInternalLinks extracts all internal links from HTML content.
// Internal links are those starting with "/" (not external URLs).
func ExtractInternalLinks(htmlContent string) []string {
	matches := internalLinkRegex.FindAllStringSubmatch(htmlContent, -1)
	links := make([]string, 0, len(matches))

	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) >= 2 {
			link := match[1]
			// Skip anchor-only links and already seen links
			if link == "" || link == "/" || seen[link] {
				continue
			}
			// Skip anchor links within the same page
			if strings.HasPrefix(link, "/#") {
				continue
			}
			seen[link] = true
			links = append(links, link)
		}
	}

	return links
}

// ValidateLinks checks if all internal links resolve to valid URLs in the site.
// validURLs should be a map of all valid URL paths in the site.
// Returns a list of broken links.
func ValidateLinks(htmlContent string, sourcePage string, validURLs map[string]bool) []BrokenLink {
	links := ExtractInternalLinks(htmlContent)
	var broken []BrokenLink

	for _, link := range links {
		// Normalize link for comparison
		normalized := normalizeLink(link)

		// Check if the link resolves
		if !validURLs[normalized] {
			// Also check without trailing slash
			withoutSlash := strings.TrimSuffix(normalized, "/")
			withSlash := normalized
			if !strings.HasSuffix(normalized, "/") {
				withSlash = normalized + "/"
			}

			if !validURLs[withoutSlash] && !validURLs[withSlash] {
				broken = append(broken, BrokenLink{
					SourcePage: sourcePage,
					LinkURL:    link,
				})
			}
		}
	}

	return broken
}

// normalizeLink normalizes a link URL for comparison
func normalizeLink(link string) string {
	// Remove anchor/fragment
	if idx := strings.Index(link, "#"); idx != -1 {
		link = link[:idx]
	}

	// Ensure trailing slash for directory-style URLs
	if link != "/" && !strings.HasSuffix(link, "/") && !strings.Contains(link, ".") {
		link = link + "/"
	}

	return link
}

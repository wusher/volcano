package markdown

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wusher/volcano/internal/tree"
)

// attachmentExtensions lists file extensions that are attachments (not pages)
// These links are not validated against the page URL list
var validationSkipExtensions = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
	".svg": true, ".bmp": true, ".ico": true, ".heic": true,
	".pdf": true,
	".mp3": true, ".mp4": true, ".wav": true, ".ogg": true, ".webm": true, ".mov": true,
	".zip": true, ".docx": true, ".xlsx": true, ".pptx": true,
}

// internalLinkRegex matches href attributes with internal links (starting with /)
var internalLinkRegex = regexp.MustCompile(`href="(/[^"]*)"`)

// BrokenLink represents a broken internal link with detailed context
type BrokenLink struct {
	SourcePage     string   // The page containing the link (URL path)
	SourceFile     string   // The source markdown file path
	LineNumber     int      // Line number in the markdown file (0 if unknown)
	LinkURL        string   // The broken link URL
	OriginalSyntax string   // The original markdown/wikilink syntax
	LinkText       string   // The display text of the link
	Suggestions    []string // Suggested similar valid URLs
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
// validURLs should be a map of all valid URL paths in the site (including base URL prefix if applicable).
// Returns a list of broken links with detailed context.
func ValidateLinks(htmlContent string, sourcePage string, validURLs map[string]bool) []BrokenLink {
	return ValidateLinksWithSource(htmlContent, sourcePage, "", "", validURLs)
}

// ValidateLinksWithSource checks if all internal links resolve to valid URLs in the site,
// including source file context for better error messages.
func ValidateLinksWithSource(htmlContent string, sourcePage string, sourceFile string, mdSource string, validURLs map[string]bool) []BrokenLink {
	links := ExtractInternalLinks(htmlContent)
	var broken []BrokenLink

	// Build a map of URL -> line info from markdown source
	linkInfo := extractLinkInfoFromMarkdown(mdSource)

	for _, link := range links {
		// Skip attachment links (images, PDFs, etc.) - they're not validated as pages
		ext := strings.ToLower(filepath.Ext(link))
		if validationSkipExtensions[ext] {
			continue
		}

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
				// Find suggestions for similar URLs
				suggestions := findSimilarURLs(link, validURLs)

				// Look up line info from markdown source
				info := linkInfo[link]

				broken = append(broken, BrokenLink{
					SourcePage:     sourcePage,
					SourceFile:     sourceFile,
					LineNumber:     info.LineNumber,
					LinkURL:        link,
					OriginalSyntax: info.OriginalSyntax,
					LinkText:       info.LinkText,
					Suggestions:    suggestions,
				})
			}
		}
	}

	return broken
}

// linkInfo holds information about a link extracted from markdown
type linkInfo struct {
	LineNumber     int
	OriginalSyntax string
	LinkText       string
}

// extractLinkInfoFromMarkdown parses markdown source to find link information
// Returns a map of resolved URL -> link info
// Uses fuzzy matching since wikilink URL conversion depends on context
func extractLinkInfoFromMarkdown(mdSource string) map[string]linkInfo {
	result := make(map[string]linkInfo)
	if mdSource == "" {
		return result
	}

	lines := strings.Split(mdSource, "\n")

	// Regex patterns for different link types
	// Wikilinks: [[page]] or [[page|text]]
	wikilinkRegex := regexp.MustCompile(`!?\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)
	// Markdown links: [text](url)
	mdLinkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

	for lineNum, line := range lines {
		// Find wikilinks - store with multiple possible URL variations
		matches := wikilinkRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				target := strings.TrimSpace(match[1])
				text := getLastPathSegment(target)
				if len(match) >= 3 && match[2] != "" {
					text = match[2]
				}

				// Generate URL variations for the wikilink
				// Since we don't know the sourceDir context, we create multiple possibilities
				urls := generateWikilinkURLVariations(target)

				info := linkInfo{
					LineNumber:     lineNum + 1,
					OriginalSyntax: match[0],
					LinkText:       text,
				}

				for _, url := range urls {
					result[url] = info
				}
			}
		}

		// Find markdown links
		matches = mdLinkRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				text := match[1]
				url := match[2]

				// Only process internal links
				if strings.HasPrefix(url, "/") {
					result[url] = linkInfo{
						LineNumber:     lineNum + 1,
						OriginalSyntax: match[0],
						LinkText:       text,
					}
				}
			}
		}
	}

	return result
}

// getLastPathSegment returns the last segment of a path
func getLastPathSegment(path string) string {
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if strings.TrimSpace(parts[i]) != "" {
			return parts[i]
		}
	}
	return path
}

// generateWikilinkURLVariations generates possible URL variations for a wikilink target
// Since we don't know sourceDir context, we generate common possibilities
func generateWikilinkURLVariations(target string) []string {
	var urls []string

	// Remove .md extension if present
	target = strings.TrimSuffix(target, ".md")

	// Handle anchors
	anchor := ""
	if idx := strings.Index(target, "#"); idx != -1 {
		anchor = target[idx:]
		target = target[:idx]
	}

	if target == "" {
		return []string{anchor}
	}

	// Split into path segments and slugify
	parts := strings.Split(target, "/")
	sluggedParts := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// Use tree.Slugify to match actual URL generation
		sluggedParts = append(sluggedParts, tree.Slugify(part))
	}

	if len(sluggedParts) == 0 {
		return []string{"/" + anchor}
	}

	// Check if last segment is index/readme
	lastSegment := sluggedParts[len(sluggedParts)-1]
	isIndex := lastSegment == "index" || lastSegment == "readme"

	// Generate variations
	fullPath := "/" + strings.Join(sluggedParts, "/") + "/" + anchor
	urls = append(urls, fullPath)

	// If last segment is index/readme, also add parent dir version
	if isIndex && len(sluggedParts) > 1 {
		parentPath := "/" + strings.Join(sluggedParts[:len(sluggedParts)-1], "/") + "/" + anchor
		urls = append(urls, parentPath)
	}

	return urls
}

// findSimilarURLs finds URLs that are similar to the broken link (for suggestions)
func findSimilarURLs(brokenURL string, validURLs map[string]bool) []string {
	var suggestions []string

	// Normalize the broken URL for comparison
	broken := strings.ToLower(strings.Trim(brokenURL, "/"))

	// Find URLs with similar paths
	for validURL := range validURLs {
		valid := strings.ToLower(strings.Trim(validURL, "/"))

		// Skip root
		if valid == "" {
			continue
		}

		// Check if they share common segments or are close matches
		if strings.Contains(valid, broken) || strings.Contains(broken, valid) {
			suggestions = append(suggestions, validURL)
			if len(suggestions) >= 3 {
				break
			}
		}
	}

	return suggestions
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

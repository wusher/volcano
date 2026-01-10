package markdown

import (
	"net/url"
	"regexp"
	"strings"
)

// linkRegex matches <a> tags
var linkRegex = regexp.MustCompile(`(?i)<a\s+([^>]*)href="([^"]*)"([^>]*)>(.*?)</a>`)

// hrefRegex matches href attributes in any tag
var hrefRegex = regexp.MustCompile(`(?i)href="(/[^"]*)"`)

// srcRegex matches src attributes in any tag
var srcRegex = regexp.MustCompile(`(?i)src="(/[^"]*)"`)

// srcsetRegex matches srcset attributes in any tag
var srcsetRegex = regexp.MustCompile(`(?i)srcset="([^"]*)"`)

// posterRegex matches poster attributes in video tags
var posterRegex = regexp.MustCompile(`(?i)poster="(/[^"]*)"`)

// dataRegex matches data-* attributes with URL values
var dataRegex = regexp.MustCompile(`(?i)data-[a-z-]+="(/[^"]*)"`)

// ProcessExternalLinks modifies external links to open in new tab and add icons
func ProcessExternalLinks(htmlContent string, siteURL string) string {
	// Parse site URL for comparison
	var siteHost string
	if siteURL != "" {
		if parsed, err := url.Parse(siteURL); err == nil {
			siteHost = strings.ToLower(parsed.Host)
		}
	}

	result := linkRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		matches := linkRegex.FindStringSubmatch(match)
		if len(matches) < 5 {
			return match
		}

		attrsBefore := matches[1]
		href := matches[2]
		attrsAfter := matches[3]
		content := matches[4]

		// Check if this is an external link
		if !isExternalURL(href, siteHost) {
			return match
		}

		// Skip if it's an image link
		if strings.Contains(strings.ToLower(content), "<img") {
			return match
		}

		// Check if already has target attribute
		hasTarget := strings.Contains(strings.ToLower(attrsBefore+attrsAfter), "target=")
		hasRel := strings.Contains(strings.ToLower(attrsBefore+attrsAfter), "rel=")

		// Build new link
		var sb strings.Builder
		sb.WriteString(`<a `)
		sb.WriteString(attrsBefore)
		sb.WriteString(`href="`)
		sb.WriteString(href)
		sb.WriteString(`"`)
		sb.WriteString(attrsAfter)

		if !hasTarget {
			sb.WriteString(` target="_blank"`)
		}
		if !hasRel {
			sb.WriteString(` rel="noopener noreferrer"`)
		}

		sb.WriteString(`>`)
		sb.WriteString(content)
		sb.WriteString(`<svg class="external-icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>`)
		sb.WriteString(`<span class="sr-only">(opens in new tab)</span>`)
		sb.WriteString(`</a>`)

		return sb.String()
	})

	return result
}

// isExternalURL checks if a URL is external to the site
func isExternalURL(href string, siteHost string) bool {
	// Skip empty, anchor-only, or relative URLs
	if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "/") {
		return false
	}

	// Skip mailto, tel, javascript
	if strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "tel:") ||
		strings.HasPrefix(href, "javascript:") {
		return false
	}

	// Parse the URL
	parsed, err := url.Parse(href)
	if err != nil {
		return false
	}

	// If no host, it's relative
	if parsed.Host == "" {
		return false
	}

	// Compare with site host
	linkHost := strings.ToLower(parsed.Host)
	if siteHost != "" && linkHost == siteHost {
		return false
	}

	// Handle www prefix
	if siteHost != "" {
		if linkHost == "www."+siteHost || siteHost == "www."+linkHost {
			return false
		}
	}

	return true
}

// PrefixInternalLinks adds a base URL path to internal links and resources in HTML content.
// Internal paths are those starting with "/" (e.g., "/guides/intro/", "/images/logo.png").
// The baseURL should be a full URL like "https://example.com/volcano/".
// If baseURL has a path component (e.g., "/volcano"), it's prepended to internal paths.
// This handles href, src, srcset, poster, data-*, and content attributes.
func PrefixInternalLinks(htmlContent string, baseURL string) string {
	if baseURL == "" {
		return htmlContent
	}

	// Extract the base path from the URL
	// e.g., "https://wusher.github.io/volcano/" -> "/volcano"
	basePath := extractBasePath(baseURL)
	if basePath == "" {
		return htmlContent
	}

	// Helper function to check if a path should be prefixed
	shouldPrefix := func(path string) bool {
		// Must start with / but not //
		if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
			return false
		}
		// Skip anchor-only links
		if strings.HasPrefix(path, "#") {
			return false
		}
		return true
	}

	// Process href attributes
	htmlContent = hrefRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		matches := hrefRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		path := matches[1]
		if !shouldPrefix(path) {
			return match
		}
		return `href="` + basePath + path + `"`
	})

	// Process src attributes
	htmlContent = srcRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		matches := srcRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		path := matches[1]
		if !shouldPrefix(path) {
			return match
		}
		return `src="` + basePath + path + `"`
	})

	// Process srcset attributes (can have multiple URLs)
	htmlContent = srcsetRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		matches := srcsetRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		srcset := matches[1]

		// Process each URL in srcset (format: "url1 1x, url2 2x")
		parts := strings.Split(srcset, ",")
		for i, part := range parts {
			part = strings.TrimSpace(part)
			// Split URL from descriptor (e.g., "2x" or "100w")
			fields := strings.Fields(part)
			if len(fields) > 0 {
				url := fields[0]
				if shouldPrefix(url) {
					fields[0] = basePath + url
					parts[i] = strings.Join(fields, " ")
				}
			}
		}
		return `srcset="` + strings.Join(parts, ", ") + `"`
	})

	// Process poster attributes
	htmlContent = posterRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		matches := posterRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		path := matches[1]
		if !shouldPrefix(path) {
			return match
		}
		return `poster="` + basePath + path + `"`
	})

	// Process data-* attributes
	htmlContent = dataRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		// Extract the attribute name and value
		parts := strings.SplitN(match, `="`, 2)
		if len(parts) < 2 {
			return match
		}
		attrName := parts[0]
		path := strings.TrimSuffix(parts[1], `"`)
		if !shouldPrefix(path) {
			return match
		}
		return attrName + `="` + basePath + path + `"`
	})

	return htmlContent
}

// extractBasePath extracts the path portion from a URL.
// e.g., "https://example.com/volcano/" -> "/volcano"
// e.g., "https://example.com/" -> ""
func extractBasePath(baseURL string) string {
	// Remove trailing slash for consistent handling
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Find the scheme separator
	schemeEnd := strings.Index(baseURL, "://")
	if schemeEnd == -1 {
		// No scheme, might just be a path
		if strings.HasPrefix(baseURL, "/") {
			return baseURL
		}
		return ""
	}

	// Find the first slash after the scheme (start of path)
	pathStart := strings.Index(baseURL[schemeEnd+3:], "/")
	if pathStart == -1 {
		// No path component
		return ""
	}

	// Extract the path
	path := baseURL[schemeEnd+3+pathStart:]
	if path == "" || path == "/" {
		return ""
	}

	return path
}

package markdown

import (
	"net/url"
	"regexp"
	"strings"
)

// linkRegex matches <a> tags
var linkRegex = regexp.MustCompile(`(?i)<a\s+([^>]*)href="([^"]*)"([^>]*)>(.*?)</a>`)

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

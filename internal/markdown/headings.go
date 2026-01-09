package markdown

import (
	"regexp"
	"strings"
	"unicode"
)

// HeadingID tracks heading IDs for uniqueness
type HeadingID struct {
	Original string
	Slug     string
	Count    int
}

// Slugify converts text to a URL-friendly slug
func Slugify(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)

	// Replace spaces and underscores with hyphens
	text = strings.ReplaceAll(text, " ", "-")
	text = strings.ReplaceAll(text, "_", "-")

	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	prevHyphen := false
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			result.WriteRune(r)
			prevHyphen = false
		} else if r == '-' && !prevHyphen {
			result.WriteRune('-')
			prevHyphen = true
		}
	}

	// Trim leading/trailing hyphens
	return strings.Trim(result.String(), "-")
}

// headingRegex matches h1-h6 tags
var anchorHeadingRegex = regexp.MustCompile(`(?i)<(h[1-6])([^>]*)>(.*?)</h[1-6]>`)
var existingIDRegex = regexp.MustCompile(`\s+id="[^"]*"`)

// AddHeadingAnchors adds anchor links to all headings in HTML content
func AddHeadingAnchors(htmlContent string) string {
	seenIDs := make(map[string]int)

	result := anchorHeadingRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		matches := anchorHeadingRegex.FindStringSubmatch(match)
		if len(matches) < 4 {
			return match
		}

		tag := matches[1]       // h1, h2, etc
		attrs := matches[2]     // existing attributes
		content := matches[3]   // heading content

		// Strip HTML tags from content to get plain text for slug
		plainText := stripHTMLTags(content)
		plainText = strings.TrimSpace(plainText)

		if plainText == "" {
			return match
		}

		// Generate unique ID
		baseSlug := Slugify(plainText)
		if baseSlug == "" {
			baseSlug = "heading"
		}

		slug := baseSlug
		if count, exists := seenIDs[baseSlug]; exists {
			slug = baseSlug + "-" + itoa(count)
			seenIDs[baseSlug] = count + 1
		} else {
			seenIDs[baseSlug] = 1
		}

		// Remove any existing id attribute
		attrs = existingIDRegex.ReplaceAllString(attrs, "")

		// Build new heading with anchor
		var sb strings.Builder
		sb.WriteString("<")
		sb.WriteString(tag)
		sb.WriteString(` id="`)
		sb.WriteString(slug)
		sb.WriteString(`"`)
		sb.WriteString(attrs)
		sb.WriteString(">")
		sb.WriteString(`<a href="#`)
		sb.WriteString(slug)
		sb.WriteString(`" class="heading-anchor" aria-label="Link to `)
		sb.WriteString(escapeAttr(plainText))
		sb.WriteString(` section">`)
		sb.WriteString(`<svg class="anchor-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path></svg>`)
		sb.WriteString(`</a>`)
		sb.WriteString(content)
		sb.WriteString("</")
		sb.WriteString(tag)
		sb.WriteString(">")

		return sb.String()
	})

	return result
}

// stripHTMLTags removes all HTML tags from a string
func stripHTMLTags(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
		} else if r == '>' {
			inTag = false
		} else if !inTag {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// escapeAttr escapes a string for use in an HTML attribute
func escapeAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// itoa converts int to string
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + itoa(-n)
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

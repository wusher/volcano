// Package seo provides SEO-related functionality like meta tags and Open Graph.
package seo

import (
	"html/template"
	"regexp"
	"strings"
)

// SEOMeta contains SEO meta tag data
type SEOMeta struct {
	Title       string // Page title
	Description string // Meta description
	Canonical   string // Canonical URL
	Robots      string // Robots directive
	Author      string // Author name
}

// OpenGraph contains Open Graph meta tag data
type OpenGraph struct {
	Title       string // og:title
	Description string // og:description
	Type        string // og:type (article, website)
	URL         string // og:url
	SiteName    string // og:site_name
	Image       string // og:image
}

// TwitterCard contains Twitter Card meta tag data
type TwitterCard struct {
	Card        string // twitter:card (summary, summary_large_image)
	Title       string // twitter:title
	Description string // twitter:description
	Image       string // twitter:image
}

// PageMeta combines all meta information for a page
type PageMeta struct {
	SEO     SEOMeta
	OG      OpenGraph
	Twitter TwitterCard
}

// Config holds site-wide SEO configuration
type Config struct {
	SiteURL       string
	SiteTitle     string
	DefaultDesc   string
	Author        string
	OGImage       string
	TwitterHandle string
}

var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)
var whitespaceRegex = regexp.MustCompile(`\s+`)

// GeneratePageMeta creates all meta data for a page
func GeneratePageMeta(pageTitle, pageContent, urlPath string, config Config) PageMeta {
	description := extractDescription(pageContent, 160)
	if description == "" && config.DefaultDesc != "" {
		description = config.DefaultDesc
	}

	fullTitle := pageTitle
	if config.SiteTitle != "" && pageTitle != config.SiteTitle {
		fullTitle = pageTitle + " - " + config.SiteTitle
	}

	canonical := ""
	if config.SiteURL != "" {
		canonical = strings.TrimSuffix(config.SiteURL, "/") + urlPath
	}

	return PageMeta{
		SEO: SEOMeta{
			Title:       fullTitle,
			Description: description,
			Canonical:   canonical,
			Robots:      "index, follow",
			Author:      config.Author,
		},
		OG: OpenGraph{
			Title:       pageTitle,
			Description: description,
			Type:        "article",
			URL:         canonical,
			SiteName:    config.SiteTitle,
			Image:       config.OGImage,
		},
		Twitter: TwitterCard{
			Card:        getTwitterCardType(config.OGImage),
			Title:       pageTitle,
			Description: description,
			Image:       config.OGImage,
		},
	}
}

// extractDescription extracts a description from HTML content
func extractDescription(htmlContent string, maxLen int) string {
	// Strip HTML tags
	text := htmlTagRegex.ReplaceAllString(htmlContent, " ")

	// Normalize whitespace
	text = whitespaceRegex.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	// Find first paragraph of reasonable content
	if len(text) == 0 {
		return ""
	}

	// Truncate to maxLen
	if len(text) > maxLen {
		// Find last space before maxLen
		truncated := text[:maxLen]
		lastSpace := strings.LastIndex(truncated, " ")
		if lastSpace > maxLen/2 {
			text = truncated[:lastSpace] + "..."
		} else {
			text = truncated + "..."
		}
	}

	return text
}

// getTwitterCardType returns the appropriate Twitter card type
func getTwitterCardType(imageURL string) string {
	if imageURL != "" {
		return "summary_large_image"
	}
	return "summary"
}

// RenderMetaTags renders all meta tags as HTML
func RenderMetaTags(meta PageMeta) template.HTML {
	var sb strings.Builder

	// SEO meta tags
	if meta.SEO.Description != "" {
		sb.WriteString(`  <meta name="description" content="`)
		sb.WriteString(template.HTMLEscapeString(meta.SEO.Description))
		sb.WriteString(`">`)
		sb.WriteString("\n")
	}

	if meta.SEO.Robots != "" {
		sb.WriteString(`  <meta name="robots" content="`)
		sb.WriteString(template.HTMLEscapeString(meta.SEO.Robots))
		sb.WriteString(`">`)
		sb.WriteString("\n")
	}

	if meta.SEO.Author != "" {
		sb.WriteString(`  <meta name="author" content="`)
		sb.WriteString(template.HTMLEscapeString(meta.SEO.Author))
		sb.WriteString(`">`)
		sb.WriteString("\n")
	}

	if meta.SEO.Canonical != "" {
		sb.WriteString(`  <link rel="canonical" href="`)
		sb.WriteString(template.HTMLEscapeString(meta.SEO.Canonical))
		sb.WriteString(`">`)
		sb.WriteString("\n")
	}

	// Open Graph tags
	sb.WriteString("\n")
	sb.WriteString(`  <!-- Open Graph -->`)
	sb.WriteString("\n")

	sb.WriteString(`  <meta property="og:title" content="`)
	sb.WriteString(template.HTMLEscapeString(meta.OG.Title))
	sb.WriteString(`">`)
	sb.WriteString("\n")

	if meta.OG.Description != "" {
		sb.WriteString(`  <meta property="og:description" content="`)
		sb.WriteString(template.HTMLEscapeString(meta.OG.Description))
		sb.WriteString(`">`)
		sb.WriteString("\n")
	}

	sb.WriteString(`  <meta property="og:type" content="`)
	sb.WriteString(template.HTMLEscapeString(meta.OG.Type))
	sb.WriteString(`">`)
	sb.WriteString("\n")

	if meta.OG.URL != "" {
		sb.WriteString(`  <meta property="og:url" content="`)
		sb.WriteString(template.HTMLEscapeString(meta.OG.URL))
		sb.WriteString(`">`)
		sb.WriteString("\n")
	}

	if meta.OG.SiteName != "" {
		sb.WriteString(`  <meta property="og:site_name" content="`)
		sb.WriteString(template.HTMLEscapeString(meta.OG.SiteName))
		sb.WriteString(`">`)
		sb.WriteString("\n")
	}

	if meta.OG.Image != "" {
		sb.WriteString(`  <meta property="og:image" content="`)
		sb.WriteString(template.HTMLEscapeString(meta.OG.Image))
		sb.WriteString(`">`)
		sb.WriteString("\n")
	}

	// Twitter Card tags
	sb.WriteString("\n")
	sb.WriteString(`  <!-- Twitter Card -->`)
	sb.WriteString("\n")

	sb.WriteString(`  <meta name="twitter:card" content="`)
	sb.WriteString(template.HTMLEscapeString(meta.Twitter.Card))
	sb.WriteString(`">`)
	sb.WriteString("\n")

	sb.WriteString(`  <meta name="twitter:title" content="`)
	sb.WriteString(template.HTMLEscapeString(meta.Twitter.Title))
	sb.WriteString(`">`)
	sb.WriteString("\n")

	if meta.Twitter.Description != "" {
		sb.WriteString(`  <meta name="twitter:description" content="`)
		sb.WriteString(template.HTMLEscapeString(meta.Twitter.Description))
		sb.WriteString(`">`)
		sb.WriteString("\n")
	}

	if meta.Twitter.Image != "" {
		sb.WriteString(`  <meta name="twitter:image" content="`)
		sb.WriteString(template.HTMLEscapeString(meta.Twitter.Image))
		sb.WriteString(`">`)
		sb.WriteString("\n")
	}

	return template.HTML(sb.String())
}

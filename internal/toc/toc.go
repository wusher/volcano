// Package toc provides table of contents generation from HTML content.
package toc

import (
	"html/template"
	"regexp"
	"strings"
)

// TOCItem represents a single item in the table of contents
type TOCItem struct {
	ID       string     // Heading ID for anchor link
	Text     string     // Heading text content
	Level    int        // 2, 3, or 4
	Children []*TOCItem // Nested headings
}

// PageTOC represents the table of contents for a page
type PageTOC struct {
	Items    []*TOCItem
	MinItems int // Minimum headings to show TOC (default: 3)
}

// headingRegex matches h2, h3, h4 headings with id and content
var headingRegex = regexp.MustCompile(`(?i)<h([2-4])[^>]*\s+id="([^"]+)"[^>]*>(.*?)</h[2-4]>`)

// stripTagsRegex removes HTML tags
var stripTagsRegex = regexp.MustCompile(`<[^>]*>`)

// ExtractTOC extracts table of contents from rendered HTML content
func ExtractTOC(htmlContent string, minItems int) *PageTOC {
	if minItems <= 0 {
		minItems = 3
	}

	matches := headingRegex.FindAllStringSubmatch(htmlContent, -1)

	if len(matches) < minItems {
		return nil
	}

	var items []*TOCItem
	var stack []*TOCItem

	for _, match := range matches {
		level := int(match[1][0] - '0') // Convert '2', '3', '4' to 2, 3, 4
		id := match[2]
		text := stripTagsRegex.ReplaceAllString(match[3], "")
		text = strings.TrimSpace(text)

		item := &TOCItem{
			ID:    id,
			Text:  text,
			Level: level,
		}

		// Find parent based on level
		for len(stack) > 0 && stack[len(stack)-1].Level >= level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			items = append(items, item)
		} else {
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, item)
		}

		stack = append(stack, item)
	}

	return &PageTOC{
		Items:    items,
		MinItems: minItems,
	}
}

// RenderTOC renders the table of contents as HTML
func RenderTOC(toc *PageTOC) template.HTML {
	if toc == nil || len(toc.Items) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<aside class="toc-sidebar" aria-label="Table of contents">`)
	sb.WriteString("\n")
	sb.WriteString(`  <nav class="toc">`)
	sb.WriteString("\n")
	sb.WriteString(`    <h2 class="toc-title">On this page</h2>`)
	sb.WriteString("\n")

	renderTOCItems(&sb, toc.Items, 2)

	sb.WriteString(`  </nav>`)
	sb.WriteString("\n")
	sb.WriteString(`</aside>`)
	sb.WriteString("\n")

	return template.HTML(sb.String())
}

func renderTOCItems(sb *strings.Builder, items []*TOCItem, indent int) {
	if len(items) == 0 {
		return
	}

	indentStr := strings.Repeat("  ", indent)

	sb.WriteString(indentStr)
	sb.WriteString("<ul>\n")

	for _, item := range items {
		sb.WriteString(indentStr)
		sb.WriteString("  <li>\n")
		sb.WriteString(indentStr)
		sb.WriteString(`    <a href="#`)
		sb.WriteString(template.HTMLEscapeString(item.ID))
		sb.WriteString(`">`)
		sb.WriteString(template.HTMLEscapeString(item.Text))
		sb.WriteString("</a>\n")

		if len(item.Children) > 0 {
			renderTOCItems(sb, item.Children, indent+2)
		}

		sb.WriteString(indentStr)
		sb.WriteString("  </li>\n")
	}

	sb.WriteString(indentStr)
	sb.WriteString("</ul>\n")
}

// HasTOC returns true if the content has enough headings for a TOC
func HasTOC(htmlContent string, minItems int) bool {
	if minItems <= 0 {
		minItems = 3
	}
	matches := headingRegex.FindAllString(htmlContent, minItems)
	return len(matches) >= minItems
}

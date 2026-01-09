package markdown

import (
	"regexp"
	"strings"
)

// AdmonitionType represents the type of admonition
type AdmonitionType string

// Admonition type constants
const (
	AdmonitionNote    AdmonitionType = "note"
	AdmonitionTip     AdmonitionType = "tip"
	AdmonitionWarning AdmonitionType = "warning"
	AdmonitionDanger  AdmonitionType = "danger"
	AdmonitionInfo    AdmonitionType = "info"
)

// Admonition represents a parsed admonition block
type Admonition struct {
	Type    AdmonitionType
	Title   string
	Content string
}

// admonitionRegex matches :::type content ::: blocks in markdown
// This is applied before markdown parsing
var admonitionStartRegex = regexp.MustCompile(`(?m)^:::(note|tip|warning|danger|info)(?:\s+(.+))?$`)
var admonitionEndRegex = regexp.MustCompile(`(?m)^:::$`)

// defaultTitles maps admonition types to default titles
var defaultTitles = map[AdmonitionType]string{
	AdmonitionNote:    "Note",
	AdmonitionTip:     "Tip",
	AdmonitionWarning: "Warning",
	AdmonitionDanger:  "Danger",
	AdmonitionInfo:    "Info",
}

// admonitionIcons maps admonition types to SVG icons
var admonitionIcons = map[AdmonitionType]string{
	AdmonitionNote:    `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>`,
	AdmonitionTip:     `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 18h6"></path><path d="M10 22h4"></path><path d="M15.09 14c.18-.98.65-1.74 1.41-2.5A4.65 4.65 0 0 0 18 8 6 6 0 0 0 6 8c0 1 .23 2.23 1.5 3.5A4.61 4.61 0 0 1 8.91 14"></path></svg>`,
	AdmonitionWarning: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>`,
	AdmonitionDanger:  `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="7.86 2 16.14 2 22 7.86 22 16.14 16.14 22 7.86 22 2 16.14 2 7.86 7.86 2"></polygon><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>`,
	AdmonitionInfo:    `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>`,
}

// ProcessAdmonitions converts :::type ::: blocks to HTML admonitions
// This should be called on the markdown content BEFORE parsing
func ProcessAdmonitions(markdown string) string {
	lines := strings.Split(markdown, "\n")
	var result []string
	var inAdmonition bool
	var currentType AdmonitionType
	var currentTitle string
	var contentLines []string

	for _, line := range lines {
		if inAdmonition {
			if admonitionEndRegex.MatchString(line) {
				// End of admonition - output the HTML
				result = append(result, renderAdmonitionHTML(currentType, currentTitle, strings.Join(contentLines, "\n")))
				inAdmonition = false
				contentLines = nil
			} else {
				contentLines = append(contentLines, line)
			}
		} else {
			matches := admonitionStartRegex.FindStringSubmatch(line)
			if len(matches) >= 2 {
				// Start of admonition
				inAdmonition = true
				currentType = AdmonitionType(matches[1])
				if len(matches) >= 3 && matches[2] != "" {
					currentTitle = strings.TrimSpace(matches[2])
				} else {
					currentTitle = defaultTitles[currentType]
				}
				contentLines = nil
			} else {
				result = append(result, line)
			}
		}
	}

	// Handle unclosed admonition
	if inAdmonition && len(contentLines) > 0 {
		result = append(result, renderAdmonitionHTML(currentType, currentTitle, strings.Join(contentLines, "\n")))
	}

	return strings.Join(result, "\n")
}

// renderAdmonitionHTML generates HTML for an admonition
// The content is still markdown and will be parsed later
func renderAdmonitionHTML(adType AdmonitionType, title string, content string) string {
	icon := admonitionIcons[adType]
	if icon == "" {
		icon = admonitionIcons[AdmonitionInfo]
	}

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(`<div class="admonition admonition-`)
	sb.WriteString(string(adType))
	sb.WriteString(`" role="note">`)
	sb.WriteString("\n")
	sb.WriteString(`  <div class="admonition-heading">`)
	sb.WriteString("\n")
	sb.WriteString(`    <span class="admonition-icon">`)
	sb.WriteString(icon)
	sb.WriteString(`</span>`)
	sb.WriteString("\n")
	sb.WriteString(`    <span class="admonition-title">`)
	sb.WriteString(escapeHTML(title))
	sb.WriteString(`</span>`)
	sb.WriteString("\n")
	sb.WriteString(`  </div>`)
	sb.WriteString("\n")
	sb.WriteString(`  <div class="admonition-content">`)
	sb.WriteString("\n")
	sb.WriteString("\n")
	sb.WriteString(content)
	sb.WriteString("\n")
	sb.WriteString("\n")
	sb.WriteString(`  </div>`)
	sb.WriteString("\n")
	sb.WriteString(`</div>`)
	sb.WriteString("\n")

	return sb.String()
}

// escapeHTML escapes HTML special characters
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

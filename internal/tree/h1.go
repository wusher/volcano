package tree

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
)

var (
	// h1Regex matches the first H1 heading in markdown
	h1Regex = regexp.MustCompile(`^#\s+(.+)$`)
	// linkRegex matches markdown links like [text](url)
	linkRegex = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`)
	// inlineMarkdownRegex matches inline markdown formatting
	inlineMarkdownRegex = regexp.MustCompile("[*_~`]")
)

// ExtractH1 extracts the first H1 heading from markdown content.
// Returns empty string if no H1 is found at the beginning of the file.
func ExtractH1(content []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	inFrontmatter := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Handle frontmatter
		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			// End of frontmatter
			inFrontmatter = false
			continue
		}

		// Skip if we're in frontmatter
		if inFrontmatter {
			continue
		}

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for H1
		if matches := h1Regex.FindStringSubmatch(line); len(matches) > 1 {
			title := matches[1]
			// Strip inline markdown formatting
			title = stripInlineMarkdown(title)
			title = strings.TrimSpace(title)
			if title != "" {
				return title
			}
		}

		// If first non-empty, non-frontmatter line isn't H1, stop looking
		// (H1 should be at the top of the content)
		break
	}

	return ""
}

// stripInlineMarkdown removes inline markdown formatting from text.
// This includes bold, italic, code, links, etc.
func stripInlineMarkdown(text string) string {
	// Replace links [text](url) with just text
	text = linkRegex.ReplaceAllString(text, "$1")
	// Remove inline formatting characters: * _ ~ `
	text = inlineMarkdownRegex.ReplaceAllString(text, "")
	return text
}

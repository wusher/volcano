package markdown

import (
	"regexp"
	"strings"

	"volcano/internal/tree"
)

// wikiLinkRegex matches Obsidian-style wiki links: [[Page]] or [[Page|Display Text]]
// Also captures optional ! prefix for embeds: ![[Page]]
var wikiLinkRegex = regexp.MustCompile(`!?\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)

// ConvertWikiLinks converts Obsidian-style [[wiki links]] to standard markdown links.
// Handles formats:
//   - [[Page Name]] -> [Page Name](/current-dir/page-name/) (relative to current page)
//   - [[Page Name|Display]] -> [Display](/current-dir/page-name/)
//   - [[folder/Page Name]] -> [Page Name](/folder/page-name/) (explicit path from root)
//   - ![[Page Name]] -> [Page Name](/page-name/) (embeds converted to links)
//
// sourceDir is the source file's directory (e.g., "/guides/") for sibling link resolution.
// Links without a path separator are resolved relative to this directory.
func ConvertWikiLinks(content []byte, sourceDir string) []byte {
	result := wikiLinkRegex.ReplaceAllFunc(content, func(match []byte) []byte {
		submatch := wikiLinkRegex.FindSubmatch(match)
		if len(submatch) < 2 {
			return match
		}

		target := string(submatch[1])
		displayText := ""
		if len(submatch) >= 3 && len(submatch[2]) > 0 {
			displayText = string(submatch[2])
		}

		// Clean up the target path
		target = strings.TrimSpace(target)

		// Remove .md extension if present
		target = strings.TrimSuffix(target, ".md")

		// Get display text (use filename if not specified)
		if displayText == "" {
			// Use the last part of the path as display text
			parts := strings.Split(target, "/")
			displayText = parts[len(parts)-1]
		}

		// Convert target to URL path, considering source file's directory
		urlPath := convertToURLPath(target, sourceDir)

		// Return standard markdown link
		return []byte("[" + displayText + "](" + urlPath + ")")
	})

	return result
}

// convertToURLPath converts a wiki link target to a URL path.
// If the target has no path separator (e.g., "Page Name"), it's resolved relative
// to the source file's directory (sibling resolution).
// If it has a path (e.g., "folder/Page"), it's resolved from the root.
// Special handling: index/readme files resolve to their parent directory.
// Anchors (e.g., #section) are preserved and appended to the final URL.
func convertToURLPath(target string, sourceDir string) string {
	// Extract anchor/fragment if present (e.g., "faq#permissions" -> "faq", "#permissions")
	anchor := ""
	if idx := strings.Index(target, "#"); idx != -1 {
		anchor = target[idx:] // includes the #
		target = target[:idx]
	}

	// Handle empty target (just an anchor like [[#section]])
	if target == "" {
		return anchor
	}

	// Check if target contains a path separator (explicit path from root)
	hasExplicitPath := strings.Contains(target, "/")

	// Split into path segments
	parts := strings.Split(target, "/")

	// Slugify each segment
	sluggedParts := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		sluggedParts = append(sluggedParts, tree.Slugify(part))
	}

	if len(sluggedParts) == 0 {
		return "/" + anchor
	}

	// Check if the last segment is index or readme (these resolve to parent directory)
	lastSegment := sluggedParts[len(sluggedParts)-1]
	isIndexFile := lastSegment == "index" || lastSegment == "readme"

	// If it's an index/readme file, remove it from the path (it resolves to parent dir)
	if isIndexFile {
		sluggedParts = sluggedParts[:len(sluggedParts)-1]
	}

	// If no explicit path, resolve relative to the source file's directory
	// This means [[sibling-page]] from a file in /guides/ resolves to /guides/sibling-page/
	if !hasExplicitPath && sourceDir != "" && sourceDir != "/" {
		// sourceDir is already the directory (e.g., "/guides/")
		currentDir := strings.Trim(sourceDir, "/")
		if currentDir != "" {
			if len(sluggedParts) == 0 {
				// [[index]] in /guides/ -> /guides/
				return "/" + currentDir + "/" + anchor
			}
			return "/" + currentDir + "/" + strings.Join(sluggedParts, "/") + "/" + anchor
		}
	}

	if len(sluggedParts) == 0 {
		return "/" + anchor
	}

	return "/" + strings.Join(sluggedParts, "/") + "/" + anchor
}

package markdown

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wusher/volcano/internal/tree"
)

// attachmentExtensions lists file extensions that should be treated as attachments
// (not markdown pages). These preserve their extension and don't get a trailing slash.
var attachmentExtensions = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
	".svg": true, ".bmp": true, ".ico": true, ".heic": true,
	".pdf": true,
	".mp3": true, ".mp4": true, ".wav": true, ".ogg": true, ".webm": true, ".mov": true,
	".zip": true, ".docx": true, ".xlsx": true, ".pptx": true,
}

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

		// Remove .md extension if present (handle case with anchor: "file.md#section")
		if idx := strings.Index(target, ".md#"); idx != -1 {
			target = target[:idx] + target[idx+3:] // Remove .md but keep #anchor
		} else {
			target = strings.TrimSuffix(target, ".md")
		}

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

// isAttachment checks if a filename has an attachment extension
func isAttachment(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return attachmentExtensions[ext]
}

// convertToURLPath converts a wiki link target to a URL path.
// If the target has no path separator (e.g., "Page Name"), it's resolved relative
// to the source file's directory (sibling resolution).
// If it has a path (e.g., "folder/Page"), it's resolved from the root.
// Special handling: index/readme files resolve to their parent directory.
// Anchors (e.g., #section) are preserved and appended to the final URL.
// Attachments (images, PDFs, etc.) preserve their extension and don't get trailing slash.
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

	// Check if this is an attachment (image, PDF, etc.)
	isAttachmentLink := isAttachment(target)

	// Check if target contains a path separator (explicit path from root)
	hasExplicitPath := strings.Contains(target, "/")

	// Split into path segments
	parts := strings.Split(target, "/")

	// Process each segment
	sluggedParts := make([]string, 0, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// For attachments, preserve the filename but slugify directory parts
		if isAttachmentLink && i == len(parts)-1 {
			// For the filename, just lowercase and replace spaces with dashes
			// but preserve the extension
			ext := filepath.Ext(part)
			name := strings.TrimSuffix(part, ext)
			name = strings.ToLower(name)
			name = strings.ReplaceAll(name, " ", "-")
			sluggedParts = append(sluggedParts, name+strings.ToLower(ext))
		} else {
			sluggedParts = append(sluggedParts, tree.Slugify(part))
		}
	}

	if len(sluggedParts) == 0 {
		return "/" + anchor
	}

	// For attachments, don't add trailing slash
	if isAttachmentLink {
		if !hasExplicitPath && sourceDir != "" && sourceDir != "/" {
			currentDir := strings.Trim(sourceDir, "/")
			if currentDir != "" {
				return "/" + currentDir + "/" + strings.Join(sluggedParts, "/") + anchor
			}
		}
		return "/" + strings.Join(sluggedParts, "/") + anchor
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

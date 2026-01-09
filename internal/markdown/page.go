package markdown

import (
	"os"
	"path/filepath"
	"strings"
)

// Page represents a parsed markdown page
type Page struct {
	Title      string // First H1 or clean filename
	Content    string // Rendered HTML content
	SourcePath string // Path to original .md file
	OutputPath string // Path for output .html file
	URLPath    string // URL path for navigation links
}

// ParseFile reads and parses a markdown file, returning a Page
func ParseFile(sourcePath string, outputPath string, urlPath string, fallbackTitle string) (*Page, error) {
	// Read the file
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return nil, err
	}

	return ParseContent(content, sourcePath, outputPath, urlPath, fallbackTitle)
}

// ParseContent parses pre-read markdown content, returning a Page.
// This allows preprocessing (e.g., admonitions) before parsing.
func ParseContent(content []byte, sourcePath string, outputPath string, urlPath string, fallbackTitle string) (*Page, error) {
	// Extract title from markdown
	title := ExtractTitle(content)
	if title == "" {
		title = fallbackTitle
	}

	// Parse markdown to HTML
	parser := NewParser()
	html, err := parser.Parse(content)
	if err != nil {
		return nil, err
	}

	return &Page{
		Title:      title,
		Content:    string(html),
		SourcePath: sourcePath,
		OutputPath: outputPath,
		URLPath:    urlPath,
	}, nil
}

// CleanFilenameTitle generates a clean title from a filename
// e.g., "getting-started.md" -> "Getting Started"
func CleanFilenameTitle(filename string) string {
	// Remove extension
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Remove leading numbers and date prefixes
	name = removeLeadingPrefix(name)

	// Replace dashes and underscores with spaces
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")

	// Title case
	return strings.Title(strings.TrimSpace(name)) //nolint:staticcheck // strings.Title is fine for this use case
}

// removeLeadingPrefix removes leading numeric or date prefixes
func removeLeadingPrefix(s string) string {
	// Handle date-like prefixes (YYYY-MM-DD-)
	if len(s) >= 11 && isDatePrefix(s[:11]) {
		return s[11:]
	}

	// Handle simple numeric prefixes (01-, 001_, etc.)
	for i, r := range s {
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '-' || r == '_' {
			if i+1 < len(s) {
				return s[i+1:]
			}
		}
		return s
	}

	return s
}

// isDatePrefix checks if the string looks like "YYYY-MM-DD-"
func isDatePrefix(s string) bool {
	if len(s) != 11 {
		return false
	}
	// Format: YYYY-MM-DD-
	for i, r := range s {
		switch i {
		case 4, 7:
			if r != '-' {
				return false
			}
		case 10:
			if r != '-' {
				return false
			}
		default:
			if r < '0' || r > '9' {
				return false
			}
		}
	}
	return true
}

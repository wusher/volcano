package tree

import (
	"path/filepath"
	"strings"
	"unicode"
)

// CleanLabel converts a filename to a clean display label
// Examples:
//   - getting-started.md → "Getting Started"
//   - api_reference.md → "Api Reference"
//   - FAQ.md → "FAQ"
//   - 01-introduction.md → "Introduction"
func CleanLabel(filename string) string {
	// Remove extension
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Remove leading numbers and separators (e.g., "01-", "01_", "2024-01-01-")
	name = removeLeadingNumbers(name)

	// Replace - and _ with spaces
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")

	// Title case each word
	name = titleCase(name)

	// Clean up multiple spaces
	name = strings.Join(strings.Fields(name), " ")

	return name
}

// removeLeadingNumbers removes leading numeric prefixes like "01-", "001_", "2024-01-01-"
func removeLeadingNumbers(s string) string {
	// Handle date-like prefixes (YYYY-MM-DD-)
	if len(s) >= 11 && isDatePrefix(s[:11]) {
		return s[11:]
	}

	// Handle simple numeric prefixes (01-, 001_, etc.)
	for i, r := range s {
		if !unicode.IsDigit(r) {
			if r == '-' || r == '_' {
				// Skip the separator too
				if i+1 < len(s) {
					return s[i+1:]
				}
			}
			// Not a numeric prefix, return as-is
			return s
		}
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
			if !unicode.IsDigit(r) {
				return false
			}
		}
	}
	return true
}

// titleCase converts a string to title case
func titleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			// Handle all-uppercase words (like FAQ, API) - keep them as-is
			if isAllUppercase(word) {
				continue
			}
			// Capitalize first letter
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

// isAllUppercase checks if a string is all uppercase letters
func isAllUppercase(s string) bool {
	hasLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				return false
			}
		}
	}
	return hasLetter
}

// IsMarkdownFile checks if the filename has a markdown extension (case insensitive)
func IsMarkdownFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".md" || ext == ".markdown"
}

// IsHidden checks if the filename or path starts with a dot
func IsHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

// IsIndexFile checks if the filename is an index file
func IsIndexFile(filename string) bool {
	name := strings.ToLower(filename)
	return name == "index.md" || name == "index.markdown" ||
		name == "readme.md" || name == "readme.markdown"
}

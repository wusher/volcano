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
//   - 0. Inbox → "Inbox" (folder with dot)
func CleanLabel(filename string) string {
	name := filename

	// Only remove .md/.markdown extensions (not arbitrary dots)
	lowerName := strings.ToLower(name)
	if strings.HasSuffix(lowerName, ".md") {
		name = name[:len(name)-3]
	} else if strings.HasSuffix(lowerName, ".markdown") {
		name = name[:len(name)-9]
	}

	// Remove leading numbers and separators (e.g., "01-", "01_", "0. ", "2024-01-01-")
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

// removeLeadingNumbers removes leading numeric prefixes like "01-", "001_", "0. ", "01 ", "6 - ", "2024-01-01-"
// Does NOT strip year-like numbers (4 digits, no leading zeros) followed by just space
func removeLeadingNumbers(s string) string {
	// Handle date-like prefixes (YYYY-MM-DD- or YYYY-MM-DD )
	if len(s) >= 11 && isDatePrefix(s[:11]) {
		return s[11:]
	}
	// Also check for date with space separator (YYYY-MM-DD )
	if len(s) >= 11 && isDatePrefixWithSpace(s[:11]) {
		return strings.TrimLeft(s[11:], " ")
	}

	// Find where the leading digits end
	digitEnd := 0
	for i, r := range s {
		if !unicode.IsDigit(r) {
			digitEnd = i
			break
		}
		digitEnd = i + 1
	}

	// No digits at start
	if digitEnd == 0 {
		return s
	}

	// No separator after digits
	if digitEnd >= len(s) {
		return s
	}

	numStr := s[:digitEnd]
	sep := s[digitEnd]

	// Handle dash or underscore separator: "01-title", "01_title"
	if sep == '-' || sep == '_' {
		if digitEnd+1 < len(s) {
			return s[digitEnd+1:]
		}
		return s
	}

	// Handle dot separator: "0. title", "0.title"
	if sep == '.' {
		rest := s[digitEnd+1:]
		rest = strings.TrimLeft(rest, " ")
		if len(rest) > 0 {
			return rest
		}
		return s
	}

	// Handle space separator: "01 title", "6 - title"
	if sep == ' ' {
		rest := s[digitEnd+1:]

		// Check for space-dash-space pattern: "6 - title"
		if len(rest) > 0 && (rest[0] == '-' || rest[0] == '_') {
			rest = strings.TrimLeft(rest[1:], " ")
			if len(rest) > 0 {
				return rest
			}
			return s
		}

		// For space-only separator, check if number looks like a year
		// Year-like: 4 digits, starts with 1-9 (no leading zeros), like 2023, 1999
		// NOT year-like: 01, 001, 0001, 1, 12 (ordering prefixes)
		isYearLike := len(numStr) == 4 && numStr[0] >= '1' && numStr[0] <= '9'

		if isYearLike {
			// Don't strip year-like numbers with space-only separator
			// "2023 Goals" stays as "2023 Goals"
			return s
		}

		// Strip ordering prefixes like "01 title"
		rest = strings.TrimLeft(rest, " ")
		if len(rest) > 0 {
			return rest
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

// isDatePrefixWithSpace checks if the string looks like "YYYY-MM-DD " (space at end)
func isDatePrefixWithSpace(s string) bool {
	if len(s) != 11 {
		return false
	}
	// Format: YYYY-MM-DD (space at position 10)
	for i, r := range s {
		switch i {
		case 4, 7:
			if r != '-' {
				return false
			}
		case 10:
			if r != ' ' {
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

// Slugify converts a string to a URL-safe slug
// Examples:
//   - "0. Inbox" → "inbox"
//   - "Hello World" → "hello-world"
//   - "API Reference" → "api-reference"
func Slugify(s string) string {
	// Remove leading number prefixes first
	s = removeLeadingNumbers(s)

	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces, underscores with dashes
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove dots and other non-URL-safe characters
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			result.WriteRune(r)
		}
	}

	// Clean up multiple dashes
	slug := result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim leading/trailing dashes
	slug = strings.Trim(slug, "-")

	return slug
}

package tree

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	// datePrefix matches YYYY-MM-DD followed by separator (-, _, or space)
	// Supports: 2024-01-15-title, 2024_01_15_title, "2024-01-15 title"
	datePrefixRegex = regexp.MustCompile(`^(\d{4})[-_](\d{2})[-_](\d{2})[-_\s](.+)$`)
	// numberPrefix matches leading digits followed by separator (-, _, ., or space)
	// Supports: 01-title, 01_title, "01 title", "0. title"
	numberPrefixRegex = regexp.MustCompile(`^(\d+)[-_.\s]\s*(.+)$`)
)

// FileMetadata holds parsed metadata from a filename
type FileMetadata struct {
	OriginalName string    // "2024-01-15-01-hello-world.md"
	Date         time.Time // Parsed from filename prefix (zero if none)
	HasDate      bool      // true if date was extracted from filename
	Number       *int      // nil if no number prefix
	Slug         string    // "hello-world"
	DisplayName  string    // "Hello World"
	IsDraft      bool      // starts with "_"
}

// ExtractFileMetadata parses date/number prefixes from filename
// Only extracts date from filename prefix - does not use file modification time
func ExtractFileMetadata(filename string, modTime time.Time) FileMetadata {
	// Only strip known markdown extensions, not arbitrary "extensions"
	// This prevents "0. Inbox" from being trimmed to "0"
	stem := filename
	lower := strings.ToLower(filename)
	if strings.HasSuffix(lower, ".md") {
		stem = filename[:len(filename)-3]
	} else if strings.HasSuffix(lower, ".markdown") {
		stem = filename[:len(filename)-9]
	}
	meta := FileMetadata{
		OriginalName: filename,
		IsDraft:      strings.HasPrefix(stem, "_"),
	}

	// Remove draft prefix for further processing
	if meta.IsDraft {
		stem = strings.TrimPrefix(stem, "_")
	}

	remaining := stem

	// Try to extract date prefix (YYYY-MM-DD with various separators)
	if matches := datePrefixRegex.FindStringSubmatch(stem); len(matches) == 5 {
		year, _ := strconv.Atoi(matches[1])
		month, _ := strconv.Atoi(matches[2])
		day, _ := strconv.Atoi(matches[3])

		meta.Date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		meta.HasDate = true
		remaining = matches[4]
	}
	// Note: If no date prefix found, HasDate stays false and Date is zero time

	// Try to extract number prefix from remaining
	if matches := numberPrefixRegex.FindStringSubmatch(remaining); len(matches) == 3 {
		num, _ := strconv.Atoi(matches[1])
		meta.Number = &num
		remaining = matches[2]
	}

	meta.Slug = Slugify(remaining)
	meta.DisplayName = titleize(remaining)

	return meta
}

// titleize converts "hello-world" to "Hello World"
func titleize(slug string) string {
	// Split on hyphens and underscores
	parts := strings.FieldsFunc(slug, func(r rune) bool {
		return r == '-' || r == '_'
	})
	for i, part := range parts {
		if len(part) > 0 {
			// Capitalize first letter
			parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, " ")
}

// SortNodes sorts tree nodes by date, number, then name
// Sorting order: files first, then folders, each sorted by date/number/name
func SortNodes(nodes []*Node, newestFirst bool) {
	sort.Slice(nodes, func(i, j int) bool {
		a, b := nodes[i], nodes[j]

		// Files before folders
		if a.IsFolder != b.IsFolder {
			return !a.IsFolder
		}

		// For files, get metadata and sort by date/number/name
		if !a.IsFolder && !b.IsFolder {
			aMeta := GetNodeMetadata(a)
			bMeta := GetNodeMetadata(b)

			// Primary: Date (from filename only)
			// Files with dates come before files without dates
			if aMeta.HasDate != bMeta.HasDate {
				return aMeta.HasDate // files with dates first
			}
			// Both have dates - sort by date
			if aMeta.HasDate && bMeta.HasDate && !aMeta.Date.Equal(bMeta.Date) {
				if newestFirst {
					return aMeta.Date.After(bMeta.Date)
				}
				return aMeta.Date.Before(bMeta.Date)
			}

			// Secondary: Number (lower numbers first, nil sorted last)
			aNum := getNumberForSort(aMeta.Number, newestFirst)
			bNum := getNumberForSort(bMeta.Number, newestFirst)
			if aNum != bNum {
				if newestFirst {
					return aNum > bNum
				}
				return aNum < bNum
			}
		}

		// Tertiary: Name (alphabetical)
		if newestFirst {
			return strings.ToLower(a.Name) > strings.ToLower(b.Name)
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
}

// GetNodeMetadata extracts metadata from a node's filename
func GetNodeMetadata(node *Node) FileMetadata {
	var modTime time.Time
	if info, err := os.Stat(node.SourcePath); err == nil {
		modTime = info.ModTime()
	} else {
		modTime = time.Now()
	}

	filename := filepath.Base(node.SourcePath)
	return ExtractFileMetadata(filename, modTime)
}

// getNumberForSort returns the number for sorting, handling nil
func getNumberForSort(n *int, newestFirst bool) int {
	if n == nil {
		if newestFirst {
			return -1 // Sort after numbered items
		}
		return 999999 // Sort after numbered items
	}
	return *n
}

// IsDraftFile checks if a filename indicates a draft
func IsDraftFile(filename string) bool {
	return strings.HasPrefix(filename, "_")
}

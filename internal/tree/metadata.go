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
	// datePrefix matches YYYY-MM-DD- prefix
	datePrefixRegex = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})-(.+)$`)
	// numberPrefix matches leading digits like 01-
	numberPrefixRegex = regexp.MustCompile(`^(\d+)-(.+)$`)
)

// FileMetadata holds parsed metadata from a filename
type FileMetadata struct {
	OriginalName string    // "2024-01-15-01-hello-world.md"
	Date         time.Time // Parsed from filename or file mtime
	DateSource   string    // "filename" or "mtime"
	Number       *int      // nil if no number prefix
	Slug         string    // "hello-world"
	DisplayName  string    // "Hello World"
	IsDraft      bool      // starts with "_"
}

// ExtractFileMetadata parses date/number prefixes from filename
func ExtractFileMetadata(filename string, modTime time.Time) FileMetadata {
	stem := strings.TrimSuffix(filename, filepath.Ext(filename))
	meta := FileMetadata{
		OriginalName: filename,
		IsDraft:      strings.HasPrefix(stem, "_"),
	}

	// Remove draft prefix for further processing
	if meta.IsDraft {
		stem = strings.TrimPrefix(stem, "_")
	}

	remaining := stem

	// Try to extract date prefix (YYYY-MM-DD-)
	if matches := datePrefixRegex.FindStringSubmatch(stem); len(matches) == 5 {
		year, _ := strconv.Atoi(matches[1])
		month, _ := strconv.Atoi(matches[2])
		day, _ := strconv.Atoi(matches[3])

		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		meta.Date = date
		meta.DateSource = "filename"
		remaining = matches[4]
	}

	// If no date from filename, use file mtime
	if meta.DateSource == "" {
		meta.Date = modTime
		meta.DateSource = "mtime"
	}

	// Try to extract number prefix from remaining
	if matches := numberPrefixRegex.FindStringSubmatch(remaining); len(matches) == 3 {
		num, _ := strconv.Atoi(matches[1])
		meta.Number = &num
		remaining = matches[2]
	}

	meta.Slug = remaining
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
func SortNodes(nodes []*Node, newestFirst bool) {
	sort.Slice(nodes, func(i, j int) bool {
		a, b := nodes[i], nodes[j]

		// Folders always before files
		if a.IsFolder != b.IsFolder {
			return a.IsFolder
		}

		// For files, get metadata
		if !a.IsFolder && !b.IsFolder {
			aMeta := getNodeMetadata(a)
			bMeta := getNodeMetadata(b)

			// Primary: Date
			if !aMeta.Date.Equal(bMeta.Date) {
				if newestFirst {
					return aMeta.Date.After(bMeta.Date)
				}
				return aMeta.Date.Before(bMeta.Date)
			}

			// Secondary: Number (nil treated as infinity when newestFirst)
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

// getNodeMetadata extracts metadata from a node's filename
func getNodeMetadata(node *Node) FileMetadata {
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

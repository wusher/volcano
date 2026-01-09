package content

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

// ModifiedDate represents the last modification date of a file
type ModifiedDate struct {
	Time     time.Time
	Source   string // "git" or "filesystem"
	Absolute string // "January 5, 2025"
	Relative string // "3 days ago"
}

// GetLastModified retrieves the last modified date for a file.
// It first tries to get the date from git, falling back to filesystem.
func GetLastModified(filePath string) ModifiedDate {
	// Try git first
	gitDate, err := getGitModifiedDate(filePath)
	if err == nil {
		return gitDate
	}

	// Fallback to filesystem
	return getFileModifiedDate(filePath)
}

// getGitModifiedDate gets the last commit date for a file from git
func getGitModifiedDate(filePath string) (ModifiedDate, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%cI", "--", filePath)
	output, err := cmd.Output()
	if err != nil {
		return ModifiedDate{}, err
	}

	dateStr := strings.TrimSpace(string(output))
	if dateStr == "" {
		return ModifiedDate{}, os.ErrNotExist
	}

	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return ModifiedDate{}, err
	}

	return ModifiedDate{
		Time:     t,
		Source:   "git",
		Absolute: formatAbsoluteDate(t),
		Relative: formatRelativeDate(t),
	}, nil
}

// getFileModifiedDate gets the modification time from the filesystem
func getFileModifiedDate(filePath string) ModifiedDate {
	info, err := os.Stat(filePath)
	if err != nil {
		// If we can't get the file info, use current time
		now := time.Now()
		return ModifiedDate{
			Time:     now,
			Source:   "filesystem",
			Absolute: formatAbsoluteDate(now),
			Relative: formatRelativeDate(now),
		}
	}

	t := info.ModTime()
	return ModifiedDate{
		Time:     t,
		Source:   "filesystem",
		Absolute: formatAbsoluteDate(t),
		Relative: formatRelativeDate(t),
	}
}

// formatAbsoluteDate formats a time as "January 5, 2025"
func formatAbsoluteDate(t time.Time) string {
	return t.Format("January 2, 2006")
}

// formatRelativeDate formats a time as a relative duration (e.g., "3 days ago")
func formatRelativeDate(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return string(rune(mins+'0')) + " minutes ago"
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return formatNumber(hours) + " hours ago"
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return formatNumber(days) + " days ago"
	case diff < 30*24*time.Hour:
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return formatNumber(weeks) + " weeks ago"
	case diff < 365*24*time.Hour:
		months := int(diff.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return formatNumber(months) + " months ago"
	default:
		years := int(diff.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return formatNumber(years) + " years ago"
	}
}

// formatNumber converts an integer to a string
func formatNumber(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	// For larger numbers, use simple conversion
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}

// FormatLastModified returns a formatted string for display
func FormatLastModified(mod ModifiedDate, useRelative bool) string {
	if useRelative {
		return mod.Relative
	}
	return mod.Absolute
}

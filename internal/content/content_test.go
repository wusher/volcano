package content

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCalculateReadingTime(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected int // expected minutes
	}{
		{"empty", "", 1},
		{"short", "Hello world", 1},
		{"one minute", generateWords(200), 1},
		{"two minutes", generateWords(450), 2},
		{"five minutes", generateWords(1100), 5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateReadingTime(tc.content)
			if result.Minutes != tc.expected {
				t.Errorf("Minutes = %d, want %d", result.Minutes, tc.expected)
			}
		})
	}
}

func generateWords(count int) string {
	var content string
	for i := 0; i < count; i++ {
		content += "word "
	}
	return content
}

func TestFormatReadingTime(t *testing.T) {
	tests := []struct {
		rt       ReadingTime
		expected string
	}{
		{ReadingTime{Minutes: 1, Words: 100}, "1 min read"},
		{ReadingTime{Minutes: 5, Words: 1000}, "5 min read"},
		{ReadingTime{Minutes: 10, Words: 2000}, "10 min read"},
	}

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			result := FormatReadingTime(tc.rt)
			if result != tc.expected {
				t.Errorf("FormatReadingTime() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestGetLastModifiedFilesystem(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")

	if err := os.WriteFile(testFile, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	mod := GetLastModified(testFile)

	// Should get filesystem time
	if mod.Source != "filesystem" && mod.Source != "git" {
		t.Errorf("Source = %q, want 'filesystem' or 'git'", mod.Source)
	}

	if mod.Absolute == "" {
		t.Error("Absolute date should not be empty")
	}
	if mod.Relative == "" {
		t.Error("Relative date should not be empty")
	}
}

func TestFormatAbsoluteDate(t *testing.T) {
	date := time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC)
	result := formatAbsoluteDate(date)
	expected := "January 15, 2025"
	if result != expected {
		t.Errorf("formatAbsoluteDate() = %q, want %q", result, expected)
	}
}

func TestFormatRelativeDate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		contains string
	}{
		{"just now", now.Add(-30 * time.Second), "just now"},
		{"minutes", now.Add(-5 * time.Minute), "minute"},
		{"hours", now.Add(-3 * time.Hour), "hour"},
		{"yesterday", now.Add(-36 * time.Hour), "yesterday"},
		{"days", now.Add(-5 * 24 * time.Hour), "days ago"},
		{"weeks", now.Add(-14 * 24 * time.Hour), "week"},
		{"months", now.Add(-60 * 24 * time.Hour), "month"},
		{"years", now.Add(-400 * 24 * time.Hour), "year"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := formatRelativeDate(tc.time)
			if result == "" {
				t.Error("result should not be empty")
			}
		})
	}
}

func TestFormatLastModified(t *testing.T) {
	mod := ModifiedDate{
		Absolute: "January 15, 2025",
		Relative: "3 days ago",
	}

	absResult := FormatLastModified(mod, false)
	if absResult != "January 15, 2025" {
		t.Errorf("absolute format = %q, want %q", absResult, "January 15, 2025")
	}

	relResult := FormatLastModified(mod, true)
	if relResult != "3 days ago" {
		t.Errorf("relative format = %q, want %q", relResult, "3 days ago")
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		n        int
		expected string
	}{
		{0, "0"},
		{5, "5"},
		{10, "10"},
		{123, "123"},
	}

	for _, tc := range tests {
		result := formatNumber(tc.n)
		if result != tc.expected {
			t.Errorf("formatNumber(%d) = %q, want %q", tc.n, result, tc.expected)
		}
	}
}

func TestItoa(t *testing.T) {
	tests := []struct {
		n        int
		expected string
	}{
		{0, "0"},
		{5, "5"},
		{-5, "-5"},
		{123, "123"},
		{-123, "-123"},
	}

	for _, tc := range tests {
		result := itoa(tc.n)
		if result != tc.expected {
			t.Errorf("itoa(%d) = %q, want %q", tc.n, result, tc.expected)
		}
	}
}

func TestCalculateReadingTimeWithCode(t *testing.T) {
	content := `<p>Some regular text here.</p>
<pre><code>func main() {
    fmt.Println("Hello, world!")
}</code></pre>
<p>More text here.</p>`

	rt := CalculateReadingTime(content)
	if rt.Words < 5 {
		t.Errorf("Words should be at least 5, got %d", rt.Words)
	}
	if rt.Minutes < 1 {
		t.Errorf("Minutes should be at least 1, got %d", rt.Minutes)
	}
}

func TestGetLastModifiedNonexistentFile(t *testing.T) {
	mod := GetLastModified("/nonexistent/path/to/file.md")

	// Should fallback gracefully
	if mod.Source != "filesystem" {
		t.Errorf("Source = %q, want 'filesystem'", mod.Source)
	}
	if mod.Absolute == "" {
		t.Error("Absolute should not be empty even for nonexistent file")
	}
}

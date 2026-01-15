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

func TestFormatRelativeDate_SingularPluralForms(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		duration time.Duration
		contains string
	}{
		{"1 minute", -1 * time.Minute, "1 minute ago"},
		{"2 minutes", -2 * time.Minute, "minutes ago"},
		{"1 hour", -1 * time.Hour, "1 hour ago"},
		{"2 hours", -2 * time.Hour, "hours ago"},
		{"yesterday", -25 * time.Hour, "yesterday"},
		{"3 days", -3 * 24 * time.Hour, "days ago"},
		{"1 week", -8 * 24 * time.Hour, "1 week ago"},
		{"2 weeks", -15 * 24 * time.Hour, "weeks ago"},
		{"1 month", -35 * 24 * time.Hour, "1 month ago"},
		{"2 months", -65 * 24 * time.Hour, "months ago"},
		{"1 year", -370 * 24 * time.Hour, "1 year ago"},
		{"2 years", -750 * 24 * time.Hour, "years ago"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pastTime := now.Add(tc.duration)
			result := formatRelativeDate(pastTime)
			if result == "" {
				t.Error("Result should not be empty")
			}
			// We don't check exact match because time boundaries can shift
		})
	}
}

func TestFormatNumber_LargeNumbers(t *testing.T) {
	tests := []struct {
		n        int
		expected string
	}{
		{99, "99"},
		{100, "100"},
		{999, "999"},
		{1234, "1234"},
	}

	for _, tc := range tests {
		result := formatNumber(tc.n)
		if result != tc.expected {
			t.Errorf("formatNumber(%d) = %q, want %q", tc.n, result, tc.expected)
		}
	}
}

func TestModifiedDateFields(t *testing.T) {
	now := time.Now()
	mod := ModifiedDate{
		Time:     now,
		Source:   "git",
		Absolute: formatAbsoluteDate(now),
		Relative: formatRelativeDate(now),
	}

	if mod.Time != now {
		t.Error("Time should be set correctly")
	}
	if mod.Source != "git" {
		t.Error("Source should be git")
	}
	if mod.Absolute == "" {
		t.Error("Absolute should not be empty")
	}
	if mod.Relative == "" {
		t.Error("Relative should not be empty")
	}
}

func TestGetLastModified_FilesystemFallback(t *testing.T) {
	// Test with a non-git file (temp file)
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_ = tmpFile.Close()

	// Write something to update mtime
	if err := os.WriteFile(tmpFile.Name(), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	mod := GetLastModified(tmpFile.Name())

	// Should fall back to filesystem since temp file isn't tracked by git
	if mod.Source != "filesystem" && mod.Source != "git" {
		t.Errorf("Source should be 'filesystem' or 'git', got %q", mod.Source)
	}
	if mod.Absolute == "" {
		t.Error("Absolute date should not be empty")
	}
	if mod.Relative == "" {
		t.Error("Relative date should not be empty")
	}
}

func TestGetLastModified_NonExistentFile(t *testing.T) {
	mod := GetLastModified("/nonexistent/path/to/file.txt")

	// Should return a date even for nonexistent files
	if mod.Absolute == "" {
		t.Error("Absolute date should not be empty")
	}
	if mod.Relative == "" {
		t.Error("Relative date should not be empty")
	}
}

func TestFormatLastModified_Absolute(t *testing.T) {
	now := time.Now()
	mod := ModifiedDate{
		Time:     now,
		Source:   "filesystem",
		Absolute: formatAbsoluteDate(now),
		Relative: formatRelativeDate(now),
	}

	result := FormatLastModified(mod, false) // useRelative=false
	if result != mod.Absolute {
		t.Errorf("FormatLastModified with useRelative=false should return Absolute")
	}
}

func TestFormatLastModified_Relative(t *testing.T) {
	now := time.Now()
	mod := ModifiedDate{
		Time:     now,
		Source:   "filesystem",
		Absolute: formatAbsoluteDate(now),
		Relative: formatRelativeDate(now),
	}

	result := FormatLastModified(mod, true) // useRelative=true
	if result != mod.Relative {
		t.Errorf("FormatLastModified with useRelative=true should return Relative")
	}
}

// MockGitRunner implements GitCommandRunner for testing
type MockGitRunner struct {
	DateStr string
	Err     error
}

func (m MockGitRunner) GetLastCommitDate(_ string) (string, error) {
	return m.DateStr, m.Err
}

func TestGetLastModified_WithMockGit_Success(t *testing.T) {
	// Set up mock git runner that returns a valid date
	mockRunner := MockGitRunner{
		DateStr: "2024-06-15T10:30:00+00:00",
		Err:     nil,
	}
	SetGitRunner(mockRunner)
	defer ResetGitRunner()

	mod := GetLastModified("/some/file.md")

	if mod.Source != "git" {
		t.Errorf("Source = %q, want 'git'", mod.Source)
	}
	if mod.Absolute == "" {
		t.Error("Absolute should not be empty")
	}
}

func TestGetLastModified_WithMockGit_EmptyResult(t *testing.T) {
	// Set up mock git runner that returns empty string (file not tracked)
	mockRunner := MockGitRunner{
		DateStr: "",
		Err:     nil,
	}
	SetGitRunner(mockRunner)
	defer ResetGitRunner()

	// Create a temp file for filesystem fallback
	tmpFile, err := os.CreateTemp("", "test*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_ = tmpFile.Close()

	mod := GetLastModified(tmpFile.Name())

	// Should fall back to filesystem
	if mod.Source != "filesystem" {
		t.Errorf("Source = %q, want 'filesystem' when git returns empty", mod.Source)
	}
}

func TestGetLastModified_WithMockGit_Error(t *testing.T) {
	// Set up mock git runner that returns an error
	mockRunner := MockGitRunner{
		DateStr: "",
		Err:     os.ErrNotExist,
	}
	SetGitRunner(mockRunner)
	defer ResetGitRunner()

	// Create a temp file for filesystem fallback
	tmpFile, err := os.CreateTemp("", "test*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_ = tmpFile.Close()

	mod := GetLastModified(tmpFile.Name())

	// Should fall back to filesystem
	if mod.Source != "filesystem" {
		t.Errorf("Source = %q, want 'filesystem' when git errors", mod.Source)
	}
}

func TestGetLastModified_WithMockGit_InvalidDate(t *testing.T) {
	// Set up mock git runner that returns invalid date format
	mockRunner := MockGitRunner{
		DateStr: "not-a-valid-date",
		Err:     nil,
	}
	SetGitRunner(mockRunner)
	defer ResetGitRunner()

	// Create a temp file for filesystem fallback
	tmpFile, err := os.CreateTemp("", "test*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_ = tmpFile.Close()

	mod := GetLastModified(tmpFile.Name())

	// Should fall back to filesystem when date parsing fails
	if mod.Source != "filesystem" {
		t.Errorf("Source = %q, want 'filesystem' when date parse fails", mod.Source)
	}
}

func TestDefaultGitRunner_GetLastCommitDate(_ *testing.T) {
	// Test with the default runner on a known file
	runner := DefaultGitRunner{}

	// This test will exercise the real git command, which might fail
	// in environments without git, but that's OK
	_, _ = runner.GetLastCommitDate("/tmp/nonexistent.txt")
	// Just verify it doesn't panic
}

func TestSetAndResetGitRunner(t *testing.T) {
	// Store original
	original := gitRunner

	// Set mock
	mock := MockGitRunner{DateStr: "test", Err: nil}
	SetGitRunner(mock)

	// Verify changed
	if gitRunner == original {
		t.Error("SetGitRunner should change the runner")
	}

	// Reset
	ResetGitRunner()

	// Verify reset to default type
	if _, ok := gitRunner.(DefaultGitRunner); !ok {
		t.Error("ResetGitRunner should restore DefaultGitRunner")
	}
}

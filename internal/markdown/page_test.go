package markdown

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	// Create a temporary markdown file
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "test.md")

	content := `# Test Page

This is a test paragraph.

## Section

- Item 1
- Item 2

` + "```go\nfunc main() {}\n```"

	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	page, err := ParseFile(mdFile, "test/index.html", "/test/", "Fallback Title")
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	// Check title extraction
	if page.Title != "Test Page" {
		t.Errorf("Title = %q, want %q", page.Title, "Test Page")
	}

	// Check HTML content
	if !strings.Contains(page.Content, "<h1") {
		t.Error("Content should contain h1 tag")
	}
	if !strings.Contains(page.Content, "<ul>") {
		t.Error("Content should contain ul tag")
	}
	if !strings.Contains(page.Content, "<pre") {
		t.Error("Content should contain pre tag")
	}

	// Check paths
	if page.SourcePath != mdFile {
		t.Errorf("SourcePath = %q, want %q", page.SourcePath, mdFile)
	}
	if page.OutputPath != "test/index.html" {
		t.Errorf("OutputPath = %q, want %q", page.OutputPath, "test/index.html")
	}
	if page.URLPath != "/test/" {
		t.Errorf("URLPath = %q, want %q", page.URLPath, "/test/")
	}
}

func TestParseFileFallbackTitle(t *testing.T) {
	// Create a temporary markdown file without h1
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "no-title.md")

	content := `This is content without a heading.

Just some paragraphs and text.`

	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	page, err := ParseFile(mdFile, "output.html", "/path/", "Fallback Title")
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	// Should use fallback title
	if page.Title != "Fallback Title" {
		t.Errorf("Title = %q, want %q (fallback)", page.Title, "Fallback Title")
	}
}

func TestParseFileNonExistent(t *testing.T) {
	_, err := ParseFile("/nonexistent/file.md", "output.html", "/path/", "Title")
	if err == nil {
		t.Error("ParseFile() should return error for non-existent file")
	}
}

func TestCleanFilenameTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"getting-started.md", "Getting Started"},
		{"api_reference.md", "Api Reference"},
		{"simple.md", "Simple"},
		{"01-introduction.md", "Introduction"},
		{"001_setup.md", "Setup"},
		{"2024-01-01-my-post.md", "My Post"},
		{"hello-world.md", "Hello World"},
		{"no-extension", "No Extension"},
		{"multiple---dashes.md", "Multiple   Dashes"},
		{"under__scores.md", "Under  Scores"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := CleanFilenameTitle(tt.input)
			if result != tt.expected {
				t.Errorf("CleanFilenameTitle(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRemoveLeadingPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"01-intro", "intro"},
		{"001_setup", "setup"},
		{"2024-01-01-post", "post"},
		{"no-numbers", "no-numbers"},
		{"123", "123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := removeLeadingPrefix(tt.input)
			if result != tt.expected {
				t.Errorf("removeLeadingPrefix(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsDatePrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"2024-01-01-", true},
		{"2023-12-31-", true},
		{"1999-05-15-", true},
		{"2024-1-01-", false},
		{"2024-01-1-", false},
		{"2024/01/01-", false},
		{"2024-01-01", false},
		{"short", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isDatePrefix(tt.input)
			if result != tt.expected {
				t.Errorf("isDatePrefix(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

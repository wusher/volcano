package search

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateSearchIndex(t *testing.T) {
	tmpDir := t.TempDir()

	index := &Index{
		Pages: []PageEntry{
			{
				Title: "Home",
				URL:   "/",
				Headings: []HeadingEntry{
					{Text: "Welcome", Anchor: "welcome", Level: 2},
					{Text: "Getting Started", Anchor: "getting-started", Level: 3},
				},
			},
			{
				Title: "About",
				URL:   "/about/",
				Headings: []HeadingEntry{
					{Text: "Our Team", Anchor: "our-team", Level: 2},
				},
			},
		},
	}

	err := GenerateSearchIndex(tmpDir, index)
	if err != nil {
		t.Fatalf("GenerateSearchIndex() error = %v", err)
	}

	// Check file was created
	indexPath := filepath.Join(tmpDir, "search-index.json")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read search-index.json: %v", err)
	}

	// Verify content
	contentStr := string(content)
	if contentStr == "" {
		t.Error("search-index.json should not be empty")
	}

	// Check JSON structure
	if !contains(contentStr, "Home") {
		t.Error("search-index.json should contain page title 'Home'")
	}
	if !contains(contentStr, "/about/") {
		t.Error("search-index.json should contain URL '/about/'")
	}
	if !contains(contentStr, "Welcome") {
		t.Error("search-index.json should contain heading 'Welcome'")
	}
	if !contains(contentStr, "getting-started") {
		t.Error("search-index.json should contain anchor 'getting-started'")
	}
}

func TestGenerateSearchIndex_EmptyIndex(t *testing.T) {
	tmpDir := t.TempDir()

	index := &Index{
		Pages: []PageEntry{},
	}

	err := GenerateSearchIndex(tmpDir, index)
	if err != nil {
		t.Fatalf("GenerateSearchIndex() error = %v", err)
	}

	// Check file was created
	indexPath := filepath.Join(tmpDir, "search-index.json")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read search-index.json: %v", err)
	}

	// Verify it's valid JSON with empty pages array
	if !contains(string(content), `"pages": []`) {
		t.Error("Empty index should produce JSON with empty pages array")
	}
}

func TestGenerateSearchIndex_InvalidDir(t *testing.T) {
	// Try to write to a non-existent directory
	err := GenerateSearchIndex("/nonexistent/path/does/not/exist", &Index{Pages: []PageEntry{}})
	if err == nil {
		t.Error("GenerateSearchIndex() should return error for invalid directory")
	}
}

func TestGenerateSearchIndex_WithNilHeadings(t *testing.T) {
	tmpDir := t.TempDir()

	index := &Index{
		Pages: []PageEntry{
			{
				Title:    "Page Without Headings",
				URL:      "/no-headings/",
				Headings: nil,
			},
		},
	}

	err := GenerateSearchIndex(tmpDir, index)
	if err != nil {
		t.Fatalf("GenerateSearchIndex() error = %v", err)
	}

	// Check file was created and is valid
	indexPath := filepath.Join(tmpDir, "search-index.json")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read search-index.json: %v", err)
	}

	if !contains(string(content), "Page Without Headings") {
		t.Error("search-index.json should contain page title")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

package tree

import (
	"testing"
	"time"
)

func TestExtractFileMetadata(t *testing.T) {
	modTime := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		filename    string
		wantDate    string
		wantSlug    string
		wantNumber  *int
		wantDateSrc string
	}{
		{
			name:        "simple filename",
			filename:    "hello-world.md",
			wantSlug:    "hello-world",
			wantDateSrc: "mtime",
		},
		{
			name:        "date prefix",
			filename:    "2024-01-15-hello-world.md",
			wantDate:    "2024-01-15",
			wantSlug:    "hello-world",
			wantDateSrc: "filename",
		},
		{
			name:        "number prefix",
			filename:    "01-getting-started.md",
			wantSlug:    "getting-started",
			wantNumber:  intPtr(1),
			wantDateSrc: "mtime",
		},
		{
			name:        "date and number prefix",
			filename:    "2024-01-15-01-intro.md",
			wantDate:    "2024-01-15",
			wantSlug:    "intro", // number is also extracted
			wantNumber:  intPtr(1),
			wantDateSrc: "filename",
		},
		{
			name:        "draft prefix",
			filename:    "_draft-post.md",
			wantSlug:    "draft-post",
			wantDateSrc: "mtime",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			meta := ExtractFileMetadata(tc.filename, modTime)

			if meta.Slug != tc.wantSlug {
				t.Errorf("Slug = %q, want %q", meta.Slug, tc.wantSlug)
			}

			if tc.wantDateSrc != "" && meta.DateSource != tc.wantDateSrc {
				t.Errorf("DateSource = %q, want %q", meta.DateSource, tc.wantDateSrc)
			}

			if tc.wantNumber != nil {
				if meta.Number == nil || *meta.Number != *tc.wantNumber {
					t.Errorf("Number = %v, want %v", meta.Number, *tc.wantNumber)
				}
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}

func TestTitleize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello-world", "Hello World"},
		{"getting_started", "Getting Started"},
		{"api-reference", "Api Reference"},
		{"FAQ", "Faq"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := titleize(tc.input)
			if result != tc.expected {
				t.Errorf("titleize(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsDraftFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"_draft.md", true},
		{"_unpublished-post.md", true},
		{"regular-post.md", false},
		{"draft.md", false},
	}

	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			result := IsDraftFile(tc.filename)
			if result != tc.expected {
				t.Errorf("IsDraftFile(%q) = %v, want %v", tc.filename, result, tc.expected)
			}
		})
	}
}

func TestSortNodes(t *testing.T) {
	// Create nodes with different dates
	node1 := &Node{Name: "Old Post", Path: "2024-01-01-old.md", IsFolder: false, SourcePath: ""}
	node2 := &Node{Name: "New Post", Path: "2024-03-15-new.md", IsFolder: false, SourcePath: ""}
	node3 := &Node{Name: "Folder", Path: "folder", IsFolder: true}

	nodes := []*Node{node1, node2, node3}

	// Sort newest first
	SortNodes(nodes, true)

	// Folders should come first
	if !nodes[0].IsFolder {
		t.Error("Folder should be first after sorting")
	}
}

func TestGetNumberForSort(t *testing.T) {
	num := 5

	// With number, newest first
	result := getNumberForSort(&num, true)
	if result != 5 {
		t.Errorf("getNumberForSort with number = %d, want 5", result)
	}

	// Without number, newest first
	result = getNumberForSort(nil, true)
	if result != -1 {
		t.Errorf("getNumberForSort nil newest first = %d, want -1", result)
	}

	// Without number, oldest first
	result = getNumberForSort(nil, false)
	if result != 999999 {
		t.Errorf("getNumberForSort nil oldest first = %d, want 999999", result)
	}
}

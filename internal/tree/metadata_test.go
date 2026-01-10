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
		wantHasDate bool
	}{
		{
			name:        "simple filename",
			filename:    "hello-world.md",
			wantSlug:    "hello-world",
			wantHasDate: false,
		},
		{
			name:        "date prefix with dash",
			filename:    "2024-01-15-hello-world.md",
			wantDate:    "2024-01-15",
			wantSlug:    "hello-world",
			wantHasDate: true,
		},
		{
			name:        "date prefix with underscore",
			filename:    "2024_01_15_hello-world.md",
			wantDate:    "2024-01-15",
			wantSlug:    "hello-world",
			wantHasDate: true,
		},
		{
			name:        "date prefix with space",
			filename:    "2024-01-15 hello-world.md",
			wantDate:    "2024-01-15",
			wantSlug:    "hello-world",
			wantHasDate: true,
		},
		{
			name:        "number prefix with dash",
			filename:    "01-getting-started.md",
			wantSlug:    "getting-started",
			wantNumber:  intPtr(1),
			wantHasDate: false,
		},
		{
			name:        "number prefix with underscore",
			filename:    "01_getting-started.md",
			wantSlug:    "getting-started",
			wantNumber:  intPtr(1),
			wantHasDate: false,
		},
		{
			name:        "number prefix with space",
			filename:    "01 getting-started.md",
			wantSlug:    "getting-started",
			wantNumber:  intPtr(1),
			wantHasDate: false,
		},
		{
			name:        "date and number prefix",
			filename:    "2024-01-15-01-intro.md",
			wantDate:    "2024-01-15",
			wantSlug:    "intro",
			wantNumber:  intPtr(1),
			wantHasDate: true,
		},
		{
			name:        "draft prefix",
			filename:    "_draft-post.md",
			wantSlug:    "draft-post",
			wantHasDate: false,
		},
		{
			name:        "number prefix with dot (folder style)",
			filename:    "0. Inbox",
			wantSlug:    "inbox",
			wantNumber:  intPtr(0),
			wantHasDate: false,
		},
		{
			name:        "number prefix with dot and space",
			filename:    "1. Projects",
			wantSlug:    "projects",
			wantNumber:  intPtr(1),
			wantHasDate: false,
		},
		{
			name:        "number prefix with dot no space",
			filename:    "8.Archive",
			wantSlug:    "archive",
			wantNumber:  intPtr(8),
			wantHasDate: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			meta := ExtractFileMetadata(tc.filename, modTime)

			if meta.Slug != tc.wantSlug {
				t.Errorf("Slug = %q, want %q", meta.Slug, tc.wantSlug)
			}

			if meta.HasDate != tc.wantHasDate {
				t.Errorf("HasDate = %v, want %v", meta.HasDate, tc.wantHasDate)
			}

			if tc.wantHasDate && tc.wantDate != "" {
				gotDate := meta.Date.Format("2006-01-02")
				if gotDate != tc.wantDate {
					t.Errorf("Date = %q, want %q", gotDate, tc.wantDate)
				}
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

	// Files should come first, then folders
	if nodes[0].IsFolder {
		t.Error("Files should come before folders after sorting")
	}
	// Last item should be the folder
	if !nodes[2].IsFolder {
		t.Error("Folder should be last after sorting")
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

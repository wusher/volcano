package tree

import (
	"os"
	"testing"
)

func TestExtractFileMetadata(t *testing.T) {
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
			meta := ExtractFileMetadata(tc.filename)

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
	t.Run("files before folders", func(t *testing.T) {
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
	})

	t.Run("files sorted by date", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create test files with date prefixes
		oldFile := tmpDir + "/2024-01-01-old.md"
		newFile := tmpDir + "/2024-03-15-new.md"
		if err := os.WriteFile(oldFile, []byte("old"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(newFile, []byte("new"), 0644); err != nil {
			t.Fatal(err)
		}

		node1 := &Node{Name: "Old Post", FileName: "2024-01-01-old.md", IsFolder: false, SourcePath: oldFile}
		node2 := &Node{Name: "New Post", FileName: "2024-03-15-new.md", IsFolder: false, SourcePath: newFile}

		nodes := []*Node{node1, node2}

		// Sort newest first
		SortNodes(nodes, true)
		if nodes[0].FileName != "2024-03-15-new.md" {
			t.Error("Newer file should come first when sorting newest first")
		}

		// Sort oldest first
		nodes = []*Node{node1, node2}
		SortNodes(nodes, false)
		if nodes[0].FileName != "2024-01-01-old.md" {
			t.Error("Older file should come first when sorting oldest first")
		}
	})

	t.Run("files sorted by number prefix", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create test files with number prefixes
		file1 := tmpDir + "/01-intro.md"
		file2 := tmpDir + "/02-setup.md"
		if err := os.WriteFile(file1, []byte("intro"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file2, []byte("setup"), 0644); err != nil {
			t.Fatal(err)
		}

		node1 := &Node{Name: "Intro", FileName: "01-intro.md", IsFolder: false, SourcePath: file1}
		node2 := &Node{Name: "Setup", FileName: "02-setup.md", IsFolder: false, SourcePath: file2}

		nodes := []*Node{node2, node1}

		// Sort oldest first (lower numbers first)
		SortNodes(nodes, false)
		if nodes[0].FileName != "01-intro.md" {
			t.Error("Lower numbered file should come first when sorting oldest first")
		}
	})

	t.Run("files sorted alphabetically by name", func(t *testing.T) {
		tmpDir := t.TempDir()

		fileA := tmpDir + "/alpha.md"
		fileB := tmpDir + "/beta.md"
		if err := os.WriteFile(fileA, []byte("a"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fileB, []byte("b"), 0644); err != nil {
			t.Fatal(err)
		}

		node1 := &Node{Name: "Alpha", FileName: "alpha.md", IsFolder: false, SourcePath: fileA}
		node2 := &Node{Name: "Beta", FileName: "beta.md", IsFolder: false, SourcePath: fileB}

		nodes := []*Node{node2, node1}

		// Sort oldest first (alphabetical A-Z)
		SortNodes(nodes, false)
		if nodes[0].Name != "Alpha" {
			t.Error("Alpha should come first when sorting oldest first (A-Z)")
		}

		// Sort newest first (alphabetical Z-A)
		nodes = []*Node{node2, node1}
		SortNodes(nodes, true)
		if nodes[0].Name != "Beta" {
			t.Error("Beta should come first when sorting newest first (Z-A)")
		}
	})

	t.Run("dated files before undated files", func(t *testing.T) {
		tmpDir := t.TempDir()

		datedFile := tmpDir + "/2024-01-15-dated.md"
		undatedFile := tmpDir + "/undated.md"
		if err := os.WriteFile(datedFile, []byte("dated"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(undatedFile, []byte("undated"), 0644); err != nil {
			t.Fatal(err)
		}

		node1 := &Node{Name: "Dated", FileName: "2024-01-15-dated.md", IsFolder: false, SourcePath: datedFile}
		node2 := &Node{Name: "Undated", FileName: "undated.md", IsFolder: false, SourcePath: undatedFile}

		nodes := []*Node{node2, node1}

		SortNodes(nodes, false)
		if nodes[0].FileName != "2024-01-15-dated.md" {
			t.Error("Dated files should come before undated files")
		}
	})
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

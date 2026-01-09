package tree

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScan(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create test structure:
	// tmpDir/
	//   index.md
	//   about.md
	//   guides/
	//     index.md
	//     getting-started.md
	//   api/
	//     endpoints.md
	//   .hidden/
	//     secret.md
	//   empty/

	files := map[string]string{
		"index.md":                  "# Home",
		"about.md":                  "# About",
		"guides/index.md":           "# Guides",
		"guides/getting-started.md": "# Getting Started",
		"api/endpoints.md":          "# Endpoints",
		".hidden/secret.md":         "# Secret",
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create empty directory
	if err := os.MkdirAll(filepath.Join(tmpDir, "empty"), 0755); err != nil {
		t.Fatal(err)
	}

	// Scan the directory
	tree, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// Verify root
	if tree.Root == nil {
		t.Fatal("Root should not be nil")
	}

	// Verify AllPages count (should be 5: index, about, guides/index, getting-started, endpoints)
	if len(tree.AllPages) != 5 {
		t.Errorf("AllPages length = %d, want 5", len(tree.AllPages))
	}

	// Verify root children (should be: api, guides, About, Index - sorted with folders first)
	if len(tree.Root.Children) != 4 {
		t.Errorf("Root children = %d, want 4", len(tree.Root.Children))
	}

	// Verify folders come before files
	if len(tree.Root.Children) >= 2 {
		if !tree.Root.Children[0].IsFolder {
			t.Error("First child should be a folder")
		}
	}

	// Verify hidden folder is not included
	for _, child := range tree.Root.Children {
		if child.Name == "Hidden" || child.Path == ".hidden" {
			t.Error("Hidden folder should not be included")
		}
	}

	// Verify empty folder is not included
	for _, child := range tree.Root.Children {
		if child.Name == "Empty" || child.Path == "empty" {
			t.Error("Empty folder should not be included")
		}
	}

	// Find guides folder and verify it has HasIndex set
	var guidesFolder *Node
	for _, child := range tree.Root.Children {
		if child.Name == "Guides" {
			guidesFolder = child
			break
		}
	}
	if guidesFolder == nil {
		t.Fatal("Guides folder not found")
	}
	if !guidesFolder.HasIndex {
		t.Error("Guides folder should have HasIndex = true")
	}
	if guidesFolder.IndexPath != filepath.Join("guides", "index.md") {
		t.Errorf("Guides IndexPath = %q, want %q", guidesFolder.IndexPath, filepath.Join("guides", "index.md"))
	}
}

func TestScanNonExistentDirectory(t *testing.T) {
	_, err := Scan("/nonexistent/directory/path")
	if err == nil {
		t.Error("Scan should return error for non-existent directory")
	}
}

func TestGetOutputPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		isFolder bool
		expected string
	}{
		{
			name:     "regular file",
			path:     "guides/intro.md",
			isFolder: false,
			expected: filepath.Join("guides", "intro", "index.html"),
		},
		{
			name:     "index file",
			path:     "index.md",
			isFolder: false,
			expected: "index.html",
		},
		{
			name:     "nested index file",
			path:     "guides/index.md",
			isFolder: false,
			expected: filepath.Join("guides", "index.html"),
		},
		{
			name:     "readme file",
			path:     "readme.md",
			isFolder: false,
			expected: "index.html", // README treated as index
		},
		{
			name:     "folder",
			path:     "guides",
			isFolder: true,
			expected: "",
		},
		{
			name:     "date prefix file",
			path:     "posts/2024-01-15-hello-world.md",
			isFolder: false,
			expected: filepath.Join("posts", "hello-world", "index.html"),
		},
		{
			name:     "number prefix file",
			path:     "docs/01-getting-started.md",
			isFolder: false,
			expected: filepath.Join("docs", "getting-started", "index.html"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewNode("Test", tt.path, tt.isFolder)
			result := GetOutputPath(node)
			if result != tt.expected {
				t.Errorf("GetOutputPath() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetURLPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		isFolder bool
		expected string
	}{
		{
			name:     "regular file",
			path:     "guides/intro.md",
			isFolder: false,
			expected: "/guides/intro/",
		},
		{
			name:     "root index file",
			path:     "index.md",
			isFolder: false,
			expected: "/",
		},
		{
			name:     "nested index file",
			path:     "guides/index.md",
			isFolder: false,
			expected: "/guides/",
		},
		{
			name:     "readme file",
			path:     "docs/readme.md",
			isFolder: false,
			expected: "/docs/",
		},
		{
			name:     "folder",
			path:     "guides",
			isFolder: true,
			expected: "",
		},
		{
			name:     "date prefix file",
			path:     "posts/2024-01-15-hello-world.md",
			isFolder: false,
			expected: "/posts/hello-world/",
		},
		{
			name:     "number prefix file",
			path:     "docs/01-getting-started.md",
			isFolder: false,
			expected: "/docs/getting-started/",
		},
		{
			name:     "root date prefix file",
			path:     "2024-01-01-news.md",
			isFolder: false,
			expected: "/news/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewNode("Test", tt.path, tt.isFolder)
			result := GetURLPath(node)
			if result != tt.expected {
				t.Errorf("GetURLPath() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSortAndPrune(t *testing.T) {
	// Create a root with unsorted children
	root := NewNode("Root", "", true)

	// Add children in wrong order (files before folders, unsorted)
	fileZ := NewNode("Zebra", "zebra.md", false)
	fileA := NewNode("Alpha", "alpha.md", false)
	folderM := NewNode("Middle", "middle", true)
	folderA := NewNode("Apple", "apple", true)

	// Add a file to folders so they're not empty
	folderM.AddChild(NewNode("File", "middle/file.md", false))
	folderA.AddChild(NewNode("File", "apple/file.md", false))

	root.AddChild(fileZ)
	root.AddChild(fileA)
	root.AddChild(folderM)
	root.AddChild(folderA)

	// Also add an empty folder that should be pruned
	emptyFolder := NewNode("Empty", "empty", true)
	root.AddChild(emptyFolder)

	sortAndPrune(root)

	// Verify order: folders first (Apple, Middle), then files (Alpha, Zebra)
	// Empty folder should be removed
	if len(root.Children) != 4 {
		t.Errorf("Children count = %d, want 4 (empty folder should be pruned)", len(root.Children))
	}

	expectedOrder := []string{"Apple", "Middle", "Alpha", "Zebra"}
	for i, expected := range expectedOrder {
		if i >= len(root.Children) {
			break
		}
		if root.Children[i].Name != expected {
			t.Errorf("Child[%d].Name = %q, want %q", i, root.Children[i].Name, expected)
		}
	}
}

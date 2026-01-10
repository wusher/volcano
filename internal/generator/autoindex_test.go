package generator

import (
	"strings"
	"testing"

	"github.com/wusher/volcano/internal/tree"
)

func TestNeedsAutoIndex(t *testing.T) {
	tests := []struct {
		name     string
		node     *tree.Node
		expected bool
	}{
		{
			name:     "file node",
			node:     &tree.Node{IsFolder: false},
			expected: false,
		},
		{
			name:     "folder with index",
			node:     &tree.Node{IsFolder: true, HasIndex: true},
			expected: false,
		},
		{
			name: "folder with index.md child",
			node: &tree.Node{
				IsFolder: true,
				Children: []*tree.Node{
					{Path: "folder/index.md", IsFolder: false},
				},
			},
			expected: false,
		},
		{
			name: "folder with readme.md child",
			node: &tree.Node{
				IsFolder: true,
				Children: []*tree.Node{
					{Path: "folder/README.md", IsFolder: false},
				},
			},
			expected: false,
		},
		{
			name: "folder without index",
			node: &tree.Node{
				IsFolder: true,
				Children: []*tree.Node{
					{Path: "folder/page.md", IsFolder: false},
				},
			},
			expected: true,
		},
		{
			name:     "empty folder",
			node:     &tree.Node{IsFolder: true, Children: []*tree.Node{}},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NeedsAutoIndex(tc.node)
			if result != tc.expected {
				t.Errorf("NeedsAutoIndex() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestBuildAutoIndex(t *testing.T) {
	folder := &tree.Node{
		Name:     "Guides",
		Path:     "guides",
		IsFolder: true,
		Children: []*tree.Node{
			{Name: "Introduction", Path: "guides/intro.md", IsFolder: false},
			{Name: "Subguides", Path: "guides/sub", IsFolder: true},
			{Name: "Advanced", Path: "guides/advanced.md", IsFolder: false},
		},
	}

	index := BuildAutoIndex(folder)

	if index.Title != "Guides" {
		t.Errorf("Title = %q, want %q", index.Title, "Guides")
	}

	if index.URLPath != "/guides/" {
		t.Errorf("URLPath = %q, want %q", index.URLPath, "/guides/")
	}

	// Files should come first
	if len(index.Children) != 3 {
		t.Fatalf("expected 3 children, got %d", len(index.Children))
	}

	if index.Children[0].IsFolder {
		t.Error("first child should be file, not folder (files sorted first)")
	}
}

func TestBuildAutoIndexRootPath(t *testing.T) {
	folder := &tree.Node{
		Name:     "Root",
		Path:     "",
		IsFolder: true,
	}

	index := BuildAutoIndex(folder)
	if index.URLPath != "/" {
		t.Errorf("URLPath for root = %q, want %q", index.URLPath, "/")
	}
}

func TestRenderAutoIndexContent(t *testing.T) {
	index := AutoIndex{
		Title: "Guides",
		Children: []IndexItem{
			{Title: "Subfolder", URL: "/guides/sub/", IsFolder: true},
			{Title: "Page", URL: "/guides/page/", IsFolder: false},
		},
	}

	html := string(RenderAutoIndexContent(index))

	if !strings.Contains(html, "Guides") {
		t.Error("should contain title")
	}
	if !strings.Contains(html, "folder-item") {
		t.Error("should contain folder-item class")
	}
	if !strings.Contains(html, "page-item") {
		t.Error("should contain page-item class")
	}
}

func TestRenderAutoIndexContentEmpty(t *testing.T) {
	index := AutoIndex{
		Title:    "Empty",
		Children: []IndexItem{},
	}

	html := string(RenderAutoIndexContent(index))

	if !strings.Contains(html, "empty-folder") {
		t.Error("should show empty folder message")
	}
}

func TestCollectFoldersNeedingAutoIndex(t *testing.T) {
	root := &tree.Node{
		Path:     "",
		IsFolder: true,
		Children: []*tree.Node{
			{
				Name:     "WithIndex",
				Path:     "withindex",
				IsFolder: true,
				HasIndex: true,
			},
			{
				Name:     "WithoutIndex",
				Path:     "withoutindex",
				IsFolder: true,
				Children: []*tree.Node{
					{Path: "withoutindex/page.md", IsFolder: false},
				},
			},
		},
	}

	folders := collectFoldersNeedingAutoIndex(root)

	if len(folders) != 1 {
		t.Errorf("expected 1 folder needing auto-index, got %d", len(folders))
	}
	if len(folders) > 0 && folders[0].Name != "WithoutIndex" {
		t.Errorf("expected WithoutIndex folder, got %s", folders[0].Name)
	}
}

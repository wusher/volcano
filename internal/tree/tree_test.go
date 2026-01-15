package tree

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	node := NewNode("Test", "path/to/test", false)

	if node.Name != "Test" {
		t.Errorf("Name = %q, want %q", node.Name, "Test")
	}
	if node.Path != "path/to/test" {
		t.Errorf("Path = %q, want %q", node.Path, "path/to/test")
	}
	if node.IsFolder {
		t.Error("IsFolder should be false")
	}
	if node.Children == nil {
		t.Error("Children should not be nil")
	}
	if len(node.Children) != 0 {
		t.Errorf("Children length = %d, want 0", len(node.Children))
	}
}

func TestNodeAddChild(t *testing.T) {
	parent := NewNode("Parent", "", true)
	child := NewNode("Child", "child", false)

	parent.AddChild(child)

	if len(parent.Children) != 1 {
		t.Errorf("Children length = %d, want 1", len(parent.Children))
	}
	if parent.Children[0] != child {
		t.Error("Child not added correctly")
	}
	if child.Parent != parent {
		t.Error("Parent reference not set correctly")
	}
}

func TestNodeFindChild(t *testing.T) {
	parent := NewNode("Parent", "", true)
	child1 := NewNode("Child1", "child1", false)
	child2 := NewNode("Child2", "child2", false)

	parent.AddChild(child1)
	parent.AddChild(child2)

	// Find existing child
	found := parent.FindChild("Child1")
	if found != child1 {
		t.Error("FindChild should return child1")
	}

	// Find non-existing child
	notFound := parent.FindChild("NonExistent")
	if notFound != nil {
		t.Error("FindChild should return nil for non-existing child")
	}
}

func TestNodeIsEmpty(t *testing.T) {
	folder := NewNode("Folder", "", true)
	file := NewNode("File", "file.md", false)

	// Empty folder
	if !folder.IsEmpty() {
		t.Error("Empty folder should return true for IsEmpty")
	}

	// File is not empty
	if file.IsEmpty() {
		t.Error("File should return false for IsEmpty")
	}

	// Non-empty folder
	folder.AddChild(file)
	if folder.IsEmpty() {
		t.Error("Non-empty folder should return false for IsEmpty")
	}
}

func TestNodeHasMarkdownContent(t *testing.T) {
	// File always has content
	file := NewNode("File", "file.md", false)
	if !file.HasMarkdownContent() {
		t.Error("File should have markdown content")
	}

	// Empty folder has no content
	emptyFolder := NewNode("Empty", "", true)
	if emptyFolder.HasMarkdownContent() {
		t.Error("Empty folder should not have markdown content")
	}

	// Folder with file has content
	folderWithFile := NewNode("Folder", "", true)
	folderWithFile.AddChild(NewNode("File", "file.md", false))
	if !folderWithFile.HasMarkdownContent() {
		t.Error("Folder with file should have markdown content")
	}

	// Nested folder with file has content
	outerFolder := NewNode("Outer", "", true)
	innerFolder := NewNode("Inner", "inner", true)
	innerFile := NewNode("File", "inner/file.md", false)
	innerFolder.AddChild(innerFile)
	outerFolder.AddChild(innerFolder)
	if !outerFolder.HasMarkdownContent() {
		t.Error("Folder with nested file should have markdown content")
	}

	// Nested empty folders have no content
	emptyOuter := NewNode("EmptyOuter", "", true)
	emptyInner := NewNode("EmptyInner", "empty", true)
	emptyOuter.AddChild(emptyInner)
	if emptyOuter.HasMarkdownContent() {
		t.Error("Folder with only empty subfolders should not have markdown content")
	}
}

func TestBuildValidURLMap(t *testing.T) {
	// Create a simple site structure
	root := NewNode("Root", "", true)
	page1 := NewNode("Page1", "page1.md", false)
	page2 := NewNode("Page2", "page2.md", false)
	folder := NewNode("Folder", "folder", true)
	page3 := NewNode("Page3", "folder/page3.md", false)

	root.AddChild(page1)
	root.AddChild(page2)
	root.AddChild(folder)
	folder.AddChild(page3)

	site := &Site{
		Root:     root,
		AllPages: []*Node{page1, page2, page3},
	}

	t.Run("no base URL", func(t *testing.T) {
		validURLs := BuildValidURLMap(site, "")

		// Should have root
		if !validURLs["/"] {
			t.Error("Should include root URL /")
		}

		// Should have pages
		if !validURLs["/page1/"] {
			t.Error("Should include /page1/")
		}
		if !validURLs["/page2/"] {
			t.Error("Should include /page2/")
		}
		if !validURLs["/folder/page3/"] {
			t.Error("Should include /folder/page3/")
		}

		// Should have folder
		if !validURLs["/folder/"] {
			t.Error("Should include /folder/")
		}

		// Should NOT have prefixed URLs
		if validURLs["/volcano/page1/"] {
			t.Error("Should NOT include prefixed URLs when no base URL")
		}
	})

	t.Run("with base URL", func(t *testing.T) {
		validURLs := BuildValidURLMap(site, "https://example.com/volcano/")

		// Should have both unprefixed and prefixed root
		if !validURLs["/"] {
			t.Error("Should include root URL /")
		}
		if !validURLs["/volcano/"] {
			t.Error("Should include prefixed root /volcano/")
		}

		// Should have both unprefixed and prefixed pages
		if !validURLs["/page1/"] {
			t.Error("Should include /page1/")
		}
		if !validURLs["/volcano/page1/"] {
			t.Error("Should include prefixed /volcano/page1/")
		}

		if !validURLs["/folder/page3/"] {
			t.Error("Should include /folder/page3/")
		}
		if !validURLs["/volcano/folder/page3/"] {
			t.Error("Should include prefixed /volcano/folder/page3/")
		}

		// Should have both unprefixed and prefixed folders
		if !validURLs["/folder/"] {
			t.Error("Should include /folder/")
		}
		if !validURLs["/volcano/folder/"] {
			t.Error("Should include prefixed /volcano/folder/")
		}
	})

	t.Run("with base URL no path", func(t *testing.T) {
		validURLs := BuildValidURLMap(site, "https://example.com/")

		// Should only have unprefixed URLs (no base path to add)
		if !validURLs["/"] {
			t.Error("Should include root URL /")
		}
		if validURLs["/example/"] {
			t.Error("Should NOT add domain as prefix")
		}
	})
}

func TestBuildValidURLMapWithAutoIndex(t *testing.T) {
	// Create nodes
	page1 := NewNode("Page1", "page1.md", false)
	page2 := NewNode("Page2", "folder/page2.md", false)
	autoFolder := NewNode("AutoFolder", "autofolder", true)

	allPages := []*Node{page1, page2}
	autoIndexFolders := []*Node{autoFolder}

	t.Run("no base URL", func(t *testing.T) {
		validURLs := BuildValidURLMapWithAutoIndex(allPages, autoIndexFolders, "")

		// Should have root
		if !validURLs["/"] {
			t.Error("Should include root URL /")
		}

		// Should have pages
		if !validURLs["/page1/"] {
			t.Error("Should include /page1/")
		}
		if !validURLs["/folder/page2/"] {
			t.Error("Should include /folder/page2/")
		}

		// Should have auto-index folder
		if !validURLs["/autofolder/"] {
			t.Error("Should include auto-index folder /autofolder/")
		}
	})

	t.Run("with base URL", func(t *testing.T) {
		validURLs := BuildValidURLMapWithAutoIndex(allPages, autoIndexFolders, "https://example.com/volcano/")

		// Should have both unprefixed and prefixed
		if !validURLs["/page1/"] {
			t.Error("Should include /page1/")
		}
		if !validURLs["/volcano/page1/"] {
			t.Error("Should include prefixed /volcano/page1/")
		}

		// Should have both unprefixed and prefixed auto-index folder
		if !validURLs["/autofolder/"] {
			t.Error("Should include /autofolder/")
		}
		if !validURLs["/volcano/autofolder/"] {
			t.Error("Should include prefixed /volcano/autofolder/")
		}
	})
}

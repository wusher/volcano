// Package tree provides functionality for building a tree structure from markdown files.
package tree

// Node represents a node in the content tree
type Node struct {
	Name       string  // Clean display label (from H1 or filename)
	FileName   string  // Original filename
	H1Title    string  // Extracted H1 title (empty if none)
	Path       string  // Relative path from input root
	SourcePath string  // Full path to source .md file
	IsFolder   bool    // Whether this is a folder
	HasIndex   bool    // True if folder contains index.md
	IndexPath  string  // Path to index.md if exists
	Children   []*Node // Sorted alphabetically
	Parent     *Node   // Parent node
}

// Site represents the full site structure
type Site struct {
	Root     *Node   // Root of the tree
	AllPages []*Node // Flat list of all pages for easy iteration
}

// NewNode creates a new Node with the given name and path
func NewNode(name, path string, isFolder bool) *Node {
	return &Node{
		Name:     name,
		Path:     path,
		IsFolder: isFolder,
		Children: make([]*Node, 0),
	}
}

// AddChild adds a child node to this node and sets the parent reference
func (n *Node) AddChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

// FindChild finds a child node by name
func (n *Node) FindChild(name string) *Node {
	for _, child := range n.Children {
		if child.Name == name {
			return child
		}
	}
	return nil
}

// IsEmpty returns true if this is a folder with no children
func (n *Node) IsEmpty() bool {
	return n.IsFolder && len(n.Children) == 0
}

// HasMarkdownContent returns true if this folder contains markdown files (directly or nested)
func (n *Node) HasMarkdownContent() bool {
	if !n.IsFolder {
		return true // Files are content
	}
	for _, child := range n.Children {
		if child.HasMarkdownContent() {
			return true
		}
	}
	return false
}

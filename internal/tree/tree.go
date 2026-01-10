// Package tree provides functionality for building a tree structure from markdown files.
package tree

import (
	"os"
	"time"
)

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

// ModTime returns the file modification time for the source file
func (n *Node) ModTime() time.Time {
	if n.SourcePath == "" {
		return time.Now()
	}
	info, err := os.Stat(n.SourcePath)
	if err != nil {
		return time.Now()
	}
	return info.ModTime()
}

// BuildValidURLMap creates a map of all valid URLs from a site structure.
// This is used for validating internal links in generated content.
func BuildValidURLMap(site *Site) map[string]bool {
	validURLs := make(map[string]bool)
	validURLs["/"] = true

	// Add all page URLs
	for _, node := range site.AllPages {
		urlPath := GetURLPath(node)
		if urlPath != "" {
			validURLs[urlPath] = true
		}
	}

	// Add folder URLs (for auto-index pages)
	addFolderURLs(site.Root, validURLs)

	return validURLs
}

// BuildValidURLMapWithAutoIndex creates a map of all valid URLs including specific auto-index folders.
// The autoIndexFolders parameter contains additional folders that will have auto-generated indexes.
func BuildValidURLMapWithAutoIndex(allPages []*Node, autoIndexFolders []*Node) map[string]bool {
	validURLs := make(map[string]bool)
	validURLs["/"] = true

	// Add all page URLs
	for _, node := range allPages {
		urlPath := GetURLPath(node)
		if urlPath != "" {
			validURLs[urlPath] = true
		}
	}

	// Add auto-index folder URLs
	for _, folder := range autoIndexFolders {
		urlPath := "/" + SlugifyPath(folder.Path) + "/"
		validURLs[urlPath] = true
	}

	return validURLs
}

// addFolderURLs recursively adds folder URLs to the map
func addFolderURLs(node *Node, validURLs map[string]bool) {
	if node == nil {
		return
	}

	if node.IsFolder && node.Path != "" {
		urlPath := "/" + SlugifyPath(node.Path) + "/"
		validURLs[urlPath] = true
	}

	for _, child := range node.Children {
		addFolderURLs(child, validURLs)
	}
}

package tree

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Scan walks the input directory and builds a tree structure of markdown files
func Scan(inputDir string) (*Site, error) {
	absPath, err := filepath.Abs(inputDir)
	if err != nil {
		return nil, err
	}

	root := NewNode("", "", true)
	root.SourcePath = absPath

	allPages := make([]*Node, 0)

	err = scanDirectory(absPath, absPath, root, &allPages)
	if err != nil {
		return nil, err
	}

	// Sort children and prune empty folders
	sortAndPrune(root)

	return &Site{
		Root:     root,
		AllPages: allPages,
	}, nil
}

// scanDirectory recursively scans a directory for markdown files
func scanDirectory(basePath, currentPath string, parent *Node, allPages *[]*Node) error {
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files and folders
		if IsHidden(name) {
			continue
		}

		fullPath := filepath.Join(currentPath, name)
		relPath, err := filepath.Rel(basePath, fullPath)
		if err != nil {
			return err
		}

		if entry.IsDir() {
			// Create folder node
			folderNode := NewNode(CleanLabel(name), relPath, true)
			folderNode.SourcePath = fullPath
			parent.AddChild(folderNode)

			// Recursively scan subdirectory
			if err := scanDirectory(basePath, fullPath, folderNode, allPages); err != nil {
				return err
			}
		} else if IsMarkdownFile(name) {
			// Create file node with clean label as default name
			fileNode := NewNode(CleanLabel(name), relPath, false)
			fileNode.SourcePath = fullPath
			fileNode.FileName = name

			// Try to extract H1 title from the file content
			if content, err := os.ReadFile(fullPath); err == nil {
				if h1 := ExtractH1(content); h1 != "" {
					fileNode.H1Title = h1
					fileNode.Name = h1 // Override display name with H1
				}
			}

			// Check if this is an index file for the parent folder
			if IsIndexFile(name) && parent.IsFolder {
				parent.HasIndex = true
				parent.IndexPath = relPath
			}

			parent.AddChild(fileNode)
			*allPages = append(*allPages, fileNode)
		}
	}

	return nil
}

// sortAndPrune sorts children (folders first, then alphabetically) and removes empty folders
func sortAndPrune(node *Node) {
	if !node.IsFolder {
		return
	}

	// Recursively process children first
	for _, child := range node.Children {
		sortAndPrune(child)
	}

	// Remove empty folders (no markdown content at any depth)
	filtered := make([]*Node, 0, len(node.Children))
	for _, child := range node.Children {
		if child.IsFolder && !child.HasMarkdownContent() {
			continue // Skip empty folders
		}
		filtered = append(filtered, child)
	}
	node.Children = filtered

	// Sort: folders first, then files, alphabetically by name
	sort.Slice(node.Children, func(i, j int) bool {
		a, b := node.Children[i], node.Children[j]

		// Folders come before files
		if a.IsFolder != b.IsFolder {
			return a.IsFolder
		}

		// Alphabetical by name (case-insensitive)
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
}

// GetOutputPath returns the output path for a file node
// Converts: guides/intro.md → guides/intro/index.html (clean URLs)
// Converts: index.md → index.html (root index)
// Converts: posts/2024-01-15-hello.md → posts/hello/index.html (strips date prefix)
func GetOutputPath(node *Node) string {
	if node.IsFolder {
		return ""
	}

	// Get directory and filename
	dir := filepath.Dir(node.Path)
	filename := filepath.Base(node.Path)

	// Handle index files - they stay as index.html
	stem := strings.TrimSuffix(filename, filepath.Ext(filename))
	if strings.ToLower(stem) == "index" || strings.ToLower(stem) == "readme" {
		if dir == "." {
			return "index.html"
		}
		return filepath.Join(dir, "index.html")
	}

	// Extract metadata to get slug (strips date/number prefixes)
	meta := ExtractFileMetadata(filename, node.ModTime())
	slug := meta.Slug

	// For non-index files, create clean URLs: file.md → file/index.html
	if dir == "." {
		return filepath.Join(slug, "index.html")
	}
	return filepath.Join(dir, slug, "index.html")
}

// GetURLPath returns the URL path for a file node
// Converts: guides/intro.md → /guides/intro/
// Converts: index.md → /
// Converts: posts/2024-01-15-hello.md → /posts/hello/ (strips date prefix)
func GetURLPath(node *Node) string {
	if node.IsFolder {
		return ""
	}

	// Get directory and filename
	dir := filepath.Dir(node.Path)
	filename := filepath.Base(node.Path)

	// Handle index files
	stem := strings.TrimSuffix(filename, filepath.Ext(filename))
	if strings.ToLower(stem) == "index" || strings.ToLower(stem) == "readme" {
		if dir == "." {
			return "/"
		}
		return "/" + filepath.ToSlash(dir) + "/"
	}

	// Extract metadata to get slug (strips date/number prefixes)
	meta := ExtractFileMetadata(filename, node.ModTime())
	slug := meta.Slug

	// For non-index files, return the path with slug as directory
	if dir == "." {
		return "/" + slug + "/"
	}
	return "/" + filepath.ToSlash(dir) + "/" + slug + "/"
}

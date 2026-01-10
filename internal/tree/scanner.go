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

			// Skip adding root-level index/readme files to the tree
			// (site title already links to home page)
			isRootIndex := IsIndexFile(name) && filepath.Dir(relPath) == "."
			if !isRootIndex {
				parent.AddChild(fileNode)
			}
			*allPages = append(*allPages, fileNode)
		}
	}

	return nil
}

// sortAndPrune sorts children (files first, then folders, by date/number/name) and removes empty folders
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

	// Sort: files first, then folders
	// Within each category: by date (from filename), then number, then alphabetically
	sort.Slice(node.Children, func(i, j int) bool {
		a, b := node.Children[i], node.Children[j]

		// Files come before folders
		if a.IsFolder != b.IsFolder {
			return !a.IsFolder
		}

		// Both are same type - sort by date/number/name
		aMeta := GetNodeMetadata(a)
		bMeta := GetNodeMetadata(b)

		// Primary: Date (from filename only)
		// Items with dates come before items without dates
		if aMeta.HasDate != bMeta.HasDate {
			return aMeta.HasDate
		}
		// Both have dates - sort by date (newest first)
		if aMeta.HasDate && bMeta.HasDate && !aMeta.Date.Equal(bMeta.Date) {
			return aMeta.Date.After(bMeta.Date)
		}

		// Secondary: Number (lower numbers first, nil sorted last)
		aNum := numberForSort(aMeta.Number)
		bNum := numberForSort(bMeta.Number)
		if aNum != bNum {
			return aNum < bNum
		}

		// Tertiary: Name (alphabetical, case-insensitive)
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
}

// numberForSort returns the number for sorting, with nil treated as max
func numberForSort(n *int) int {
	if n == nil {
		return 999999
	}
	return *n
}

// GetOutputPath returns the output path for a file node
// Converts: guides/intro.md → guides/intro/index.html (clean URLs)
// Converts: index.md → index.html (root index)
// Converts: posts/2024-01-15-hello.md → posts/hello/index.html (strips date prefix)
// Converts: "0. Inbox/notes.md" → inbox/notes/index.html (slugified paths)
func GetOutputPath(node *Node) string {
	if node.IsFolder {
		return ""
	}

	// Get directory and filename
	dir := filepath.Dir(node.Path)
	filename := filepath.Base(node.Path)

	// Slugify the directory path to match URL paths
	slugDir := SlugifyPath(dir)

	// Handle index files - they stay as index.html
	stem := strings.TrimSuffix(filename, filepath.Ext(filename))
	if strings.ToLower(stem) == "index" || strings.ToLower(stem) == "readme" {
		if slugDir == "" {
			return "index.html"
		}
		return filepath.Join(slugDir, "index.html")
	}

	// Extract metadata to get slug (strips date/number prefixes)
	meta := ExtractFileMetadata(filename, node.ModTime())
	slug := meta.Slug

	// For non-index files, create clean URLs: file.md → file/index.html
	if slugDir == "" {
		return filepath.Join(slug, "index.html")
	}
	return filepath.Join(slugDir, slug, "index.html")
}

// GetURLPath returns the URL path for a file node
// Converts: guides/intro.md → /guides/intro/
// Converts: index.md → /
// Converts: posts/2024-01-15-hello.md → /posts/hello/ (strips date prefix)
// Converts: "0. Inbox/notes.md" → /inbox/notes/
func GetURLPath(node *Node) string {
	if node.IsFolder {
		// For folders, return the slugified path
		slugPath := SlugifyPath(node.Path)
		if slugPath == "" {
			return "/"
		}
		return "/" + slugPath + "/"
	}

	// Get directory and filename
	dir := filepath.Dir(node.Path)
	filename := filepath.Base(node.Path)

	// Slugify the directory path
	slugDir := SlugifyPath(dir)

	// Handle index files
	stem := strings.TrimSuffix(filename, filepath.Ext(filename))
	if strings.ToLower(stem) == "index" || strings.ToLower(stem) == "readme" {
		if slugDir == "" {
			return "/"
		}
		return "/" + slugDir + "/"
	}

	// Extract metadata to get slug (strips date/number prefixes)
	meta := ExtractFileMetadata(filename, node.ModTime())
	slug := meta.Slug

	// For non-index files, return the path with slug as directory
	if slugDir == "" {
		return "/" + slug + "/"
	}
	return "/" + slugDir + "/" + slug + "/"
}

// SlugifyPath slugifies each segment of a path
// Converts: "0. Inbox/1. Health" → "inbox/health"
func SlugifyPath(path string) string {
	if path == "." || path == "" {
		return ""
	}

	// Split path into segments
	segments := strings.Split(filepath.ToSlash(path), "/")
	slugged := make([]string, 0, len(segments))

	for _, seg := range segments {
		if seg == "" || seg == "." {
			continue
		}
		slugged = append(slugged, Slugify(seg))
	}

	return strings.Join(slugged, "/")
}

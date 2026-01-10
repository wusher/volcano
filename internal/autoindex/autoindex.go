// Package autoindex provides auto-generated folder index functionality.
package autoindex

import (
	"html/template"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wusher/volcano/internal/tree"
)

// Item represents an item in an auto-generated folder index
type Item struct {
	Title    string
	URL      string
	IsFolder bool
}

// Index represents an auto-generated index page for a folder
type Index struct {
	FolderNode *tree.Node
	Title      string
	Children   []Item
	OutputPath string
	URLPath    string
}

// NeedsAutoIndex determines if a folder needs an auto-generated index
func NeedsAutoIndex(node *tree.Node) bool {
	if !node.IsFolder {
		return false
	}

	// Check if folder already has an index file
	if node.HasIndex {
		return false
	}

	// Check for index.md or readme.md in children by looking at the Path (filename)
	for _, child := range node.Children {
		if !child.IsFolder {
			// Use the path's base name instead of the display name (which could be H1 title)
			baseName := strings.TrimSuffix(filepath.Base(child.Path), filepath.Ext(child.Path))
			lower := strings.ToLower(baseName)
			if lower == "index" || lower == "readme" {
				return false
			}
		}
	}

	return true
}

// Build creates an Index for a folder
func Build(node *tree.Node) Index {
	return BuildWithBaseURL(node, "")
}

// BuildWithBaseURL creates an Index for a folder with base URL prefixing
func BuildWithBaseURL(node *tree.Node, baseURL string) Index {
	var items []Item

	for _, child := range node.Children {
		url := tree.GetURLPath(child)
		// For folders, construct the URL from the path
		if child.IsFolder {
			url = "/" + tree.SlugifyPath(child.Path) + "/"
		}
		// Apply base URL prefix
		prefixedURL := tree.PrefixURL(baseURL, url)
		items = append(items, Item{
			Title:    child.Name,
			URL:      prefixedURL,
			IsFolder: child.IsFolder,
		})
	}

	// Sort: files first, then folders, then alphabetically
	// (matches the tree navigation sort order)
	sort.Slice(items, func(i, j int) bool {
		if items[i].IsFolder != items[j].IsFolder {
			return !items[i].IsFolder // files first
		}
		return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
	})

	slugPath := tree.SlugifyPath(node.Path)
	urlPath := "/" + slugPath + "/"
	outputPath := filepath.Join(slugPath, "index.html")
	if node.Path == "" || node.Path == "." {
		urlPath = "/"
		outputPath = "index.html"
	}

	return Index{
		FolderNode: node,
		Title:      node.Name,
		Children:   items,
		OutputPath: outputPath,
		URLPath:    urlPath,
	}
}

// RenderContent generates HTML content for an auto-generated index page
func RenderContent(index Index) template.HTML {
	var sb strings.Builder

	sb.WriteString(`<article class="auto-index-page">`)
	sb.WriteString("\n")
	sb.WriteString(`<h1>`)
	sb.WriteString(template.HTMLEscapeString(index.Title))
	sb.WriteString(`</h1>`)
	sb.WriteString("\n")

	if len(index.Children) > 0 {
		sb.WriteString(`<ul class="folder-index">`)
		sb.WriteString("\n")
		for _, item := range index.Children {
			itemClass := "page-item"
			if item.IsFolder {
				itemClass = "folder-item"
			}
			sb.WriteString(`<li class="`)
			sb.WriteString(itemClass)
			sb.WriteString(`">`)
			sb.WriteString("\n")
			sb.WriteString(`<a href="`)
			sb.WriteString(item.URL)
			sb.WriteString(`">`)
			sb.WriteString(template.HTMLEscapeString(item.Title))
			sb.WriteString(`</a>`)
			sb.WriteString("\n")
			sb.WriteString(`</li>`)
			sb.WriteString("\n")
		}
		sb.WriteString(`</ul>`)
		sb.WriteString("\n")
	} else {
		sb.WriteString(`<p class="empty-folder">This folder is empty.</p>`)
		sb.WriteString("\n")
	}

	sb.WriteString(`</article>`)

	return template.HTML(sb.String())
}

// CollectFoldersNeedingAutoIndex returns all folders that need auto-generated indexes
func CollectFoldersNeedingAutoIndex(node *tree.Node) []*tree.Node {
	var folders []*tree.Node

	if node.IsFolder && NeedsAutoIndex(node) && node.Path != "" && node.Path != "." {
		folders = append(folders, node)
	}

	for _, child := range node.Children {
		if child.IsFolder {
			folders = append(folders, CollectFoldersNeedingAutoIndex(child)...)
		}
	}

	return folders
}

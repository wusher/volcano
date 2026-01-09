package generator

import (
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"volcano/internal/navigation"
	"volcano/internal/seo"
	"volcano/internal/templates"
	"volcano/internal/tree"
)

// IndexItem represents an item in an auto-generated folder index
type IndexItem struct {
	Title    string
	URL      string
	IsFolder bool
}

// AutoIndex represents an auto-generated index page for a folder
type AutoIndex struct {
	FolderNode *tree.Node
	Title      string
	Children   []IndexItem
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

// BuildAutoIndex creates an AutoIndex for a folder
func BuildAutoIndex(node *tree.Node) AutoIndex {
	var items []IndexItem

	for _, child := range node.Children {
		url := tree.GetURLPath(child)
		// For folders, construct the URL from the path
		if child.IsFolder {
			url = "/" + filepath.ToSlash(child.Path) + "/"
		}
		items = append(items, IndexItem{
			Title:    child.Name,
			URL:      url,
			IsFolder: child.IsFolder,
		})
	}

	// Sort: folders first, then alphabetically
	sort.Slice(items, func(i, j int) bool {
		if items[i].IsFolder != items[j].IsFolder {
			return items[i].IsFolder
		}
		return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
	})

	urlPath := "/" + node.Path + "/"
	if node.Path == "" || node.Path == "." {
		urlPath = "/"
	}

	return AutoIndex{
		FolderNode: node,
		Title:      node.Name,
		Children:   items,
		OutputPath: filepath.Join(node.Path, "index.html"),
		URLPath:    urlPath,
	}
}

// RenderAutoIndexContent generates HTML content for an auto-generated index page
func RenderAutoIndexContent(index AutoIndex) template.HTML {
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

// generateAutoIndex generates an auto-index page for a folder without an index.md
func (g *Generator) generateAutoIndex(node *tree.Node, root *tree.Node) error {
	index := BuildAutoIndex(node)
	fullOutputPath := filepath.Join(g.config.OutputDir, index.OutputPath)

	// Build content
	htmlContent := RenderAutoIndexContent(index)

	// Build breadcrumbs
	breadcrumbs := navigation.BuildBreadcrumbs(node, g.config.Title)
	breadcrumbsHTML := navigation.RenderBreadcrumbs(breadcrumbs)

	// Generate SEO meta tags
	seoConfig := seo.Config{
		SiteURL:   g.config.SiteURL,
		SiteTitle: g.config.Title,
		Author:    g.config.Author,
		OGImage:   g.config.OGImage,
	}
	pageMeta := seo.GeneratePageMeta(index.Title, string(htmlContent), index.URLPath, seoConfig)
	metaTagsHTML := seo.RenderMetaTags(pageMeta)

	// Render navigation
	nav := templates.RenderNavigation(root, index.URLPath)

	// Prepare template data
	data := templates.PageData{
		SiteTitle:   g.config.Title,
		PageTitle:   index.Title,
		Content:     htmlContent,
		Navigation:  nav,
		CurrentPath: index.URLPath,
		Breadcrumbs: breadcrumbsHTML,
		MetaTags:    metaTagsHTML,
		ShowSearch:  true,
	}

	// Create output directory
	outputDir := filepath.Dir(fullOutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Write file
	f, err := os.Create(fullOutputPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	return g.renderer.Render(f, data)
}

// collectFoldersNeedingAutoIndex returns all folders that need auto-generated indexes
func collectFoldersNeedingAutoIndex(node *tree.Node) []*tree.Node {
	var folders []*tree.Node

	if node.IsFolder && NeedsAutoIndex(node) && node.Path != "" && node.Path != "." {
		folders = append(folders, node)
	}

	for _, child := range node.Children {
		if child.IsFolder {
			folders = append(folders, collectFoldersNeedingAutoIndex(child)...)
		}
	}

	return folders
}

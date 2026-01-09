package navigation

import (
	"html/template"
	"path/filepath"
	"strings"

	"volcano/internal/tree"
)

// NavLink represents a navigation link (previous or next page)
type NavLink struct {
	Title   string // Page title
	URL     string // URL path
	Section string // Parent folder name (optional)
}

// PageNavigation contains previous and next links for a page
type PageNavigation struct {
	Previous *NavLink
	Next     *NavLink
}

// BuildPageNavigation creates prev/next navigation for a page
func BuildPageNavigation(currentNode *tree.Node, allPages []*tree.Node) PageNavigation {
	if currentNode == nil || len(allPages) == 0 {
		return PageNavigation{}
	}

	// Find current page index
	currentIdx := -1
	for i, page := range allPages {
		if page.SourcePath == currentNode.SourcePath {
			currentIdx = i
			break
		}
	}

	if currentIdx == -1 {
		return PageNavigation{}
	}

	nav := PageNavigation{}

	// Previous page
	if currentIdx > 0 {
		prevNode := allPages[currentIdx-1]
		nav.Previous = &NavLink{
			Title:   prevNode.Name,
			URL:     tree.GetURLPath(prevNode),
			Section: getSection(prevNode),
		}
	}

	// Next page
	if currentIdx < len(allPages)-1 {
		nextNode := allPages[currentIdx+1]
		nav.Next = &NavLink{
			Title:   nextNode.Name,
			URL:     tree.GetURLPath(nextNode),
			Section: getSection(nextNode),
		}
	}

	return nav
}

// getSection returns the parent folder name if different from root
func getSection(node *tree.Node) string {
	if node.Parent != nil && node.Parent.Parent != nil { // Has parent that isn't root
		return node.Parent.Name
	}
	return ""
}

// RenderPageNavigation renders the prev/next navigation as HTML
func RenderPageNavigation(nav PageNavigation) template.HTML {
	if nav.Previous == nil && nav.Next == nil {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`<nav class="page-nav" aria-label="Page navigation">`)
	sb.WriteString("\n")

	if nav.Previous != nil {
		sb.WriteString(`  <a href="`)
		sb.WriteString(template.HTMLEscapeString(nav.Previous.URL))
		sb.WriteString(`" class="page-nav-prev">`)
		sb.WriteString("\n")
		sb.WriteString(`    <span class="page-nav-label">Previous</span>`)
		sb.WriteString("\n")
		sb.WriteString(`    <span class="page-nav-title">`)
		sb.WriteString(template.HTMLEscapeString("← " + nav.Previous.Title))
		sb.WriteString(`</span>`)
		sb.WriteString("\n")
		sb.WriteString(`  </a>`)
		sb.WriteString("\n")
	} else {
		sb.WriteString(`  <span class="page-nav-placeholder"></span>`)
		sb.WriteString("\n")
	}

	if nav.Next != nil {
		sb.WriteString(`  <a href="`)
		sb.WriteString(template.HTMLEscapeString(nav.Next.URL))
		sb.WriteString(`" class="page-nav-next">`)
		sb.WriteString("\n")
		sb.WriteString(`    <span class="page-nav-label">Next</span>`)
		sb.WriteString("\n")
		sb.WriteString(`    <span class="page-nav-title">`)
		sb.WriteString(template.HTMLEscapeString(nav.Next.Title + " →"))
		sb.WriteString(`</span>`)
		sb.WriteString("\n")
		sb.WriteString(`  </a>`)
		sb.WriteString("\n")
	}

	sb.WriteString(`</nav>`)
	sb.WriteString("\n")

	return template.HTML(sb.String())
}

// FlattenTreeForPagination returns all pages in depth-first order for pagination
func FlattenTreeForPagination(root *tree.Node) []*tree.Node {
	var pages []*tree.Node
	flattenNode(root, &pages)
	return pages
}

func flattenNode(node *tree.Node, pages *[]*tree.Node) {
	if node == nil {
		return
	}

	for _, child := range node.Children {
		if child.IsFolder {
			// If folder has index, add it first
			if child.HasIndex {
				// Find the index node
				for _, grandchild := range child.Children {
					if tree.IsIndexFile(filepath.Base(grandchild.SourcePath)) {
						*pages = append(*pages, grandchild)
						break
					}
				}
			}
			// Recurse into folder
			flattenNode(child, pages)
		} else {
			// Skip index files that were already added
			if !tree.IsIndexFile(filepath.Base(node.SourcePath)) || node.Parent == nil || !node.Parent.HasIndex {
				*pages = append(*pages, child)
			}
		}
	}
}

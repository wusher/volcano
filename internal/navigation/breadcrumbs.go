// Package navigation provides navigation components like breadcrumbs and pagination.
package navigation

import (
	"html/template"
	"strconv"
	"strings"

	"volcano/internal/tree"
)

// Breadcrumb represents a single breadcrumb item
type Breadcrumb struct {
	Label   string // Display label
	URL     string // URL path (empty for current page)
	Current bool   // True if this is the current page
}

// BuildBreadcrumbs generates breadcrumb trail for a page
func BuildBreadcrumbs(node *tree.Node, siteTitle string) []Breadcrumb {
	if node == nil {
		return nil
	}

	var crumbs []Breadcrumb

	// Add home link
	crumbs = append(crumbs, Breadcrumb{
		Label:   siteTitle,
		URL:     "/",
		Current: false,
	})

	// Get path segments
	urlPath := tree.GetURLPath(node)
	if urlPath == "/" {
		// We're on the home page, just return home as current
		crumbs[0].Current = true
		crumbs[0].URL = ""
		return crumbs
	}

	// Build path from node's parent chain
	var ancestors []*tree.Node
	current := node.Parent
	for current != nil && current.Parent != nil { // Stop before root
		ancestors = append([]*tree.Node{current}, ancestors...)
		current = current.Parent
	}

	// Add ancestor folders
	for _, ancestor := range ancestors {
		if ancestor.IsFolder {
			// All folders get a slugified URL (they either have an index or auto-index)
			url := "/" + tree.SlugifyPath(ancestor.Path) + "/"
			if url == "/./" {
				url = "/"
			}
			crumbs = append(crumbs, Breadcrumb{
				Label:   ancestor.Name,
				URL:     url,
				Current: false,
			})
		}
	}

	// Add current page
	crumbs = append(crumbs, Breadcrumb{
		Label:   node.Name,
		URL:     "",
		Current: true,
	})

	return crumbs
}

// RenderBreadcrumbs renders breadcrumbs as HTML
func RenderBreadcrumbs(crumbs []Breadcrumb) template.HTML {
	if len(crumbs) <= 1 {
		return "" // Don't show breadcrumbs on home page or if only one item
	}

	var sb strings.Builder
	sb.WriteString(`<nav class="breadcrumbs" aria-label="Breadcrumb">`)
	sb.WriteString("\n")
	sb.WriteString(`  <ol itemscope itemtype="https://schema.org/BreadcrumbList">`)
	sb.WriteString("\n")

	for i, crumb := range crumbs {
		sb.WriteString(`    <li itemprop="itemListElement" itemscope itemtype="https://schema.org/ListItem">`)
		sb.WriteString("\n")

		if crumb.Current {
			sb.WriteString(`      <span itemprop="name" aria-current="page">`)
			sb.WriteString(template.HTMLEscapeString(crumb.Label))
			sb.WriteString(`</span>`)
		} else if crumb.URL != "" {
			sb.WriteString(`      <a itemprop="item" href="`)
			sb.WriteString(template.HTMLEscapeString(crumb.URL))
			sb.WriteString(`"><span itemprop="name">`)
			sb.WriteString(template.HTMLEscapeString(crumb.Label))
			sb.WriteString(`</span></a>`)
		} else {
			sb.WriteString(`      <span itemprop="name">`)
			sb.WriteString(template.HTMLEscapeString(crumb.Label))
			sb.WriteString(`</span>`)
		}
		sb.WriteString("\n")

		sb.WriteString(`      <meta itemprop="position" content="`)
		sb.WriteString(template.HTMLEscapeString(strconv.Itoa(i + 1)))
		sb.WriteString(`" />`)
		sb.WriteString("\n")
		sb.WriteString(`    </li>`)
		sb.WriteString("\n")
	}

	sb.WriteString(`  </ol>`)
	sb.WriteString("\n")
	sb.WriteString(`</nav>`)
	sb.WriteString("\n")

	return template.HTML(sb.String())
}

// Package templates provides HTML template rendering for generated pages.
package templates

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"strings"

	"volcano/internal/tree"
)

//go:embed layout.html
var layoutFS embed.FS

// TopNavItem represents an item in the top navigation bar
type TopNavItem struct {
	Name string // Display name
	URL  string // URL path
}

// PageData contains all data needed to render a page
type PageData struct {
	SiteTitle    string        // Site title for header
	PageTitle    string        // Current page title
	Content      template.HTML // Rendered HTML content
	Navigation   template.HTML // Rendered navigation HTML
	CurrentPath  string        // Current page URL path for active state
	CSS          template.CSS  // Embedded CSS styles
	Breadcrumbs  template.HTML // Breadcrumb navigation
	PageNav      template.HTML // Previous/Next navigation
	TOC          template.HTML // Table of contents
	MetaTags     template.HTML // SEO meta tags
	FaviconLinks template.HTML // Favicon link tags
	ReadingTime  string        // Reading time display (e.g., "5 min read")
	LastModified string        // Last modified date (e.g., "January 5, 2025")
	HasTOC       bool          // Whether to show TOC sidebar
	ShowSearch   bool          // Whether to show nav search
	TopNavItems  []TopNavItem  // Items for top navigation bar (when --top-nav enabled)
}

// Renderer handles HTML template rendering
type Renderer struct {
	tmpl *template.Template
	css  string
}

// NewRenderer creates a new template renderer
func NewRenderer(css string) (*Renderer, error) {
	tmplContent, err := layoutFS.ReadFile("layout.html")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("layout").Parse(string(tmplContent))
	if err != nil {
		return nil, err
	}

	return &Renderer{
		tmpl: tmpl,
		css:  css,
	}, nil
}

// Render renders a page with the given data
func (r *Renderer) Render(w io.Writer, data PageData) error {
	data.CSS = template.CSS(r.css)
	return r.tmpl.Execute(w, data)
}

// RenderToString renders a page and returns the result as a string
func (r *Renderer) RenderToString(data PageData) (string, error) {
	var buf bytes.Buffer
	if err := r.Render(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderNavigation renders the navigation tree as HTML
func RenderNavigation(root *tree.Node, currentPath string) template.HTML {
	var buf bytes.Buffer
	renderNavNode(&buf, root.Children, currentPath, 0)
	return template.HTML(buf.String())
}

// renderNavNode recursively renders navigation nodes
func renderNavNode(buf *bytes.Buffer, nodes []*tree.Node, currentPath string, depth int) {
	if len(nodes) == 0 {
		return
	}

	buf.WriteString("<ul role=\"tree\">\n")

	for _, node := range nodes {
		if node.IsFolder {
			renderFolderNode(buf, node, currentPath, depth)
		} else {
			renderFileNode(buf, node, currentPath)
		}
	}

	buf.WriteString("</ul>\n")
}

// renderFolderNode renders a folder node with its children
func renderFolderNode(buf *bytes.Buffer, node *tree.Node, currentPath string, depth int) {
	buf.WriteString("<li role=\"treeitem\" class=\"folder\" data-search-text=\"")
	buf.WriteString(template.HTMLEscapeString(node.Name))
	buf.WriteString("\">\n")
	buf.WriteString("<div class=\"folder-header\">\n")

	// Toggle button
	buf.WriteString("<button class=\"folder-toggle\" aria-label=\"Toggle folder\">\n")
	buf.WriteString("<svg class=\"chevron\" xmlns=\"http://www.w3.org/2000/svg\" width=\"16\" height=\"16\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\"><polyline points=\"9 18 15 12 9 6\"></polyline></svg>\n")
	buf.WriteString("</button>\n")

	// Folder link (if has index) or just label
	if node.HasIndex {
		urlPath := tree.GetURLPath(&tree.Node{Path: node.IndexPath})
		active := ""
		if urlPath == currentPath {
			active = " active"
		}
		buf.WriteString("<a href=\"" + template.HTMLEscapeString(urlPath) + "\" class=\"folder-link" + active + "\">" + template.HTMLEscapeString(node.Name) + "</a>\n")
	} else {
		buf.WriteString("<span class=\"folder-label\">" + template.HTMLEscapeString(node.Name) + "</span>\n")
	}

	buf.WriteString("</div>\n")

	// Render children
	if len(node.Children) > 0 {
		buf.WriteString("<div class=\"folder-children\">\n")
		renderNavNode(buf, node.Children, currentPath, depth+1)
		buf.WriteString("</div>\n")
	}

	buf.WriteString("</li>\n")
}

// renderFileNode renders a file node as a link
func renderFileNode(buf *bytes.Buffer, node *tree.Node, currentPath string) {
	urlPath := tree.GetURLPath(node)
	active := ""
	if urlPath == currentPath {
		active = " active"
	}

	buf.WriteString("<li role=\"treeitem\" data-search-text=\"")
	buf.WriteString(template.HTMLEscapeString(node.Name))
	buf.WriteString("\">\n")
	buf.WriteString("<a href=\"" + template.HTMLEscapeString(urlPath) + "\" class=\"file-link" + active + "\">" + template.HTMLEscapeString(node.Name) + "</a>\n")
	buf.WriteString("</li>\n")
}

// BuildTopNavItems extracts root-level files for top navigation bar
// Returns nil if topNav is disabled or there are more than 5 root files
func BuildTopNavItems(root *tree.Node, topNav bool) []TopNavItem {
	if !topNav || root == nil {
		return nil
	}

	// Count root files (excluding index/readme files)
	var rootFiles []*tree.Node
	for _, child := range root.Children {
		if child.IsFolder {
			continue
		}
		// Skip index files using the same logic as tree package
		filename := strings.ToLower(child.FileName)
		if filename == "index.md" || filename == "readme.md" {
			continue
		}
		rootFiles = append(rootFiles, child)
	}

	// Only use top nav if between 1 and 5 root files
	if len(rootFiles) == 0 || len(rootFiles) > 5 {
		return nil
	}

	// Build top nav items
	var items []TopNavItem
	for _, node := range rootFiles {
		items = append(items, TopNavItem{
			Name: node.Name,
			URL:  tree.GetURLPath(node),
		})
	}

	return items
}

// RenderNavigationWithTopNav renders navigation excluding root files when top nav is enabled
func RenderNavigationWithTopNav(root *tree.Node, currentPath string, topNavItems []TopNavItem) template.HTML {
	if len(topNavItems) == 0 {
		return RenderNavigation(root, currentPath)
	}

	// Create a set of URLs to exclude from sidebar
	topNavURLs := make(map[string]bool)
	for _, item := range topNavItems {
		topNavURLs[item.URL] = true
	}

	var buf bytes.Buffer
	renderNavNodeFiltered(&buf, root.Children, currentPath, 0, topNavURLs)
	return template.HTML(buf.String())
}

// renderNavNodeFiltered renders nav nodes, filtering out specified URLs
func renderNavNodeFiltered(buf *bytes.Buffer, nodes []*tree.Node, currentPath string, depth int, excludeURLs map[string]bool) {
	if len(nodes) == 0 {
		return
	}

	buf.WriteString("<ul role=\"tree\">\n")

	for _, node := range nodes {
		// Skip excluded files
		if !node.IsFolder && excludeURLs[tree.GetURLPath(node)] {
			continue
		}

		if node.IsFolder {
			renderFolderNode(buf, node, currentPath, depth)
		} else {
			renderFileNode(buf, node, currentPath)
		}
	}

	buf.WriteString("</ul>\n")
}


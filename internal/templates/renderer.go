// Package templates provides HTML template rendering for generated pages.
package templates

import (
	"bytes"
	"embed"
	"html/template"
	"io"

	"volcano/internal/tree"
)

//go:embed layout.html
var layoutFS embed.FS

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

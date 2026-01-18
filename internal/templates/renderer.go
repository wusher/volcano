// Package templates provides HTML template rendering for generated pages.
package templates

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"sort"
	"strings"

	"github.com/wusher/volcano/internal/minify"
	"github.com/wusher/volcano/internal/tree"
)

//go:embed layout.html layout.js
var layoutFS embed.FS

// TopNavItem represents an item in the top navigation bar
type TopNavItem struct {
	Name string // Display name
	URL  string // URL path
}

// PageData contains all data needed to render a page
type PageData struct {
	SiteTitle       string        // Site title for header
	PageTitle       string        // Current page title
	Content         template.HTML // Rendered HTML content
	Navigation      template.HTML // Rendered navigation HTML
	CurrentPath     string        // Current page URL path for active state
	CSS             template.CSS  // Embedded CSS styles (used when CSSURL is empty)
	CSSURL          string        // External CSS file URL (when set, CSS is ignored)
	JSURL           string        // External JS file URL (when set, InstantNavJS is ignored)
	Breadcrumbs     template.HTML // Breadcrumb navigation
	PageNav         template.HTML // Previous/Next navigation
	TOC             template.HTML // Table of contents
	MetaTags        template.HTML // SEO meta tags
	FaviconLinks    template.HTML // Favicon link tags
	ReadingTime     string        // Reading time display (e.g., "5 min read")
	HasTOC          bool          // Whether to show TOC sidebar
	ShowSearch      bool          // Whether to show nav search
	TopNavItems     []TopNavItem  // Items for top navigation bar (when --top-nav enabled)
	BaseURL         string        // Base URL path prefix for all links (e.g., "/volcano")
	InstantNavJS    template.JS   // Instant navigation JavaScript (when --instant-nav enabled)
	ViewTransitions bool          // Enable browser view transitions API (when --view-transitions enabled)
	PWAEnabled      bool          // Whether PWA is enabled (adds manifest link + SW registration)
	SearchEnabled   bool          // Whether search is enabled (adds command palette + lazy load)
	InlineJS        template.JS   // Minified inline JavaScript for page functionality
}

// Renderer handles HTML template rendering
type Renderer struct {
	tmpl     *template.Template
	css      string
	inlineJS string
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

	// Load and minify inline JavaScript
	jsContent, err := layoutFS.ReadFile("layout.js")
	if err != nil {
		return nil, err
	}
	inlineJS := minify.JS(string(jsContent))

	return &Renderer{
		tmpl:     tmpl,
		css:      css,
		inlineJS: inlineJS,
	}, nil
}

// Render renders a page with the given data
func (r *Renderer) Render(w io.Writer, data PageData) error {
	data.CSS = template.CSS(r.css)
	data.InlineJS = template.JS(r.inlineJS)
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
	return RenderNavigationWithBaseURL(root, currentPath, "")
}

// RenderNavigationWithBaseURL renders the navigation tree as HTML with base URL prefixing
func RenderNavigationWithBaseURL(root *tree.Node, currentPath, baseURL string) template.HTML {
	var buf bytes.Buffer
	renderNavNode(&buf, root.Children, currentPath, 0, baseURL)
	return template.HTML(buf.String())
}

// renderNavNode recursively renders navigation nodes
func renderNavNode(buf *bytes.Buffer, nodes []*tree.Node, currentPath string, depth int, baseURL string) {
	if len(nodes) == 0 {
		return
	}

	buf.WriteString("<ul role=\"tree\">\n")

	for _, node := range nodes {
		if node.IsFolder {
			renderFolderNode(buf, node, currentPath, depth, baseURL)
		} else {
			renderFileNode(buf, node, currentPath, baseURL)
		}
	}

	buf.WriteString("</ul>\n")
}

// renderFolderNode renders a folder node with its children
func renderFolderNode(buf *bytes.Buffer, node *tree.Node, currentPath string, depth int, baseURL string) {
	buf.WriteString("<li role=\"treeitem\" class=\"folder\" data-search-text=\"")
	buf.WriteString(template.HTMLEscapeString(node.Name))
	buf.WriteString("\">\n")
	buf.WriteString("<div class=\"folder-header\">\n")

	// Toggle button
	buf.WriteString("<button class=\"folder-toggle\" aria-label=\"Toggle folder\">\n")
	buf.WriteString("<svg class=\"chevron\" xmlns=\"http://www.w3.org/2000/svg\" width=\"16\" height=\"16\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\"><polyline points=\"9 18 15 12 9 6\"></polyline></svg>\n")
	buf.WriteString("</button>\n")

	// Folder link - either has an index.md OR will have auto-generated index
	// All folders with children are clickable (auto-index will be generated)
	folderURL := "/" + tree.SlugifyPath(node.Path) + "/"
	if node.HasIndex {
		folderURL = tree.GetURLPath(&tree.Node{Path: node.IndexPath})
	}
	// Apply base URL prefix
	prefixedURL := tree.PrefixURL(baseURL, folderURL)
	active := ""
	if folderURL == currentPath {
		active = " active"
	}
	buf.WriteString("<a href=\"" + template.HTMLEscapeString(prefixedURL) + "\" class=\"folder-link" + active + "\">" + template.HTMLEscapeString(node.Name) + "</a>\n")

	buf.WriteString("</div>\n")

	// Render children
	if len(node.Children) > 0 {
		buf.WriteString("<div class=\"folder-children\">\n")
		renderNavNode(buf, node.Children, currentPath, depth+1, baseURL)
		buf.WriteString("</div>\n")
	}

	buf.WriteString("</li>\n")
}

// renderFileNode renders a file node as a link
func renderFileNode(buf *bytes.Buffer, node *tree.Node, currentPath string, baseURL string) {
	urlPath := tree.GetURLPath(node)
	prefixedURL := tree.PrefixURL(baseURL, urlPath)
	active := ""
	if urlPath == currentPath {
		active = " active"
	}

	buf.WriteString("<li role=\"treeitem\" data-search-text=\"")
	buf.WriteString(template.HTMLEscapeString(node.Name))
	buf.WriteString("\">\n")
	buf.WriteString("<a href=\"" + template.HTMLEscapeString(prefixedURL) + "\" class=\"file-link" + active + "\">" + template.HTMLEscapeString(node.Name) + "</a>\n")
	buf.WriteString("</li>\n")
}

// BuildTopNavItems extracts root-level items for top navigation bar
// Returns nil if topNav is disabled or there are no eligible items
// Items are sorted: files first, then folders, each sorted by date/number/name
func BuildTopNavItems(root *tree.Node, topNav bool) []TopNavItem {
	return BuildTopNavItemsWithBaseURL(root, topNav, "")
}

// BuildTopNavItemsWithBaseURL extracts root-level items for top navigation bar with base URL prefixing
func BuildTopNavItemsWithBaseURL(root *tree.Node, topNav bool, baseURL string) []TopNavItem {
	if !topNav || root == nil {
		return nil
	}

	// Collect root items (excluding index/readme files)
	var rootItems []*tree.Node
	for _, child := range root.Children {
		// Skip index files using the same logic as tree package
		if !child.IsFolder {
			filename := strings.ToLower(child.FileName)
			if filename == "index.md" || filename == "readme.md" {
				continue
			}
		}
		rootItems = append(rootItems, child)
	}

	// Only use top nav if between 1 and 8 root items
	// (increased from 5 to accommodate both files and folders)
	if len(rootItems) == 0 || len(rootItems) > 8 {
		return nil
	}

	// Sort: files first, then folders
	// Within each category: by date (newest first), then number, then name alphabetically
	// Date is extracted from filename prefix only (e.g., 2024-01-15-title.md)
	sort.Slice(rootItems, func(i, j int) bool {
		a, b := rootItems[i], rootItems[j]

		// Files before folders
		if a.IsFolder != b.IsFolder {
			return !a.IsFolder // files (false) before folders (true)
		}

		// Both are same type - sort by date/number/name using tree metadata
		aMeta := tree.GetNodeMetadata(a)
		bMeta := tree.GetNodeMetadata(b)

		// Primary: Date (from filename only)
		// Items with dates come before items without dates
		if aMeta.HasDate != bMeta.HasDate {
			return aMeta.HasDate // items with dates first
		}
		// Both have dates - sort by date (newest first)
		if aMeta.HasDate && bMeta.HasDate && !aMeta.Date.Equal(bMeta.Date) {
			return aMeta.Date.After(bMeta.Date)
		}

		// Secondary: Number (lower numbers first, nil sorted last)
		aNum := topNavNumberForSort(aMeta.Number)
		bNum := topNavNumberForSort(bMeta.Number)
		if aNum != bNum {
			return aNum < bNum
		}

		// Tertiary: Name (alphabetical)
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	// Build top nav items with base URL prefixing
	var items []TopNavItem
	for _, node := range rootItems {
		urlPath := tree.GetURLPath(node)
		items = append(items, TopNavItem{
			Name: node.Name,
			URL:  tree.PrefixURL(baseURL, urlPath),
		})
	}

	return items
}

// topNavNumberForSort returns the number for sorting, with nil treated as max
func topNavNumberForSort(n *int) int {
	if n == nil {
		return 999999 // Sort after numbered items
	}
	return *n
}

// RenderNavigationWithTopNav renders navigation excluding root files when top nav is enabled
func RenderNavigationWithTopNav(root *tree.Node, currentPath string, topNavItems []TopNavItem) template.HTML {
	return RenderNavigationWithTopNavAndBaseURL(root, currentPath, topNavItems, "")
}

// RenderNavigationWithTopNavAndBaseURL renders navigation with top nav filtering and base URL prefixing
func RenderNavigationWithTopNavAndBaseURL(root *tree.Node, currentPath string, topNavItems []TopNavItem, baseURL string) template.HTML {
	if len(topNavItems) == 0 {
		return RenderNavigationWithBaseURL(root, currentPath, baseURL)
	}

	// Create a set of URLs to exclude from sidebar (use unprefixed URLs for comparison)
	topNavURLs := make(map[string]bool)
	for _, item := range topNavItems {
		// Store the original URL path (without base prefix) for exclusion matching
		topNavURLs[item.URL] = true
	}

	var buf bytes.Buffer
	renderNavNodeFiltered(&buf, root.Children, currentPath, 0, topNavURLs, baseURL)
	return template.HTML(buf.String())
}

// renderNavNodeFiltered renders nav nodes, filtering out specified URLs
func renderNavNodeFiltered(buf *bytes.Buffer, nodes []*tree.Node, currentPath string, depth int, excludeURLs map[string]bool, baseURL string) {
	if len(nodes) == 0 {
		return
	}

	buf.WriteString("<ul role=\"tree\">\n")

	for _, node := range nodes {
		// Skip excluded files - compare against prefixed URL since that's what's in excludeURLs
		if !node.IsFolder {
			urlPath := tree.GetURLPath(node)
			prefixedURL := tree.PrefixURL(baseURL, urlPath)
			if excludeURLs[prefixedURL] {
				continue
			}
		}

		if node.IsFolder {
			renderFolderNode(buf, node, currentPath, depth, baseURL)
		} else {
			renderFileNode(buf, node, currentPath, baseURL)
		}
	}

	buf.WriteString("</ul>\n")
}

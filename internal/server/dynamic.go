// Package server provides HTTP file server and dynamic rendering functionality.
package server

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/wusher/volcano/internal/assets"
	"github.com/wusher/volcano/internal/autoindex"
	"github.com/wusher/volcano/internal/content"
	"github.com/wusher/volcano/internal/instant"
	"github.com/wusher/volcano/internal/markdown"
	"github.com/wusher/volcano/internal/navigation"
	"github.com/wusher/volcano/internal/styles"
	"github.com/wusher/volcano/internal/templates"
	"github.com/wusher/volcano/internal/toc"
	"github.com/wusher/volcano/internal/tree"
)

// DynamicConfig holds configuration for the dynamic server
type DynamicConfig struct {
	SourceDir   string
	Title       string
	Port        int
	Quiet       bool
	Verbose     bool
	TopNav      bool
	ShowPageNav bool
	Theme       string
	CSSPath     string
	AccentColor string // Custom accent color in hex format (e.g., "#ff6600")
	FaviconPath string // Path to favicon file
	InstantNav  bool   // Enable instant navigation with hover prefetching
}

// DynamicServer serves markdown files with live rendering
type DynamicServer struct {
	config       DynamicConfig
	renderer     *templates.Renderer
	transformer  *markdown.ContentTransformer
	writer       io.Writer
	server       *http.Server
	fs           FileSystem
	scanner      TreeScanner
	cssLoader    styles.CSSLoader
	faviconLinks template.HTML // Favicon HTML tags
	faviconData  []byte        // Favicon file content (in memory)
	faviconMime  string        // Favicon MIME type
	faviconName  string        // Favicon filename (e.g., "logo.png")
	instantNavJS template.JS   // Instant navigation JavaScript (if enabled)
}

// NewDynamicServer creates a new dynamic server
func NewDynamicServer(config DynamicConfig, writer io.Writer) (*DynamicServer, error) {
	cssConfig := styles.CSSConfig{
		Theme:       config.Theme,
		CSSPath:     config.CSSPath,
		AccentColor: config.AccentColor,
	}
	cssLoader := styles.NewCSSLoader(cssConfig, os.ReadFile)
	css, err := cssLoader.LoadCSS()
	if err != nil {
		return nil, fmt.Errorf("failed to load CSS: %w", err)
	}

	renderer, err := templates.NewRenderer(css)
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	srv := &DynamicServer{
		config:      config,
		renderer:    renderer,
		transformer: markdown.NewContentTransformer(""), // Dynamic server doesn't use site URL for external links
		writer:      writer,
		fs:          osFileSystem{},
		scanner:     defaultScanner{},
		cssLoader:   cssLoader,
	}

	// Initialize instant navigation JS if enabled
	if config.InstantNav {
		srv.instantNavJS = template.JS(instant.InstantNavJS)
	}

	// Process favicon if configured
	if config.FaviconPath != "" {
		if err := srv.loadFavicon(); err != nil {
			// Log warning but don't fail server startup
			_, _ = fmt.Fprintf(writer, "Warning: Failed to load favicon: %v\n", err)
		}
	}

	return srv, nil
}

// loadFavicon loads the favicon file into memory and generates HTML tags
func (s *DynamicServer) loadFavicon() error {
	// Validate file exists
	if _, err := os.Stat(s.config.FaviconPath); os.IsNotExist(err) {
		return fmt.Errorf("favicon file not found: %s", s.config.FaviconPath)
	}

	// Read file into memory
	data, err := os.ReadFile(s.config.FaviconPath)
	if err != nil {
		return fmt.Errorf("failed to read favicon: %w", err)
	}

	// Get filename and MIME type
	filename := filepath.Base(s.config.FaviconPath)
	mimeType := assets.GetFaviconMimeType(filename)
	if mimeType == "" {
		return fmt.Errorf("unsupported favicon format: %s", filename)
	}

	// Store in memory
	s.faviconData = data
	s.faviconMime = mimeType
	s.faviconName = filename

	// Generate HTML link tags
	links := []assets.FaviconLink{{
		Rel:  "icon",
		Type: mimeType,
		Href: "/" + filename,
	}}
	s.faviconLinks = assets.RenderFaviconLinks(links)

	return nil
}

// serveFavicon serves the favicon from memory
func (s *DynamicServer) serveFavicon(w http.ResponseWriter, urlPath string) bool {
	if s.faviconData == nil || s.faviconName == "" {
		return false
	}

	// Check if the request is for the favicon
	requestedFile := strings.TrimPrefix(urlPath, "/")
	if requestedFile != s.faviconName {
		return false
	}

	// Serve from memory
	w.Header().Set("Content-Type", s.faviconMime)
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	_, _ = w.Write(s.faviconData)
	return true
}

// getRenderer returns a renderer, re-reading CSS on each request for live reload
func (s *DynamicServer) getRenderer() (*templates.Renderer, error) {
	// Always reload CSS for live reload (works for both custom CSS and theme files in development)
	css, err := s.cssLoader.LoadCSS()
	if err != nil {
		// Fall back to cached renderer if CSS load fails
		s.logError("Failed to load CSS, using cached: %v", err)
		return s.renderer, nil
	}
	return templates.NewRenderer(css)
}

// WithFileSystem sets a custom FileSystem (for testing)
func (s *DynamicServer) WithFileSystem(fs FileSystem) *DynamicServer {
	s.fs = fs
	return s
}

// WithScanner sets a custom TreeScanner (for testing)
func (s *DynamicServer) WithScanner(scanner TreeScanner) *DynamicServer {
	s.scanner = scanner
	return s
}

// Handler returns the HTTP handler for this server
func (s *DynamicServer) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRequest)
	return mux
}

// Start starts the dynamic HTTP server and blocks until shutdown
func (s *DynamicServer) Start() error {
	handler := s.Handler()

	addr := fmt.Sprintf(":%d", s.config.Port)
	s.server = &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		s.log("Live server - changes are reflected immediately")
		s.log("Serving %s at http://localhost:%d", s.config.SourceDir, s.config.Port)
		s.log("Press Ctrl+C to stop")
		s.log("")
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-stop:
		s.log("")
		s.log("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.server.Shutdown(ctx)
	}
}

func (s *DynamicServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

	// Set cache control headers for development (no caching)
	rec.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	rec.Header().Set("Pragma", "no-cache")
	rec.Header().Set("Expires", "0")

	urlPath := r.URL.Path

	// Serve favicon from memory if requested
	if s.serveFavicon(rec, urlPath) {
		s.logRequest(r.Method, urlPath, rec.statusCode, time.Since(start))
		return
	}

	// Try to serve static files first (images, CSS, JS, etc.)
	if s.serveStaticFile(rec, r, urlPath) {
		s.logRequest(r.Method, urlPath, rec.statusCode, time.Since(start))
		return
	}

	// Try to render a markdown page
	if s.renderPage(rec, r, urlPath) {
		s.logRequest(r.Method, urlPath, rec.statusCode, time.Since(start))
		return
	}

	// Try to render an auto-generated folder index
	if s.tryAutoIndex(rec, urlPath) {
		s.logRequest(r.Method, urlPath, rec.statusCode, time.Since(start))
		return
	}

	// Serve 404
	s.serve404(rec, r)
	s.logRequest(r.Method, urlPath, rec.statusCode, time.Since(start))
}

// tryAutoIndex tries to render an auto-generated index for a folder
func (s *DynamicServer) tryAutoIndex(w http.ResponseWriter, urlPath string) bool {
	// Only handle paths that look like directories (ending with /)
	if !strings.HasSuffix(urlPath, "/") && urlPath != "" {
		return false
	}

	// Scan the tree
	site, err := s.scanner.Scan(s.config.SourceDir)
	if err != nil {
		return false
	}

	// Find the folder
	folderNode := findFolderByPath(site.Root, urlPath)
	if folderNode == nil {
		return false
	}

	// Check if it needs auto-index
	if !autoindex.NeedsAutoIndex(folderNode) {
		return false
	}

	// Render the auto-index
	return s.renderAutoIndex(w, urlPath, folderNode, site)
}

// serveStaticFile tries to serve a static file (non-markdown)
func (s *DynamicServer) serveStaticFile(w http.ResponseWriter, r *http.Request, urlPath string) bool {
	// Clean the path
	cleanPath := strings.TrimPrefix(urlPath, "/")
	if cleanPath == "" {
		return false // Let markdown handler deal with root
	}

	fullPath := filepath.Join(s.config.SourceDir, cleanPath)

	// Check if file exists and is not a directory
	stat, err := s.fs.Stat(fullPath)
	if err != nil || stat.IsDir() {
		return false
	}

	// Check if it's a markdown file - don't serve raw markdown
	if strings.HasSuffix(strings.ToLower(fullPath), ".md") {
		return false
	}

	// Serve the static file
	http.ServeFile(w, r, fullPath)
	return true
}

// renderPage tries to render a markdown page for the given URL
func (s *DynamicServer) renderPage(w http.ResponseWriter, _ *http.Request, urlPath string) bool {
	// Find the markdown file for this URL
	mdPath := s.resolveMarkdownPath(urlPath)
	if mdPath == "" {
		return false
	}

	fullMdPath := filepath.Join(s.config.SourceDir, mdPath)

	// Check if the file exists
	if _, err := s.fs.Stat(fullMdPath); err != nil {
		return false
	}

	// Scan the directory tree for navigation (fresh on every request)
	site, err := s.scanner.Scan(s.config.SourceDir)
	if err != nil {
		s.logError("Failed to scan directory: %v", err)
		return false
	}

	// Find the tree node for this page
	node := findNodeBySourcePath(site.Root, mdPath)
	if node == nil {
		return false
	}

	// Get paths
	outputPath := tree.GetOutputPath(node)
	nodeURLPath := tree.GetURLPath(node)

	// Compute source directory for wikilink resolution
	relDir := filepath.Dir(node.Path)
	sourceDir := "/"
	if relDir != "." && relDir != "" {
		sourceDir = "/" + tree.SlugifyPath(relDir) + "/"
	}

	// Read and transform the markdown file
	mdContent, err := s.fs.ReadFile(fullMdPath)
	if err != nil {
		s.logError("Failed to read markdown: %v", err)
		return false
	}

	// Transform markdown to HTML with all enhancements
	page, err := s.transformer.TransformMarkdown(
		mdContent,
		sourceDir,
		fullMdPath,
		outputPath,
		nodeURLPath,
		node.Name,
	)
	if err != nil {
		s.logError("Failed to parse markdown: %v", err)
		return false
	}

	htmlContent := page.Content

	// Validate internal links (no base URL for dev server)
	validURLs := tree.BuildValidURLMap(site, "")
	brokenLinks := markdown.ValidateLinks(htmlContent, nodeURLPath, validURLs)
	if len(brokenLinks) > 0 {
		s.serveBrokenLinksError(w, nodeURLPath, brokenLinks, site)
		return true // We handled the request (with an error page)
	}

	// Calculate reading time
	rt := content.CalculateReadingTime(htmlContent)
	readingTime := content.FormatReadingTime(rt)

	// Build breadcrumbs
	breadcrumbs := navigation.BuildBreadcrumbs(node, s.config.Title)
	breadcrumbsHTML := navigation.RenderBreadcrumbs(breadcrumbs)

	// Build page navigation (only if enabled)
	var pageNavHTML template.HTML
	if s.config.ShowPageNav {
		allPages := collectAllPages(site.Root)
		pageNav := navigation.BuildPageNavigation(node, allPages)
		pageNavHTML = navigation.RenderPageNavigation(pageNav)
	}

	// Extract TOC
	pageTOC := toc.ExtractTOC(htmlContent, 3)
	tocHTML := toc.RenderTOC(pageTOC)
	hasTOC := pageTOC != nil && len(pageTOC.Items) > 0

	// Build top nav items if enabled
	topNavItems := templates.BuildTopNavItems(site.Root, s.config.TopNav)

	// Render navigation (filtered when top nav is enabled)
	nav := templates.RenderNavigationWithTopNav(site.Root, nodeURLPath, topNavItems)

	// Prepare template data
	data := templates.PageData{
		SiteTitle:    s.config.Title,
		PageTitle:    page.Title,
		Content:      template.HTML(htmlContent),
		Navigation:   nav,
		CurrentPath:  nodeURLPath,
		Breadcrumbs:  breadcrumbsHTML,
		PageNav:      pageNavHTML,
		TOC:          tocHTML,
		FaviconLinks: s.faviconLinks,
		ReadingTime:  readingTime,
		HasTOC:       hasTOC,
		ShowSearch:   true,
		TopNavItems:  topNavItems,
		InstantNavJS: s.instantNavJS,
	}

	// Get renderer (re-reads CSS if using custom CSS file)
	renderer, err := s.getRenderer()
	if err != nil {
		s.logError("Failed to get renderer: %v", err)
		return false
	}

	// Render the page
	var buf bytes.Buffer
	if err := renderer.Render(&buf, data); err != nil {
		s.logError("Failed to render page: %v", err)
		return false
	}

	// Write response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
	return true
}

// resolveMarkdownPath resolves a URL path to a markdown file path
func (s *DynamicServer) resolveMarkdownPath(urlPath string) string {
	// Remove leading slash
	urlPath = strings.TrimPrefix(urlPath, "/")

	// Remove trailing slash for processing
	urlPath = strings.TrimSuffix(urlPath, "/")

	// Root path - look for index.md
	if urlPath == "" {
		fullPath := filepath.Join(s.config.SourceDir, "index.md")
		if _, err := s.fs.Stat(fullPath); err == nil {
			return "index.md"
		}
		return ""
	}

	// Find the actual filesystem path from the slugified URL
	actualDir := s.findActualDir(urlPath)

	// Try as a file with .md extension (clean URLs: /about/ -> about.md)
	// First try the URL path directly (for non-prefixed files)
	mdPath := urlPath + ".md"
	fullPath := filepath.Join(s.config.SourceDir, mdPath)
	if _, err := s.fs.Stat(fullPath); err == nil {
		return mdPath
	}

	// Try as directory with index.md using actual filesystem path
	if actualDir != "" {
		// Search for any index file (case-insensitive)
		dirPath := filepath.Join(s.config.SourceDir, actualDir)
		if entries, err := os.ReadDir(dirPath); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				name := entry.Name()
				baseName := strings.TrimSuffix(name, filepath.Ext(name))
				if strings.EqualFold(baseName, "index") || strings.EqualFold(baseName, "readme") {
					return filepath.Join(actualDir, name)
				}
			}
		}
	}

	// Try to find a file with date/number prefix
	// e.g., /posts/asset-helpers/ might map to posts/2024-01-09-asset-helpers.md
	dir := filepath.Dir(urlPath)
	slug := filepath.Base(urlPath)
	if found := s.findPrefixedFile(dir, slug); found != "" {
		return found
	}

	return ""
}

// findPrefixedFile searches for a markdown file with date/number prefixes
// that matches the given slug in the directory (dir is a slugified URL path)
func (s *DynamicServer) findPrefixedFile(slugDir, slug string) string {
	// Find the actual filesystem directory from the slugified path
	actualDir := s.findActualDir(slugDir)
	if actualDir == "" && slugDir != "." && slugDir != "" {
		return ""
	}

	searchDir := s.config.SourceDir
	if actualDir != "" {
		searchDir = filepath.Join(s.config.SourceDir, actualDir)
	}

	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".md") {
			continue
		}

		// Extract metadata to get the slug
		info, err := entry.Info()
		if err != nil {
			continue
		}
		meta := tree.ExtractFileMetadata(name, info.ModTime())

		// Check if the slug matches
		if meta.Slug == slug {
			if actualDir == "" {
				return name
			}
			return filepath.Join(actualDir, name)
		}
	}

	return ""
}

// findActualDir finds the actual filesystem directory path from a slugified URL path
// e.g., "inbox/health" → "0. Inbox/1. Health"
func (s *DynamicServer) findActualDir(slugPath string) string {
	if slugPath == "" || slugPath == "." {
		return ""
	}

	segments := strings.Split(slugPath, "/")
	currentPath := s.config.SourceDir
	var actualSegments []string

	for _, slugSeg := range segments {
		if slugSeg == "" {
			continue
		}

		entries, err := os.ReadDir(currentPath)
		if err != nil {
			return ""
		}

		found := false
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			// Check if this directory's slugified name matches
			if tree.Slugify(entry.Name()) == slugSeg {
				actualSegments = append(actualSegments, entry.Name())
				currentPath = filepath.Join(currentPath, entry.Name())
				found = true
				break
			}
		}

		if !found {
			return ""
		}
	}

	return filepath.Join(actualSegments...)
}

// findNodeBySourcePath finds a node in the tree by its source path
func findNodeBySourcePath(node *tree.Node, sourcePath string) *tree.Node {
	if node == nil {
		return nil
	}

	// Check if this node matches (case-insensitive for cross-platform compatibility)
	if strings.EqualFold(node.Path, sourcePath) {
		return node
	}

	// For folders with index, check the index path
	if node.IsFolder && node.HasIndex && strings.EqualFold(node.IndexPath, sourcePath) {
		return &tree.Node{
			Path:       node.IndexPath,
			Name:       node.Name,
			SourcePath: filepath.Join(node.SourcePath, "index.md"),
		}
	}

	// Search children
	for _, child := range node.Children {
		if found := findNodeBySourcePath(child, sourcePath); found != nil {
			return found
		}
	}

	return nil
}

// serveBrokenLinksError renders an error page showing broken internal links
func (s *DynamicServer) serveBrokenLinksError(w http.ResponseWriter, sourcePage string, brokenLinks []markdown.BrokenLink, site *tree.Site) {
	// Log the error
	s.logError("Page %s has %d broken internal links:", sourcePage, len(brokenLinks))
	for _, bl := range brokenLinks {
		s.logError("  -> %s", bl.LinkURL)
	}

	// Build error content
	var sb strings.Builder
	sb.WriteString(`<h1>Build Error: Broken Links</h1>`)
	sb.WriteString("\n")
	sb.WriteString(`<p>The page <code>`)
	sb.WriteString(template.HTMLEscapeString(sourcePage))
	sb.WriteString(`</code> contains broken internal links:</p>`)
	sb.WriteString("\n")
	sb.WriteString(`<ul class="broken-links-list">`)
	sb.WriteString("\n")
	for _, bl := range brokenLinks {
		sb.WriteString(`<li><code>`)
		sb.WriteString(template.HTMLEscapeString(bl.LinkURL))
		sb.WriteString(`</code></li>`)
		sb.WriteString("\n")
	}
	sb.WriteString(`</ul>`)
	sb.WriteString("\n")
	sb.WriteString(`<p>Please fix these links to continue.</p>`)

	// Render navigation
	var nav template.HTML
	if site != nil {
		nav = templates.RenderNavigation(site.Root, "")
	}

	data := templates.PageData{
		SiteTitle:    s.config.Title,
		PageTitle:    "Build Error",
		Content:      template.HTML(sb.String()),
		Navigation:   nav,
		CurrentPath:  "",
		FaviconLinks: s.faviconLinks,
		InstantNavJS: s.instantNavJS,
	}

	// Get renderer
	renderer, err := s.getRenderer()
	if err != nil {
		http.Error(w, "Build Error: Broken Links", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	var buf bytes.Buffer
	if err := renderer.Render(&buf, data); err != nil {
		http.Error(w, "Build Error: Broken Links", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// serve404 renders a 404 error page
func (s *DynamicServer) serve404(w http.ResponseWriter, _ *http.Request) {
	// Try to scan for navigation
	var nav template.HTML
	site, err := s.scanner.Scan(s.config.SourceDir)
	if err == nil {
		nav = templates.RenderNavigation(site.Root, "")
	}

	content := `<h1>404 - Page Not Found</h1>
<p>The page you're looking for doesn't exist.</p>
<p><a href="/">Return to home</a></p>`

	data := templates.PageData{
		SiteTitle:    s.config.Title,
		PageTitle:    "Page Not Found",
		Content:      template.HTML(content),
		Navigation:   nav,
		CurrentPath:  "",
		FaviconLinks: s.faviconLinks,
		InstantNavJS: s.instantNavJS,
	}

	// Get renderer (re-reads CSS if using custom CSS file)
	renderer, err := s.getRenderer()
	if err != nil {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	var buf bytes.Buffer
	if err := renderer.Render(&buf, data); err != nil {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

// log prints a message if not in quiet mode
func (s *DynamicServer) log(format string, args ...interface{}) {
	if !s.config.Quiet {
		_, _ = fmt.Fprintf(s.writer, format+"\n", args...)
	}
}

// logError logs an error message
func (s *DynamicServer) logError(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(s.writer, "Error: "+format+"\n", args...)
}

// logRequest logs an HTTP request
func (s *DynamicServer) logRequest(method, path string, status int, duration time.Duration) {
	if !s.config.Quiet {
		statusColor := ""
		resetColor := ""

		if status >= 200 && status < 300 {
			statusColor = "\033[32m" // Green
			resetColor = "\033[0m"
		} else if status >= 400 {
			statusColor = "\033[31m" // Red
			resetColor = "\033[0m"
		}

		_, _ = fmt.Fprintf(s.writer, "%s%-4s%s %s %s%d%s %s\n",
			statusColor, method, resetColor,
			path,
			statusColor, status, resetColor,
			duration.Round(time.Millisecond),
		)
	}
}

// collectAllPages collects all non-folder nodes from the tree for prev/next navigation
func collectAllPages(node *tree.Node) []*tree.Node {
	var pages []*tree.Node
	collectPagesRecursive(node, &pages)
	return pages
}

func collectPagesRecursive(node *tree.Node, pages *[]*tree.Node) {
	if node == nil {
		return
	}

	// Add non-folder nodes (but skip the root if it has no path)
	if !node.IsFolder && node.Path != "" {
		*pages = append(*pages, node)
	}

	// If folder has an index, add a virtual node for it
	if node.IsFolder && node.HasIndex {
		indexNode := &tree.Node{
			Path:       node.IndexPath,
			Name:       node.Name,
			SourcePath: filepath.Join(node.SourcePath, "index.md"),
		}
		*pages = append(*pages, indexNode)
	}

	// Recurse into children
	for _, child := range node.Children {
		collectPagesRecursive(child, pages)
	}
}

// renderAutoIndex renders an auto-generated index page for a folder
func (s *DynamicServer) renderAutoIndex(w http.ResponseWriter, urlPath string, node *tree.Node, site *tree.Site) bool {
	// Build index using shared autoindex package
	index := autoindex.Build(node)
	htmlContent := autoindex.RenderContent(index)

	// Build breadcrumbs
	breadcrumbs := navigation.BuildBreadcrumbs(node, s.config.Title)
	breadcrumbsHTML := navigation.RenderBreadcrumbs(breadcrumbs)

	// Build top nav items if enabled
	topNavItems := templates.BuildTopNavItems(site.Root, s.config.TopNav)

	// Render navigation (filtered when top nav is enabled)
	nav := templates.RenderNavigationWithTopNav(site.Root, urlPath, topNavItems)

	// Prepare template data
	data := templates.PageData{
		SiteTitle:    s.config.Title,
		PageTitle:    node.Name,
		Content:      htmlContent,
		Navigation:   nav,
		CurrentPath:  urlPath,
		Breadcrumbs:  breadcrumbsHTML,
		FaviconLinks: s.faviconLinks,
		ShowSearch:   true,
		TopNavItems:  topNavItems,
		InstantNavJS: s.instantNavJS,
	}

	// Get renderer (re-reads CSS if using custom CSS file)
	renderer, err := s.getRenderer()
	if err != nil {
		s.logError("Failed to get renderer: %v", err)
		return false
	}

	// Render the page
	var buf bytes.Buffer
	if err := renderer.Render(&buf, data); err != nil {
		s.logError("Failed to render auto-index: %v", err)
		return false
	}

	// Write response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
	return true
}

// findFolderByPath finds a folder node by its URL path (slugified)
func findFolderByPath(node *tree.Node, urlPath string) *tree.Node {
	if node == nil {
		return nil
	}

	// Normalize the URL path
	urlPath = strings.TrimPrefix(urlPath, "/")
	urlPath = strings.TrimSuffix(urlPath, "/")

	// Check if this folder matches (compare slugified paths)
	if node.IsFolder && tree.SlugifyPath(node.Path) == urlPath {
		return node
	}

	// Search children
	for _, child := range node.Children {
		if found := findFolderByPath(child, urlPath); found != nil {
			return found
		}
	}

	return nil
}

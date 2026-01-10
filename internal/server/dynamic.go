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
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/wusher/volcano/internal/content"
	"github.com/wusher/volcano/internal/markdown"
	"github.com/wusher/volcano/internal/navigation"
	"github.com/wusher/volcano/internal/styles"
	"github.com/wusher/volcano/internal/templates"
	"github.com/wusher/volcano/internal/toc"
	"github.com/wusher/volcano/internal/tree"
)

// buildValidURLMap creates a map of all valid URLs from the site
func buildValidURLMap(site *tree.Site) map[string]bool {
	validURLs := make(map[string]bool)
	validURLs["/"] = true

	// Add all page URLs
	for _, node := range site.AllPages {
		urlPath := tree.GetURLPath(node)
		if urlPath != "" {
			validURLs[urlPath] = true
		}
	}

	// Add folder URLs (for auto-index)
	addFolderURLs(site.Root, validURLs)

	return validURLs
}

// addFolderURLs recursively adds folder URLs to the map
func addFolderURLs(node *tree.Node, validURLs map[string]bool) {
	if node == nil {
		return
	}

	if node.IsFolder && node.Path != "" {
		urlPath := "/" + tree.SlugifyPath(node.Path) + "/"
		validURLs[urlPath] = true
	}

	for _, child := range node.Children {
		addFolderURLs(child, validURLs)
	}
}

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
	DevMode     bool // Enable keyboard shortcuts for toggling settings
}

// DynamicServer serves markdown files with live rendering
type DynamicServer struct {
	config      DynamicConfig
	renderer    *templates.Renderer
	writer      io.Writer
	server      *http.Server
	fs          FileSystem
	scanner     TreeScanner
	sse         *SSEBroadcaster
	keyboard    *KeyboardHandler
	validThemes []string
}

// NewDynamicServer creates a new dynamic server
func NewDynamicServer(config DynamicConfig, writer io.Writer) (*DynamicServer, error) {
	css, err := getCSSContent(config)
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
		writer:      writer,
		fs:          osFileSystem{},
		scanner:     defaultScanner{},
		sse:         NewSSEBroadcaster(),
		validThemes: []string{"docs", "blog", "vanilla"},
	}

	// Initialize keyboard handler for dev mode
	if config.DevMode {
		srv.keyboard = NewKeyboardHandler(writer, srv.handleKeyPress)
	}

	return srv, nil
}

// getCSSContent returns minified CSS from custom file or embedded theme
func getCSSContent(config DynamicConfig) (string, error) {
	var css string
	if config.CSSPath != "" {
		content, err := os.ReadFile(config.CSSPath)
		if err != nil {
			return "", err
		}
		css = string(content)
	} else {
		css = styles.GetCSS(config.Theme)
	}
	return styles.MinifyCSS(css)
}

// getRenderer returns a renderer, re-reading CSS file if using custom CSS
func (s *DynamicServer) getRenderer() (*templates.Renderer, error) {
	// If using custom CSS, re-read on each request for live reload
	if s.config.CSSPath != "" {
		css, err := getCSSContent(s.config)
		if err != nil {
			// Fall back to cached renderer if file read fails
			s.logError("Failed to read CSS file, using cached: %v", err)
			return s.renderer, nil
		}
		return templates.NewRenderer(css)
	}
	// Use cached renderer for embedded themes
	return s.renderer, nil
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
	mux.HandleFunc("/__volcano/events", s.sse.Handler())
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

	// Start SSE broadcaster
	s.sse.Start()
	defer s.sse.Stop()

	// Start keyboard handler if in dev mode
	if s.keyboard != nil {
		if err := s.keyboard.Start(); err != nil {
			s.logError("Keyboard handler failed: %v", err)
		} else {
			defer s.keyboard.Stop()
		}
	}

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		s.log("Live server - changes are reflected immediately")
		s.log("Serving %s at http://localhost:%d", s.config.SourceDir, s.config.Port)
		if s.config.DevMode {
			s.log("")
			s.printDevModeHelp()
			s.printCurrentSettings()
		} else {
			s.log("Press Ctrl+C to stop")
		}
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
	if !needsAutoIndex(folderNode) {
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

	// Parse the markdown file
	page, err := markdown.ParseFile(
		fullMdPath,
		outputPath,
		nodeURLPath,
		sourceDir,
		node.Name,
	)
	if err != nil {
		s.logError("Failed to parse markdown: %v", err)
		return false
	}

	// Process the HTML content
	htmlContent := page.Content

	// Wrap code blocks with copy button
	htmlContent = markdown.WrapCodeBlocks(htmlContent)

	// Validate internal links
	validURLs := buildValidURLMap(site)
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
		SiteTitle:   s.config.Title,
		PageTitle:   page.Title,
		Content:     template.HTML(htmlContent),
		Navigation:  nav,
		CurrentPath: nodeURLPath,
		Breadcrumbs: breadcrumbsHTML,
		PageNav:     pageNavHTML,
		TOC:         tocHTML,
		ReadingTime: readingTime,
		HasTOC:      hasTOC,
		ShowSearch:  true,
		TopNavItems: topNavItems,
		DevMode:     s.config.DevMode,
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
		SiteTitle:   s.config.Title,
		PageTitle:   "Build Error",
		Content:     template.HTML(sb.String()),
		Navigation:  nav,
		CurrentPath: "",
		DevMode:     s.config.DevMode,
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
		SiteTitle:   s.config.Title,
		PageTitle:   "Page Not Found",
		Content:     template.HTML(content),
		Navigation:  nav,
		CurrentPath: "",
		DevMode:     s.config.DevMode,
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

// handleKeyPress handles keyboard input for toggling settings
func (s *DynamicServer) handleKeyPress(key rune) {
	switch key {
	case '1':
		s.setTheme("docs")
	case '2':
		s.setTheme("blog")
	case '3':
		s.setTheme("vanilla")
	case 't', 'T':
		s.cycleTheme()
	case 'n', 'N':
		s.config.TopNav = !s.config.TopNav
		s.notifyReload("topNav", s.config.TopNav)
		s.printSettingChange("top-nav", s.config.TopNav)
	case 'p', 'P':
		s.config.ShowPageNav = !s.config.ShowPageNav
		s.notifyReload("pageNav", s.config.ShowPageNav)
		s.printSettingChange("page-nav", s.config.ShowPageNav)
	case 'h', 'H', '?':
		s.printDevModeHelp()
		s.printCurrentSettings()
	case 'r', 'R':
		s.notifyReload("reload", nil)
		s.logRaw("\r\033[K\033[33m↻\033[0m Reload triggered")
	}
}

// setTheme sets the theme and notifies clients
func (s *DynamicServer) setTheme(theme string) {
	if s.config.Theme == theme {
		return
	}
	s.config.Theme = theme
	s.updateRenderer()
	s.notifyReload("theme", theme)
	s.logRaw("\r\033[K\033[36m◆\033[0m Theme: %s", theme)
}

// cycleTheme cycles through available themes
func (s *DynamicServer) cycleTheme() {
	currentIdx := 0
	for i, t := range s.validThemes {
		if t == s.config.Theme {
			currentIdx = i
			break
		}
	}
	nextIdx := (currentIdx + 1) % len(s.validThemes)
	s.setTheme(s.validThemes[nextIdx])
}

// updateRenderer updates the renderer with new CSS
func (s *DynamicServer) updateRenderer() {
	css, err := getCSSContent(s.config)
	if err != nil {
		s.logError("Failed to load CSS: %v", err)
		return
	}
	renderer, err := templates.NewRenderer(css)
	if err != nil {
		s.logError("Failed to create renderer: %v", err)
		return
	}
	s.renderer = renderer
}

// notifyReload sends a reload notification to all connected browsers
func (s *DynamicServer) notifyReload(eventType string, data interface{}) {
	s.sse.Broadcast(eventType, map[string]interface{}{
		"type":  eventType,
		"value": data,
	})
}

// printDevModeHelp prints the dev mode help
func (s *DynamicServer) printDevModeHelp() {
	s.logRaw("\033[1mDev Mode Shortcuts:\033[0m")
	s.logRaw("  \033[33mt\033[0m       Cycle theme (docs → blog → vanilla)")
	s.logRaw("  \033[33m1/2/3\033[0m   Switch to docs/blog/vanilla theme")
	s.logRaw("  \033[33mn\033[0m       Toggle top navigation")
	s.logRaw("  \033[33mp\033[0m       Toggle page navigation")
	s.logRaw("  \033[33mr\033[0m       Force reload all browsers")
	s.logRaw("  \033[33mh/?/\033[0m    Show this help")
	s.logRaw("  \033[33mCtrl+C\033[0m  Stop server")
}

// printCurrentSettings prints the current settings
func (s *DynamicServer) printCurrentSettings() {
	s.logRaw("")
	s.logRaw("\033[1mCurrent Settings:\033[0m")
	s.logRaw("  theme:     \033[36m%s\033[0m", s.config.Theme)
	s.logRaw("  top-nav:   %s", s.formatBool(s.config.TopNav))
	s.logRaw("  page-nav:  %s", s.formatBool(s.config.ShowPageNav))
}

// printSettingChange prints a setting change notification
func (s *DynamicServer) printSettingChange(setting string, value bool) {
	s.logRaw("\r\033[K\033[36m◆\033[0m %s: %s", setting, s.formatBool(value))
}

// formatBool formats a boolean for display
func (s *DynamicServer) formatBool(b bool) string {
	if b {
		return "\033[32mon\033[0m"
	}
	return "\033[90moff\033[0m"
}

// logRaw prints a message without newline handling (for raw terminal output)
func (s *DynamicServer) logRaw(format string, args ...interface{}) {
	if !s.config.Quiet {
		_, _ = fmt.Fprintf(s.writer, format+"\n", args...)
	}
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

// AutoIndexItem represents an item in an auto-generated folder index
type AutoIndexItem struct {
	Title    string
	URL      string
	IsFolder bool
}

// renderAutoIndex renders an auto-generated index page for a folder
func (s *DynamicServer) renderAutoIndex(w http.ResponseWriter, urlPath string, node *tree.Node, site *tree.Site) bool {
	// Build list of children
	var items []AutoIndexItem
	for _, child := range node.Children {
		url := tree.GetURLPath(child)
		if child.IsFolder {
			url = "/" + tree.SlugifyPath(child.Path) + "/"
		}
		items = append(items, AutoIndexItem{
			Title:    child.Name,
			URL:      url,
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

	// Build HTML content
	var sb strings.Builder
	sb.WriteString(`<article class="auto-index-page">`)
	sb.WriteString("\n")
	sb.WriteString(`<h1>`)
	sb.WriteString(template.HTMLEscapeString(node.Name))
	sb.WriteString(`</h1>`)
	sb.WriteString("\n")

	if len(items) > 0 {
		sb.WriteString(`<ul class="folder-index">`)
		sb.WriteString("\n")
		for _, item := range items {
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
	htmlContent := sb.String()

	// Build breadcrumbs
	breadcrumbs := navigation.BuildBreadcrumbs(node, s.config.Title)
	breadcrumbsHTML := navigation.RenderBreadcrumbs(breadcrumbs)

	// Build top nav items if enabled
	topNavItems := templates.BuildTopNavItems(site.Root, s.config.TopNav)

	// Render navigation (filtered when top nav is enabled)
	nav := templates.RenderNavigationWithTopNav(site.Root, urlPath, topNavItems)

	// Prepare template data
	data := templates.PageData{
		SiteTitle:   s.config.Title,
		PageTitle:   node.Name,
		Content:     template.HTML(htmlContent),
		Navigation:  nav,
		CurrentPath: urlPath,
		Breadcrumbs: breadcrumbsHTML,
		ShowSearch:  true,
		TopNavItems: topNavItems,
		DevMode:     s.config.DevMode,
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

// needsAutoIndex checks if a folder needs an auto-generated index
func needsAutoIndex(node *tree.Node) bool {
	if !node.IsFolder {
		return false
	}

	// Already has an index
	if node.HasIndex {
		return false
	}

	// Check for index.md in children
	for _, child := range node.Children {
		if !child.IsFolder {
			baseName := strings.TrimSuffix(filepath.Base(child.Path), filepath.Ext(child.Path))
			lower := strings.ToLower(baseName)
			if lower == "index" || lower == "readme" {
				return false
			}
		}
	}

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

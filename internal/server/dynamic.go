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

	"volcano/internal/markdown"
	"volcano/internal/styles"
	"volcano/internal/templates"
	"volcano/internal/tree"
)

// DynamicConfig holds configuration for the dynamic server
type DynamicConfig struct {
	SourceDir string
	Title     string
	Port      int
	Quiet     bool
	Verbose   bool
}

// DynamicServer serves markdown files with live rendering
type DynamicServer struct {
	config   DynamicConfig
	renderer *templates.Renderer
	writer   io.Writer
	server   *http.Server
	fs       FileSystem
	scanner  TreeScanner
}

// NewDynamicServer creates a new dynamic server
func NewDynamicServer(config DynamicConfig, writer io.Writer) (*DynamicServer, error) {
	renderer, err := templates.NewRenderer(styles.GetCSS())
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	return &DynamicServer{
		config:   config,
		renderer: renderer,
		writer:   writer,
		fs:       osFileSystem{},
		scanner:  defaultScanner{},
	}, nil
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

	// Serve 404
	s.serve404(rec, r)
	s.logRequest(r.Method, urlPath, rec.statusCode, time.Since(start))
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

	// Parse the markdown file
	page, err := markdown.ParseFile(
		fullMdPath,
		outputPath,
		nodeURLPath,
		node.Name,
	)
	if err != nil {
		s.logError("Failed to parse markdown: %v", err)
		return false
	}

	// Render navigation
	nav := templates.RenderNavigation(site.Root, nodeURLPath)

	// Prepare template data
	data := templates.PageData{
		SiteTitle:   s.config.Title,
		PageTitle:   page.Title,
		Content:     template.HTML(page.Content),
		Navigation:  nav,
		CurrentPath: nodeURLPath,
	}

	// Render the page
	var buf bytes.Buffer
	if err := s.renderer.Render(&buf, data); err != nil {
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

	// Try as a file with .md extension (clean URLs: /about/ -> about.md)
	mdPath := urlPath + ".md"
	fullPath := filepath.Join(s.config.SourceDir, mdPath)
	if _, err := s.fs.Stat(fullPath); err == nil {
		return mdPath
	}

	// Try as directory with index.md
	indexPath := filepath.Join(urlPath, "index.md")
	fullPath = filepath.Join(s.config.SourceDir, indexPath)
	if _, err := s.fs.Stat(fullPath); err == nil {
		return indexPath
	}

	return ""
}

// findNodeBySourcePath finds a node in the tree by its source path
func findNodeBySourcePath(node *tree.Node, sourcePath string) *tree.Node {
	if node == nil {
		return nil
	}

	// Check if this node matches
	if node.Path == sourcePath {
		return node
	}

	// For folders with index, check the index path
	if node.IsFolder && node.HasIndex && node.IndexPath == sourcePath {
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
	}

	w.WriteHeader(http.StatusNotFound)
	var buf bytes.Buffer
	if err := s.renderer.Render(&buf, data); err != nil {
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

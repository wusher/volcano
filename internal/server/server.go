// Package server provides an HTTP file server for previewing generated sites.
package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Config holds server configuration
type Config struct {
	Dir     string
	Port    int
	Quiet   bool
	Verbose bool
}

// Server is an HTTP file server with clean URL support
type Server struct {
	config Config
	writer io.Writer
	server *http.Server
}

// New creates a new Server
func New(config Config, writer io.Writer) *Server {
	return &Server{
		config: config,
		writer: writer,
	}
}

// Handler returns the HTTP handler for this server
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRequest)
	return mux
}

// Start starts the HTTP server and blocks until shutdown
func (s *Server) Start() error {
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
		s.log("Serving %s at http://localhost:%d", s.config.Dir, s.config.Port)
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

// responseRecorder wraps http.ResponseWriter to capture status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

	// Set cache control headers for development
	rec.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	rec.Header().Set("Pragma", "no-cache")
	rec.Header().Set("Expires", "0")

	// Resolve the file path
	urlPath := r.URL.Path
	filePath := s.resolvePath(urlPath)

	// Check if file exists
	fullPath := filepath.Join(s.config.Dir, filePath)
	stat, err := os.Stat(fullPath)

	if err != nil || stat.IsDir() {
		// Try to serve 404.html
		s.serve404(rec, r)
	} else {
		// Serve the file
		http.ServeFile(rec, r, fullPath)
	}

	// Log the request
	duration := time.Since(start)
	s.logRequest(r.Method, urlPath, rec.statusCode, duration)
}

// resolvePath resolves a URL path to a file path with clean URL support
func (s *Server) resolvePath(urlPath string) string {
	// Remove leading slash
	urlPath = strings.TrimPrefix(urlPath, "/")

	// If it's empty, serve index.html
	if urlPath == "" {
		return "index.html"
	}

	// If it ends with a slash, serve index.html from that directory
	if strings.HasSuffix(urlPath, "/") {
		return filepath.Join(urlPath, "index.html")
	}

	// Check if the path exists as-is
	fullPath := filepath.Join(s.config.Dir, urlPath)
	if stat, err := os.Stat(fullPath); err == nil && !stat.IsDir() {
		return urlPath
	}

	// Check if it's a directory with index.html
	indexPath := filepath.Join(urlPath, "index.html")
	fullIndexPath := filepath.Join(s.config.Dir, indexPath)
	if stat, err := os.Stat(fullIndexPath); err == nil && !stat.IsDir() {
		return indexPath
	}

	// Check if adding .html extension works
	htmlPath := urlPath + ".html"
	fullHTMLPath := filepath.Join(s.config.Dir, htmlPath)
	if stat, err := os.Stat(fullHTMLPath); err == nil && !stat.IsDir() {
		return htmlPath
	}

	// Return original path (will result in 404)
	return urlPath
}

// serve404 serves the 404.html page if it exists
func (s *Server) serve404(w http.ResponseWriter, r *http.Request) {
	notFoundPath := filepath.Join(s.config.Dir, "404.html")
	if stat, err := os.Stat(notFoundPath); err == nil && !stat.IsDir() {
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, notFoundPath)
		return
	}

	// Fallback to simple 404 response
	http.Error(w, "404 - Page Not Found", http.StatusNotFound)
}

// log prints a message if not in quiet mode
func (s *Server) log(format string, args ...interface{}) {
	if !s.config.Quiet {
		_, _ = fmt.Fprintf(s.writer, format+"\n", args...)
	}
}

// logRequest logs an HTTP request
func (s *Server) logRequest(method, path string, status int, duration time.Duration) {
	if !s.config.Quiet {
		statusColor := ""
		resetColor := ""

		// Add color codes if writing to a terminal
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

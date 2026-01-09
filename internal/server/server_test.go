package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"index.html":              "<html>Home</html>",
		"about.html":              "<html>About</html>",
		"guides/index.html":       "<html>Guides</html>",
		"guides/intro/index.html": "<html>Intro</html>",
		"assets/style.css":        "body {}",
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	s := New(Config{Dir: tmpDir, Port: 8080}, io.Discard)

	tests := []struct {
		name     string
		urlPath  string
		expected string
	}{
		{"root", "/", "index.html"},
		{"root index", "/index.html", "index.html"},
		{"page", "/about", "about.html"},
		{"page with extension", "/about.html", "about.html"},
		{"directory", "/guides/", "guides/index.html"},
		{"directory without slash", "/guides", "guides/index.html"},
		{"nested directory", "/guides/intro/", "guides/intro/index.html"},
		{"asset", "/assets/style.css", "assets/style.css"},
		{"missing file", "/nonexistent", "nonexistent"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.resolvePath(tt.urlPath)
			if result != tt.expected {
				t.Errorf("resolvePath(%q) = %q, want %q", tt.urlPath, result, tt.expected)
			}
		})
	}
}

func TestHandleRequest(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"index.html":        "<!DOCTYPE html><html><body>Home</body></html>",
		"guides/index.html": "<!DOCTYPE html><html><body>Guides</body></html>",
		"404.html":          "<!DOCTYPE html><html><body>Not Found</body></html>",
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	var buf bytes.Buffer
	s := New(Config{Dir: tmpDir, Port: 8080}, &buf)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{"home page", "/", http.StatusOK, "Home"},
		{"guides page", "/guides/", http.StatusOK, "Guides"},
		{"guides no slash", "/guides", http.StatusOK, "Guides"},
		{"404 page", "/nonexistent", http.StatusNotFound, "Not Found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rec := httptest.NewRecorder()

			s.handleRequest(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Status = %d, want %d", rec.Code, tt.expectedStatus)
			}

			body := rec.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("Body = %q, want to contain %q", body, tt.expectedBody)
			}

			// Check cache control headers
			cacheControl := rec.Header().Get("Cache-Control")
			if !strings.Contains(cacheControl, "no-cache") {
				t.Errorf("Cache-Control should contain no-cache, got %q", cacheControl)
			}
		})
	}
}

func TestHandleRequestWithout404Page(t *testing.T) {
	tmpDir := t.TempDir()

	// Create only index.html, no 404.html
	indexPath := filepath.Join(tmpDir, "index.html")
	if err := os.WriteFile(indexPath, []byte("<html>Home</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	s := New(Config{Dir: tmpDir, Port: 8080, Quiet: true}, io.Discard)

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	rec := httptest.NewRecorder()

	s.handleRequest(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "404") {
		t.Errorf("Body should contain '404', got %q", body)
	}
}

func TestLogging(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.html
	indexPath := filepath.Join(tmpDir, "index.html")
	if err := os.WriteFile(indexPath, []byte("<html>Home</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	s := New(Config{Dir: tmpDir, Port: 8080}, &buf)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	s.handleRequest(rec, req)

	output := buf.String()
	if !strings.Contains(output, "GET") {
		t.Errorf("Log should contain method 'GET', got %q", output)
	}
	if !strings.Contains(output, "/") {
		t.Errorf("Log should contain path '/', got %q", output)
	}
	if !strings.Contains(output, "200") {
		t.Errorf("Log should contain status '200', got %q", output)
	}
}

func TestQuietMode(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.html
	indexPath := filepath.Join(tmpDir, "index.html")
	if err := os.WriteFile(indexPath, []byte("<html>Home</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	s := New(Config{Dir: tmpDir, Port: 8080, Quiet: true}, &buf)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	s.handleRequest(rec, req)

	output := buf.String()
	if output != "" {
		t.Errorf("Quiet mode should produce no output, got %q", output)
	}
}

func TestNew(t *testing.T) {
	config := Config{
		Dir:  "/some/dir",
		Port: 9999,
	}

	s := New(config, io.Discard)

	if s == nil {
		t.Fatal("New() should not return nil")
	}
	if s.config.Dir != "/some/dir" {
		t.Errorf("Dir = %q, want %q", s.config.Dir, "/some/dir")
	}
	if s.config.Port != 9999 {
		t.Errorf("Port = %d, want %d", s.config.Port, 9999)
	}
}

func TestServerLog(t *testing.T) {
	var buf bytes.Buffer
	s := New(Config{Dir: ".", Port: 8080, Quiet: false}, &buf)

	s.log("test message %s", "arg")

	output := buf.String()
	if !strings.Contains(output, "test message arg") {
		t.Errorf("log output missing expected content, got %q", output)
	}
}

func TestServerLogQuiet(t *testing.T) {
	var buf bytes.Buffer
	s := New(Config{Dir: ".", Port: 8080, Quiet: true}, &buf)

	s.log("test message")

	if buf.String() != "" {
		t.Errorf("quiet mode should suppress log, got %q", buf.String())
	}
}

func TestLogRequest404(t *testing.T) {
	tmpDir := t.TempDir()

	// Create only index.html
	if err := os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte("<html>Home</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	s := New(Config{Dir: tmpDir, Port: 8080, Quiet: false}, &buf)

	req := httptest.NewRequest("GET", "/missing", nil)
	rec := httptest.NewRecorder()

	s.handleRequest(rec, req)

	output := buf.String()
	if !strings.Contains(output, "404") {
		t.Errorf("log should contain 404 status, got %q", output)
	}
}

func TestLogRequest300(t *testing.T) {
	var buf bytes.Buffer
	s := New(Config{Dir: ".", Port: 8080, Quiet: false}, &buf)

	// Test 3xx status code path in logRequest
	s.logRequest("GET", "/redirect", 302, 0)

	output := buf.String()
	if !strings.Contains(output, "302") {
		t.Errorf("log should contain 302 status, got %q", output)
	}
}

func TestResponseRecorder(t *testing.T) {
	rec := httptest.NewRecorder()
	rr := &responseRecorder{ResponseWriter: rec, statusCode: http.StatusOK}

	// Test default status
	if rr.statusCode != http.StatusOK {
		t.Errorf("default status should be 200, got %d", rr.statusCode)
	}

	// Test WriteHeader
	rr.WriteHeader(http.StatusNotFound)
	if rr.statusCode != http.StatusNotFound {
		t.Errorf("status should be 404 after WriteHeader, got %d", rr.statusCode)
	}
}

func TestServerHandler(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "index.html"), []byte("<html>Home</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	s := New(Config{Dir: tmpDir, Port: 8080, Quiet: true}, io.Discard)
	handler := s.Handler()

	if handler == nil {
		t.Fatal("Handler() should not return nil")
	}

	// Test the handler with httptest.Server
	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

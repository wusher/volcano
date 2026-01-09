package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"volcano/internal/tree"
)

func TestNewDynamicServer(t *testing.T) {
	config := DynamicConfig{
		SourceDir: ".",
		Title:     "Test",
		Port:      8080,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer failed: %v", err)
	}

	if srv.config.Title != "Test" {
		t.Errorf("expected title 'Test', got %q", srv.config.Title)
	}
}

func TestDynamicServer_HandleRequest_Index(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	indexContent := "# Welcome\n\nThis is the home page."
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte(indexContent), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", http.NoBody)
	rec := httptest.NewRecorder()

	srv.handleRequest(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Welcome") {
		t.Error("response should contain 'Welcome'")
	}
	if !strings.Contains(body, "Test Site") {
		t.Error("response should contain site title")
	}
}

func TestDynamicServer_HandleRequest_CleanURL(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md and about.md
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "about.md"), []byte("# About Us\n\nAbout page content."), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	// Test /about/ (clean URL)
	req := httptest.NewRequest("GET", "/about/", http.NoBody)
	rec := httptest.NewRecorder()

	srv.handleRequest(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "About Us") {
		t.Error("response should contain 'About Us'")
	}
}

func TestDynamicServer_HandleRequest_NestedPage(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure
	guidesDir := filepath.Join(tmpDir, "guides")
	if err := os.MkdirAll(guidesDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(guidesDir, "intro.md"), []byte("# Introduction\n\nGuide content."), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/guides/intro/", http.NoBody)
	rec := httptest.NewRecorder()

	srv.handleRequest(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Introduction") {
		t.Error("response should contain 'Introduction'")
	}
}

func TestDynamicServer_HandleRequest_DirectoryIndex(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory with index.md
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "index.md"), []byte("# Documentation\n\nDocs index."), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/docs/", http.NoBody)
	rec := httptest.NewRecorder()

	srv.handleRequest(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Documentation") {
		t.Error("response should contain 'Documentation'")
	}
}

func TestDynamicServer_HandleRequest_404(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/nonexistent/", http.NoBody)
	rec := httptest.NewRecorder()

	srv.handleRequest(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Page Not Found") {
		t.Error("response should contain 'Page Not Found'")
	}
}

func TestDynamicServer_HandleRequest_StaticFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md and a static file
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "image.png"), []byte("PNG data"), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/image.png", http.NoBody)
	rec := httptest.NewRecorder()

	srv.handleRequest(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if rec.Body.String() != "PNG data" {
		t.Error("static file content mismatch")
	}
}

func TestDynamicServer_HandleRequest_MarkdownNotServedRaw(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "test.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	// Request the .md file directly - should render as HTML, not serve raw
	req := httptest.NewRequest("GET", "/test.md", http.NoBody)
	rec := httptest.NewRecorder()

	srv.handleRequest(rec, req)

	// Should return 404 for direct .md access (or could render - depends on implementation)
	// The static file handler rejects .md files, so it falls through to page renderer
	// which won't match /test.md as a URL, so 404
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for direct .md access, got %d", rec.Code)
	}
}

func TestDynamicServer_ResolveMarkdownPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create file structure
	guidesDir := filepath.Join(tmpDir, "guides")
	if err := os.MkdirAll(guidesDir, 0755); err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"index.md":        "# Home",
		"about.md":        "# About",
		"guides/index.md": "# Guides",
		"guides/intro.md": "# Intro",
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	config := DynamicConfig{SourceDir: tmpDir}
	srv := &DynamicServer{config: config}

	tests := []struct {
		urlPath  string
		expected string
	}{
		{"/", "index.md"},
		{"/about/", "about.md"},
		{"/about", "about.md"},
		{"/guides/", "guides/index.md"},
		{"/guides/intro/", "guides/intro.md"},
		{"/nonexistent/", ""},
	}

	for _, tc := range tests {
		result := srv.resolveMarkdownPath(tc.urlPath)
		if result != tc.expected {
			t.Errorf("resolveMarkdownPath(%q) = %q, expected %q", tc.urlPath, result, tc.expected)
		}
	}
}

func TestFindNodeBySourcePath(t *testing.T) {
	// Build a simple tree
	root := tree.NewNode("", "", true)

	about := tree.NewNode("About", "about.md", false)
	about.SourcePath = "/tmp/about.md"
	root.AddChild(about)

	guides := tree.NewNode("Guides", "guides", true)
	guides.HasIndex = true
	guides.IndexPath = "guides/index.md"
	guides.SourcePath = "/tmp/guides"
	root.AddChild(guides)

	intro := tree.NewNode("Intro", "guides/intro.md", false)
	intro.SourcePath = "/tmp/guides/intro.md"
	guides.AddChild(intro)

	tests := []struct {
		sourcePath string
		expectNil  bool
		expectName string
	}{
		{"about.md", false, "About"},
		{"guides/intro.md", false, "Intro"},
		{"guides/index.md", false, "Guides"},
		{"nonexistent.md", true, ""},
	}

	for _, tc := range tests {
		result := findNodeBySourcePath(root, tc.sourcePath)
		if tc.expectNil {
			if result != nil {
				t.Errorf("findNodeBySourcePath(%q) expected nil, got %v", tc.sourcePath, result)
			}
		} else {
			if result == nil {
				t.Errorf("findNodeBySourcePath(%q) expected node, got nil", tc.sourcePath)
			} else if result.Name != tc.expectName {
				t.Errorf("findNodeBySourcePath(%q) expected name %q, got %q", tc.sourcePath, tc.expectName, result.Name)
			}
		}
	}
}

func TestFindNodeBySourcePath_NilNode(t *testing.T) {
	result := findNodeBySourcePath(nil, "test.md")
	if result != nil {
		t.Error("expected nil for nil input")
	}
}

func TestDynamicServer_ServeStaticFile_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/subdir", http.NoBody)
	rec := httptest.NewRecorder()

	// serveStaticFile should return false for directories
	result := srv.serveStaticFile(rec, req, "/subdir")
	if result {
		t.Error("serveStaticFile should return false for directories")
	}
}

func TestDynamicServer_ServeStaticFile_RootPath(t *testing.T) {
	tmpDir := t.TempDir()

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", http.NoBody)
	rec := httptest.NewRecorder()

	// serveStaticFile should return false for root path
	result := srv.serveStaticFile(rec, req, "/")
	if result {
		t.Error("serveStaticFile should return false for root path")
	}
}

func TestDynamicServer_Log(t *testing.T) {
	var buf strings.Builder

	config := DynamicConfig{
		SourceDir: ".",
		Title:     "Test",
		Port:      8080,
		Quiet:     false,
	}

	srv, _ := NewDynamicServer(config, &buf)
	srv.log("test message %s", "arg")

	if !strings.Contains(buf.String(), "test message arg") {
		t.Errorf("log output missing: %s", buf.String())
	}
}

func TestDynamicServer_Log_Quiet(t *testing.T) {
	var buf strings.Builder

	config := DynamicConfig{
		SourceDir: ".",
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, _ := NewDynamicServer(config, &buf)
	srv.log("test message")

	if buf.String() != "" {
		t.Errorf("quiet mode should suppress log output, got: %s", buf.String())
	}
}

func TestDynamicServer_LogError(t *testing.T) {
	var buf strings.Builder

	config := DynamicConfig{
		SourceDir: ".",
		Title:     "Test",
		Port:      8080,
		Quiet:     true, // Even in quiet mode, errors should be logged
	}

	srv, _ := NewDynamicServer(config, &buf)
	srv.logError("error: %s", "details")

	if !strings.Contains(buf.String(), "Error: error: details") {
		t.Errorf("error log output missing: %s", buf.String())
	}
}

func TestDynamicServer_LogRequest(t *testing.T) {
	var buf strings.Builder

	config := DynamicConfig{
		SourceDir: ".",
		Title:     "Test",
		Port:      8080,
		Quiet:     false,
	}

	srv, _ := NewDynamicServer(config, &buf)

	// Test 2xx response
	srv.logRequest("GET", "/test", 200, 0)
	if !strings.Contains(buf.String(), "GET") || !strings.Contains(buf.String(), "/test") {
		t.Error("log request should contain method and path")
	}

	buf.Reset()

	// Test 4xx response
	srv.logRequest("GET", "/notfound", 404, 0)
	if !strings.Contains(buf.String(), "404") {
		t.Error("log request should contain status code")
	}
}

func TestDynamicServer_LogRequest_Quiet(t *testing.T) {
	var buf strings.Builder

	config := DynamicConfig{
		SourceDir: ".",
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, _ := NewDynamicServer(config, &buf)
	srv.logRequest("GET", "/test", 200, 0)

	if buf.String() != "" {
		t.Errorf("quiet mode should suppress request log, got: %s", buf.String())
	}
}

func TestDynamicServer_CacheHeaders(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", http.NoBody)
	rec := httptest.NewRecorder()

	srv.handleRequest(rec, req)

	// Check cache control headers
	if rec.Header().Get("Cache-Control") != "no-cache, no-store, must-revalidate" {
		t.Error("missing Cache-Control header")
	}
	if rec.Header().Get("Pragma") != "no-cache" {
		t.Error("missing Pragma header")
	}
	if rec.Header().Get("Expires") != "0" {
		t.Error("missing Expires header")
	}
}

func TestDynamicServer_ServeStaticFile_Nonexistent(t *testing.T) {
	tmpDir := t.TempDir()

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/nonexistent.js", http.NoBody)
	rec := httptest.NewRecorder()

	result := srv.serveStaticFile(rec, req, "/nonexistent.js")
	if result {
		t.Error("serveStaticFile should return false for nonexistent file")
	}
}

func TestDynamicServer_Serve404_WithNav(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md so we have navigation
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/notfound/", http.NoBody)
	rec := httptest.NewRecorder()

	srv.serve404(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	body := rec.Body.String()
	// Should have navigation from the scanned tree
	if !strings.Contains(body, "tree-nav") {
		t.Error("404 page should contain navigation")
	}
}

func TestDynamicServer_LogRequest_3xx(t *testing.T) {
	var buf strings.Builder

	config := DynamicConfig{
		SourceDir: ".",
		Title:     "Test",
		Port:      8080,
		Quiet:     false,
	}

	srv, _ := NewDynamicServer(config, &buf)

	// Test 3xx response (no special coloring)
	srv.logRequest("GET", "/redirect", 302, 0)
	if !strings.Contains(buf.String(), "302") {
		t.Error("log request should contain 302 status")
	}
}

func TestDynamicServer_Handler(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test",
		Port:      8080,
		Quiet:     true,
	}

	srv, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	handler := srv.Handler()
	if handler == nil {
		t.Fatal("Handler() should not return nil")
	}

	// Test with httptest.Server
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

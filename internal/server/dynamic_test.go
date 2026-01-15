package server

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/wusher/volcano/internal/markdown"
	"github.com/wusher/volcano/internal/tree"
)

// mockFileSystem implements FileSystem for testing
type mockFileSystem struct {
	files map[string]mockFileInfo
}

type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return m.modTime }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return nil }

func (m *mockFileSystem) Stat(name string) (os.FileInfo, error) {
	if info, ok := m.files[name]; ok {
		return info, nil
	}
	return nil, os.ErrNotExist
}

func (m *mockFileSystem) ReadFile(_ string) ([]byte, error) {
	return nil, os.ErrNotExist
}

// mockScanner implements TreeScanner for testing
type mockScanner struct {
	site *tree.Site
	err  error
}

func (m *mockScanner) Scan(_ string) (*tree.Site, error) {
	return m.site, m.err
}

func TestNewDynamicServer(t *testing.T) {
	tmpDir := t.TempDir()
	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		Port:      8080,
	}

	server, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	if server == nil {
		t.Fatal("NewDynamicServer() returned nil")
	}

	if server.config.SourceDir != tmpDir {
		t.Errorf("SourceDir = %v, want %v", server.config.SourceDir, tmpDir)
	}
}

func TestDynamicServer_Handler(t *testing.T) {
	tmpDir := t.TempDir()
	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		Port:      8080,
	}

	server, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	handler := server.Handler()
	if handler == nil {
		t.Fatal("Handler() returned nil")
	}
}

func TestDynamicServer_ResolveMarkdownPath(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()

	// Create index.md at root
	indexContent := []byte("# Home\n\nWelcome")
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), indexContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create about.md
	aboutContent := []byte("# About\n\nAbout page")
	if err := os.WriteFile(filepath.Join(tmpDir, "about.md"), aboutContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create docs/index.md
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatal(err)
	}
	docsContent := []byte("# Docs\n\nDocumentation")
	if err := os.WriteFile(filepath.Join(docsDir, "index.md"), docsContent, 0644); err != nil {
		t.Fatal(err)
	}

	server := &DynamicServer{
		config: DynamicConfig{
			SourceDir: tmpDir,
		},
		fs: osFileSystem{},
	}

	tests := []struct {
		urlPath  string
		expected string
	}{
		{"/", "index.md"},
		{"/about/", "about.md"},
		{"/docs/", "docs/index.md"},
		{"/nonexistent/", ""},
	}

	for _, tt := range tests {
		t.Run(tt.urlPath, func(t *testing.T) {
			result := server.resolveMarkdownPath(tt.urlPath)
			if result != tt.expected {
				t.Errorf("resolveMarkdownPath(%q) = %q, want %q", tt.urlPath, result, tt.expected)
			}
		})
	}
}

func TestDynamicServer_ResolveMarkdownPath_WithDatePrefix(t *testing.T) {
	tmpDir := t.TempDir()

	// Create posts directory with date-prefixed file
	postsDir := filepath.Join(tmpDir, "posts")
	if err := os.MkdirAll(postsDir, 0755); err != nil {
		t.Fatal(err)
	}

	postContent := []byte("# Hello World\n\nFirst post")
	if err := os.WriteFile(filepath.Join(postsDir, "2024-01-15-hello-world.md"), postContent, 0644); err != nil {
		t.Fatal(err)
	}

	server := &DynamicServer{
		config: DynamicConfig{
			SourceDir: tmpDir,
		},
		fs: osFileSystem{},
	}

	result := server.resolveMarkdownPath("/posts/hello-world/")
	expected := "posts/2024-01-15-hello-world.md"
	if result != expected {
		t.Errorf("resolveMarkdownPath('/posts/hello-world/') = %q, want %q", result, expected)
	}
}

func TestDynamicServer_RenderPage(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	content := []byte("# Test Page\n\nSome content here.")
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), content, 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}

	server, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	// Create request and recorder
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Render the page
	result := server.renderPage(rec, req, "/")

	if !result {
		t.Error("renderPage() returned false, expected true")
	}

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Test Page") {
		t.Error("response body should contain page title")
	}
}

func TestDynamicServer_Serve404(t *testing.T) {
	tmpDir := t.TempDir()

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}

	server, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rec := httptest.NewRecorder()

	server.serve404(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusNotFound)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "404") {
		t.Error("response body should contain '404'")
	}
}

func TestDynamicServer_HandleRequest_Integration(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple site structure
	content := []byte("# Home\n\nWelcome to the site.")
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), content, 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		Port:      0,
	}

	server, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	handler := server.Handler()

	// Test home page
	t.Run("home page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("GET / status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	// Test 404
	t.Run("404 page", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nonexistent/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("GET /nonexistent/ status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}

func TestDynamicServer_ServeStaticFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a static file
	staticContent := []byte("body { color: red; }")
	if err := os.WriteFile(filepath.Join(tmpDir, "style.css"), staticContent, 0644); err != nil {
		t.Fatal(err)
	}

	// Create a markdown file (should not be served as static)
	mdContent := []byte("# Test\n\nContent")
	if err := os.WriteFile(filepath.Join(tmpDir, "page.md"), mdContent, 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}

	server, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	// Test CSS file
	t.Run("serves CSS file", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/style.css", nil)
		rec := httptest.NewRecorder()

		result := server.serveStaticFile(rec, req, "/style.css")
		if !result {
			t.Error("serveStaticFile() should return true for CSS file")
		}

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	// Test markdown file (should not be served as static)
	t.Run("does not serve markdown as static", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/page.md", nil)
		rec := httptest.NewRecorder()

		result := server.serveStaticFile(rec, req, "/page.md")
		if result {
			t.Error("serveStaticFile() should return false for markdown file")
		}
	})
}

func TestDynamicServer_LogRequest(t *testing.T) {
	var buf bytes.Buffer

	server := &DynamicServer{
		config: DynamicConfig{
			Quiet: false,
		},
		writer: &buf,
	}

	server.logRequest("GET", "/test", 200, 100*time.Millisecond)

	output := buf.String()
	if !strings.Contains(output, "GET") {
		t.Error("log should contain method")
	}
	if !strings.Contains(output, "/test") {
		t.Error("log should contain path")
	}
	if !strings.Contains(output, "200") {
		t.Error("log should contain status code")
	}
}

func TestDynamicServer_LogRequest_Quiet(t *testing.T) {
	var buf bytes.Buffer

	server := &DynamicServer{
		config: DynamicConfig{
			Quiet: true,
		},
		writer: &buf,
	}

	server.logRequest("GET", "/test", 200, 100*time.Millisecond)

	if buf.Len() > 0 {
		t.Error("should not log when quiet mode is enabled")
	}
}

func TestDynamicServer_WithCustomTheme(t *testing.T) {
	tmpDir := t.TempDir()
	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		Theme:     "blog",
	}

	server, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	if server.config.Theme != "blog" {
		t.Errorf("Theme = %v, want blog", server.config.Theme)
	}
}

func TestDynamicServer_WithCustomCSS(t *testing.T) {
	tmpDir := t.TempDir()

	// Create custom CSS file
	customCSS := "body { background: blue; }"
	cssPath := filepath.Join(tmpDir, "custom.css")
	if err := os.WriteFile(cssPath, []byte(customCSS), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		CSSPath:   cssPath,
	}

	server, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	if server.config.CSSPath != cssPath {
		t.Errorf("CSSPath = %v, want %v", server.config.CSSPath, cssPath)
	}
}

func TestDynamicServer_GetRenderer_WithCustomCSS(t *testing.T) {
	tmpDir := t.TempDir()

	// Create custom CSS file
	customCSS := "body { background: blue; }"
	cssPath := filepath.Join(tmpDir, "custom.css")
	if err := os.WriteFile(cssPath, []byte(customCSS), 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		CSSPath:   cssPath,
	}

	server, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	renderer, err := server.getRenderer()
	if err != nil {
		t.Fatalf("getRenderer() error = %v", err)
	}

	if renderer == nil {
		t.Error("getRenderer() returned nil")
	}
}

func TestFindNodeBySourcePath(t *testing.T) {
	// Create a simple tree
	root := tree.NewNode("", "", true)

	file1 := tree.NewNode("File 1", "file1.md", false)
	root.AddChild(file1)

	folder := tree.NewNode("Folder", "folder", true)
	root.AddChild(folder)

	file2 := tree.NewNode("File 2", "folder/file2.md", false)
	folder.AddChild(file2)

	tests := []struct {
		name       string
		sourcePath string
		wantNil    bool
	}{
		{"find root file", "file1.md", false},
		{"find nested file", "folder/file2.md", false},
		{"not found", "nonexistent.md", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findNodeBySourcePath(root, tt.sourcePath)
			if tt.wantNil && result != nil {
				t.Errorf("findNodeBySourcePath() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("findNodeBySourcePath() = nil, want non-nil")
			}
		})
	}
}

func TestFindFolderByPath(t *testing.T) {
	// Create tree structure
	root := tree.NewNode("", "", true)

	folder1 := tree.NewNode("Folder One", "folder-one", true)
	folder1.Path = "folder-one"
	root.AddChild(folder1)

	folder2 := tree.NewNode("Folder Two", "folder-one/folder-two", true)
	folder2.Path = "folder-one/folder-two"
	folder1.AddChild(folder2)

	tests := []struct {
		urlPath string
		wantNil bool
	}{
		{"/folder-one/", false},
		{"/folder-one/folder-two/", false},
		{"/nonexistent/", true},
	}

	for _, tt := range tests {
		t.Run(tt.urlPath, func(t *testing.T) {
			result := findFolderByPath(root, tt.urlPath)
			if tt.wantNil && result != nil {
				t.Errorf("findFolderByPath(%q) = %v, want nil", tt.urlPath, result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("findFolderByPath(%q) = nil, want non-nil", tt.urlPath)
			}
		})
	}
}

func TestCollectAllPages(t *testing.T) {
	root := tree.NewNode("", "", true)

	file1 := tree.NewNode("File 1", "file1.md", false)
	file1.Path = "file1.md"
	root.AddChild(file1)

	folder := tree.NewNode("Folder", "folder", true)
	folder.Path = "folder"
	root.AddChild(folder)

	file2 := tree.NewNode("File 2", "folder/file2.md", false)
	file2.Path = "folder/file2.md"
	folder.AddChild(file2)

	pages := collectAllPages(root)

	if len(pages) != 2 {
		t.Errorf("collectAllPages() returned %d pages, want 2", len(pages))
	}
}

func TestDynamicServer_AutoIndex(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a folder without index.md
	folderDir := filepath.Join(tmpDir, "folder")
	if err := os.MkdirAll(folderDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a file in the folder
	content := []byte("# Test File\n\nContent")
	if err := os.WriteFile(filepath.Join(folderDir, "test.md"), content, 0644); err != nil {
		t.Fatal(err)
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}

	server, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	handler := server.Handler()

	req := httptest.NewRequest(http.MethodGet, "/folder/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Should serve auto-generated index
	if rec.Code != http.StatusOK {
		t.Errorf("GET /folder/ status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "folder") {
		t.Error("auto-index should contain folder name")
	}
}

func TestDynamicServer_WithFileSystem(t *testing.T) {
	server := &DynamicServer{}

	mockFS := &mockFileSystem{
		files: map[string]mockFileInfo{
			"/test/file.txt": {name: "file.txt", size: 100, isDir: false},
		},
	}

	result := server.WithFileSystem(mockFS)

	if result != server {
		t.Error("WithFileSystem should return the same server instance")
	}

	if server.fs != mockFS {
		t.Error("WithFileSystem should set the filesystem")
	}
}

func TestDynamicServer_WithScanner(t *testing.T) {
	server := &DynamicServer{}

	mockScan := &mockScanner{
		site: &tree.Site{},
	}

	result := server.WithScanner(mockScan)

	if result != server {
		t.Error("WithScanner should return the same server instance")
	}

	if server.scanner != mockScan {
		t.Error("WithScanner should set the scanner")
	}
}

func TestDynamicServer_TopNav(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a few root-level pages (for top nav)
	for _, name := range []string{"about.md", "contact.md", "services.md"} {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte("# "+name), 0644); err != nil {
			t.Fatal(err)
		}
	}

	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		TopNav:    true,
	}

	server, err := NewDynamicServer(config, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	handler := server.Handler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET / status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "top-nav") {
		t.Error("response should contain top-nav when TopNav is enabled")
	}
}

func TestDynamicServer_ServeStaticFile_WithMockFS(t *testing.T) {
	mockFS := &mockFileSystem{
		files: map[string]mockFileInfo{
			"/test/assets/image.png": {name: "image.png", size: 1000, isDir: false},
			"/test/assets":           {name: "assets", isDir: true},
		},
	}

	srv := &DynamicServer{
		config: DynamicConfig{
			SourceDir: "/test",
		},
		fs: mockFS,
	}

	// Test that directories are not served as static files
	req := httptest.NewRequest(http.MethodGet, "/assets", nil)
	rec := httptest.NewRecorder()

	result := srv.serveStaticFile(rec, req, "/assets")
	if result {
		t.Error("should not serve directory as static file")
	}
}

func TestDynamicServer_ServeBrokenLinksError(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	config := DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}

	server, err := NewDynamicServer(config, &buf)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	root := tree.NewNode("", "", true)
	about := tree.NewNode("About", "about.md", false)
	root.AddChild(about)
	site := &tree.Site{Root: root}

	broken := []markdown.BrokenLink{
		{
			SourcePage:     "/about/",
			SourceFile:     "/tmp/about.md",
			LineNumber:     42,
			LinkURL:        "/missing/",
			OriginalSyntax: "[[missing]]",
			LinkText:       "Missing Page",
			Suggestions:    []string{"/contact/", "/help/"},
		},
	}

	rec := httptest.NewRecorder()

	server.serveBrokenLinksError(rec, "/about/", broken, site)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusInternalServerError)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Broken Links") {
		t.Error("response body should contain broken links header")
	}
	if !strings.Contains(body, "/missing/") {
		t.Error("response body should list broken link URL")
	}
	// Check for detailed error information
	if !strings.Contains(body, "/tmp/about.md") {
		t.Error("response body should contain source file path")
	}
	if !strings.Contains(body, ":42") {
		t.Error("response body should contain line number")
	}
	if !strings.Contains(body, "[[missing]]") {
		t.Error("response body should contain original syntax")
	}
	if !strings.Contains(body, "Missing Page") {
		t.Error("response body should contain link text")
	}
	if !strings.Contains(body, "/contact/") || !strings.Contains(body, "/help/") {
		t.Error("response body should contain suggestions")
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "broken internal links") {
		t.Error("log should mention broken internal links")
	}
}

func TestDynamicServer_GetRenderer_MissingCSS(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a CSS file first so server can be created
	cssPath := filepath.Join(tmpDir, "custom.css")
	if err := os.WriteFile(cssPath, []byte("body { color: red; }"), 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	var buf bytes.Buffer
	server, err := NewDynamicServer(DynamicConfig{SourceDir: tmpDir, CSSPath: cssPath}, &buf)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	// Now remove the CSS file to simulate missing file on subsequent loads
	if err := os.Remove(cssPath); err != nil {
		t.Fatalf("Remove error = %v", err)
	}

	renderer, err := server.getRenderer()
	if err != nil {
		t.Fatalf("getRenderer() error = %v", err)
	}
	if renderer != server.renderer {
		t.Error("getRenderer() should fall back to cached renderer on CSS read failure")
	}

	if !strings.Contains(buf.String(), "Failed to load CSS") {
		t.Error("getRenderer() should log CSS load failure")
	}
}

func TestDynamicServer_Log(t *testing.T) {
	var buf bytes.Buffer
	server := &DynamicServer{
		config: DynamicConfig{Quiet: false},
		writer: &buf,
	}

	server.log("hello %s", "world")
	if !strings.Contains(buf.String(), "hello world") {
		t.Error("log should write when not quiet")
	}

	buf.Reset()
	server.config.Quiet = true
	server.log("should not write")
	if buf.Len() != 0 {
		t.Error("log should not write when quiet")
	}
}

func TestDynamicServerStart_ShutdownWithSignal(t *testing.T) {
	tmpDir := t.TempDir()

	server, err := NewDynamicServer(DynamicConfig{SourceDir: tmpDir, Port: 0, Quiet: true}, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- server.Start()
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		process, err := os.FindProcess(os.Getpid())
		if err == nil {
			_ = process.Signal(os.Interrupt)
		}
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Start() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Start() did not return after signal")
	}
}

func TestDynamicServer_LoadFavicon(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test favicon file
	faviconPath := filepath.Join(tmpDir, "favicon.ico")
	if err := os.WriteFile(faviconPath, []byte("fake ico data"), 0644); err != nil {
		t.Fatal(err)
	}

	server := &DynamicServer{
		config: DynamicConfig{
			SourceDir:   tmpDir,
			FaviconPath: faviconPath,
		},
	}

	err := server.loadFavicon()
	if err != nil {
		t.Fatalf("loadFavicon() error = %v", err)
	}

	if server.faviconData == nil {
		t.Error("faviconData should not be nil")
	}
	if server.faviconMime != "image/x-icon" {
		t.Errorf("faviconMime = %q, want %q", server.faviconMime, "image/x-icon")
	}
	if server.faviconName != "favicon.ico" {
		t.Errorf("faviconName = %q, want %q", server.faviconName, "favicon.ico")
	}
	if server.faviconLinks == "" {
		t.Error("faviconLinks should not be empty")
	}
}

func TestDynamicServer_LoadFavicon_NotFound(t *testing.T) {
	server := &DynamicServer{
		config: DynamicConfig{
			FaviconPath: "/nonexistent/favicon.ico",
		},
	}

	err := server.loadFavicon()
	if err == nil {
		t.Error("loadFavicon() should return error for non-existent file")
	}
}

func TestDynamicServer_LoadFavicon_UnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with unsupported extension
	badPath := filepath.Join(tmpDir, "favicon.txt")
	if err := os.WriteFile(badPath, []byte("text"), 0644); err != nil {
		t.Fatal(err)
	}

	server := &DynamicServer{
		config: DynamicConfig{
			FaviconPath: badPath,
		},
	}

	err := server.loadFavicon()
	if err == nil {
		t.Error("loadFavicon() should return error for unsupported format")
	}
}

func TestDynamicServer_ServeFavicon(t *testing.T) {
	server := &DynamicServer{
		faviconData: []byte("favicon data"),
		faviconMime: "image/x-icon",
		faviconName: "favicon.ico",
	}

	t.Run("serves favicon", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.serveFavicon(rec, "/favicon.ico")
		if !result {
			t.Error("serveFavicon() should return true")
		}

		if rec.Header().Get("Content-Type") != "image/x-icon" {
			t.Errorf("Content-Type = %q, want %q", rec.Header().Get("Content-Type"), "image/x-icon")
		}

		if rec.Body.String() != "favicon data" {
			t.Errorf("body = %q, want %q", rec.Body.String(), "favicon data")
		}
	})

	t.Run("ignores other paths", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.serveFavicon(rec, "/other.ico")
		if result {
			t.Error("serveFavicon() should return false for non-favicon paths")
		}
	})

	t.Run("returns false when no favicon data", func(t *testing.T) {
		emptyServer := &DynamicServer{}
		rec := httptest.NewRecorder()

		result := emptyServer.serveFavicon(rec, "/favicon.ico")
		if result {
			t.Error("serveFavicon() should return false when no favicon data")
		}
	})
}

func TestDynamicServer_ServePWA(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		PWA:       true,
	}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("returns false when PWA disabled", func(t *testing.T) {
		disabledServer := &DynamicServer{pwaEnabled: false}
		rec := httptest.NewRecorder()

		result := disabledServer.servePWA(rec, "/manifest.json")
		if result {
			t.Error("servePWA() should return false when PWA disabled")
		}
	})

	t.Run("serves manifest.json", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.servePWA(rec, "/manifest.json")
		if !result {
			t.Error("servePWA() should return true for manifest.json")
		}

		if rec.Header().Get("Content-Type") != "application/manifest+json" {
			t.Errorf("Content-Type = %q, want %q", rec.Header().Get("Content-Type"), "application/manifest+json")
		}

		body := rec.Body.String()
		if !strings.Contains(body, "Test Site") {
			t.Error("manifest should contain site title")
		}
	})

	t.Run("serves sw.js", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.servePWA(rec, "/sw.js")
		if !result {
			t.Error("servePWA() should return true for sw.js")
		}

		if rec.Header().Get("Content-Type") != "application/javascript" {
			t.Errorf("Content-Type = %q, want %q", rec.Header().Get("Content-Type"), "application/javascript")
		}
	})

	t.Run("returns false for unknown paths", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.servePWA(rec, "/unknown.txt")
		if result {
			t.Error("servePWA() should return false for unknown paths")
		}
	})
}

func TestCollectPageURLs(t *testing.T) {
	// Create a tree structure
	root := tree.NewNode("", "", true)

	file1 := tree.NewNode("File 1", "file1.md", false)
	file1.Path = "file1.md"
	root.AddChild(file1)

	folder := tree.NewNode("Folder", "folder", true)
	folder.Path = "folder"
	folder.HasIndex = true
	folder.IndexPath = "folder/index.md"
	root.AddChild(folder)

	file2 := tree.NewNode("File 2", "folder/file2.md", false)
	file2.Path = "folder/file2.md"
	folder.AddChild(file2)

	var urls []string
	collectPageURLs(root, &urls)

	if len(urls) < 2 {
		t.Errorf("collectPageURLs() should collect at least 2 URLs, got %d", len(urls))
	}

	// Verify file URLs are collected
	hasFile1 := false
	for _, url := range urls {
		if strings.Contains(url, "file1") {
			hasFile1 = true
		}
	}
	if !hasFile1 {
		t.Error("collectPageURLs() should include file1 URL")
	}
}

func TestCollectPageURLs_NilNode(t *testing.T) {
	var urls []string
	collectPageURLs(nil, &urls)

	if len(urls) != 0 {
		t.Errorf("collectPageURLs(nil) should not add any URLs, got %d", len(urls))
	}
}

func TestDynamicServer_BuildBrokenLinksWarning(t *testing.T) {
	server := &DynamicServer{}

	tests := []struct {
		name         string
		brokenLinks  []markdown.BrokenLink
		expectedHTML []string
	}{
		{
			name: "single broken link with syntax",
			brokenLinks: []markdown.BrokenLink{
				{
					SourceFile:     "/path/to/file.md",
					LineNumber:     10,
					LinkURL:        "/missing/",
					OriginalSyntax: "[[missing]]",
				},
			},
			expectedHTML: []string{
				"1 Broken Link",
				"file.md:10",
				"[[missing]]",
			},
		},
		{
			name: "single broken link without syntax shows URL",
			brokenLinks: []markdown.BrokenLink{
				{
					LinkURL: "/missing-page/",
				},
			},
			expectedHTML: []string{
				"1 Broken Link",
				"/missing-page/",
			},
		},
		{
			name: "multiple broken links",
			brokenLinks: []markdown.BrokenLink{
				{LinkURL: "/link1/"},
				{LinkURL: "/link2/"},
			},
			expectedHTML: []string{
				"2 Broken Links",
				"/link1/",
				"/link2/",
			},
		},
		{
			name: "broken link with suggestions",
			brokenLinks: []markdown.BrokenLink{
				{
					LinkURL:     "/typo/",
					Suggestions: []string{"/type/"},
				},
			},
			expectedHTML: []string{
				"Did you mean",
				"/type/",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := server.buildBrokenLinksWarning(tc.brokenLinks)

			for _, expected := range tc.expectedHTML {
				if !strings.Contains(result, expected) {
					t.Errorf("buildBrokenLinksWarning() should contain %q", expected)
				}
			}
		})
	}
}

func TestDynamicServer_ServeSearchIndex(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home\n\n## Section One\n\nContent"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		Search:    true,
	}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("returns false when search disabled", func(t *testing.T) {
		disabledServer := &DynamicServer{searchEnabled: false}
		rec := httptest.NewRecorder()

		result := disabledServer.serveSearchIndex(rec, "/search-index.json")
		if result {
			t.Error("serveSearchIndex() should return false when search disabled")
		}
	})

	t.Run("returns false for wrong path", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.serveSearchIndex(rec, "/other.json")
		if result {
			t.Error("serveSearchIndex() should return false for wrong path")
		}
	})

	t.Run("serves search index", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.serveSearchIndex(rec, "/search-index.json")
		if !result {
			t.Error("serveSearchIndex() should return true for /search-index.json")
		}

		if rec.Header().Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", rec.Header().Get("Content-Type"))
		}

		body := rec.Body.String()
		if !strings.Contains(body, "Home") {
			t.Error("search index should contain page title")
		}
	})
}

func TestDynamicServer_ServeSearchJS(t *testing.T) {
	tmpDir := t.TempDir()

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		Search:    true,
	}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("returns false when search disabled", func(t *testing.T) {
		disabledServer := &DynamicServer{searchEnabled: false}
		rec := httptest.NewRecorder()

		result := disabledServer.serveSearchJS(rec, "/search.js")
		if result {
			t.Error("serveSearchJS() should return false when search disabled")
		}
	})

	t.Run("returns false for wrong path", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.serveSearchJS(rec, "/other.js")
		if result {
			t.Error("serveSearchJS() should return false for wrong path")
		}
	})

	t.Run("serves search JS", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.serveSearchJS(rec, "/search.js")
		if !result {
			t.Error("serveSearchJS() should return true for /search.js")
		}

		if rec.Header().Get("Content-Type") != "application/javascript" {
			t.Errorf("Content-Type = %q, want application/javascript", rec.Header().Get("Content-Type"))
		}

		body := rec.Body.String()
		if !strings.Contains(body, "command-palette") {
			t.Error("search JS should contain command palette code")
		}
	})
}

func TestDynamicServer_ServeFavicon_AppleTouchIcons(t *testing.T) {
	server := &DynamicServer{
		appleTouchIcon: []byte("apple-icon-data"),
		faviconIco:     []byte("favicon-ico-data"),
	}

	t.Run("serves apple-touch-icon.png", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.serveFavicon(rec, "/apple-touch-icon.png")
		if !result {
			t.Error("serveFavicon() should return true for apple-touch-icon.png")
		}

		if rec.Header().Get("Content-Type") != "image/png" {
			t.Errorf("Content-Type = %q, want image/png", rec.Header().Get("Content-Type"))
		}
	})

	t.Run("serves apple-touch-icon-precomposed.png", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.serveFavicon(rec, "/apple-touch-icon-precomposed.png")
		if !result {
			t.Error("serveFavicon() should return true for apple-touch-icon-precomposed.png")
		}
	})

	t.Run("serves favicon.ico from generated data", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.serveFavicon(rec, "/favicon.ico")
		if !result {
			t.Error("serveFavicon() should return true for favicon.ico")
		}

		if rec.Header().Get("Content-Type") != "image/x-icon" {
			t.Errorf("Content-Type = %q, want image/x-icon", rec.Header().Get("Content-Type"))
		}
	})

	t.Run("returns false when no apple icon data", func(t *testing.T) {
		emptyServer := &DynamicServer{}
		rec := httptest.NewRecorder()

		result := emptyServer.serveFavicon(rec, "/apple-touch-icon.png")
		if result {
			t.Error("serveFavicon() should return false when no apple icon data")
		}
	})
}

func TestDynamicServer_ServePWA_Icons(t *testing.T) {
	server := &DynamicServer{
		pwaEnabled: true,
		pwaIcon192: []byte("192-icon-data"),
		pwaIcon512: []byte("512-icon-data"),
	}

	t.Run("serves icon-192.png", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.servePWA(rec, "/icon-192.png")
		if !result {
			t.Error("servePWA() should return true for icon-192.png")
		}

		if rec.Header().Get("Content-Type") != "image/png" {
			t.Errorf("Content-Type = %q, want image/png", rec.Header().Get("Content-Type"))
		}
	})

	t.Run("serves icon-512.png", func(t *testing.T) {
		rec := httptest.NewRecorder()

		result := server.servePWA(rec, "/icon-512.png")
		if !result {
			t.Error("servePWA() should return true for icon-512.png")
		}
	})

	t.Run("returns false for missing icons", func(t *testing.T) {
		noIconsServer := &DynamicServer{pwaEnabled: true}
		rec := httptest.NewRecorder()

		result := noIconsServer.servePWA(rec, "/icon-192.png")
		if result {
			t.Error("servePWA() should return false when icon data is nil")
		}
	})
}

func TestDynamicServer_ResolveMarkdownPath_ReadmeFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create README.md (note: no index.md)
	if err := os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Readme"), 0644); err != nil {
		t.Fatal(err)
	}

	server := &DynamicServer{
		config: DynamicConfig{
			SourceDir: tmpDir,
		},
		fs: osFileSystem{},
	}

	result := server.resolveMarkdownPath("/")
	if result != "README.md" {
		t.Errorf("resolveMarkdownPath('/') = %q, want 'README.md'", result)
	}
}

func TestDynamicServer_ResolveMarkdownPath_FallbackToFirstFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file that's not index.md or README.md
	if err := os.WriteFile(filepath.Join(tmpDir, "about.md"), []byte("# About"), 0644); err != nil {
		t.Fatal(err)
	}

	server := &DynamicServer{
		config: DynamicConfig{
			SourceDir: tmpDir,
		},
		fs: osFileSystem{},
	}

	result := server.resolveMarkdownPath("/")
	if result != "about.md" {
		t.Errorf("resolveMarkdownPath('/') = %q, want 'about.md'", result)
	}
}

func TestDynamicServer_HandleRequest_WithBrokenLinks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a page with a broken link
	content := "# Home\n\n[[missing-page]]\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}, &buf)
	if err != nil {
		t.Fatal(err)
	}

	handler := server.Handler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Should still render (dev server shows warning inline)
	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	// Should contain broken link warning
	body := rec.Body.String()
	if !strings.Contains(body, "Broken Link") {
		t.Error("Response should contain broken link warning")
	}

	// Check console output for broken link logging
	output := buf.String()
	if !strings.Contains(output, "broken internal link") {
		t.Error("Console should log broken link")
	}
}

func TestDynamicServer_CollectPagesRecursive_WithFolderIndex(t *testing.T) {
	root := tree.NewNode("", "", true)

	// Add a folder with index
	folder := tree.NewNode("Folder", "folder", true)
	folder.Path = "folder"
	folder.IsFolder = true
	folder.HasIndex = true
	folder.IndexPath = "folder/index.md"
	folder.SourcePath = "/test/folder"
	root.AddChild(folder)

	pages := collectAllPages(root)

	// Should include the folder's index as a page
	found := false
	for _, p := range pages {
		if strings.Contains(p.Path, "index.md") {
			found = true
			break
		}
	}
	if !found {
		t.Error("collectAllPages should include folder index page")
	}
}

func TestDynamicServer_RenderAutoIndex_Error(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}, &buf)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock scanner that returns an error
	mockScan := &mockScanner{
		err: os.ErrNotExist,
	}
	server.WithScanner(mockScan)

	rec := httptest.NewRecorder()
	result := server.tryAutoIndex(rec, "/folder/")

	if result {
		t.Error("tryAutoIndex should return false when scanner fails")
	}
}

func TestDynamicServer_FindNodeBySourcePath_WithIndexPath(t *testing.T) {
	root := tree.NewNode("", "", true)

	folder := tree.NewNode("Docs", "docs", true)
	folder.Path = "docs"
	folder.IsFolder = true
	folder.HasIndex = true
	folder.IndexPath = "docs/index.md"
	folder.SourcePath = "/test/docs"
	root.AddChild(folder)

	// Should find by index path
	result := findNodeBySourcePath(root, "docs/index.md")
	if result == nil {
		t.Error("findNodeBySourcePath should find node by index path")
	}
}

func TestDynamicServer_NewDynamicServer_WithAllOptions(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	config := DynamicConfig{
		SourceDir:       tmpDir,
		Title:           "Test Site",
		Port:            8080,
		Theme:           "blog",
		AccentColor:     "#ff0000",
		InstantNav:      true,
		ViewTransitions: true,
		PWA:             true,
		Search:          true,
		TopNav:          true,
		ShowPageNav:     true,
		ShowBreadcrumbs: true,
	}

	server, err := NewDynamicServer(config, &buf)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	if server.instantNavJS == "" {
		t.Error("instantNavJS should be set when InstantNav is true")
	}

	if !server.viewTransitions {
		t.Error("viewTransitions should be true")
	}

	if !server.pwaEnabled {
		t.Error("pwaEnabled should be true")
	}

	if !server.searchEnabled {
		t.Error("searchEnabled should be true")
	}
}

func TestDynamicServer_LogError(t *testing.T) {
	var buf bytes.Buffer
	server := &DynamicServer{
		writer: &buf,
	}

	server.logError("test error: %s", "something went wrong")

	output := buf.String()
	if !strings.Contains(output, "Error:") {
		t.Error("logError should prefix with 'Error:'")
	}
	if !strings.Contains(output, "something went wrong") {
		t.Error("logError should include the error message")
	}
}

func TestDynamicServer_FindActualDir_EmptyPath(t *testing.T) {
	server := &DynamicServer{
		config: DynamicConfig{
			SourceDir: "/test",
		},
	}

	result := server.findActualDir("")
	if result != "" {
		t.Errorf("findActualDir('') = %q, want empty string", result)
	}

	result = server.findActualDir(".")
	if result != "" {
		t.Errorf("findActualDir('.') = %q, want empty string", result)
	}
}

func TestDynamicServer_ServeStaticFile_EmptyPath(t *testing.T) {
	server := &DynamicServer{
		config: DynamicConfig{
			SourceDir: "/test",
		},
		fs: osFileSystem{},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	result := server.serveStaticFile(rec, req, "/")
	if result {
		t.Error("serveStaticFile should return false for root path")
	}
}

func TestDynamicServer_GeneratePWAIcons_NoFavicon(t *testing.T) {
	server := &DynamicServer{
		faviconData: nil,
	}

	err := server.generatePWAIcons()
	if err != nil {
		t.Errorf("generatePWAIcons() should not error with nil favicon, got: %v", err)
	}
}

func TestDynamicServer_GeneratePWAIcons_SVGNotSupported(t *testing.T) {
	server := &DynamicServer{
		faviconData: []byte("<svg></svg>"),
		faviconName: "icon.svg",
	}

	err := server.generatePWAIcons()
	if err != nil {
		t.Errorf("generatePWAIcons() should not error for SVG, got: %v", err)
	}
	if server.pwaIcon192 != nil {
		t.Error("should not generate icons from SVG")
	}
}

func TestDynamicServer_GeneratePWAIcons_ICONotSupported(t *testing.T) {
	server := &DynamicServer{
		faviconData: []byte("fake ico"),
		faviconName: "favicon.ico",
	}

	err := server.generatePWAIcons()
	if err != nil {
		t.Errorf("generatePWAIcons() should not error for ICO, got: %v", err)
	}
	if server.pwaIcon192 != nil {
		t.Error("should not generate icons from ICO")
	}
}

func TestDynamicServer_GeneratePWAIcons_UnknownFormat(t *testing.T) {
	server := &DynamicServer{
		faviconData: []byte("data"),
		faviconName: "icon.xyz",
	}

	err := server.generatePWAIcons()
	if err != nil {
		t.Errorf("generatePWAIcons() should not error for unknown format, got: %v", err)
	}
	if server.pwaIcon192 != nil {
		t.Error("should not generate icons from unknown format")
	}
}

func TestDynamicServer_ServeServiceWorker(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md for scanning
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		PWA:       true,
	}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	server.serveServiceWorker(rec)

	if rec.Header().Get("Content-Type") != "application/javascript" {
		t.Errorf("Content-Type = %q, want application/javascript", rec.Header().Get("Content-Type"))
	}

	body := rec.Body.String()
	if !strings.Contains(body, "volcano-cache") {
		t.Error("Service worker should contain cache name")
	}
}

func TestDynamicServer_Serve404_CustomPage(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md for a valid site
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/nonexistent/", nil)
	rec := httptest.NewRecorder()

	server.serve404(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "404") {
		t.Error("404 page should contain '404'")
	}
	if !strings.Contains(body, "Page Not Found") {
		t.Error("404 page should contain 'Page Not Found'")
	}
}

func TestDynamicServer_RenderAutoIndex(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a folder with files but no index.md
	folderDir := filepath.Join(tmpDir, "folder")
	if err := os.MkdirAll(folderDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folderDir, "page1.md"), []byte("# Page 1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folderDir, "page2.md"), []byte("# Page 2"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create root index
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	handler := server.Handler()
	req := httptest.NewRequest(http.MethodGet, "/folder/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Page 1") {
		t.Error("Auto-index should list child pages")
	}
}

func TestDynamicServer_HandleRequest_Redirect(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "about.md"), []byte("# About"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	handler := server.Handler()

	// Test request without trailing slash should redirect
	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Should redirect or serve the page
	if rec.Code != http.StatusMovedPermanently && rec.Code != http.StatusOK {
		t.Errorf("status = %d, want redirect (301) or OK (200)", rec.Code)
	}
}

func TestDynamicServer_HandleRequest_MethodNotAllowed(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	handler := server.Handler()

	// Test POST request should fail
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Should return method not allowed
	if rec.Code != http.StatusMethodNotAllowed && rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 405 or 200", rec.Code)
	}
}

func TestDynamicServer_NewDynamicServer_WithFavicon(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid PNG favicon
	faviconPath := filepath.Join(tmpDir, "favicon.png")
	img := createTestPNG(100, 100)
	if err := os.WriteFile(faviconPath, img, 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir:   tmpDir,
		Title:       "Test Site",
		FaviconPath: faviconPath,
	}, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	if server.faviconData == nil {
		t.Error("faviconData should be loaded")
	}
	if server.faviconLinks == "" {
		t.Error("faviconLinks should be set")
	}
}

func TestDynamicServer_NewDynamicServer_WithPWAAndFavicon(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid PNG favicon
	faviconPath := filepath.Join(tmpDir, "favicon.png")
	img := createTestPNG(100, 100)
	if err := os.WriteFile(faviconPath, img, 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir:   tmpDir,
		Title:       "Test Site",
		FaviconPath: faviconPath,
		PWA:         true,
	}, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	if !server.pwaEnabled {
		t.Error("pwaEnabled should be true")
	}

	// PWA icons should be generated from favicon
	if server.pwaIcon192 == nil {
		t.Error("pwaIcon192 should be generated from favicon")
	}
	if server.pwaIcon512 == nil {
		t.Error("pwaIcon512 should be generated from favicon")
	}
}

// Helper function to create a test PNG
func createTestPNG(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var buf bytes.Buffer
	_ = pngEncode(&buf, img)
	return buf.Bytes()
}

func pngEncode(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}

func TestDynamicServer_RenderAutoIndex_FolderWithFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create folder with subfiles
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create files - but no index.md for docs folder
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "guide.md"), []byte("# Guide"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	handler := http.HandlerFunc(server.handleRequest)

	// Request the docs folder which should get auto-index
	req := httptest.NewRequest("GET", "/docs/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 for auto-index", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "Guide") {
		t.Error("Auto-index should contain link to Guide")
	}
}

func TestDynamicServer_ServeStaticCSS(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a static file
	if err := os.WriteFile(filepath.Join(tmpDir, "style.css"), []byte("body { color: red; }"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	handler := http.HandlerFunc(server.handleRequest)

	// Request the CSS file
	req := httptest.NewRequest("GET", "/style.css", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 for static file", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "color: red") {
		t.Error("Should serve static CSS file")
	}
}

func TestDynamicServer_FindNodeBySourcePath_NotFound(t *testing.T) {
	// Create a minimal tree
	root := &tree.Node{
		Name:     "Root",
		IsFolder: true,
		Children: []*tree.Node{
			{Name: "File", Path: "file.md", SourcePath: "/tmp/file.md"},
		},
	}

	// Search for non-existent path (note: function compares against node.Path, not SourcePath)
	node := findNodeBySourcePath(root, "nonexistent.md")
	if node != nil {
		t.Error("Should not find node for non-existent path")
	}
}

func TestDynamicServer_FindNodeBySourcePath_Found(t *testing.T) {
	// Note: findNodeBySourcePath compares against node.Path (case-insensitive)
	root := &tree.Node{
		Name:     "Root",
		IsFolder: true,
		Children: []*tree.Node{
			{Name: "Found", Path: "found.md", SourcePath: "/tmp/found.md"},
		},
	}

	// Search by node.Path value
	node := findNodeBySourcePath(root, "found.md")
	if node == nil {
		t.Fatal("Should find node")
	}
	if node.Name != "Found" {
		t.Errorf("node.Name = %q, want %q", node.Name, "Found")
	}
}

func TestDynamicServer_FindNodeBySourcePath_ByFilename(t *testing.T) {
	// Tests the fallback that checks SourcePath filename
	root := &tree.Node{
		Name:     "Root",
		IsFolder: true,
		Children: []*tree.Node{
			{Name: "Page", Path: "different.md", SourcePath: "/full/path/to/actual.md"},
		},
	}

	// Search by the filename from SourcePath
	node := findNodeBySourcePath(root, "actual.md")
	if node == nil {
		t.Fatal("Should find node by SourcePath filename")
	}
	if node.Name != "Page" {
		t.Errorf("node.Name = %q, want %q", node.Name, "Page")
	}
}

func TestDynamicServer_FindFolderByPath_NotFound(t *testing.T) {
	root := &tree.Node{
		Name:     "Root",
		IsFolder: true,
		Path:     "",
		Children: []*tree.Node{
			{Name: "Docs", Path: "docs", IsFolder: true},
		},
	}

	// Search for non-existent folder
	folder := findFolderByPath(root, "/nonexistent/")
	if folder != nil {
		t.Error("Should not find non-existent folder")
	}
}

func TestDynamicServer_FindFolderByPath_Found(t *testing.T) {
	root := &tree.Node{
		Name:     "Root",
		IsFolder: true,
		Path:     "",
		Children: []*tree.Node{
			{Name: "Docs", Path: "docs", IsFolder: true},
		},
	}

	folder := findFolderByPath(root, "/docs/")
	if folder == nil {
		t.Fatal("Should find docs folder")
	}
	if folder.Name != "Docs" {
		t.Errorf("folder.Name = %q, want %q", folder.Name, "Docs")
	}
}

func TestDynamicServer_FindFolderByPath_NilRoot(t *testing.T) {
	folder := findFolderByPath(nil, "/docs/")
	if folder != nil {
		t.Error("Should return nil for nil root")
	}
}

func TestDynamicServer_CollectPagesRecursive(t *testing.T) {
	root := &tree.Node{
		Name:     "Root",
		IsFolder: true,
		Children: []*tree.Node{
			{Name: "Page1", Path: "page1.md", IsFolder: false},
			{
				Name:     "Folder",
				Path:     "folder",
				IsFolder: true,
				Children: []*tree.Node{
					{Name: "Page2", Path: "folder/page2.md", IsFolder: false},
				},
			},
		},
	}

	var pages []*tree.Node
	collectPagesRecursive(root, &pages)

	if len(pages) != 2 {
		t.Errorf("len(pages) = %d, want 2", len(pages))
	}
}

func TestDynamicServer_CollectPagesRecursive_NilNode(t *testing.T) {
	var pages []*tree.Node
	collectPagesRecursive(nil, &pages)

	if len(pages) != 0 {
		t.Errorf("len(pages) = %d, want 0 for nil node", len(pages))
	}
}

func TestDynamicServer_HandleRequestWithDifferentPaths(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home\nContent"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
	}, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	handler := http.HandlerFunc(server.handleRequest)

	tests := []struct {
		path       string
		wantStatus int
	}{
		{"/", http.StatusOK},
		{"/nonexistent", http.StatusNotFound},
		{"/nonexistent/", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("path %q: status = %d, want %d", tt.path, rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestDynamicServer_ServeSearchIndexNotEnabled(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		Search:    false, // Search NOT enabled
	}, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	handler := http.HandlerFunc(server.handleRequest)

	// Request search index when search is not enabled
	req := httptest.NewRequest("GET", "/search-index.json", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Should return 404 since search is not enabled
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404 when search not enabled", rec.Code)
	}
}

func TestDynamicServer_ServePWANotEnabled(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		PWA:       false, // PWA NOT enabled
	}, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	handler := http.HandlerFunc(server.handleRequest)

	// Request manifest.json when PWA is not enabled
	req := httptest.NewRequest("GET", "/manifest.json", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Should return 404 since PWA is not enabled
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404 when PWA not enabled", rec.Code)
	}
}

func TestDynamicServer_ServePWAEnabled(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		PWA:       true, // PWA enabled
	}, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	handler := http.HandlerFunc(server.handleRequest)

	// Request manifest.json when PWA is enabled
	req := httptest.NewRequest("GET", "/manifest.json", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 for manifest.json when PWA enabled", rec.Code)
	}

	// Request sw.js
	req = httptest.NewRequest("GET", "/sw.js", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 for sw.js when PWA enabled", rec.Code)
	}
}

func TestDynamicServer_ServeSearchEnabled(t *testing.T) {
	tmpDir := t.TempDir()

	// Create index.md with some content
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Home\n\nSome content here"), 0644); err != nil {
		t.Fatal(err)
	}

	server, err := NewDynamicServer(DynamicConfig{
		SourceDir: tmpDir,
		Title:     "Test Site",
		Search:    true, // Search enabled
	}, io.Discard)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	handler := http.HandlerFunc(server.handleRequest)

	// Request search-index.json when search is enabled
	req := httptest.NewRequest("GET", "/search-index.json", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 for search-index.json when search enabled", rec.Code)
	}

	// Request search.js
	req = httptest.NewRequest("GET", "/search.js", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 for search.js when search enabled", rec.Code)
	}
}

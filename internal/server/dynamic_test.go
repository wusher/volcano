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

func TestNeedsAutoIndex(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *tree.Node
		expected bool
	}{
		{
			name: "folder without index needs auto-index",
			setup: func() *tree.Node {
				folder := tree.NewNode("Folder", "folder", true)
				file := tree.NewNode("File", "folder/file.md", false)
				folder.AddChild(file)
				return folder
			},
			expected: true,
		},
		{
			name: "folder with index.md does not need auto-index",
			setup: func() *tree.Node {
				folder := tree.NewNode("Folder", "folder", true)
				index := tree.NewNode("Index", "folder/index.md", false)
				index.Path = "folder/index.md"
				folder.AddChild(index)
				return folder
			},
			expected: false,
		},
		{
			name: "folder with HasIndex=true does not need auto-index",
			setup: func() *tree.Node {
				folder := tree.NewNode("Folder", "folder", true)
				folder.HasIndex = true
				return folder
			},
			expected: false,
		},
		{
			name: "non-folder does not need auto-index",
			setup: func() *tree.Node {
				return tree.NewNode("File", "file.md", false)
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.setup()
			result := needsAutoIndex(node)
			if result != tt.expected {
				t.Errorf("needsAutoIndex() = %v, want %v", result, tt.expected)
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
		{SourcePage: "/about/", LinkURL: "/missing/"},
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

	logOutput := buf.String()
	if !strings.Contains(logOutput, "broken internal links") {
		t.Error("log should mention broken internal links")
	}
}

func TestDynamicServer_GetRenderer_MissingCSS(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	server, err := NewDynamicServer(DynamicConfig{SourceDir: tmpDir}, &buf)
	if err != nil {
		t.Fatalf("NewDynamicServer() error = %v", err)
	}

	server.config.CSSPath = filepath.Join(tmpDir, "missing.css")
	renderer, err := server.getRenderer()
	if err != nil {
		t.Fatalf("getRenderer() error = %v", err)
	}
	if renderer != server.renderer {
		t.Error("getRenderer() should fall back to cached renderer on CSS read failure")
	}

	if !strings.Contains(buf.String(), "Failed to read CSS file") {
		t.Error("getRenderer() should log CSS read failure")
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

package assets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetMimeType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"favicon.ico", "image/x-icon"},
		{"favicon.png", "image/png"},
		{"favicon.svg", "image/svg+xml"},
		{"favicon.gif", "image/gif"},
		{"favicon.ICO", "image/x-icon"},
		{"favicon.PNG", "image/png"},
		{"favicon.txt", ""},
		{"favicon.jpg", ""},
	}

	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			result := getMimeType(tc.filename)
			if result != tc.expected {
				t.Errorf("getMimeType(%q) = %q, want %q", tc.filename, result, tc.expected)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	content := "test content"
	if err := os.WriteFile(srcPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Copy file
	dstPath := filepath.Join(tmpDir, "dest.txt")
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile error: %v", err)
	}

	// Verify content
	data, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Errorf("copied content = %q, want %q", string(data), content)
	}
}

func TestCopyFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	err := copyFile(filepath.Join(tmpDir, "nonexistent"), filepath.Join(tmpDir, "dest"))
	if err == nil {
		t.Error("expected error for nonexistent source")
	}
}

func TestProcessFavicon(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a test favicon
	faviconPath := filepath.Join(tmpDir, "favicon.ico")
	if err := os.WriteFile(faviconPath, []byte("fake ico"), 0644); err != nil {
		t.Fatal(err)
	}

	config := FaviconConfig{IconPath: faviconPath}
	links, err := ProcessFavicon(config, outputDir)
	if err != nil {
		t.Fatalf("ProcessFavicon error: %v", err)
	}

	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}

	if links[0].Rel != "icon" {
		t.Errorf("Rel = %q, want %q", links[0].Rel, "icon")
	}
	if links[0].Type != "image/x-icon" {
		t.Errorf("Type = %q, want %q", links[0].Type, "image/x-icon")
	}
	if links[0].Href != "/favicon.ico" {
		t.Errorf("Href = %q, want %q", links[0].Href, "/favicon.ico")
	}

	// Check file was copied
	if _, err := os.Stat(filepath.Join(outputDir, "favicon.ico")); os.IsNotExist(err) {
		t.Error("favicon.ico not copied to output")
	}
}

func TestProcessFaviconNotFound(t *testing.T) {
	config := FaviconConfig{IconPath: "/nonexistent/favicon.ico"}
	_, err := ProcessFavicon(config, "/tmp")
	if err == nil {
		t.Error("expected error for nonexistent favicon")
	}
}

func TestProcessFaviconUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with unsupported extension
	badPath := filepath.Join(tmpDir, "favicon.txt")
	if err := os.WriteFile(badPath, []byte("text"), 0644); err != nil {
		t.Fatal(err)
	}

	config := FaviconConfig{IconPath: badPath}
	_, err := ProcessFavicon(config, tmpDir)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestProcessFaviconEmpty(t *testing.T) {
	config := FaviconConfig{}
	links, err := ProcessFavicon(config, "/tmp")
	if err != nil {
		t.Fatalf("ProcessFavicon error: %v", err)
	}
	if len(links) != 0 {
		t.Errorf("expected 0 links for empty config, got %d", len(links))
	}
}

func TestRenderFaviconLinks(t *testing.T) {
	links := []FaviconLink{
		{Rel: "icon", Type: "image/x-icon", Href: "/favicon.ico"},
		{Rel: "apple-touch-icon", Type: "image/png", Sizes: "180x180", Href: "/apple-touch-icon.png"},
	}

	html := string(RenderFaviconLinks(links))

	if !strings.Contains(html, `rel="icon"`) {
		t.Error("missing icon rel")
	}
	if !strings.Contains(html, `type="image/x-icon"`) {
		t.Error("missing icon type")
	}
	if !strings.Contains(html, `href="/favicon.ico"`) {
		t.Error("missing icon href")
	}
	if !strings.Contains(html, `rel="apple-touch-icon"`) {
		t.Error("missing apple-touch-icon rel")
	}
	if !strings.Contains(html, `sizes="180x180"`) {
		t.Error("missing sizes")
	}
}

func TestRenderFaviconLinksEmpty(t *testing.T) {
	html := RenderFaviconLinks(nil)
	if html != "" {
		t.Errorf("expected empty string, got %q", html)
	}
}

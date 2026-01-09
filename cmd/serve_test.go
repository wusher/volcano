package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestServeStaticDirectory(t *testing.T) {
	// Create a temp directory with static HTML (not a source directory)
	tmpDir := t.TempDir()
	indexPath := filepath.Join(tmpDir, "index.html")
	if err := os.WriteFile(indexPath, []byte("<html><body>Test Server</body></html>"), 0644); err != nil {
		t.Fatal(err)
	}

	port := 18765

	cfg := &Config{
		InputDir: tmpDir,
		Port:     port,
		Quiet:    true,
	}

	// Start server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- Serve(cfg, io.Discard)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Make a test request
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%d/", port), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	if !strings.Contains(string(body), "Test Server") {
		t.Errorf("Body should contain 'Test Server', got %q", string(body))
	}
}

func TestServeSourceDirectory(t *testing.T) {
	// Create a temp directory with markdown files (source directory)
	tmpDir := t.TempDir()
	indexPath := filepath.Join(tmpDir, "index.md")
	if err := os.WriteFile(indexPath, []byte("# Welcome\n\nThis is markdown."), 0644); err != nil {
		t.Fatal(err)
	}

	port := 18766

	cfg := &Config{
		InputDir: tmpDir,
		Title:    "Test Site",
		Port:     port,
		Quiet:    true,
	}

	// Start server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- Serve(cfg, io.Discard)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Make a test request
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%d/", port), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	// Should be rendered HTML, not raw markdown
	if !strings.Contains(string(body), "<h1") {
		t.Errorf("Body should contain rendered HTML h1 tag, got %q", string(body))
	}
	if !strings.Contains(string(body), "Welcome") {
		t.Errorf("Body should contain 'Welcome', got %q", string(body))
	}
}

func TestIsSourceDirectory(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(dir string) error
		expected bool
	}{
		{
			name: "has markdown files",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "test.md"), []byte("# Test"), 0644)
			},
			expected: true,
		},
		{
			name: "has index.html",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html></html>"), 0644)
			},
			expected: false,
		},
		{
			name: "has both md and index.html",
			setup: func(dir string) error {
				if err := os.WriteFile(filepath.Join(dir, "test.md"), []byte("# Test"), 0644); err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html></html>"), 0644)
			},
			expected: false, // index.html takes precedence
		},
		{
			name: "empty directory",
			setup: func(_ string) error {
				return nil
			},
			expected: false,
		},
		{
			name: "only non-markdown files",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "test.txt"), []byte("text"), 0644)
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			if err := tc.setup(tmpDir); err != nil {
				t.Fatal(err)
			}

			result := isSourceDirectory(tmpDir)
			if result != tc.expected {
				t.Errorf("isSourceDirectory() = %v, expected %v", result, tc.expected)
			}
		})
	}
}

func TestIsSourceDirectory_NonexistentDir(t *testing.T) {
	result := isSourceDirectory("/nonexistent/directory/path")
	if result {
		t.Error("isSourceDirectory should return false for nonexistent directory")
	}
}

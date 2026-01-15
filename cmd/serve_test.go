package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// syncBuffer is a thread-safe buffer for testing
type syncBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (sb *syncBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.Write(p)
}

func (sb *syncBuffer) String() string {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.String()
}

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

func TestServeCommandHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := ServeCommand([]string{"-h"}, &stdout, &stderr)
	if err != nil {
		t.Errorf("ServeCommand with -h should not return error, got: %v", err)
	}
	output := stdout.String()
	if !strings.Contains(output, "Usage:") {
		t.Error("Help output should contain 'Usage:'")
	}
	if !strings.Contains(output, "volcano serve") {
		t.Error("Help output should contain 'volcano serve'")
	}
}

func TestServeCommandHelpLong(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := ServeCommand([]string{"--help"}, &stdout, &stderr)
	if err != nil {
		t.Errorf("ServeCommand with --help should not return error, got: %v", err)
	}
	output := stdout.String()
	if !strings.Contains(output, "Usage:") {
		t.Error("Help output should contain 'Usage:'")
	}
}

func TestServeCommandMissingInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := ServeCommand([]string{}, &stdout, &stderr)
	if err == nil {
		t.Error("ServeCommand without input should return error")
	}
	if !strings.Contains(err.Error(), "input folder is required") {
		t.Errorf("Error should mention 'input folder is required', got: %v", err)
	}
}

func TestServeCommandNonexistentInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := ServeCommand([]string{"/nonexistent/directory"}, &stdout, &stderr)
	if err == nil {
		t.Error("ServeCommand with nonexistent input should return error")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Error should mention 'does not exist', got: %v", err)
	}
}

func TestServeCommandInvalidTheme(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := ServeCommand([]string{"--theme=nonexistent", tmpDir}, &stdout, &stderr)
	if err == nil {
		t.Error("ServeCommand with invalid theme should return error")
	}
	if !strings.Contains(err.Error(), "invalid theme") {
		t.Errorf("Error should mention 'invalid theme', got: %v", err)
	}
}

func TestServeCommandNonexistentCSS(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := ServeCommand([]string{"--css=/nonexistent/style.css", tmpDir}, &stdout, &stderr)
	if err == nil {
		t.Error("ServeCommand with nonexistent CSS should return error")
	}
	if !strings.Contains(err.Error(), "CSS file not found") {
		t.Errorf("Error should mention 'CSS file not found', got: %v", err)
	}
}

func TestServeCommandNonexistentFavicon(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := ServeCommand([]string{"--favicon=/nonexistent/favicon.ico", tmpDir}, &stdout, &stderr)
	if err == nil {
		t.Error("ServeCommand with nonexistent favicon should return error")
	}
	if !strings.Contains(err.Error(), "favicon file not found") {
		t.Errorf("Error should mention 'favicon file not found', got: %v", err)
	}
}

func TestPrintServeUsage(t *testing.T) {
	var buf bytes.Buffer
	printServeUsage(&buf)
	output := buf.String()

	expectedPhrases := []string{
		"Start development server",
		"Usage:",
		"volcano serve",
		"volcano server",
		"--port",
		"--title",
		"--theme",
		"--css",
		"--favicon",
		"--quiet",
		"--help",
		"Examples:",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Usage output should contain %q", phrase)
		}
	}
}

func TestServeCommandInputNotDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := ServeCommand([]string{filePath}, &stdout, &stderr)
	if err == nil {
		t.Error("ServeCommand with file as input should return error")
	}
	if !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("Error should mention 'not a directory', got: %v", err)
	}
}

func TestServeCommandDeprecatedViewTransitions(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	stderr := &syncBuffer{}
	// Use a port that won't actually start (we're testing the deprecation warning, not the server)
	// We need the command to fail before starting the server, so use an invalid port
	go func() {
		// This will start a server - we just want to verify the warning is logged
		_ = ServeCommand([]string{"--view-transitions", "-p", "18799", tmpDir}, io.Discard, stderr)
	}()

	// Give it time to process flags and log warning
	time.Sleep(200 * time.Millisecond)

	stderrOutput := stderr.String()
	if !strings.Contains(stderrOutput, "deprecated") {
		t.Errorf("Stderr should contain deprecation warning, got: %q", stderrOutput)
	}
}

func TestServeCommandInvalidFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := ServeCommand([]string{"--invalid-flag-xyz"}, &stdout, &stderr)
	if err == nil {
		t.Error("ServeCommand with invalid flag should return error")
	}
}

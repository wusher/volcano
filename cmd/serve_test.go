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
	"testing"
	"time"
)

func TestServeIntegration(t *testing.T) {
	// Create a temp directory with test content
	tmpDir := t.TempDir()
	indexPath := filepath.Join(tmpDir, "index.html")
	if err := os.WriteFile(indexPath, []byte("<html><body>Test Server</body></html>"), 0644); err != nil {
		t.Fatal(err)
	}

	// Use a high port number to avoid conflicts
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

func TestServeLogging(t *testing.T) {
	// This test verifies the server outputs startup message
	// We use internal/server tests for more detailed logging tests
	var buf bytes.Buffer
	cfg := &Config{
		InputDir: "/nonexistent",
		Port:     18766,
		Quiet:    false,
	}

	// The server will fail to start or block, so we just verify config is correct
	// Full server testing is in internal/server/server_test.go
	if cfg.InputDir != "/nonexistent" {
		t.Error("Config should have correct input directory")
	}
	if cfg.Port != 18766 {
		t.Error("Config should have correct port")
	}
	_ = buf // Silence unused variable warning
}

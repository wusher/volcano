package main

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

func TestValidateInputDir(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid directory",
			path:    tmpDir,
			wantErr: false,
		},
		{
			name:    "non-existent directory",
			path:    filepath.Join(tmpDir, "nonexistent"),
			wantErr: true,
			errMsg:  "does not exist",
		},
		{
			name:    "file instead of directory",
			path:    createTempFile(t, tmpDir),
			wantErr: true,
			errMsg:  "not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInputDir(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInputDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateInputDir() error = %v, should contain %q", err, tt.errMsg)
				}
			}
		})
	}
}

func createTempFile(t *testing.T, dir string) string {
	t.Helper()
	f, err := os.CreateTemp(dir, "test")
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
}

// TestCLIHelp tests that --help flag works correctly
func TestCLIHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{"--help"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "Volcano - Static site generator") {
		t.Errorf("Help output should contain description, got: %s", output)
	}
	if !strings.Contains(output, "build") {
		t.Error("Help output should mention build command")
	}
	if !strings.Contains(output, "serve") {
		t.Error("Help output should mention serve command")
	}
}

// TestCLIShortHelp tests that -h flag works correctly
func TestCLIShortHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{"-h"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "Volcano - Static site generator") {
		t.Errorf("Help output should contain description, got: %s", output)
	}
}

// TestCLIVersion tests that --version flag works correctly
func TestCLIVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{"--version"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "volcano version") {
		t.Errorf("Version output should contain 'volcano version', got: %s", output)
	}
}

// TestCLIShortVersion tests that -v flag works correctly
func TestCLIShortVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{"-v"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	output := stdout.String()
	if !strings.Contains(output, "volcano version") {
		t.Errorf("Version output should contain 'volcano version', got: %s", output)
	}
}

// TestCLIMissingInput tests error when no input folder provided
func TestCLIMissingInput(t *testing.T) {
	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{}, &stdout, &stderr)

	if exitCode == 0 {
		t.Fatal("Expected non-zero exit code when no input folder provided")
	}

	stderrStr := stderr.String()
	if !strings.Contains(stderrStr, "input folder is required") {
		t.Errorf("Error should mention missing input folder, got: %s", stderrStr)
	}
}

// TestCLINonExistentInput tests error when input folder doesn't exist
func TestCLINonExistentInput(t *testing.T) {
	var stdout, stderr bytes.Buffer

	exitCode := Run([]string{"./nonexistent-folder-12345"}, &stdout, &stderr)

	if exitCode == 0 {
		t.Fatal("Expected non-zero exit code when input folder doesn't exist")
	}

	stderrStr := stderr.String()
	if !strings.Contains(stderrStr, "does not exist") {
		t.Errorf("Error should mention folder doesn't exist, got: %s", stderrStr)
	}
}

// TestCLIGenerate tests basic generate command
func TestCLIGenerate(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	mustMkdirAll(t, inputDir)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{inputDir}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Generate command failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Generating site") {
		t.Errorf("Output should contain 'Generating site', got: %s", output)
	}
}

// TestCLIServe tests basic serve command with serve subcommand
func TestCLIServe(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	mustMkdirAll(t, inputDir)

	// Create an index.html for the server to serve
	if err := os.WriteFile(filepath.Join(inputDir, "index.html"), []byte("<html>Test</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	port := 19001
	var stdout bytes.Buffer

	// Start server in a goroutine (it blocks)
	go func() {
		_ = Run([]string{"serve", "-p", fmt.Sprintf("%d", port), inputDir}, &stdout, io.Discard)
	}()

	// Wait for server to start
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
		t.Fatalf("Server should be serving, got error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check stdout contains "Serving"
	output := stdout.String()
	if !strings.Contains(output, "Serving") {
		t.Errorf("Output should contain 'Serving', got: %s", output)
	}
}

// TestCLIServeLong tests serve command with server alias
func TestCLIServeLong(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	mustMkdirAll(t, inputDir)

	// Create an index.html for the server to serve
	if err := os.WriteFile(filepath.Join(inputDir, "index.html"), []byte("<html>Test</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	port := 19002
	var stdout bytes.Buffer

	// Start server in a goroutine (it blocks)
	go func() {
		_ = Run([]string{"server", "-p", fmt.Sprintf("%d", port), inputDir}, &stdout, io.Discard)
	}()

	// Wait for server to start
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
		t.Fatalf("Server should be serving, got error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check stdout contains "Serving"
	output := stdout.String()
	if !strings.Contains(output, "Serving") {
		t.Errorf("Output should contain 'Serving', got: %s", output)
	}
}

// TestCLIWithTitle tests --title flag
func TestCLIWithTitle(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	mustMkdirAll(t, inputDir)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"--title=Custom Title", inputDir}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Command with title failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "Custom Title") {
		t.Errorf("Output should contain custom title, got: %s", output)
	}
}

// TestCLIWithOutput tests -o flag
func TestCLIWithOutput(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "custom-output")
	mustMkdirAll(t, inputDir)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, inputDir}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Command with output flag failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, outputDir) {
		t.Errorf("Output should contain custom output dir, got: %s", output)
	}
}

// TestCLIWithOutputLong tests --output flag
func TestCLIWithOutputLong(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "custom-output")
	mustMkdirAll(t, inputDir)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"--output", outputDir, inputDir}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Command with output flag failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, outputDir) {
		t.Errorf("Output should contain custom output dir, got: %s", output)
	}
}

// TestCLIWithPort tests -p flag
func TestCLIWithPort(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	mustMkdirAll(t, inputDir)

	// Create an index.html for the server to serve
	if err := os.WriteFile(filepath.Join(inputDir, "index.html"), []byte("<html>Test</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	port := 19003
	var stdout bytes.Buffer

	// Start server in a goroutine (it blocks)
	go func() {
		_ = Run([]string{"serve", "-p", fmt.Sprintf("%d", port), inputDir}, &stdout, io.Discard)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Verify server is running on the specified port
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%d/", port), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Server should be serving on port %d, got error: %v", port, err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Check stdout contains the port number
	output := stdout.String()
	if !strings.Contains(output, fmt.Sprintf("%d", port)) {
		t.Errorf("Output should contain custom port, got: %s", output)
	}
}

// TestCLIWithPortLong tests --port flag
func TestCLIWithPortLong(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	mustMkdirAll(t, inputDir)

	// Create an index.html for the server to serve
	if err := os.WriteFile(filepath.Join(inputDir, "index.html"), []byte("<html>Test</html>"), 0644); err != nil {
		t.Fatal(err)
	}

	port := 19004
	var stdout bytes.Buffer

	// Start server in a goroutine (it blocks)
	go func() {
		_ = Run([]string{"serve", "--port", fmt.Sprintf("%d", port), inputDir}, &stdout, io.Discard)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Verify server is running on the specified port
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:%d/", port), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Server should be serving on port %d, got error: %v", port, err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Check stdout contains the port number
	output := stdout.String()
	if !strings.Contains(output, fmt.Sprintf("%d", port)) {
		t.Errorf("Output should contain custom port, got: %s", output)
	}
}

// TestCLIWithQuiet tests -q flag
func TestCLIWithQuiet(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	mustMkdirAll(t, inputDir)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-q", inputDir}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Command with quiet flag failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}
}

// TestCLIWithVerbose tests --verbose flag
func TestCLIWithVerbose(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	mustMkdirAll(t, inputDir)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"--verbose", inputDir}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Command with verbose flag failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}
}

// TestCLIFileAsInput tests error when file given instead of directory
func TestCLIFileAsInput(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := createTempFile(t, tmpDir)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{filePath}, &stdout, &stderr)

	if exitCode == 0 {
		t.Fatal("Expected non-zero exit code when file given as input")
	}

	stderrStr := stderr.String()
	if !strings.Contains(stderrStr, "not a directory") {
		t.Errorf("Error should mention not a directory, got: %s", stderrStr)
	}
}

// TestCLIAllFlagsCombined tests using multiple flags together
func TestCLIAllFlagsCombined(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")
	mustMkdirAll(t, inputDir)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--title=Test Site", inputDir}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Command with multiple flags failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, outputDir) {
		t.Errorf("Output should contain output dir, got: %s", output)
	}
	if !strings.Contains(output, "Test Site") {
		t.Errorf("Output should contain title, got: %s", output)
	}
}

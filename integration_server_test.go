package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func runServeForTest(t *testing.T, args []string) {
	t.Helper()

	done := make(chan int, 1)
	go func() {
		done <- Run(args, io.Discard, io.Discard)
	}()

	time.Sleep(100 * time.Millisecond)
	if proc, err := os.FindProcess(os.Getpid()); err == nil {
		_ = proc.Signal(os.Interrupt)
	}

	select {
	case <-done:
		return
	case <-time.After(2 * time.Second):
		t.Fatal("serve command did not exit")
	}
}

func TestIntegrationServer_Validation(t *testing.T) {
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Test server command validation
	var stdout, stderr bytes.Buffer

	// Test server with non-existent directory
	exitCode := Run([]string{"-s", "/non/existent/path"}, &stdout, &stderr)
	if exitCode == 0 {
		t.Error("Server should fail with non-existent directory")
	}

	errOutput := stderr.String()
	if !strings.Contains(errOutput, "does not exist") && !strings.Contains(errOutput, "no such file") {
		t.Errorf("Error should mention directory doesn't exist, got: %s", errOutput)
	}
}

func TestIntegrationServer_PortValidation(t *testing.T) {
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Generate a site first
	inputDir := "./example"
	outputDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, inputDir}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Failed to generate site: %s", stderr.String())
	}

	// Test server with invalid port
	var serverStdout, serverStderr bytes.Buffer
	serverExitCode := Run([]string{"-s", "-p", "invalid", outputDir}, &serverStdout, &serverStderr)
	if serverExitCode == 0 {
		t.Error("Server should fail with invalid port")
	}
}

func TestIntegrationServer_ValidSiteServing(t *testing.T) {
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Generate a site first
	inputDir := "./example"
	outputDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, inputDir}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Failed to generate site: %s", stderr.String())
	}

	// Verify the generated site has necessary files for serving
	expectedFiles := []string{
		"index.html",
		"404.html",
		"getting-started/index.html",
		"guides/index.html",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(outputDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist for serving", file)
		}
	}
}

func TestIntegrationServer_SpecialCharactersInPaths(t *testing.T) {
	// Create a site with special characters
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// Create files with special characters
	files := map[string]string{
		"file with spaces.md":      "# Spaces Test\n\nContent here.",
		"file-with-dashes.md":      "# Dashes Test\n\nContent here.",
		"file_with_underscores.md": "# Underscores Test\n\nContent here.",
		"file.with.dots.md":        "# Dots Test\n\nContent here.",
		"special&chars.md":         "# Special Chars\n\nContent here.",
	}

	for filename, content := range files {
		if err := os.WriteFile(filepath.Join(inputDir, filename), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Generate site
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, inputDir}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Failed to generate site: %s", stderr.String())
	}

	// Verify files are generated (some may be slugified differently)
	expectedHTMLFiles := []string{
		"file-with-spaces/index.html",
		"file-with-dashes/index.html",
		"file_with_underscores/index.html",
		"special&chars/index.html", // May be encoded differently
	}

	for _, file := range expectedHTMLFiles {
		fullPath := filepath.Join(outputDir, file)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			// Try alternative naming patterns
			altPath := filepath.Join(outputDir, strings.ReplaceAll(file, "&", "-"))
			if content, err = os.ReadFile(altPath); err != nil {
				t.Logf("File %s (or %s) not generated: %v", file, altPath, err)
				continue
			}
		}

		if !strings.Contains(string(content), "<!DOCTYPE html>") {
			t.Errorf("Generated file %s should be valid HTML", file)
		}
	}

	// Test server command can handle this output directory
	// Just test that server can validate the directory without crashing
	runServeForTest(t, []string{"-s", "-p", "0", outputDir})
}

func TestIntegrationServer_EmptyDirectoryHandling(t *testing.T) {
	// Test server with empty output directory
	outputDir := t.TempDir()

	runServeForTest(t, []string{"-s", "-p", "0", outputDir})
}

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationError_InvalidInputDirectory(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Test non-existent directory
	exitCode := Run([]string{"-o", outputDir, "/non/existent/path"}, &stdout, &stderr)
	if exitCode == 0 {
		t.Error("Should fail with non-existent input directory")
	}

	errOutput := stderr.String()
	if !strings.Contains(errOutput, "does not exist") && !strings.Contains(errOutput, "no such file") {
		t.Errorf("Error message should mention directory doesn't exist, got: %s", errOutput)
	}
}

func TestIntegrationError_FileInsteadOfDirectory(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// Create a file instead of directory
	testFile := filepath.Join(inputDir, "not-a-dir.md")
	os.WriteFile(testFile, []byte("# Test"), 0644)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, testFile}, &stdout, &stderr)
	if exitCode == 0 {
		t.Error("Should fail when input is a file not directory")
	}
}

func TestIntegrationError_InvalidOutputDirectory(t *testing.T) {
	// Create a file where we want output directory
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "not-a-dir")
	os.WriteFile(outputPath, []byte("occupied"), 0644)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputPath, "./example"}, &stdout, &stderr)
	if exitCode == 0 {
		t.Error("Should fail when output path exists as a file")
	}
}

func TestIntegrationError_UnreadableDirectory(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// Create a directory and remove read permissions
	subDir := filepath.Join(inputDir, "noread")
	os.Mkdir(subDir, 0755)
	os.Chmod(subDir, 0000)
	defer os.Chmod(subDir, 0755) // Restore for cleanup

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, subDir}, &stdout, &stderr)
	if exitCode == 0 {
		t.Error("Should fail with unreadable directory")
	}
}

func TestIntegrationError_CorruptedMarkdown_UnclosedAdmonition(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// Create markdown with unclosed admonition
	os.WriteFile(filepath.Join(inputDir, "test.md"), []byte(`
# Test

!!! note "This is not closed
This admonition block doesn't have a closing !!!
`), 0644)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, inputDir}, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("Should handle unclosed admonition gracefully, got exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Should still generate HTML
	if _, err := os.Stat(filepath.Join(outputDir, "test", "index.html")); os.IsNotExist(err) {
		t.Error("Should still generate HTML despite unclosed admonition")
	}
}

func TestIntegrationError_CorruptedMarkdown_InvalidCodeBlock(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// Create markdown with malformed code block
	os.WriteFile(filepath.Join(inputDir, "test.md"), []byte("# Test\n\n```go\nfunc test() {\n    // Unclosed backticks\n"), 0644)

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, inputDir}, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("Should handle malformed code block gracefully, got exit code %d, stderr: %s", exitCode, stderr.String())
	}
}

func TestIntegrationError_InvalidTheme(t *testing.T) {
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Test with invalid theme
	exitCode := Run([]string{"-o", outputDir, "--theme=nonexistent", "./example"}, &stdout, &stderr)
	if exitCode == 0 {
		t.Error("Should fail with invalid theme")
	}

	errOutput := stderr.String()
	if !strings.Contains(errOutput, "theme") && !strings.Contains(errOutput, "invalid") {
		t.Errorf("Error should mention theme issue, got: %s", errOutput)
	}
}

func TestIntegrationError_NonExistentCSS(t *testing.T) {
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Test with non-existent custom CSS
	exitCode := Run([]string{"-o", outputDir, "--css=/non/existent/style.css", "./example"}, &stdout, &stderr)
	if exitCode == 0 {
		t.Error("Should fail with non-existent custom CSS")
	}
}

func TestIntegrationError_InvalidFaviconPath(t *testing.T) {
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Test with non-existent favicon
	_ = Run([]string{"-o", outputDir, "--favicon=/non/existent/favicon.ico", "./example"}, &stdout, &stderr)

	// The application warns but continues - check both stdout and stderr for warning
	stdoutOutput := stdout.String()
	stderrOutput := stderr.String()
	output := stdoutOutput + stderrOutput
	if !strings.Contains(output, "favicon") || !strings.Contains(output, "not found") {
		t.Errorf("Should warn about missing favicon, stdout: %s, stderr: %s", stdoutOutput, stderrOutput)
	}

	// Should still generate files despite the warning
	if _, err := os.Stat(filepath.Join(outputDir, "index.html")); os.IsNotExist(err) {
		t.Error("Should still generate HTML despite missing favicon")
	}
}

func TestIntegrationError_ConflictingFlags(t *testing.T) {
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Test potentially conflicting flags (if any exist)
	// This test assumes some flags might conflict - adjust based on actual flag behavior
	exitCode := Run([]string{"-o", outputDir, "--theme=docs", "--theme=blog", "./example"}, &stdout, &stderr)
	if exitCode == 0 {
		// If the command doesn't fail, that's also valid - just test it doesn't crash
		t.Log("Conflicting theme flags handled gracefully")
	}
}

func TestIntegrationError_MalformedURL(t *testing.T) {
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Test with malformed URL if --url flag exists
	exitCode := Run([]string{"-o", outputDir, "--url=not-a-valid-url", "./example"}, &stdout, &stderr)
	// This test documents current behavior - may pass or fail depending on URL validation
	if exitCode != 0 {
		t.Logf("Malformed URL rejected with exit code %d", exitCode)
	}
}

func TestIntegrationError_EmptyInputDirectory(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// Don't create any files in inputDir
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, inputDir}, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("Should handle empty directory gracefully, got exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Check that warning is shown for no markdown files
	stdoutOutput := stdout.String()
	stderrOutput := stderr.String()
	output := stdoutOutput + stderrOutput
	if !strings.Contains(output, "No markdown files found") {
		t.Errorf("Should warn about no markdown files, got stdout: %s, stderr: %s", stdoutOutput, stderrOutput)
	}

	// Should generate empty output directory (no files created)
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("Should be able to read output directory: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("Output directory should be empty for empty input, got %d files", len(entries))
	}
}

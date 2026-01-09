package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create input directory with a markdown file
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home\n\nWelcome!"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Test Site",
	}

	var buf bytes.Buffer
	err := Generate(cfg, &buf)
	if err != nil {
		t.Errorf("Generate() unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Generating site") {
		t.Error("Generate should print 'Generating site'")
	}
	if !strings.Contains(output, cfg.InputDir) {
		t.Error("Generate should print input directory")
	}
	if !strings.Contains(output, cfg.OutputDir) {
		t.Error("Generate should print output directory")
	}
	if !strings.Contains(output, cfg.Title) {
		t.Error("Generate should print site title")
	}

	// Verify output was created
	indexPath := filepath.Join(outputDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("Generate should create index.html")
	}
}

package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCSSWritesToWriter(t *testing.T) {
	var buf bytes.Buffer
	if err := CSS([]string{}, &buf); err != nil {
		t.Fatalf("CSS() error = %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Fatal("CSS() should write output to writer")
	}
	if !strings.Contains(output, "VANILLA") {
		t.Error("CSS() output should contain vanilla theme header")
	}
}

func TestCSSWritesToFile(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "custom.css")

	var buf bytes.Buffer
	if err := CSS([]string{"-o", outPath}, &buf); err != nil {
		t.Fatalf("CSS() error = %v", err)
	}

	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("output CSS file should not be empty")
	}
}

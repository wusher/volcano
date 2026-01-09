package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	cfg := &Config{
		InputDir:  "/tmp/test-input",
		OutputDir: "/tmp/test-output",
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
}

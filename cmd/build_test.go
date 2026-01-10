package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Build([]string{"-h"}, &stdout, &stderr)
	if err != nil {
		t.Errorf("Build with -h should not return error, got: %v", err)
	}
	output := stdout.String()
	if !strings.Contains(output, "Usage:") {
		t.Error("Help output should contain 'Usage:'")
	}
	if !strings.Contains(output, "volcano build") {
		t.Error("Help output should contain 'volcano build'")
	}
}

func TestBuildHelpLong(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Build([]string{"--help"}, &stdout, &stderr)
	if err != nil {
		t.Errorf("Build with --help should not return error, got: %v", err)
	}
	output := stdout.String()
	if !strings.Contains(output, "Usage:") {
		t.Error("Help output should contain 'Usage:'")
	}
}

func TestBuildMissingInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Build([]string{}, &stdout, &stderr)
	if err == nil {
		t.Error("Build without input should return error")
	}
	if !strings.Contains(err.Error(), "input folder is required") {
		t.Errorf("Error should mention 'input folder is required', got: %v", err)
	}
}

func TestBuildNonexistentInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := Build([]string{"/nonexistent/directory"}, &stdout, &stderr)
	if err == nil {
		t.Error("Build with nonexistent input should return error")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Error should mention 'does not exist', got: %v", err)
	}
}

func TestBuildInputNotDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := Build([]string{filePath}, &stdout, &stderr)
	if err == nil {
		t.Error("Build with file as input should return error")
	}
	if !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("Error should mention 'not a directory', got: %v", err)
	}
}

func TestBuildInvalidTheme(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := Build([]string{"--theme=nonexistent", inputDir}, &stdout, &stderr)
	if err == nil {
		t.Error("Build with invalid theme should return error")
	}
	if !strings.Contains(err.Error(), "invalid theme") {
		t.Errorf("Error should mention 'invalid theme', got: %v", err)
	}
}

func TestBuildNonexistentCSS(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := Build([]string{"--css=/nonexistent/style.css", inputDir}, &stdout, &stderr)
	if err == nil {
		t.Error("Build with nonexistent CSS should return error")
	}
	if !strings.Contains(err.Error(), "CSS file not found") {
		t.Errorf("Error should mention 'CSS file not found', got: %v", err)
	}
}

func TestBuildSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home\n\nWelcome!"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := Build([]string{"-o", outputDir, "--title=Test", inputDir}, &stdout, &stderr)
	if err != nil {
		t.Errorf("Build should succeed, got error: %v", err)
	}

	// Check that output was created
	indexPath := filepath.Join(outputDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("Build should create index.html")
	}
}

func TestBuildWithOutputFlag(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "custom-output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := Build([]string{"--output", outputDir, inputDir}, &stdout, &stderr)
	if err != nil {
		t.Errorf("Build should succeed, got error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outputDir, "index.html")); os.IsNotExist(err) {
		t.Error("Build should create output in custom directory")
	}
}

func TestBuildWithQuietFlag(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := Build([]string{"-q", "-o", outputDir, inputDir}, &stdout, &stderr)
	if err != nil {
		t.Errorf("Build should succeed, got error: %v", err)
	}

	// Quiet mode should suppress most output
	output := stdout.String()
	if strings.Contains(output, "Generating site") {
		t.Error("Quiet mode should suppress 'Generating site' message")
	}
}

func TestBuildWithCustomCSS(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")
	cssPath := filepath.Join(tmpDir, "custom.css")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cssPath, []byte("body { color: red; }"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	err := Build([]string{"--css", cssPath, "-o", outputDir, inputDir}, &stdout, &stderr)
	if err != nil {
		t.Errorf("Build with custom CSS should succeed, got error: %v", err)
	}
}

func TestReorderArgs(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		valueFlags map[string]bool
		expected   []string
	}{
		{
			name:       "flags before positional",
			args:       []string{"-q", "./docs"},
			valueFlags: buildValueFlags,
			expected:   []string{"-q", "./docs"},
		},
		{
			name:       "positional before flags",
			args:       []string{"./docs", "-q"},
			valueFlags: buildValueFlags,
			expected:   []string{"-q", "./docs"},
		},
		{
			name:       "value flag with argument",
			args:       []string{"./docs", "-o", "./output"},
			valueFlags: buildValueFlags,
			expected:   []string{"-o", "./output", "./docs"},
		},
		{
			name:       "mixed flags and positional",
			args:       []string{"./docs", "-q", "-o", "./output"},
			valueFlags: buildValueFlags,
			expected:   []string{"-q", "-o", "./output", "./docs"},
		},
		{
			name:       "flag with equals",
			args:       []string{"./docs", "--title=Test"},
			valueFlags: buildValueFlags,
			expected:   []string{"--title=Test", "./docs"},
		},
		{
			name:       "empty args",
			args:       []string{},
			valueFlags: buildValueFlags,
			expected:   []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := reorderArgs(tc.args, tc.valueFlags)
			if len(result) != len(tc.expected) {
				t.Errorf("reorderArgs(%v) length = %d, want %d", tc.args, len(result), len(tc.expected))
				return
			}
			for i, v := range result {
				if v != tc.expected[i] {
					t.Errorf("reorderArgs(%v)[%d] = %q, want %q", tc.args, i, v, tc.expected[i])
				}
			}
		})
	}
}

func TestIsValueFlagInSet(t *testing.T) {
	valueFlags := map[string]bool{
		"o": true, "output": true,
		"title": true,
	}

	tests := []struct {
		flag     string
		expected bool
	}{
		{"-o", true},
		{"--output", true},
		{"-title", true},
		{"--title", true},
		{"-q", false},
		{"--quiet", false},
		{"--unknown", false},
	}

	for _, tc := range tests {
		t.Run(tc.flag, func(t *testing.T) {
			result := isValueFlagInSet(tc.flag, valueFlags)
			if result != tc.expected {
				t.Errorf("isValueFlagInSet(%q) = %v, want %v", tc.flag, result, tc.expected)
			}
		})
	}
}

func TestValidateInputDir(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() string
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid directory",
			setup: func() string {
				dir := t.TempDir()
				return dir
			},
			expectError: false,
		},
		{
			name: "nonexistent path",
			setup: func() string {
				return "/nonexistent/path/to/dir"
			},
			expectError: true,
			errorMsg:    "does not exist",
		},
		{
			name: "file instead of directory",
			setup: func() string {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "file.txt")
				_ = os.WriteFile(filePath, []byte("test"), 0644)
				return filePath
			},
			expectError: true,
			errorMsg:    "not a directory",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := tc.setup()
			err := validateInputDir(path)
			if tc.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				} else if !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("Error should contain %q, got: %v", tc.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPrintBuildUsage(t *testing.T) {
	var buf bytes.Buffer
	printBuildUsage(&buf)
	output := buf.String()

	expectedPhrases := []string{
		"Generate a static site",
		"Usage:",
		"volcano build",
		"--output",
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

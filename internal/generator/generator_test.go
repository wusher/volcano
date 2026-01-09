package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"volcano/internal/tree"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		InputDir:  "/tmp/input",
		OutputDir: "/tmp/output",
		Title:     "Test Site",
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if g == nil {
		t.Fatal("New() returned nil")
	}
}

func TestGenerateEmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Test Site",
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if result.PagesGenerated != 0 {
		t.Errorf("PagesGenerated = %d, want 0", result.PagesGenerated)
	}
	if len(result.Warnings) == 0 {
		t.Error("Should have warning about no markdown files")
	}
}

func TestGenerateSingleFile(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a markdown file
	mdContent := `# Hello World

This is a test page.
`
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Test Site",
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if result.PagesGenerated != 1 {
		t.Errorf("PagesGenerated = %d, want 1", result.PagesGenerated)
	}

	// Check output file exists
	indexPath := filepath.Join(outputDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("index.html should exist")
	}

	// Check 404 page exists
	notFoundPath := filepath.Join(outputDir, "404.html")
	if _, err := os.Stat(notFoundPath); os.IsNotExist(err) {
		t.Error("404.html should exist")
	}

	// Check content
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	html := string(content)
	if !strings.Contains(html, "Hello World") {
		t.Error("Output should contain page title")
	}
	if !strings.Contains(html, "Test Site") {
		t.Error("Output should contain site title")
	}
	if !strings.Contains(html, "This is a test page") {
		t.Error("Output should contain page content")
	}
}

func TestGenerateMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create directory structure
	files := map[string]string{
		"index.md":        "# Home\n\nWelcome!",
		"about.md":        "# About\n\nAbout page.",
		"guides/index.md": "# Guides\n\nGuide index.",
		"guides/intro.md": "# Introduction\n\nIntro content.",
		"guides/setup.md": "# Setup\n\nSetup guide.",
	}

	for path, content := range files {
		fullPath := filepath.Join(inputDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "My Docs",
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if result.PagesGenerated != 5 {
		t.Errorf("PagesGenerated = %d, want 5", result.PagesGenerated)
	}

	// Check output structure
	expectedFiles := []string{
		"index.html",
		"about/index.html",
		"guides/index.html",
		"guides/intro/index.html",
		"guides/setup/index.html",
		"404.html",
	}

	for _, file := range expectedFiles {
		path := filepath.Join(outputDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", file)
		}
	}
}

func TestGenerateWithClean(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create existing output directory with a file
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}
	oldFile := filepath.Join(outputDir, "old.html")
	if err := os.WriteFile(oldFile, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create markdown file
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Test",
		Clean:     true,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Old file should be gone
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("Old file should be removed when Clean is true")
	}

	// New file should exist
	if _, err := os.Stat(filepath.Join(outputDir, "index.html")); os.IsNotExist(err) {
		t.Error("New index.html should exist")
	}
}

func TestGenerateQuietMode(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Test",
		Quiet:     true,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// In quiet mode, buffer should be empty
	if buf.Len() != 0 {
		t.Errorf("In quiet mode, output should be empty, got: %s", buf.String())
	}
}

func TestGenerateVerboseMode(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Test",
		Verbose:   true,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// In verbose mode, should see checkmarks
	output := buf.String()
	if !strings.Contains(output, "✓") {
		t.Error("Verbose mode should show checkmarks")
	}
}

func TestGenerateNavigationLinks(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create files
	files := map[string]string{
		"index.md": "# Home",
		"about.md": "# About",
	}

	for path, content := range files {
		fullPath := filepath.Join(inputDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Test",
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check that pages contain navigation links
	indexContent, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	if !strings.Contains(string(indexContent), "About") {
		t.Error("Index page should contain link to About")
	}

	aboutContent, _ := os.ReadFile(filepath.Join(outputDir, "about", "index.html"))
	// The navigation now uses H1 titles, so index.md with "# Home" shows as "Home"
	if !strings.Contains(string(aboutContent), "Home") {
		t.Error("About page should contain link to Home")
	}
}

func TestGenerateColoredOutput(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Test",
		Colored:   true,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Should have some output
	if buf.Len() == 0 {
		t.Error("With colored=true, should have output")
	}
}

func TestCountFolders(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *tree.Node
		expected int
	}{
		{
			name: "nil node",
			setup: func() *tree.Node {
				return nil
			},
			expected: 0,
		},
		{
			name: "single folder",
			setup: func() *tree.Node {
				return &tree.Node{IsFolder: true}
			},
			expected: 1,
		},
		{
			name: "folder with children",
			setup: func() *tree.Node {
				root := &tree.Node{IsFolder: true}
				root.Children = []*tree.Node{
					{IsFolder: true},
					{IsFolder: false},
					{IsFolder: true, Children: []*tree.Node{{IsFolder: true}}},
				}
				return root
			},
			expected: 4, // root + 2 direct folders + 1 nested
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := countFolders(tc.setup())
			if result != tc.expected {
				t.Errorf("countFolders() = %d, expected %d", result, tc.expected)
			}
		})
	}
}

func TestGenerateWithNestedFolders(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create deeply nested structure
	files := map[string]string{
		"index.md":                    "# Home",
		"level1/index.md":             "# Level 1",
		"level1/level2/index.md":      "# Level 2",
		"level1/level2/level3/doc.md": "# Deep Doc",
	}

	for path, content := range files {
		fullPath := filepath.Join(inputDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Nested Test",
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if result.PagesGenerated != 4 {
		t.Errorf("PagesGenerated = %d, want 4", result.PagesGenerated)
	}

	// Check deeply nested file
	deepPath := filepath.Join(outputDir, "level1", "level2", "level3", "doc", "index.html")
	if _, err := os.Stat(deepPath); os.IsNotExist(err) {
		t.Error("Deeply nested file should exist")
	}
}

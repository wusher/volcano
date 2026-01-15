package generator

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wusher/volcano/internal/tree"
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

func TestPrepareOutputDirWithClean(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "stale.txt"), []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	g, err := New(Config{OutputDir: outputDir, Clean: true}, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := g.prepareOutputDir(); err != nil {
		t.Fatalf("prepareOutputDir() error = %v", err)
	}

	if _, err := os.Stat(outputDir); err != nil {
		t.Fatalf("output dir should exist, error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "stale.txt")); !os.IsNotExist(err) {
		t.Errorf("stale file should be removed, got error: %v", err)
	}
}

func TestCSSLoading(t *testing.T) {
	tmpDir := t.TempDir()
	cssPath := filepath.Join(tmpDir, "custom.css")
	if err := os.WriteFile(cssPath, []byte("body { color: red; }"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test custom CSS file loading via generator
	var buf bytes.Buffer
	config := Config{
		InputDir:  tmpDir,
		OutputDir: filepath.Join(tmpDir, "out"),
		CSSPath:   cssPath,
	}
	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() with custom CSS error = %v", err)
	}
	if g == nil {
		t.Error("Generator should be created with custom CSS")
	}

	// Test theme loading via generator
	config2 := Config{
		InputDir:  tmpDir,
		OutputDir: filepath.Join(tmpDir, "out2"),
		Theme:     "vanilla",
	}
	g2, err := New(config2, &buf)
	if err != nil {
		t.Fatalf("New() with theme error = %v", err)
	}
	if g2 == nil {
		t.Error("Generator should be created with theme")
	}

	// Test missing CSS file
	config3 := Config{
		InputDir:  tmpDir,
		OutputDir: filepath.Join(tmpDir, "out3"),
		CSSPath:   filepath.Join(tmpDir, "missing.css"),
	}
	_, err = New(config3, &buf)
	if err == nil {
		t.Error("New() should return error for missing CSS file")
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
		InputDir:    inputDir,
		OutputDir:   outputDir,
		Title:       "Test",
		ShowPageNav: true, // Enable page navigation for this test
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
		t.Error("Index page should contain link to About (in page navigation)")
	}

	aboutContent, _ := os.ReadFile(filepath.Join(outputDir, "about", "index.html"))
	// The page navigation uses H1 titles, so index.md with "# Home" shows as "Home"
	if !strings.Contains(string(aboutContent), "Home") {
		t.Error("About page should contain link to Home (in page navigation)")
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

func TestGenerateWithAllOptionsEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create test files with links to exercise link validation
	files := map[string]string{
		"index.md":        "# Home\n\nSee [[about]] and [[guides/intro]]",
		"about.md":        "# About\n\nBack to [[index|Home]]",
		"guides/index.md": "# Guides",
		"guides/intro.md": "# Introduction\n\nVisit [[/about/]]",
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
		InputDir:    inputDir,
		OutputDir:   outputDir,
		Title:       "Full Test",
		SiteURL:     "https://example.com",
		Theme:       "docs",
		ShowPageNav: true,
		ShowLastMod: true,
		InstantNav:  true,
		AccentColor: "#ff0000",
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if result.PagesGenerated < 4 {
		t.Errorf("PagesGenerated = %d, want at least 4", result.PagesGenerated)
	}

	// Verify base URL is used in generated files
	indexPath := filepath.Join(outputDir, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}

	if !bytes.Contains(content, []byte("https://example.com")) {
		t.Error("Generated HTML should contain base URL")
	}

	// Verify accent color is applied (now in external CSS file)
	assetsDir := filepath.Join(outputDir, "assets")
	entries, err := os.ReadDir(assetsDir)
	if err != nil {
		t.Fatalf("Failed to read assets directory: %v", err)
	}
	var cssContent []byte
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "styles.") && strings.HasSuffix(entry.Name(), ".css") {
			cssContent, err = os.ReadFile(filepath.Join(assetsDir, entry.Name()))
			if err != nil {
				t.Fatalf("Failed to read CSS file: %v", err)
			}
			break
		}
	}
	if !bytes.Contains(cssContent, []byte("#ff0000")) {
		t.Error("Generated CSS should contain accent color")
	}
}

func TestGenerateWithBrokenLinks(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create test files with broken links
	files := map[string]string{
		"index.md": "# Home\n\nSee [[missing]] page",
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
		Title:     "Broken Links Test",
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err == nil {
		t.Error("Generate() should return error for broken links")
	}
	if err != nil && !strings.Contains(err.Error(), "broken") {
		t.Errorf("Error should mention broken links, got: %v", err)
	}
}

func TestGenerateWithBreadcrumbsDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create test files
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home\n\nTest page"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:        inputDir,
		OutputDir:       outputDir,
		Title:           "Test",
		ShowBreadcrumbs: false, // Disable breadcrumbs
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check that generated HTML doesn't contain breadcrumbs
	indexPath := filepath.Join(outputDir, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(content), "class=\"breadcrumbs\"") {
		t.Error("Generated HTML should not contain breadcrumbs when disabled")
	}
}

func TestGenerateWithBreadcrumbsEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create test files with a subdirectory
	if err := os.MkdirAll(filepath.Join(inputDir, "guides"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "guides", "intro.md"), []byte("# Introduction"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:        inputDir,
		OutputDir:       outputDir,
		Title:           "Test",
		ShowBreadcrumbs: true, // Enable breadcrumbs
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check that generated HTML contains breadcrumbs for the nested page
	guidePath := filepath.Join(outputDir, "guides", "intro", "index.html")
	content, err := os.ReadFile(guidePath)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "class=\"breadcrumbs\"") {
		t.Error("Generated HTML should contain breadcrumbs when enabled")
	}
}

func TestGenerateWithTopNav(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create test files - TopNav requires root files
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "about.md"), []byte("# About"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "docs.md"), []byte("# Docs"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Test",
		TopNav:    true, // Enable top navigation
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check that generated HTML contains top nav
	indexPath := filepath.Join(outputDir, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	// TopNav should include links to About and Docs
	contentStr := string(content)
	if !strings.Contains(contentStr, "About") {
		t.Error("Generated HTML with TopNav should contain About link")
	}
	if !strings.Contains(contentStr, "Docs") {
		t.Error("Generated HTML with TopNav should contain Docs link")
	}
}

func TestGenerateWithFaviconAndOGImage(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create favicon file
	faviconPath := filepath.Join(tmpDir, "favicon.ico")
	if err := os.WriteFile(faviconPath, []byte("fake-icon"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create OG image file
	ogImagePath := filepath.Join(tmpDir, "og-image.png")
	if err := os.WriteFile(ogImagePath, []byte("fake-image"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create markdown file
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:    inputDir,
		OutputDir:   outputDir,
		Title:       "Test",
		FaviconPath: faviconPath,
		OGImage:     ogImagePath,
		Author:      "Test Author",
		SiteURL:     "https://example.com",
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check that favicon was copied
	faviconDest := filepath.Join(outputDir, "favicon.ico")
	if _, err := os.Stat(faviconDest); os.IsNotExist(err) {
		t.Error("Favicon should be copied to output directory")
	}

	// Check that OG image was copied
	ogImageDest := filepath.Join(outputDir, "og-image.png")
	if _, err := os.Stat(ogImageDest); os.IsNotExist(err) {
		t.Error("OG image should be copied to output directory")
	}

	// Check that HTML contains favicon and OG image references
	indexPath := filepath.Join(outputDir, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "favicon") {
		t.Error("Generated HTML should contain favicon link")
	}
	if !strings.Contains(contentStr, "og:image") {
		t.Error("Generated HTML should contain OG image meta tag")
	}
	if !strings.Contains(contentStr, "Test Author") {
		t.Error("Generated HTML should contain author meta tag")
	}
}

func TestGenerateWithPWA(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create markdown files
	files := map[string]string{
		"index.md":        "# Home\n\nWelcome!",
		"about.md":        "# About\n\nAbout page.",
		"guides/index.md": "# Guides\n\nGuide section.",
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
		Title:     "PWA Test Site",
		PWA:       true,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if result.PagesGenerated < 3 {
		t.Errorf("PagesGenerated = %d, want at least 3", result.PagesGenerated)
	}

	// Check manifest.json exists
	manifestPath := filepath.Join(outputDir, "manifest.json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Error("manifest.json should exist when PWA is enabled")
	}

	// Check service worker exists
	swPath := filepath.Join(outputDir, "sw.js")
	if _, err := os.Stat(swPath); os.IsNotExist(err) {
		t.Error("sw.js should exist when PWA is enabled")
	}

	// Check manifest.json content
	manifestContent, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifestContent), "PWA Test Site") {
		t.Error("manifest.json should contain site title")
	}

	// Check service worker content
	swContent, err := os.ReadFile(swPath)
	if err != nil {
		t.Fatal(err)
	}
	swStr := string(swContent)
	if !strings.Contains(swStr, "volcano-cache-") {
		t.Error("Service worker should contain cache name")
	}
	if !strings.Contains(swStr, "URLS_TO_CACHE") {
		t.Error("Service worker should contain URLs to cache")
	}

	// Check HTML contains PWA meta tags
	indexPath := filepath.Join(outputDir, "index.html")
	htmlContent, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	htmlStr := string(htmlContent)
	if !strings.Contains(htmlStr, "manifest.json") {
		t.Error("HTML should link to manifest.json")
	}
	if !strings.Contains(htmlStr, "serviceWorker") {
		t.Error("HTML should contain service worker registration")
	}
}

func TestGenerateWithPWAAndFavicon(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create a valid PNG favicon using Go's image library
	faviconPath := filepath.Join(tmpDir, "favicon.png")
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	f, err := os.Create(faviconPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	_ = f.Close()

	// Create markdown file
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:    inputDir,
		OutputDir:   outputDir,
		Title:       "PWA With Icons",
		PWA:         true,
		FaviconPath: faviconPath,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check that PWA icons were generated
	icon192Path := filepath.Join(outputDir, "icon-192.png")
	icon512Path := filepath.Join(outputDir, "icon-512.png")

	if _, err := os.Stat(icon192Path); os.IsNotExist(err) {
		t.Error("icon-192.png should be generated when PNG favicon provided")
	}
	if _, err := os.Stat(icon512Path); os.IsNotExist(err) {
		t.Error("icon-512.png should be generated when PNG favicon provided")
	}

	// Check manifest.json contains icon references
	manifestContent, err := os.ReadFile(filepath.Join(outputDir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifestContent), "icon-192.png") {
		t.Error("manifest.json should contain icon-192.png reference")
	}
}

func TestGenerateWithPWAAndBaseURL(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create markdown files
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "PWA with Base URL",
		SiteURL:   "https://example.com/docs",
		PWA:       true,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check manifest.json contains base URL
	manifestContent, err := os.ReadFile(filepath.Join(outputDir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	manifestStr := string(manifestContent)
	if !strings.Contains(manifestStr, "/docs/") {
		t.Error("manifest.json should contain base URL path")
	}

	// Check service worker contains base URL in URLs
	swContent, err := os.ReadFile(filepath.Join(outputDir, "sw.js"))
	if err != nil {
		t.Fatal(err)
	}
	swStr := string(swContent)
	if !strings.Contains(swStr, "/docs/") {
		t.Error("Service worker should contain base URL in cached URLs")
	}
}

func TestCollectPageURLs(t *testing.T) {
	tests := []struct {
		name              string
		allPages          []*tree.Node
		autoIndexFolders  []*tree.Node
		baseURL           string
		expectedContains  []string
		expectedMinLength int
	}{
		{
			name: "basic pages without base URL",
			allPages: []*tree.Node{
				{Path: "index.md"},
				{Path: "about.md"},
			},
			autoIndexFolders:  nil,
			baseURL:           "",
			expectedContains:  []string{"/", "/about/"},
			expectedMinLength: 2,
		},
		{
			name: "pages with base URL",
			allPages: []*tree.Node{
				{Path: "index.md"},
			},
			autoIndexFolders:  nil,
			baseURL:           "/docs",
			expectedContains:  []string{"/docs/"},
			expectedMinLength: 1,
		},
		{
			name:     "with auto-index folders",
			allPages: []*tree.Node{},
			autoIndexFolders: []*tree.Node{
				{Path: "guides"},
				{Path: "api"},
			},
			baseURL:           "",
			expectedContains:  []string{"/guides/", "/api/"},
			expectedMinLength: 2,
		},
		{
			name: "mixed pages and folders with base URL",
			allPages: []*tree.Node{
				{Path: "index.md"},
			},
			autoIndexFolders: []*tree.Node{
				{Path: "guides"},
			},
			baseURL:           "/site",
			expectedContains:  []string{"/site/", "/site/guides/"},
			expectedMinLength: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			urls := collectPageURLs(tc.allPages, tc.autoIndexFolders, tc.baseURL)

			if len(urls) < tc.expectedMinLength {
				t.Errorf("Expected at least %d URLs, got %d", tc.expectedMinLength, len(urls))
			}

			urlSet := make(map[string]bool)
			for _, u := range urls {
				urlSet[u] = true
			}

			for _, expected := range tc.expectedContains {
				if !urlSet[expected] {
					t.Errorf("Expected URLs to contain %q, got: %v", expected, urls)
				}
			}
		})
	}
}

func TestGenerateWithInlineAssets(t *testing.T) {
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
		InputDir:     inputDir,
		OutputDir:    outputDir,
		Title:        "Test",
		InlineAssets: true,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check that HTML contains inline CSS
	indexPath := filepath.Join(outputDir, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	// When inline assets is enabled, CSS should be in a <style> tag
	if !strings.Contains(contentStr, "<style>") {
		t.Error("Generated HTML should contain inline CSS when InlineAssets is true")
	}
}

func TestGenerateWithAllowBrokenLinks(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create a page with a broken link
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# Home\n\nSee [[missing-page]] for details."
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:         inputDir,
		OutputDir:        outputDir,
		Title:            "Test",
		AllowBrokenLinks: true, // Allow broken links
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v (expected success with AllowBrokenLinks=true)", err)
	}

	if result.PagesGenerated != 1 {
		t.Errorf("PagesGenerated = %d, want 1", result.PagesGenerated)
	}

	// Should log a warning about broken link
	output := buf.String()
	if !strings.Contains(output, "broken") {
		t.Error("Should log warning about broken links")
	}
}

func TestGenerateWithSearch(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create some pages
	files := map[string]string{
		"index.md": "# Home\n\n## Welcome\n\nThis is the home page.",
		"about.md": "# About\n\n## Our Team\n\nAbout our team.",
		"guide.md": "# Guide\n\n## Getting Started\n\nHow to get started.",
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
		Title:     "Search Test",
		Search:    true,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check search-index.json was created
	searchIndexPath := filepath.Join(outputDir, "search-index.json")
	if _, err := os.Stat(searchIndexPath); os.IsNotExist(err) {
		t.Error("search-index.json should be created when Search=true")
	}

	// Check search.js was created
	searchJSPath := filepath.Join(outputDir, "search.js")
	if _, err := os.Stat(searchJSPath); os.IsNotExist(err) {
		t.Error("search.js should be created when Search=true")
	}

	// Check search index contains page data
	indexContent, err := os.ReadFile(searchIndexPath)
	if err != nil {
		t.Fatal(err)
	}
	indexStr := string(indexContent)
	if !strings.Contains(indexStr, "Home") {
		t.Error("Search index should contain page title 'Home'")
	}
	if !strings.Contains(indexStr, "Welcome") {
		t.Error("Search index should contain heading 'Welcome'")
	}
}

func TestGenerateWithAutoIndex(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create a folder without index.md (should get auto-index)
	folderDir := filepath.Join(inputDir, "folder")
	if err := os.MkdirAll(folderDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create root index
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create files in folder
	if err := os.WriteFile(filepath.Join(folderDir, "page1.md"), []byte("# Page 1\n\nContent"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(folderDir, "page2.md"), []byte("# Page 2\n\nContent"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Auto Index Test",
		Verbose:   true, // Enable verbose to see auto-index messages
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Should generate pages for the folder contents
	if result.PagesGenerated < 3 {
		t.Errorf("PagesGenerated = %d, want at least 3 (index + 2 pages)", result.PagesGenerated)
	}

	// Check folder index was auto-generated
	folderIndexPath := filepath.Join(outputDir, "folder", "index.html")
	if _, err := os.Stat(folderIndexPath); os.IsNotExist(err) {
		t.Error("folder/index.html should be auto-generated")
	}

	// Check auto-index contains links to child pages
	indexContent, err := os.ReadFile(folderIndexPath)
	if err != nil {
		t.Fatal(err)
	}
	indexStr := string(indexContent)
	if !strings.Contains(indexStr, "Page 1") {
		t.Error("Auto-index should contain link to Page 1")
	}
	if !strings.Contains(indexStr, "Page 2") {
		t.Error("Auto-index should contain link to Page 2")
	}
}

func TestGenerateWithLastMod(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create a test page
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# Test Article\n\nSome content here."
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:    inputDir,
		OutputDir:   outputDir,
		Title:       "Test",
		ShowLastMod: true,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check output file was created
	indexPath := filepath.Join(outputDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("index.html should be created")
	}
}

func TestGenerateErrorOnMissingInputDir(t *testing.T) {
	var buf bytes.Buffer
	config := Config{
		InputDir:  "/nonexistent/path/to/input",
		OutputDir: "/tmp/output",
		Title:     "Test",
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err == nil {
		t.Error("Generate() should return error for missing input directory")
	}
}

func TestPrepareOutputDir_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "new", "nested", "output")

	var buf bytes.Buffer
	g, err := New(Config{OutputDir: outputDir}, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := g.prepareOutputDir(); err != nil {
		t.Fatalf("prepareOutputDir() error = %v", err)
	}

	// Check directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("prepareOutputDir() should create nested output directory")
	}
}

func TestGenerateWithViewTransitions(t *testing.T) {
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
		InputDir:        inputDir,
		OutputDir:       outputDir,
		Title:           "Test",
		ViewTransitions: true,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check output contains view transition meta tag
	indexPath := filepath.Join(outputDir, "index.html")
	htmlContent, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	htmlStr := string(htmlContent)

	if !strings.Contains(htmlStr, "view-transition") {
		t.Error("Generated HTML should contain view-transition meta tag")
	}
}

func TestGenerateWithPWAAndSearch(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create a favicon for PWA
	faviconPath := filepath.Join(tmpDir, "favicon.png")
	if err := createTestPNGFile(faviconPath, 100, 100); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home\n\n## Section"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:    inputDir,
		OutputDir:   outputDir,
		Title:       "PWA Test",
		PWA:         true,
		Search:      true,
		FaviconPath: faviconPath,
		SiteURL:     "/docs",
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check PWA files
	manifestPath := filepath.Join(outputDir, "manifest.json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Error("manifest.json should be created with PWA=true")
	}

	swPath := filepath.Join(outputDir, "sw.js")
	if _, err := os.Stat(swPath); os.IsNotExist(err) {
		t.Error("sw.js should be created with PWA=true")
	}

	// Check search files
	searchIndexPath := filepath.Join(outputDir, "search-index.json")
	if _, err := os.Stat(searchIndexPath); os.IsNotExist(err) {
		t.Error("search-index.json should be created with Search=true")
	}

	// Check 404.html
	notFoundPath := filepath.Join(outputDir, "404.html")
	if _, err := os.Stat(notFoundPath); os.IsNotExist(err) {
		t.Error("404.html should be created")
	}
}

func TestGenerateWithBrokenNavLinks(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create a folder structure
	folderDir := filepath.Join(inputDir, "folder")
	if err := os.MkdirAll(folderDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create index at root
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a page in folder
	if err := os.WriteFile(filepath.Join(folderDir, "page.md"), []byte("# Page"), 0644); err != nil {
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

	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Should succeed - auto-index will be generated for folder
	if result.PagesGenerated < 2 {
		t.Errorf("PagesGenerated = %d, want at least 2", result.PagesGenerated)
	}
}

func TestGenerateWithTopNavMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create multiple root-level files
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "about.md"), []byte("# About"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "contact.md"), []byte("# Contact"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Top Nav Test",
		TopNav:    true,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check output contains top nav
	indexPath := filepath.Join(outputDir, "index.html")
	htmlContent, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	htmlStr := string(htmlContent)

	// With TopNav enabled, should see nav items
	if !strings.Contains(htmlStr, "About") {
		t.Error("Top nav should contain link to About")
	}
}

func TestPrepareOutputDir_Clean(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")

	// Create directory with existing files
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}
	existingFile := filepath.Join(outputDir, "existing.html")
	if err := os.WriteFile(existingFile, []byte("<html>"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	g, err := New(Config{OutputDir: outputDir, Clean: true}, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if err := g.prepareOutputDir(); err != nil {
		t.Fatalf("prepareOutputDir() error = %v", err)
	}

	// Existing file should be removed
	if _, err := os.Stat(existingFile); !os.IsNotExist(err) {
		t.Error("prepareOutputDir() with Clean=true should remove existing files")
	}
}

func TestGenerateWithEmptyFolder(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create empty folder
	emptyFolder := filepath.Join(inputDir, "empty")
	if err := os.MkdirAll(emptyFolder, 0755); err != nil {
		t.Fatal(err)
	}

	// Create root index
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
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

	// Should succeed - empty folders are skipped
	_, err = g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
}

func TestGenerateWithNoMarkdownFiles(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	// Create directory with only non-md files
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(inputDir, "image.png"), []byte("png"), 0644); err != nil {
		t.Fatal(err)
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

// Helper function to create a test PNG file
func createTestPNGFile(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return png.Encode(f, img)
}

func TestGenerateWithInvalidFaviconPath(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create index.md
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:    inputDir,
		OutputDir:   outputDir,
		Title:       "Test",
		FaviconPath: "/nonexistent/favicon.png", // Invalid path
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Should succeed with warning about favicon
	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() should succeed despite invalid favicon: %v", err)
	}

	if result.PagesGenerated != 1 {
		t.Errorf("PagesGenerated = %d, want 1", result.PagesGenerated)
	}

	// Check that warning was logged (output contains warning message)
	output := buf.String()
	if !strings.Contains(output, "favicon") {
		t.Log("Expected warning about favicon in output")
	}
}

func TestGenerateWithInvalidOGImagePath(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create index.md
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Title:     "Test",
		OGImage:   "/nonexistent/og-image.png", // Invalid path
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Should succeed with warning about OG image
	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() should succeed despite invalid OG image: %v", err)
	}

	if result.PagesGenerated != 1 {
		t.Errorf("PagesGenerated = %d, want 1", result.PagesGenerated)
	}
}

func TestVerifyLinksWithMissingOutputFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		OutputDir: outputDir,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Create a node that doesn't have an output file
	node := tree.NewNode("Test Page", "test.md", false)
	allPages := []*tree.Node{node}

	// This should find the page as "broken" since no output file exists
	broken := g.verifyLinks(allPages)

	if len(broken) != 1 {
		t.Errorf("verifyLinks() returned %d broken links, want 1", len(broken))
	}
}

func TestVerifyLinksWithFolderNode(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		OutputDir: outputDir,
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Create a folder node (should be skipped since output path is empty)
	folderNode := tree.NewNode("Folder", "folder", true)
	allPages := []*tree.Node{folderNode}

	// Folder nodes have empty output path, so should not be flagged as broken
	broken := g.verifyLinks(allPages)

	if len(broken) != 0 {
		t.Errorf("verifyLinks() returned %d broken links for folder, want 0", len(broken))
	}
}

func TestPrepareOutputDirWithReadOnlyParent(t *testing.T) {
	// Skip on systems where we can't create read-only directories reliably
	if os.Geteuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")

	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Make directory read-only
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chmod(readOnlyDir, 0755) }()

	var buf bytes.Buffer
	config := Config{
		OutputDir: filepath.Join(readOnlyDir, "subdir", "output"),
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Should fail to create directory
	err = g.prepareOutputDir()
	if err == nil {
		t.Error("prepareOutputDir() should fail when parent is read-only")
	}
}

func TestCountFoldersWithNilNode(t *testing.T) {
	count := countFolders(nil)
	if count != 0 {
		t.Errorf("countFolders(nil) = %d, want 0", count)
	}
}

func TestCountFoldersWithNestedStructure(t *testing.T) {
	// Create a tree structure with multiple folders
	root := tree.NewNode("Root", "", true)
	folder1 := tree.NewNode("Folder1", "folder1", true)
	folder2 := tree.NewNode("Folder2", "folder2", true)
	subfolder := tree.NewNode("SubFolder", "folder1/subfolder", true)
	file := tree.NewNode("File", "file.md", false)

	root.AddChild(folder1)
	root.AddChild(folder2)
	root.AddChild(file)
	folder1.AddChild(subfolder)

	count := countFolders(root)
	// root + folder1 + folder2 + subfolder = 4
	if count != 4 {
		t.Errorf("countFolders() = %d, want 4", count)
	}
}

func TestGenerateWithBrokenContentLinks_Fail(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create index.md with broken internal link
	mdContent := `# Home

This page has a [broken link](/nonexistent/).
`
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:         inputDir,
		OutputDir:        outputDir,
		Title:            "Test",
		AllowBrokenLinks: false, // Should fail
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Should fail due to broken link
	_, err = g.Generate()
	if err == nil {
		t.Error("Generate() should fail with broken content link when AllowBrokenLinks=false")
	}
	if !strings.Contains(err.Error(), "broken") {
		t.Errorf("Error should mention broken links: %v", err)
	}
}

func TestGenerateWithBrokenContentLinks_AllowBrokenLinks(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create index.md with broken internal link
	mdContent := `# Home

This page has a [broken link](/nonexistent/).
`
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:         inputDir,
		OutputDir:        outputDir,
		Title:            "Test",
		AllowBrokenLinks: true, // Should succeed with warning
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Should succeed despite broken link
	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() should succeed with AllowBrokenLinks=true: %v", err)
	}

	if result.PagesGenerated != 1 {
		t.Errorf("PagesGenerated = %d, want 1", result.PagesGenerated)
	}

	// Check that warning was logged
	output := buf.String()
	if !strings.Contains(output, "broken") && !strings.Contains(output, "Broken") {
		t.Log("Expected warning about broken links in output")
	}
}

func TestGenerateWithBrokenWikiLink(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create index.md with broken wiki link
	mdContent := `# Home

This page has a [[nonexistent-page|broken wiki link]].
`
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:         inputDir,
		OutputDir:        outputDir,
		Title:            "Test",
		AllowBrokenLinks: true, // Allow it to pass
	}

	g, err := New(config, &buf)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Should succeed with warning
	result, err := g.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if result.PagesGenerated != 1 {
		t.Errorf("PagesGenerated = %d, want 1", result.PagesGenerated)
	}
}

func TestGenerateWithInlineAssets_NoHashedFiles(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create index.md
	if err := os.WriteFile(filepath.Join(inputDir, "index.md"), []byte("# Home"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	config := Config{
		InputDir:     inputDir,
		OutputDir:    outputDir,
		Title:        "Test",
		InlineAssets: true, // Skip writing hashed CSS/JS files
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

	// With InlineAssets, there should be no hashed .js or .css files in output
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".js") && strings.Contains(name, "-") {
			// This would be a hashed file like "main-abc123.js"
			t.Errorf("Found hashed asset file %s when InlineAssets=true", name)
		}
	}
}

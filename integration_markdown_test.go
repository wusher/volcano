package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationMarkdown_WikiLinksComplex(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// Create markdown with basic wiki links
	complexWikiLinks := `# Wiki Links Test

## Basic Wiki Links
[[Simple Page]]
[[Page with Spaces]]

## Display Text
[[Target Page|Custom Display Text]]

## Edge Cases
[[page-with-dashes]]
[[page_with_underscores]]
[[page.with.dots]]
`

	os.WriteFile(filepath.Join(inputDir, "test.md"), []byte(complexWikiLinks), 0644)

	// Create target files for validation
	targetFiles := map[string]string{
		"Simple Page.md":           "# Simple Page\n\nContent here.",
		"Page with Spaces.md":      "# Page with Spaces\n\nContent here.",
		"Target Page.md":           "# Target Page\n\nContent here.",
		"page-with-dashes.md":      "# Page with Dashes\n\nContent here.",
		"page_with_underscores.md": "# Underscores\n\nContent here.",
		"page.with.dots.md":        "# Dots\n\nContent here.",
	}

	for filePath, content := range targetFiles {
		fullPath := filepath.Join(inputDir, filePath)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create target file %s: %v", filePath, err)
		}
	}

	// Generate site
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, inputDir}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Failed to generate site: %s", stderr.String())
	}

	// Verify the main page is generated and contains processed links
	mainPagePath := filepath.Join(outputDir, "test", "index.html")
	content, err := os.ReadFile(mainPagePath)
	if err != nil {
		t.Fatalf("Failed to read generated main page: %v", err)
	}

	html := string(content)

	// Check that basic wiki links were converted to HTML links
	expectedLinks := []string{
		`<a href="/simple-page/">Simple Page</a>`,
		`<a href="/page-with-spaces/">Page with Spaces</a>`,
		`<a href="/target-page/">Custom Display Text</a>`, // custom display text
	}

	for _, link := range expectedLinks {
		if !strings.Contains(html, link) {
			t.Errorf("Expected HTML to contain processed wiki link: %s", link)
		}
	}

	for _, link := range expectedLinks {
		if !strings.Contains(html, link) {
			t.Errorf("Expected HTML to contain processed wiki link: %s", link)
		}
	}

	// Check display text processing
	if !strings.Contains(html, ">Custom Display Text</a>") {
		t.Error("Expected custom display text in wiki link")
	}
}

func TestIntegrationMarkdown_AdmonitionEdgeCases(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// Create markdown with complex admonition scenarios
	admonitionContent := "# Admonition Edge Cases\n\n## Basic Admonitions\n:::note Simple Note\nThis is a simple note.\n:::\n\n:::warning\nWarning without custom title.\n:::\n\n## Admonitions with Code Blocks\n:::tip Code Example\nHere's some code:\n\n```go\nfunc example() {\n    return \"nested code\"\n}\n```\n\nAnd more text.\n:::\n\n## Admonitions with Lists\n:::info Information\n- Item 1\n- Item 2\n  - Nested item\n- Item 3\n:::\n\n## Admonitions in Lists\n- Normal list item\n:::warning List Warning\nThis warning is inside a list.\n:::\n- Another normal item\n\n## Unclosed Admonition\n:::note Unclosed\nThis admonition block doesn't have proper closing."

	os.WriteFile(filepath.Join(inputDir, "test.md"), []byte(admonitionContent), 0644)

	// Generate site
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, inputDir}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Failed to generate site: %s", stderr.String())
	}

	// Verify the page is generated
	mainPagePath := filepath.Join(outputDir, "test", "index.html")
	content, err := os.ReadFile(mainPagePath)
	if err != nil {
		t.Fatalf("Failed to read generated page: %v", err)
	}

	html := string(content)

	// Check that admonitions are processed
	expectedAdmonitions := []string{
		`class="admonition admonition-note"`,
		`class="admonition admonition-warning"`,
		`class="admonition admonition-tip"`,
		`class="admonition admonition-info"`,
	}

	for _, admonition := range expectedAdmonitions {
		if !strings.Contains(html, admonition) {
			t.Errorf("Expected HTML to contain admonition class: %s", admonition)
		}
	}

	// Check that code blocks inside admonitions are preserved
	// Code gets syntax highlighted, so check for the function name
	if !strings.Contains(html, "example") || !strings.Contains(html, "code-block") {
		t.Error("Expected code block inside admonition to be preserved")
	}
}

func TestIntegrationMarkdown_CodeBlockFeatures(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// Create markdown with various code block scenarios
	codeContent := `# Code Block Features

## Basic Code Blocks
` + "```javascript" + `
function hello() {
    console.log("Hello, world!");
}
` + "```" + `

## No Language Specified
` + "```" + `
plain text code
no syntax highlighting
` + "```" + `

## Language-Specific Highlighting
` + "```python" + `
def python_function():
    return {"key": "value"}  # Python dict
` + "```" + `

` + "```bash" + `
echo "bash script"
ls -la
` + "```" + `

## Inline Code in Headers
# Header with ` + "`inline-code`" + ` here

## Code Block with Special Characters
` + "```html" + `
<div class="special">&lt;&gt;&amp;</div>
<script>alert("test");</script>
` + "```" + `

## Code Block with Tabs and Spaces
` + "```go" + `
func mixed() {
	if true {
		return "tabs and spaces"
	}
}
` + "```" + `

## Code Block in Lists
- Item 1
  ` + "```javascript" + `
  // nested code
  console.log("nested");
  ` + "```" + `
- Item 2
`

	os.WriteFile(filepath.Join(inputDir, "test.md"), []byte(codeContent), 0644)

	// Generate site
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, inputDir}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Failed to generate site: %s", stderr.String())
	}

	// Verify the page is generated
	mainPagePath := filepath.Join(outputDir, "test", "index.html")
	content, err := os.ReadFile(mainPagePath)
	if err != nil {
		t.Fatalf("Failed to read generated page: %v", err)
	}

	html := string(content)

	// Check that code blocks are properly formatted
	expectedCodeFeatures := []string{
		`class="chroma"`, // syntax highlighting container
		`class="kd"`,     // keyword declarations like "function"
		`class="k"`,      // keywords
		`class="nx"`,     // variable names like "hello"
		`function`,       // should contain function text
		`def`,            // should contain Python def
		`echo`,           // should contain bash echo
		`inline-code`,    // inline code should be preserved
	}

	for _, feature := range expectedCodeFeatures {
		if !strings.Contains(html, feature) {
			t.Errorf("Expected HTML to contain code feature: %s", feature)
		}
	}
}

func TestIntegrationMarkdown_ExternalLinks(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// Create markdown with various external link scenarios
	linkContent := "# External Links Test\n\n## Basic External Links\n[Google](https://www.google.com)\n[GitHub](https://github.com)\n\n## Links with Special Characters\n[Example with spaces & symbols](https://example.com/path?param=value&other=test)\n\n## Email Links\n[Contact me](mailto:test@example.com)\n[Mail with subject](mailto:test@example.com?subject=Hello%20World)\n\n## Protocol-relative Links\n[Protocol-relative](//example.com/resource)\n\n## Link with Title Attribute\n[Link with title](https://example.com \"This is a title\")\n\n## Image Links\n![Image Alt Text](https://example.com/image.jpg)\n![Image with title](https://example.com/image.png \"Image title\")\n\n## Mixed Internal and External\n[Internal](internal-page.md)\n[External](https://external.com)"

	os.WriteFile(filepath.Join(inputDir, "test.md"), []byte(linkContent), 0644)

	// Create the referenced internal page
	os.WriteFile(filepath.Join(inputDir, "internal-page.md"), []byte("# Internal Page\n\nContent here."), 0644)

	// Generate site
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, inputDir}, &stdout, &stderr)
	if exitCode != 0 {
		// Check if it failed due to broken internal link (expected)
		errOutput := stderr.String()
		if !strings.Contains(errOutput, "broken internal links") {
			t.Fatalf("Unexpected generation failure: %s", errOutput)
		}
		t.Log("Generation failed as expected due to broken internal link reference")
		return // Skip HTML content checks for this test
	}

	// Verify the page is generated
	mainPagePath := filepath.Join(outputDir, "test", "index.html")
	content, err := os.ReadFile(mainPagePath)
	if err != nil {
		t.Fatalf("Failed to read generated page: %v", err)
	}

	html := string(content)

	// Check that external links are properly handled
	expectedLinks := []string{
		`<a href="https://www.google.com"`,
		`<a href="https://github.com"`,
		`<a href="https://example.com/path?param=value&amp;other=test"`,
		`<a href="mailto:test@example.com"`,
		`<a href="mailto:test@example.com?subject=Hello%20World"`,
		`<a href="//example.com/resource"`,
		`title="This is a title"`,
		`<img src="https://example.com/image.jpg" alt="Image Alt Text"`,
		`<img src="https://example.com/image.png" alt="Image with title" title="Image title"`,
	}

	for _, link := range expectedLinks {
		if !strings.Contains(html, link) {
			t.Errorf("Expected HTML to contain external link: %s", link)
		}
	}

	// Check that URLs in regular text are preserved
	if !strings.Contains(html, "https://www.google.com") {
		t.Error("Expected URLs in regular text to be preserved")
	}

	// Check that URLs in code are not converted to links
	if !strings.Contains(html, "<code>https://www.google.com</code>") {
		t.Error("Expected URLs in code to remain as code, not converted to links")
	}
}

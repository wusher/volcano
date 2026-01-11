package markdown

import (
	"strings"
	"testing"
)

func TestNewContentTransformer(t *testing.T) {
	transformer := NewContentTransformer("https://example.com")
	if transformer == nil {
		t.Fatal("expected transformer to be non-nil")
	}
	if transformer.siteURL != "https://example.com" {
		t.Errorf("expected siteURL to be 'https://example.com', got %s", transformer.siteURL)
	}
}

func TestTransform(t *testing.T) {
	transformer := NewContentTransformer("")

	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "heading anchors",
			input:    "<h2>Hello World</h2>",
			contains: []string{"id=\"hello-world\"", "Hello World"},
		},
		{
			name:     "external link",
			input:    `<a href="https://google.com">Google</a>`,
			contains: []string{"external-icon", "google.com", "_blank"},
		},
		{
			name:     "code block wrapping",
			input:    "<pre><code>test code</code></pre>",
			contains: []string{"code-block", "test code", "copy-button"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformer.Transform(tt.input)

			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("Transform() result missing %q\nGot: %s", want, result)
				}
			}
		})
	}
}

func TestTransformMarkdown(t *testing.T) {
	transformer := NewContentTransformer("")

	tests := []struct {
		name          string
		input         string
		sourceDir     string
		sourcePath    string
		outputPath    string
		urlPath       string
		fallbackTitle string
		wantTitle     string
		wantContains  []string
	}{
		{
			name:          "markdown with h1",
			input:         "# My Title\n\nContent here.",
			sourceDir:     "/",
			sourcePath:    "/test.md",
			outputPath:    "test/index.html",
			urlPath:       "/test/",
			fallbackTitle: "Test",
			wantTitle:     "My Title",
			wantContains:  []string{"<h1", "My Title", "Content here"},
		},
		{
			name:          "markdown without h1",
			input:         "Just some content.",
			sourceDir:     "/",
			sourcePath:    "/test.md",
			outputPath:    "test/index.html",
			urlPath:       "/test/",
			fallbackTitle: "Test Page",
			wantTitle:     "Test Page",
			wantContains:  []string{"Just some content"},
		},
		{
			name:          "markdown with frontmatter",
			input:         "---\ntitle: Frontmatter Title\n---\n\n# Heading\n\nContent",
			sourceDir:     "/",
			sourcePath:    "/test.md",
			outputPath:    "test/index.html",
			urlPath:       "/test/",
			fallbackTitle: "Test",
			wantTitle:     "Heading",
			wantContains:  []string{"<h1", "Heading", "Content"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, err := transformer.TransformMarkdown(
				[]byte(tt.input),
				tt.sourceDir,
				tt.sourcePath,
				tt.outputPath,
				tt.urlPath,
				tt.fallbackTitle,
			)

			if err != nil {
				t.Fatalf("TransformMarkdown() error = %v", err)
			}

			if page.Title != tt.wantTitle {
				t.Errorf("TransformMarkdown() title = %q, want %q", page.Title, tt.wantTitle)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(page.Content, want) {
					t.Errorf("TransformMarkdown() content missing %q\nGot: %s", want, page.Content)
				}
			}
		})
	}
}

func TestTransformMarkdownWithWikilinks(t *testing.T) {
	transformer := NewContentTransformer("")

	input := "Link to [[Another Page]] and [[../parent]]"
	page, err := transformer.TransformMarkdown(
		[]byte(input),
		"/docs/",
		"/docs/test.md",
		"docs/test/index.html",
		"/docs/test/",
		"Test",
	)

	if err != nil {
		t.Fatalf("TransformMarkdown() error = %v", err)
	}

	// Should contain converted wikilinks
	if !strings.Contains(page.Content, "href") {
		t.Error("TransformMarkdown() should contain href for wikilinks")
	}
}

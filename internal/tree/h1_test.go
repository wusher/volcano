package tree

import (
	"testing"
)

func TestExtractH1(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "simple h1",
			content:  "# Hello World\n\nContent here",
			expected: "Hello World",
		},
		{
			name:     "h1 with inline markdown",
			content:  "# **Bold** Title\n\nContent",
			expected: "Bold Title",
		},
		{
			name:     "h1 with link",
			content:  "# [Link Text](url) Title\n\nContent",
			expected: "Link Text Title",
		},
		{
			name:     "h1 after frontmatter",
			content:  "---\ntitle: Test\n---\n\n# Actual Title\n\nContent",
			expected: "Actual Title",
		},
		{
			name:     "no h1",
			content:  "## Subheading\n\nContent without H1",
			expected: "",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
		{
			name:     "h1 not at start",
			content:  "Some intro text\n\n# Title\n\nMore content",
			expected: "",
		},
		{
			name:     "h1 with empty lines before",
			content:  "\n\n# Title\n\nContent",
			expected: "Title",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ExtractH1([]byte(tc.content))
			if result != tc.expected {
				t.Errorf("ExtractH1() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestStripInlineMarkdown(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"plain text", "plain text"},
		{"**bold**", "bold"},
		{"*italic*", "italic"},
		{"`code`", "code"},
		{"[link](url)", "link"},
		{"~~strikethrough~~", "strikethrough"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := stripInlineMarkdown(tc.input)
			if result != tc.expected {
				t.Errorf("stripInlineMarkdown(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

package markdown

import (
	"testing"
)

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple heading",
			input:    "# Hello World",
			expected: "Hello World",
		},
		{
			name:     "heading with extra spaces",
			input:    "#   Spaced   Title   ",
			expected: "Spaced   Title",
		},
		{
			name:     "heading after content",
			input:    "Some intro text\n\n# Main Title\n\nMore content",
			expected: "Main Title",
		},
		{
			name:     "multiple headings returns first",
			input:    "# First Title\n\n# Second Title",
			expected: "First Title",
		},
		{
			name:     "no heading",
			input:    "Just some content without a heading",
			expected: "",
		},
		{
			name:     "h2 is not h1",
			input:    "## Not H1\n\nContent",
			expected: "",
		},
		{
			name:     "heading with markdown",
			input:    "# Title with **bold** and *italic*",
			expected: "Title with **bold** and *italic*",
		},
		{
			name:     "empty content",
			input:    "",
			expected: "",
		},
		{
			name:     "heading at end",
			input:    "Content first\n# Title at End",
			expected: "Title at End",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractTitle([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("ExtractTitle() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractTitleFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple h1",
			input:    "<h1>Hello World</h1>",
			expected: "Hello World",
		},
		{
			name:     "h1 with class",
			input:    "<h1 class=\"title\">Styled Title</h1>",
			expected: "Styled Title",
		},
		{
			name:     "h1 with id",
			input:    "<h1 id=\"main-title\">ID Title</h1>",
			expected: "ID Title",
		},
		{
			name:     "multiple h1 returns first",
			input:    "<h1>First</h1><h1>Second</h1>",
			expected: "First",
		},
		{
			name:     "no h1",
			input:    "<h2>Not H1</h2><p>Content</p>",
			expected: "",
		},
		{
			name:     "h1 uppercase tag",
			input:    "<H1>Uppercase Tag</H1>",
			expected: "Uppercase Tag",
		},
		{
			name:     "empty content",
			input:    "",
			expected: "",
		},
		{
			name:     "h1 with whitespace",
			input:    "<h1>  Spaced Title  </h1>",
			expected: "Spaced Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractTitleFromHTML([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("ExtractTitleFromHTML() = %q, want %q", result, tt.expected)
			}
		})
	}
}

package markdown

import (
	"strings"
	"testing"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Getting Started", "getting-started"},
		{"API Reference", "api-reference"},
		{"What's New?", "whats-new"},
		{"Test 123", "test-123"},
		{"  Spaces  ", "spaces"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"Special!@#Characters", "specialcharacters"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := Slugify(tc.input)
			if result != tc.expected {
				t.Errorf("Slugify(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestAddHeadingAnchors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:  "h2 heading",
			input: "<h2>Getting Started</h2>",
			contains: []string{
				`id="getting-started"`,
				`class="heading-anchor"`,
				`href="#getting-started"`,
			},
		},
		{
			name:  "h3 heading",
			input: "<h3>Installation</h3>",
			contains: []string{
				`id="installation"`,
			},
		},
		{
			name:  "replaces existing id with generated one",
			input: `<h2 id="custom-id">Title</h2>`,
			contains: []string{
				`id="title"`,
				`href="#title"`,
			},
		},
		{
			name:     "no headings",
			input:    "<p>Just a paragraph</p>",
			contains: []string{"<p>Just a paragraph</p>"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := AddHeadingAnchors(tc.input)
			for _, expected := range tc.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("result should contain %q\ngot: %s", expected, result)
				}
			}
		})
	}
}

func TestAddHeadingAnchorsMultiple(t *testing.T) {
	input := `
<h2>First Section</h2>
<p>Content</p>
<h3>Subsection</h3>
<p>More content</p>
<h2>Second Section</h2>
`

	result := AddHeadingAnchors(input)

	if !strings.Contains(result, `id="first-section"`) {
		t.Error("should contain first-section id")
	}
	if !strings.Contains(result, `id="subsection"`) {
		t.Error("should contain subsection id")
	}
	if !strings.Contains(result, `id="second-section"`) {
		t.Error("should contain second-section id")
	}
}

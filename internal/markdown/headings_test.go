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

func TestEscapePseudoHTMLTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "escapes pseudo tags",
			input:    "-o, --output <dir>",
			expected: "-o, --output &lt;dir&gt;",
		},
		{
			name:     "escapes multiple pseudo tags",
			input:    "--port <port> --host <host>",
			expected: "--port &lt;port&gt; --host &lt;host&gt;",
		},
		{
			name:     "preserves valid HTML tags",
			input:    "Use <code>command</code> here",
			expected: "Use <code>command</code> here",
		},
		{
			name:     "preserves strong and em",
			input:    "<strong>bold</strong> and <em>italic</em>",
			expected: "<strong>bold</strong> and <em>italic</em>",
		},
		{
			name:     "no tags",
			input:    "Simple text",
			expected: "Simple text",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := escapePseudoHTMLTags(tc.input)
			if result != tc.expected {
				t.Errorf("escapePseudoHTMLTags(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestAddHeadingAnchorsWithPseudoTags(t *testing.T) {
	input := "<h3>-o, --output <dir></h3>"

	result := AddHeadingAnchors(input)

	// Should escape <dir> in the heading content
	if !strings.Contains(result, "&lt;dir&gt;") {
		t.Errorf("should escape <dir> to &lt;dir&gt;\ngot: %s", result)
	}

	// Should include escaped version in the slug
	if !strings.Contains(result, "ltdirgt") {
		t.Errorf("slug should include ltdirgt\ngot: %s", result)
	}
}

func TestStripHTMLTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no tags",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "simple tag",
			input:    "<b>bold</b>",
			expected: "bold",
		},
		{
			name:     "nested tags",
			input:    "<div><span>content</span></div>",
			expected: "content",
		},
		{
			name:     "tag with attributes",
			input:    `<a href="http://example.com">link</a>`,
			expected: "link",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only tags",
			input:    "<br/><hr/>",
			expected: "",
		},
		{
			name:     "mixed content",
			input:    "Text <b>bold</b> and <i>italic</i> text",
			expected: "Text bold and italic text",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := stripHTMLTags(tc.input)
			if result != tc.expected {
				t.Errorf("stripHTMLTags(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

package markdown

import (
	"strings"
	"testing"
)

func TestNewParser(t *testing.T) {
	p := NewParser()
	if p == nil {
		t.Fatal("NewParser() returned nil")
	}
	if p.md == nil {
		t.Fatal("Parser.md is nil")
	}
}

func TestParserParse(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:     "heading",
			input:    "# Hello World",
			contains: []string{"<h1", "Hello World", "</h1>"},
		},
		{
			name:     "paragraph",
			input:    "This is a paragraph.",
			contains: []string{"<p>", "This is a paragraph.", "</p>"},
		},
		{
			name:     "bold",
			input:    "This is **bold** text.",
			contains: []string{"<strong>", "bold", "</strong>"},
		},
		{
			name:     "italic",
			input:    "This is *italic* text.",
			contains: []string{"<em>", "italic", "</em>"},
		},
		{
			name:     "link",
			input:    "[Click here](https://example.com)",
			contains: []string{"<a", "href=\"https://example.com\"", "Click here", "</a>"},
		},
		{
			name:     "unordered list",
			input:    "- Item 1\n- Item 2\n- Item 3",
			contains: []string{"<ul>", "<li>", "Item 1", "Item 2", "Item 3", "</li>", "</ul>"},
		},
		{
			name:     "ordered list",
			input:    "1. First\n2. Second\n3. Third",
			contains: []string{"<ol>", "<li>", "First", "Second", "Third", "</li>", "</ol>"},
		},
		{
			name:     "code block",
			input:    "```go\nfunc main() {}\n```",
			contains: []string{"<pre", "<code", "func", "main"},
		},
		{
			name:     "inline code",
			input:    "Use the `fmt.Println` function.",
			contains: []string{"<code>", "fmt.Println", "</code>"},
		},
		{
			name:     "blockquote",
			input:    "> This is a quote",
			contains: []string{"<blockquote>", "This is a quote"},
		},
		{
			name:     "horizontal rule",
			input:    "Above\n\n---\n\nBelow",
			contains: []string{"<hr"},
		},
		{
			name:     "table (GFM)",
			input:    "| A | B |\n|---|---|\n| 1 | 2 |",
			contains: []string{"<table>", "<thead>", "<tbody>", "<tr>", "<th>", "<td>"},
		},
		{
			name:     "strikethrough (GFM)",
			input:    "This is ~~deleted~~ text.",
			contains: []string{"<del>", "deleted", "</del>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := p.Parse([]byte(tt.input))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			html := string(result)
			for _, expected := range tt.contains {
				if !strings.Contains(html, expected) {
					t.Errorf("Parse() result should contain %q, got: %s", expected, html)
				}
			}
		})
	}
}

func TestParserParseString(t *testing.T) {
	p := NewParser()

	input := "# Test Heading\n\nSome paragraph text."
	result, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	if !strings.Contains(result, "<h1") {
		t.Error("ParseString() should contain h1 tag")
	}
	if !strings.Contains(result, "Test Heading") {
		t.Error("ParseString() should contain heading text")
	}
	if !strings.Contains(result, "<p>") {
		t.Error("ParseString() should contain p tag")
	}
}

func TestParserSyntaxHighlighting(t *testing.T) {
	p := NewParser()

	input := "```go\npackage main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n```"
	result, err := p.Parse([]byte(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	html := string(result)
	// Check that code block is rendered with syntax highlighting
	if !strings.Contains(html, "<pre") {
		t.Error("Should contain pre tag")
	}
	if !strings.Contains(html, "<code") {
		t.Error("Should contain code tag")
	}
	// Check for highlighted tokens (chroma adds spans)
	if !strings.Contains(html, "package") {
		t.Error("Should contain 'package' keyword")
	}
}

func TestParserMultipleHeadings(t *testing.T) {
	p := NewParser()

	input := `# Heading 1
## Heading 2
### Heading 3
#### Heading 4
##### Heading 5
###### Heading 6`

	result, err := p.ParseString(input)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	for i := 1; i <= 6; i++ {
		tag := "<h" + string(rune('0'+i))
		if !strings.Contains(result, tag) {
			t.Errorf("Should contain %s tag", tag)
		}
	}
}

package markdown

import (
	"strings"
	"testing"
)

func TestWrapCodeBlocks(t *testing.T) {
	input := `<pre><code class="language-go">package main

func main() {
    fmt.Println("Hello")
}
</code></pre>`

	result := WrapCodeBlocks(input)

	if !strings.Contains(result, `class="code-block"`) {
		t.Error("should wrap code block")
	}
	if !strings.Contains(result, "copy-button") {
		t.Error("should add copy button")
	}
	if !strings.Contains(result, "Copy") {
		t.Error("should have Copy text")
	}
}

func TestWrapCodeBlocksMultiple(t *testing.T) {
	input := `
<pre><code class="language-go">go code</code></pre>
<p>Some text</p>
<pre><code class="language-python">python code</code></pre>
`

	result := WrapCodeBlocks(input)

	if strings.Count(result, "copy-button") != 2 {
		t.Errorf("should have 2 copy buttons, got %d", strings.Count(result, "copy-button"))
	}
}

func TestWrapCodeBlocksNoLanguage(t *testing.T) {
	input := `<pre><code>plain code</code></pre>`
	result := WrapCodeBlocks(input)

	if !strings.Contains(result, `class="code-block"`) {
		t.Error("should wrap code blocks without language class")
	}
}

func TestParseLineSpec(t *testing.T) {
	tests := []struct {
		spec     string
		expected []int
	}{
		{"1", []int{1}},
		{"1,3,5", []int{1, 3, 5}},
		{"1-3", []int{1, 2, 3}},
		{"1-3,5,7-9", []int{1, 2, 3, 5, 7, 8, 9}},
		{"", nil},
		{"invalid", nil},
	}

	for _, tc := range tests {
		t.Run(tc.spec, func(t *testing.T) {
			result := ParseLineSpec(tc.spec)
			if len(result.Lines) != len(tc.expected) {
				t.Errorf("ParseLineSpec(%q) = %v, want %v", tc.spec, result.Lines, tc.expected)
				return
			}
			for i, line := range result.Lines {
				if line != tc.expected[i] {
					t.Errorf("line %d: got %d, want %d", i, line, tc.expected[i])
				}
			}
		})
	}
}

func TestApplyLineHighlighting(t *testing.T) {
	code := `line 1
line 2
line 3
line 4`

	result := ApplyLineHighlighting(code, []int{2, 4})

	if !strings.Contains(result, `class="line highlight"`) {
		t.Error("should have line highlight class")
	}

	// Count highlighted lines
	count := strings.Count(result, "line highlight")
	if count != 2 {
		t.Errorf("should have 2 highlighted lines, got %d", count)
	}
}

func TestApplyLineHighlightingEmpty(t *testing.T) {
	code := "some code"
	result := ApplyLineHighlighting(code, nil)

	if result != code {
		t.Error("with no lines to highlight, should return original")
	}
}

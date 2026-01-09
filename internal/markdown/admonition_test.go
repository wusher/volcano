package markdown

import (
	"strings"
	"testing"
)

func TestProcessAdmonitions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "note admonition",
			input: `:::note
This is a note.
:::`,
			contains: []string{
				`class="admonition admonition-note"`,
				`class="admonition-title">Note</span>`,
				"This is a note.",
			},
		},
		{
			name: "tip with custom title",
			input: `:::tip Custom Title
This is a tip.
:::`,
			contains: []string{
				`class="admonition admonition-tip"`,
				`class="admonition-title">Custom Title</span>`,
			},
		},
		{
			name: "warning admonition",
			input: `:::warning
Be careful!
:::`,
			contains: []string{
				`class="admonition admonition-warning"`,
			},
		},
		{
			name: "danger admonition",
			input: `:::danger
Danger zone!
:::`,
			contains: []string{
				`class="admonition admonition-danger"`,
			},
		},
		{
			name: "info admonition",
			input: `:::info
Some info.
:::`,
			contains: []string{
				`class="admonition admonition-info"`,
			},
		},
		{
			name:     "no admonition",
			input:    "Regular markdown content",
			contains: []string{"Regular markdown content"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ProcessAdmonitions(tc.input)
			for _, expected := range tc.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("result should contain %q\ngot: %s", expected, result)
				}
			}
		})
	}
}

func TestProcessAdmonitionsUnclosed(t *testing.T) {
	input := `:::note
Unclosed admonition content`

	result := ProcessAdmonitions(input)

	// Should still render the admonition
	if !strings.Contains(result, "admonition-note") {
		t.Error("unclosed admonition should still be rendered")
	}
}

func TestProcessAdmonitionsMultiple(t *testing.T) {
	input := `:::note
First note.
:::

Some text between.

:::warning
A warning.
:::`

	result := ProcessAdmonitions(input)

	if !strings.Contains(result, "admonition-note") {
		t.Error("should contain note admonition")
	}
	if !strings.Contains(result, "admonition-warning") {
		t.Error("should contain warning admonition")
	}
	if !strings.Contains(result, "Some text between") {
		t.Error("should contain text between admonitions")
	}
}

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"<script>", "&lt;script&gt;"},
		{"a & b", "a &amp; b"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"it's", "it&#39;s"},
	}

	for _, tc := range tests {
		result := escapeHTML(tc.input)
		if result != tc.expected {
			t.Errorf("escapeHTML(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

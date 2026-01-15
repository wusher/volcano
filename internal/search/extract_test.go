package search

import (
	"testing"
)

func TestExtractHeadings(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected []HeadingEntry
	}{
		{
			name: "extracts h2-h4 with ids",
			html: `<h2 id="intro">Introduction</h2>
<p>Some text</p>
<h3 id="setup">Setup Guide</h3>
<h4 id="config">Configuration</h4>`,
			expected: []HeadingEntry{
				{Text: "Introduction", Anchor: "intro", Level: 2},
				{Text: "Setup Guide", Anchor: "setup", Level: 3},
				{Text: "Configuration", Anchor: "config", Level: 4},
			},
		},
		{
			name: "strips inner HTML tags",
			html: `<h2 id="test"><a href="#">Link</a> <code>Code</code> Text</h2>`,
			expected: []HeadingEntry{
				{Text: "Link Code Text", Anchor: "test", Level: 2},
			},
		},
		{
			name:     "returns nil for no headings",
			html:     `<p>Just a paragraph</p>`,
			expected: nil,
		},
		{
			name:     "ignores h1 headings",
			html:     `<h1 id="title">Title</h1>`,
			expected: nil,
		},
		{
			name:     "ignores headings without id",
			html:     `<h2>No ID</h2>`,
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ExtractHeadings(tc.html)
			if len(result) != len(tc.expected) {
				t.Errorf("expected %d headings, got %d", len(tc.expected), len(result))
				return
			}
			for i, h := range result {
				if h != tc.expected[i] {
					t.Errorf("heading %d: expected %+v, got %+v", i, tc.expected[i], h)
				}
			}
		})
	}
}

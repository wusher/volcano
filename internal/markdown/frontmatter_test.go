package markdown

import (
	"testing"
)

func TestStripFrontMatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no front matter",
			input:    "# Hello\n\nContent here",
			expected: "# Hello\n\nContent here",
		},
		{
			name: "simple front matter",
			input: `---
title: My Page
date: 2024-01-01
---

# Hello

Content here`,
			expected: `# Hello

Content here`,
		},
		{
			name: "front matter with empty value",
			input: `---
title: Test
tags:
---

Content`,
			expected: `Content`,
		},
		{
			name:     "only opening delimiter",
			input:    "---\ntitle: Test\nNo closing delimiter",
			expected: "---\ntitle: Test\nNo closing delimiter",
		},
		{
			name:     "delimiter not at start",
			input:    "Some text\n---\ntitle: Test\n---\nContent",
			expected: "Some text\n---\ntitle: Test\n---\nContent",
		},
		{
			name: "multiple dashes in content",
			input: `---
title: Test
---

Some content with --- dashes`,
			expected: `Some content with --- dashes`,
		},
		{
			name: "windows line endings",
			input:    "---\r\ntitle: Test\r\n---\r\n\r\nContent",
			expected: "Content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripFrontMatter([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("StripFrontMatter() = %q, want %q", string(result), tt.expected)
			}
		})
	}
}

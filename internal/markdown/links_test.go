package markdown

import (
	"strings"
	"testing"
)

func TestProcessExternalLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		siteURL  string
		contains []string
		excludes []string
	}{
		{
			name:    "external link",
			input:   `<a href="https://example.com">Example</a>`,
			siteURL: "https://mysite.com",
			contains: []string{
				`target="_blank"`,
				`rel="noopener noreferrer"`,
				`class="external-icon"`,
			},
		},
		{
			name:    "internal link",
			input:   `<a href="/about/">About</a>`,
			siteURL: "https://mysite.com",
			excludes: []string{
				`target="_blank"`,
				`external-icon`,
			},
		},
		{
			name:    "same domain link",
			input:   `<a href="https://mysite.com/page/">Page</a>`,
			siteURL: "https://mysite.com",
			excludes: []string{
				`target="_blank"`,
				`external-icon`,
			},
		},
		{
			name:    "relative link",
			input:   `<a href="../other/">Other</a>`,
			siteURL: "https://mysite.com",
			excludes: []string{
				`target="_blank"`,
			},
		},
		{
			name:    "anchor link",
			input:   `<a href="#section">Section</a>`,
			siteURL: "https://mysite.com",
			excludes: []string{
				`target="_blank"`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ProcessExternalLinks(tc.input, tc.siteURL)

			for _, expected := range tc.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("result should contain %q\ngot: %s", expected, result)
				}
			}

			for _, excluded := range tc.excludes {
				if strings.Contains(result, excluded) {
					t.Errorf("result should not contain %q\ngot: %s", excluded, result)
				}
			}
		})
	}
}

func TestProcessExternalLinksNoSiteURL(t *testing.T) {
	input := `<a href="https://example.com">Example</a>`
	result := ProcessExternalLinks(input, "")

	// Without site URL, all http/https links are external
	if !strings.Contains(result, `target="_blank"`) {
		t.Error("should mark as external when no site URL")
	}
}

func TestProcessExternalLinksMultiple(t *testing.T) {
	input := `
<p>Check out <a href="https://external.com">External</a> and <a href="/internal/">Internal</a>.</p>
`
	result := ProcessExternalLinks(input, "https://mysite.com")

	// Should only add external attributes to external link
	if strings.Count(result, `target="_blank"`) != 1 {
		t.Error("should have exactly one external link")
	}
}

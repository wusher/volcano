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

func TestPrefixInternalLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		baseURL  string
		expected string
	}{
		{
			name:     "prefix internal link with subpath",
			input:    `<a href="/guides/intro/">Intro</a>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<a href="/volcano/guides/intro/">Intro</a>`,
		},
		{
			name:     "prefix multiple internal links",
			input:    `<a href="/docs/">Docs</a> and <a href="/about/">About</a>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<a href="/volcano/docs/">Docs</a> and <a href="/volcano/about/">About</a>`,
		},
		{
			name:     "skip external links",
			input:    `<a href="https://external.com">External</a>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<a href="https://external.com">External</a>`,
		},
		{
			name:     "skip anchor links",
			input:    `<a href="#section">Section</a>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<a href="#section">Section</a>`,
		},
		{
			name:     "skip protocol-relative links",
			input:    `<a href="//cdn.example.com/file.js">CDN</a>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<a href="//cdn.example.com/file.js">CDN</a>`,
		},
		{
			name:     "no prefix when no base path",
			input:    `<a href="/docs/">Docs</a>`,
			baseURL:  "https://example.com/",
			expected: `<a href="/docs/">Docs</a>`,
		},
		{
			name:     "no prefix when empty baseURL",
			input:    `<a href="/docs/">Docs</a>`,
			baseURL:  "",
			expected: `<a href="/docs/">Docs</a>`,
		},
		{
			name:     "preserve link attributes",
			input:    `<a class="nav-link" href="/docs/" id="link1">Docs</a>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<a class="nav-link" href="/volcano/docs/" id="link1">Docs</a>`,
		},
		{
			name:     "prefix img src",
			input:    `<img src="/images/logo.png" alt="Logo">`,
			baseURL:  "https://example.com/volcano/",
			expected: `<img src="/volcano/images/logo.png" alt="Logo">`,
		},
		{
			name:     "prefix multiple src attributes",
			input:    `<img src="/img1.png"><script src="/js/app.js"></script>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<img src="/volcano/img1.png"><script src="/volcano/js/app.js"></script>`,
		},
		{
			name:     "prefix srcset",
			input:    `<img srcset="/img-1x.png 1x, /img-2x.png 2x">`,
			baseURL:  "https://example.com/volcano/",
			expected: `<img srcset="/volcano/img-1x.png 1x, /volcano/img-2x.png 2x">`,
		},
		{
			name:     "prefix video poster",
			input:    `<video poster="/thumb.jpg" src="/video.mp4"></video>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<video poster="/volcano/thumb.jpg" src="/volcano/video.mp4"></video>`,
		},
		{
			name:     "skip external images",
			input:    `<img src="https://example.com/img.png">`,
			baseURL:  "https://example.com/volcano/",
			expected: `<img src="https://example.com/img.png">`,
		},
		{
			name:     "mixed internal and external resources",
			input:    `<img src="/logo.png"><img src="https://cdn.com/icon.png"><a href="/docs/">Docs</a>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<img src="/volcano/logo.png"><img src="https://cdn.com/icon.png"><a href="/volcano/docs/">Docs</a>`,
		},
		{
			name:     "prefix data-* attributes",
			input:    `<div data-url="/api/endpoint" data-image="/img.png">Content</div>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<div data-url="/volcano/api/endpoint" data-image="/volcano/img.png">Content</div>`,
		},
		{
			name:     "srcset with single URL",
			input:    `<img srcset="/img.png">`,
			baseURL:  "https://example.com/volcano/",
			expected: `<img srcset="/volcano/img.png">`,
		},
		{
			name:     "srcset with width descriptors",
			input:    `<img srcset="/img-sm.png 400w, /img-lg.png 800w">`,
			baseURL:  "https://example.com/volcano/",
			expected: `<img srcset="/volcano/img-sm.png 400w, /volcano/img-lg.png 800w">`,
		},
		{
			name:     "skip data-* with external URLs",
			input:    `<div data-url="https://example.com/api">Content</div>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<div data-url="https://example.com/api">Content</div>`,
		},
		{
			name:     "complex mixed content",
			input:    `<a href="/page/">Link</a><img src="/img.png" srcset="/img-1x.png 1x, /img-2x.png 2x"><video poster="/poster.jpg" src="/video.mp4"></video><div data-path="/data">Text</div>`,
			baseURL:  "https://example.com/volcano/",
			expected: `<a href="/volcano/page/">Link</a><img src="/volcano/img.png" srcset="/volcano/img-1x.png 1x, /volcano/img-2x.png 2x"><video poster="/volcano/poster.jpg" src="/volcano/video.mp4"></video><div data-path="/volcano/data">Text</div>`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := PrefixInternalLinks(tc.input, tc.baseURL)
			if result != tc.expected {
				t.Errorf("PrefixInternalLinks(%q, %q)\ngot:  %q\nwant: %q", tc.input, tc.baseURL, result, tc.expected)
			}
		})
	}
}

func TestExtractBasePath(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		expected string
	}{
		{
			name:     "URL with subpath",
			baseURL:  "https://example.com/volcano/",
			expected: "/volcano",
		},
		{
			name:     "URL with nested subpath",
			baseURL:  "https://example.com/docs/v2/",
			expected: "/docs/v2",
		},
		{
			name:     "URL without subpath",
			baseURL:  "https://example.com/",
			expected: "",
		},
		{
			name:     "URL without trailing slash",
			baseURL:  "https://example.com/volcano",
			expected: "/volcano",
		},
		{
			name:     "URL root without trailing slash",
			baseURL:  "https://example.com",
			expected: "",
		},
		{
			name:     "just a path",
			baseURL:  "/volcano",
			expected: "/volcano",
		},
		{
			name:     "empty string",
			baseURL:  "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractBasePath(tc.baseURL)
			if result != tc.expected {
				t.Errorf("extractBasePath(%q) = %q, want %q", tc.baseURL, result, tc.expected)
			}
		})
	}
}

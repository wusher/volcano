package markdown

import (
	"testing"
)

func TestExtractInternalLinks(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected []string
	}{
		{
			name:     "no links",
			html:     "<p>Hello world</p>",
			expected: []string{},
		},
		{
			name:     "single internal link",
			html:     `<a href="/about/">About</a>`,
			expected: []string{"/about/"},
		},
		{
			name:     "multiple internal links",
			html:     `<a href="/about/">About</a> <a href="/contact/">Contact</a>`,
			expected: []string{"/about/", "/contact/"},
		},
		{
			name:     "skip external links",
			html:     `<a href="https://example.com">External</a> <a href="/about/">About</a>`,
			expected: []string{"/about/"},
		},
		{
			name:     "skip root link",
			html:     `<a href="/">Home</a> <a href="/about/">About</a>`,
			expected: []string{"/about/"},
		},
		{
			name:     "skip anchor-only links",
			html:     `<a href="/#section">Section</a> <a href="/about/">About</a>`,
			expected: []string{"/about/"},
		},
		{
			name:     "deduplicate links",
			html:     `<a href="/about/">About</a> <a href="/about/">About Again</a>`,
			expected: []string{"/about/"},
		},
		{
			name:     "links with anchors",
			html:     `<a href="/about/#section">About Section</a>`,
			expected: []string{"/about/#section"},
		},
		{
			name:     "nested path links",
			html:     `<a href="/guides/intro/">Intro</a> <a href="/api/endpoints/">API</a>`,
			expected: []string{"/guides/intro/", "/api/endpoints/"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractInternalLinks(tt.html)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d links, got %d: %v", len(tt.expected), len(result), result)
				return
			}
			for i, link := range result {
				if link != tt.expected[i] {
					t.Errorf("link %d: expected %q, got %q", i, tt.expected[i], link)
				}
			}
		})
	}
}

func TestValidateLinks(t *testing.T) {
	validURLs := map[string]bool{
		"/":               true,
		"/about/":         true,
		"/guides/intro/":  true,
		"/api/endpoints/": true,
	}

	tests := []struct {
		name           string
		html           string
		sourcePage     string
		expectedBroken int
	}{
		{
			name:           "all links valid",
			html:           `<a href="/about/">About</a> <a href="/guides/intro/">Intro</a>`,
			sourcePage:     "/",
			expectedBroken: 0,
		},
		{
			name:           "one broken link",
			html:           `<a href="/about/">About</a> <a href="/nonexistent/">Missing</a>`,
			sourcePage:     "/",
			expectedBroken: 1,
		},
		{
			name:           "multiple broken links",
			html:           `<a href="/broken1/">B1</a> <a href="/broken2/">B2</a>`,
			sourcePage:     "/about/",
			expectedBroken: 2,
		},
		{
			name:           "link with anchor to valid page",
			html:           `<a href="/about/#section">About Section</a>`,
			sourcePage:     "/",
			expectedBroken: 0,
		},
		{
			name:           "link without trailing slash to valid page",
			html:           `<a href="/about">About</a>`,
			sourcePage:     "/",
			expectedBroken: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			broken := ValidateLinks(tt.html, tt.sourcePage, validURLs)
			if len(broken) != tt.expectedBroken {
				t.Errorf("expected %d broken links, got %d: %v", tt.expectedBroken, len(broken), broken)
			}
		})
	}
}

func TestNormalizeLink(t *testing.T) {
	tests := []struct {
		link     string
		expected string
	}{
		{"/about/", "/about/"},
		{"/about", "/about/"},
		{"/about/#section", "/about/"},
		{"/guides/intro/", "/guides/intro/"},
		{"/file.html", "/file.html"},
		{"/", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.link, func(t *testing.T) {
			result := normalizeLink(tt.link)
			if result != tt.expected {
				t.Errorf("normalizeLink(%q) = %q, expected %q", tt.link, result, tt.expected)
			}
		})
	}
}

func TestBrokenLinkStruct(t *testing.T) {
	bl := BrokenLink{
		SourcePage: "/about/",
		LinkURL:    "/nonexistent/",
	}

	if bl.SourcePage != "/about/" {
		t.Errorf("SourcePage = %q, expected %q", bl.SourcePage, "/about/")
	}
	if bl.LinkURL != "/nonexistent/" {
		t.Errorf("LinkURL = %q, expected %q", bl.LinkURL, "/nonexistent/")
	}
}

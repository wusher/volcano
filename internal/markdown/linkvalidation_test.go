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

func TestExtractLinkInfoFromMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		wantURLs []string
	}{
		{
			name:     "empty",
			markdown: "",
			wantURLs: []string{},
		},
		{
			name:     "wikilink simple",
			markdown: "Check out [[about]] for more info",
			wantURLs: []string{"/about/"},
		},
		{
			name:     "wikilink with pipe",
			markdown: "See [[about|About Page]] for details",
			wantURLs: []string{"/about/"},
		},
		{
			name:     "wikilink with path",
			markdown: "Read [[guides/intro]]",
			wantURLs: []string{"/guides/intro/"},
		},
		{
			name:     "markdown link",
			markdown: "Check [about](/about/) page",
			wantURLs: []string{"/about/"},
		},
		{
			name:     "multiple links",
			markdown: "See [[about]] and [[contact]]",
			wantURLs: []string{"/about/", "/contact/"},
		},
		{
			name:     "wikilink with anchor",
			markdown: "See [[page#section]]",
			wantURLs: []string{"/page/#section"},
		},
		{
			name:     "wikilink index/readme",
			markdown: "Check [[folder/index]]",
			wantURLs: []string{"/folder/index/", "/folder/"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := extractLinkInfoFromMarkdown(tt.markdown)
			if len(info) != len(tt.wantURLs) {
				t.Errorf("got %d URLs, want %d", len(info), len(tt.wantURLs))
				return
			}
			for _, url := range tt.wantURLs {
				if _, found := info[url]; !found {
					t.Errorf("URL %q not found in extracted info", url)
				}
			}
		})
	}
}

func TestExtractLinkInfoFromMarkdown_LineNumbers(t *testing.T) {
	markdown := `# Title

Check out [[about]] for more info.
Another link to [[contact]].`

	info := extractLinkInfoFromMarkdown(markdown)

	// Check about link
	aboutInfo, found := info["/about/"]
	if !found {
		t.Fatal("about link not found")
	}
	if aboutInfo.LineNumber != 3 {
		t.Errorf("about link line number = %d, want 3", aboutInfo.LineNumber)
	}
	if aboutInfo.OriginalSyntax != "[[about]]" {
		t.Errorf("about link syntax = %q, want [[about]]", aboutInfo.OriginalSyntax)
	}

	// Check contact link
	contactInfo, found := info["/contact/"]
	if !found {
		t.Fatal("contact link not found")
	}
	if contactInfo.LineNumber != 4 {
		t.Errorf("contact link line number = %d, want 4", contactInfo.LineNumber)
	}
}

func TestGetLastPathSegment(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"about", "about"},
		{"guides/intro", "intro"},
		{"folder/subfolder/page", "page"},
		{"", ""},
		{"path/", "path"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := getLastPathSegment(tt.path)
			if got != tt.want {
				t.Errorf("getLastPathSegment(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestGenerateWikilinkURLVariations(t *testing.T) {
	tests := []struct {
		target string
		want   []string
	}{
		{
			target: "about",
			want:   []string{"/about/"},
		},
		{
			target: "guides/intro",
			want:   []string{"/guides/intro/"},
		},
		{
			target: "folder/index",
			want:   []string{"/folder/index/", "/folder/"},
		},
		{
			target: "folder/readme",
			want:   []string{"/folder/readme/", "/folder/"},
		},
		{
			target: "page#section",
			want:   []string{"/page/#section"},
		},
		{
			target: "#anchor",
			want:   []string{"#anchor"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			got := generateWikilinkURLVariations(tt.target)
			if len(got) != len(tt.want) {
				t.Errorf("got %d variations, want %d: %v", len(got), len(tt.want), got)
				return
			}
			for i, url := range got {
				if url != tt.want[i] {
					t.Errorf("variation %d: got %q, want %q", i, url, tt.want[i])
				}
			}
		})
	}
}

func TestValidateLinksWithSource(t *testing.T) {
	validURLs := map[string]bool{
		"/":      true,
		"/about/": true,
	}

	markdown := "See [[contact]] for info"
	html := `<a href="/contact/">Contact</a>`

	broken := ValidateLinksWithSource(html, "/test/", "/test.md", markdown, validURLs)

	if len(broken) != 1 {
		t.Fatalf("expected 1 broken link, got %d", len(broken))
	}

	bl := broken[0]
	if bl.SourceFile != "/test.md" {
		t.Errorf("SourceFile = %q, want /test.md", bl.SourceFile)
	}
	if bl.LineNumber != 1 {
		t.Errorf("LineNumber = %d, want 1", bl.LineNumber)
	}
	if bl.OriginalSyntax != "[[contact]]" {
		t.Errorf("OriginalSyntax = %q, want [[contact]]", bl.OriginalSyntax)
	}
	if bl.LinkURL != "/contact/" {
		t.Errorf("LinkURL = %q, want /contact/", bl.LinkURL)
	}
}

func TestFindSimilarURLs(t *testing.T) {
	validURLs := map[string]bool{
		"/":                  true,
		"/about/":            true,
		"/about-us/":         true,
		"/contact/":          true,
		"/guides/intro/":     true,
		"/guides/advanced/":  true,
	}

	tests := []struct {
		brokenURL       string
		wantSuggestions int
	}{
		{"/about-page/", 1},        // Should suggest /about/ or /about-us/
		{"/guides/", 2},            // Should suggest /guides/intro/ and /guides/advanced/
		{"/xyz/", 0},               // No similar URLs
		{"/intro/", 1},             // Should suggest /guides/intro/
	}

	for _, tt := range tests {
		t.Run(tt.brokenURL, func(t *testing.T) {
			suggestions := findSimilarURLs(tt.brokenURL, validURLs)
			if len(suggestions) == 0 && tt.wantSuggestions > 0 {
				t.Errorf("findSimilarURLs(%q) returned no suggestions, want at least %d", tt.brokenURL, tt.wantSuggestions)
			}
		})
	}
}

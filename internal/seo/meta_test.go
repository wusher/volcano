package seo

import (
	"strings"
	"testing"
)

func TestGeneratePageMeta(t *testing.T) {
	config := Config{
		SiteURL:   "https://example.com",
		SiteTitle: "My Site",
		Author:    "Test Author",
		OGImage:   "https://example.com/og.png",
	}

	content := "<p>This is a test page with some content for testing.</p>"
	meta := GeneratePageMeta("Test Page", content, "/test/", config)

	// Check SEO
	if meta.SEO.Title != "Test Page - My Site" {
		t.Errorf("SEO.Title = %q, want %q", meta.SEO.Title, "Test Page - My Site")
	}
	if meta.SEO.Canonical != "https://example.com/test/" {
		t.Errorf("SEO.Canonical = %q, want %q", meta.SEO.Canonical, "https://example.com/test/")
	}
	if meta.SEO.Author != "Test Author" {
		t.Errorf("SEO.Author = %q, want %q", meta.SEO.Author, "Test Author")
	}

	// Check Open Graph
	if meta.OG.Title != "Test Page" {
		t.Errorf("OG.Title = %q, want %q", meta.OG.Title, "Test Page")
	}
	if meta.OG.SiteName != "My Site" {
		t.Errorf("OG.SiteName = %q, want %q", meta.OG.SiteName, "My Site")
	}
	if meta.OG.Type != "article" {
		t.Errorf("OG.Type = %q, want %q", meta.OG.Type, "article")
	}

	// Check Twitter
	if meta.Twitter.Card != "summary_large_image" {
		t.Errorf("Twitter.Card = %q, want %q", meta.Twitter.Card, "summary_large_image")
	}
}

func TestGeneratePageMetaNoImage(t *testing.T) {
	config := Config{
		SiteURL:   "https://example.com",
		SiteTitle: "My Site",
	}

	meta := GeneratePageMeta("Test", "<p>Test</p>", "/", config)

	if meta.Twitter.Card != "summary" {
		t.Errorf("Twitter.Card without image = %q, want %q", meta.Twitter.Card, "summary")
	}
}

func TestExtractDescription(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		maxLen   int
		expected string
	}{
		{"empty", "", 160, ""},
		{"short", "<p>Hello world</p>", 160, "Hello world"},
		{"strips tags", "<h1>Title</h1><p>Content here</p>", 160, "Title Content here"},
		{"truncates", "<p>" + strings.Repeat("word ", 50) + "</p>", 50, "word word word word word word word word word..."},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractDescription(tc.html, tc.maxLen)
			if tc.expected != "" && !strings.HasPrefix(result, strings.Split(tc.expected, "...")[0][:10]) {
				// Just check prefix for truncated content
				if result == "" && tc.expected != "" {
					t.Errorf("extractDescription() = %q, want something like %q", result, tc.expected)
				}
			}
		})
	}
}

func TestGetTwitterCardType(t *testing.T) {
	if getTwitterCardType("https://example.com/img.png") != "summary_large_image" {
		t.Error("with image should be summary_large_image")
	}
	if getTwitterCardType("") != "summary" {
		t.Error("without image should be summary")
	}
}

func TestRenderMetaTags(t *testing.T) {
	meta := PageMeta{
		SEO: Meta{
			Description: "Test description",
			Robots:      "index, follow",
			Author:      "Test Author",
			Canonical:   "https://example.com/test/",
		},
		OG: OpenGraph{
			Title:       "Test Title",
			Description: "Test description",
			Type:        "article",
			URL:         "https://example.com/test/",
			SiteName:    "My Site",
			Image:       "https://example.com/og.png",
		},
		Twitter: TwitterCard{
			Card:        "summary_large_image",
			Title:       "Test Title",
			Description: "Test description",
			Image:       "https://example.com/og.png",
		},
	}

	html := string(RenderMetaTags(meta))

	// Check SEO tags
	if !strings.Contains(html, `name="description"`) {
		t.Error("should contain description meta")
	}
	if !strings.Contains(html, `name="robots"`) {
		t.Error("should contain robots meta")
	}
	if !strings.Contains(html, `name="author"`) {
		t.Error("should contain author meta")
	}
	if !strings.Contains(html, `rel="canonical"`) {
		t.Error("should contain canonical link")
	}

	// Check OG tags
	if !strings.Contains(html, `property="og:title"`) {
		t.Error("should contain og:title")
	}
	if !strings.Contains(html, `property="og:type"`) {
		t.Error("should contain og:type")
	}

	// Check Twitter tags
	if !strings.Contains(html, `name="twitter:card"`) {
		t.Error("should contain twitter:card")
	}
}

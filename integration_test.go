package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationGenerateFromExample(t *testing.T) {
	// Skip if example folder doesn't exist
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Create a temporary output directory
	outputDir := t.TempDir()

	// Run the generator
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--title=Integration Test", "./example"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Generation failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Verify output contains expected messages
	output := stdout.String()
	if !strings.Contains(output, "Generating site") {
		t.Error("Output should contain 'Generating site'")
	}
	if !strings.Contains(output, "Generated") {
		t.Error("Output should contain 'Generated'")
	}

	// Verify expected files exist
	expectedFiles := []string{
		"index.html",
		"404.html",
		"getting-started/index.html",
		"faq/index.html",
		"guides/index.html",
		"guides/installation/index.html",
		"guides/configuration/index.html",
		"api/index.html",
		"api/endpoints/index.html",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(outputDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", file)
		}
	}
}

func TestIntegrationHTMLContent(t *testing.T) {
	// Skip if example folder doesn't exist
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Create a temporary output directory
	outputDir := t.TempDir()

	// Run the generator
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--title=Content Test", "./example"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Generation failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Read the generated index.html
	indexPath := filepath.Join(outputDir, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}

	html := string(content)

	// Verify HTML structure
	checks := []struct {
		name    string
		content string
	}{
		{"doctype", "<!DOCTYPE html>"},
		{"title", "<title>"},
		{"site title", "Content Test"},
		{"navigation", "class=\"tree-nav\""},
		{"sidebar", "class=\"sidebar\""},
		{"content", "class=\"content\""},
		{"theme toggle", "class=\"theme-toggle\""},
		{"mobile menu", "class=\"mobile-menu-btn\""},
		{"page content", "Welcome to Volcano"},
	}

	for _, check := range checks {
		if !strings.Contains(html, check.content) {
			t.Errorf("index.html should contain %s (%q)", check.name, check.content)
		}
	}
}

func TestIntegrationSyntaxHighlighting(t *testing.T) {
	// Skip if example folder doesn't exist
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Create a temporary output directory
	outputDir := t.TempDir()

	// Run the generator
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Generation failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Read a page that has code blocks
	guidesPath := filepath.Join(outputDir, "guides/index.html")
	content, err := os.ReadFile(guidesPath)
	if err != nil {
		t.Fatalf("Failed to read guides/index.html: %v", err)
	}

	html := string(content)

	// Verify code blocks have syntax highlighting classes
	// Chroma adds class="highlight" or class="chroma"
	if !strings.Contains(html, "highlight") && !strings.Contains(html, "chroma") {
		t.Error("Code blocks should have syntax highlighting")
	}
}

func TestIntegrationNavigation(t *testing.T) {
	// Skip if example folder doesn't exist
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Create a temporary output directory
	outputDir := t.TempDir()

	// Run the generator
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Generation failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Read index.html
	indexPath := filepath.Join(outputDir, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}

	html := string(content)

	// Verify navigation contains links to pages
	navLinks := []string{
		"Getting Started",
		"Guides",
		"Api", // "api" folder becomes "Api" per CleanLabel
	}

	for _, link := range navLinks {
		if !strings.Contains(html, link) {
			t.Errorf("Navigation should contain link to %q", link)
		}
	}
}

func TestIntegration404Page(t *testing.T) {
	// Skip if example folder doesn't exist
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Create a temporary output directory
	outputDir := t.TempDir()

	// Run the generator
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Generation failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Read 404.html
	notFoundPath := filepath.Join(outputDir, "404.html")
	content, err := os.ReadFile(notFoundPath)
	if err != nil {
		t.Fatalf("Failed to read 404.html: %v", err)
	}

	html := string(content)

	// Verify 404 page content
	if !strings.Contains(html, "404") {
		t.Error("404 page should contain '404'")
	}
	if !strings.Contains(html, "Page Not Found") {
		t.Error("404 page should contain 'Page Not Found'")
	}
	if !strings.Contains(html, "Return to home") {
		t.Error("404 page should contain link to home")
	}
}

// Story 14: Table of Contents Component
func TestIntegrationStory14_TOC(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	// getting-started has multiple headings, should have TOC
	content, _ := os.ReadFile(filepath.Join(outputDir, "getting-started/index.html"))
	html := string(content)

	// Story 14 acceptance: TOC sidebar present
	if !strings.Contains(html, "toc-sidebar") {
		t.Error("Story 14: should have TOC sidebar")
	}
	// Story 14 acceptance: TOC scroll spy script
	if !strings.Contains(html, "IntersectionObserver") {
		t.Error("Story 14: should have scroll spy using IntersectionObserver")
	}
}

// Story 15: Breadcrumb Navigation
func TestIntegrationStory15_Breadcrumbs(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--title=Test Site", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "guides/configuration/index.html"))
	html := string(content)

	// Story 15 acceptance: breadcrumb navigation
	if !strings.Contains(html, "class=\"breadcrumbs\"") {
		t.Error("Story 15: should have breadcrumbs")
	}
	// Story 15 acceptance: schema.org structured data
	if !strings.Contains(html, "schema.org/BreadcrumbList") {
		t.Error("Story 15: should have schema.org BreadcrumbList")
	}
}

// Story 16: Previous/Next Page Navigation
func TestIntegrationStory16_PageNavigation(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--page-nav", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "getting-started/index.html"))
	html := string(content)

	// Story 16 acceptance: prev/next navigation
	if !strings.Contains(html, "class=\"page-nav\"") {
		t.Error("Story 16: should have page navigation")
	}
	if !strings.Contains(html, "page-nav-prev") || !strings.Contains(html, "page-nav-next") {
		t.Error("Story 16: should have prev/next links")
	}
}

// Story 17: Heading Anchor Links
func TestIntegrationStory17_HeadingAnchors(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "getting-started/index.html"))
	html := string(content)

	// Story 17 acceptance: heading anchors
	if !strings.Contains(html, "heading-anchor") {
		t.Error("Story 17: should have heading anchor links")
	}
	if !strings.Contains(html, "id=\"installation\"") {
		t.Error("Story 17: headings should have ID attributes")
	}
}

// Story 18: External Link Indicators
func TestIntegrationStory18_ExternalLinks(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "getting-started/index.html"))
	html := string(content)

	// Story 18 acceptance: external link icon styling
	if !strings.Contains(html, "external-icon") {
		t.Error("Story 18: should have external link icon class in CSS")
	}
}

// Story 19: Code Block Copy Button
func TestIntegrationStory19_CopyButton(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "getting-started/index.html"))
	html := string(content)

	// Story 19 acceptance: copy button
	if !strings.Contains(html, "copy-button") {
		t.Error("Story 19: should have copy button")
	}
	if !strings.Contains(html, "navigator.clipboard") {
		t.Error("Story 19: should have clipboard API usage")
	}
}

// Story 20: Keyboard Navigation Shortcuts
func TestIntegrationStory20_KeyboardShortcuts(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 20 acceptance: keyboard shortcuts modal
	if !strings.Contains(html, "shortcuts-modal") {
		t.Error("Story 20: should have shortcuts modal")
	}
	// Story 20 acceptance: shortcut keys defined
	if !strings.Contains(html, "keydown") {
		t.Error("Story 20: should have keydown event listener")
	}
}

// Story 21: Print Stylesheet
func TestIntegrationStory21_PrintStyles(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 21 acceptance: print stylesheet
	if !strings.Contains(html, "@media print") {
		t.Error("Story 21: should have print media query")
	}
	// Story 21 acceptance: sidebar hidden in print (CSS is minified)
	if !strings.Contains(html, "display:none!important") {
		t.Error("Story 21: should hide navigation in print")
	}
}

// Story 22: Reading Time Indicator
func TestIntegrationStory22_ReadingTime(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "getting-started/index.html"))
	html := string(content)

	// Story 22 acceptance: reading time
	if !strings.Contains(html, "reading-time") {
		t.Error("Story 22: should have reading time")
	}
	if !strings.Contains(html, "min read") {
		t.Error("Story 22: should show 'min read'")
	}
}

// Story 23: Last Modified Display
func TestIntegrationStory23_LastModified(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	// Use --last-modified flag
	exitCode := Run([]string{"-o", outputDir, "--last-modified", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "getting-started/index.html"))
	html := string(content)

	// Story 23 acceptance: last modified
	if !strings.Contains(html, "last-modified") {
		t.Error("Story 23: should have last-modified class")
	}
	if !strings.Contains(html, "Updated") {
		t.Error("Story 23: should show 'Updated' label")
	}
}

// Story 24: Scroll Progress Indicator
func TestIntegrationStory24_ScrollProgress(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 24 acceptance: scroll progress
	if !strings.Contains(html, "scroll-progress") {
		t.Error("Story 24: should have scroll progress element")
	}
	if !strings.Contains(html, "scroll-progress-bar") {
		t.Error("Story 24: should have scroll progress bar")
	}
}

// Story 25: Back to Top Button
func TestIntegrationStory25_BackToTop(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 25 acceptance: back to top
	if !strings.Contains(html, "back-to-top") {
		t.Error("Story 25: should have back-to-top button")
	}
	if !strings.Contains(html, "scrollTo") {
		t.Error("Story 25: should have scrollTo functionality")
	}
}

// Story 26: SEO Meta Tags Generation
func TestIntegrationStory26_SEOMeta(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--url=https://example.com", "--author=Test Author", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 26 acceptance: meta description
	if !strings.Contains(html, "name=\"description\"") {
		t.Error("Story 26: should have meta description")
	}
	// Story 26 acceptance: meta robots
	if !strings.Contains(html, "name=\"robots\"") {
		t.Error("Story 26: should have meta robots")
	}
}

// Story 27: Open Graph Support
func TestIntegrationStory27_OpenGraph(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--url=https://example.com", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 27 acceptance: Open Graph tags
	if !strings.Contains(html, "og:title") {
		t.Error("Story 27: should have og:title")
	}
	if !strings.Contains(html, "og:description") {
		t.Error("Story 27: should have og:description")
	}
	if !strings.Contains(html, "og:type") {
		t.Error("Story 27: should have og:type")
	}
}

// Story 28: Custom Favicon Support
func TestIntegrationStory28_Favicon(t *testing.T) {
	// Create a temp favicon
	outputDir := t.TempDir()
	faviconDir := t.TempDir()
	faviconPath := filepath.Join(faviconDir, "favicon.ico")
	if err := os.WriteFile(faviconPath, []byte{0, 0, 1, 0}, 0644); err != nil { // Minimal ICO header
		t.Fatalf("failed to write favicon: %v", err)
	}

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--favicon=" + faviconPath, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 28 acceptance: favicon link
	if !strings.Contains(html, "rel=\"icon\"") {
		t.Error("Story 28: should have favicon link tag")
	}
}

// Story 29: Admonition/Callout Blocks
func TestIntegrationStory29_Admonitions(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 29 acceptance: admonition styles
	if !strings.Contains(html, ".admonition") {
		t.Error("Story 29: should have admonition styles")
	}
	if !strings.Contains(html, ".admonition-note") {
		t.Error("Story 29: should have admonition-note style")
	}
}

// Story 30: Code Line Highlighting
func TestIntegrationStory30_LineHighlighting(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 30 acceptance: line highlight styles
	if !strings.Contains(html, ".line.highlight") {
		t.Error("Story 30: should have line highlight styles")
	}
}

// Story 31: Smooth Scroll Behavior
func TestIntegrationStory31_SmoothScroll(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 31 acceptance: smooth scroll (CSS is minified)
	if !strings.Contains(html, "scroll-behavior:smooth") {
		t.Error("Story 31: should have smooth scroll behavior")
	}
	// Story 31 acceptance: prefers-reduced-motion
	if !strings.Contains(html, "prefers-reduced-motion") {
		t.Error("Story 31: should respect prefers-reduced-motion")
	}
}

// Story 32: Clickable Folder Navigation
func TestIntegrationStory32_ClickableFolders(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 32 acceptance: folder links (folders with index)
	if !strings.Contains(html, "folder-link") {
		t.Error("Story 32: should have folder-link class for clickable folders")
	}
}

// Story 33: Navigation Tree Search
func TestIntegrationStory33_NavSearch(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 33 acceptance: navigation search
	if !strings.Contains(html, "nav-search") {
		t.Error("Story 33: should have nav-search")
	}
	if !strings.Contains(html, "nav-search-input") {
		t.Error("Story 33: should have search input")
	}
	if !strings.Contains(html, "data-search-text") {
		t.Error("Story 33: should have data-search-text attributes")
	}
}

// Story 34: Top Navigation Bar
func TestIntegrationStory34_TopNav(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--top-nav", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 34 acceptance: top nav bar
	if !strings.Contains(html, "class=\"top-nav\"") {
		t.Error("Story 34: should have top-nav class")
	}
	if !strings.Contains(html, "top-nav-list") {
		t.Error("Story 34: should have top-nav-list")
	}
	if !strings.Contains(html, "has-top-nav") {
		t.Error("Story 34: body should have has-top-nav class")
	}
}

// Story 35: Auto-Generated Folder Index
func TestIntegrationStory35_AutoIndex(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	// Posts folder has no index.md, should have auto-generated index
	postsIndex := filepath.Join(outputDir, "posts/index.html")
	if _, err := os.Stat(postsIndex); os.IsNotExist(err) {
		t.Error("Story 35: should have auto-generated posts index")
	}

	content, _ := os.ReadFile(postsIndex)
	html := string(content)

	// Story 35 acceptance: auto-index lists children
	if !strings.Contains(html, "Features") {
		t.Error("Story 35: auto-index should list child pages")
	}
}

// Story 36: H1-Based Tree Labels
func TestIntegrationStory36_H1Labels(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Story 36 acceptance: H1 title used in navigation
	// faq.md has H1 "Frequently Asked Questions"
	if !strings.Contains(html, "Frequently Asked Questions") {
		t.Error("Story 36: navigation should use H1 title")
	}
}

// Story 37: Filename Date/Number Prefixes
func TestIntegrationStory37_FilenamePrefix(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	// Posts have date prefixes like 2024-01-01-features.md
	// Should generate clean URLs without the date prefix
	featuresPath := filepath.Join(outputDir, "posts/features/index.html")
	if _, err := os.Stat(featuresPath); os.IsNotExist(err) {
		t.Error("Story 37: should strip date prefix from URL (posts/features/)")
	}

	// Also check that the navigation doesn't show the date prefix
	content, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	html := string(content)

	// Should show "Features" not "2024-01-01-features"
	if strings.Contains(html, "2024-01-01") {
		t.Error("Story 37: navigation should not show date prefix")
	}
}

// Test: Base URL Prefixing for Subpath Deployment
func TestIntegrationBaseURLPrefixing(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	// Use --url flag with a subpath to test URL prefixing
	exitCode := Run([]string{"-o", outputDir, "--url=https://wusher.github.io/volcano/", "--title=Volcano Docs", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	// Read index.html to verify URL prefixing
	content, err := os.ReadFile(filepath.Join(outputDir, "index.html"))
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}
	html := string(content)

	// Verify navigation links are prefixed with /volcano/
	if !strings.Contains(html, "href=\"/volcano/getting-started/\"") {
		t.Error("Base URL: navigation links should be prefixed with /volcano/")
	}

	// Verify site title link is prefixed
	if !strings.Contains(html, "href=\"/volcano/\"") {
		t.Error("Base URL: site title link should be prefixed with /volcano/")
	}

	// Verify keyboard shortcut home URL is set correctly
	// Note: Go templates escape forward slashes in JavaScript contexts
	if !strings.Contains(html, "const baseURL = '\\/volcano'") && !strings.Contains(html, "const baseURL = '/volcano'") {
		t.Error("Base URL: JavaScript baseURL should be set to /volcano")
	}

	// Read 404.html to verify base URL in error page
	content404, err := os.ReadFile(filepath.Join(outputDir, "404.html"))
	if err != nil {
		t.Fatalf("Failed to read 404.html: %v", err)
	}
	html404 := string(content404)

	if !strings.Contains(html404, "href=\"/volcano/\"") {
		t.Error("Base URL: 404 page home link should be prefixed with /volcano/")
	}
}

// Test: Base URL Prefixing with Breadcrumbs
func TestIntegrationBaseURLBreadcrumbs(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--url=https://wusher.github.io/volcano/", "--title=Volcano Docs", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	// Read a nested page to check breadcrumbs
	content, err := os.ReadFile(filepath.Join(outputDir, "guides/configuration/index.html"))
	if err != nil {
		t.Fatalf("Failed to read guides/configuration/index.html: %v", err)
	}
	html := string(content)

	// Verify breadcrumb home link is prefixed
	if !strings.Contains(html, "href=\"/volcano/\"") {
		t.Error("Base URL: breadcrumb home link should be prefixed with /volcano/")
	}

	// Verify breadcrumb intermediate links are prefixed
	if !strings.Contains(html, "href=\"/volcano/guides/\"") {
		t.Error("Base URL: breadcrumb guides link should be prefixed with /volcano/")
	}
}

// Test: Base URL Prefixing with Page Navigation
func TestIntegrationBaseURLPageNav(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--url=https://wusher.github.io/volcano/", "--page-nav", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	// Read a page that has prev/next navigation
	content, err := os.ReadFile(filepath.Join(outputDir, "getting-started/index.html"))
	if err != nil {
		t.Fatalf("Failed to read getting-started/index.html: %v", err)
	}
	html := string(content)

	// Verify page navigation links contain the base URL prefix
	// The prev/next links should be prefixed with /volcano/
	if !strings.Contains(html, "class=\"page-nav\"") {
		t.Error("Base URL: should have page navigation")
	}

	// Check that at least some navigation links are prefixed
	if !strings.Contains(html, "/volcano/") {
		t.Error("Base URL: page navigation links should be prefixed with /volcano/")
	}
}

// Test: Base URL Prefixing with Auto-Index
func TestIntegrationBaseURLAutoIndex(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--url=https://wusher.github.io/volcano/", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	// Posts folder has no index.md, should have auto-generated index
	content, err := os.ReadFile(filepath.Join(outputDir, "posts/index.html"))
	if err != nil {
		t.Fatalf("Failed to read posts/index.html: %v", err)
	}
	html := string(content)

	// Verify auto-index links are prefixed
	if !strings.Contains(html, "href=\"/volcano/posts/features/\"") {
		t.Error("Base URL: auto-index links should be prefixed with /volcano/")
	}
}

// Test: No Base URL Prefixing when URL has no subpath
func TestIntegrationBaseURLNoSubpath(t *testing.T) {
	outputDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	// Use --url flag without a subpath
	exitCode := Run([]string{"-o", outputDir, "--url=https://example.com/", "--title=Test Site", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	// Read index.html to verify no unnecessary prefixing
	content, err := os.ReadFile(filepath.Join(outputDir, "index.html"))
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}
	html := string(content)

	// Verify navigation links start with just /
	if !strings.Contains(html, "href=\"/getting-started/\"") {
		t.Error("No subpath: navigation links should start with just /")
	}

	// Verify site title link is just /
	if !strings.Contains(html, "href=\"/\"") {
		t.Error("No subpath: site title link should be /")
	}

	// Verify baseURL is empty in JavaScript
	// Note: empty baseURL should produce const baseURL = '';
	if !strings.Contains(html, "const baseURL = '';") {
		t.Error("No subpath: JavaScript baseURL should be empty")
	}
}

// Test: Base URL Prefixing for Content Links (Wiki Links)
func TestIntegrationBaseURLContentLinks(t *testing.T) {
	// Create a temp directory with markdown files that have wiki links
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// Create a markdown file with wiki links
	pageContent := `# Test Page

This page has [[other-page|a wiki link]] and a regular [markdown link](/docs/page/).

![Image](/images/logo.png)
`
	if err := os.WriteFile(filepath.Join(inputDir, "test.md"), []byte(pageContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create the other-page so links are valid
	otherContent := `# Other Page

Content here.
`
	if err := os.WriteFile(filepath.Join(inputDir, "other-page.md"), []byte(otherContent), 0644); err != nil {
		t.Fatalf("Failed to write other-page file: %v", err)
	}

	// Create a docs folder with a page
	docsDir := filepath.Join(inputDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("Failed to create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "page.md"), []byte("# Docs Page\n"), 0644); err != nil {
		t.Fatalf("Failed to write docs/page file: %v", err)
	}

	// Create an images directory (even if empty, for link validation)
	imagesDir := filepath.Join(outputDir, "images")
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		t.Fatalf("Failed to create images dir: %v", err)
	}

	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--url=https://wusher.github.io/volcano/", "--title=Test Site", inputDir}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("Generation failed: %s", stderr.String())
	}

	// Read the generated HTML
	content, err := os.ReadFile(filepath.Join(outputDir, "test", "index.html"))
	if err != nil {
		t.Fatalf("Failed to read test/index.html: %v", err)
	}
	html := string(content)

	// Verify wiki link is prefixed with /volcano/
	if !strings.Contains(html, `href="/volcano/other-page/"`) {
		t.Error("Content links: wiki link should be prefixed with /volcano/")
	}

	// Verify markdown link is prefixed with /volcano/
	if !strings.Contains(html, `href="/volcano/docs/page/"`) {
		t.Error("Content links: markdown link should be prefixed with /volcano/")
	}

	// Verify image src is prefixed with /volcano/
	if !strings.Contains(html, `src="/volcano/images/logo.png"`) {
		t.Error("Content resources: image src should be prefixed with /volcano/")
	}
}
